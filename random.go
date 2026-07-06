// SPDX-License-Identifier: EUPL-1.2

// Random values for the Core framework.
// Provides cryptographically secure helpers and fast non-crypto selection.

package core

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	fastrand "math/rand/v2"
)

// RandomBytes returns n cryptographically secure random bytes wrapped in
// a Result. Returns OK=false when n is negative or the OS entropy source
// fails. Code is "random.length.invalid" or "random.entropy.failed".
//
//	r := core.RandomBytes(32)
//	if !r.OK { return r }
//	token := r.Value.([]byte)
func RandomBytes(n int) Result {
	if n < 0 {
		return Result{Value: NewCode("random.length.invalid", "RandomBytes: negative length"), OK: false}
	}
	out := make([]byte, n)
	if _, err := cryptorand.Read(out); err != nil {
		return Result{Value: WrapCode(err, "random.entropy.failed", "RandomBytes", "OS entropy source failed"), OK: false}
	}
	return Result{Value: out, OK: true}
}

// RandomString returns n cryptographically secure random bytes encoded
// as lowercase hex. The Value is a string of length n*2. Code mirrors
// RandomBytes when the underlying source fails.
//
//	r := core.RandomString(16)
//	if r.OK { token := r.Value.(string) }
func RandomString(n int) Result {
	if n < 0 {
		return Result{Value: NewCode("random.length.invalid", "RandomString: negative length"), OK: false}
	}
	// Inline the bytes + hex pipeline; the previous version went through
	// RandomBytes() and paid the intermediate Result.Value boxing plus
	// the []byte interface assertion. One alloc less per call.
	bytes := make([]byte, n)
	if _, err := cryptorand.Read(bytes); err != nil {
		return Result{Value: WrapCode(err, "random.entropy.failed", "RandomString", "OS entropy source failed"), OK: false}
	}
	return Result{Value: HexEncode(bytes), OK: true}
}

// RandomInt returns a cryptographically secure integer in the half-open
// range [min, max). Returns OK=false when max <= min (Code
// "random.range.empty") or when crypto/rand fails (Code
// "random.entropy.failed").
//
//	r := core.RandomInt(10, 20)
//	if r.OK { delay := r.Value.(int) }
func RandomInt(min, max int) Result {
	if max <= min {
		return Result{Value: NewCode("random.range.empty", "RandomInt: empty range"), OK: false}
	}
	// uint64 rejection sampling. The previous big.Int path allocated 3-4
	// big.Ints per call; for any int range (which by definition fits in
	// int64), uint64 arithmetic is exact and free.
	//
	// Sample 8 random bytes, reject values >= threshold (the largest
	// multiple of span that fits in uint64) to eliminate modulo bias.
	// For typical spans (≪ uint64Max) the rejection probability is
	// effectively zero — the loop runs once.
	span := uint64(max - min)
	threshold := ^uint64(0) - (^uint64(0) % span)
	var buf [8]byte
	for {
		if _, err := cryptorand.Read(buf[:]); err != nil {
			return Result{Value: WrapCode(err, "random.entropy.failed", "RandomInt", "OS entropy source failed"), OK: false}
		}
		n := binary.BigEndian.Uint64(buf[:])
		if n < threshold {
			return Result{Value: int(n%span) + min, OK: true}
		}
	}
}

// RandPick returns a pseudo-random item from items.
// It panics when items is empty; use it for fast non-crypto selection only.
//
//	choice := core.RandPick([]string{"red", "green", "blue"})
func RandPick[T any](items []T) T {
	if len(items) == 0 {
		panic("core.RandPick: empty slice")
	}
	return items[RandIntn(len(items))]
}

// RandIntn returns a pseudo-random integer in [0, n).
// It panics when n is less than or equal to zero; use it for fast non-crypto values.
//
//	index := core.RandIntn(len(items))
func RandIntn(n int) int {
	return fastrand.IntN(n)
}

// RandRead fills b with cryptographically secure random bytes. Returns
// Result.OK true on success; on failure r.Value holds the underlying
// error. Use this when callers need to fill an existing slice rather
// than allocate via RandomBytes.
//
//	buf := make([]byte, 32)
//	if r := core.RandRead(buf); !r.OK {
//	    return core.Fail(core.E("seed", "rand: "+r.Error(), nil))
//	}
func RandRead(b []byte) Result {
	if _, err := cryptorand.Read(b); err != nil {
		return Result{err, false}
	}
	return Result{nil, true}
}
