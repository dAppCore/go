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

// StringReader is an alias for strings.Reader — the io.Reader/Seeker
// over an in-memory string returned by NewReader. Lets consumers
// declare reader-typed fields without importing strings. Named
// StringReader (not Reader) because core.Reader already aliases
// io.Reader.
//
//	var r *core.StringReader = core.NewReader("payload")
type StringReader = strings.Reader

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

// IndexAny returns the byte position of the first occurrence in s of
// any Unicode code point in chars, or -1 when none are present. The
// character-class form of Index.
//
//	core.IndexAny("a/b\\c", "/\\")  // 1
func IndexAny(s, chars string) int {
	return strings.IndexAny(s, chars)
}

// ContainsAny reports whether any Unicode code point in chars is in s.
//
//	core.ContainsAny("user@host", "@:")  // true
func ContainsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

// ContainsRune reports whether the Unicode code point r is in s.
//
//	core.ContainsRune("café", 'é')  // true
func ContainsRune(s string, r rune) bool {
	return strings.ContainsRune(s, r)
}

// Count returns the number of non-overlapping instances of substr in s.
// An empty substr returns 1 + the rune count of s (stdlib semantics).
//
//	core.Count("a.b.c", ".")  // 2
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// EqualFold reports whether s and t are equal under simple Unicode
// case-folding — the case-insensitive comparison that avoids allocating
// two Lower copies just to compare them.
//
//	core.EqualFold("Bearer", "bearer")  // true
func EqualFold(s, t string) bool {
	return strings.EqualFold(s, t)
}

// Repeat returns a new string consisting of count copies of s. It
// panics when count is negative or the result overflows (stdlib
// contract); callers control both inputs so this stays infallible.
//
//	core.Repeat("=", 8)  // "========"
func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

// Fields splits s around runs of whitespace, returning the non-empty
// substrings. The whitespace-delimited tokeniser — Split needs an
// explicit separator and keeps empties, Fields collapses any run.
//
//	core.Fields("  go   test ./... ")  // ["go", "test", "./..."]
func Fields(s string) []string {
	return strings.Fields(s)
}

// Cut slices s around the first occurrence of sep, returning the text
// before and after it and whether sep was found. The idiomatic
// replacement for SplitN(s, sep, 2) when both halves are needed.
//
//	key, value, ok := core.Cut("port=8080", "=")  // "port", "8080", true
func Cut(s, sep string) (before, after string, found bool) {
	return strings.Cut(s, sep)
}

// CutPrefix returns s without the leading prefix and reports whether
// the prefix was present. Unlike TrimPrefix it tells the caller whether
// a cut actually happened.
//
//	rest, ok := core.CutPrefix("--verbose", "--")  // "verbose", true
func CutPrefix(s, prefix string) (after string, found bool) {
	return strings.CutPrefix(s, prefix)
}

// CutSuffix returns s without the trailing suffix and reports whether
// the suffix was present — the trailing-end sibling of CutPrefix.
//
//	base, ok := core.CutSuffix("main.go", ".go")  // "main", true
func CutSuffix(s, suffix string) (before string, found bool) {
	return strings.CutSuffix(s, suffix)
}
