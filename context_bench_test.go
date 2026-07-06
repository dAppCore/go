// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the context primitives in context.go.
// Per AX-11 — Background / TODO / WithCancel / WithTimeout / WithValue
// are called by every Action handler and Service lifecycle entry point.
// Bare passthroughs to context but the bench harness gates the contract.
//
// Run:    go test -bench='BenchmarkCtx' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	ctxSinkContext Context
	ctxSinkCancel  CancelFunc
	ctxSinkAny     any
)

type ctxKey struct{}

// --- Roots ---

func BenchmarkCtx_Background(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkContext = Background()
	}
}

func BenchmarkCtx_TODO(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkContext = TODO()
	}
}

// --- Cancellable variants ---
//
// WithCancel + WithTimeout + WithDeadline allocate a child ctx node
// each call. The cancel func must be invoked to release the goroutine
// the runtime starts behind it — otherwise the bench leaks one per iter.

func BenchmarkCtx_WithCancel(b *B) {
	parent := Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkContext, ctxSinkCancel = WithCancel(parent)
		ctxSinkCancel()
	}
}

func BenchmarkCtx_WithTimeout(b *B) {
	parent := Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkContext, ctxSinkCancel = WithTimeout(parent, 1*Second)
		ctxSinkCancel()
	}
}

func BenchmarkCtx_WithDeadline(b *B) {
	parent := Background()
	deadline := Now().Add(1 * Second)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkContext, ctxSinkCancel = WithDeadline(parent, deadline)
		ctxSinkCancel()
	}
}

// --- Values ---

func BenchmarkCtx_WithValue(b *B) {
	parent := Background()
	k := ctxKey{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkContext = WithValue(parent, k, "req-12345")
	}
}

func BenchmarkCtx_Value_Lookup(b *B) {
	k := ctxKey{}
	ctx := WithValue(Background(), k, "req-12345")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctxSinkAny = ctx.Value(k)
	}
}
