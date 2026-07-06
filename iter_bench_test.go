// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the iter primitives in iter.go.
// Per AX-11 — Seq / Seq2 are type aliases (zero cost); Pull / Pull2
// convert push-style to pull-style by spawning a coordination goroutine
// (heavyweight setup, but useful when consumers need pull semantics).
//
// Run:    go test -bench='BenchmarkIter|BenchmarkPull' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	iterSinkInt int
	iterSinkOK  bool
)

// --- range over Seq ---

func BenchmarkIter_RangeSeq_Small(b *B) {
	seq := func(yield func(int) bool) {
		for i := 0; i < 10; i++ {
			if !yield(i) {
				return
			}
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sum := 0
		for n := range seq {
			sum += n
		}
		iterSinkInt = sum
	}
}

func BenchmarkIter_RangeSeq_Large(b *B) {
	seq := func(yield func(int) bool) {
		for i := 0; i < 1000; i++ {
			if !yield(i) {
				return
			}
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sum := 0
		for n := range seq {
			sum += n
		}
		iterSinkInt = sum
	}
}

// --- Pull ---

func BenchmarkIter_Pull_Cycle(b *B) {
	seq := func(yield func(int) bool) {
		for i := 0; i < 10; i++ {
			if !yield(i) {
				return
			}
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		next, stop := Pull(seq)
		for {
			_, ok := next()
			if !ok {
				break
			}
		}
		stop()
	}
}

func BenchmarkIter_Pull_OnlyFirst(b *B) {
	seq := func(yield func(int) bool) {
		for i := 0; i < 1000; i++ {
			if !yield(i) {
				return
			}
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		next, stop := Pull(seq)
		iterSinkInt, iterSinkOK = next()
		stop()
	}
}

// --- Pull2 ---

func BenchmarkIter_Pull2_Cycle(b *B) {
	seq2 := func(yield func(string, int) bool) {
		for i := 0; i < 10; i++ {
			if !yield("k", i) {
				return
			}
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		next, stop := Pull2(seq2)
		for {
			_, _, ok := next()
			if !ok {
				break
			}
		}
		stop()
	}
}
