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

import "unsafe"

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
