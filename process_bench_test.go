// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the process façade in process.go. Run/RunIn/RunWithEnv/
// Start/Kill dispatch to the process.* actions; on a bare Core no handler
// is registered (the executor ships as dappco.re/go-process), so these
// measure the dispatch + Result path without spawning a subprocess.
//
// Run:    go test -bench='BenchmarkProcess' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var (
	procSink     Result
	procSinkP    *Process
	procSinkBool bool
)

func BenchmarkProcess(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSinkP = c.Process()
	}
}

func BenchmarkProcess_Run(b *B) {
	p := New().Process()
	ctx := Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSink = p.Run(ctx, "true")
	}
}

func BenchmarkProcess_RunIn(b *B) {
	p := New().Process()
	ctx := Background()
	dir := TempDir()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSink = p.RunIn(ctx, dir, "true")
	}
}

func BenchmarkProcess_RunWithEnv(b *B) {
	p := New().Process()
	ctx := Background()
	dir := TempDir()
	env := Environ()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSink = p.RunWithEnv(ctx, dir, env, "true")
	}
}

func BenchmarkProcess_Start(b *B) {
	p := New().Process()
	ctx := Background()
	opts := NewOptions(Option{Key: "command", Value: "true"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSink = p.Start(ctx, opts)
	}
}

func BenchmarkProcess_Kill(b *B) {
	p := New().Process()
	ctx := Background()
	opts := NewOptions(Option{Key: "id", Value: "bench"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSink = p.Kill(ctx, opts)
	}
}

func BenchmarkProcess_Exists(b *B) {
	p := New().Process()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		procSinkBool = p.Exists()
	}
}
