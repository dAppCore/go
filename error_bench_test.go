// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the error primitives in error.go.
// Per AX-11 — E / Wrap / WrapCode / NewCode are on every core.Fail
// path; Err.Error() formatting + Operation / ErrorCode accessors land
// on every log line that prints an error. Even modest per-call overhead
// here compounds across the ecosystem.
//
// Run:    go test -bench='BenchmarkErr' -benchmem -run='^$' .

package core_test

import (
	"errors"
	"testing"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	errSinkErr  error
	errSinkStr  string
	errSinkBool bool
)

// Fixtures
var (
	errSentinel = errors.New("connection refused")
	errCoreLeaf = E("net.Dial", "dial failed", errSentinel)
	errCoreDeep = Wrap(Wrap(errCoreLeaf, "api.Call", "remote unreachable"), "agent.Ping", "homelab probe failed")
	errCoded    = NewCode("CONFIG_MISSING", "missing config.host")
)

// --- Construction ---

func BenchmarkErr_E_WithCause(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = E("user.Save", "failed", errSentinel)
	}
}

func BenchmarkErr_E_NoCause(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = E("api.Call", "rate limited", nil)
	}
}

func BenchmarkErr_Wrap_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = Wrap(errSentinel, "db.Query", "query failed")
	}
}

func BenchmarkErr_Wrap_NilPassthrough(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = Wrap(nil, "db.Query", "query failed")
	}
}

func BenchmarkErr_WrapCode(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = WrapCode(errSentinel, "DB_FAIL", "db.Query", "query failed")
	}
}

func BenchmarkErr_NewCode(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = NewCode("CONFIG_MISSING", "missing config.host")
	}
}

func BenchmarkErr_NewError(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = NewError("connection refused")
	}
}

// --- Error() formatting ---

func BenchmarkErr_Error_Leaf(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = errCoreLeaf.Error()
	}
}

func BenchmarkErr_Error_Deep(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = errCoreDeep.Error()
	}
}

func BenchmarkErr_Error_Coded(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = errCoded.Error()
	}
}

// --- Inspection ---

func BenchmarkErr_Operation(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = Operation(errCoreDeep)
	}
}

func BenchmarkErr_ErrorCode(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = ErrorCode(errCoded)
	}
}

func BenchmarkErr_ErrorMessage(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = ErrorMessage(errCoreLeaf)
	}
}

func BenchmarkErr_Root(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = Root(errCoreDeep)
	}
}

func BenchmarkErr_Is(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkBool = Is(errCoreDeep, errSentinel)
	}
}

func BenchmarkErr_As(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var target *Err
		errSinkBool = As(errCoreDeep, &target)
	}
}

// --- ErrorJoin ---

func BenchmarkErr_ErrorJoin_Two(b *B) {
	a := errors.New("a failed")
	c := errors.New("c failed")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkErr = ErrorJoin(a, c)
	}
}

// --- (*Err).Unwrap ---

func BenchmarkErr_Unwrap(b *B) {
	b.ReportAllocs()
	leaf := errCoreLeaf.(*Err)
	for i := 0; i < b.N; i++ {
		errSinkErr = leaf.Unwrap()
	}
}

// --- AllOperations / StackTrace / FormatStackTrace ---

func BenchmarkErr_AllOperations_Iter(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		count := 0
		for range AllOperations(errCoreDeep) {
			count++
		}
		errSinkBool = count > 0
	}
}

func BenchmarkErr_StackTrace_Deep(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stack := StackTrace(errCoreDeep)
		errSinkBool = len(stack) > 0
	}
}

func BenchmarkErr_FormatStackTrace_Deep(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		errSinkStr = FormatStackTrace(errCoreDeep)
	}
}

// --- ErrorPanic surface ---
//
// Recover / SafeGo are not hot paths — their numbers are informational
// (a panic + recover, or a goroutine spawn, per iteration), benched here
// for coverage completeness. Reports() is the inspection accessor.

func BenchmarkErrorLog_Must(b *B) {
	el := New().Log()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		el.Must(nil, "Bench", "no error") // nil err = no panic
	}
}

func BenchmarkErrorPanic_Recover(b *B) {
	h := New().Error()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		func() {
			defer h.Recover()
			panic("bench")
		}()
	}
}

func BenchmarkErrorPanic_SafeGo(b *B) {
	h := New().Error()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h.SafeGo(func() {})
	}
}

func BenchmarkErr_ErrorPanic_Reports(b *B) {
	c := New()
	ep := c.Error()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ep.Reports(10)
	}
}

// Quiet unused-imports for "testing" — used implicitly via core.B alias chain.
var _ = testing.Short
