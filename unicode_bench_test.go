// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the rune helpers in unicode.go.
// Per AX-11 — IsLetter / IsDigit / IsSpace / ToUpper / ToLower are
// inner-loop ops on every tokeniser pass + every CLI prompt scanner.
// Bare passthroughs to unicode but the bench harness gates the contract.
//
// Run:    go test -bench='BenchmarkUnicode' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	unicodeSinkBool bool
	unicodeSinkRune rune
)

// --- Predicates (ASCII fast path) ---

func BenchmarkUnicode_IsLetter_ASCII(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsLetter('A')
	}
}

func BenchmarkUnicode_IsLetter_NonASCII(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsLetter('Ω')
	}
}

func BenchmarkUnicode_IsDigit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsDigit('9')
	}
}

func BenchmarkUnicode_IsSpace_Newline(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsSpace('\n')
	}
}

func BenchmarkUnicode_IsSpace_NonSpace(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsSpace('A')
	}
}

func BenchmarkUnicode_IsUpper(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsUpper('A')
	}
}

func BenchmarkUnicode_IsLower(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkBool = IsLower('a')
	}
}

// --- Case conversion ---

func BenchmarkUnicode_ToUpper_ASCII(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkRune = ToUpper('a')
	}
}

func BenchmarkUnicode_ToUpper_NonASCII(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkRune = ToUpper('ω')
	}
}

func BenchmarkUnicode_ToLower_ASCII(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unicodeSinkRune = ToLower('A')
	}
}
