// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the zero-copy view primitives in unsafe.go.
// Per AX-11 — AsBytes / AsString are the SPOR for unsafe-package
// zero-copy conversion. Hot path: hash digest of a string, Writer.Write
// of a string body, ReadAll → string return. Each shaves the source
// length in alloc + copy cost from the call.
//
// Run:    go test -bench='BenchmarkAs(Bytes|String)' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	unsafeSinkBytes  []byte
	unsafeSinkString string
)

var (
	unsafeStr16   = "aaaaaaaaaaaaaaaa"
	unsafeStr1K   = makeString(1024)
	unsafeBytes16 = []byte("aaaaaaaaaaaaaaaa")
	unsafeBytes1K = func() []byte {
		out := make([]byte, 1024)
		for i := range out {
			out[i] = byte('a' + i%26)
		}
		return out
	}()
)

func makeString(n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte('a' + i%26)
	}
	return string(out)
}

// --- AsBytes (zero-copy) vs. []byte(s) baseline ---

func BenchmarkAsBytes_16B(b *B) {
	s := unsafeStr16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = AsBytes(s)
	}
}

func BenchmarkAsBytes_1KB(b *B) {
	s := unsafeStr1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = AsBytes(s)
	}
}

func BenchmarkAsBytes_Baseline_16B(b *B) {
	s := unsafeStr16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = []byte(s)
	}
}

func BenchmarkAsBytes_Baseline_1KB(b *B) {
	s := unsafeStr1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = []byte(s)
	}
}

// --- AsString (zero-copy) vs. string(b) baseline ---

func BenchmarkAsString_16B(b *B) {
	bs := unsafeBytes16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = AsString(bs)
	}
}

func BenchmarkAsString_1KB(b *B) {
	bs := unsafeBytes1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = AsString(bs)
	}
}

func BenchmarkAsString_Baseline_16B(b *B) {
	bs := unsafeBytes16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = string(bs)
	}
}

func BenchmarkAsString_Baseline_1KB(b *B) {
	bs := unsafeBytes1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = string(bs)
	}
}

// --- Empty input fast paths ---

func BenchmarkAsBytes_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = AsBytes("")
	}
}

func BenchmarkAsString_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = AsString(nil)
	}
}
