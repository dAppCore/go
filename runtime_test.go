package core_test

import (
	. "dappco.re/go"
)

// --- ServiceRuntime ---

type testOpts struct {
	URL     string
	Timeout int
}

// --- NewWithFactories ---

func TestRuntime_NewWithFactories_Good(t *T) {
	r := NewWithFactories(nil, map[string]ServiceFactory{
		"svc1": func() Result { return Result{Value: Service{}, OK: true} },
		"svc2": func() Result { return Result{Value: Service{}, OK: true} },
	})
	AssertTrue(t, r.OK)
	rt := r.Value.(*Runtime)
	AssertNotNil(t, rt.Core)
}

func TestRuntime_NewWithFactories_NilFactory_Good(t *T) {
	r := NewWithFactories(nil, map[string]ServiceFactory{
		"bad": nil,
	})
	AssertTrue(t, r.OK) // nil factories skipped
}

func TestRuntime_NewRuntime_Good(t *T) {
	r := NewRuntime(nil)
	AssertTrue(t, r.OK)
}

// --- Runtime lifecycle edge: shutdown on a zero-value Runtime ---

func TestRuntime_ServiceShutdown_NilCore_Good(t *T) {
	rt := &Runtime{}
	result := rt.ServiceShutdown(Background())
	AssertTrue(t, result.OK)
}

func TestCore_Context_Good(t *T) {
	c := New()
	c.ServiceStartup(Background(), nil)
	AssertNotNil(t, c.Context())
	c.ServiceShutdown(Background())
}

func TestRuntime_Core_ServiceShutdown_Good(t *T) {
	stopped := false
	c := New()
	c.Service("agent", Service{
		OnStart: func() Result { return Result{OK: true} },
		OnStop:  func() Result { stopped = true; return Result{OK: true} },
	})
	AssertTrue(t, c.ServiceStartup(Background(), nil).OK)

	r := c.ServiceShutdown(Background())

	AssertTrue(t, r.OK)
	AssertTrue(t, stopped)
}

func TestRuntime_Core_ServiceShutdown_Bad(t *T) {
	c := New()
	c.Service("agent", Service{
		OnStop: func() Result { return Result{Value: NewError("shutdown refused"), OK: false} },
	})

	r := c.ServiceShutdown(Background())

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "shutdown refused")
}

func TestRuntime_Core_ServiceShutdown_Ugly(t *T) {
	c := New()
	c.Service("agent", Service{
		OnStop: func() Result { return Result{OK: true} },
	})
	ctx, cancel := WithCancel(Background())
	cancel()

	r := c.ServiceShutdown(ctx)

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestRuntime_Core_ServiceStartup_Good(t *T) {
	started := false
	c := New()
	c.Service("agent", Service{OnStart: func() Result { started = true; return Result{OK: true} }})

	r := c.ServiceStartup(Background(), nil)

	AssertTrue(t, r.OK)
	AssertTrue(t, started)
}

func TestRuntime_Core_ServiceStartup_Bad(t *T) {
	c := New()
	c.Service("agent", Service{OnStart: func() Result { return Result{OK: true} }})
	ctx, cancel := WithCancel(Background())
	cancel()

	r := c.ServiceStartup(ctx, nil)

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestRuntime_Core_ServiceStartup_Ugly(t *T) {
	c := New()
	var sawStartup bool
	c.RegisterAction(func(_ *Core, msg Message) Result {
		_, sawStartup = msg.(ActionServiceStartup)
		return Result{OK: true}
	})

	r := c.ServiceStartup(Background(), nil)

	AssertTrue(t, r.OK)
	AssertTrue(t, sawStartup)
}

func TestRuntime_Core_Go_Good(t *T) {
	c := New()
	done := make(chan bool, 1)
	c.Go(func() { done <- true })
	select {
	case <-done:
	case <-After(2 * Second):
		t.Fatal("c.Go did not run the function")
	}
}

func TestRuntime_Core_Go_Bad(t *T) {
	// Go returns immediately; the fn runs asynchronously, not inline.
	c := New()
	release := make(chan bool)
	observed := make(chan bool, 1)
	c.Go(func() {
		<-release // still blocked when Go has already returned
		observed <- true
	})
	close(release)
	select {
	case <-observed:
	case <-After(2 * Second):
		t.Fatal("goroutine did not complete after release")
	}
}

func TestRuntime_Core_Go_Ugly(t *T) {
	// Many concurrent goroutines all run.
	c := New()
	const n = 16
	done := make(chan int, n)
	for i := 0; i < n; i++ {
		c.Go(func() { done <- 1 })
	}
	sum := 0
	for i := 0; i < n; i++ {
		select {
		case <-done:
			sum++
		case <-After(2 * Second):
			t.Fatal("a goroutine did not run")
		}
	}
	AssertEqual(t, n, sum)
}

func TestRuntime_Core_IsShutdown_Good(t *T) {
	// A fresh core is running, not shut down.
	AssertFalse(t, New().IsShutdown())
}

func TestRuntime_Core_IsShutdown_Bad(t *T) {
	c := New()
	c.ServiceShutdown(Background())
	AssertTrue(t, c.IsShutdown())
}

func TestRuntime_Core_IsShutdown_Ugly(t *T) {
	// Shutdown is sticky and idempotent.
	c := New()
	AssertFalse(t, c.IsShutdown())
	c.ServiceShutdown(Background())
	AssertTrue(t, c.IsShutdown())
	c.ServiceShutdown(Background())
	AssertTrue(t, c.IsShutdown())
}

func TestRuntime_ServiceRuntime_Config_Good(t *T) {
	c := New()
	rt := NewServiceRuntime(c, testOpts{})

	AssertSame(t, c.Config(), rt.Config())
}

func TestRuntime_ServiceRuntime_Config_Bad(t *T) {
	rt := NewServiceRuntime[testOpts](nil, testOpts{})
	AssertPanics(t, func() {
		_ = rt.Config()
	})
}

func TestRuntime_ServiceRuntime_Config_Ugly(t *T) {
	c := New()
	c.Config().Set("agent.region", "lab")
	rt := NewServiceRuntime(c, testOpts{})

	AssertEqual(t, "lab", rt.Config().String("agent.region"))
}

func TestRuntime_ServiceRuntime_Core_Good(t *T) {
	c := New()
	rt := NewServiceRuntime(c, testOpts{})
	AssertSame(t, c, rt.Core())
}

func TestRuntime_ServiceRuntime_Core_Bad(t *T) {
	rt := NewServiceRuntime[testOpts](nil, testOpts{})
	AssertNil(t, rt.Core())
}

func TestRuntime_ServiceRuntime_Core_Ugly(t *T) {
	c := New()
	rt := NewServiceRuntime(c, testOpts{})
	c.ServiceShutdown(Background())
	AssertSame(t, c, rt.Core())
}

func TestRuntime_ServiceRuntime_Options_Good(t *T) {
	opts := testOpts{URL: "https://api.lthn.ai", Timeout: 30}
	rt := NewServiceRuntime(New(), opts)
	AssertEqual(t, opts, rt.Options())
}

func TestRuntime_ServiceRuntime_Options_Bad(t *T) {
	rt := NewServiceRuntime(New(), testOpts{})
	AssertEqual(t, testOpts{}, rt.Options())
}

func TestRuntime_ServiceRuntime_Options_Ugly(t *T) {
	type queueOptions struct {
		Topics []string
	}
	opts := queueOptions{Topics: []string{}}
	rt := NewServiceRuntime(New(), opts)
	AssertEqual(t, opts, rt.Options())
}

func TestRuntime_NewRuntime_Bad(t *T) {
	r := NewRuntime("gui-runtime")
	AssertTrue(t, r.OK)
	AssertNotNil(t, r.Value.(*Runtime).Core)
}

func TestRuntime_NewRuntime_Ugly(t *T) {
	app := struct{ Name string }{Name: "agent-ui"}
	r := NewRuntime(app)
	AssertTrue(t, r.OK)
	AssertEqual(t, "core", r.Value.(*Runtime).Core.App().Name)
}

func TestRuntime_NewServiceRuntime_Good(t *T) {
	c := New()
	rt := NewServiceRuntime(c, testOpts{URL: "https://api.lthn.ai", Timeout: 30})
	AssertSame(t, c, rt.Core())
	AssertEqual(t, 30, rt.Options().Timeout)
}

func TestRuntime_NewServiceRuntime_Bad(t *T) {
	rt := NewServiceRuntime[testOpts](nil, testOpts{URL: "offline"})
	AssertNil(t, rt.Core())
	AssertEqual(t, "offline", rt.Options().URL)
}

func TestRuntime_NewServiceRuntime_Ugly(t *T) {
	// Options is captured by value at construction — mutating the caller's
	// struct afterwards must not be visible through the runtime.
	opts := testOpts{URL: "https://api.lthn.ai", Timeout: 30}
	rt := NewServiceRuntime(New(), opts)

	opts.URL = "mutated-after-construction"

	AssertEqual(t, "https://api.lthn.ai", rt.Options().URL)
}

func TestRuntime_NewWithFactories_Bad(t *T) {
	r := NewWithFactories(nil, map[string]ServiceFactory{
		"agent": func() Result { return Result{Value: NewError("factory refused"), OK: false} },
	})
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "factory \"agent\" failed")
}

func TestRuntime_NewWithFactories_Ugly(t *T) {
	r := NewWithFactories(nil, map[string]ServiceFactory{
		"agent": func() Result { return Result{Value: "not a service", OK: true} },
	})
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "non-Service")
}

func TestRuntime_Runtime_ServiceName_Good(t *T) {
	AssertEqual(t, "Core", (&Runtime{}).ServiceName())
}

func TestRuntime_Runtime_ServiceName_Bad(t *T) {
	var r *Runtime
	AssertEqual(t, "Core", r.ServiceName())
}

func TestRuntime_Runtime_ServiceName_Ugly(t *T) {
	// ServiceName is a fixed identity — even after the wrapped Core has
	// shut down, it still reports "Core".
	r := &Runtime{Core: New()}
	r.Core.ServiceShutdown(Background())

	AssertEqual(t, "Core", r.ServiceName())
}

func TestRuntime_Runtime_ServiceShutdown_Good(t *T) {
	stopped := false
	rt := &Runtime{Core: New()}
	rt.Core.Service("agent", Service{OnStop: func() Result { stopped = true; return Result{OK: true} }})

	r := rt.ServiceShutdown(Background())

	AssertTrue(t, r.OK)
	AssertTrue(t, stopped)
}

func TestRuntime_Runtime_ServiceShutdown_Bad(t *T) {
	rt := &Runtime{Core: New()}
	rt.Core.Service("agent", Service{OnStop: func() Result { return Result{Value: NewError("stop refused"), OK: false} }})

	r := rt.ServiceShutdown(Background())

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "stop refused")
}

func TestRuntime_Runtime_ServiceShutdown_Ugly(t *T) {
	rt := &Runtime{}
	r := rt.ServiceShutdown(Background())
	AssertTrue(t, r.OK)
}

func TestRuntime_Runtime_ServiceStartup_Good(t *T) {
	started := false
	rt := &Runtime{Core: New()}
	rt.Core.Service("agent", Service{OnStart: func() Result { started = true; return Result{OK: true} }})

	r := rt.ServiceStartup(Background(), nil)

	AssertTrue(t, r.OK)
	AssertTrue(t, started)
}

func TestRuntime_Runtime_ServiceStartup_Bad(t *T) {
	rt := &Runtime{Core: New()}
	rt.Core.Service("agent", Service{OnStart: func() Result { return Result{Value: NewError("start refused"), OK: false} }})

	r := rt.ServiceStartup(Background(), nil)

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "start refused")
}

func TestRuntime_Runtime_ServiceStartup_Ugly(t *T) {
	rt := &Runtime{}
	AssertPanics(t, func() {
		_ = rt.ServiceStartup(Background(), nil)
	})
}

// --- W2-3: ServiceStartup seals the conclave under WithServiceLock ---

func TestRuntime_ServiceStartup_Good_SealsConclave(t *T) {
	c := New(WithServiceLock())
	AssertTrue(t, c.ServiceStartup(Background(), nil).OK)
	// Post-startup the capability surface is frozen: late registration
	// does not land in the registry.
	c.Action("late.register", func(Context, Options) Result { return Ok(nil) })
	AssertFalse(t, c.Action("late.register").Exists())
}

func TestRuntime_ServiceStartup_Bad_NoLockNoSeal(t *T) {
	c := New()
	AssertTrue(t, c.ServiceStartup(Background(), nil).OK)
	// Without WithServiceLock, the surface stays open.
	c.Action("late.register", func(Context, Options) Result { return Ok(nil) })
	AssertTrue(t, c.Action("late.register").Exists())
}

func TestRuntime_ServiceStartup_Ugly_FeaturesSealedNotFrozen(t *T) {
	c := New(WithServiceLock())
	c.Feature("dark-mode").Enable()
	AssertTrue(t, c.ServiceStartup(Background(), nil).OK)
	// Existing flags stay toggleable (Sealed), new flags cannot appear.
	c.Feature("dark-mode").Disable()
	AssertFalse(t, c.Feature("dark-mode").Enabled())
	c.Feature("brand-new").Enable()
	AssertFalse(t, c.Feature("brand-new").Enabled())
}

// --- W3-5: the OnReload runner ---

type reloadProbe struct{ count int }

func (r *reloadProbe) OnReload(Context) Result {
	r.count++
	return Ok(nil)
}

type failingReloader struct{}

func (f *failingReloader) OnReload(Context) Result {
	return Fail(NewError("reload broke"))
}

func TestRuntime_Core_ServiceReload_Good(t *T) {
	c := New()
	p := &reloadProbe{}
	AssertTrue(t, c.RegisterService("probe", p).OK)
	AssertTrue(t, c.ServiceReload(Background()).OK)
	AssertEqual(t, 1, p.count)
}

func TestRuntime_Core_ServiceReload_Bad(t *T) {
	c := New()
	AssertTrue(t, c.RegisterService("failing", &failingReloader{}).OK)
	AssertFalse(t, c.ServiceReload(Background()).OK)
}

func TestRuntime_Core_ServiceReload_Ugly(t *T) {
	// A cancelled context stops the chain before any hook runs.
	ctx, cancel := WithCancel(Background())
	cancel()
	c := New()
	p := &reloadProbe{}
	AssertTrue(t, c.RegisterService("probe", p).OK)
	AssertFalse(t, c.ServiceReload(ctx).OK)
	AssertEqual(t, 0, p.count)
}

// --- W4-6: optional services degrade instead of aborting boot ---

func TestRuntime_Core_ServiceStartup_Good_OptionalDegrades(t *T) {
	c := New()
	started := false
	AssertTrue(t, c.Service("flaky", Service{
		Optional: true,
		OnStart:  func() Result { return Fail(NewError("telemetry down")) },
	}).OK)
	AssertTrue(t, c.Service("essential", Service{
		OnStart: func() Result { started = true; return Ok(nil) },
	}).OK)
	AssertTrue(t, c.ServiceStartup(Background(), nil).OK)
	AssertTrue(t, started)
}

func TestRuntime_Core_ServiceStartup_Bad_EssentialStillAborts(t *T) {
	c := New()
	AssertTrue(t, c.Service("essential", Service{
		OnStart: func() Result { return Fail(NewError("db down")) },
	}).OK)
	AssertFalse(t, c.ServiceStartup(Background(), nil).OK)
}
