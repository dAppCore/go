// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the encoding primitives in encode.go.
// Per AX-11 — hex + base64 are on every fingerprint emit, every API
// auth-header build, every model-hash compare, every keys.Pack write.
// The wrappers are encoding/{hex,base64} passthroughs; the bench
// harness gates the contract so a Core reroute can't slow them silently.
//
// Run:    go test -bench='BenchmarkEncode' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	encodeSinkString string
	encodeSinkResult Result
)

// --- Fixtures ---
//
// 16-byte payload mirrors SHA-256 truncated or a small auth token.
// 1KB payload covers larger blobs (model header chunk, jwt body).

var (
	encodeBytes16 = []byte("aaaaaaaaaaaaaaaa")
	encodeBytes1K = func() []byte {
		out := make([]byte, 1024)
		for i := range out {
			out[i] = byte('a' + i%26)
		}
		return out
	}()
)

// --- Hex ---

func BenchmarkEncode_HexEncode_16B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkString = HexEncode(encodeBytes16)
	}
}

func BenchmarkEncode_HexEncode_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkString = HexEncode(encodeBytes1K)
	}
}

func BenchmarkEncode_HexDecode_16B(b *B) {
	encoded := HexEncode(encodeBytes16)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkResult = HexDecode(encoded)
	}
}

func BenchmarkEncode_HexDecode_1KB(b *B) {
	encoded := HexEncode(encodeBytes1K)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkResult = HexDecode(encoded)
	}
}

// --- Base64 std ---

func BenchmarkEncode_Base64Encode_16B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkString = Base64Encode(encodeBytes16)
	}
}

func BenchmarkEncode_Base64Encode_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkString = Base64Encode(encodeBytes1K)
	}
}

func BenchmarkEncode_Base64Decode_16B(b *B) {
	encoded := Base64Encode(encodeBytes16)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkResult = Base64Decode(encoded)
	}
}

func BenchmarkEncode_Base64Decode_1KB(b *B) {
	encoded := Base64Encode(encodeBytes1K)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkResult = Base64Decode(encoded)
	}
}

// --- Base64 url ---

func BenchmarkEncode_Base64URLEncode_16B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkString = Base64URLEncode(encodeBytes16)
	}
}

func BenchmarkEncode_Base64URLDecode_16B(b *B) {
	encoded := Base64URLEncode(encodeBytes16)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encodeSinkResult = Base64URLDecode(encoded)
	}
}
