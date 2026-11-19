// SPDX-License-Identifier: EUPL-1.2

package core

import "unsafe"

// SliceDataPtr is &s[0] as an unsafe.Pointer, so unsafe_test.go can check
// that a PinnedView points at the caller's own backing array (zero-copy)
// without importing unsafe itself. Exported test vars and funcs in package
// core are visible to core_test.
//
//	AssertTrue(t, view.Ptr() == SliceDataPtr(slice))
func SliceDataPtr[E any](s []E) unsafe.Pointer {
	return unsafe.Pointer(&s[0])
}

// PinSinkPtr keeps the PinSlice benches' results alive (no dead-code elimination).
var PinSinkPtr unsafe.Pointer
