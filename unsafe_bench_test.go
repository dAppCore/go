// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the zero-copy view primitives in unsafe.go.
// Per AX-11 — AsBytes / AsString are the SPOR for unsafe-package
// zero-copy conversion. Hot path: hash digest of a string, Writer.Write
// of a string body, ReadAll → string return. Each shaves the source
// length in alloc + copy cost from the call.
//
// Run:    go test -bench='BenchmarkAs(Bytes|String)' -benchmem -run='^$' .

package core_test

import (
	"unsafe"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	unsafeSinkBytes  []byte
	unsafeSinkString string
)

var (
	unsafeStr16   = "aaaaaaaaaaaaaaaa"
	unsafeStr1K   = makeString(1024)
	unsafeBytes16 = []byte("aaaaaaaaaaaaaaaa")
	unsafeBytes1K = func() []byte {
		out := make([]byte, 1024)
		for i := range out {
			out[i] = byte('a' + i%26)
		}
		return out
	}()
)

func makeString(n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte('a' + i%26)
	}
	return string(out)
}

// --- AsBytes (zero-copy) vs. []byte(s) baseline ---

func BenchmarkAsBytes_16B(b *B) {
	s := unsafeStr16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = AsBytes(s)
	}
}

func BenchmarkAsBytes_1KB(b *B) {
	s := unsafeStr1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = AsBytes(s)
	}
}

func BenchmarkAsBytes_Baseline_16B(b *B) {
	s := unsafeStr16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = []byte(s)
	}
}

func BenchmarkAsBytes_Baseline_1KB(b *B) {
	s := unsafeStr1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = []byte(s)
	}
}

// --- AsString (zero-copy) vs. string(b) baseline ---

func BenchmarkAsString_16B(b *B) {
	bs := unsafeBytes16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = AsString(bs)
	}
}

func BenchmarkAsString_1KB(b *B) {
	bs := unsafeBytes1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = AsString(bs)
	}
}

func BenchmarkAsString_Baseline_16B(b *B) {
	bs := unsafeBytes16
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = string(bs)
	}
}

func BenchmarkAsString_Baseline_1KB(b *B) {
	bs := unsafeBytes1K
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = string(bs)
	}
}

// --- Empty input fast paths ---

func BenchmarkAsBytes_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkBytes = AsBytes("")
	}
}

func BenchmarkAsString_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		unsafeSinkString = AsString(nil)
	}
}

// --- PinnedView: zero-copy slice handoff to C across cgo boundary ---
// The point of comparison is "make + copy into a fresh []T of matching
// element width" — the pattern open-coded across go-mlx's metal package
// for every tensor op. PinnedView is one Pin (one small alloc) plus a
// pointer view; copy-into-buf is one make plus N element copies. The
// crossover by element count is where the substrate decision lives.

var pinSinkPtr unsafe.Pointer

func BenchmarkPinSlice_Int32_4(b *B) {
	slice := []int32{1, 2, 3, 4}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var view PinnedView
		PinSlice(slice, &view)
		pinSinkPtr = view.Ptr()
		view.Release()
	}
}

func BenchmarkPinSlice_Int32_64(b *B) {
	slice := make([]int32, 64)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var view PinnedView
		PinSlice(slice, &view)
		pinSinkPtr = view.Ptr()
		view.Release()
	}
}

func BenchmarkPinSlice_Float32_2048(b *B) {
	slice := make([]float32, 2048)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var view PinnedView
		PinSlice(slice, &view)
		pinSinkPtr = view.Ptr()
		view.Release()
	}
}

// Baseline: open-coded make + copy that go-mlx uses for shape arrays
// today. Same element width on both sides, so this is the minimal
// cost the substrate replaces.
func BenchmarkPinSlice_Baseline_MakeCopy_Int32_4(b *B) {
	slice := []int32{1, 2, 3, 4}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		copied := make([]int32, len(slice))
		copy(copied, slice)
		pinSinkPtr = unsafe.Pointer(&copied[0])
	}
}

func BenchmarkPinSlice_Baseline_MakeCopy_Int32_64(b *B) {
	slice := make([]int32, 64)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		copied := make([]int32, len(slice))
		copy(copied, slice)
		pinSinkPtr = unsafe.Pointer(&copied[0])
	}
}

func BenchmarkPinSlice_Baseline_MakeCopy_Float32_2048(b *B) {
	slice := make([]float32, 2048)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		copied := make([]float32, len(slice))
		copy(copied, slice)
		pinSinkPtr = unsafe.Pointer(&copied[0])
	}
}

func BenchmarkPinSlice_Empty(b *B) {
	var slice []int32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var view PinnedView
		PinSlice(slice, &view)
		pinSinkPtr = view.Ptr()
		view.Release()
	}
}

var (
	pinSinkInt  int
	pinSinkBool bool
)

func BenchmarkPinnedView_Len(b *B) {
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pinSinkInt = view.Len()
	}
}

func BenchmarkPinnedView_Bytes(b *B) {
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pinSinkInt = view.Bytes()
	}
}

func BenchmarkPinnedView_Active(b *B) {
	slice := []int32{1, 2, 3, 4}
	var view PinnedView
	PinSlice(slice, &view)
	defer view.Release()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pinSinkBool = view.Active()
	}
}
