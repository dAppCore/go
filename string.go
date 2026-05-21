// SPDX-License-Identifier: EUPL-1.2

// String operations for the Core framework.
// Provides safe, predictable string helpers that downstream packages
// use directly — same pattern as Array[T] for slices.

package core

import (
	"html"
	"strings"
	"unicode/utf8"
)

// HasPrefix returns true if s starts with prefix.
//
//	core.HasPrefix("--verbose", "--")  // true
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix returns true if s ends with suffix.
//
//	core.HasSuffix("test.go", ".go")  // true
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// TrimPrefix removes prefix from s.
//
//	core.TrimPrefix("--verbose", "--")  // "verbose"
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// TrimSuffix removes suffix from s.
//
//	core.TrimSuffix("test.go", ".go")  // "test"
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// Clone returns a fresh copy of s, detached from its backing memory.
// Use this when s aliases a reusable buffer (e.g. via core.AsString
// over a scratch slice) and the result must outlive the buffer's
// next reuse — map keys, struct fields, channel sends.
//
//	key := core.Clone(core.AsString(scratch))  // map[key] = ... is safe
func Clone(s string) string {
	return strings.Clone(s)
}

// Contains returns true if s contains substr.
//
//	core.Contains("hello world", "world")  // true
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Split splits s by separator.
//
//	core.Split("a/b/c", "/")  // ["a", "b", "c"]
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// SplitN splits s by separator into at most n parts.
//
//	core.SplitN("key=value=extra", "=", 2)  // ["key", "value=extra"]
func SplitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

// Join joins parts with a separator into a single string.
//
//	core.Join("/", "deploy", "to", "homelab")      // "deploy/to/homelab"
//	core.Join(".", "cmd", "deploy", "description")  // "cmd.deploy.description"
func Join(sep string, parts ...string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	}
	// Pre-size the Builder to the exact final length so WriteString
	// never grows the internal buffer. The earlier implementation
	// chained Concat(result, sep, p) per pair which produced O(N²)
	// allocations and bytes-copied for N parts.
	n := len(sep) * (len(parts) - 1)
	for _, p := range parts {
		n += len(p)
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(parts[0])
	for _, p := range parts[1:] {
		b.WriteString(sep)
		b.WriteString(p)
	}
	return b.String()
}

// Replace replaces all occurrences of old with new in s.
//
//	core.Replace("deploy/to/homelab", "/", ".")  // "deploy.to.homelab"
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Lower returns s in lowercase.
//
//	core.Lower("HELLO")  // "hello"
func Lower(s string) string {
	// ASCII no-op fast path. strings.ToLower walks the Unicode case
	// table for every byte regardless of input — costly even when the
	// answer is "no change". Scan once; if the input is pure ASCII
	// without any uppercase, return it unchanged with zero allocations
	// and ~0.2ns/byte (a simple byte compare in a tight loop) instead
	// of strings.ToLower's full Unicode walk.
	//
	// For inputs that DO need work (ASCII with uppercase, or any non-
	// ASCII), fall through to strings.ToLower — it's already optimised
	// for that path and our Builder-based ASCII variant lost to it.
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' || c >= 0x80 {
			return strings.ToLower(s)
		}
	}
	return s
}

// Upper returns s in uppercase.
//
//	core.Upper("hello")  // "HELLO"
func Upper(s string) string {
	// ASCII no-op fast path — symmetric to Lower. See Lower's comment
	// for the rationale.
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' || c >= 0x80 {
			return strings.ToUpper(s)
		}
	}
	return s
}

// Trim removes leading and trailing whitespace.
//
//	core.Trim("  hello  ")  // "hello"
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimCutset removes any leading and trailing characters in cutset
// from s. The character-class form of Trim — TrimSpace handles only
// whitespace, TrimCutset handles arbitrary character sets.
//
//	core.TrimCutset("//path//", "/")     // "path"
//	core.TrimCutset("[name]", "[]")      // "name"
func TrimCutset(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// TrimLeft removes any leading characters in cutset from s.
//
//	core.TrimLeft("///path", "/")        // "path"
//	core.TrimLeft("---verbose", "-")     // "verbose"
func TrimLeft(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

// TrimRight removes any trailing characters in cutset from s.
//
//	core.TrimRight("path///", "/")       // "path"
//	core.TrimRight("hello!!!", "!")      // "hello"
func TrimRight(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}

// Index returns the byte position of the first occurrence of sep in s,
// or -1 when sep is absent.
//
//	core.Index("key=value", "=")         // 3
//	core.Index("nothing here", "?")      // -1
func Index(s, sep string) int {
	return strings.Index(s, sep)
}

// RuneCount returns the number of runes (unicode characters) in s.
//
//	core.RuneCount("hello")  // 5
//	core.RuneCount("🔥")     // 1
func RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}

// Builder is an alias for strings.Builder — the byte-by-byte string
// assembly type. Lets consumers declare builder-typed fields and locals
// without importing strings.
//
//	type Sink struct {
//	    out core.Builder
//	}
//
//	var b core.Builder
//	b.WriteString("ready")
type Builder = strings.Builder

// NewBuilder returns a new strings.Builder.
//
//	b := core.NewBuilder()
//	b.WriteString("hello")
//	b.String() // "hello"
func NewBuilder() *strings.Builder {
	return &strings.Builder{}
}

// NewReader returns a strings.NewReader for the given string.
//
//	r := core.NewReader("hello world")
func NewReader(s string) *strings.Reader {
	return strings.NewReader(s)
}

// Concat joins variadic string parts into one string.
// Hook point for validation, sanitisation, and security checks.
//
//	core.Concat("cmd.", "deploy.to.homelab", ".description")
//	core.Concat("https://", host, "/api/v1")
func Concat(parts ...string) string {
	// Pre-size to the exact final length so WriteString never grows
	// the internal buffer. Without this the Builder doubles its
	// backing array as parts append, costing 1-3 extra allocations
	// on every call.
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	var b strings.Builder
	b.Grow(n)
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

// HTMLEscape returns s with special HTML characters escaped.
//
//	escaped := core.HTMLEscape(`<a href="/search?q=go&lang=en">Go</a>`)
func HTMLEscape(s string) string {
	return html.EscapeString(s)
}

// HTMLUnescape returns s with HTML character references unescaped.
//
//	unescaped := core.HTMLUnescape("&lt;strong&gt;Go&lt;/strong&gt;")
func HTMLUnescape(s string) string {
	return html.UnescapeString(s)
}

// LastIndex returns the index of the last instance of substr in s, or
// -1 if substr is not present. Mirrors strings.LastIndex; pair with
// Index when consumer code needs both ends of a delimiter.
//
//	colon := core.LastIndex("host.example.com:8080", ":")  // 16
func LastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}
