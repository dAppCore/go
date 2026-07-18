// SPDX-License-Identifier: EUPL-1.2

// Package core is a dependency injection and service lifecycle framework for Go.
// This file defines the Core struct, accessors, and IPC/error wrappers.

package core

// --- Core Struct ---

// Core is the central application object that manages services, assets, and communication.
//
//	c := core.New(core.WithOption("name", "homelab"))
//	ctx := c.Context()
//	_ = ctx
type Core struct {
	options *Options    // c.Options()        — Input configuration used to create this Core
	app     *App        // c.App()            — Application identity + optional GUI runtime
	data    *Data       // c.Data()           — Embedded/stored content from packages
	drive   *Drive      // c.Drive()          — Resource handle registry (transports)
	fs      *Fs         // c.Fs()             — Local filesystem I/O (sandboxable)
	config  *Config     // c.Config()         — Configuration, settings, feature flags
	error   *ErrorPanic // c.Error()          — Panic recovery and crash reporting
	log     *ErrorLog   // c.Log()            — Structured logging + error wrapping
	// cli accessed via c.Cli() — CLI command framework (opt-in service "cli", core.WithCli())
	commands *CommandRegistry // c.Command("path")  — Command tree
	services *ServiceRegistry // c.Service("name")  — Service registry
	locks    *Registry[*Lock] // c.Lock("name")     — Named mutexes
	ipc      *Ipc             // c.IPC()            — Message bus for IPC
	api      *API             // c.API()            — Remote streams
	info     *SysInfo         // c.Env("key")        — Read-only system/environment information
	i18n     *I18n            // c.I18n()           — Internationalisation and locale collection

	entitlementChecker EntitlementChecker // default: everything permitted
	usageRecorder      UsageRecorder      // default: nil (no-op)

	context       Context
	cancel        CancelFunc
	taskIDCounter AtomicUint64
	waitGroup     WaitGroup
	shutdown      AtomicBool
}

// --- Accessors ---

// Options returns the input configuration passed to core.New().
//
//	opts := c.Options()
//	name := opts.String("name")
func (c *Core) Options() *Options { return c.options }

// App returns application identity metadata.
//
//	c.App().Name     // "my-app"
//	c.App().Version  // "1.0.0"
func (c *Core) App() *App { return c.app }

// Data returns the embedded asset registry (Registry[*Embed]).
//
//	r := c.Data().ReadString("prompts/coding.md")
func (c *Core) Data() *Data { return c.data }

// Drive returns the transport handle registry (Registry[*DriveHandle]).
//
//	r := c.Drive().Get("forge")
func (c *Core) Drive() *Drive { return c.drive }

// Fs returns the sandboxed filesystem.
//
//	r := c.Fs().Read("/path/to/file")
//	c.Fs().WriteAtomic("/status.json", data)
func (c *Core) Fs() *Fs { return c.fs }

// Config returns runtime settings and feature flags. With a group argument it
// returns a Config view scoped under that key prefix (see Config.Group).
//
//	host := c.Config().String("database.host")
//	c.Config().Enable("dark-mode")
//	c.Config("database").Set("host", "localhost") // stores "database.host"
func (c *Core) Config(group ...string) *Config {
	if len(group) == 0 || group[0] == "" {
		return c.config
	}
	return c.config.Group(group[0])
}

// Feature returns a keyed handle to a single feature flag (name required).
//
//	c.Feature("dark-mode").Enable()
//	if c.Feature("dark-mode").Enabled() { core.Println("on") }
func (c *Core) Feature(name string) Feature {
	return Feature{cfg: c.config, name: name}
}

// Error returns the panic recovery subsystem.
//
//	c.Error().Recover()
func (c *Core) Error() *ErrorPanic { return c.error }

// Log returns the structured logging subsystem.
//
//	c.Log().Info("started", "port", 8080)
func (c *Core) Log() *ErrorLog { return c.log }

// Cli returns the CLI command framework (registered as service "cli").
//
//	c.Cli().Run("deploy", "to", "homelab")
func (c *Core) Cli() *Cli {
	cl, _ := ServiceFor[*Cli](c, "cli")
	return cl
}

// IPC returns the message bus internals.
//
//	c.IPC()
func (c *Core) IPC() *Ipc { return c.ipc }

// I18n returns the internationalisation subsystem.
//
//	tr := c.I18n().Translate("cmd.deploy.description")
func (c *Core) I18n() *I18n { return c.i18n }

// Env returns an environment variable by key (cached at init, falls back to os.Getenv).
//
//	home := c.Env("DIR_HOME")
//	token := c.Env("FORGE_TOKEN")
func (c *Core) Env(key string) string { return Env(key) }

// Context returns Core's lifecycle context (cancelled on shutdown).
//
//	ctx := c.Context()
func (c *Core) Context() Context { return c.context }

// WithContext returns a shallow clone of c whose lifecycle context is
// derived from ctx — for request-scoped context derivation (auth
// substrate etc.). The derived Core shares services/data/config (every
// heavy subsystem) with the parent by pointer; only the context, its
// cancel, and the per-Core lifecycle bookkeeping (waitGroup / shutdown
// flag / task counter) are fresh.
//
// Cancellation is one-directional: the derived Core's cancel does NOT
// cancel the parent, but a parent shutdown propagates DOWN to the
// derived context when ctx chains from the parent (the usual case —
// pass core.WithValue(c.Context(), …)). The clone re-wraps ctx in a
// fresh WithCancel so callers that retain the parent stay unaffected.
//
//	type userKey struct{}
//	rc := c.WithContext(core.WithValue(c.Context(), userKey{}, user))
//	id := rc.Context().Value(userKey{})  // round-trips on the clone
func (c *Core) WithContext(ctx Context) *Core {
	derivedCtx, derivedCancel := WithCancel(ctx)
	return &Core{
		options:            c.options,
		app:                c.app,
		data:               c.data,
		drive:              c.drive,
		fs:                 c.fs,
		config:             c.config,
		error:              c.error,
		log:                c.log,
		commands:           c.commands,
		services:           c.services,
		locks:              c.locks,
		ipc:                c.ipc,
		api:                c.api,
		info:               c.info,
		i18n:               c.i18n,
		entitlementChecker: c.entitlementChecker,
		usageRecorder:      c.usageRecorder,
		context:            derivedCtx,
		cancel:             derivedCancel,
		// taskIDCounter / waitGroup / shutdown intentionally start fresh —
		// the derived Core owns its own request-scoped lifecycle.
	}
}

// Core returns self — satisfies the ServiceRuntime interface.
//
//	c := s.Core()
func (c *Core) Core() *Core { return c }

// --- Lifecycle ---

// RunResult starts all services, runs the CLI, then shuts down.
// Returns Result so main() can decide how to handle failure.
// ServiceShutdown is always called via defer, even on startup failure or panic.
//
//	r := c.RunResult()
//	if !r.OK { core.Exit(1) }
func (c *Core) RunResult() Result {
	defer c.ServiceShutdown(Background())

	r := c.ServiceStartup(c.context, nil)
	if !r.OK {
		if _, ok := r.Value.(error); !ok {
			return Result{Value: NewCode("core.run.startup", "startup failed"), OK: false}
		}
		return r
	}

	if cli := c.Cli(); cli != nil {
		r = cli.Run()
	}

	// "cli.noop" is the CLI's benign no-op sentinel (banner/help shown,
	// nothing to run) — success. Everything else propagates: a valueless
	// Result{OK: false} is a real failure, no longer inferred benign
	// from its nil Value (W2-1; the old inference silently converted any
	// bare failure into success).
	if !r.OK && r.Code() == "cli.noop" {
		return Result{OK: true}
	}
	return r
}

// Run starts all services, runs the CLI, then shuts down. Calls
// c.Exit(1) on failure (graceful shutdown chain, 30s timeout). For
// programmatic error handling use RunResult().
//
//	c := core.New(core.WithService(myService.Register))
//	c.Run()
func (c *Core) Run() {
	r := c.RunResult()
	if !r.OK {
		Error(r.Error())
		c.Exit(1)
	}
}

// --- IPC (uppercase aliases) ---

// ACTION broadcasts a message to all registered handlers (fire-and-forget).
// Each handler is wrapped in panic recovery. All handlers fire regardless.
//
//	c.ACTION(messages.AgentCompleted{Agent: "codex", Status: "completed"})
func (c *Core) ACTION(msg Message) Result { return c.broadcast(msg) }

// QUERY sends a request — first handler to return OK wins.
//
//	r := c.QUERY(MyQuery{Name: "brain"})
func (c *Core) QUERY(q Query) Result { return c.Query(q) }

// QUERYALL sends a request — collects all OK responses.
//
//	r := c.QUERYALL(countQuery{})
//	results := r.Value.([]any)
func (c *Core) QUERYALL(q Query) Result { return c.QueryAll(q) }

// --- Error+Log ---

// LogError logs an error and returns the Result from ErrorLog.
//
//	c := core.New()
//	err := core.NewError("homelab unreachable")
//	r := c.LogError(err, "agent.Ping", "health check failed")
//	if !r.OK { return r }
func (c *Core) LogError(err error, op, msg string) Result {
	return c.log.Error(err, op, msg)
}

// LogWarn logs a warning and returns the Result from ErrorLog.
//
//	c := core.New()
//	err := core.NewError("config.host missing")
//	r := c.LogWarn(err, "config.Load", "using default host")
//	if !r.OK { return r }
func (c *Core) LogWarn(err error, op, msg string) Result {
	return c.log.Warn(err, op, msg)
}

// Must logs and panics if err is not nil.
//
//	c := core.New()
//	c.Must(nil, "agent.Start", "startup failed")
func (c *Core) Must(err error, op, msg string) {
	c.log.Must(err, op, msg)
}

// --- Registry Accessor ---

// RegistryOf returns a point-in-time snapshot of a named registry, wrapped in a
// Result, for cross-cutting queries (Names/Len/Has/List). Known names:
// "services", "commands", "actions", "tasks", "locks", "data", "drive", "api",
// "embed", "features". An unknown name yields Result{OK:false} — distinguishable
// from a known-but-empty registry. The snapshot is a copy taken at call time and
// does NOT track later changes; call again for a fresh view.
//
//	r := c.RegistryOf("services")
//	if r.OK { names := r.Value.(*core.Registry[any]).Names() }
func (c *Core) RegistryOf(name string) Result {
	// Bridge typed registries to untyped access for cross-cutting queries.
	// Each registry is wrapped in a read-only snapshot proxy.
	switch name {
	case "services":
		return Ok(registryProxy(c.services.Registry))
	case "commands":
		return Ok(registryProxy(c.commands.Registry))
	case "actions":
		return Ok(registryProxy(c.ipc.actions))
	case "tasks":
		return Ok(registryProxy(c.ipc.tasks))
	case "locks":
		return Ok(registryProxy(c.locks))
	case "data":
		return Ok(registryProxy(c.data.Registry))
	case "drive":
		return Ok(registryProxy(c.drive.Registry))
	case "api":
		return Ok(registryProxy(c.api.protocols))
	case "embed":
		return Ok(registryProxy(assetGroups))
	case "features":
		return Ok(registryProxy(c.config.featureFlags()))
	default:
		return Result{E("core.RegistryOf", Concat("unknown registry: \"", name, "\""), nil), false}
	}
}

// registryProxy creates a read-only any-typed snapshot of a typed registry.
// Copies current state — not a live view (avoids type parameter leaking).
func registryProxy[T any](src *Registry[T]) *Registry[any] {
	proxy := NewRegistry[any]()
	src.Each(func(name string, item T) {
		proxy.Set(name, item)
	})
	return proxy
}

// --- Global Instance ---
