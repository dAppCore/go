// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the integer conversion primitives in int.go.
// Per AX-11 — Atoi / Itoa / ParseInt / FormatInt land on every config
// parse, every HTTP query-string read, every protocol-port handler.
// These are strconv passthroughs but the bench harness gates the
// contract: if a Core reroute slows them down, the regression surfaces
// against a published floor.
//
// Run:    go test -bench='BenchmarkInt' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	intSinkResult Result
	intSinkString string
)

// --- Atoi / Itoa ---

func BenchmarkInt_Atoi(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkResult = Atoi("42")
	}
}

func BenchmarkInt_Atoi_Bad(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkResult = Atoi("not-a-number")
	}
}

func BenchmarkInt_Itoa(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkString = Itoa(123456)
	}
}

// --- FormatInt / FormatUint ---

func BenchmarkInt_FormatInt_Decimal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkString = FormatInt(1234567890, 10)
	}
}

func BenchmarkInt_FormatInt_Hex(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkString = FormatInt(255, 16)
	}
}

func BenchmarkInt_FormatUint_Hex(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkString = FormatUint(0xdeadbeef, 16)
	}
}

// --- ParseInt ---

func BenchmarkInt_ParseInt_Hex(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkResult = ParseInt("deadbeef", 16, 64)
	}
}

func BenchmarkInt_ParseInt_Decimal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSinkResult = ParseInt("1234567890", 10, 64)
	}
}
