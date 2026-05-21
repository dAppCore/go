// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the runtime primitives in runtime.go.
// Per AX-11 — runtime.go covers the service-lifecycle + tracked-
// goroutine surface that every consumer hits on startup, on every
// background worker spawn, and on every shutdown-check inside a
// long-running loop:
//
//   * ServiceRuntime[T] accessors — Core(), Options(), Config() — sit
//     on every service method call.
//   * Core.Go / Core.IsShutdown — wrap the tracked-goroutine + shutdown
//     poll path. IsShutdown is the inner-loop check in go-mlx workers.
//   * NewRuntime / NewWithFactories — boot-time, but worth gating.
//
// Run:    go test -bench='BenchmarkRuntime|BenchmarkServiceRuntime|BenchmarkCoreGo|BenchmarkIsShutdown|BenchmarkNewRuntime' -benchmem -run='^$' .

package core_test

import (
	"sync"

	. "dappco.re/go"
)

// Sinks defeat compiler dead-code elimination on the read paths.
var (
	runtimeSinkCore        *Core
	runtimeSinkOpts        CliOptions
	runtimeSinkConfig      *Config
	runtimeSinkBool        bool
	runtimeSinkString      string
	runtimeSinkSvcRuntime  *ServiceRuntime[CliOptions]
	runtimeSinkRuntimeRes  Result
)

// --- ServiceRuntime[T] ---

func BenchmarkNewServiceRuntime(b *B) {
	c := New()
	opts := CliOptions{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtimeSinkSvcRuntime = NewServiceRuntime(c, opts)
	}
}

func BenchmarkServiceRuntime_Core(b *B) {
	c := New()
	r := NewServiceRuntime(c, CliOptions{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtimeSinkCore = r.Core()
	}
}

func BenchmarkServiceRuntime_Options(b *B) {
	c := New()
	r := NewServiceRuntime(c, CliOptions{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtimeSinkOpts = r.Options()
	}
}

func BenchmarkServiceRuntime_Config(b *B) {
	c := New()
	r := NewServiceRuntime(c, CliOptions{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtimeSinkConfig = r.Config()
	}
}

// --- Core.Go / Core.IsShutdown (the worker-loop hot path) ---

func BenchmarkCoreIsShutdown(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtimeSinkBool = c.IsShutdown()
	}
}

func BenchmarkCoreGo_Spawn(b *B) {
	c := New()
	var wg sync.WaitGroup
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		c.Go(func() {
			wg.Done()
		})
	}
	wg.Wait()
}

// --- Runtime DTO ---

func BenchmarkNewRuntime(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewRuntime(nil)
	}
}

func BenchmarkNewWithFactories_Empty(b *B) {
	factories := map[string]ServiceFactory{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewWithFactories(nil, factories)
	}
}

func BenchmarkRuntime_ServiceName(b *B) {
	r := &Runtime{Core: New()}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtimeSinkString = r.ServiceName()
	}
}
