package core_test

import (
	. "dappco.re/go"
)

// ExampleAssertCLI demonstrates the canonical CLI test shape: build a
// Core, register a process.run handler (typically dappco.re/go-process
// in production code, an in-test fake here), then dispatch a CLITest
// case that asserts on stdout substring and OK status.
func ExampleAssertCLI() {
	c := New()
	c.Action("process.run", func(ctx Context, opts Options) Result {
		return Result{Value: "go version go1.26.0\n", OK: true}
	})

	tc := CLITest{
		Cmd:      "go",
		Args:     []string{"version"},
		WantOK:   true,
		Contains: "go1.26",
	}
	r := c.Process().Run(Background(), tc.Cmd, tc.Args...)
	Println(r.OK)
	Println(Contains(r.Value.(string), tc.Contains))
	// Output:
	// true
	// true
}

// ExampleCliRegister wires the CLI action bridge into a Core through `CliRegister`.
func ExampleCliRegister() {
	_ = CliRegister(New(WithCli())) // r.OK once the CLI commands are registered
}

// ExampleCli_SetOutput redirects CLI output to a writer through `Cli.SetOutput`.
func ExampleCli_SetOutput() {
	cli := New(WithCli()).Cli()
	buf := NewBuffer()
	cli.SetOutput(buf)
	cli.Print("ready")
	Println(buf.String())
	// Output: ready
}

// ExampleCli_Print writes formatted text to the CLI output through `Cli.Print`.
func ExampleCli_Print() {
	cli := New(WithCli()).Cli()
	buf := NewBuffer()
	cli.SetOutput(buf)
	cli.Print("port=%d", 8080)
	Println(buf.String())
	// Output: port=8080
}

// ExampleCli_SetBanner sets the CLI banner renderer through `Cli.SetBanner`.
func ExampleCli_SetBanner() {
	cli := New(WithCli()).Cli()
	cli.SetBanner(func(*Cli) string { return "Core v1" })
	Println(cli.Banner())
	// Output: Core v1
}

// ExampleCli_Banner returns the rendered CLI banner through `Cli.Banner`.
func ExampleCli_Banner() {
	cli := New(WithCli()).Cli()
	cli.SetBanner(func(*Cli) string { return "Lethean" })
	Println(cli.Banner())
	// Output: Lethean
}

// ExampleCli_Run dispatches CLI arguments to registered commands through `Cli.Run`.
func ExampleCli_Run() {
	cli := New(WithCli()).Cli()
	cli.SetOutput(NewBuffer())
	if false {
		cli.Run("version") // dispatches the named command
	}
}

// ExampleCli_PrintHelp writes usage text to the CLI output through `Cli.PrintHelp`.
func ExampleCli_PrintHelp() {
	cli := New(WithCli()).Cli()
	cli.SetOutput(NewBuffer())
	cli.PrintHelp()
}
