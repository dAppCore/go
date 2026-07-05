package core_test

import (
	. "dappco.re/go"
)

// --- Command DTO ---

func TestCommand_Run_Good(t *T) {
	c := New()
	c.Command("greet", Command{Action: func(opts Options) Result {
		return Result{Value: Concat("hello ", opts.String("name")), OK: true}
	}})
	cmd := c.Command("greet").Value.(*Command)
	r := cmd.Run(NewOptions(Option{Key: "name", Value: "world"}))
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello world", r.Value)
}

func TestCommand_Run_NoAction_Good(t *T) {
	c := New()
	c.Command("empty", Command{Description: "no action"})
	cmd := c.Command("empty").Value.(*Command)
	r := cmd.Run(NewOptions())
	AssertFalse(t, r.OK)
}

// --- I18n Key Derivation ---

func TestCommand_I18nKey_Good(t *T) {
	c := New()
	c.Command("deploy/to/homelab", Command{})
	cmd := c.Command("deploy/to/homelab").Value.(*Command)
	AssertEqual(t, "cmd.deploy.to.homelab.description", cmd.I18nKey())
}

func TestCommand_I18nKey_Custom_Good(t *T) {
	c := New()
	c.Command("deploy", Command{Description: "custom.deploy.key"})
	cmd := c.Command("deploy").Value.(*Command)
	AssertEqual(t, "custom.deploy.key", cmd.I18nKey())
}

func TestCommand_I18nKey_Simple_Good(t *T) {
	c := New()
	c.Command("serve", Command{})
	cmd := c.Command("serve").Value.(*Command)
	AssertEqual(t, "cmd.serve.description", cmd.I18nKey())
}

// --- Managed ---

func TestCommand_IsManaged_Good(t *T) {
	c := New()
	c.Command("serve", Command{
		Action:  func(_ Options) Result { return Result{Value: "running", OK: true} },
		Managed: "process.daemon",
	})
	cmd := c.Command("serve").Value.(*Command)
	AssertTrue(t, cmd.IsManaged())
}

func TestCommand_IsManaged_Bad_NotManaged(t *T) {
	c := New()
	c.Command("deploy", Command{
		Action: func(_ Options) Result { return Result{OK: true} },
	})
	cmd := c.Command("deploy").Value.(*Command)
	AssertFalse(t, cmd.IsManaged())
}

// --- Cli Run with Managed ---

func TestCli_Run_Managed_Good(t *T) {
	c := New(WithCli())
	ran := false
	c.Command("serve", Command{
		Action:  func(_ Options) Result { ran = true; return Result{OK: true} },
		Managed: "process.daemon",
	})
	r := c.Cli().Run("serve")
	AssertTrue(t, r.OK)
	AssertTrue(t, ran)
}

func TestCli_Run_NoAction_Bad(t *T) {
	c := New(WithCli())
	c.Command("empty", Command{})
	r := c.Cli().Run("empty")
	AssertFalse(t, r.OK)
}

// --- AX-7 canonical triplets ---

func TestCommand_Command_I18nKey_Good(t *T) {
	cmd := &Command{Path: "deploy/to/homelab"}
	AssertEqual(t, "cmd.deploy.to.homelab.description", cmd.I18nKey())
}

func TestCommand_Command_I18nKey_Bad(t *T) {
	cmd := &Command{Description: "cmd.custom.description", Path: "deploy"}
	AssertEqual(t, "cmd.custom.description", cmd.I18nKey())
}

func TestCommand_Command_I18nKey_Ugly(t *T) {
	cmd := &Command{}
	AssertEqual(t, "cmd..description", cmd.I18nKey())
}

func TestCommand_Command_Run_Good(t *T) {
	cmd := &Command{Path: "agent/run", Action: func(opts Options) Result {
		return Result{Value: opts.String("agent"), OK: true}
	}}
	r := cmd.Run(NewOptions(Option{Key: "agent", Value: "codex"}))
	AssertTrue(t, r.OK)
	AssertEqual(t, "codex", r.Value)
}

func TestCommand_Command_Run_Bad(t *T) {
	cmd := &Command{Path: "agent/run"}
	r := cmd.Run(NewOptions())
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "not executable")
}

func TestCommand_Command_Run_Ugly(t *T) {
	cmd := &Command{Path: "agent/run", Action: func(opts Options) Result {
		return Result{Value: opts.Len(), OK: true}
	}}
	r := cmd.Run(NewOptions())
	AssertTrue(t, r.OK)
	AssertEqual(t, 0, r.Value)
}

func TestCommand_Command_IsManaged_Good(t *T) {
	cmd := &Command{Managed: "process.daemon"}
	AssertTrue(t, cmd.IsManaged())
}

func TestCommand_Command_IsManaged_Bad(t *T) {
	cmd := &Command{}
	AssertFalse(t, cmd.IsManaged())
}

func TestCommand_Command_IsManaged_Ugly(t *T) {
	cmd := &Command{Managed: " "}
	AssertTrue(t, cmd.IsManaged())
}

func TestCommand_Core_Command_Good(t *T) {
	c := New()
	r := c.Command("deploy/to/homelab", Command{Action: func(_ Options) Result {
		return Result{Value: "deployed", OK: true}
	}})
	AssertTrue(t, r.OK)
	AssertTrue(t, c.Command("deploy/to/homelab").OK)
	// Parent paths are auto-created and retrievable.
	AssertTrue(t, c.Command("deploy").OK)
	AssertTrue(t, c.Command("deploy/to").OK)
}

func TestCommand_Core_Command_Bad(t *T) {
	c := New()
	// Malformed paths are rejected.
	AssertFalse(t, c.Command("/deploy", Command{}).OK)         // leading slash
	AssertFalse(t, c.Command("trailing/", Command{}).OK)       // trailing slash
	AssertFalse(t, c.Command("deploy//homelab", Command{}).OK) // double slash
	AssertFalse(t, c.Command("", Command{}).OK)                // empty path
	// Duplicate registration of the same real (action-bearing) path is rejected.
	dup := Command{Action: func(_ Options) Result { return Result{OK: true} }}
	AssertTrue(t, c.Command("once", dup).OK)
	AssertFalse(t, c.Command("once", dup).OK)
	// Retrieving an unregistered path misses.
	AssertFalse(t, c.Command("never/registered").OK)
}

func TestCommand_Core_Command_Ugly(t *T) {
	c := New()
	c.Command("deploy/to/homelab", Command{})
	r := c.Command("deploy", Command{Action: func(_ Options) Result {
		return Result{Value: "parent", OK: true}
	}})
	AssertTrue(t, r.OK)
	parent := c.Command("deploy").Value.(*Command)
	AssertNotNil(t, parent)
	AssertNotEmpty(t, parent.I18nKey())
}

func TestCommand_Core_Commands_Good(t *T) {
	c := New()
	c.Command("agent/prepare", Command{Action: func(_ Options) Result { return Result{OK: true} }})
	c.Command("agent/dispatch", Command{Action: func(_ Options) Result { return Result{OK: true} }})
	commands := c.Commands()
	AssertContains(t, commands, "agent/prepare")
	AssertContains(t, commands, "agent/dispatch")
}

func TestCommand_Core_Commands_Bad(t *T) {
	c := New()
	AssertEmpty(t, c.Commands())
}

func TestCommand_Core_Commands_Ugly(t *T) {
	c := New()
	c.Command("agent/prepare", Command{Action: func(_ Options) Result { return Result{OK: true} }})
	commands := c.Commands()
	commands[0] = "mutated"
	AssertContains(t, c.Commands(), "agent/prepare")
	AssertNotContains(t, c.Commands(), "mutated")
}
