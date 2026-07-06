// SPDX-License-Identifier: EUPL-1.2

// Encoding helpers for the Core framework.
// Wraps encoding/hex so consumers can use core primitives for common
// byte-string encodings.
package core

import (
	"encoding/base64"
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
