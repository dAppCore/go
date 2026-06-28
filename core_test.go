package core_test

import (
	. "dappco.re/go"
)

// --- New ---

func TestCore_New_Good(t *T) {
	c := New()
	AssertNotNil(t, c)
}

func TestCore_New_WithOptions_Good(t *T) {
	c := New(WithOptions(NewOptions(Option{Key: "name", Value: "myapp"})))
	AssertNotNil(t, c)
	AssertEqual(t, "myapp", c.App().Name)
}

func TestCore_New_WithOptions_Bad(t *T) {
	// Empty options — should still create a valid Core
	c := New(WithOptions(NewOptions()))
	AssertNotNil(t, c)
}

func TestCore_New_WithService_Good(t *T) {
	started := false
	c := New(
		WithOptions(NewOptions(Option{Key: "name", Value: "myapp"})),
		WithService(func(c *Core) Result {
			c.Service("test", Service{
				OnStart: func() Result { started = true; return Result{OK: true} },
			})
			return Result{OK: true}
		}),
	)

	svc := c.Service("test")
	AssertTrue(t, svc.OK)

	c.ServiceStartup(Background(), nil)
	AssertTrue(t, started)
}

func TestCore_New_WithServiceLock_Good(t *T) {
	c := New(
		WithService(func(c *Core) Result {
			c.Service("allowed", Service{})
			return Result{OK: true}
		}),
		WithServiceLock(),
	)

	// Registration after lock should fail
	reg := c.Service("blocked", Service{})
	AssertFalse(t, reg.OK)
}

func TestCore_New_WithService_Bad_FailingOption(t *T) {
	secondCalled := false
	_ = New(
		WithService(func(c *Core) Result {
			return Result{Value: E("test", "intentional failure", nil), OK: false}
		}),
		WithService(func(c *Core) Result {
			secondCalled = true
			return Result{OK: true}
		}),
	)
	AssertFalse(t, secondCalled, "second option should not run after first fails")
}

// --- Core Error/Log Helpers ---

func TestCore_LogError_Good(t *T) {
	c := New()
	cause := AnError
	r := c.LogError(cause, "test.Operation", "something broke")

	err, ok := r.Value.(error)
	AssertTrue(t, ok)
	AssertErrorIs(t, err, cause)
}

func TestCore_LogWarn_Good(t *T) {
	c := New()
	r := c.LogWarn(AnError, "test.Operation", "heads up")

	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

func TestCore_Must_Ugly(t *T) {
	c := New()
	AssertPanics(t, func() {
		c.Must(AnError, "test.Operation", "fatal")
	})
}

func TestCore_Must_Nil_Good(t *T) {
	c := New()
	AssertNotPanics(t, func() {
		c.Must(nil, "test.Operation", "no error")
	})
}

// --- RegistryOf ---

func TestCore_RegistryOf_Good_Services(t *T) {
	c := New(
		WithCli(),
		WithService(func(c *Core) Result {
			return c.Service("alpha", Service{})
		}),
		WithService(func(c *Core) Result {
			return c.Service("bravo", Service{})
		}),
	)
	reg := c.RegistryOf("services").Value.(*Registry[any])
	// cli (via WithCli) + our 2
	AssertTrue(t, reg.Has("alpha"))
	AssertTrue(t, reg.Has("bravo"))
	AssertTrue(t, reg.Has("cli"))
}

func TestCore_RegistryOf_Good_Commands(t *T) {
	c := New()
	c.Command("deploy", Command{Action: func(_ Options) Result { return Result{OK: true} }})
	c.Command("test", Command{Action: func(_ Options) Result { return Result{OK: true} }})

	reg := c.RegistryOf("commands").Value.(*Registry[any])
	AssertTrue(t, reg.Has("deploy"))
	AssertTrue(t, reg.Has("test"))
}

func TestCore_RegistryOf_Good_Actions(t *T) {
	c := New()
	c.Action("process.run", func(_ Context, _ Options) Result { return Result{OK: true} })
	c.Action("brain.recall", func(_ Context, _ Options) Result { return Result{OK: true} })

	reg := c.RegistryOf("actions").Value.(*Registry[any])
	AssertTrue(t, reg.Has("process.run"))
	AssertTrue(t, reg.Has("brain.recall"))
	AssertEqual(t, 2, reg.Len())
}

func TestCore_RegistryOf_Bad_Unknown(t *T) {
	c := New()
	r := c.RegistryOf("nonexistent")
	AssertFalse(t, r.OK, "unknown registry returns OK=false")
}

// --- RunResult ---

func TestCore_RunResult_Good(t *T) {
	c := New(
		WithService(func(c *Core) Result {
			return c.Service("healthy", Service{
				OnStart: func() Result { return Result{OK: true} },
				OnStop:  func() Result { return Result{OK: true} },
			})
		}),
	)
	r := c.RunResult()
	AssertTrue(t, r.OK)
}

func TestCore_RunResult_Bad(t *T) {
	c := New(
		WithService(func(c *Core) Result {
			return c.Service("broken", Service{
				OnStart: func() Result {
					return Result{Value: NewError("startup failed"), OK: false}
				},
			})
		}),
	)
	r := c.RunResult()
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "startup failed")
}

func TestCore_RunResult_Ugly(t *T) {
	shutdownCalled := false
	c := New(
		WithService(func(c *Core) Result {
			return c.Service("cleanup", Service{
				OnStart: func() Result { return Result{OK: true} },
				OnStop:  func() Result { shutdownCalled = true; return Result{OK: true} },
			})
		}),
		WithService(func(c *Core) Result {
			return c.Service("broken", Service{
				OnStart: func() Result {
					return Result{Value: NewError("boom"), OK: false}
				},
			})
		}),
	)
	r := c.RunResult()
	AssertFalse(t, r.OK)
	AssertTrue(t, shutdownCalled, "ServiceShutdown must be called even when startup fails — cleanup service must get OnStop")
}

// Run() delegates to RunResult() — tested via RunResult tests above.
// c.Exit(1) behaviour is verified by RunResult returning OK=false correctly.

func TestCore_Core_ACTION_Good(t *T) {
	c := New()
	var got Message
	c.RegisterAction(func(_ *Core, msg Message) Result {
		got = msg
		return Result{OK: true}
	})

	r := c.ACTION("agent.dispatch")

	AssertTrue(t, r.OK)
	AssertEqual(t, Message("agent.dispatch"), got)
}

func TestCore_Core_ACTION_Bad(t *T) {
	c := New()
	c.RegisterAction(func(_ *Core, _ Message) Result {
		return Result{Value: NewError("handler refused"), OK: false}
	})

	r := c.ACTION("agent.dispatch")

	AssertTrue(t, r.OK)
}

func TestCore_Core_ACTION_Ugly(t *T) {
	c := New()
	c.RegisterAction(func(_ *Core, _ Message) Result {
		panic("handler panic")
	})

	r := c.ACTION(nil)

	AssertTrue(t, r.OK)
}

func TestCore_Core_App_Good(t *T) {
	c := New(WithOption("name", "homelab"))
	AssertEqual(t, "homelab", c.App().Name)
}

func TestCore_Core_App_Bad(t *T) {
	c := New()
	AssertEqual(t, "", c.App().Name)
}

func TestCore_Core_App_Ugly(t *T) {
	c := New()
	AssertSame(t, c.App(), c.App())
}

func TestCore_Core_Cli_Good(t *T) {
	c := New(WithCli())
	AssertNotNil(t, c.Cli())
}

func TestCore_Core_Cli_Bad(t *T) {
	c := New(WithCli(), WithServiceLock())
	AssertNotNil(t, c.Cli())
}

func TestCore_Core_Cli_Ugly(t *T) {
	c := New(WithCli())
	c.Cli().SetOutput(NewBuffer())
	AssertNotPanics(t, func() {
		_ = c.Cli().Run()
	})
}

func TestCore_Core_Config_Good(t *T) {
	c := New()
	c.Config().Set("agent.region", "lab")
	AssertEqual(t, "lab", c.Config().String("agent.region"))
}

func TestCore_Core_Config_Bad(t *T) {
	c := New()
	AssertEqual(t, "", c.Config().String("missing"))
}

func TestCore_Core_Config_Ugly(t *T) {
	c := New()
	AssertSame(t, c.Config(), c.Config())
}

func TestCore_Core_Context_Good(t *T) {
	c := New()
	AssertNotNil(t, c.Context())
}

func TestCore_Core_Context_Bad(t *T) {
	c := New()
	c.ServiceShutdown(Background())
	AssertError(t, c.Context().Err())
}

func TestCore_Core_Context_Ugly(t *T) {
	c := New()
	before := c.Context()
	c.ServiceStartup(Background(), nil)
	AssertNotNil(t, before)
	AssertNotNil(t, c.Context())
}

func TestCore_Core_Core_Good(t *T) {
	c := New()
	AssertSame(t, c, c.Core())
}

func TestCore_Core_Core_Bad(t *T) {
	c := New(WithOption("name", "agent"))
	AssertSame(t, c, c.Core())
}

func TestCore_Core_Core_Ugly(t *T) {
	var c *Core
	AssertNil(t, c.Core())
}

func TestCore_Core_Data_Good(t *T) {
	c := New()
	AssertNotNil(t, c.Data())
}

func TestCore_Core_Data_Bad(t *T) {
	c := New()
	AssertSame(t, c.Data(), c.Data())
}

func TestCore_Core_Data_Ugly(t *T) {
	c := New()
	AssertEqual(t, 0, c.Data().Len())
}

func TestCore_Core_Drive_Good(t *T) {
	c := New()
	AssertNotNil(t, c.Drive())
}

func TestCore_Core_Drive_Bad(t *T) {
	c := New()
	AssertSame(t, c.Drive(), c.Drive())
}

func TestCore_Core_Drive_Ugly(t *T) {
	c := New()
	AssertEqual(t, 0, c.Drive().Len())
}

func TestCore_Core_Env_Good(t *T) {
	c := New()
	AssertEqual(t, OS(), c.Env("OS"))
}

func TestCore_Core_Env_Bad(t *T) {
	c := New()
	AssertEqual(t, "", c.Env("CORE_TEST_MISSING"))
}

func TestCore_Core_Env_Ugly(t *T) {
	t.Setenv("CORE_TEST_SESSION", "token")
	c := New()
	AssertEqual(t, "token", c.Env("CORE_TEST_SESSION"))
}

func TestCore_Core_Feature_Good(t *T) {
	c := New()
	f := c.Feature("dark-mode")
	AssertEqual(t, "dark-mode", f.Name())
	f.Enable()
	AssertTrue(t, c.Feature("dark-mode").Enabled())
}

func TestCore_Core_Feature_Bad(t *T) {
	c := New()
	// An empty name yields a usable handle that is never enabled.
	f := c.Feature("")
	AssertEqual(t, "", f.Name())
	AssertFalse(t, f.Enabled())
}

func TestCore_Core_Feature_Ugly(t *T) {
	c := New()
	// Separate handles to the same name share one backing store.
	c.Feature("beta").Enable()
	AssertTrue(t, c.Feature("beta").Enabled())
	c.Feature("beta").Disable()
	AssertFalse(t, c.Feature("beta").Enabled())
}

func TestCore_Core_Error_Good(t *T) {
	c := New()
	AssertNotNil(t, c.Error())
}

func TestCore_Core_Error_Bad(t *T) {
	c := New()
	AssertNotPanics(t, func() {
		c.Error().Recover()
	})
}

func TestCore_Core_Error_Ugly(t *T) {
	c := New()
	AssertSame(t, c.Error(), c.Error())
}

func TestCore_Core_Fs_Good(t *T) {
	c := New()
	AssertNotNil(t, c.Fs())
}

func TestCore_Core_Fs_Bad(t *T) {
	c := New()
	AssertFalse(t, c.Fs().Read(Path(t.TempDir(), "missing")).OK)
}

func TestCore_Core_Fs_Ugly(t *T) {
	c := New()
	AssertSame(t, c.Fs(), c.Fs())
}

func TestCore_Core_I18n_Good(t *T) {
	c := New()
	AssertNotNil(t, c.I18n())
}

func TestCore_Core_I18n_Bad(t *T) {
	c := New()
	AssertSame(t, c.I18n(), c.I18n())
}

func TestCore_Core_I18n_Ugly(t *T) {
	c := New()
	r := c.I18n().Translate("missing.key")
	AssertTrue(t, r.OK)
	AssertEqual(t, "missing.key", r.Value)
}

func TestCore_Core_IPC_Good(t *T) {
	c := New()
	AssertNotNil(t, c.IPC())
}

func TestCore_Core_IPC_Bad(t *T) {
	c := New()
	AssertSame(t, c.IPC(), c.IPC())
}

func TestCore_Core_IPC_Ugly(t *T) {
	c := New()
	// Broadcasting through the IPC handle with no handlers is a no-op success.
	AssertNotNil(t, c.IPC())
	AssertTrue(t, c.ACTION(nil).OK)
}

func TestCore_Core_Log_Good(t *T) {
	c := New()
	AssertNotNil(t, c.Log())
}

func TestCore_Core_Log_Bad(t *T) {
	c := New()
	AssertTrue(t, c.Log().Error(nil, "agent.Dispatch", "ok").OK)
}

func TestCore_Core_Log_Ugly(t *T) {
	c := New()
	AssertSame(t, c.Log(), c.Log())
}

func TestCore_Core_LogError_Good(t *T) {
	r := New().LogError(AnError, "agent.Dispatch", "failed")
	AssertFalse(t, r.OK)
	AssertErrorIs(t, r.Value.(error), AnError)
}

func TestCore_Core_LogError_Bad(t *T) {
	r := New().LogError(nil, "agent.Dispatch", "ok")
	AssertTrue(t, r.OK)
}

func TestCore_Core_LogError_Ugly(t *T) {
	r := New().LogError(NewCode("agent.refused", "dispatch refused"), "agent.Dispatch", "failed")
	AssertFalse(t, r.OK)
	AssertEqual(t, "agent.refused", ErrorCode(r.Value.(error)))
}

func TestCore_Core_LogWarn_Good(t *T) {
	r := New().LogWarn(AnError, "agent.Dispatch", "degraded")
	AssertFalse(t, r.OK)
	AssertErrorIs(t, r.Value.(error), AnError)
}

func TestCore_Core_LogWarn_Bad(t *T) {
	r := New().LogWarn(nil, "agent.Dispatch", "ok")
	AssertTrue(t, r.OK)
}

func TestCore_Core_LogWarn_Ugly(t *T) {
	r := New().LogWarn(NewCode("agent.degraded", "dispatch degraded"), "agent.Dispatch", "degraded")
	AssertFalse(t, r.OK)
	AssertEqual(t, "agent.degraded", ErrorCode(r.Value.(error)))
}

func TestCore_Core_Must_Good(t *T) {
	AssertNotPanics(t, func() {
		New().Must(nil, "agent.Dispatch", "ok")
	})
}

func TestCore_Core_Must_Bad(t *T) {
	AssertPanicsWithError(t, "dispatch failed", func() {
		New().Must(AnError, "agent.Dispatch", "dispatch failed")
	})
}

func TestCore_Core_Must_Ugly(t *T) {
	AssertPanicsWithError(t, "agent.refused", func() {
		New().Must(NewCode("agent.refused", "dispatch refused"), "agent.Dispatch", "failed")
	})
}

func TestCore_Core_Options_Good(t *T) {
	c := New(WithOptions(NewOptions(
		Option{Key: "name", Value: "agent"},
		Option{Key: "port", Value: 8080},
		Option{Key: "debug", Value: true},
	)))
	AssertEqual(t, "agent", c.Options().String("name"))
	AssertEqual(t, 8080, c.Options().Int("port"))
	AssertTrue(t, c.Options().Bool("debug"))
}

func TestCore_Core_Options_Bad(t *T) {
	c := New()
	AssertNil(t, c.Options())
}

func TestCore_Core_Options_Ugly(t *T) {
	c := New(WithOptions(NewOptions()))
	AssertEqual(t, 0, c.Options().Len())
}

func TestCore_Core_QUERY_Good(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, q Query) Result {
		return Result{Value: Sprintf("ack:%v", q), OK: true}
	})

	r := c.QUERY("agent.status")

	AssertTrue(t, r.OK)
	AssertEqual(t, "ack:agent.status", r.Value)
}

func TestCore_Core_QUERY_Bad(t *T) {
	r := New().QUERY("missing")
	AssertFalse(t, r.OK)
}

func TestCore_Core_QUERY_Ugly(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result {
		return Result{Value: nil, OK: true}
	})

	r := c.QUERY(nil)

	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
}

func TestCore_Core_QUERYALL_Good(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "agent", OK: true} })
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "health", OK: true} })

	r := c.QUERYALL("status")

	AssertTrue(t, r.OK)
	AssertElementsMatch(t, []any{"agent", "health"}, r.Value.([]any))
}

func TestCore_Core_QUERYALL_Bad(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "ignored", OK: false} })

	r := c.QUERYALL("status")

	AssertTrue(t, r.OK)
	AssertEmpty(t, r.Value.([]any))
}

func TestCore_Core_QUERYALL_Ugly(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: nil, OK: true} })

	r := c.QUERYALL(nil)

	AssertTrue(t, r.OK)
	AssertEmpty(t, r.Value.([]any))
}

func TestCore_Core_RegistryOf_Good(t *T) {
	c := New()
	c.Action("agent.dispatch", func(_ Context, _ Options) Result { return Result{OK: true} })

	reg := c.RegistryOf("actions").Value.(*Registry[any])

	AssertTrue(t, reg.Has("agent.dispatch"))
}

func TestCore_Core_RegistryOf_Bad(t *T) {
	AssertFalse(t, New().RegistryOf("missing").OK)
}

func TestCore_Core_RegistryOf_Ugly(t *T) {
	c := New()
	// Snapshot taken before the action is registered does not see it.
	reg := c.RegistryOf("actions").Value.(*Registry[any])
	c.Action("agent.dispatch", func(_ Context, _ Options) Result { return Result{OK: true} })

	AssertFalse(t, reg.Has("agent.dispatch"))
}

func TestCore_Core_Run_Good(t *T) {
	c := New(WithCli())
	c.Cli().SetOutput(NewBuffer())
	AssertNotPanics(t, func() {
		c.Run()
	})
}

func TestCore_Core_Run_Bad(t *T) {
	c := New(WithCli(), WithService(func(c *Core) Result {
		return c.Service("agent", Service{OnStart: func() Result { return Result{OK: true} }})
	}))
	c.Cli().SetOutput(NewBuffer())
	AssertNotPanics(t, func() {
		c.Run()
	})
}

func TestCore_Core_Run_Ugly(t *T) {
	// A core carrying a registered command still runs without panic.
	c := New(WithCli())
	c.Command("noop", Command{Action: func(_ Options) Result { return Result{OK: true} }})
	c.Cli().SetOutput(NewBuffer())
	AssertNotPanics(t, func() {
		c.Run()
	})
}

func TestCore_Core_RunResult_Good(t *T) {
	c := New(WithCli())
	c.Cli().SetOutput(NewBuffer())
	AssertTrue(t, c.RunResult().OK)
}

func TestCore_Core_RunResult_Bad(t *T) {
	c := New(WithService(func(c *Core) Result {
		return c.Service("agent", Service{
			OnStart: func() Result { return Result{Value: NewError("startup refused"), OK: false} },
		})
	}))

	r := c.RunResult()

	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "startup refused")
}

func TestCore_Core_RunResult_Ugly(t *T) {
	stopped := false
	c := New(
		WithService(func(c *Core) Result {
			return c.Service("cleanup", Service{
				OnStart: func() Result { return Result{OK: true} },
				OnStop:  func() Result { stopped = true; return Result{OK: true} },
			})
		}),
		WithService(func(c *Core) Result {
			return c.Service("broken", Service{
				OnStart: func() Result { return Result{Value: NewError("boom"), OK: false} },
			})
		}),
	)

	r := c.RunResult()

	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "boom")
	AssertTrue(t, stopped)
}

// --- WithContext ---

type withContextKey struct{}

// Good — the request-scoped value round-trips on the clone, the parent's
// Context is left untouched, and the heavy subsystems are shared by
// pointer (here: the same *Config).
func TestCore_WithContext_Good(t *T) {
	c := New()
	rc := c.WithContext(WithValue(c.Context(), withContextKey{}, "user-42"))

	AssertEqual(t, "user-42", rc.Context().Value(withContextKey{}))
	AssertNil(t, c.Context().Value(withContextKey{})) // parent not mutated
	AssertSame(t, c.Config(), rc.Config())            // shared subsystem
}

// Bad — cancelling the derived Core must NOT cancel the parent: the
// parent's lifecycle context stays live (Err() == nil) afterwards.
func TestCore_WithContext_Bad(t *T) {
	c := New()
	rc := c.WithContext(c.Context())

	rc.ServiceShutdown(Background()) // fires the derived cancel

	AssertNotNil(t, rc.Context().Err()) // derived is cancelled
	AssertNil(t, c.Context().Err())     // parent is not
}

// Ugly — a parent shutdown propagates DOWN to the derived context when
// the derived ctx chains from the parent (the documented usual case).
func TestCore_WithContext_Ugly(t *T) {
	c := New()
	rc := c.WithContext(WithValue(c.Context(), withContextKey{}, "scoped"))

	c.ServiceShutdown(Background()) // cancels the parent's context

	AssertNotNil(t, rc.Context().Err()) // propagated to the derived ctx
}
