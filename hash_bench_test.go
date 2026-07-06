// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the hash primitives in hash.go.
// Per AX-11 — SHA256 / HMAC / HKDF sit on every config-integrity
// check, every cache key, every signed-payload validation. SHA256
// over small inputs is a fingerprint primitive; over MB-class
// inputs it's a model-integrity primitive.
//
// Run:    go test -bench='BenchmarkSHA|BenchmarkHMAC|BenchmarkHKDF' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	hashSinkArr32  [32]byte
	hashSinkString string
	hashSinkResult Result
)

// Fixtures
var (
	hashSmall  = []byte("hello")
	hashMedium = []byte("the quick brown fox jumps over the lazy dog")
	hash1KB    = makeBytes(1024)
	hash64KB   = makeBytes(64 * 1024)
	hash1MB    = makeBytes(1024 * 1024)

	hashSmallStr  = string(hashSmall)
	hashMediumStr = string(hashMedium)

	hashKey  = []byte("secret-key-32-bytes-for-test----")
	hashSalt = []byte("salt-16-bytes---")
	hashInfo = []byte("session")
)

func makeBytes(n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte('a' + (i % 26))
	}
	return out
}

// --- SHA256 ---

func BenchmarkSHA256_5B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256(hashSmall)
	}
}

func BenchmarkSHA256_43B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256(hashMedium)
	}
}

func BenchmarkSHA256_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256(hash1KB)
	}
}

func BenchmarkSHA256_64KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256(hash64KB)
	}
}

func BenchmarkSHA256_1MB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256(hash1MB)
	}
}

// --- SHA256 Hex variants ---

func BenchmarkSHA256Hex_43B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkString = SHA256Hex(hashMedium)
	}
}

// --- SHA256String / SHA256HexString ---

func BenchmarkSHA256String_5B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256String(hashSmallStr)
	}
}

func BenchmarkSHA256String_43B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkArr32 = SHA256String(hashMediumStr)
	}
}

func BenchmarkSHA256HexString_43B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkString = SHA256HexString(hashMediumStr)
	}
}

// --- HMAC ---

func BenchmarkHMAC_SHA256_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkResult = HMAC("sha256", hashKey, hashMedium)
	}
}

func BenchmarkHMAC_SHA256_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkResult = HMAC("sha256", hashKey, hash1KB)
	}
}

func BenchmarkHMAC_SHA512_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkResult = HMAC("sha512", hashKey, hashMedium)
	}
}

// --- HKDF ---

func BenchmarkHKDF_SHA256_32B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkResult = HKDF("sha256", hashKey, hashSalt, hashInfo, 32)
	}
}

func BenchmarkHKDF_SHA256_256B(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hashSinkResult = HKDF("sha256", hashKey, hashSalt, hashInfo, 256)
	}
}
