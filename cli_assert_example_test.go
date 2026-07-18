// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

// ExampleAssertCLI_stubProcess dispatches a CLITest case through `AssertCLI` for AX-10
// binary validation. Like assert_example_test.go's helpers, this is compiled but not
// executed (no "Output:" comment): AssertCLI takes a live TB, which an Example function
// signature can never supply. cli_example_test.go already declares `ExampleAssertCLI`
// (attached to cli.go, documenting the same dispatch by calling c.Process() directly) —
// this file's file-aware requirement is for the cli_assert.go symbol itself, so this
// variant calls the real AssertCLI helper the symbol names instead of redeclaring an
// identical duplicate under a different file.
func ExampleAssertCLI_stubProcess() {
	var t *T
	c := New()
	c.Action("process.run", func(ctx Context, opts Options) Result {
		return Result{Value: "go version go1.26.0\n", OK: true}
	})

	AssertCLI(t, c, CLITest{
		Cmd: "go", Args: []string{"version"},
		WantOK: true, Contains: "go1.26",
	})
}

// ExampleAssertCLIs runs a table of CLI scenarios through `AssertCLIs` for AX-10 binary
// validation. Compiled but not executed (no "Output:" comment): AssertCLIs takes a live
// *T, which an Example function signature can never supply.
func ExampleAssertCLIs() {
	var t *T
	c := New()
	c.Action("process.run", func(ctx Context, opts Options) Result {
		return Result{Value: "go version go1.26.0\n", OK: true}
	})

	AssertCLIs(t, c, []CLITest{
		{Name: "version", Cmd: "go", Args: []string{"version"}, WantOK: true, Contains: "go1."},
		{Name: "vet", Cmd: "go", Args: []string{"vet", "./..."}, WantOK: true},
	})
}
