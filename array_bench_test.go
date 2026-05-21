// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Array[T comparable] primitive in array.go.
// Per AX-11 — Array sits on every typed-collection path consumers reach
// for: agent lists, vocab sets, registered protocols, route tables.
// Each public method gets a perf contract here.
//
// Run:    go test -bench='BenchmarkArray' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var benchArrayStrings = []string{
	"apple", "banana", "cherry", "date", "elderberry",
	"fig", "grape", "honeydew",
}

// --- Constructors ---

func BenchmarkNewArray_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewArray[string]()
	}
}

func BenchmarkNewArray_WithItems(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewArray(benchArrayStrings...)
	}
}

// --- Add / AddUnique ---

func BenchmarkArray_Add(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := NewArray[string]()
		a.Add(benchArrayStrings...)
	}
}

func BenchmarkArray_AddUnique_AllNew(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := NewArray[string]()
		a.AddUnique(benchArrayStrings...)
	}
}

func BenchmarkArray_AddUnique_AllDup(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := NewArray[string]("apple")
		a.AddUnique("apple", "apple", "apple", "apple")
	}
}

// --- Lookup ---

func BenchmarkArray_Contains_Hit(b *B) {
	a := NewArray(benchArrayStrings...)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.Contains("elderberry")
	}
}

func BenchmarkArray_Contains_Miss(b *B) {
	a := NewArray(benchArrayStrings...)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.Contains("nope")
	}
}

// --- Filter / Each ---

func BenchmarkArray_Filter_HalfHit(b *B) {
	a := NewArray(benchArrayStrings...)
	pred := func(s string) bool { return len(s) >= 5 }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.Filter(pred)
	}
}

func BenchmarkArray_Each(b *B) {
	a := NewArray(benchArrayStrings...)
	sink := 0
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Each(func(s string) { sink += len(s) })
	}
	_ = sink
}

// --- Remove ---

func BenchmarkArray_Remove_Found(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := NewArray(benchArrayStrings...)
		a.Remove("cherry")
	}
}

// --- Deduplicate ---

func BenchmarkArray_Deduplicate_NoDups(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := NewArray(benchArrayStrings...)
		a.Deduplicate()
	}
}

func BenchmarkArray_Deduplicate_AllDups(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := NewArray("a", "a", "a", "a", "a", "a", "a", "a")
		a.Deduplicate()
	}
}

// --- Size / Clear / AsSlice ---

func BenchmarkArray_Len(b *B) {
	a := NewArray(benchArrayStrings...)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.Len()
	}
}

func BenchmarkArray_AsSlice(b *B) {
	a := NewArray(benchArrayStrings...)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.AsSlice()
	}
}
