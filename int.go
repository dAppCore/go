// SPDX-License-Identifier: EUPL-1.2

// Integer conversion helpers for the Core framework.
// Wraps strconv so consumers don't import it directly.
package core

import "strconv"

// Atoi converts a decimal string to an int.
//
//	r := core.Atoi("42")
//	if r.OK { n := r.Value.(int) }
func Atoi(s string) Result {
	i, err := strconv.Atoi(s)
	if err != nil {
		return Result{err, false}
	}
	return Result{i, true}
}

// Itoa converts an int to a decimal string.
//
//	s := core.Itoa(42)
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// FormatInt converts an int64 to a string in the given base.
//
//	s := core.FormatInt(255, 16) // "ff"
func FormatInt(i int64, base int) string {
	return strconv.FormatInt(i, base)
}

// FormatUint converts a uint64 to a string in the given base.
//
//	s := core.FormatUint(255, 16) // "ff"
func FormatUint(i uint64, base int) string {
	return strconv.FormatUint(i, base)
}

// ParseInt converts a string in the given base and bit size to an int64.
//
//	r := core.ParseInt("ff", 16, 64)
//	if r.OK { n := r.Value.(int64) }
func ParseInt(s string, base int, bitSize int) Result {
	i, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		return Result{err, false}
	}
	return Result{i, true}
}

// The Append* family formats into a caller-supplied buffer and returns the
// extended slice, so a caller assembling a line stays at one allocation where
// the Format*/Itoa forms would each add their own. int.go is strconv's sole
// owner (SPOR), so they live here rather than at the call sites that need them
// — the structured logger's field encoder and ID(). All are inlinable
// pass-throughs, so routing through core costs nothing over the direct call.

// AppendInt appends the base-`base` text of i to dst.
//
//	buf = core.AppendInt(buf[:0], -42, 10)
func AppendInt(dst []byte, i int64, base int) []byte { return strconv.AppendInt(dst, i, base) }

// AppendUint appends the base-`base` text of i to dst.
//
//	buf = core.AppendUint(buf, counter.Add(1), 10)
func AppendUint(dst []byte, i uint64, base int) []byte { return strconv.AppendUint(dst, i, base) }

// AppendBool appends "true" or "false" to dst.
//
//	buf = core.AppendBool(buf[:0], ok)
func AppendBool(dst []byte, b bool) []byte { return strconv.AppendBool(dst, b) }

// AppendFloat appends the text of f to dst, formatted as by FormatFloat.
//
//	buf = core.AppendFloat(buf[:0], f, 'g', -1, 64)
func AppendFloat(dst []byte, f float64, fmt byte, prec, bitSize int) []byte {
	return strconv.AppendFloat(dst, f, fmt, prec, bitSize)
}

// AppendQuote appends the double-quoted Go string literal of s to dst.
//
//	buf = core.AppendQuote(buf[:0], v)
func AppendQuote(dst []byte, s string) []byte { return strconv.AppendQuote(dst, s) }
