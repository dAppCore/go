// SPDX-License-Identifier: EUPL-1.2

// Encoding helpers for the Core framework.
// Wraps encoding/base64, encoding/binary and encoding/hex so consumers can use
// core primitives for common byte-string encodings — this file is their sole
// owner (SPOR), so the fixed-width readers live here beside the text ones.
package core

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
)

// HexEncode returns src encoded as a lowercase hexadecimal string.
//
//	s := core.HexEncode([]byte("hello"))
//
// Zero-copy: skips the stdlib EncodeToString return-side copy by
// aliasing the freshly-allocated dst buffer via AsString. Saves one
// alloc per call — load-bearing on SHA256HexString and friends.
func HexEncode(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return AsString(dst)
}

// HexDecode decodes a hexadecimal string into bytes.
//
//	r := core.HexDecode("68656c6c6f")
//	if r.OK { b := r.Value.([]byte) }
func HexDecode(s string) Result {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Result{err, false}
	}
	return Result{b, true}
}

// Base64Encode returns src encoded as a standard base64 string.
//
//	s := core.Base64Encode([]byte("hello"))
//
// Zero-copy return: pre-allocates the encoded buffer and aliases via
// AsString to skip stdlib's return-side copy. Saves one alloc per call.
func Base64Encode(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return AsString(dst)
}

// Base64Decode decodes a standard base64 string into bytes.
//
//	r := core.Base64Decode("aGVsbG8=")
//	if r.OK { b := r.Value.([]byte) }
func Base64Decode(s string) Result {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Result{err, false}
	}
	return Result{b, true}
}

// Base64URLEncode returns src encoded as a URL-safe base64 string.
//
//	s := core.Base64URLEncode([]byte("hello"))
func Base64URLEncode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

// Base64URLDecode decodes a URL-safe base64 string into bytes.
//
//	r := core.Base64URLDecode("aGVsbG8=")
//	if r.OK { b := r.Value.([]byte) }
func Base64URLDecode(s string) Result {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return Result{err, false}
	}
	return Result{b, true}
}

// BigEndianUint64 reads the first eight bytes of b as a big-endian uint64.
//
//	n := core.BigEndianUint64(buf[:])
//
// A thin pass-through so callers on a hot path need not import
// encoding/binary — encode.go is its sole owner (SPOR). Inlinable, so the
// wrapper costs nothing over the direct call. Panics if b is shorter than
// eight bytes, exactly as the stdlib does.
func BigEndianUint64(b []byte) uint64 { return binary.BigEndian.Uint64(b) }

// LittleEndianUint64 reads the first eight bytes of b as a little-endian uint64.
//
//	lane := core.LittleEndianUint64(data[i*8:])
//
// The little-endian counterpart to BigEndianUint64; see it for why this lives
// here. Keccak absorbs its state in this order.
func LittleEndianUint64(b []byte) uint64 { return binary.LittleEndian.Uint64(b) }

// HexAppendEncode appends the hex encoding of src to dst and returns the
// extended slice.
//
//	buf = core.HexAppendEncode(buf, rnd[:])
//
// The append-style form matters: it lets a caller building a string in one
// buffer stay at a single allocation, where HexEncode would add its own.
func HexAppendEncode(dst, src []byte) []byte { return hex.AppendEncode(dst, src) }

// PutLittleEndianUint64 writes v into the first eight bytes of b, little-endian.
//
//	core.PutLittleEndianUint64(sum[i*8:], state[i])
//
// The write-side counterpart to LittleEndianUint64, here for the same reason.
func PutLittleEndianUint64(b []byte, v uint64) { binary.LittleEndian.PutUint64(b, v) }
