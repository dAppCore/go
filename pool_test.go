// SPDX-License-Identifier: EUPL-1.2

package core

import (
	"sync"
	"testing"
)

// TestPoolGetEmptyReturnsZero: Get on an empty pool yields the zero T (nil for a
// pointer element), the sentinel the caller reads to build a fresh value.
func TestPoolGetEmptyReturnsZero(t *testing.T) {
	var pool Pool[*int]
	if got := pool.Get(); got != nil {
		t.Fatalf("Get on empty pool = %v, want nil", got)
	}
	if pool.Len() != 0 {
		t.Fatalf("Len on empty pool = %d, want 0", pool.Len())
	}
}

// TestPoolPutGetRoundTrip: a Put value comes back out of Get.
func TestPoolPutGetRoundTrip(t *testing.T) {
	var pool Pool[*int]
	v := new(int)
	*v = 42
	pool.Put(v)
	if pool.Len() != 1 {
		t.Fatalf("Len after one Put = %d, want 1", pool.Len())
	}
	got := pool.Get()
	if got != v || *got != 42 {
		t.Fatalf("Get = %v, want the pooled %v", got, v)
	}
	if pool.Len() != 0 {
		t.Fatalf("Len after Get = %d, want 0", pool.Len())
	}
}

// TestPoolIsLIFO: the most-recently-Put value is returned first, keeping a small
// working set hot.
func TestPoolIsLIFO(t *testing.T) {
	var pool Pool[int]
	pool.Put(1)
	pool.Put(2)
	pool.Put(3)
	for _, want := range []int{3, 2, 1} {
		if got := pool.Get(); got != want {
			t.Fatalf("LIFO Get = %d, want %d", got, want)
		}
	}
}

// TestPoolGetDropsReference: Get must not keep pinning a handed-out value (a
// leak of a device buffer the caller now owns). After Get empties the pool, the
// backing slot no longer references the value.
func TestPoolGetDropsReference(t *testing.T) {
	var pool Pool[*int]
	v := new(int)
	pool.Put(v)
	_ = pool.Get()
	// The pool is empty and holds no reference; a second Get yields nil, not v.
	if got := pool.Get(); got != nil {
		t.Fatalf("second Get = %v, want nil (pool must not retain a handed-out value)", got)
	}
}

// TestPoolConcurrent: Get/Put are safe under concurrent use (race detector is the
// real assertion; the count check catches a lost update).
func TestPoolConcurrent(t *testing.T) {
	var pool Pool[int]
	const workers, each = 8, 1000
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < each; i++ {
				pool.Put(i)
				pool.Get()
			}
		}()
	}
	wg.Wait()
	// Every Put is matched by a Get, so the pool nets to empty.
	if pool.Len() != 0 {
		t.Fatalf("Len after balanced concurrent Put/Get = %d, want 0", pool.Len())
	}
}
