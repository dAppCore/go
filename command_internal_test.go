// SPDX-License-Identifier: EUPL-1.2

package core

func TestCommand_pathName_Good(t *T) {
	AssertEqual(t, "homelab", pathName("deploy/to/homelab"))
}
func TestCommand_pathName_Bad(t *T) {
	AssertEqual(t, "", pathName(""))
}
func TestCommand_pathName_Ugly(t *T) {
	AssertEqual(t, "", pathName("deploy/"))
}

// TestCommand_Command_DisabledSkipsDispatch_Good proves F1: a registry-disabled
// command resolves as absent through c.Command and does not fire on dispatch.
func TestCommand_Command_DisabledSkipsDispatch_Good(t *T) {
	c := New(WithCli())
	ran := false
	c.Command("agent/run", Command{Action: func(_ Options) Result { ran = true; return Result{OK: true} }})

	// Enabled: resolves and dispatches.
	AssertTrue(t, c.Command("agent/run").OK)

	// Soft-disable at the registry level.
	AssertTrue(t, c.commands.Disable("agent/run").OK)

	// Now absent to resolution...
	AssertFalse(t, c.Command("agent/run").OK)
	// ...and dispatch does not fire the action.
	r := c.Cli().Run("agent", "run")
	AssertFalse(t, r.OK)
	AssertFalse(t, ran)
}
