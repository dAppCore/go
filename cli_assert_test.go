// SPDX-License-Identifier: EUPL-1.2

package core

// cliStubCore returns a Core with a stub "process.run" action so AssertCLI
// can be exercised without an external process service: a command named
// "fail" reports !OK, anything else succeeds with "ran: <command>" output.
func cliStubCore() *Core {
	c := New()
	c.Action("process.run", func(_ Context, opts Options) Result {
		cmd := opts.String("command")
		if cmd == "fail" {
			return Result{Value: NewError("command failed"), OK: false}
		}
		return Result{Value: "ran: " + cmd, OK: true}
	})
	return c
}

func TestCliAssert_AssertCLI_Good(t *T) {
	c := cliStubCore()
	AssertCLI(t, c, CLITest{Cmd: "echo", WantOK: true, Contains: "echo"})
}

func TestCliAssert_AssertCLI_Bad(t *T) {
	// WantOK:true but the command fails → AssertCLI records the mismatch.
	c := cliStubCore()
	st := assertStub(t)
	AssertCLI(st, c, CLITest{Cmd: "fail", WantOK: true})
	assertOneMessage(t, st, "AssertCLI")
}

func TestCliAssert_AssertCLI_Ugly(t *T) {
	// OK matches but the required substring is absent → AssertCLI records it.
	c := cliStubCore()
	st := assertStub(t)
	AssertCLI(st, c, CLITest{Cmd: "echo", WantOK: true, Contains: "absent-substring"})
	assertOneMessage(t, st, "stdout-contains")
}

func TestCliAssert_AssertCLIs_Good(t *T) {
	c := cliStubCore()
	AssertCLIs(t, c, []CLITest{
		{Name: "echo", Cmd: "echo", WantOK: true, Contains: "echo"},
		{Name: "list", Cmd: "ls", WantOK: true},
	})
}

func TestCliAssert_AssertCLIs_Bad(t *T) {
	// A command expected to fail (WantOK:false) is handled as a pass.
	c := cliStubCore()
	AssertCLIs(t, c, []CLITest{
		{Name: "expected-failure", Cmd: "fail", WantOK: false},
	})
}

func TestCliAssert_AssertCLIs_Ugly(t *T) {
	// An empty table runs no sub-tests and passes trivially.
	AssertCLIs(t, cliStubCore(), nil)
}
