package core_test

import (
	. "dappco.re/go"
)

func TestSlice_SliceContains_Good(t *T) {
	AssertTrue(t, SliceContains([]string{"a", "b"}, "b"))
	// First element.
	AssertTrue(t, SliceContains([]string{"a", "b", "c"}, "a"))
	// Last element.
	AssertTrue(t, SliceContains([]string{"a", "b", "c"}, "c"))
	// Single-element slice.
	AssertTrue(t, SliceContains([]string{"only"}, "only"))
	// Duplicates still report present.
	AssertTrue(t, SliceContains([]int{7, 7, 7}, 7))
}

func TestSlice_SliceContains_Bad(t *T) {
	AssertFalse(t, SliceContains([]string{"a", "b"}, "c"))
	// Empty slice contains nothing.
	AssertFalse(t, SliceContains([]string{}, "a"))
	// Match is case-sensitive: "A" is not "a".
	AssertFalse(t, SliceContains([]string{"a", "b"}, "A"))
	// Absent numeric value.
	AssertFalse(t, SliceContains([]int{1, 2, 3}, 4))
}

func TestSlice_SliceContains_Ugly(t *T) {
	// Nil slice never contains the zero value.
	AssertFalse(t, SliceContains([]int(nil), 0))
	AssertFalse(t, SliceContains([]string(nil), ""))
	// But a real zero value present in a non-nil slice is found.
	AssertTrue(t, SliceContains([]int{0, 1}, 0))
	AssertTrue(t, SliceContains([]string{""}, ""))
}

func TestSlice_SliceEqual_Good(t *T) {
	AssertTrue(t, SliceEqual([]int{1, 2, 3}, []int{1, 2, 3}))
}

func TestSlice_SliceEqual_Bad(t *T) {
	AssertFalse(t, SliceEqual([]int{1, 2, 3}, []int{1, 2, 4}))
}

func TestSlice_SliceEqual_Ugly(t *T) {
	// different lengths are unequal; nil and empty are equal
	AssertFalse(t, SliceEqual([]byte("ab"), []byte("abc")))
	AssertTrue(t, SliceEqual([]int(nil), []int{}))
}

func TestSlice_SliceIndex_Good(t *T) {
	AssertEqual(t, 1, SliceIndex([]string{"a", "b"}, "b"))
	// First position.
	AssertEqual(t, 0, SliceIndex([]string{"a", "b", "c"}, "a"))
	// Last position.
	AssertEqual(t, 2, SliceIndex([]string{"a", "b", "c"}, "c"))
	// Single-element slice.
	AssertEqual(t, 0, SliceIndex([]string{"only"}, "only"))
}

func TestSlice_SliceIndex_Bad(t *T) {
	AssertEqual(t, -1, SliceIndex([]string{"a", "b"}, "c"))
	// Empty and nil slices yield -1.
	AssertEqual(t, -1, SliceIndex([]string{}, "a"))
	AssertEqual(t, -1, SliceIndex([]int(nil), 0))
	// Absent numeric value.
	AssertEqual(t, -1, SliceIndex([]int{1, 2, 3}, 9))
}

func TestSlice_SliceIndex_Ugly(t *T) {
	// Duplicates: the first matching index wins.
	AssertEqual(t, 0, SliceIndex([]int{7, 7, 7}, 7))
	AssertEqual(t, 1, SliceIndex([]int{3, 7, 7, 7}, 7))
	// A zero value is located at its real position.
	AssertEqual(t, 2, SliceIndex([]int{1, 2, 0, 0}, 0))
}

func TestSlice_SliceSort_Good(t *T) {
	items := []int{3, 1, 2}
	SliceSort(items)

	AssertEqual(t, []int{1, 2, 3}, items)
}

func TestSlice_SliceSort_Bad(t *T) {
	var items []int
	SliceSort(items)

	AssertNil(t, items)
}

func TestSlice_SliceSort_Ugly(t *T) {
	items := []string{"beta", "alpha", "beta"}
	SliceSort(items)

	AssertEqual(t, []string{"alpha", "beta", "beta"}, items)
}

func TestSlice_SliceUniq_Good(t *T) {
	items := SliceUniq([]string{"a", "b", "a"})

	AssertEqual(t, []string{"a", "b"}, items)
	// Order of first appearance is preserved.
	AssertEqual(t, []string{"c", "a", "b"}, SliceUniq([]string{"c", "a", "c", "b", "a"}))
	// A slice with no duplicates keeps every element.
	AssertEqual(t, []string{"x", "y", "z"}, SliceUniq([]string{"x", "y", "z"}))
	// Single element.
	AssertEqual(t, []int{5}, SliceUniq([]int{5}))
}

func TestSlice_SliceUniq_Bad(t *T) {
	// Nil and empty inputs both collapse to nil.
	AssertNil(t, SliceUniq([]string(nil)))
	AssertNil(t, SliceUniq([]int{}))
	// The nil result reports length zero.
	AssertLen(t, SliceUniq([]int{}), 0)
	// A non-empty input is never nil.
	AssertNotNil(t, SliceUniq([]int{1}))
}

func TestSlice_SliceUniq_Ugly(t *T) {
	items := []int{3, 3, 2, 1, 2, 1}

	AssertEqual(t, []int{3, 2, 1}, SliceUniq(items))
	// The original input is not mutated.
	AssertEqual(t, []int{3, 3, 2, 1, 2, 1}, items)
	// All-identical collapses to a single element.
	AssertEqual(t, []int{9}, SliceUniq([]int{9, 9, 9, 9}))
	// Inputs longer than 16 take the map-based dedupe path while still
	// preserving first-appearance order.
	large := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 0, 1, 16}
	AssertEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, SliceUniq(large))
}

func TestSlice_SliceReverse_Good(t *T) {
	items := []int{1, 2, 3}
	SliceReverse(items)

	AssertEqual(t, []int{3, 2, 1}, items)
}

func TestSlice_SliceReverse_Bad(t *T) {
	var items []int
	SliceReverse(items)

	AssertNil(t, items)
}

func TestSlice_SliceReverse_Ugly(t *T) {
	items := []string{"a", "b", "c", "d"}
	SliceReverse(items)

	AssertEqual(t, []string{"d", "c", "b", "a"}, items)
}

func TestSlice_SliceAll_Good(t *T) {
	AssertTrue(t, SliceAll([]string{"agent", "agent.dispatch"}, func(name string) bool {
		return HasPrefix(name, "agent")
	}))
}

func TestSlice_SliceAll_Bad(t *T) {
	AssertFalse(t, SliceAll([]int{2, 4, 5}, func(n int) bool { return n%2 == 0 }))
	// A failure on the very first element short-circuits to false.
	AssertFalse(t, SliceAll([]int{1, 2, 4}, func(n int) bool { return n%2 == 0 }))
	// Single element that fails the predicate.
	AssertFalse(t, SliceAll([]int{3}, func(n int) bool { return n%2 == 0 }))
}

func TestSlice_SliceAll_Ugly(t *T) {
	// Vacuously true for an empty slice: the predicate is never called.
	AssertTrue(t, SliceAll([]string{}, func(string) bool { return false }))
	// A nil slice is also vacuously true.
	AssertTrue(t, SliceAll([]int(nil), func(int) bool { return false }))
	// Every element satisfies the predicate.
	AssertTrue(t, SliceAll([]int{2, 4, 6}, func(n int) bool { return n%2 == 0 }))
}

func TestSlice_SliceAny_Good(t *T) {
	AssertTrue(t, SliceAny([]string{"codex", "hades"}, func(name string) bool {
		return name == "hades"
	}))
}

func TestSlice_SliceAny_Bad(t *T) {
	AssertFalse(t, SliceAny([]string{"codex", "hades"}, func(name string) bool {
		return name == "homelab"
	}))
}

func TestSlice_SliceAny_Ugly(t *T) {
	// Empty slice: the predicate is never invoked, so a nil pred is safe.
	AssertFalse(t, SliceAny([]string{}, nil))
	// A nil slice behaves the same.
	AssertFalse(t, SliceAny([]int(nil), nil))
	// Non-empty input where nothing matches.
	AssertFalse(t, SliceAny([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 }))
}

func TestSlice_SliceClone_Good(t *T) {
	clone := SliceClone([]string{"codex", "hades"})

	AssertEqual(t, []string{"codex", "hades"}, clone)
	// Length matches the source.
	AssertLen(t, clone, 2)
	// Cloning an empty (non-nil) slice yields an empty, non-nil slice.
	empty := SliceClone([]int{})
	AssertNotNil(t, empty)
	AssertLen(t, empty, 0)
	// Numeric slices clone element-for-element.
	AssertEqual(t, []int{1, 2, 3}, SliceClone([]int{1, 2, 3}))
}

func TestSlice_SliceClone_Bad(t *T) {
	// Nil in, nil out — nil-ness is preserved.
	AssertNil(t, SliceClone[string](nil))
	AssertNil(t, SliceClone[int](nil))
	// The nil clone reports length zero.
	AssertLen(t, SliceClone[int](nil), 0)
	// A non-nil source clones to a non-nil result.
	AssertNotNil(t, SliceClone([]int{1}))
}

func TestSlice_SliceClone_Ugly(t *T) {
	original := []string{"codex"}
	clone := SliceClone(original)
	clone[0] = "hades"

	AssertEqual(t, []string{"codex"}, original)
	AssertEqual(t, []string{"hades"}, clone)
}

func TestSlice_SliceDrop_Good(t *T) {
	AssertEqual(t, []string{"dispatch", "ready"}, SliceDrop([]string{"agent", "dispatch", "ready"}, 1))
	// Dropping two leaves only the tail.
	AssertEqual(t, []string{"ready"}, SliceDrop([]string{"agent", "dispatch", "ready"}, 2))
	// Numeric slice, drop the head.
	AssertEqual(t, []int{2, 3, 4}, SliceDrop([]int{1, 2, 3, 4}, 1))
}

func TestSlice_SliceDrop_Bad(t *T) {
	// n greater than length drops everything.
	AssertNil(t, SliceDrop([]string{"agent"}, 2))
	// n exactly equal to length also yields nil.
	AssertNil(t, SliceDrop([]int{1, 2, 3}, 3))
	// Dropping from an empty slice yields nil.
	AssertNil(t, SliceDrop([]int{}, 1))
}

func TestSlice_SliceDrop_Ugly(t *T) {
	items := []string{"agent", "dispatch"}

	// n == 0 is a no-op: the whole slice comes back.
	AssertEqual(t, items, SliceDrop(items, 0))
	// Negative n is also a no-op.
	AssertEqual(t, items, SliceDrop(items, -3))
	// The source is never mutated by Drop.
	AssertEqual(t, []string{"agent", "dispatch"}, items)
}

func TestSlice_SliceFilter_Good(t *T) {
	got := SliceFilter([]string{"codex", "hades", "homelab"}, func(name string) bool {
		return HasPrefix(name, "h")
	})

	AssertEqual(t, []string{"hades", "homelab"}, got)
}

func TestSlice_SliceFilter_Bad(t *T) {
	got := SliceFilter([]string{"codex", "hades"}, func(name string) bool {
		return name == "missing"
	})

	AssertEmpty(t, got)
}

func TestSlice_SliceFilter_Ugly(t *T) {
	// Nil input yields nil regardless of the predicate.
	AssertNil(t, SliceFilter([]string(nil), func(string) bool { return true }))
	// Empty (non-nil) input also yields nil.
	AssertNil(t, SliceFilter([]int{}, func(int) bool { return true }))
	// A predicate that keeps everything returns all elements in order.
	AssertEqual(t, []int{1, 2, 3}, SliceFilter([]int{1, 2, 3}, func(int) bool { return true }))
}

func TestSlice_SliceFlatMap_Good(t *T) {
	got := SliceFlatMap([]string{"agent dispatch", "homelab ready"}, func(line string) []string {
		return Split(line, " ")
	})

	AssertEqual(t, []string{"agent", "dispatch", "homelab", "ready"}, got)
}

func TestSlice_SliceFlatMap_Bad(t *T) {
	// Empty input maps to nil.
	AssertNil(t, SliceFlatMap([]string{}, func(line string) []string { return []string{line} }))
	// Nil input maps to nil.
	AssertNil(t, SliceFlatMap([]string(nil), func(line string) []string { return []string{line} }))
	// Non-empty input whose fn yields nothing for every element also
	// collapses to nil.
	AssertNil(t, SliceFlatMap([]string{"a", "b"}, func(string) []string { return nil }))
}

func TestSlice_SliceFlatMap_Ugly(t *T) {
	got := SliceFlatMap([]string{"agent", "", "ready"}, func(line string) []string {
		if line == "" {
			return nil
		}
		return []string{line}
	})

	AssertEqual(t, []string{"agent", "ready"}, got)
}

func TestSlice_SliceMap_Good(t *T) {
	got := SliceMap([]string{"codex", "hades"}, func(name string) string { return Upper(name) })

	AssertEqual(t, []string{"CODEX", "HADES"}, got)
	// Output length always matches input length.
	AssertLen(t, got, 2)
	// Single element.
	AssertEqual(t, []string{"X"}, SliceMap([]string{"x"}, func(s string) string { return Upper(s) }))
	// Type-changing map: string -> rune count.
	AssertEqual(t, []int{5, 5}, SliceMap([]string{"codex", "hades"}, func(s string) int { return RuneCount(s) }))
}

func TestSlice_SliceMap_Bad(t *T) {
	// Empty input maps to nil (no allocation).
	AssertNil(t, SliceMap([]string{}, func(name string) string { return Upper(name) }))
	// Nil input maps to nil.
	AssertNil(t, SliceMap([]int(nil), func(n int) int { return n * 2 }))
	// The nil result reports length zero.
	AssertLen(t, SliceMap([]int(nil), func(n int) int { return n * 2 }), 0)
	// A non-empty input maps to a non-nil result.
	AssertNotNil(t, SliceMap([]int{1}, func(n int) int { return n * 2 }))
}

func TestSlice_SliceMap_Ugly(t *T) {
	got := SliceMap([]string{"", "agent"}, func(name string) int { return RuneCount(name) })

	AssertEqual(t, []int{0, 5}, got)
	// Positions are preserved one-to-one, including the zero-length entry.
	AssertLen(t, got, 2)
	// Map into a different type (predicate result) keeps order.
	AssertEqual(t, []bool{true, false, true}, SliceMap([]int{2, 3, 4}, func(n int) bool { return n%2 == 0 }))
}

func TestSlice_SliceReduce_Good(t *T) {
	got := SliceReduce([]int{1, 2, 3}, 0, func(total, n int) int { return total + n })

	AssertEqual(t, 6, got)
	// The initial accumulator is included in the fold.
	AssertEqual(t, 16, SliceReduce([]int{1, 2, 3}, 10, func(total, n int) int { return total + n }))
	// Product fold starting from 1.
	AssertEqual(t, 24, SliceReduce([]int{1, 2, 3, 4}, 1, func(total, n int) int { return total * n }))
	// Single element folds to init + element.
	AssertEqual(t, 5, SliceReduce([]int{5}, 0, func(total, n int) int { return total + n }))
}

func TestSlice_SliceReduce_Bad(t *T) {
	got := SliceReduce([]int{}, 42, func(total, n int) int { return total + n })

	AssertEqual(t, 42, got)
	// A nil slice also returns the initial accumulator untouched.
	AssertEqual(t, 7, SliceReduce([]int(nil), 7, func(total, n int) int { return total + n }))
	// Empty input never invokes fn, so the init string is returned as-is.
	AssertEqual(t, "seed", SliceReduce([]string{}, "seed", func(total, s string) string { return total + s }))
}

func TestSlice_SliceReduce_Ugly(t *T) {
	got := SliceReduce([]string{"agent", "dispatch"}, "", func(total, part string) string {
		if total == "" {
			return part
		}
		return Join(".", total, part)
	})

	AssertEqual(t, "agent.dispatch", got)
}

func TestSlice_SliceSorted_Good(t *T) {
	seq := func(yield func(string) bool) {
		yield("hades")
		yield("codex")
		yield("homelab")
	}

	AssertEqual(t, []string{"codex", "hades", "homelab"}, SliceSorted(seq))
}

func TestSlice_SliceSorted_Bad(t *T) {
	seq := func(yield func(int) bool) { /* empty sequence yields nothing */ }

	AssertEmpty(t, SliceSorted(seq))
	// An empty sequence collects to nil.
	AssertNil(t, SliceSorted(seq))
	// A single-element sequence collects to that one element.
	one := func(yield func(int) bool) { yield(42) }
	AssertEqual(t, []int{42}, SliceSorted(one))
}

func TestSlice_SliceSorted_Ugly(t *T) {
	seq := func(yield func(int) bool) {
		yield(3)
		yield(1)
		yield(3)
	}

	AssertEqual(t, []int{1, 3, 3}, SliceSorted(seq))
}

func TestSlice_SliceTake_Good(t *T) {
	AssertEqual(t, []string{"agent", "dispatch"}, SliceTake([]string{"agent", "dispatch", "ready"}, 2))
	// Take one yields just the head.
	AssertEqual(t, []string{"agent"}, SliceTake([]string{"agent", "dispatch", "ready"}, 1))
	// Take exactly the length returns the whole slice.
	AssertEqual(t, []int{1, 2, 3}, SliceTake([]int{1, 2, 3}, 3))
}

func TestSlice_SliceTake_Bad(t *T) {
	// n == 0 yields nil.
	AssertNil(t, SliceTake([]string{"agent"}, 0))
	// Negative n also yields nil.
	AssertNil(t, SliceTake([]int{1, 2, 3}, -2))
	// n == 0 on a multi-element slice still yields nil.
	AssertNil(t, SliceTake([]string{"a", "b", "c"}, 0))
}

func TestSlice_SliceSortFunc_Good(t *T) {
	// A custom comparator (descending) orders the slice in place.
	s := []int{3, 1, 4, 1, 5}
	SliceSortFunc(s, func(a, b int) bool { return a > b })
	AssertEqual(t, []int{5, 4, 3, 1, 1}, s)
}

func TestSlice_SliceSortFunc_Bad(t *T) {
	// Empty and single-element slices are no-ops, not panics.
	var empty []int
	SliceSortFunc(empty, func(a, b int) bool { return a < b })
	AssertEqual(t, 0, len(empty))

	one := []int{42}
	SliceSortFunc(one, func(a, b int) bool { return a < b })
	AssertEqual(t, []int{42}, one)
}

func TestSlice_SliceSortFunc_Ugly(t *T) {
	// A comparator on a derived key (string length) gives a total order.
	s := []string{"ccc", "a", "bb"}
	SliceSortFunc(s, func(a, b string) bool { return len(a) < len(b) })
	AssertEqual(t, []string{"a", "bb", "ccc"}, s)
}

func TestSlice_SliceTake_Ugly(t *T) {
	items := []string{"agent", "dispatch"}

	// n larger than length returns the entire slice.
	AssertEqual(t, items, SliceTake(items, 5))
	// n exactly equal to length also returns the whole slice.
	AssertEqual(t, items, SliceTake(items, 2))
	// The source is never mutated.
	AssertEqual(t, []string{"agent", "dispatch"}, items)
}
