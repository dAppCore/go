// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the command tree in command.go — i18n-key derivation,
// action dispatch, managed-lifecycle check, and registration on Core.
//
// Run:    go test -bench='BenchmarkCommand|BenchmarkCore_Command' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var (
	cmdSinkStr    string
	cmdSinkResult Result
	cmdSinkBool   bool
	cmdSinkSlice  []string
)

func BenchmarkCommand_I18nKey(b *B) {
	cmd := &Command{Path: "deploy/to/homelab"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cmdSinkStr = cmd.I18nKey()
	}
}

func BenchmarkCommand_Run(b *B) {
	cmd := &Command{Action: func(Options) Result { return Result{OK: true} }}
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cmdSinkResult = cmd.Run(opts)
	}
}

func BenchmarkCommand_IsManaged(b *B) {
	cmd := &Command{Managed: "process.daemon"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cmdSinkBool = cmd.IsManaged()
	}
}

func BenchmarkCore_Command(b *B) {
	c := New()
	cmd := Command{Name: "deploy", Path: "deploy", Action: func(Options) Result { return Result{OK: true} }}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cmdSinkResult = c.Command("deploy", cmd)
	}
}

func BenchmarkCore_Commands(b *B) {
	c := New()
	c.Command("deploy", Command{Name: "deploy", Path: "deploy", Action: func(Options) Result { return Result{OK: true} }})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cmdSinkSlice = c.Commands()
	}
}
