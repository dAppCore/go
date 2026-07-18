// SPDX-License-Identifier: EUPL-1.2

// Contracts, options, and type definitions for the Core framework.

package core

// Message is the type for IPC broadcasts (fire-and-forget).
//
//	c := core.New()
//	var msg core.Message = core.ActionTaskStarted{TaskIdentifier: "task-42", Action: "agent.run"}
//	c.ACTION(msg)
type Message any

// Query is the type for read-only IPC requests.
//
//	c := core.New()
//	var query core.Query = core.NewOptions(core.Option{Key: "name", Value: "agent"})
//	_ = c.Query(query)
type Query any

// QueryHandler handles Query requests. Returns Result{Value, OK}.
//
//	handler := func(c *core.Core, q core.Query) core.Result {
//	    opts := q.(core.Options)
//	    return core.Result{Value: opts.String("name"), OK: true}
//	}
//	core.New().RegisterQuery(core.QueryHandler(handler))
type QueryHandler func(*Core, Query) Result

// Startable is implemented by services that need startup initialisation.
//
//	func (s *MyService) OnStartup(ctx Context) core.Result {
//	    return core.Result{OK: true}
//	}
type Startable interface {
	OnStartup(ctx Context) Result
}

// Stoppable is implemented by services that need shutdown cleanup.
//
//	func (s *MyService) OnShutdown(ctx Context) core.Result {
//	    return core.Result{OK: true}
//	}
type Stoppable interface {
	OnShutdown(ctx Context) Result
}

// Reloadable is implemented by services that support configuration
// reload. ServiceReload is the top-level runner; trigger it from a
// signal handler, a config watcher, or an admin action.
//
//	func (s *MyService) OnReload(ctx Context) core.Result {
//	    return s.reconnect(ctx)
//	}
type Reloadable interface {
	OnReload(ctx Context) Result
}

// --- Action Messages ---

// ActionServiceStartup is broadcast when a Core service finishes startup.
//
//	c := core.New()
//	c.Action("service.startup", func(ctx Context, opts core.Options) core.Result {
//	    name := opts.String("name")
//	    core.Println(core.Sprintf("service %s started", name))
//	    return core.Result{OK: true}
//	})
type ActionServiceStartup struct{}

// ActionServiceReload is broadcast when ServiceReload completes.
//
//	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
//	    if _, ok := msg.(core.ActionServiceReload); ok {
//	        core.Info("configuration reloaded")
//	    }
//	    return core.Result{OK: true}
//	})
type ActionServiceReload struct{}

// ActionServiceShutdown is broadcast when Core begins service shutdown.
//
//	c := core.New()
//	c.Action("service.shutdown", func(ctx Context, opts core.Options) core.Result {
//	    name := opts.String("name")
//	    core.Println(core.Sprintf("service %s stopped", name))
//	    return core.Result{OK: true}
//	})
type ActionServiceShutdown struct{}

// ActionTaskStarted is broadcast when an asynchronous task begins.
//
//	event := core.ActionTaskStarted{TaskIdentifier: "task-42", Action: "agent.run"}
//	core.New().ACTION(event)
type ActionTaskStarted struct {
	TaskIdentifier string
	Action         string
	Options        Options
}

// ActionTaskProgress is broadcast when an asynchronous task reports progress.
//
//	event := core.ActionTaskProgress{TaskIdentifier: "task-42", Action: "agent.run", Progress: 0.75, Message: "indexing repo"}
//	core.New().ACTION(event)
type ActionTaskProgress struct {
	TaskIdentifier string
	Action         string
	Progress       float64
	Message        string
}

// ActionTaskCompleted is broadcast when an asynchronous task finishes.
//
//	event := core.ActionTaskCompleted{TaskIdentifier: "task-42", Action: "agent.run", Result: core.Result{Value: "done", OK: true}}
//	core.New().ACTION(event)
type ActionTaskCompleted struct {
	TaskIdentifier string
	Action         string
	Result         Result
}

// --- Constructor ---

// CoreOption is a functional option applied during Core construction.
// Returns Result — if !OK, New() stops and returns the error.
//
//	core.New(
//	    core.WithService(agentic.Register),
//	    core.WithService(monitor.Register),
//	    core.WithServiceLock(),
//	)
type CoreOption func(*Core) Result

// New initialises a Core instance by applying options in order.
// Services registered here form the application conclave — they share
// IPC access and participate in the lifecycle (ServiceStartup/ServiceShutdown).
//
//	c := core.New(
//	    core.WithOption("name", "myapp"),
//	    core.WithService(auth.Register),
//	    core.WithServiceLock(),
//	)
//	c.Run()
func New(opts ...CoreOption) *Core {
	c := newCore()

	for _, opt := range opts {
		if r := opt(c); !r.OK {
			Error("core.New failed", "err", r.Value)
			break
		}
	}

	// Apply service lock after all opts — v0.3.3 parity
	c.LockApply()

	return c
}

// MustNew is New for package-var bundles: the first failed option
// panics, so a broken bundle fails at import time instead of
// half-constructing silently (New logs and continues — fine in main(),
// fatal in a package var nobody inspects).
//
//	var Widgets = core.MustNew(
//	    core.WithService(widgets.Register),
//	    core.WithServiceLock(),
//	)
func MustNew(opts ...CoreOption) *Core {
	c := newCore()

	for _, opt := range opts {
		if r := opt(c); !r.OK {
			panic(E("core.MustNew", Sprint("option failed: ", r.Value), nil))
		}
	}

	c.LockApply()

	return c
}

// newCore builds the bare Core all constructors share — subsystems
// wired, no options applied.
func newCore() *Core {
	c := &Core{
		app:                &App{},
		data:               &Data{Registry: NewRegistry[*Embed]()},
		drive:              &Drive{Registry: NewRegistry[*DriveHandle]()},
		fs:                 (&Fs{}).New("/"),
		config:             (&Config{}).New(),
		error:              &ErrorPanic{},
		log:                &ErrorLog{},
		locks:              NewRegistry[*Lock](),
		ipc:                &Ipc{actions: NewRegistry[*Action](), tasks: NewRegistry[*Task]()},
		info:               systemInfo,
		i18n:               &I18n{},
		api:                &API{protocols: NewRegistry[StreamFactory]()},
		services:           &ServiceRegistry{Registry: NewRegistry[*Service]()},
		commands:           &CommandRegistry{Registry: NewRegistry[*Command]()},
		entitlementChecker: defaultChecker,
	}
	c.context, c.cancel = WithCancel(Background())
	c.api.core = c
	c.bootTime = Now()
	registerBuiltins(c)
	return c
}

// WithOptions applies key-value configuration to Core.
//
//	core.WithOptions(core.NewOptions(core.Option{Key: "name", Value: "myapp"}))
func WithOptions(opts Options) CoreOption {
	return func(c *Core) Result {
		// Merge, never clobber (W2-4): WithOption keys set earlier in the
		// option list survive a later WithOptions.
		if c.options == nil {
			c.options = &opts
		} else {
			for _, opt := range opts.Items() {
				c.options.Set(opt.Key, opt.Value)
			}
		}
		if name := opts.String("name"); name != "" {
			c.app.Name = name
		}
		return Result{OK: true}
	}
}

// WithService registers a service via its factory function.
// If the factory returns a non-nil Value, WithService auto-discovers the
// service name from the factory's package path (last path segment, lowercase,
// with any "_test" suffix stripped) and calls RegisterService on the instance.
// IPC handler auto-registration is handled by RegisterService.
//
// If the factory returns nil Value (it registered itself), WithService
// returns success without a second registration.
//
//	core.WithService(agentic.Register)
//	core.WithService(display.Register(nil))
func WithService(factory func(*Core) Result) CoreOption {
	return func(c *Core) Result {
		r := factory(c)
		if !r.OK {
			return r
		}
		if r.Value == nil {
			// Factory self-registered — nothing more to do.
			return Result{OK: true}
		}
		// Auto-discover the service name from the instance's package path.
		instance := r.Value
		typeOf := TypeOf(instance)
		if typeOf.Kind() == KindPointer {
			typeOf = typeOf.Elem()
		}
		pkgPath := typeOf.PkgPath()
		parts := Split(pkgPath, "/")
		name := Lower(parts[len(parts)-1])
		if name == "" {
			return Result{E("core.WithService", Sprintf("service name could not be discovered for type %T", instance), nil), false}
		}

		// RegisterService handles Startable/Stoppable/HandleIPCEvents discovery
		return c.RegisterService(name, instance)
	}
}

// WithName registers a service with an explicit name (no reflect discovery).
//
//	core.WithName("ws", func(c *Core) Result {
//	    return Result{Value: hub, OK: true}
//	})
func WithName(name string, factory func(*Core) Result) CoreOption {
	return func(c *Core) Result {
		r := factory(c)
		if !r.OK {
			return r
		}
		if r.Value == nil {
			return Result{E("core.WithName", Sprintf("failed to create service %q", name), nil), false}
		}
		return c.RegisterService(name, r.Value)
	}
}

// WithOption is a convenience for setting a single key-value option.
//
//	core.New(
//	    core.WithOption("name", "myapp"),
//	    core.WithOption("port", 8080),
//	)
func WithOption(key string, value any) CoreOption {
	return func(c *Core) Result {
		if c.options == nil {
			opts := NewOptions()
			c.options = &opts
		}
		c.options.Set(key, value)
		if key == "name" {
			if s, ok := value.(string); ok {
				c.app.Name = s
			}
		}
		return Result{OK: true}
	}
}

// WithServiceLock prevents further service registration after construction.
//
//	core.New(
//	    core.WithService(auth.Register),
//	    core.WithServiceLock(),
//	)
func WithServiceLock() CoreOption {
	return func(c *Core) Result {
		c.LockEnable()
		return Result{OK: true}
	}
}

// WithCrashFile sets the crash-report file for panic recovery — the
// exported seam for ErrorPanic's report sink. Recover appends reports
// there; Reports reads them back.
//
//	core.New(core.WithCrashFile("/var/log/myapp/crash.json"))
func WithCrashFile(path string) CoreOption {
	return func(c *Core) Result {
		if path == "" {
			return Result{E("core.WithCrashFile", "path cannot be empty", nil), false}
		}
		c.error.filePath = path
		return Result{OK: true}
	}
}

// WithCli registers the built-in CLI command framework as service "cli".
// core.New no longer auto-registers it — opt in here for a package's basic
// compile-and-run binary, or bring an intentional CLI (dappco.re/go/cli).
//
//	core.New(core.WithCli())
func WithCli() CoreOption {
	return func(c *Core) Result {
		return CliRegister(c)
	}
}

// WithConfigFile loads a JSON config file into c.Config() during
// construction. A missing or malformed file fails the option — MustNew
// panics, New logs and continues. For optional files call
// c.Config().Load at runtime instead.
//
//	core.New(core.WithConfigFile("/etc/myapp/config.json"))
func WithConfigFile(path string) CoreOption {
	return func(c *Core) Result {
		return c.config.Load(path)
	}
}

// WithEnvConfig imports prefixed environment variables into c.Config()
// during construction (MYAPP_DATABASE_HOST → "database.host").
//
//	core.New(core.WithEnvConfig("MYAPP_"))
func WithEnvConfig(prefix string) CoreOption {
	return func(c *Core) Result {
		return c.config.FromEnv(prefix)
	}
}

// WithReloadOnSIGHUP completes the daemon loop: when a signal service
// (e.g. go-process) broadcasts signal.received with SIGHUP, run
// ServiceReload. Fails when "signal.received" already has a handler —
// call ServiceReload from that handler instead (the signal contract is
// single-handler by design; see signal.go).
//
//	core.New(core.WithService(process.Register), core.WithReloadOnSIGHUP())
func WithReloadOnSIGHUP() CoreOption {
	return func(c *Core) Result {
		if c.Action("signal.received").Exists() {
			return Result{E("core.WithReloadOnSIGHUP", "signal.received already handled — call ServiceReload from your handler", nil), false}
		}
		c.Action("signal.received", func(ctx Context, opts Options) Result {
			if opts.String("name") == "SIGHUP" {
				return c.ServiceReload(ctx)
			}
			return Result{OK: true}
		})
		return Result{OK: true}
	}
}

// WithBundle mounts another Core's capability surface under a dotted
// prefix — sealed toolkits composing into applications. The bundle's
// actions become "<prefix>.<name>" (delegated: the bundle's own
// entitlements and metering fire first, then the host gates the
// prefixed name — double-gated enclave semantics). Its data mounts and
// drive handles become "<prefix>.<name>". Collisions fail the option
// loudly. Mount before WithServiceLock so the seal freezes the composed
// surface.
//
//	var Widgets = core.MustNew(core.WithService(widgets.Register), core.WithServiceLock())
//
//	c := core.New(
//	    core.WithBundle("widgets", Widgets),
//	    core.WithServiceLock(),
//	)
//	c.Action("widgets.render").Run(ctx, opts)
func WithBundle(prefix string, other *Core) CoreOption {
	return func(c *Core) Result {
		if prefix == "" || other == nil {
			return Result{E("core.WithBundle", "prefix and bundle are both required", nil), false}
		}

		// Collision pre-pass — fail before mutating anything.
		for _, name := range other.Actions() {
			if c.Action(Concat(prefix, ".", name)).Exists() {
				return Result{E("core.WithBundle", Concat("action collision: ", prefix, ".", name), nil), false}
			}
		}
		for _, name := range other.data.Names() {
			if c.data.Has(Concat(prefix, ".", name)) {
				return Result{E("core.WithBundle", Concat("data mount collision: ", prefix, ".", name), nil), false}
			}
		}
		for _, name := range other.drive.Names() {
			if c.drive.Has(Concat(prefix, ".", name)) {
				return Result{E("core.WithBundle", Concat("drive handle collision: ", prefix, ".", name), nil), false}
			}
		}

		for _, name := range other.Actions() {
			orig := name
			c.Action(Concat(prefix, ".", name), func(ctx Context, opts Options) Result {
				return other.Action(orig).Run(ctx, opts)
			})
		}
		other.data.Each(func(name string, emb *Embed) {
			c.data.Set(Concat(prefix, ".", name), emb)
		})
		other.drive.Each(func(name string, handle *DriveHandle) {
			c.drive.Set(Concat(prefix, ".", name), handle)
		})
		return Result{OK: true}
	}
}
