package core_test

import . "dappco.re/go"

// --- Mutex ---

func TestSync_Mutex_Good(t *T) {
	var m Mutex
	m.Lock()
	m.Unlock()

	// Unlock actually released it — a fresh acquisition succeeds.
	r := m.TryLock()
	AssertTrue(t, r.OK)
	m.Unlock()
}

func TestSync_Mutex_Bad(t *T) {
	// Bad: TryLock on already-held mutex returns Result{OK: false}.
	var m Mutex
	m.Lock()
	defer m.Unlock()
	r := m.TryLock()
	AssertFalse(t, r.OK)
}

func TestSync_Mutex_Ugly(t *T) {
	// Ugly: contention. Two goroutines incrementing under the same Mutex
	// must produce 1000 increments without races (-race must pass).
	var m Mutex
	count := 0
	var wg WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 500 {
				m.Lock()
				count++
				m.Unlock()
			}
		}()
	}
	wg.Wait()
	AssertEqual(t, 1000, count)
}

// --- RWMutex ---

func TestSync_RWMutex_Good(t *T) {
	var m RWMutex
	m.Lock()
	m.Unlock()
	m.RLock()
	m.RUnlock()

	// Both the write and read lock cycles fully released — a fresh
	// write acquisition succeeds.
	r := m.TryLock()
	AssertTrue(t, r.OK)
	m.Unlock()
}

func TestSync_RWMutex_Bad(t *T) {
	// Bad: TryLock fails when write-held.
	var m RWMutex
	m.Lock()
	defer m.Unlock()
	r := m.TryLock()
	AssertFalse(t, r.OK)
}

func TestSync_RWMutex_Ugly(t *T) {
	// Ugly: many readers + occasional writer; -race must remain clean.
	var m RWMutex
	value := 0
	var wg WaitGroup
	// 5 readers
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				m.RLock()
				_ = value
				m.RUnlock()
			}
		}()
	}
	// 2 writers
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 50 {
				m.Lock()
				value++
				m.Unlock()
			}
		}()
	}
	wg.Wait()
	AssertEqual(t, 100, value)
}

func TestSync_RWMutex_TryRLock_Good(t *T) {
	var m RWMutex
	r := m.TryRLock()
	AssertTrue(t, r.OK)
	m.RUnlock()
}

// --- Once ---

func TestSync_Once_Good(t *T) {
	var o Once
	count := 0
	o.Do(func() { count++ })
	o.Do(func() { count++ })
	o.Do(func() { count++ })
	AssertEqual(t, 1, count, "Once.Do must execute the function exactly once")
}

func TestSync_Once_Bad(t *T) {
	// Bad: caller passes nil. Stdlib Once panics on nil; we pass through.
	var o Once
	AssertPanics(t, func() { o.Do(nil) })
}

func TestSync_Once_Ugly(t *T) {
	// Ugly: Reset between invocations re-arms the Once.
	var o Once
	count := 0
	o.Do(func() { count++ })
	o.Do(func() { count++ })
	AssertEqual(t, 1, count)
	o.Reset()
	o.Do(func() { count++ })
	o.Do(func() { count++ })
	AssertEqual(t, 2, count, "After Reset, Do must fire once more")
}

// --- WaitGroup ---

func TestSync_WaitGroup_Good(t *T) {
	var wg WaitGroup
	var mu Mutex
	done := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		Sleep(10 * Millisecond)
		mu.Lock()
		done = true
		mu.Unlock()
	}()
	wg.Wait()
	mu.Lock()
	defer mu.Unlock()
	AssertTrue(t, done)
}

func TestSync_WaitGroup_Bad(t *T) {
	// Bad: Done called more times than Add. Stdlib panics; we pass through.
	var wg WaitGroup
	wg.Add(1)
	wg.Done()
	AssertPanics(t, func() { wg.Done() })
}

func TestSync_WaitGroup_Ugly(t *T) {
	// Ugly: many goroutines, all must complete before Wait returns.
	var wg WaitGroup
	var mu Mutex
	counter := 0
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	mu.Lock()
	defer mu.Unlock()
	AssertEqual(t, 100, counter)
}

func TestSync_Mutex_Lock_Good(t *T) {
	var mu Mutex

	mu.Lock()
	mu.Unlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_Mutex_Lock_Bad(t *T) {
	var mu Mutex
	mu.Lock()
	defer mu.Unlock()

	r := mu.TryLock()

	AssertFalse(t, r.OK)
}

func TestSync_Mutex_Lock_Ugly(t *T) {
	var mu Mutex
	count := 0
	var wg WaitGroup

	for range 16 {
		wg.Go(func() {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}
	wg.Wait()

	AssertEqual(t, 16, count)
}

func TestSync_Mutex_Unlock_Good(t *T) {
	var mu Mutex

	mu.Lock()
	mu.Unlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_Mutex_Unlock_Bad(t *T) {
	var mu Mutex
	mu.Lock()

	// While held, a second acquisition fails; Unlock is what frees it.
	AssertFalse(t, mu.TryLock().OK)
	mu.Unlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_Mutex_Unlock_Ugly(t *T) {
	var mu Mutex

	for range 3 {
		mu.Lock()
		mu.Unlock()
	}

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_Mutex_TryLock_Good(t *T) {
	var mu Mutex

	r := mu.TryLock()

	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_Mutex_TryLock_Bad(t *T) {
	var mu Mutex
	mu.Lock()
	defer mu.Unlock()

	r := mu.TryLock()

	AssertFalse(t, r.OK)
}

func TestSync_Mutex_TryLock_Ugly(t *T) {
	var mu Mutex

	first := mu.TryLock()
	second := mu.TryLock()

	AssertTrue(t, first.OK)
	AssertFalse(t, second.OK)
	mu.Unlock()
}

func TestSync_RWMutex_Lock_Good(t *T) {
	var mu RWMutex

	mu.Lock()
	mu.Unlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_Lock_Bad(t *T) {
	var mu RWMutex
	mu.RLock()
	defer mu.RUnlock()

	r := mu.TryLock()

	AssertFalse(t, r.OK)
}

func TestSync_RWMutex_Lock_Ugly(t *T) {
	var mu RWMutex
	count := 0
	var wg WaitGroup

	for range 8 {
		wg.Go(func() {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}
	wg.Wait()

	AssertEqual(t, 8, count)
}

func TestSync_RWMutex_Unlock_Good(t *T) {
	var mu RWMutex

	mu.Lock()
	mu.Unlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_Unlock_Bad(t *T) {
	var mu RWMutex
	mu.Lock()

	// A write-lock blocks another writer until Unlock releases it.
	AssertFalse(t, mu.TryLock().OK)
	mu.Unlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_Unlock_Ugly(t *T) {
	var mu RWMutex

	for range 3 {
		mu.Lock()
		mu.Unlock()
	}

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_RLock_Good(t *T) {
	var mu RWMutex

	mu.RLock()
	mu.RUnlock()

	// RUnlock actually released the read lock — a writer can proceed.
	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_RLock_Bad(t *T) {
	var mu RWMutex
	mu.Lock()
	defer mu.Unlock()

	r := mu.TryRLock()

	AssertFalse(t, r.OK)
}

func TestSync_RWMutex_RLock_Ugly(t *T) {
	var mu RWMutex

	mu.RLock()
	mu.RLock()
	mu.RUnlock()
	mu.RUnlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_RUnlock_Good(t *T) {
	var mu RWMutex

	mu.RLock()
	mu.RUnlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_RUnlock_Bad(t *T) {
	var mu RWMutex
	mu.RLock()

	// A read-lock blocks a writer until RUnlock releases it.
	AssertFalse(t, mu.TryLock().OK)
	mu.RUnlock()

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_RUnlock_Ugly(t *T) {
	var mu RWMutex

	for range 3 {
		mu.RLock()
		mu.RUnlock()
	}

	r := mu.TryLock()
	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_TryLock_Good(t *T) {
	var mu RWMutex

	r := mu.TryLock()

	AssertTrue(t, r.OK)
	mu.Unlock()
}

func TestSync_RWMutex_TryLock_Bad(t *T) {
	var mu RWMutex
	mu.Lock()
	defer mu.Unlock()

	r := mu.TryLock()

	AssertFalse(t, r.OK)
}

func TestSync_RWMutex_TryLock_Ugly(t *T) {
	var mu RWMutex
	mu.RLock()
	defer mu.RUnlock()

	r := mu.TryLock()

	AssertFalse(t, r.OK)
}

func TestSync_RWMutex_TryRLock_Bad(t *T) {
	var mu RWMutex
	mu.Lock()
	defer mu.Unlock()

	r := mu.TryRLock()

	AssertFalse(t, r.OK)
}

func TestSync_RWMutex_TryRLock_Ugly(t *T) {
	var mu RWMutex

	first := mu.TryRLock()
	second := mu.TryRLock()

	AssertTrue(t, first.OK)
	AssertTrue(t, second.OK)
	mu.RUnlock()
	mu.RUnlock()
}

func TestSync_Once_Do_Good(t *T) {
	var once Once
	count := 0

	once.Do(func() { count++ })
	once.Do(func() { count++ })

	AssertEqual(t, 1, count)
}

func TestSync_Once_Do_Bad(t *T) {
	var once Once

	AssertPanics(t, func() { once.Do(nil) })
}

func TestSync_Once_Do_Ugly(t *T) {
	var once Once
	var count AtomicInt32
	var wg WaitGroup

	for range 16 {
		wg.Go(func() {
			once.Do(func() { count.Add(1) })
		})
	}
	wg.Wait()

	AssertEqual(t, int32(1), count.Load())
}

func TestSync_Once_Reset_Good(t *T) {
	var once Once
	count := 0

	once.Do(func() { count++ })
	once.Reset()
	once.Do(func() { count++ })

	AssertEqual(t, 2, count)
}

func TestSync_Once_Reset_Bad(t *T) {
	var once Once

	once.Reset()
	ran := false
	once.Do(func() { ran = true })

	AssertTrue(t, ran) // Do runs again after Reset
}

func TestSync_Once_Reset_Ugly(t *T) {
	var once Once
	count := 0

	for range 3 {
		once.Do(func() { count++ })
		once.Reset()
	}

	AssertEqual(t, 3, count)
}

func TestSync_WaitGroup_Add_Good(t *T) {
	var wg WaitGroup
	done := make(chan bool, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		done <- true
	}()
	wg.Wait()

	AssertTrue(t, <-done)
}

func TestSync_WaitGroup_Add_Bad(t *T) {
	var wg WaitGroup

	AssertPanics(t, func() { wg.Add(-1) })
}

func TestSync_WaitGroup_Add_Ugly(t *T) {
	var wg WaitGroup
	var done AtomicBool

	wg.Add(0) // a zero delta leaves the group immediately ready
	wg.Go(func() { done.Store(true) })
	wg.Wait()

	AssertTrue(t, done.Load())
}

func TestSync_WaitGroup_Done_Good(t *T) {
	var wg WaitGroup
	var count AtomicInt32

	wg.Add(1)
	go func() { count.Add(1); wg.Done() }()
	wg.Wait()

	AssertEqual(t, int32(1), count.Load()) // Wait returns after Done; the increment is visible
}

func TestSync_WaitGroup_Done_Bad(t *T) {
	var wg WaitGroup
	wg.Add(1)
	wg.Done()

	AssertPanics(t, func() { wg.Done() })
}

func TestSync_WaitGroup_Done_Ugly(t *T) {
	var wg WaitGroup
	var count AtomicInt32

	for range 3 {
		wg.Add(1)
		go func() { count.Add(1); wg.Done() }()
	}
	wg.Wait()

	AssertEqual(t, int32(3), count.Load())
}

func TestSync_WaitGroup_Wait_Good(t *T) {
	var wg WaitGroup
	var done AtomicBool

	wg.Go(func() { done.Store(true) })
	wg.Wait()

	AssertTrue(t, done.Load())
}

func TestSync_WaitGroup_Wait_Bad(t *T) {
	var wg WaitGroup
	done := make(chan bool, 1)

	go func() { wg.Wait(); done <- true }()

	select {
	case <-done:
	case <-After(2 * Second):
		t.Fatal("Wait blocked on a zero-counter group")
	}
}

func TestSync_WaitGroup_Wait_Ugly(t *T) {
	var wg WaitGroup
	var count AtomicInt32

	for range 32 {
		wg.Go(func() { count.Add(1) })
	}
	wg.Wait()

	AssertEqual(t, int32(32), count.Load())
}

func TestSync_WaitGroup_Go_Good(t *T) {
	var wg WaitGroup
	var ran AtomicBool

	wg.Go(func() { ran.Store(true) })
	wg.Wait()

	AssertTrue(t, ran.Load())
}

func TestSync_WaitGroup_Go_Bad(t *T) {
	var wg WaitGroup
	var ran AtomicBool

	wg.Go(func() { ran.Store(false) })
	wg.Wait()

	AssertFalse(t, ran.Load())
}

func TestSync_WaitGroup_Go_Ugly(t *T) {
	var wg WaitGroup
	var count AtomicInt32

	for range 64 {
		wg.Go(func() { count.Add(1) })
	}
	wg.Wait()

	AssertEqual(t, int32(64), count.Load())
}

func TestSync_SyncMap_Load_Good(t *T) {
	var cache SyncMap
	cache.Store("agent.dispatch.status", "ready")

	value, ok := cache.Load("agent.dispatch.status")

	AssertTrue(t, ok)
	AssertEqual(t, "ready", value)
}

func TestSync_SyncMap_Load_Bad(t *T) {
	var cache SyncMap

	value, ok := cache.Load("agent.dispatch.status")

	AssertFalse(t, ok)
	AssertNil(t, value)
}

func TestSync_SyncMap_Load_Ugly(t *T) {
	var cache SyncMap
	cache.Store("agent.dispatch.status", nil)

	value, ok := cache.Load("agent.dispatch.status")

	AssertTrue(t, ok)
	AssertNil(t, value)
}

func TestSync_SyncMap_Store_Good(t *T) {
	var cache SyncMap

	cache.Store("agent.dispatch.status", "ready")
	value, ok := cache.Load("agent.dispatch.status")

	AssertTrue(t, ok)
	AssertEqual(t, "ready", value)
}

func TestSync_SyncMap_Store_Bad(t *T) {
	var cache SyncMap
	cache.Store("agent.dispatch.status", "starting")

	cache.Store("agent.dispatch.status", "ready")
	value, ok := cache.Load("agent.dispatch.status")

	AssertTrue(t, ok)
	AssertEqual(t, "ready", value)
}

func TestSync_SyncMap_Store_Ugly(t *T) {
	var cache SyncMap
	key := struct {
		agent string
		shard int
	}{agent: "dispatch", shard: 0}

	cache.Store(key, "ready")
	value, ok := cache.Load(key)

	AssertTrue(t, ok)
	AssertEqual(t, "ready", value)
}

func TestSync_SyncMap_LoadOrStore_Good(t *T) {
	var cache SyncMap

	value, loaded := cache.LoadOrStore("session.token", "fresh")

	AssertFalse(t, loaded)
	AssertEqual(t, "fresh", value)
}

func TestSync_SyncMap_LoadOrStore_Bad(t *T) {
	var cache SyncMap
	cache.Store("session.token", "existing")

	value, loaded := cache.LoadOrStore("session.token", "fresh")

	AssertTrue(t, loaded)
	AssertEqual(t, "existing", value)
}

func TestSync_SyncMap_LoadOrStore_Ugly(t *T) {
	var cache SyncMap

	value, loaded := cache.LoadOrStore("session.token", nil)

	AssertFalse(t, loaded)
	AssertNil(t, value)
}

func TestSync_SyncMap_LoadAndDelete_Good(t *T) {
	var cache SyncMap
	cache.Store("session.token", "fresh")

	value, loaded := cache.LoadAndDelete("session.token")
	_, ok := cache.Load("session.token")

	AssertTrue(t, loaded)
	AssertEqual(t, "fresh", value)
	AssertFalse(t, ok)
}

func TestSync_SyncMap_LoadAndDelete_Bad(t *T) {
	var cache SyncMap

	value, loaded := cache.LoadAndDelete("session.token")

	AssertFalse(t, loaded)
	AssertNil(t, value)
}

func TestSync_SyncMap_LoadAndDelete_Ugly(t *T) {
	var cache SyncMap
	cache.Store("session.token", nil)

	value, loaded := cache.LoadAndDelete("session.token")

	AssertTrue(t, loaded)
	AssertNil(t, value)
}

func TestSync_SyncMap_Delete_Good(t *T) {
	var cache SyncMap
	cache.Store("agent", "online")

	cache.Delete("agent")
	_, ok := cache.Load("agent")

	AssertFalse(t, ok)
}

func TestSync_SyncMap_Delete_Bad(t *T) {
	var cache SyncMap

	cache.Delete("agent") // deleting an absent key is a safe no-op

	_, ok := cache.Load("agent")
	AssertFalse(t, ok)
}

func TestSync_SyncMap_Delete_Ugly(t *T) {
	var cache SyncMap
	cache.Store("agent", "online")

	cache.Delete("agent")
	cache.Delete("agent")
	_, ok := cache.Load("agent")

	AssertFalse(t, ok)
}

func TestSync_SyncMap_Swap_Good(t *T) {
	var cache SyncMap
	cache.Store("agent", "starting")

	previous, loaded := cache.Swap("agent", "ready")
	current, ok := cache.Load("agent")

	AssertTrue(t, loaded)
	AssertEqual(t, "starting", previous)
	AssertTrue(t, ok)
	AssertEqual(t, "ready", current)
}

func TestSync_SyncMap_Swap_Bad(t *T) {
	var cache SyncMap

	previous, loaded := cache.Swap("agent", "ready")

	AssertFalse(t, loaded)
	AssertNil(t, previous)
}

func TestSync_SyncMap_Swap_Ugly(t *T) {
	var cache SyncMap
	cache.Store("agent", nil)

	previous, loaded := cache.Swap("agent", "ready")

	AssertTrue(t, loaded)
	AssertNil(t, previous)
}

func TestSync_SyncMap_CompareAndSwap_Good(t *T) {
	var cache SyncMap
	cache.Store("agent", "starting")

	swapped := cache.CompareAndSwap("agent", "starting", "ready")

	AssertTrue(t, swapped)
	value, ok := cache.Load("agent")
	AssertTrue(t, ok)
	AssertEqual(t, "ready", value)
}

func TestSync_SyncMap_CompareAndSwap_Bad(t *T) {
	var cache SyncMap
	cache.Store("agent", "starting")

	swapped := cache.CompareAndSwap("agent", "stale", "ready")

	AssertFalse(t, swapped)
	value, ok := cache.Load("agent")
	AssertTrue(t, ok)
	AssertEqual(t, "starting", value)
}

func TestSync_SyncMap_CompareAndSwap_Ugly(t *T) {
	var cache SyncMap
	cache.Store("agent", nil)

	swapped := cache.CompareAndSwap("agent", nil, "ready")

	AssertTrue(t, swapped)
	value, ok := cache.Load("agent")
	AssertTrue(t, ok)
	AssertEqual(t, "ready", value)
}

func TestSync_SyncMap_CompareAndDelete_Good(t *T) {
	var cache SyncMap
	cache.Store("agent", "ready")

	deleted := cache.CompareAndDelete("agent", "ready")
	_, ok := cache.Load("agent")

	AssertTrue(t, deleted)
	AssertFalse(t, ok)
}

func TestSync_SyncMap_CompareAndDelete_Bad(t *T) {
	var cache SyncMap
	cache.Store("agent", "ready")

	deleted := cache.CompareAndDelete("agent", "stale")
	_, ok := cache.Load("agent")

	AssertFalse(t, deleted)
	AssertTrue(t, ok)
}

func TestSync_SyncMap_CompareAndDelete_Ugly(t *T) {
	var cache SyncMap

	deleted := cache.CompareAndDelete("agent", "ready")

	AssertFalse(t, deleted)
}

func TestSync_SyncMap_Range_Good(t *T) {
	var cache SyncMap
	cache.Store("agent.a", "ready")
	cache.Store("agent.b", "ready")
	count := 0

	cache.Range(func(_, _ any) bool {
		count++
		return true
	})

	AssertEqual(t, 2, count)
}

func TestSync_SyncMap_Range_Bad(t *T) {
	var cache SyncMap
	cache.Store("agent.a", "ready")
	cache.Store("agent.b", "ready")
	count := 0

	cache.Range(func(_, _ any) bool {
		count++
		return false
	})

	AssertEqual(t, 1, count)
}

func TestSync_SyncMap_Range_Ugly(t *T) {
	var cache SyncMap
	count := 0

	cache.Range(func(_, _ any) bool {
		count++
		return true
	})

	AssertEqual(t, 0, count)
}

func TestSync_SyncMap_Clear_Good(t *T) {
	var cache SyncMap
	cache.Store("agent.a", "ready")
	cache.Store("agent.b", "ready")

	cache.Clear()
	_, ok := cache.Load("agent.a")

	AssertFalse(t, ok)
}

func TestSync_SyncMap_Clear_Bad(t *T) {
	var cache SyncMap

	cache.Clear() // clearing an empty map is a safe no-op

	count := 0
	cache.Range(func(_, _ any) bool { count++; return true })
	AssertEqual(t, 0, count)
}

func TestSync_SyncMap_Clear_Ugly(t *T) {
	var cache SyncMap
	cache.Store("agent.a", "ready")

	cache.Clear()
	cache.Clear()
	_, ok := cache.Load("agent.a")

	AssertFalse(t, ok)
}
