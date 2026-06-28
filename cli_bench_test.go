// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the CLI surface in cli.go — banner rendering, output
// redirection, formatted printing, and command dispatch. Output is sent
// to Discard so the numbers measure the CLI machinery, not the terminal.
//
// Run:    go test -bench='BenchmarkCli' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var (
	cliSink    Result
	cliSinkStr string
)

func BenchmarkCliRegister(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cliSink = CliRegister(New(WithCli()))
	}
}

func BenchmarkCli_Print(b *B) {
	cli := New(WithCli()).Cli()
	cli.SetOutput(Discard)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cli.Print("port=%d", 8080)
	}
}

func BenchmarkCli_SetOutput(b *B) {
	cli := New(WithCli()).Cli()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cli.SetOutput(Discard)
	}
}

func BenchmarkCli_SetBanner(b *B) {
	cli := New(WithCli()).Cli()
	fn := func(*Cli) string { return "Core" }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cli.SetBanner(fn)
	}
}

func BenchmarkCli_Banner(b *B) {
	cli := New(WithCli()).Cli()
	cli.SetBanner(func(*Cli) string { return "Core" })
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cliSinkStr = cli.Banner()
	}
}

func BenchmarkCli_PrintHelp(b *B) {
	cli := New(WithCli()).Cli()
	cli.SetOutput(Discard)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cli.PrintHelp()
	}
}

func BenchmarkCli_Run(b *B) {
	cli := New(WithCli()).Cli()
	cli.SetOutput(Discard)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cliSink = cli.Run("status")
	}
}
