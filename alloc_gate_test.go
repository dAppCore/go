// SPDX-License-Identifier: EUPL-1.2

// Allocation regression gates. Unlike a Benchmark (which only prints a
// number nobody compares), each case here asserts a hard ceiling on
// allocations-per-call via testing.AllocsPerRun — so the suite FAILS if a
// hot primitive starts allocating where it didn't before. ns/op is
// deliberately not gated (machine-dependent); allocs are stable.
//
// Ceilings are the measured floor as of this commit:
//   0 — pure / zero-copy primitive (no heap touch)
//   1 — one unavoidable alloc: Result{Value:<non-pointer>} boxes the value
//       into `any` (string/int → 1 alloc; pointer values box for free)
//
// Tightening a ceiling that's been beaten is welcome; raising one is a
// regression that needs justifying in review.
//
// Run:    go test -run 'TestAllocs_' .

package core_test

import (
	"testing"

	. "dappco.re/go"
)

// gate sinks — keep the measured call from being elided as dead code.
var (
	gateInt    int
	gateBool   bool
	gateStr    string
	gateBytes  []byte
	gateResult Result
	gateAny    any
	gateDiags  []LSPDiagnostic
)

func TestAllocs_HotPrimitives(t *T) {
	ptr := &struct{ n int }{}
	bs := []byte("homelab")
	var ai32 AtomicInt32
	var ab AtomicBool
	var au64 AtomicUint64
	cases := []struct {
		name    string
		ceiling int
		fn      func()
	}{
		// math / ordering — pure, must stay 0
		{"Compare", 0, func() { gateInt = Compare(3, 7) }},
		{"Min", 0, func() { gateInt = Min(3, 7) }},
		{"Max", 0, func() { gateInt = Max(3, 7) }},
		{"Abs", 0, func() { gateInt = Abs(-42) }},
		{"Clamp", 0, func() { gateInt = Clamp(15, 0, 10) }},
		{"Sign", 0, func() { gateInt = Sign(-3) }},

		// string predicates / index — pure, must stay 0
		{"Contains", 0, func() { gateBool = Contains("agent ready", "ready") }},
		{"HasPrefix", 0, func() { gateBool = HasPrefix("agent.go", "agent") }},
		{"HasSuffix", 0, func() { gateBool = HasSuffix("agent.go", ".go") }},
		{"EqualFold", 0, func() { gateBool = EqualFold("Bearer", "bearer") }},
		{"Index", 0, func() { gateInt = Index("agent", "e") }},
		{"Count", 0, func() { gateInt = Count("banana", "a") }},
		{"Trim", 0, func() { gateStr = Trim("  x  ") }},

		// slice search — pure, must stay 0
		{"SliceContains", 0, func() { gateBool = SliceContains([]int{1, 2, 3}, 2) }},
		{"SliceIndex", 0, func() { gateInt = SliceIndex([]int{1, 2, 3}, 2) }},

		// unsafe zero-copy conversions — must stay 0
		{"AsBytes", 0, func() { gateBytes = AsBytes("homelab") }},
		{"AsString", 0, func() { gateStr = AsString(bs) }},

		// Result constructors — pointer value boxes for free (0); a
		// non-pointer value pays the one inherent Result-boxing alloc.
		{"Ok_Pointer", 0, func() { gateResult = Ok(ptr) }},
		{"Fail", 0, func() { gateResult = Fail(AnError) }},
		{"Ok_String", 1, func() { gateResult = Ok("ready") }},
		{"ResultOf_String", 1, func() { gateResult = ResultOf("ready", nil) }},

		// Result accessors — pure reads, must stay 0
		{"Result_Or", 0, func() { gateAny = (Result{OK: false}).Or("fallback") }},

		// atomics — lock-free, must stay 0
		{"AtomicInt32_Load", 0, func() { gateInt = int(ai32.Load()) }},
		{"AtomicInt32_Store", 0, func() { ai32.Store(7) }},
		{"AtomicInt32_Add", 0, func() { gateInt = int(ai32.Add(1)) }},
		{"AtomicBool_Load", 0, func() { gateBool = ab.Load() }},
		{"AtomicBool_Store", 0, func() { ab.Store(true) }},
		{"AtomicUint64_Add", 0, func() { gateInt = int(au64.Add(1)) }},

		// more math — pure, must stay 0
		{"Pow", 0, func() { gateInt = int(Pow(2, 8)) }},
		{"Floor", 0, func() { gateInt = int(Floor(3.7)) }},
		{"Ceil", 0, func() { gateInt = int(Ceil(3.1)) }},
		{"Round", 0, func() { gateInt = int(Round(3.5)) }},

		// slice — Clone copies (1 alloc for the backing array); Reverse is in-place (0)
		{"SliceClone", 1, func() { gateBytes = SliceClone(bs) }},
	}
	for _, c := range cases {
		avg := int(testing.AllocsPerRun(1000, c.fn))
		AssertLessOrEqual(t, avg, c.ceiling, c.name)
	}
}

// TestAllocs_LSPComputeDiagnostics locks the regex-hoist win (c1cf230):
// with per-call regexp.Compile it allocated ~186/call; compiling the
// naming-diagnostic patterns once dropped it to ~47. A regression to
// per-call compilation jumps back over 100, so the ceiling catches it.
func TestAllocs_LSPComputeDiagnostics(t *T) {
	uri := "file:///agent/worker.go"
	content := []byte("package worker\n\nfunc process(in string) (string, error) {\n\treturn in, nil\n}\n")
	LSPComputeDiagnostics(uri, content) // warm the per-dir cache → measure steady state
	avg := int(testing.AllocsPerRun(200, func() {
		gateDiags = LSPComputeDiagnostics(uri, content)
	}))
	AssertLessOrEqual(t, avg, 70, "LSPComputeDiagnostics")
}
