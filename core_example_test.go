package core_test

import . "dappco.re/go"

// ExampleCore_accessors reads the grouped accessor methods through `Core` for Core
// orchestration. Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_accessors() {
	c := New(WithCli(), WithOption("name", "ops"))

	Println(c.Options().String("name"))
	Println(c.App().Name)
	Println(c.Data() != nil)
	Println(c.Drive() != nil)
	Println(c.Fs() != nil)
	Println(c.Config() != nil)
	Println(c.Error() != nil)
	Println(c.Log() != nil)
	Println(c.Cli() != nil)
	Println(c.IPC() != nil)
	Println(c.I18n() != nil)
	// Output:
	// ops
	// ops
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
}

// ExampleCore_Env reads environment access through `Core.Env` for Core orchestration. Core
// keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_Env() {
	c := New()
	Println(c.Env("OS") != "")
	// Output: true
}

// ExampleCore_Context reads the active context through `Core.Context` for Core
// orchestration. Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_Context() {
	c := New()
	Println(c.Context() != nil)
	// Output: true
}

// ExampleCore_Core returns the Core instance through `Core.Core` for Core orchestration.
// Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_Core() {
	c := New()
	Println(c.Core() == c)
	// Output: true
}

// ExampleCore_RunResult runs `Core.RunResult` through the Result-returning startup path
// for Core orchestration. Core keeps orchestration helpers reachable from one
// predictable facade.
func ExampleCore_RunResult() {
	c := New()
	Println(c.RunResult().OK)
	// Output: true
}

// ExampleCore_Run runs `Core.Run` with representative caller inputs for Core
// orchestration. Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_Run() {
	_ = New().Run
}

// ExampleCore_ACTION calls the uppercase action helper through `Core.ACTION` for Core
// orchestration. Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_ACTION() {
	c := New()
	seen := ""
	c.RegisterAction(func(_ *Core, msg Message) Result {
		seen = msg.(string)
		return Result{OK: true}
	})

	c.ACTION("started")
	Println(seen)
	// Output: started
}

// ExampleCore_QUERY calls the uppercase query helper through `Core.QUERY` for Core
// orchestration. Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_QUERY() {
	c := New()
	c.RegisterQuery(func(_ *Core, q Query) Result {
		return Result{Value: Concat("query:", q.(string)), OK: true}
	})

	r := c.QUERY("status")
	Println(r.Value)
	// Output: query:status
}

// ExampleCore_QUERYALL calls every matching uppercase query through `Core.QUERYALL` for
// Core orchestration. Core keeps orchestration helpers reachable from one predictable
// facade.
func ExampleCore_QUERYALL() {
	c := New()
	c.RegisterQuery(func(_ *Core, q Query) Result {
		return Result{Value: Concat("a:", q.(string)), OK: true}
	})
	c.RegisterQuery(func(_ *Core, q Query) Result {
		return Result{Value: Concat("b:", q.(string)), OK: true}
	})

	r := c.QUERYALL("status")
	Println(r.Value)
	// Output: [a:status b:status]
}

// ExampleCore_LogError logs an error through `Core.LogError` for Core orchestration. Core
// keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_LogError() {
	c := New()
	r := c.LogError(nil, "example", "nothing to log")
	Println(r.OK)
	// Output: true
}

// ExampleCore_LogWarn logs a warning through `Core.LogWarn` for Core orchestration. Core
// keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_LogWarn() {
	c := New()
	r := c.LogWarn(nil, "example", "nothing to warn")
	Println(r.OK)
	// Output: true
}

// ExampleCore_Must unwraps a successful Result through `Core.Must` for Core orchestration.
// Core keeps orchestration helpers reachable from one predictable facade.
func ExampleCore_Must() {
	c := New()
	c.Must(nil, "example", "no panic")
	Println("ok")
	// Output: ok
}

// ExampleCore_RegistryOf retrieves a named registry through `Core.RegistryOf` for Core
// orchestration. Core keeps orchestration helpers reachable from one predictable facade.
// ExampleCore_Options returns the key-value options store through `Core.Options`.
func ExampleCore_Options() {
	c := New(WithOption("env", "prod"))
	Println(c.Options() != nil)
	// Output: true
}

// ExampleCore_App returns application metadata through `Core.App`.
func ExampleCore_App() {
	Println(New().App() != nil)
	// Output: true
}

// ExampleCore_Data returns the embedded-asset registry through `Core.Data`.
func ExampleCore_Data() {
	Println(New().Data() != nil)
	// Output: true
}

// ExampleCore_Drive returns the remote-endpoint registry through `Core.Drive`.
func ExampleCore_Drive() {
	Println(New().Drive() != nil)
	// Output: true
}

// ExampleCore_Fs returns the sandboxed filesystem through `Core.Fs`.
func ExampleCore_Fs() {
	Println(New().Fs() != nil)
	// Output: true
}

// ExampleCore_Config returns the configuration subsystem through `Core.Config`.
func ExampleCore_Config() {
	Println(New().Config() != nil)
	// Output: true
}

// ExampleCore_Feature returns a feature-flag handle through `Core.Feature`.
func ExampleCore_Feature() {
	c := New()
	c.Feature("dark-mode").Enable()
	Println(c.Feature("dark-mode").Enabled())
	// Output: true
}

// ExampleCore_Error returns the panic-recovery subsystem through `Core.Error`.
func ExampleCore_Error() {
	Println(New().Error() != nil)
	// Output: true
}

// ExampleCore_Log returns the structured logger through `Core.Log`.
func ExampleCore_Log() {
	Println(New().Log() != nil)
	// Output: true
}

// ExampleCore_Cli returns the CLI subsystem through `Core.Cli`.
func ExampleCore_Cli() {
	Println(New(WithCli()).Cli() != nil)
	// Output: true
}

// ExampleCore_IPC returns the IPC action/task subsystem through `Core.IPC`.
func ExampleCore_IPC() {
	Println(New().IPC() != nil)
	// Output: true
}

// ExampleCore_I18n returns the translation subsystem through `Core.I18n`.
func ExampleCore_I18n() {
	Println(New().I18n() != nil)
	// Output: true
}

// ExampleCore_WithContext derives a Core bound to a context through `Core.WithContext`.
func ExampleCore_WithContext() {
	Println(New().WithContext(Background()) != nil)
	// Output: true
}

func ExampleCore_RegistryOf() {
	c := New()
	c.Action("agent.deploy", func(_ Context, _ Options) Result { return Result{OK: true} })
	r := c.RegistryOf("actions")
	Println(r.OK)
	Println(r.Value.(*Registry[any]).Names())
	Println(c.RegistryOf("missing").OK)
	// Output:
	// true
	// [agent.deploy]
	// false
}
