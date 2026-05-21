// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the time primitives in time.go.
// Per AX-11 — time helpers sit on every batching, timeout, and
// metric-stamp path. Now() is hit per-token in go-mlx inference
// when streaming with timestamps; TimeFormat lands in every log
// line; ParseDuration in every config-driven timeout.
//
// Sleep / After / NewTicker are excluded — they consume real time and
// would corrupt other concurrent benches.
//
// Run:    go test -bench='BenchmarkTime|BenchmarkNow|BenchmarkUnix|BenchmarkSince|BenchmarkUntil|BenchmarkParseDuration' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE on the read paths.
var (
	timeSinkTime     Time
	timeSinkDuration Duration
	timeSinkInt64    int64
	timeSinkString   string
	timeSinkResult   Result
)

// --- Now / UnixNow ---

func BenchmarkNow(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkTime = Now()
	}
}

func BenchmarkUnixNow(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkInt64 = UnixNow()
	}
}

// --- Since / Until ---

func BenchmarkSince(b *B) {
	start := Now()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkDuration = Since(start)
	}
}

func BenchmarkUntil(b *B) {
	future := Now().Add(Hour)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkDuration = Until(future)
	}
}

// --- Unix constructors ---

func BenchmarkUnixTime(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkTime = UnixTime(int64(i))
	}
}

func BenchmarkUnix(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkTime = Unix(int64(i), 0)
	}
}

func BenchmarkUnixMilli(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkTime = UnixMilli(int64(i) * 1000)
	}
}

// --- Format / Parse ---

func BenchmarkTimeFormat_RFC3339(b *B) {
	t := Now()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkString = TimeFormat(t, TimeRFC3339)
	}
}

func BenchmarkTimeFormat_DateTime(b *B) {
	t := Now()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkString = TimeFormat(t, TimeDateTime)
	}
}

func BenchmarkTimeParse_RFC3339(b *B) {
	value := "2026-04-28T07:00:00Z"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkResult = TimeParse(TimeRFC3339, value)
	}
}

// --- ParseDuration ---

func BenchmarkParseDuration_Simple(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkResult = ParseDuration("250ms")
	}
}

func BenchmarkParseDuration_Compound(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timeSinkResult = ParseDuration("1h30m45s")
	}
}
