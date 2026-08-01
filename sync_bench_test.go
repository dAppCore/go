// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the sync primitives in sync.go.
// Per AX-11 — Mutex / RWMutex / Once / WaitGroup / SyncMap are stdlib
// passthroughs at the floor; the bench harness gates against accidental
// pessimisation and documents the contention floors that load-bearing
// consumer code is sensitive to.
//
// Run:    go test -bench='BenchmarkMutex|BenchmarkRWMutex|BenchmarkOnce|BenchmarkWaitGroup|BenchmarkSyncMap' -benchmem -run='^$' .

package core_test

import (
	"sync"
	"testing"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE on the read paths.
var (
	syncSinkAny    any
	syncSinkBool   bool
	syncSinkResult Result
)

// --- Mutex ---

func BenchmarkMutex_LockUnlock(b *B) {
	var m Mutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}

func BenchmarkMutex_TryLock(b *B) {
	var m Mutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkResult = m.TryLock()
		if syncSinkResult.OK {
			m.Unlock()
		}
	}
}

func BenchmarkMutex_LockUnlock_Parallel(b *B) {
	var m Mutex
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Lock()
			m.Unlock()
		}
	})
}

// --- RWMutex ---

func BenchmarkRWMutex_LockUnlock(b *B) {
	var m RWMutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}

func BenchmarkRWMutex_RLockRUnlock(b *B) {
	var m RWMutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.RLock()
		m.RUnlock()
	}
}

func BenchmarkRWMutex_TryLock(b *B) {
	var m RWMutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkResult = m.TryLock()
		if syncSinkResult.OK {
			m.Unlock()
		}
	}
}

func BenchmarkRWMutex_TryRLock(b *B) {
	var m RWMutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkResult = m.TryRLock()
		if syncSinkResult.OK {
			m.RUnlock()
		}
	}
}

func BenchmarkRWMutex_RLockRUnlock_Parallel(b *B) {
	var m RWMutex
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.RLock()
			m.RUnlock()
		}
	})
}

// --- Once ---

func BenchmarkOnce_Do_FirstCall(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var o Once
		o.Do(func() {})
	}
}

func BenchmarkOnce_Do_AfterFirst(b *B) {
	var o Once
	o.Do(func() {})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		o.Do(func() {})
	}
}

// --- WaitGroup ---

func BenchmarkWaitGroup_AddDone(b *B) {
	var w WaitGroup
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w.Add(1)
		w.Done()
	}
}

func BenchmarkWaitGroup_Go(b *B) {
	var w WaitGroup
	var done sync.WaitGroup
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		done.Add(1)
		w.Go(func() { done.Done() })
	}
	done.Wait()
}

// --- SyncMap ---

func BenchmarkSyncMap_Store(b *B) {
	var m SyncMap
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
}

func BenchmarkSyncMap_Load_Hit(b *B) {
	var m SyncMap
	for i := range 1000 {
		m.Store(i, i)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkAny, syncSinkBool = m.Load(i % 1000)
	}
}

func BenchmarkSyncMap_Load_Miss(b *B) {
	var m SyncMap
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkAny, syncSinkBool = m.Load(i)
	}
}

func BenchmarkSyncMap_LoadOrStore_Hit(b *B) {
	var m SyncMap
	m.Store("k", "v")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkAny, syncSinkBool = m.LoadOrStore("k", "v")
	}
}

func BenchmarkSyncMap_Delete(b *B) {
	var m SyncMap
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Delete(i)
	}
}

func BenchmarkSyncMap_Load_Parallel(b *B) {
	var m SyncMap
	for i := range 1000 {
		m.Store(i, i)
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var k int
		for pb.Next() {
			syncSinkAny, syncSinkBool = m.Load(k)
			k = (k + 1) % 1000
		}
	})
}

// --- SyncMap atomic ops ---

func BenchmarkSyncMap_Swap(b *B) {
	var m SyncMap
	m.Store("k", 0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkAny, syncSinkBool = m.Swap("k", i)
	}
}

func BenchmarkSyncMap_LoadAndDelete(b *B) {
	var m SyncMap
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkAny, syncSinkBool = m.LoadAndDelete(i)
	}
}

func BenchmarkSyncMap_CompareAndSwap_Hit(b *B) {
	var m SyncMap
	m.Store("k", 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkBool = m.CompareAndSwap("k", 1, 1)
	}
}

func BenchmarkSyncMap_CompareAndDelete_Miss(b *B) {
	var m SyncMap
	m.Store("k", 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		syncSinkBool = m.CompareAndDelete("k", 99)
	}
}

func BenchmarkSyncMap_Range_1000(b *B) {
	var m SyncMap
	for i := range 1000 {
		m.Store(i, i)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		count := 0
		m.Range(func(_, _ any) bool {
			count++
			return true
		})
		syncSinkAny = count
	}
}

func BenchmarkSyncMap_Clear_1000(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var m SyncMap
		for j := range 1000 {
			m.Store(j, j)
		}
		m.Clear()
	}
}

// --- Once.Reset ---

func BenchmarkOnce_Reset(b *B) {
	var once Once
	once.Do(func() {})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		once.Reset()
	}
}
