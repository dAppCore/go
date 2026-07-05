// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func ExampleAsBytes() {
	b := AsBytes("homelab")
	Println(len(b))
	// Output: 7
}

func ExampleAsString() {
	s := AsString([]byte("codex"))
	Println(s)
	// Output: codex
}

// ExamplePinSlice pins a slice's backing array for the cgo boundary through `PinSlice`.
func ExamplePinSlice() {
	var view PinnedView
	PinSlice([]int32{1, 2, 3, 4}, &view)
	defer view.Release()
	Println(view.Len())
	// Output: 4
}

// ExamplePinnedView_Ptr returns the pinned start-of-slice pointer through `PinnedView.Ptr`.
func ExamplePinnedView_Ptr() {
	var view PinnedView
	PinSlice([]int32{1, 2, 3}, &view)
	defer view.Release()
	_ = view.Ptr() // hand to C across the cgo boundary
}

// ExamplePinnedView_Len returns the pinned element count through `PinnedView.Len`.
func ExamplePinnedView_Len() {
	var view PinnedView
	PinSlice([]int32{1, 2, 3}, &view)
	defer view.Release()
	Println(view.Len())
	// Output: 3
}

// ExamplePinnedView_Bytes returns the pinned byte length through `PinnedView.Bytes`.
func ExamplePinnedView_Bytes() {
	var view PinnedView
	PinSlice([]int32{1, 2, 3, 4}, &view)
	defer view.Release()
	Println(view.Bytes())
	// Output: 16
}

// ExamplePinnedView_Active reports whether the view holds a live pin through `PinnedView.Active`.
func ExamplePinnedView_Active() {
	var view PinnedView
	PinSlice([]int32{1}, &view)
	defer view.Release()
	Println(view.Active())
	// Output: true
}

// ExamplePinnedView_Release unpins the slice through `PinnedView.Release`.
func ExamplePinnedView_Release() {
	var view PinnedView
	PinSlice([]int32{1}, &view)
	view.Release()
	Println(view.Active())
	// Output: false
}
