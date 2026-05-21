// SPDX-License-Identifier: EUPL-1.2

// Correctness tests for the unsafe.go primitives, focused on
// PinnedView semantics — the contract is more subtle than the
// zero-copy view conversions, so the Good/Bad/Ugly triple is
// worth its weight.

package core_test

import (
	"testing"
	"unsafe"

	. "dappco.re/go"
)

// --- Good ---

// PinSlice gives a stable pointer to slice[0] that the caller can hand
// to C; Release unpins so the GC can reclaim the slice.
func TestPinSlice_Int32_PinsAddress(t *testing.T) {
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	if !view.Active() {
		t.Fatal("Active() returned false on a freshly pinned view")
	}
	if view.Len() != 4 {
		t.Fatalf("Len() = %d, want 4", view.Len())
	}
	if view.Bytes() != 16 {
		t.Fatalf("Bytes() = %d, want 16 (4 * sizeof(int32))", view.Bytes())
	}
	if view.Ptr() != unsafe.Pointer(&slice[0]) {
		t.Fatal("Ptr() does not match &slice[0] — pin failed")
	}
}

// PinSlice on []float32 reports the correct byte width (4) per element.
func TestPinSlice_Float32_Bytes(t *testing.T) {
	slice := make([]float32, 16)
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	if view.Bytes() != 64 {
		t.Fatalf("Bytes() = %d, want 64 (16 * sizeof(float32))", view.Bytes())
	}
}

// PinSlice on []byte reports the correct byte width (1) per element.
func TestPinSlice_Byte_Bytes(t *testing.T) {
	slice := []byte("payload")
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	if view.Bytes() != len(slice) {
		t.Fatalf("Bytes() = %d, want %d", view.Bytes(), len(slice))
	}
}

// --- Bad ---

// PinSlice on an empty slice leaves the view zero-valued — Active is
// false, Ptr is nil, Release is a no-op.
func TestPinSlice_EmptySlice_ZeroView(t *testing.T) {
	var slice []int32
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()

	if view.Active() {
		t.Fatal("Active() returned true for empty input")
	}
	if view.Ptr() != nil {
		t.Fatal("Ptr() not nil for empty input")
	}
	if view.Len() != 0 || view.Bytes() != 0 {
		t.Fatalf("Len/Bytes non-zero: Len=%d Bytes=%d", view.Len(), view.Bytes())
	}
}

// PinSlice with a nil view pointer must not panic — defensive contract
// for callers that pass through a forgotten initialiser.
func TestPinSlice_NilOutPointer_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("PinSlice panicked on nil out pointer: %v", r)
		}
	}()
	slice := []int32{1, 2, 3}
	PinSlice[int32](slice, nil)
}

// --- Ugly ---

// Release on a zero view is safe. Release on a released view is also
// safe — the second call is a no-op, leaving the GC free to reclaim
// the slice exactly once.
func TestPinSlice_DoubleRelease_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Release double-call panicked: %v", r)
		}
	}()
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	view.Release()
	view.Release()
	if view.Active() {
		t.Fatal("Active() true after Release")
	}
}

// Nil receiver methods must report safe zero values rather than
// panicking — lets callers route through helpers that may return a
// nil PinnedView for failure paths without crashing.
func TestPinSlice_NilReceiver_SafeZero(t *testing.T) {
	var view *PinnedView
	if view.Active() {
		t.Fatal("Active() true on nil receiver")
	}
	if view.Ptr() != nil {
		t.Fatal("Ptr() not nil on nil receiver")
	}
	if view.Len() != 0 {
		t.Fatal("Len() not 0 on nil receiver")
	}
	if view.Bytes() != 0 {
		t.Fatal("Bytes() not 0 on nil receiver")
	}
	view.Release()
}
