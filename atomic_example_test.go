package core_test

import (
	. "dappco.re/go"
)

// ExampleAtomicBool updates a boolean atomically through `AtomicBool` for shared runtime
// state. Concurrent state changes use explicit load, store, swap, and compare-and-swap
// shapes.
func ExampleAtomicBool() {
	var ready AtomicBool
	if !ready.Load() {
		Println("not ready")
	}
	ready.Store(true)
	if ready.Load() {
		Println("ready")
	}
	// Output:
	// not ready
	// ready
}

// ExampleAtomicBool_Swap swaps a value through `AtomicBool.Swap` and returns the previous
// value for shared runtime state. Concurrent state changes use explicit load, store, swap,
// and compare-and-swap shapes.
func ExampleAtomicBool_Swap() {
	var ready AtomicBool
	ready.Store(true)
	Println(ready.Swap(false))
	Println(ready.Load())
	// Output:
	// true
	// false
}

// ExampleAtomicBool_CompareAndSwap updates `AtomicBool.CompareAndSwap` only when the
// previous value matches for shared runtime state. Concurrent state changes use explicit
// load, store, swap, and compare-and-swap shapes.
func ExampleAtomicBool_CompareAndSwap() {
	var ready AtomicBool
	Println(ready.CompareAndSwap(false, true))
	Println(ready.Load())
	// Output:
	// true
	// true
}

// ExampleAtomicInt32 updates a 32-bit integer atomically through `AtomicInt32` for shared
// runtime state. Concurrent state changes use explicit load, store, swap, and
// compare-and-swap shapes.
func ExampleAtomicInt32() {
	var counter AtomicInt32
	counter.Store(10)
	Println(counter.Add(5))
	Println(counter.Swap(1))
	Println(counter.Load())
	// Output:
	// 15
	// 15
	// 1
}

// ExampleAtomicInt64 updates a 64-bit integer atomically through `AtomicInt64` for shared
// runtime state. Concurrent state changes use explicit load, store, swap, and
// compare-and-swap shapes.
func ExampleAtomicInt64() {
	var counter AtomicInt64
	counter.Add(1)
	counter.Add(1)
	counter.Add(1)
	Println(counter.Load())
	// Output:
	// 3
}

// ExampleAtomicInt64_Swap swaps a value through `AtomicInt64.Swap` and returns the
// previous value for shared runtime state. Concurrent state changes use explicit load,
// store, swap, and compare-and-swap shapes.
func ExampleAtomicInt64_Swap() {
	var counter AtomicInt64
	counter.Store(7)
	Println(counter.Swap(9))
	Println(counter.Load())
	// Output:
	// 7
	// 9
}

// ExampleAtomicInt64_CompareAndSwap updates `AtomicInt64.CompareAndSwap` only when the
// previous value matches for shared runtime state. Concurrent state changes use explicit
// load, store, swap, and compare-and-swap shapes.
func ExampleAtomicInt64_CompareAndSwap() {
	var counter AtomicInt64
	counter.Store(7)
	Println(counter.CompareAndSwap(7, 9))
	Println(counter.Load())
	// Output:
	// true
	// 9
}

// ExampleAtomicInt32_CompareAndSwap updates `AtomicInt32.CompareAndSwap` only when the
// previous value matches for shared runtime state. Concurrent state changes use explicit
// load, store, swap, and compare-and-swap shapes.
func ExampleAtomicInt32_CompareAndSwap() {
	var state AtomicInt32
	if state.CompareAndSwap(0, 1) {
		Println("claimed")
	}
	if !state.CompareAndSwap(0, 2) {
		Println("already claimed")
	}
	// Output:
	// claimed
	// already claimed
}

// ExampleAtomicUint32 updates an unsigned 32-bit integer atomically through `AtomicUint32`
// for shared runtime state. Concurrent state changes use explicit load, store, swap, and
// compare-and-swap shapes.
func ExampleAtomicUint32() {
	var counter AtomicUint32
	counter.Store(10)
	Println(counter.Add(5))
	Println(counter.Swap(1))
	Println(counter.Load())
	// Output:
	// 15
	// 15
	// 1
}

// ExampleAtomicUint32_CompareAndSwap updates `AtomicUint32.CompareAndSwap` only when the
// previous value matches for shared runtime state. Concurrent state changes use explicit
// load, store, swap, and compare-and-swap shapes.
func ExampleAtomicUint32_CompareAndSwap() {
	var counter AtomicUint32
	counter.Store(10)
	Println(counter.CompareAndSwap(10, 11))
	Println(counter.Load())
	// Output:
	// true
	// 11
}

// ExampleAtomicUint64 updates an unsigned 64-bit integer atomically through `AtomicUint64`
// for shared runtime state. Concurrent state changes use explicit load, store, swap, and
// compare-and-swap shapes.
func ExampleAtomicUint64() {
	var counter AtomicUint64
	counter.Store(10)
	Println(counter.Add(5))
	Println(counter.Swap(1))
	Println(counter.Load())
	// Output:
	// 15
	// 15
	// 1
}

// ExampleAtomicUint64_CompareAndSwap updates `AtomicUint64.CompareAndSwap` only when the
// previous value matches for shared runtime state. Concurrent state changes use explicit
// load, store, swap, and compare-and-swap shapes.
func ExampleAtomicUint64_CompareAndSwap() {
	var counter AtomicUint64
	counter.Store(10)
	Println(counter.CompareAndSwap(10, 11))
	Println(counter.Load())
	// Output:
	// true
	// 11
}

type config struct {
	name string
}

// ExampleAtomicPointer updates a typed pointer atomically through `AtomicPointer` for
// shared runtime state. Concurrent state changes use explicit load, store, swap, and
// compare-and-swap shapes.
func ExampleAtomicPointer() {
	var current AtomicPointer[config]
	current.Store(&config{name: "v1"})
	cfg := current.Load()
	Println(cfg.name)
	current.Store(&config{name: "v2"})
	Println(current.Load().name)
	// Output:
	// v1
	// v2
}

// ExampleAtomicPointer_Swap swaps a value through `AtomicPointer.Swap` and returns the
// previous value for shared runtime state. Concurrent state changes use explicit load,
// store, swap, and compare-and-swap shapes.
func ExampleAtomicPointer_Swap() {
	var current AtomicPointer[config]
	first := &config{name: "v1"}
	second := &config{name: "v2"}
	current.Store(first)

	Println(current.Swap(second).name)
	Println(current.Load().name)
	// Output:
	// v1
	// v2
}

// ExampleAtomicPointer_CompareAndSwap updates `AtomicPointer.CompareAndSwap` only when the
// previous value matches for shared runtime state. Concurrent state changes use explicit
// load, store, swap, and compare-and-swap shapes.
// ExampleAtomicBool_Load reads a boolean atomically through `AtomicBool.Load`.
func ExampleAtomicBool_Load() {
	var b AtomicBool
	b.Store(true)
	Println(b.Load())
	// Output: true
}

// ExampleAtomicBool_Store writes a boolean atomically through `AtomicBool.Store`.
func ExampleAtomicBool_Store() {
	var b AtomicBool
	b.Store(true)
	Println(b.Load())
	// Output: true
}

// ExampleAtomicInt32_Load reads an int32 atomically through `AtomicInt32.Load`.
func ExampleAtomicInt32_Load() {
	var n AtomicInt32
	n.Store(5)
	Println(n.Load())
	// Output: 5
}

// ExampleAtomicInt32_Store writes an int32 atomically through `AtomicInt32.Store`.
func ExampleAtomicInt32_Store() {
	var n AtomicInt32
	n.Store(42)
	Println(n.Load())
	// Output: 42
}

// ExampleAtomicInt32_Add atomically adds and returns the new value through `AtomicInt32.Add`.
func ExampleAtomicInt32_Add() {
	var n AtomicInt32
	n.Store(5)
	Println(n.Add(3))
	// Output: 8
}

// ExampleAtomicInt32_Swap stores and returns the previous value through `AtomicInt32.Swap`.
func ExampleAtomicInt32_Swap() {
	var n AtomicInt32
	n.Store(5)
	Println(n.Swap(9))
	// Output: 5
}

// ExampleAtomicInt64_Load reads an int64 atomically through `AtomicInt64.Load`.
func ExampleAtomicInt64_Load() {
	var n AtomicInt64
	n.Store(5)
	Println(n.Load())
	// Output: 5
}

// ExampleAtomicInt64_Store writes an int64 atomically through `AtomicInt64.Store`.
func ExampleAtomicInt64_Store() {
	var n AtomicInt64
	n.Store(42)
	Println(n.Load())
	// Output: 42
}

// ExampleAtomicInt64_Add atomically adds and returns the new value through `AtomicInt64.Add`.
func ExampleAtomicInt64_Add() {
	var n AtomicInt64
	n.Store(5)
	Println(n.Add(3))
	// Output: 8
}

// ExampleAtomicUint32_Load reads a uint32 atomically through `AtomicUint32.Load`.
func ExampleAtomicUint32_Load() {
	var n AtomicUint32
	n.Store(5)
	Println(n.Load())
	// Output: 5
}

// ExampleAtomicUint32_Store writes a uint32 atomically through `AtomicUint32.Store`.
func ExampleAtomicUint32_Store() {
	var n AtomicUint32
	n.Store(42)
	Println(n.Load())
	// Output: 42
}

// ExampleAtomicUint32_Add atomically adds and returns the new value through `AtomicUint32.Add`.
func ExampleAtomicUint32_Add() {
	var n AtomicUint32
	n.Store(5)
	Println(n.Add(3))
	// Output: 8
}

// ExampleAtomicUint32_Swap stores and returns the previous value through `AtomicUint32.Swap`.
func ExampleAtomicUint32_Swap() {
	var n AtomicUint32
	n.Store(5)
	Println(n.Swap(9))
	// Output: 5
}

// ExampleAtomicUint64_Load reads a uint64 atomically through `AtomicUint64.Load`.
func ExampleAtomicUint64_Load() {
	var n AtomicUint64
	n.Store(5)
	Println(n.Load())
	// Output: 5
}

// ExampleAtomicUint64_Store writes a uint64 atomically through `AtomicUint64.Store`.
func ExampleAtomicUint64_Store() {
	var n AtomicUint64
	n.Store(42)
	Println(n.Load())
	// Output: 42
}

// ExampleAtomicUint64_Add atomically adds and returns the new value through `AtomicUint64.Add`.
func ExampleAtomicUint64_Add() {
	var n AtomicUint64
	n.Store(5)
	Println(n.Add(3))
	// Output: 8
}

// ExampleAtomicUint64_Swap stores and returns the previous value through `AtomicUint64.Swap`.
func ExampleAtomicUint64_Swap() {
	var n AtomicUint64
	n.Store(5)
	Println(n.Swap(9))
	// Output: 5
}

// ExampleAtomicPointer_Load reads a pointer atomically through `AtomicPointer.Load`.
func ExampleAtomicPointer_Load() {
	var p AtomicPointer[int]
	v := 42
	p.Store(&v)
	Println(*p.Load())
	// Output: 42
}

// ExampleAtomicPointer_Store writes a pointer atomically through `AtomicPointer.Store`.
func ExampleAtomicPointer_Store() {
	var p AtomicPointer[string]
	s := "agent"
	p.Store(&s)
	Println(*p.Load())
	// Output: agent
}

func ExampleAtomicPointer_CompareAndSwap() {
	var current AtomicPointer[config]
	first := &config{name: "v1"}
	second := &config{name: "v2"}
	current.Store(first)

	Println(current.CompareAndSwap(first, second))
	Println(current.Load().name)
	// Output:
	// true
	// v2
}
