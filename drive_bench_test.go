// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Drive primitive in drive.go.
// Per AX-11 — Drive is the transport handle registry. Every consumer
// that talks to "api" / "ssh" / "mcp" goes through Drive().Get(name)
// to resolve the handle, then invokes the transport. Lookups happen
// per-call; registration happens at boot.
//
// Run:    go test -bench='BenchmarkDrive' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	driveSinkResult Result
)

// driveFixture returns a Core with one transport handle registered.
func driveFixture() *Core {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "bench.api"},
		Option{Key: "transport", Value: "https://bench.lthn.ai"},
	))
	return c
}

// --- New / Get ---

func BenchmarkDrive_New(b *B) {
	c := New()
	opts := NewOptions(
		Option{Key: "name", Value: "bench.api"},
		Option{Key: "transport", Value: "https://bench.lthn.ai"},
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		driveSinkResult = c.Drive().New(opts)
	}
}

func BenchmarkDrive_Get_Hit(b *B) {
	c := driveFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		driveSinkResult = c.Drive().Get("bench.api")
	}
}

func BenchmarkDrive_Get_Miss(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		driveSinkResult = c.Drive().Get("noexist")
	}
}
