// SPDX-License-Identifier: EUPL-1.2

package core

func TestPool_Pool_Get_Good(t *T) {
	// A Put value comes back out of Get, and the pool empties.
	var pool Pool[*int]
	v := new(int)
	*v = 42
	pool.Put(v)

	got := pool.Get()

	AssertTrue(t, got == v)
	AssertEqual(t, 42, *got)
	AssertEqual(t, 0, pool.Len())
}

func TestPool_Pool_Get_Bad(t *T) {
	// Empty pool: the zero T (nil for a pointer) is the build-a-fresh-one sentinel.
	var pool Pool[*int]

	AssertNil(t, pool.Get())
	AssertEqual(t, 0, pool.Len())
}

func TestPool_Pool_Get_Ugly(t *T) {
	// Get drops the backing slot's reference, so the pool never pins a value
	// it has handed out (a device buffer the caller now owns).
	var pool Pool[*int]
	pool.Put(new(int))

	_ = pool.Get()

	AssertNil(t, pool.items[:1][0])
}

func TestPool_Pool_Put_Good(t *T) {
	// LIFO: the most recent Put comes out first, keeping the working set hot.
	var pool Pool[int]
	pool.Put(1)
	pool.Put(2)
	pool.Put(3)

	AssertEqual(t, 3, pool.Get())
	AssertEqual(t, 2, pool.Get())
	AssertEqual(t, 1, pool.Get())
}

func TestPool_Pool_Put_Bad(t *T) {
	// Pool is value-agnostic: a nil pointer is stored as handed (guarding is the caller's job).
	var pool Pool[*int]

	pool.Put(nil)

	AssertEqual(t, 1, pool.Len())
	AssertNil(t, pool.Get())
}

func TestPool_Pool_Put_Ugly(t *T) {
	// Concurrent Put/Get: every Put matched by a Get nets to empty (the race
	// detector is the real assertion; the count catches a lost update).
	var pool Pool[int]
	var wg WaitGroup
	for range 8 {
		wg.Go(func() {
			for i := range 1000 {
				pool.Put(i)
				pool.Get()
			}
		})
	}

	wg.Wait()

	AssertEqual(t, 0, pool.Len())
}

func TestPool_Pool_Len_Good(t *T) {
	var pool Pool[int]
	pool.Put(1)
	pool.Put(2)

	AssertEqual(t, 2, pool.Len())
}

func TestPool_Pool_Len_Bad(t *T) {
	// The zero Pool is ready to use and empty.
	var pool Pool[string]

	AssertEqual(t, 0, pool.Len())
}

func TestPool_Pool_Len_Ugly(t *T) {
	// Get on an empty pool doesn't go negative.
	var pool Pool[int]
	pool.Get()
	pool.Get()

	AssertEqual(t, 0, pool.Len())
}
