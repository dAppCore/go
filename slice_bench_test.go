// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the hot slice primitives in slice.go.
// Per AX-11 — slice helpers are exercised by every consumer that
// holds collections (token lists, config arrays, registry entries,
// route tables, etc.) so the contract here is load-bearing across
// the ecosystem.
//
// Run:    go test -bench='BenchmarkSlice.*' -benchmem -run='^$' .
// Filter: go test -bench='BenchmarkSliceSort' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// --- Fixture corpora ---

var (
	benchSliceSmall  = []int{1, 2, 3, 4, 5, 6, 7, 8}
	benchSliceMedium = func() []int {
		out := make([]int, 100)
		for i := range out {
			out[i] = i * 7 % 50
		}
		return out
	}()
	benchSliceStrings = []string{
		"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew",
	}
)

// --- Contains / Index ---

func BenchmarkSliceContains_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceContains(benchSliceMedium, 35)
	}
}

func BenchmarkSliceContains_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceContains(benchSliceMedium, 999)
	}
}

func BenchmarkSliceIndex_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceIndex(benchSliceMedium, 35)
	}
}

// --- Clone / Map / Filter / Reduce ---

func BenchmarkSliceClone_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceClone(benchSliceSmall)
	}
}

func BenchmarkSliceClone_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceClone(benchSliceMedium)
	}
}

func BenchmarkSliceMap_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceMap(benchSliceMedium, func(n int) int { return n * 2 })
	}
}

func BenchmarkSliceFilter_HalfHit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceFilter(benchSliceMedium, func(n int) bool { return n%2 == 0 })
	}
}

func BenchmarkSliceReduce_Sum(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceReduce(benchSliceMedium, 0, func(acc, n int) int { return acc + n })
	}
}

// --- Sort / SortFunc ---

func BenchmarkSliceSort_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := SliceClone(benchSliceMedium)
		SliceSort(s)
	}
}

func BenchmarkSliceSortFunc_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := SliceClone(benchSliceMedium)
		SliceSortFunc(s, func(a, b int) bool { return a < b })
	}
}

// --- Uniq ---

func BenchmarkSliceUniq_NoDups(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceUniq(benchSliceStrings)
	}
}

func BenchmarkSliceUniq_AllDups(b *B) {
	all := []string{"a", "a", "a", "a", "a", "a", "a", "a"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceUniq(all)
	}
}

// --- Predicates ---

func BenchmarkSliceAny_EarlyHit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceAny(benchSliceMedium, func(n int) bool { return n > 5 })
	}
}

func BenchmarkSliceAll_Pass(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceAll(benchSliceMedium, func(n int) bool { return n >= 0 })
	}
}

// --- Take / Drop ---

func BenchmarkSliceTake_Half(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceTake(benchSliceMedium, 50)
	}
}

func BenchmarkSliceDrop_Half(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceDrop(benchSliceMedium, 50)
	}
}

// --- FlatMap ---

func BenchmarkSliceFlatMap(b *B) {
	src := []int{1, 2, 3, 4, 5}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceFlatMap(src, func(n int) []int { return []int{n, n * 2, n * 3} })
	}
}

// --- Reverse / Sorted ---

func BenchmarkSliceReverse_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := SliceClone(benchSliceMedium)
		SliceReverse(s)
	}
}

func BenchmarkSliceSorted_Iter(b *B) {
	src := benchSliceMedium
	seq := Seq[int](func(yield func(int) bool) {
		for _, v := range src {
			if !yield(v) {
				return
			}
		}
	})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SliceSorted(seq)
	}
}
