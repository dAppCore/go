// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the SHA-3 primitives in sha3.go.
// Per AX-11 — SHA3 + Keccak land on every drive-handle fingerprint,
// every Web3 address recovery, every R1 record content hash. Keccak256
// in particular is hot — it's the in-house pure-Go implementation
// (no Cgo) consumers reach for over the deprecated x/crypto/sha3 path.
// SHAKE variants are XOF outputs used by KDF + random-seed expansion.
//
// Run:    go test -bench='BenchmarkSHA3' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	sha3SinkSum32 [32]byte
	sha3SinkBytes []byte
	sha3SinkStr   string
)

// Fixtures — three size points to expose per-byte costs.
// makeBytes lives in hash_bench_test.go (same package); reuse it.
var (
	sha3Bytes16  = []byte("aaaaaaaaaaaaaaaa")
	sha3Bytes1K  = makeBytes(1024)
	sha3Bytes64K = makeBytes(64 * 1024)
)

// --- SHA3-256 (FIPS-202) ---

func BenchmarkSHA3_256_16B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkSum32 = SHA3_256(sha3Bytes16)
	}
}

func BenchmarkSHA3_256_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkSum32 = SHA3_256(sha3Bytes1K)
	}
}

func BenchmarkSHA3_256_64KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkSum32 = SHA3_256(sha3Bytes64K)
	}
}

func BenchmarkSHA3_256Hex_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkStr = SHA3_256Hex(sha3Bytes1K)
	}
}

// --- Keccak-256 (legacy pre-NIST) ---

func BenchmarkSHA3_Keccak256_16B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkSum32 = Keccak256(sha3Bytes16)
	}
}

func BenchmarkSHA3_Keccak256_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkSum32 = Keccak256(sha3Bytes1K)
	}
}

func BenchmarkSHA3_Keccak256_64KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkSum32 = Keccak256(sha3Bytes64K)
	}
}

func BenchmarkSHA3_Keccak256Hex_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkStr = Keccak256Hex(sha3Bytes1K)
	}
}

// --- SHAKE128 / SHAKE256 (XOF) ---

func BenchmarkSHA3_Shake128_32(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkBytes = SHA3Shake128(sha3Bytes16, 32)
	}
}

func BenchmarkSHA3_Shake256_64(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sha3SinkBytes = SHA3Shake256(sha3Bytes16, 64)
	}
}
