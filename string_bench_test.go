// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the hot string primitives in string.go.
// Per AX-11 — benchmarks gate the performance contract of every
// primitive that downstream packages call on hot paths. The string
// helpers are exercised by every consumer (tokenisers, config keys,
// CLI dispatch, URL building, log formatting) so a regression here
// cascades across the ecosystem.
//
// Run:    go test -bench='Benchmark.*' -benchmem -run='^$' .
// Filter: go test -bench='BenchmarkJoin' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// --- Concat ---

func BenchmarkConcat_TwoParts(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Concat("https://", "example.com")
	}
}

func BenchmarkConcat_FourParts(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Concat("cmd.", "deploy.", "to.", "homelab")
	}
}

func BenchmarkConcat_EightParts(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Concat("a", "b", "c", "d", "e", "f", "g", "h")
	}
}

// --- Join ---

func BenchmarkJoin_TwoParts(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Join("/", "deploy", "homelab")
	}
}

func BenchmarkJoin_FourParts(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Join(".", "cmd", "deploy", "to", "homelab")
	}
}

func BenchmarkJoin_EightParts(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Join(" ", "the", "quick", "brown", "fox", "jumps", "over", "the", "lazy")
	}
}

// --- Lower / Upper ---

func BenchmarkLower_AlreadyLower(b *B) {
	s := "the quick brown fox jumps over the lazy dog"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Lower(s)
	}
}

func BenchmarkLower_Mixed(b *B) {
	s := "The Quick Brown Fox Jumps Over The Lazy Dog"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Lower(s)
	}
}

func BenchmarkUpper_AlreadyUpper(b *B) {
	s := "THE QUICK BROWN FOX"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Upper(s)
	}
}

func BenchmarkUpper_Mixed(b *B) {
	s := "the Quick brown FOX"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Upper(s)
	}
}

// --- Trim / TrimPrefix / TrimSuffix ---

func BenchmarkTrim_NoChange(b *B) {
	s := "hello"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Trim(s)
	}
}

func BenchmarkTrim_BothSides(b *B) {
	s := "   hello world   "
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Trim(s)
	}
}

func BenchmarkTrimPrefix_Hit(b *B) {
	s := "--verbose"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TrimPrefix(s, "--")
	}
}

func BenchmarkTrimSuffix_Hit(b *B) {
	s := "test.go"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TrimSuffix(s, ".go")
	}
}

// --- HasPrefix / HasSuffix / Contains ---

func BenchmarkHasPrefix(b *B) {
	s := "--verbose"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HasPrefix(s, "--")
	}
}

func BenchmarkHasSuffix(b *B) {
	s := "test.go"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HasSuffix(s, ".go")
	}
}

func BenchmarkContains(b *B) {
	s := "the quick brown fox jumps over the lazy dog"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Contains(s, "brown")
	}
}

// --- Split / SplitN / Index / LastIndex ---

func BenchmarkSplit_Small(b *B) {
	s := "a/b/c/d"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Split(s, "/")
	}
}

func BenchmarkSplit_Medium(b *B) {
	s := "the/quick/brown/fox/jumps/over/the/lazy/dog/today"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Split(s, "/")
	}
}

func BenchmarkSplitN(b *B) {
	s := "key=value=extra"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SplitN(s, "=", 2)
	}
}

func BenchmarkIndex(b *B) {
	s := "the quick brown fox jumps over the lazy dog"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Index(s, "lazy")
	}
}

func BenchmarkLastIndex(b *B) {
	s := "host.example.com:8080"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = LastIndex(s, ":")
	}
}

// --- Replace ---

func BenchmarkReplace_NoOp(b *B) {
	s := "no slashes here"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Replace(s, "/", ".")
	}
}

func BenchmarkReplace_Hit(b *B) {
	s := "deploy/to/homelab"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Replace(s, "/", ".")
	}
}

// --- RuneCount + HTML escape ---

func BenchmarkRuneCount_ASCII(b *B) {
	s := "the quick brown fox"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RuneCount(s)
	}
}

func BenchmarkHTMLEscape_NoSpecials(b *B) {
	s := "no special chars here just a plain sentence"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTMLEscape(s)
	}
}

func BenchmarkHTMLEscape_WithSpecials(b *B) {
	s := `<a href="/search?q=go&lang=en">Go</a>`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTMLEscape(s)
	}
}

func BenchmarkHTMLUnescape_NoSpecials(b *B) {
	s := "no special chars here just a plain sentence"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTMLUnescape(s)
	}
}

func BenchmarkHTMLUnescape_WithSpecials(b *B) {
	s := `&lt;a href=&quot;/search?q=go&amp;lang=en&quot;&gt;Go&lt;/a&gt;`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTMLUnescape(s)
	}
}

// --- Trim variants ---

func BenchmarkTrimCutset(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TrimCutset("[[task-id]]", "[]")
	}
}

func BenchmarkTrimLeft(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TrimLeft("---verbose", "-")
	}
}

func BenchmarkTrimRight(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TrimRight("hello!!!", "!")
	}
}

// --- Builder / Reader factories ---

func BenchmarkNewBuilder(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bld := NewBuilder()
		bld.WriteString("hello")
		_ = bld.String()
	}
}

// NewReader bench lives in io_bench_test.go (it is io-oriented).

// --- Repeat / Count / Fields / EqualFold / Cut (hot scans) ---

func BenchmarkRepeat(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Repeat("=", 64)
	}
}

func BenchmarkCount(b *B) {
	s := "a.b.c.d.e.f.g.h"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Count(s, ".")
	}
}

func BenchmarkFields(b *B) {
	s := "the quick brown fox jumps over the lazy dog"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Fields(s)
	}
}

func BenchmarkEqualFold_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EqualFold("Authorization", "authorization")
	}
}

func BenchmarkCut(b *B) {
	s := "Authorization: Bearer abc123"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _ = Cut(s, ": ")
	}
}
