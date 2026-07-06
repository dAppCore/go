// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the I/O primitives in io.go.
// Per AX-11 — io helpers sit on every payload-shaping path: HTTP
// body handling, file load/save, IPC, log collection. Even modest
// per-call overhead compounds across heavy workloads.
//
// The bench surface covers the constructors (NewBuffer / NewBufferString
// / NewBufferReader), the copy primitives (Copy / CopyN / WriteString),
// the LimitReader bound, and the ReadAll terminator. Fixture sizes span
// 1 KB → 64 KB → 1 MB to expose any size-sensitivity (Copy pre-sized
// vs growing, ReadAll buffer doubling, etc.).
//
// Run:    go test -bench='BenchmarkIO|BenchmarkCopy|BenchmarkReadAll' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// --- Fixtures ---

var (
	bench1KB  = makePayload(1 * 1024)
	bench64KB = makePayload(64 * 1024)
	bench1MB  = makePayload(1024 * 1024)

	bench1KBStr = string(bench1KB)
)

func makePayload(n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte('a' + (i % 26))
	}
	return out
}

// --- Buffer constructors ---

func BenchmarkNewBuffer_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewBuffer()
	}
}

func BenchmarkNewBuffer_WithBytes(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewBuffer(bench1KB)
	}
}

func BenchmarkNewBufferString(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewBufferString(bench1KBStr)
	}
}

func BenchmarkNewBufferReader(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewBufferReader(bench1KB)
	}
}

func BenchmarkNewReader(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewReader(bench1KBStr)
	}
}

// --- Copy ---

func BenchmarkCopy_1KB(b *B) {
	dst := NewBuffer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst.Reset()
		_ = Copy(dst, NewBufferReader(bench1KB))
	}
}

func BenchmarkCopy_64KB(b *B) {
	dst := NewBuffer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst.Reset()
		_ = Copy(dst, NewBufferReader(bench64KB))
	}
}

func BenchmarkCopyN_1KB(b *B) {
	dst := NewBuffer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst.Reset()
		_ = CopyN(dst, NewBufferReader(bench64KB), 1024)
	}
}

// --- WriteString ---

func BenchmarkWriteString_Short(b *B) {
	dst := NewBuffer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst.Reset()
		_ = WriteString(dst, "agent ready\n")
	}
}

func BenchmarkWriteString_1KB(b *B) {
	dst := NewBuffer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst.Reset()
		_ = WriteString(dst, bench1KBStr)
	}
}

// --- ReadAll ---

func BenchmarkReadAll_1KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ReadAll(NewBufferReader(bench1KB))
	}
}

func BenchmarkReadAll_64KB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ReadAll(NewBufferReader(bench64KB))
	}
}

func BenchmarkReadAll_1MB(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ReadAll(NewBufferReader(bench1MB))
	}
}

// --- LimitReader ---

func BenchmarkLimitReader_Construct(b *B) {
	r := NewBufferReader(bench1KB)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = LimitReader(r, 512)
	}
}

func BenchmarkLimitReader_ReadAll_Bounded(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ReadAll(LimitReader(NewBufferReader(bench64KB), 1024))
	}
}
