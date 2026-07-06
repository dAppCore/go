// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the App primitive in app.go.
// Per AX-11 — App is a small DTO carrying app identity. App.New() is
// called once during c.New() to materialise it from Options; App.Find()
// is used by consumers to locate an executable on PATH. Both are
// boot-time-ish but worth gating against regressions.
//
// Run:    go test -bench='BenchmarkApp' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	appSinkApp    App
	appSinkResult Result
)

// --- App.New ---

func BenchmarkApp_New_Empty(b *B) {
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		appSinkApp = App{}.New(opts)
	}
}

func BenchmarkApp_New_AllFields(b *B) {
	opts := NewOptions(
		Option{Key: "name", Value: "Cladius"},
		Option{Key: "version", Value: "1.0.0"},
		Option{Key: "description", Value: "agentic team leader"},
		Option{Key: "filename", Value: "cladius"},
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		appSinkApp = App{}.New(opts)
	}
}

// --- App.Find ---
//
// `sh` is on /bin on every dev box and POSIX-mandated on PATH. It
// makes a deterministic Find target that exercises the syscalls
// without depending on a build artifact.

func BenchmarkApp_Find_OnPATH(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		appSinkResult = App{}.Find("sh", "Shell")
	}
}

func BenchmarkApp_Find_NotFound(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		appSinkResult = App{}.Find("__core_no_such_binary__", "Missing")
	}
}

func BenchmarkApp_Find_AbsolutePath(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		appSinkResult = App{}.Find("/bin/sh", "Shell")
	}
}
