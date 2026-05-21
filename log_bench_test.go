// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the log primitives in log.go.
// Per AX-11 — log.Debug / log.Info / log.Warn / log.Error are on
// every observability path. Below-threshold logs (Debug at Info level)
// should be nearly free; above-threshold logs land in the formatter
// + writer path.
//
// Run:    go test -bench='BenchmarkLog' -benchmem -run='^$' .

package core_test

import (
	"io"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var logSinkInt int

// logBenchFixture builds a Log at the given level that writes to
// io.Discard so the bench measures formatter cost without IO noise.
func logBenchFixture(level Level) *Log {
	return NewLog(LogOptions{
		Level:  level,
		Output: io.Discard,
	})
}

// --- Below-threshold (gated, should be cheap) ---

func BenchmarkLog_Debug_BelowThreshold(b *B) {
	l := logBenchFixture(LevelInfo)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Debug("token sampled", "id", 42, "score", 0.87)
	}
}

func BenchmarkLog_Info_BelowThreshold(b *B) {
	l := logBenchFixture(LevelWarn)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("agent ready", "host", "homelab.lan")
	}
}

// --- Above-threshold (formatter + writer) ---

func BenchmarkLog_Info_NoKeyvals(b *B) {
	l := logBenchFixture(LevelDebug)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("agent ready")
	}
}

func BenchmarkLog_Info_TwoKeyvals(b *B) {
	l := logBenchFixture(LevelDebug)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("agent ready", "host", "homelab.lan", "port", 9000)
	}
}

func BenchmarkLog_Warn_SixKeyvals(b *B) {
	l := logBenchFixture(LevelDebug)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Warn("rate limited",
			"endpoint", "agentic.status",
			"limit", 100,
			"window", "1m",
			"retry_after", "30s",
			"upstream", "homelab",
			"backoff", 5.0,
		)
	}
}

// --- Level introspection ---

func BenchmarkLog_Level(b *B) {
	l := logBenchFixture(LevelInfo)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logSinkInt = int(l.Level())
	}
}

// --- Package-level functions ---

func BenchmarkLog_PackageInfo_BelowThreshold(b *B) {
	prev := Default()
	SetDefault(logBenchFixture(LevelWarn))
	defer SetDefault(prev)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Info("agent ready", "host", "homelab.lan")
	}
}
