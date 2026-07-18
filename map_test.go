package core_test

import (
	. "dappco.re/go"
)

func TestMap_MapKeys_Good(t *T) {
	keys := MapKeys(map[string]int{"a": 1, "b": 2})

	// Order is not guaranteed, so match by membership and assert length.
	AssertElementsMatch(t, []string{"a", "b"}, keys)
	AssertLen(t, keys, 2)
	AssertContains(t, keys, "a")
	AssertContains(t, keys, "b")

	// A single-entry map yields exactly that one key.
	single := MapKeys(map[string]int{"only": 42})
	AssertEqual(t, []string{"only"}, single)
}

func TestMap_MapKeys_Bad(t *T) {
	keys := MapKeys[string, int](nil)

	// make([]K, len(nil)) returns a non-nil, zero-length slice.
	AssertEmpty(t, keys)
	AssertNotNil(t, keys)
	AssertLen(t, keys, 0)

	// An empty (non-nil) map also produces zero keys.
	AssertEmpty(t, MapKeys(map[string]int{}))
}

func TestMap_MapKeys_Ugly(t *T) {
	keys := MapKeys(map[int]string{2: "b", 1: "a", 3: "c"})

	SliceSort(keys)
	AssertEqual(t, []int{1, 2, 3}, keys)
}

func TestMap_MapValues_Good(t *T) {
	values := MapValues(map[string]int{"a": 1, "b": 2})

	AssertElementsMatch(t, []int{1, 2}, values)
	AssertLen(t, values, 2)
	AssertContains(t, values, 1)
	AssertContains(t, values, 2)

	// One value is emitted per key, so duplicates are preserved.
	dups := MapValues(map[string]int{"a": 7, "b": 7})
	AssertElementsMatch(t, []int{7, 7}, dups)
	AssertLen(t, dups, 2)
}

func TestMap_MapValues_Bad(t *T) {
	values := MapValues[string, int](nil)

	// make([]V, len(nil)) returns a non-nil, zero-length slice.
	AssertEmpty(t, values)
	AssertNotNil(t, values)
	AssertLen(t, values, 0)

	// An empty (non-nil) map also produces zero values.
	AssertEmpty(t, MapValues(map[string]int{}))
}

func TestMap_MapValues_Ugly(t *T) {
	values := MapValues(map[int]string{2: "b", 1: "a", 3: "c"})

	SliceSort(values)
	AssertEqual(t, []string{"a", "b", "c"}, values)
}

func TestMap_MapClone_Good(t *T) {
	clone := MapClone(map[string]int{"a": 1})

	AssertEqual(t, map[string]int{"a": 1}, clone)
	AssertLen(t, clone, 1)

	// Every entry of a multi-key map is copied across.
	multi := MapClone(map[string]int{"a": 1, "b": 2, "c": 3})
	AssertEqual(t, map[string]int{"a": 1, "b": 2, "c": 3}, multi)
	AssertLen(t, multi, 3)
	AssertTrue(t, MapHasKey(multi, "b"))
}

func TestMap_MapClone_Bad(t *T) {
	clone := MapClone[string, int](nil)

	// maps.Clone preserves nil: a nil map clones to a nil map.
	AssertNil(t, clone)
	AssertEmpty(t, clone)
	AssertLen(t, clone, 0)
}

func TestMap_MapClone_Ugly(t *T) {
	original := map[string]int{"a": 1}
	clone := MapClone(original)
	clone["a"] = 2

	AssertEqual(t, 1, original["a"])
	AssertEqual(t, 2, clone["a"])
}

func TestMap_MapFilter_Good(t *T) {
	filtered := MapFilter(map[string]bool{"agent": true, "archive": false}, func(_ string, enabled bool) bool {
		return enabled
	})

	AssertEqual(t, map[string]bool{"agent": true}, filtered)
}

func TestMap_MapFilter_Bad(t *T) {
	filtered := MapFilter(map[string]bool{"agent": false}, func(_ string, enabled bool) bool {
		return enabled
	})

	AssertEmpty(t, filtered)
}

func TestMap_MapFilter_Ugly(t *T) {
	// Nil input short-circuits (len==0) to a nil result, regardless of pred.
	AssertNil(t, MapFilter[string, bool](nil, func(string, bool) bool { return true }))

	// An empty (non-nil) map also short-circuits to nil.
	AssertNil(t, MapFilter(map[string]bool{}, func(string, bool) bool { return true }))

	// A pred that rejects everything yields an empty (nil) result too.
	AssertEmpty(t, MapFilter(map[string]int{"a": 1, "b": 2}, func(string, int) bool { return false }))
}

func TestMap_MapHasKey_Good(t *T) {
	m := map[string]int{"agent": 1, "zero": 0}

	AssertTrue(t, MapHasKey(m, "agent"))
	// A key mapped to the zero value is still reported as present.
	AssertTrue(t, MapHasKey(m, "zero"))
}

func TestMap_MapHasKey_Bad(t *T) {
	AssertFalse(t, MapHasKey(map[string]int{"agent": 1}, "missing"))
	// An empty (non-nil) map contains no keys.
	AssertFalse(t, MapHasKey(map[string]int{}, "agent"))
}

func TestMap_MapHasKey_Ugly(t *T) {
	// Lookup on a nil map is safe and reports absence.
	AssertFalse(t, MapHasKey(map[string]int(nil), "agent"))
	// Same for a nil map with a different key/value type — no panic.
	AssertFalse(t, MapHasKey(map[int]string(nil), 1))
}

func TestMap_MapMerge_Good(t *T) {
	a := map[string]string{"agent": "codex"}
	b := map[string]string{"region": "homelab"}
	merged := MapMerge(a, b)

	AssertEqual(t, map[string]string{"agent": "codex", "region": "homelab"}, merged)
	AssertLen(t, merged, 2)
	AssertEqual(t, "codex", merged["agent"])
	AssertEqual(t, "homelab", merged["region"])

	// Merge builds a fresh map and leaves both inputs untouched.
	AssertLen(t, a, 1)
	AssertLen(t, b, 1)
}

func TestMap_MapMerge_Bad(t *T) {
	// On a key collision, b's value wins.
	merged := MapMerge(map[string]string{"agent": "codex"}, map[string]string{"agent": "hades"})

	AssertEqual(t, map[string]string{"agent": "hades"}, merged)
	AssertLen(t, merged, 1)

	// Non-colliding keys from a survive; only the shared key is overridden.
	mixed := MapMerge(
		map[string]string{"agent": "codex", "region": "homelab"},
		map[string]string{"agent": "hades"},
	)
	AssertEqual(t, "hades", mixed["agent"])
	AssertEqual(t, "homelab", mixed["region"])
	AssertLen(t, mixed, 2)
}

func TestMap_MapMerge_Ugly(t *T) {
	merged := MapMerge[string, string](nil, nil)

	AssertEmpty(t, merged)
}
