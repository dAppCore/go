// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the random primitives in random.go.
// Per AX-11 — random helpers split into two perf classes:
//
//   * RandIntn / RandPick — fast non-crypto math/rand/v2 path; sits on
//     every token-sampling step in inference loops (go-mlx), every
//     load-balancer pick, every shuffle. Must be tight.
//   * RandomBytes / RandomString / RandomInt / RandRead — crypto-secure
//     paths for nonces, IDs, secret material. OS-entropy bound but the
//     wrapper overhead and big.Int math in RandomInt are worth measuring.
//
// Run:    go test -bench='BenchmarkRand' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// --- Fast non-crypto path ---

func BenchmarkRandIntn_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandIntn(16)
	}
}

func BenchmarkRandIntn_Large(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandIntn(1 << 20)
	}
}

func BenchmarkRandPick_Strings(b *B) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandPick(items)
	}
}

func BenchmarkRandPick_Ints(b *B) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandPick(items)
	}
}

// --- Crypto-secure path ---

func BenchmarkRandomBytes_16(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomBytes(16)
	}
}

func BenchmarkRandomBytes_32(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomBytes(32)
	}
}

func BenchmarkRandomBytes_256(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomBytes(256)
	}
}

func BenchmarkRandomString_16(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomString(16)
	}
}

func BenchmarkRandomInt_SmallRange(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomInt(0, 100)
	}
}

func BenchmarkRandomInt_LargeRange(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandomInt(0, 1<<30)
	}
}

func BenchmarkRandRead_32(b *B) {
	buf := make([]byte, 32)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandRead(buf)
	}
}

func BenchmarkRandRead_256(b *B) {
	buf := make([]byte, 256)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandRead(buf)
	}
}
