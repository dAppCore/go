// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Result primitives in result.go.
// Per AX-11 — Result is on EVERY Core return. Ok / Fail / ResultOf are
// constructors; Code / Error are inspection accessors; Must / MustCast
// / Cast / Or / Try are call-site collapse helpers. Each one runs in
// every consumer's hot path. The harness gates the contract floor.
//
// Run:    go test -bench='BenchmarkResult' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	resultSinkRes    Result
	resultSinkAny    any
	resultSinkString string
	resultSinkInt    int
	resultSinkBool   bool
)

// Fixtures
var (
	resultOK    = Ok(42)
	resultFail  = Fail(NewError("boom"))
	resultCoded = Fail(NewCode("FS_NOTFOUND", "missing"))
)

// --- Constructors ---

func BenchmarkResult_Ok(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkRes = Ok(42)
	}
}

func BenchmarkResult_Fail(b *B) {
	err := NewError("boom")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkRes = Fail(err)
	}
}

func BenchmarkResult_ResultOf_OK(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkRes = ResultOf("data", nil)
	}
}

func BenchmarkResult_ResultOf_Err(b *B) {
	err := NewError("boom")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkRes = ResultOf(nil, err)
	}
}

// --- Inspection ---

func BenchmarkResult_Error_OK(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkString = resultOK.Error()
	}
}

func BenchmarkResult_Error_Fail(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkString = resultFail.Error()
	}
}

func BenchmarkResult_Code_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkString = resultCoded.Code()
	}
}

func BenchmarkResult_Code_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkString = resultOK.Code()
	}
}

// --- Collapse helpers ---

func BenchmarkResult_Or_OK(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkAny = resultOK.Or(0)
	}
}

func BenchmarkResult_Or_Fallback(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkAny = resultFail.Or(99)
	}
}

func BenchmarkResult_Must_OK(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkAny = resultOK.Must()
	}
}

// --- Cast / MustCast ---

func BenchmarkResult_Cast_Hit(b *B) {
	r := Ok(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkInt, resultSinkBool = Cast[int](r)
	}
}

func BenchmarkResult_Cast_Miss(b *B) {
	r := Ok("string")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkInt, resultSinkBool = Cast[int](r)
	}
}

func BenchmarkResult_MustCast(b *B) {
	r := Ok(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkInt = MustCast[int](r)
	}
}

// --- Try ---

func BenchmarkResult_Try_NoPanic(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkRes = Try(func() any { return 42 })
	}
}

func BenchmarkResult_Try_ErrReturn(b *B) {
	err := NewError("returned err")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultSinkRes = Try(func() any { return err })
	}
}
