package core_test

import (
	. "dappco.re/go"
)

// ExampleMutex uses a mutual exclusion lock through `Mutex` for concurrent service
// coordination. Concurrency helpers mirror the stdlib shapes while keeping ownership in
// core.
func ExampleMutex() {
	var mu Mutex
	mu.Lock()
	Println("locked")
	mu.Unlock()
	Println("unlocked")
	// Output:
	// locked
	// unlocked
}

// ExampleMutex_TryLock attempts a non-blocking write lock through `Mutex.TryLock` for
// concurrent service coordination. Concurrency helpers mirror the stdlib shapes while
// keeping ownership in core.
func ExampleMutex_TryLock() {
	var mu Mutex
	if mu.TryLock().OK {
		Println("acquired")
		mu.Unlock()
	}
	// Output: acquired
}

// ExampleRWMutex uses a read-write lock through `RWMutex` for concurrent service
// coordination. Concurrency helpers mirror the stdlib shapes while keeping ownership in
// core.
func ExampleRWMutex() {
	var mu RWMutex
	mu.RLock()
	Println("read-locked")
	mu.RUnlock()
	mu.Lock()
	Println("write-locked")
	mu.Unlock()
	// Output:
	// read-locked
	// write-locked
}

// ExampleRWMutex_TryLock attempts a non-blocking write lock through `RWMutex.TryLock` for
// concurrent service coordination. Concurrency helpers mirror the stdlib shapes while
// keeping ownership in core.
func ExampleRWMutex_TryLock() {
	var mu RWMutex
	if mu.TryLock().OK {
		Println("write")
		mu.Unlock()
	}
	// Output: write
}

// ExampleRWMutex_TryRLock attempts a non-blocking read lock through `RWMutex.TryRLock` for
// concurrent service coordination. Concurrency helpers mirror the stdlib shapes while
// keeping ownership in core.
func ExampleRWMutex_TryRLock() {
	var mu RWMutex
	if mu.TryRLock().OK {
		Println("read")
		mu.RUnlock()
	}
	// Output: read
}

// ExampleOnce runs work once through `Once` for concurrent service coordination.
// Concurrency helpers mirror the stdlib shapes while keeping ownership in core.
func ExampleOnce() {
	var o Once
	o.Do(func() { Println("ran") })
	o.Do(func() { Println("not printed") })
	// Output:
	// ran
}

// ExampleOnce_Reset resets one-shot state through `Once.Reset` for concurrent service
// coordination. Concurrency helpers mirror the stdlib shapes while keeping ownership in
// core.
func ExampleOnce_Reset() {
	var o Once
	o.Do(func() { Println("first") })
	o.Reset()
	o.Do(func() { Println("again") })
	// Output:
	// first
	// again
}

// ExampleWaitGroup waits for concurrent work through `WaitGroup` for concurrent service
// coordination. Concurrency helpers mirror the stdlib shapes while keeping ownership in
// core.
// ExampleMutex_Lock acquires exclusive access through `Mutex.Lock`.
func ExampleMutex_Lock() {
	var mu Mutex
	mu.Lock()
	defer mu.Unlock()
	// critical section
}

// ExampleMutex_Unlock releases exclusive access through `Mutex.Unlock`.
func ExampleMutex_Unlock() {
	var mu Mutex
	mu.Lock()
	mu.Unlock()
}

// ExampleRWMutex_Lock acquires the write lock through `RWMutex.Lock`.
func ExampleRWMutex_Lock() {
	var mu RWMutex
	mu.Lock()
	defer mu.Unlock()
	// exclusive write section
}

// ExampleRWMutex_Unlock releases the write lock through `RWMutex.Unlock`.
func ExampleRWMutex_Unlock() {
	var mu RWMutex
	mu.Lock()
	mu.Unlock()
}

// ExampleRWMutex_RLock acquires a shared read lock through `RWMutex.RLock`.
func ExampleRWMutex_RLock() {
	var mu RWMutex
	mu.RLock()
	defer mu.RUnlock()
	// concurrent read section
}

// ExampleRWMutex_RUnlock releases a shared read lock through `RWMutex.RUnlock`.
func ExampleRWMutex_RUnlock() {
	var mu RWMutex
	mu.RLock()
	mu.RUnlock()
}

// ExampleOnce_Do runs an initialiser exactly once through `Once.Do`.
func ExampleOnce_Do() {
	var once Once
	for range 3 {
		once.Do(func() { Println("init") })
	}
	// Output: init
}

// ExampleWaitGroup_Add registers pending work through `WaitGroup.Add`.
func ExampleWaitGroup_Add() {
	var wg WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Wait()
}

// ExampleWaitGroup_Done marks one unit of work complete through `WaitGroup.Done`.
func ExampleWaitGroup_Done() {
	var wg WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Wait()
}

// ExampleWaitGroup_Wait blocks until the counter reaches zero through `WaitGroup.Wait`.
func ExampleWaitGroup_Wait() {
	var wg WaitGroup
	wg.Wait() // returns immediately on a zero counter
}

// ExampleWaitGroup_Go runs a function in a tracked goroutine through `WaitGroup.Go`.
func ExampleWaitGroup_Go() {
	var wg WaitGroup
	wg.Go(func() { Println("done") })
	wg.Wait()
	// Output: done
}

// ExampleSyncMap_Store writes a key through `SyncMap.Store`.
func ExampleSyncMap_Store() {
	var m SyncMap
	m.Store("agents", 3)
	v, _ := m.Load("agents")
	Println(v)
	// Output: 3
}

// ExampleSyncMap_Load reads a key through `SyncMap.Load`.
func ExampleSyncMap_Load() {
	var m SyncMap
	m.Store("key", "value")
	v, ok := m.Load("key")
	Println(v, ok)
	// Output: value true
}

// ExampleSyncMap_LoadOrStore returns the existing value or stores a new one through `SyncMap.LoadOrStore`.
func ExampleSyncMap_LoadOrStore() {
	var m SyncMap
	m.Store("key", 1)
	actual, loaded := m.LoadOrStore("key", 2)
	Println(actual, loaded)
	// Output: 1 true
}

// ExampleSyncMap_LoadAndDelete reads then removes a key through `SyncMap.LoadAndDelete`.
func ExampleSyncMap_LoadAndDelete() {
	var m SyncMap
	m.Store("key", 1)
	v, loaded := m.LoadAndDelete("key")
	Println(v, loaded)
	// Output: 1 true
}

// ExampleSyncMap_Delete removes a key through `SyncMap.Delete`.
func ExampleSyncMap_Delete() {
	var m SyncMap
	m.Store("key", 1)
	m.Delete("key")
	_, ok := m.Load("key")
	Println(ok)
	// Output: false
}

// ExampleSyncMap_Swap replaces a value and returns the previous one through `SyncMap.Swap`.
func ExampleSyncMap_Swap() {
	var m SyncMap
	m.Store("key", 1)
	prev, loaded := m.Swap("key", 2)
	Println(prev, loaded)
	// Output: 1 true
}

// ExampleSyncMap_CompareAndSwap swaps only if the current value matches through `SyncMap.CompareAndSwap`.
func ExampleSyncMap_CompareAndSwap() {
	var m SyncMap
	m.Store("key", 1)
	Println(m.CompareAndSwap("key", 1, 2))
	// Output: true
}

// ExampleSyncMap_CompareAndDelete deletes only if the current value matches through `SyncMap.CompareAndDelete`.
func ExampleSyncMap_CompareAndDelete() {
	var m SyncMap
	m.Store("key", 1)
	Println(m.CompareAndDelete("key", 1))
	// Output: true
}

// ExampleSyncMap_Range iterates over all entries through `SyncMap.Range`.
func ExampleSyncMap_Range() {
	var m SyncMap
	m.Store("only", 1)
	count := 0
	m.Range(func(_, _ any) bool {
		count++
		return true
	})
	Println(count)
	// Output: 1
}

// ExampleSyncMap_Clear removes all entries through `SyncMap.Clear`.
func ExampleSyncMap_Clear() {
	var m SyncMap
	m.Store("key", 1)
	m.Clear()
	_, ok := m.Load("key")
	Println(ok)
	// Output: false
}

func ExampleWaitGroup() {
	var wg WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Println("worker done")
	}()
	wg.Wait()
	// Output:
	// worker done
}
