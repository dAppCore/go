// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the named-mutex primitives in lock.go.
// Per AX-11 — c.Lock(name) is the cooperative mutex consumers reach for
// when serialising drain, cache reload, and shared-resource access.
// The Get-or-create lookup happens per-call; the actual lock/unlock is
// stdlib sync.RWMutex underneath. Startables / Stoppables are boot-
// time queries that walk the service registry.
//
// Run:    go test -bench='BenchmarkLock' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	lockSinkLock   *Lock
	lockSinkResult Result
)

// --- Lock lookup / acquire / release ---

func BenchmarkLock_Lock_Lookup(b *B) {
	c := New()
	c.Lock("warm") // pre-create
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lockSinkLock = c.Lock("warm")
	}
}

func BenchmarkLock_Lock_LookupOrCreate(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Each iteration creates a fresh lock entry — gates the Set path.
		lockSinkLock = c.Lock(Sprintf("lock.%d", i))
	}
}

func BenchmarkLock_Lock_Unlock(b *B) {
	c := New()
	l := c.Lock("warm")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkLock_RLock_RUnlock(b *B) {
	c := New()
	l := c.Lock("warm")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.RLock()
		l.RUnlock()
	}
}

func BenchmarkLock_TryLock_Free(b *B) {
	c := New()
	l := c.Lock("warm")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := l.TryLock()
		if r.OK {
			l.Unlock()
		}
		lockSinkResult = r
	}
}

// --- LockEnable + LockApply ---

func BenchmarkLock_LockEnable(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.LockEnable()
	}
}

func BenchmarkLock_LockApply_Enabled(b *B) {
	c := New()
	c.LockEnable()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.LockApply()
	}
}

// --- Startables / Stoppables ---

func BenchmarkLock_Startables_Empty(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lockSinkResult = c.Startables()
	}
}

func BenchmarkLock_Startables_FiveRegistered(b *B) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Service(Sprintf("bench.%d", i), Service{
			Name:    Sprintf("bench.%d", i),
			OnStart: func() Result { return Result{OK: true} },
		})
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lockSinkResult = c.Startables()
	}
}

func BenchmarkLock_Stoppables_FiveRegistered(b *B) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Service(Sprintf("bench.%d", i), Service{
			Name:   Sprintf("bench.%d", i),
			OnStop: func() Result { return Result{OK: true} },
		})
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lockSinkResult = c.Stoppables()
	}
}
