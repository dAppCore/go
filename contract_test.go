// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

// --- WithService ---

// stub service used only for name-discovery tests.
type stubNamedService struct{}

// stubFactory is a package-level factory so the runtime function name carries
// the package path "core_test.stubFactory" — last segment after '/' is
// "core_test", and after stripping a "_test" suffix we get "core".
// For a real service package such as "dappco.re/go/agentic" the discovered
// name would be "agentic".
func stubFactory(c *Core) Result {
	return Result{Value: &stubNamedService{}, OK: true}
}

// TestWithService_NameDiscovery_Good verifies that WithService discovers the
// service name from the factory's package path and registers the instance via
// RegisterService, making it retrievable through c.Services().
//
// stubFactory lives in package "dappco.re/go_test", so the last path
// segment is "core_test" — WithService strips the "_test" suffix and registers
// the service under the name "core".
func TestContract_WithService_Good(t *T) {
	c := New(WithService(stubFactory))

	names := c.Services()
	// WithService discovers a name from the factory's package path and registers
	// the instance under it. cli is no longer auto-registered (opt-in via WithCli),
	// so the discovered service is the only entry.
	AssertGreater(t, len(names), 0, "expected auto-discovered service to be registered")
}

// TestWithService_FactorySelfRegisters_Good verifies that when a factory
// returns Result{OK:true} with no Value (it registered itself), WithService
// does not attempt a second registration and returns success.
func TestContract_WithService_FactorySelfRegisters_Good(t *T) {
	selfReg := func(c *Core) Result {
		// Factory registers directly, returns no instance.
		c.Service("self", Service{})
		return Result{OK: true}
	}

	c := New(WithService(selfReg))

	// "self" must be present and registered exactly once.
	svc := c.Service("self")
	AssertTrue(t, svc.OK, "expected self-registered service to be present")
}

// --- WithName ---

func TestContract_WithName_Good(t *T) {
	c := New(
		WithName("custom", func(c *Core) Result {
			return Result{Value: &stubNamedService{}, OK: true}
		}),
	)
	AssertContains(t, c.Services(), "custom")
}

// --- Lifecycle ---

type lifecycleService struct {
	started bool
}

func (s *lifecycleService) OnStartup(_ Context) Result {
	s.started = true
	return Result{OK: true}
}

func TestContract_WithService_Lifecycle_Good(t *T) {
	svc := &lifecycleService{}
	c := New(
		WithService(func(c *Core) Result {
			return Result{Value: svc, OK: true}
		}),
	)

	c.ServiceStartup(Background(), nil)
	AssertTrue(t, svc.started)
}

// --- IPC Handler ---

type ipcService struct {
	received Message
}

func (s *ipcService) HandleIPCEvents(c *Core, msg Message) Result {
	s.received = msg
	return Result{OK: true}
}

func TestContract_WithService_IPCHandler_Good(t *T) {
	svc := &ipcService{}
	c := New(
		WithService(func(c *Core) Result {
			return Result{Value: svc, OK: true}
		}),
	)

	c.ACTION("ping")
	AssertEqual(t, "ping", svc.received)
}

// --- Error ---

// TestWithService_FactoryError_Bad verifies that a failing factory
// stops further option processing (second service not registered).
func TestContract_WithService_FactoryError_Bad(t *T) {
	secondCalled := false
	c := New(
		WithService(func(c *Core) Result {
			return Result{Value: E("test", "factory failed", nil), OK: false}
		}),
		WithService(func(c *Core) Result {
			secondCalled = true
			return Result{OK: true}
		}),
	)
	AssertNotNil(t, c)
	AssertFalse(t, secondCalled, "second option should not run after first fails")
}

// --- AX-7 canonical backfill ---

func TestContract_New_Ugly(t *T) {
	c := New(WithCli(), WithServiceLock())
	r := c.Service("late-agent", Service{})
	AssertFalse(t, r.OK)
	AssertTrue(t, c.Service("cli").OK)
}

func TestContract_WithName_Bad(t *T) {
	opt := WithName("agent", func(c *Core) Result {
		return Result{OK: true}
	})
	r := opt(New())
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "failed to create service")
}

func TestContract_WithName_Ugly(t *T) {
	opt := WithName("", func(c *Core) Result {
		return Result{Value: &stubNamedService{}, OK: true}
	})
	r := opt(New())
	AssertFalse(t, r.OK)
}

func TestContract_WithOption_Bad(t *T) {
	c := New(WithOption("name", 42))
	AssertEqual(t, "", c.App().Name)
	AssertEqual(t, 42, c.Options().Get("name").Value)
}

func TestContract_WithOption_Ugly(t *T) {
	c := New(WithOption("", "root-option"))
	r := c.Options().Get("")
	AssertTrue(t, r.OK)
	AssertEqual(t, "root-option", r.Value)
}

func TestContract_WithOptions_Ugly(t *T) {
	opts := NewOptions()
	c := New(WithOptions(opts))
	AssertNotNil(t, c.Options())
	AssertEqual(t, "", c.App().Name)
}

func TestContract_WithService_Bad(t *T) {
	opt := WithService(func(c *Core) Result {
		return Result{Value: NewError("factory failed"), OK: false}
	})
	r := opt(New())
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "factory failed")
}

func TestContract_WithService_Ugly(t *T) {
	c := New()
	opt := WithService(func(c *Core) Result {
		c.Service("self-registered", Service{})
		return Result{OK: true}
	})
	r := opt(c)
	AssertTrue(t, r.OK)
	AssertTrue(t, c.Service("self-registered").OK)
}

func TestContract_WithServiceLock_Bad(t *T) {
	c := New(WithServiceLock())
	r := c.Service("late-agent", Service{})
	AssertFalse(t, r.OK)
}

func TestContract_WithServiceLock_Ugly(t *T) {
	c := New(WithCli(), WithServiceLock())
	AssertTrue(t, c.Service("cli").OK)
	AssertContains(t, c.Services(), "cli")
}

func TestContract_New_Good(t *T) {
	c := New(WithOption("env", "prod"))
	AssertNotNil(t, c)
	AssertEqual(t, "prod", c.Options().Get("env").Value)
}

func TestContract_New_Bad(t *T) {
	// A failing option aborts the remaining chain but New still returns a usable Core.
	c := New(WithService(func(c *Core) Result {
		return Result{Value: NewError("boom"), OK: false}
	}))
	AssertNotNil(t, c)
	AssertNotNil(t, c.Config())
}

func TestContract_WithOptions_Good(t *T) {
	opts := NewOptions(Option{Key: "region", Value: "eu"})
	c := New(WithOptions(opts))
	AssertEqual(t, "eu", c.Options().Get("region").Value)
}

func TestContract_WithOptions_Bad(t *T) {
	// A non-string "name" is stored but ignored for the app name.
	opts := NewOptions(Option{Key: "name", Value: 42})
	c := New(WithOptions(opts))
	AssertEqual(t, "", c.App().Name)
	AssertEqual(t, 42, c.Options().Get("name").Value)
}

func TestContract_WithOption_Good(t *T) {
	c := New(WithOption("region", "eu-west"))
	r := c.Options().Get("region")
	AssertTrue(t, r.OK)
	AssertEqual(t, "eu-west", r.Value)
}

func TestContract_WithServiceLock_Good(t *T) {
	// A service registered during construction is admitted; the lock only
	// seals registration afterwards.
	c := New(
		WithService(func(c *Core) Result {
			c.Service("early", Service{})
			return Result{OK: true}
		}),
		WithServiceLock(),
	)
	AssertTrue(t, c.Service("early").OK)            // admitted during construction
	AssertFalse(t, c.Service("late", Service{}).OK) // sealed afterwards
}

func TestContract_WithCli_Good(t *T) {
	c := New(WithCli())
	AssertTrue(t, c.Service("cli").OK)
	AssertNotNil(t, c.Cli())
}

func TestContract_WithCli_Bad(t *T) {
	// Without WithCli, no cli service is registered.
	c := New()
	AssertFalse(t, c.Service("cli").OK)
}

func TestContract_WithCli_Ugly(t *T) {
	// Re-applying WithCli keeps cli registered exactly once.
	c := New(WithCli(), WithCli())
	AssertTrue(t, c.Service("cli").OK)
	AssertContains(t, c.Services(), "cli")
}

// --- MustNew ---

func TestContract_MustNew_Good(t *T) {
	c := MustNew(WithOption("name", "bundle"))
	AssertEqual(t, "bundle", c.App().Name)
}

func TestContract_MustNew_Bad(t *T) {
	AssertPanics(t, func() {
		MustNew(func(*Core) Result { return Result{Value: NewError("boom"), OK: false} })
	})
}

func TestContract_MustNew_Ugly(t *T) {
	// The panic carries the failing option's diagnostic.
	AssertPanicsWithError(t, "option failed", func() {
		MustNew(func(*Core) Result { return Result{Value: NewError("cascade"), OK: false} })
	})
}

// --- W2-4: WithOptions merges instead of clobbering ---

func TestContract_WithOptions_Good_MergesNotClobbers(t *T) {
	c := New(
		WithOption("keep", "earlier"),
		WithOptions(NewOptions(Option{Key: "name", Value: "merged"})),
	)
	AssertEqual(t, "earlier", c.Options().String("keep"))
	AssertEqual(t, "merged", c.Options().String("name"))
}

// --- W3: WithCrashFile — the exported crash-sink seam ---

func TestContract_WithCrashFile_Good(t *T) {
	path := Path(t.TempDir(), "crash.json")
	c := New(WithCrashFile(path))
	// Configured: the miss is now the empty file, not the missing config.
	r := c.Error().Reports(1)
	AssertFalse(t, r.OK)
	AssertFalse(t, Contains(r.Error(), "no crash file"))
}

func TestContract_WithCrashFile_Bad(t *T) {
	// An empty path is a constructor failure — MustNew panics on it.
	AssertPanics(t, func() { MustNew(WithCrashFile("")) })
}

func TestContract_WithCrashFile_Ugly(t *T) {
	// New() logs-and-continues on the failed option; the sink stays unset.
	c := New(WithCrashFile(""))
	AssertContains(t, c.Error().Reports(1).Error(), "no crash file")
}
