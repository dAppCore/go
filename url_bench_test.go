// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the URL primitives in url.go.
// Per AX-11 — URL parse / encode / decode sit on inference endpoints,
// HF model URLs, route dispatch, redirect handling.
//
// Run:    go test -bench='BenchmarkURL' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	urlSinkResult Result
	urlSinkString string
)

// --- Parse / Normalize ---

func BenchmarkURL_Parse_Simple(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkResult = URLParse("https://forge.lthn.sh/agent/cladius")
	}
}

func BenchmarkURL_Parse_WithQuery(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkResult = URLParse("https://hf.co/qwen/qwen3.6-27b?file=qwen3.6-q4.gguf&revision=main")
	}
}

func BenchmarkURL_Parse_Complex(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkResult = URLParse("https://user:pass@host.example.com:8443/api/v1/agents?page=1&size=50#anchor")
	}
}

func BenchmarkURL_Normalize(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkString = URLNormalize("https://example.com/a b")
	}
}

// --- Encode / Decode ---

func BenchmarkURL_Encode_NoSpecials(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkString = URLEncode("simple-key-value")
	}
}

func BenchmarkURL_Encode_WithSpecials(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkString = URLEncode("model=qwen/qwen3.6 token=abc&def")
	}
}

func BenchmarkURL_Decode_NoSpecials(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkResult = URLDecode("simple-key-value")
	}
}

func BenchmarkURL_Decode_WithSpecials(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkResult = URLDecode("model%3Dqwen%2Fqwen3.6+token%3Dabc%26def")
	}
}

// --- PathEscape ---

func BenchmarkURL_PathEscape_NoSpecials(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkString = URLPathEscape("simple-path-segment")
	}
}

func BenchmarkURL_PathEscape_WithSpecials(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		urlSinkString = URLPathEscape("path with spaces/and slashes")
	}
}
