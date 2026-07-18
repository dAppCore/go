// SPDX-License-Identifier: EUPL-1.2

// First-match value helpers for the Core framework — pick the first value from a
// list that passes a simple test, the "fallback chain" pattern (config → env →
// default) reduced to one call.

package core

// Coalesce returns the first non-zero value, or the zero value when every
// argument is zero — the fallback chain (primary, secondary, default) as one
// call. Replaces the hand-rolled firstNonEmpty/firstNonNil helpers scattered
// across consumers.
//
//	name := core.Coalesce(req.Name, cfg.Name, "anonymous")
//	timeout := core.Coalesce(flag, env, 30)
func Coalesce[T comparable](vals ...T) T {
	var zero T
	for _, v := range vals {
		if v != zero {
			return v
		}
	}
	return zero
}

// FirstNonBlank returns the first value whose trimmed form is non-empty — a
// whitespace-only string ("   ") counts as blank and falls through — returning
// the ORIGINAL untrimmed value, or "" when every value is blank. The
// whitespace-aware sibling of Coalesce for the "first meaningful string" pattern,
// where a blank field should defer to the next fallback.
//
//	title := core.FirstNonBlank(userTitle, cfg.Title, "untitled")
func FirstNonBlank(vals ...string) string {
	for _, v := range vals {
		// The v != "" fast path skips the Trim call for the common empty arg.
		if v != "" && Trim(v) != "" {
			return v
		}
	}
	return ""
}

// FirstPositive returns the first value strictly greater than the zero value, or
// the zero value when none is — the numeric fallback chain (a value of 0 means
// "unset, try the next"). For numbers this reads as "first positive"; the
// Ordered constraint also admits strings, where it degenerates to Coalesce.
//
//	patchSize := core.FirstPositive(cfg.PatchSize, defaultPatchSize)
//	softTokens := core.FirstPositive(override, computed, 256)
func FirstPositive[T Ordered](vals ...T) T {
	var zero T
	for _, v := range vals {
		if v > zero {
			return v
		}
	}
	return zero
}
