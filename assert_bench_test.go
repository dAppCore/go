// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the assertion helpers in assert.go. Each benchmark
// drives the success path (the predicate holds, so no t.Error/Fatal
// fires) — that is the cost every passing test pays, and the core test
// suite runs tens of thousands of these. The `any` parameters box, so
// these numbers also track the reflect/interface overhead of the
// comparison helpers.
//
// Run:    go test -bench='BenchmarkAssert|BenchmarkRequire' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

func BenchmarkAssertEqual(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertEqual(b, 42, 42)
	}
}

func BenchmarkAssertNotEqual(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNotEqual(b, 1, 2)
	}
}

func BenchmarkAssertTrue(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertTrue(b, true)
	}
}

func BenchmarkAssertFalse(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertFalse(b, false)
	}
}

func BenchmarkAssertNil(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNil(b, nil)
	}
}

func BenchmarkAssertNotNil(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNotNil(b, "ready")
	}
}

func BenchmarkAssertNoError(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNoError(b, nil)
	}
}

func BenchmarkAssertError(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertError(b, AnError)
	}
}

func BenchmarkAssertContains(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertContains(b, "agent ready", "ready")
	}
}

func BenchmarkAssertNotContains(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNotContains(b, "agent ready", "offline")
	}
}

func BenchmarkAssertLen(b *B) {
	v := []int{1, 2, 3}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertLen(b, v, 3)
	}
}

func BenchmarkAssertEmpty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertEmpty(b, "")
	}
}

func BenchmarkAssertNotEmpty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNotEmpty(b, "x")
	}
}

func BenchmarkAssertGreater(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertGreater(b, 5, 3)
	}
}

func BenchmarkAssertGreaterOrEqual(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertGreaterOrEqual(b, 5, 5)
	}
}

func BenchmarkAssertLess(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertLess(b, 3, 5)
	}
}

func BenchmarkAssertLessOrEqual(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertLessOrEqual(b, 5, 5)
	}
}

func BenchmarkAssertPanics(b *B) {
	fn := func() { panic("boom") }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertPanics(b, fn)
	}
}

func BenchmarkAssertNotPanics(b *B) {
	fn := func() {}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertNotPanics(b, fn)
	}
}

func BenchmarkAssertPanicsWithError(b *B) {
	fn := func() { panic(NewError("boom token")) }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertPanicsWithError(b, "boom", fn)
	}
}

func BenchmarkAssertErrorIs(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertErrorIs(b, AnError, AnError)
	}
}

func BenchmarkAssertInDelta(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertInDelta(b, 1.0, 1.0, 0.01)
	}
}

func BenchmarkAssertSame(b *B) {
	p := &struct{ n int }{n: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertSame(b, p, p)
	}
}

func BenchmarkAssertElementsMatch(b *B) {
	want := []int{1, 2, 3}
	got := []int{3, 2, 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertElementsMatch(b, want, got)
	}
}

func BenchmarkRequireNoError(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RequireNoError(b, nil)
	}
}

func BenchmarkRequireTrue(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RequireTrue(b, true)
	}
}

func BenchmarkRequireNotEmpty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RequireNotEmpty(b, "x")
	}
}

// Typed comparisons drive the uint64 / float64 branches of the ordering
// helpers (assertCmpUint64 / assertCmpFloat64).

func BenchmarkAssertGreater_Uint64(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertGreater(b, uint64(5), uint64(3))
	}
}

func BenchmarkAssertGreater_Float64(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertGreater(b, 5.0, 3.0)
	}
}
