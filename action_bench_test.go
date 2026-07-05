// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Action / Task primitives in action.go.
// Per AX-11 — Action.Run + the entitlement check + panic recover sit
// on every named-action dispatch (`c.Action("process.run").Run(...)`).
// Action lookup happens on every Wails RPC, every CLI subcommand, every
// agentic dispatch — the bench gates the contract floor.
//
// Run:    go test -bench='BenchmarkAction' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	actionSinkResult  Result
	actionSinkBool    bool
	actionSinkStrings []string
	actionSinkAction  *Action
)

// noopHandler returns a successful Result without touching ctx/opts.
// Mirrors the minimum-work shape consumers register.
func noopActionHandler() ActionHandler {
	return func(ctx Context, opts Options) Result { return Result{OK: true} }
}

// actionFixture builds a Core with one registered action ready to run.
func actionFixture() (*Core, *Action) {
	c := New()
	c.Action("bench.action", noopActionHandler())
	return c, c.Action("bench.action")
}

// --- Lookup ---

func BenchmarkAction_Lookup_Hit(b *B) {
	c, _ := actionFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkAction = c.Action("bench.action")
	}
}

func BenchmarkAction_Lookup_Miss(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkAction = c.Action("noexist.action")
	}
}

// --- Register ---

func BenchmarkAction_Register(b *B) {
	c := New()
	h := noopActionHandler()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.Action("bench.register", h)
	}
}

// --- Run ---

func BenchmarkAction_Run(b *B) {
	c, action := actionFixture()
	ctx := c.Context()
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkResult = action.Run(ctx, opts)
	}
}

func BenchmarkAction_Run_NoHandler(b *B) {
	c := New()
	action := c.Action("noexist.action")
	ctx := c.Context()
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkResult = action.Run(ctx, opts)
	}
}

func BenchmarkAction_Run_Disabled(b *B) {
	c, action := actionFixture()
	action.Disable()
	ctx := c.Context()
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkResult = action.Run(ctx, opts)
	}
}

// --- Predicates / state ---

func BenchmarkAction_Exists_Hit(b *B) {
	_, action := actionFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkBool = action.Exists()
	}
}

func BenchmarkAction_Exists_Miss(b *B) {
	c := New()
	action := c.Action("noexist.action")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkBool = action.Exists()
	}
}

func BenchmarkAction_Enabled(b *B) {
	_, action := actionFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkBool = action.Enabled()
	}
}

func BenchmarkAction_Enable(b *B) {
	_, action := actionFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		action.Enable()
	}
}

func BenchmarkAction_Disable(b *B) {
	_, action := actionFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		action.Disable()
	}
}

// --- Listing ---

func BenchmarkAction_Actions_Empty(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkStrings = c.Actions()
	}
}

func BenchmarkAction_Actions_TenRegistered(b *B) {
	c := New()
	h := noopActionHandler()
	for i := 0; i < 10; i++ {
		c.Action(Sprintf("bench.action.%d", i), h)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkStrings = c.Actions()
	}
}

// --- Task ---

func BenchmarkAction_Task_Register(b *B) {
	c := New()
	c.Action("bench.step", noopActionHandler())
	def := Task{
		Name:  "bench.task",
		Steps: []Step{{Action: "bench.step"}, {Action: "bench.step"}},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.Task("bench.task", def)
	}
}

func BenchmarkAction_Task_Run_TwoSteps(b *B) {
	c := New()
	c.Action("bench.step", noopActionHandler())
	task := c.Task("bench.task", Task{
		Name:  "bench.task",
		Steps: []Step{{Action: "bench.step"}, {Action: "bench.step"}},
	})
	ctx := c.Context()
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkResult = task.Run(ctx, c, opts)
	}
}

func BenchmarkAction_Tasks(b *B) {
	c := New()
	c.Action("bench.step", noopActionHandler())
	for i := 0; i < 5; i++ {
		c.Task(Sprintf("bench.task.%d", i), Task{
			Name:  Sprintf("bench.task.%d", i),
			Steps: []Step{{Action: "bench.step"}},
		})
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkStrings = c.Tasks()
	}
}

func BenchmarkCore_PerformAsync(b *B) {
	c := New()
	c.Action("bench.async", noopActionHandler())
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		actionSinkResult = c.PerformAsync("bench.async", opts)
	}
}

func BenchmarkCore_Progress(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Progress("bench.task", 0.5, "halfway", "bench.async")
	}
}
