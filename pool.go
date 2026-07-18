// SPDX-License-Identifier: EUPL-1.2

// Generic object pool for the Core framework — a mutex-guarded LIFO free-list
// that recycles expensive-to-build values across calls without per-call
// allocation.

package core

import "sync"

// Pool is a concurrency-safe LIFO free-list of reusable T values — the zero-alloc
// recycling primitive behind hot-path scratch buffers and device handles. Get
// pops a recycled value (or the zero T when empty, so the caller builds a fresh
// one); Put returns one for reuse. Unlike sync.Pool it never drops entries on GC,
// so a warm pool stays warm — the right trade when the pooled value owns a
// resource (a device buffer, a mapping) the caller Closes explicitly rather than
// letting the GC reclaim.
//
// The zero Pool is ready to use. T is typically a pointer, so Get's empty
// sentinel (the zero T == nil) reads naturally at the call site:
//
//	var pool core.Pool[*scratch]
//	s := pool.Get()
//	if s == nil {
//	    s = newScratch()
//	}
//	defer pool.Put(s)
type Pool[T any] struct {
	mu    sync.Mutex
	items []T
}

// Get pops and returns the most-recently-Put value, or the zero T when the pool
// is empty (the caller then builds a fresh one). LIFO keeps a small working set
// hot.
//
//	s := pool.Get()
//	if s == nil {
//	    s = newScratch()
//	}
func (p *Pool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := len(p.items)
	if n == 0 {
		var zero T
		return zero
	}
	s := p.items[n-1]
	var zero T
	p.items[n-1] = zero // drop the reference so the pool doesn't pin a handed-out value
	p.items = p.items[:n-1]
	return s
}

// Put returns a value to the pool for reuse. The caller guards against pooling an
// invalid value (a nil pointer, a closed buffer) before calling Put — Pool is
// value-agnostic and stores whatever it is handed.
//
//	pool.Put(s)
func (p *Pool[T]) Put(s T) {
	p.mu.Lock()
	p.items = append(p.items, s)
	p.mu.Unlock()
}

// Len reports how many values are currently pooled — for diagnostics and tests.
//
//	core.Println(pool.Len())
func (p *Pool[T]) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.items)
}
