// SPDX-License-Identifier: EUPL-1.2

// Zero-copy string ↔ []byte views for the Core framework.
// SPOR file for unsafe-package usage — anywhere in core/go that
// reaches for unsafe.String / unsafe.Slice / unsafe.StringData /
// unsafe.SliceData routes through here.
//
// These primitives let consumers shave a copy when handing a string
// to a read-only byte consumer (hash, hmac, gzip.Write, Writer.Write,
// json.Unmarshal) and vice-versa when materialising a fresh buffer
// into a returned string.
//
//	digest := sha256.Sum256(core.AsBytes(payload))   // no string→[]byte copy
//	out    := core.AsString(buf.Bytes())             // no []byte→string copy
//
// SAFETY CONTRACT — read this before using:
//
//   - AsBytes: the returned []byte is a *view* into the source string's
//     memory. NEVER mutate the slice. NEVER let it outlive the source
//     string. Use it only when the immediate consumer is read-only
//     (hash, hmac, Write*, encode).
//
//   - AsString: the returned string is a view into the source []byte's
//     memory. NEVER mutate the source []byte after taking the view, and
//     NEVER let the caller see the original []byte. Use it only when
//     converting a freshly built, single-owner buffer to a returned
//     string (e.g. strings.Builder.String() does exactly this internally).
//
// Misuse corrupts memory or returns wrong results from string operations.
// When in doubt, use `[]byte(s)` / `string(b)` — the explicit copy is
// cheap up to a few hundred bytes and bulletproof.

package core

import (
	"unsafe"
)

// AsBytes returns a read-only []byte view of s without copying.
// See the package-level safety contract above before using.
//
//	digest := sha256.Sum256(core.AsBytes("payload"))
//	_, _ = w.Write(core.AsBytes(line))
func AsBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// AsString returns a string view of b without copying. The source b
// must be freshly built and never referenced again by the caller.
// See the package-level safety contract above before using.
//
//	var buf bytes.Buffer
//	buf.WriteString("agent ")
//	buf.WriteString(name)
//	return core.AsString(buf.Bytes())
func AsString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// PinnedView pins a Go slice to a stable address so its first-element
// pointer can be safely handed to C across the cgo boundary. The GC
// is prevented from moving the backing array while the view is held,
// which lets C callers retain the pointer beyond a single call (e.g.
// async kernels, model weights mlx keeps a reference to). Always pair
// with Release; the pinner table holds the slice live until then.
//
// Use for slices that:
//   - C may retain across more than one cgo invocation, OR
//   - are large enough that copy-into-C-memory dominates the call.
//
// For one-shot reads where C consumes the pointer during the call,
// the cgo runtime already prevents GC movement — use
// unsafe.SliceData / unsafe.Pointer directly without a PinnedView.
//
// SAFETY CONTRACT:
//
//   - The slice must NOT contain Go pointers (only numeric/byte
//     elements). cgo memory rules forbid passing nested Go pointers
//     through C; runtime.Pinner will panic in race mode if violated.
//   - The slice must outlive every C use of Ptr(). Caller owns
//     lifetime; PinnedView keeps the slice live, but does not own
//     it. Don't reslice or grow the source after pinning.
//   - Release exactly once. Double-release is a no-op; missing
//     release leaks the pin until process exit.
//
// Example:
//
//	var view core.PinnedView
//	core.PinSlice(weights, &view)
//	defer view.Release()
//	C.kernel_run(view.Ptr(), C.size_t(view.Bytes()))
type PinnedView struct {
	pinner Pinner
	ptr    unsafe.Pointer
	length int
	bytes  int
	active bool
}

// PinSlice pins slice's backing array and populates view in place.
// view must point to a stack or heap PinnedView (typically a local
// variable). The function is safe to call with an empty slice; in
// that case view is left zero-valued and Release is a no-op.
//
//	var view core.PinnedView
//	core.PinSlice(indices, &view)
//	defer view.Release()
//	C.fn(view.Ptr(), C.size_t(view.Len()))
func PinSlice[T any](slice []T, view *PinnedView) {
	if view == nil {
		return
	}
	if len(slice) == 0 {
		*view = PinnedView{}
		return
	}
	first := &slice[0]
	view.pinner.Pin(first)
	view.ptr = unsafe.Pointer(first)
	view.length = len(slice)
	var zero T
	view.bytes = int(unsafe.Sizeof(zero)) * len(slice)
	view.active = true
}

// Ptr returns the pinned start-of-slice pointer for C consumption.
// Returns nil on a zero or released view.
func (p *PinnedView) Ptr() unsafe.Pointer {
	if p == nil || !p.active {
		return nil
	}
	return p.ptr
}

// Len returns the element count of the pinned slice.
func (p *PinnedView) Len() int {
	if p == nil {
		return 0
	}
	return p.length
}

// Bytes returns the byte length of the pinned slice.
func (p *PinnedView) Bytes() int {
	if p == nil {
		return 0
	}
	return p.bytes
}

// Active reports whether the view holds a live pin.
func (p *PinnedView) Active() bool {
	if p == nil {
		return false
	}
	return p.active
}

// Release unpins the slice. Safe to call on a zero view or repeatedly.
func (p *PinnedView) Release() {
	if p == nil || !p.active {
		return
	}
	p.pinner.Unpin()
	p.ptr = nil
	p.length = 0
	p.bytes = 0
	p.active = false
}
