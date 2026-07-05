// SPDX-License-Identifier: EUPL-1.2

// Correctness tests for the unsafe.go primitives. PinnedView's contract
// (pin / inspect / release, plus nil-receiver and zero-view safety) is
// the subtle part, so each accessor carries its own Good/Bad/Ugly.

package core_test

import (
	"unsafe"

	. "dappco.re/go"
)

// --- AsBytes ---

func TestUnsafe_AsBytes_Good(t *T) {
	AssertEqual(t, []byte("payload"), AsBytes("payload"))
}

func TestUnsafe_AsBytes_Bad(t *T) {
	// An empty string yields a nil slice, not an empty non-nil one.
	AssertNil(t, AsBytes(""))
}

func TestUnsafe_AsBytes_Ugly(t *T) {
	// Multi-byte UTF-8 is viewed byte-for-byte, no decoding.
	AssertEqual(t, []byte("café"), AsBytes("café"))
}

// --- AsString ---

func TestUnsafe_AsString_Good(t *T) {
	AssertEqual(t, "payload", AsString([]byte("payload")))
}

func TestUnsafe_AsString_Bad(t *T) {
	// Empty and nil byte slices both yield the empty string.
	AssertEqual(t, "", AsString(nil))
	AssertEqual(t, "", AsString([]byte{}))
}

func TestUnsafe_AsString_Ugly(t *T) {
	// Round-trips through AsBytes byte-for-byte, including multi-byte runes.
	AssertEqual(t, "café", AsString(AsBytes("café")))
}

// --- PinSlice ---

func TestUnsafe_PinSlice_Good(t *T) {
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	AssertTrue(t, view.Active())
	AssertTrue(t, view.Ptr() == unsafe.Pointer(&slice[0]))
	AssertEqual(t, 4, view.Len())
	AssertEqual(t, 16, view.Bytes())
}

func TestUnsafe_PinSlice_Bad(t *T) {
	// An empty slice leaves the view zero-valued; Release is a no-op.
	var slice []int32
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	AssertFalse(t, view.Active())
	AssertTrue(t, view.Ptr() == nil)
	AssertEqual(t, 0, view.Len())
	AssertEqual(t, 0, view.Bytes())
}

func TestUnsafe_PinSlice_Ugly(t *T) {
	// A nil out-pointer is a defensive no-op, not a panic.
	AssertNotPanics(t, func() {
		PinSlice[int32]([]int32{1, 2, 3}, nil)
	})
}

// --- PinnedView.Ptr ---

func TestUnsafe_PinnedView_Ptr_Good(t *T) {
	slice := []int32{1, 2, 3}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	AssertTrue(t, view.Ptr() == unsafe.Pointer(&slice[0]))
}

func TestUnsafe_PinnedView_Ptr_Bad(t *T) {
	// A zero view has no pinned address.
	var view PinnedView
	AssertTrue(t, view.Ptr() == nil)
}

func TestUnsafe_PinnedView_Ptr_Ugly(t *T) {
	// A nil receiver returns nil rather than panicking.
	var view *PinnedView
	AssertTrue(t, view.Ptr() == nil)
}

// --- PinnedView.Len ---

func TestUnsafe_PinnedView_Len_Good(t *T) {
	slice := []int32{1, 2, 3, 4, 5}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	AssertEqual(t, 5, view.Len())
}

func TestUnsafe_PinnedView_Len_Bad(t *T) {
	var view PinnedView
	AssertEqual(t, 0, view.Len())
}

func TestUnsafe_PinnedView_Len_Ugly(t *T) {
	var view *PinnedView
	AssertEqual(t, 0, view.Len())
}

// --- PinnedView.Bytes ---

func TestUnsafe_PinnedView_Bytes_Good(t *T) {
	// Byte width scales with element type.
	var i32 PinnedView
	PinSlice([]int32{1, 2, 3, 4}, &i32)
	defer i32.Release()
	AssertEqual(t, 16, i32.Bytes()) // 4 * sizeof(int32)

	var f32 PinnedView
	PinSlice(make([]float32, 16), &f32)
	defer f32.Release()
	AssertEqual(t, 64, f32.Bytes()) // 16 * sizeof(float32)

	var b PinnedView
	PinSlice([]byte("payload"), &b)
	defer b.Release()
	AssertEqual(t, 7, b.Bytes()) // 1 byte per element
}

func TestUnsafe_PinnedView_Bytes_Bad(t *T) {
	var view PinnedView
	AssertEqual(t, 0, view.Bytes())
}

func TestUnsafe_PinnedView_Bytes_Ugly(t *T) {
	var view *PinnedView
	AssertEqual(t, 0, view.Bytes())
}

// --- PinnedView.Active ---

func TestUnsafe_PinnedView_Active_Good(t *T) {
	var view PinnedView
	PinSlice([]int32{1, 2, 3}, &view)
	defer view.Release()

	AssertTrue(t, view.Active())
}

func TestUnsafe_PinnedView_Active_Bad(t *T) {
	var view PinnedView
	AssertFalse(t, view.Active())
}

func TestUnsafe_PinnedView_Active_Ugly(t *T) {
	var view *PinnedView
	AssertFalse(t, view.Active())
}

// --- PinnedView.Release ---

func TestUnsafe_PinnedView_Release_Good(t *T) {
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)

	view.Release()

	AssertFalse(t, view.Active())
	AssertTrue(t, view.Ptr() == nil)
}

func TestUnsafe_PinnedView_Release_Bad(t *T) {
	// Releasing a zero view is a safe no-op.
	var view PinnedView
	AssertNotPanics(t, func() { view.Release() })
}

func TestUnsafe_PinnedView_Release_Ugly(t *T) {
	// Double release is a no-op, not a double-unpin panic.
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	AssertNotPanics(t, func() {
		view.Release()
		view.Release()
	})
	AssertFalse(t, view.Active())
}
