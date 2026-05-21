// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the atomic primitives in atomic.go.
// Per AX-11 — atomic ops sit on every counter, every ready-flag,
// every shared registry slot. Wrappers add no measurable overhead
// in normal use (compiler inlines through to the stdlib atomic.*
// intrinsic) but a regression here would cascade across every
// concurrent path in the ecosystem, so each one gets a perf gate.
//
// Each typed wrapper (AtomicBool / Int32 / Int64 / Uint32 / Uint64 /
// Pointer[T]) is covered for Load / Store / Add / Swap / CAS.
// Parallel variants on the int64 counter expose the contention floor.
//
// Run:    go test -bench='BenchmarkAtomic' -benchmem -run='^$' .

package core_test

import (
	"testing"

	. "dappco.re/go"
)

// Sinks prevent compiler dead-code elimination on the read paths.
var (
	atomicSinkBool    bool
	atomicSinkInt32   int32
	atomicSinkInt64   int64
	atomicSinkUint32  uint32
	atomicSinkUint64  uint64
	atomicSinkPointer *atomicBenchPayload
)

type atomicBenchPayload struct {
	value int64
}

// --- AtomicBool ---

func BenchmarkAtomicBool_Load(b *B) {
	var a AtomicBool
	a.Store(true)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkBool = a.Load()
	}
}

func BenchmarkAtomicBool_Store(b *B) {
	var a AtomicBool
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Store(i&1 == 0)
	}
}

func BenchmarkAtomicBool_Swap(b *B) {
	var a AtomicBool
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkBool = a.Swap(i&1 == 0)
	}
}

func BenchmarkAtomicBool_CompareAndSwap(b *B) {
	var a AtomicBool
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwap(false, true)
		a.Store(false)
	}
}

// --- AtomicInt32 ---

func BenchmarkAtomicInt32_Load(b *B) {
	var a AtomicInt32
	a.Store(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkInt32 = a.Load()
	}
}

func BenchmarkAtomicInt32_Store(b *B) {
	var a AtomicInt32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Store(int32(i))
	}
}

func BenchmarkAtomicInt32_Add(b *B) {
	var a AtomicInt32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkInt32 = a.Add(1)
	}
}

func BenchmarkAtomicInt32_Swap(b *B) {
	var a AtomicInt32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkInt32 = a.Swap(int32(i))
	}
}

func BenchmarkAtomicInt32_CompareAndSwap(b *B) {
	var a AtomicInt32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwap(0, 1)
		a.Store(0)
	}
}

// --- AtomicInt64 ---

func BenchmarkAtomicInt64_Load(b *B) {
	var a AtomicInt64
	a.Store(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkInt64 = a.Load()
	}
}

func BenchmarkAtomicInt64_Store(b *B) {
	var a AtomicInt64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Store(int64(i))
	}
}

func BenchmarkAtomicInt64_Add(b *B) {
	var a AtomicInt64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkInt64 = a.Add(1)
	}
}

func BenchmarkAtomicInt64_Swap(b *B) {
	var a AtomicInt64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkInt64 = a.Swap(int64(i))
	}
}

func BenchmarkAtomicInt64_CompareAndSwap(b *B) {
	var a AtomicInt64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwap(0, 1)
		a.Store(0)
	}
}

// --- AtomicUint32 ---

func BenchmarkAtomicUint32_Load(b *B) {
	var a AtomicUint32
	a.Store(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkUint32 = a.Load()
	}
}

func BenchmarkAtomicUint32_Add(b *B) {
	var a AtomicUint32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkUint32 = a.Add(1)
	}
}

func BenchmarkAtomicUint32_Swap(b *B) {
	var a AtomicUint32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkUint32 = a.Swap(uint32(i))
	}
}

func BenchmarkAtomicUint32_CompareAndSwap(b *B) {
	var a AtomicUint32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwap(0, 1)
		a.Store(0)
	}
}

// --- AtomicUint64 ---

func BenchmarkAtomicUint64_Load(b *B) {
	var a AtomicUint64
	a.Store(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkUint64 = a.Load()
	}
}

func BenchmarkAtomicUint64_Add(b *B) {
	var a AtomicUint64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkUint64 = a.Add(1)
	}
}

func BenchmarkAtomicUint64_Swap(b *B) {
	var a AtomicUint64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkUint64 = a.Swap(uint64(i))
	}
}

func BenchmarkAtomicUint64_CompareAndSwap(b *B) {
	var a AtomicUint64
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwap(0, 1)
		a.Store(0)
	}
}

// --- AtomicPointer[T] ---

func BenchmarkAtomicPointer_Load(b *B) {
	var a AtomicPointer[atomicBenchPayload]
	a.Store(&atomicBenchPayload{value: 42})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkPointer = a.Load()
	}
}

func BenchmarkAtomicPointer_Store(b *B) {
	var a AtomicPointer[atomicBenchPayload]
	p := &atomicBenchPayload{value: 42}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Store(p)
	}
}

func BenchmarkAtomicPointer_Swap(b *B) {
	var a AtomicPointer[atomicBenchPayload]
	p := &atomicBenchPayload{value: 42}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomicSinkPointer = a.Swap(p)
	}
}

func BenchmarkAtomicPointer_CompareAndSwap(b *B) {
	var a AtomicPointer[atomicBenchPayload]
	old := &atomicBenchPayload{value: 1}
	next := &atomicBenchPayload{value: 2}
	a.Store(old)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwap(old, next)
		a.Store(old)
	}
}

// --- Parallel contention floor ---
//
// AtomicInt64.Add under -cpu=32 contention shows the cache-line bounce
// cost — the realistic floor for any "shared counter" pattern.

func BenchmarkAtomicInt64_Add_Parallel(b *B) {
	var a AtomicInt64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Add(1)
		}
	})
}

func BenchmarkAtomicInt64_Load_Parallel(b *B) {
	var a AtomicInt64
	a.Store(42)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomicSinkInt64 = a.Load()
		}
	})
}
