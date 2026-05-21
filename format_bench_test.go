// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the formatting primitives in format.go.
// Per AX-11 — Sprintf is on every error message, every log key=val
// pair, every CLI output line. Bare passthroughs to fmt but the bench
// harness gates the contract so a Core reroute can't slow them.
//
// Run:    go test -bench='BenchmarkFormat' -benchmem -run='^$' .

package core_test

import (
	"bytes"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	formatSinkString string
	formatSinkErr    error
)

// --- Sprint ---

func BenchmarkFormat_Sprint_OneArg(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkString = Sprint(42)
	}
}

func BenchmarkFormat_Sprint_ThreeArgs(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkString = Sprint("agent", "ready", 9000)
	}
}

// --- Sprintf ---

func BenchmarkFormat_Sprintf_Verbs(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkString = Sprintf("%s connected on %d", "homelab", 9000)
	}
}

func BenchmarkFormat_Sprintf_KeyVal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkString = Sprintf("%v=%q", "key", "value")
	}
}

// --- Sprintln ---

func BenchmarkFormat_Sprintln(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkString = Sprintln("agent", "ready")
	}
}

// --- Print to a buffer (writer-mediated, avoids stdout in bench) ---

func BenchmarkFormat_Print_Buffer(b *B) {
	var buf bytes.Buffer
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		Print(&buf, "port: %d", 9000)
	}
}

// --- Errorf ---

func BenchmarkFormat_Errorf_NoWrap(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkErr = Errorf("connect %s failed", "homelab")
	}
}

func BenchmarkFormat_Errorf_Wrap(b *B) {
	cause := NewError("connection refused")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatSinkErr = Errorf("connect %s: %w", "homelab", cause)
	}
}
