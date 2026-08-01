// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the map primitives in map.go.
// Per AX-11 — map helpers sit on every config-key extraction, every
// vocab pivot, every registry-driven dispatch. go-mlx hits MapKeys
// on tokeniser vocabs (~50k entries), MapFilter on feature gates,
// MapHasKey on every cache probe.
//
// Run:    go test -bench='BenchmarkMap' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	mapSinkKeys   []string
	mapSinkValues []int
	mapSinkMap    map[string]int
	mapSinkBool   bool
)

// Fixtures
var (
	mapSmall = map[string]int{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
	}
	mapMedium = func() map[string]int {
		m := make(map[string]int, 100)
		for i := range 100 {
			m[string(rune('a'+i%26))+string(rune('a'+(i/26)%26))] = i
		}
		return m
	}()
	mapLarge = func() map[string]int {
		m := make(map[string]int, 10000)
		for i := range 10000 {
			m[string(rune(i))+string(rune(i+1))] = i
		}
		return m
	}()
)

// --- MapKeys / MapValues ---

func BenchmarkMapKeys_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkKeys = MapKeys(mapSmall)
	}
}

func BenchmarkMapKeys_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkKeys = MapKeys(mapMedium)
	}
}

func BenchmarkMapKeys_Large(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkKeys = MapKeys(mapLarge)
	}
}

func BenchmarkMapValues_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkValues = MapValues(mapSmall)
	}
}

func BenchmarkMapValues_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkValues = MapValues(mapMedium)
	}
}

// --- MapClone ---

func BenchmarkMapClone_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkMap = MapClone(mapSmall)
	}
}

func BenchmarkMapClone_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkMap = MapClone(mapMedium)
	}
}

// --- MapFilter ---

func BenchmarkMapFilter_HalfHit(b *B) {
	pred := func(k string, v int) bool { return v%2 == 0 }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkMap = MapFilter(mapMedium, pred)
	}
}

// --- MapMerge ---

func BenchmarkMapMerge_Disjoint(b *B) {
	a := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	c := map[string]int{"e": 5, "f": 6, "g": 7, "h": 8}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkMap = MapMerge(a, c)
	}
}

func BenchmarkMapMerge_Overlap(b *B) {
	a := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	c := map[string]int{"a": 10, "b": 20, "e": 5, "f": 6}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkMap = MapMerge(a, c)
	}
}

// --- MapHasKey ---

func BenchmarkMapHasKey_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkBool = MapHasKey(mapMedium, "aa")
	}
}

func BenchmarkMapHasKey_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapSinkBool = MapHasKey(mapMedium, "zz_nope")
	}
}
