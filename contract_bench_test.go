// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the framework-contract primitives in contract.go.
// Per AX-11 — the With* CoreOption functions are applied once per
// core.New() call. The Action* contract structs (ActionServiceStartup,
// ActionTaskStarted/Progress/Completed) flow across the IPC bus on
// every lifecycle event. The bench gates the option-application path
// and the contract-construction cost.
//
// Run:    go test -bench='BenchmarkContract' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	contractSinkCore    *Core
	contractSinkOption  CoreOption
	contractSinkMessage Message
)

// --- New + WithOption flow ---

func BenchmarkContract_New_WithOptions(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkCore = New(
			WithOption("name", "bench"),
			WithOption("port", 9000),
		)
	}
}

func BenchmarkContract_New_WithFullOptions(b *B) {
	opts := NewOptions(
		Option{Key: "name", Value: "bench"},
		Option{Key: "port", Value: 9000},
		Option{Key: "debug", Value: true},
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkCore = New(WithOptions(opts))
	}
}

func BenchmarkContract_New_WithServiceLock(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkCore = New(WithServiceLock())
	}
}

// --- CoreOption constructors (closure allocation) ---

func BenchmarkContract_WithOption(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkOption = WithOption("name", "bench")
	}
}

func BenchmarkContract_WithOptions(b *B) {
	opts := NewOptions(Option{Key: "name", Value: "bench"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkOption = WithOptions(opts)
	}
}

func BenchmarkContract_WithName(b *B) {
	factory := func(c *Core) Result {
		return Result{Value: &struct{ id int }{id: 1}, OK: true}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkOption = WithName("bench", factory)
	}
}

// --- Lifecycle action contracts ---

func BenchmarkContract_ActionServiceStartup(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkMessage = ActionServiceStartup{}
	}
}

func BenchmarkContract_ActionServiceShutdown(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkMessage = ActionServiceShutdown{}
	}
}

// --- Task lifecycle contracts (carry payload) ---

func BenchmarkContract_ActionTaskStarted(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkMessage = ActionTaskStarted{
			TaskIdentifier: "task-42",
			Action:         "agentic.dispatch",
		}
	}
}

func BenchmarkContract_ActionTaskProgress(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkMessage = ActionTaskProgress{
			TaskIdentifier: "task-42",
			Action:         "agentic.dispatch",
			Progress:       0.5,
			Message:        "halfway",
		}
	}
}

func BenchmarkContract_ActionTaskCompleted(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		contractSinkMessage = ActionTaskCompleted{
			TaskIdentifier: "task-42",
			Action:         "agentic.dispatch",
			Result:         Result{OK: true},
		}
	}
}
