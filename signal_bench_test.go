// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the OS-signal façade in signal.go — accessor, handler
// teardown, and capability probe. On a bare Core no handler is installed,
// so these measure the dispatch path without trapping real signals.
//
// Run:    go test -bench='BenchmarkSignal' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var (
	signalSinkP      *Signal
	signalSinkResult Result
	signalSinkBool   bool
)

func BenchmarkSignal(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		signalSinkP = c.Signal()
	}
}

func BenchmarkSignal_Stop(b *B) {
	s := New().Signal()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		signalSinkResult = s.Stop()
	}
}

func BenchmarkSignal_Exists(b *B) {
	s := New().Signal()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		signalSinkBool = s.Exists()
	}
}
