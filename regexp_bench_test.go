// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the regex primitives in regexp.go.
// Per AX-11 — compiled regex matching sits on URL routing, tokeniser
// pattern matching, log parsing, and format validators. Compilation
// is one-shot at boot; matching is per-call.
//
// Run:    go test -bench='BenchmarkRegex' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	regexSinkBool    bool
	regexSinkString  string
	regexSinkStrings []string
	regexSinkResult  Result
)

// Compiled regex fixtures (compiled once, reused across benches —
// the realistic shape; consumers compile at init).
var (
	regexSimple  *Regexp
	regexDigits  *Regexp
	regexCapture *Regexp
	regexCSV     *Regexp
)

func init() {
	if r := Regex(`agent-[0-9]+`); r.OK {
		regexSimple = r.Value.(*Regexp)
	}
	if r := Regex(`\d+`); r.OK {
		regexDigits = r.Value.(*Regexp)
	}
	if r := Regex(`(\w+)=(\d+)`); r.OK {
		regexCapture = r.Value.(*Regexp)
	}
	if r := Regex(`,+`); r.OK {
		regexCSV = r.Value.(*Regexp)
	}
}

// --- Compile ---

func BenchmarkRegex_Compile_Simple(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkResult = Regex(`agent-[0-9]+`)
	}
}

func BenchmarkRegex_Compile_Complex(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkResult = Regex(`^(?P<scheme>https?)://(?P<host>[^/]+)(?P<path>/.*)?$`)
	}
}

// --- MatchString ---

func BenchmarkRegex_MatchString_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkBool = regexSimple.MatchString("hello agent-42 world")
	}
}

func BenchmarkRegex_MatchString_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkBool = regexSimple.MatchString("hello world")
	}
}

// --- FindString / FindAllString ---

func BenchmarkRegex_FindString(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkString = regexDigits.FindString("hello 42 world 1234")
	}
}

func BenchmarkRegex_FindAllString_Few(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkStrings = regexDigits.FindAllString("a1 b2 c3", -1)
	}
}

func BenchmarkRegex_FindAllString_Many(b *B) {
	input := "tok1 tok2 tok3 tok4 tok5 tok6 tok7 tok8 tok9 tok10 tok11 tok12"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkStrings = regexDigits.FindAllString(input, -1)
	}
}

// --- FindStringSubmatch ---

func BenchmarkRegex_FindStringSubmatch(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkStrings = regexCapture.FindStringSubmatch("count=42")
	}
}

// --- ReplaceAllString ---

func BenchmarkRegex_ReplaceAllString(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkString = regexDigits.ReplaceAllString("a1 b2 c3 d4 e5", "X")
	}
}

// --- Split ---

func BenchmarkRegex_Split(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkStrings = regexCSV.Split("a,b,,c,,,d,e", -1)
	}
}

// --- String (source) ---

func BenchmarkRegex_String(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexSinkString = regexSimple.String()
	}
}
