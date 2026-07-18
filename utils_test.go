package core_test

import (
	. "dappco.re/go"
)

// --- ID ---

func TestUtils_ID_Good(t *T) {
	id := ID()
	AssertTrue(t, HasPrefix(id, "id-"))
	AssertTrue(t, len(id) > 5, "ID should have counter + random suffix")
}

func TestUtils_ID_Good_Unique(t *T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := ID()
		AssertFalse(t, seen[id], "ID collision: %s", id)
		seen[id] = true
	}
}

func TestUtils_ID_Ugly_CounterMonotonic(t *T) {
	// IDs should contain increasing counter values
	id1 := ID()
	id2 := ID()
	// Both should start with "id-" and have different counter parts
	AssertNotEqual(t, id1, id2)
	AssertTrue(t, HasPrefix(id1, "id-"))
	AssertTrue(t, HasPrefix(id2, "id-"))
}

// --- ValidateName ---

func TestUtils_ValidateName_Good(t *T) {
	r := ValidateName("brain")
	AssertTrue(t, r.OK)
	AssertEqual(t, "brain", r.Value)
}

func TestUtils_ValidateName_Good_WithDots(t *T) {
	// Dots are valid — used for action namespacing. On success the Result
	// carries the unchanged name back verbatim, not a copy or normalised form.
	r := ValidateName("process.run")
	AssertTrue(t, r.OK, "dots in names are valid — used for action namespacing")
	AssertEqual(t, "process.run", r.Value)

	// Multiple dots and version-style names are equally valid.
	multi := ValidateName("a.b.c.d")
	AssertTrue(t, multi.OK)
	AssertEqual(t, "a.b.c.d", multi.Value)

	ver := ValidateName("plugin.v1.2.3")
	AssertTrue(t, ver.OK)
	AssertEqual(t, "plugin.v1.2.3", ver.Value)
}

func TestUtils_ValidateName_Bad_Empty(t *T) {
	// The empty name is rejected; the Result carries a non-nil *Err whose
	// message names the failure and whose operation is "validate".
	r := ValidateName("")
	AssertFalse(t, r.OK)
	AssertNotNil(t, r.Value)
	err := r.Value.(error)
	AssertError(t, err, "invalid name")
	AssertContains(t, err.Error(), "validate")
	// Empty is an "invalid name", not a "path separator" rejection.
	AssertNotContains(t, err.Error(), "path separator")
}

func TestUtils_ValidateName_Bad_Dot(t *T) {
	// A lone "." (current dir) is rejected as an invalid name.
	r := ValidateName(".")
	AssertFalse(t, r.OK)
	err := r.Value.(error)
	AssertError(t, err, "invalid name")
	// A single dot is NOT a parent reference, so the message has no "..".
	AssertNotContains(t, err.Error(), "..")
	AssertNotContains(t, err.Error(), "path separator")
}

func TestUtils_ValidateName_Bad_DotDot(t *T) {
	// ".." (parent dir) is rejected as an invalid name, and the offending
	// value is echoed into the message — distinguishing it from the "." case.
	r := ValidateName("..")
	AssertFalse(t, r.OK)
	err := r.Value.(error)
	AssertError(t, err, "invalid name")
	AssertContains(t, err.Error(), "..")
	AssertNotContains(t, err.Error(), "path separator")
}

func TestUtils_ValidateName_Bad_Slash(t *T) {
	// Forward slashes are rejected as path separators, whether part of a
	// traversal, a relative path, or an absolute path.
	r := ValidateName("../escape")
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "path separator")

	rel := ValidateName("foo/bar")
	AssertFalse(t, rel.OK)
	AssertError(t, rel.Value.(error), "path separator")

	abs := ValidateName("/etc/passwd")
	AssertFalse(t, abs.OK)
	AssertError(t, abs.Value.(error), "path separator")
}

func TestUtils_ValidateName_Ugly_Backslash(t *T) {
	// Windows-style backslashes are rejected as path separators too, so a
	// name cannot smuggle a directory component on any platform.
	r := ValidateName("windows\\path")
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "path separator")

	deep := ValidateName("a\\b\\c")
	AssertFalse(t, deep.OK)
	AssertError(t, deep.Value.(error), "path separator")
}

// --- SanitisePath ---

func TestUtils_SanitisePath_Good(t *T) {
	// Returns the final path element for ordinary paths.
	AssertEqual(t, "file.txt", SanitisePath("/some/path/file.txt"))
	AssertEqual(t, "readme.md", SanitisePath("readme.md"))
	AssertEqual(t, "agent.json", SanitisePath("deploy/to/agent.json"))
}

func TestUtils_SanitisePath_Bad_Empty(t *T) {
	// Inputs whose final element is empty or a dot segment collapse to the
	// "invalid" sentinel.
	AssertEqual(t, "invalid", SanitisePath(""))
	AssertEqual(t, "invalid", SanitisePath("deploy/.."))
	// A trailing parent reference anywhere in the path also collapses, as
	// does a trailing current-dir segment.
	AssertEqual(t, "invalid", SanitisePath("a/b/.."))
	AssertEqual(t, "invalid", SanitisePath("foo/."))
}

func TestUtils_SanitisePath_Bad_DotDot(t *T) {
	// Bare dot segments are rejected as the "invalid" sentinel.
	AssertEqual(t, "invalid", SanitisePath(".."))
	AssertEqual(t, "invalid", SanitisePath("."))
	// A trailing separator does not save them — the base is still a dot.
	AssertEqual(t, "invalid", SanitisePath("../"))
	AssertEqual(t, "invalid", SanitisePath("./"))
}

func TestUtils_SanitisePath_Ugly_Traversal(t *T) {
	// Traversal prefixes are stripped — only the final element survives.
	AssertEqual(t, "passwd", SanitisePath("../../etc/passwd"))
	AssertEqual(t, "secret.key", SanitisePath("../../../home/user/.ssh/secret.key"))
	// A trailing slash is trimmed before taking the base, and interior
	// ".." segments are likewise discarded — only the last real element wins.
	AssertEqual(t, "c", SanitisePath("a/b/c/"))
	AssertEqual(t, "bar", SanitisePath("../foo/../bar"))
}

// --- FilterArgs ---

func TestUtils_FilterArgs_Good(t *T) {
	args := []string{"deploy", "", "to", "-test.v", "homelab", "-test.paniconexit0"}
	clean := FilterArgs(args)
	AssertEqual(t, []string{"deploy", "to", "homelab"}, clean)
}

func TestUtils_FilterArgs_Empty_Good(t *T) {
	// A nil slice yields a nil slice (no allocation), as does an empty,
	// non-nil slice — there is nothing to keep.
	AssertNil(t, FilterArgs(nil))
	AssertNil(t, FilterArgs([]string{}))
	AssertEmpty(t, FilterArgs([]string{}))
}

// --- ParseFlag ---

func TestUtils_ParseFlag_ShortValid_Good(t *T) {
	// Single letter
	k, v, ok := ParseFlag("-v")
	AssertTrue(t, ok)
	AssertEqual(t, "v", k)
	AssertEqual(t, "", v)

	// Single emoji
	k, v, ok = ParseFlag("-🔥")
	AssertTrue(t, ok)
	AssertEqual(t, "🔥", k)
	AssertEqual(t, "", v)

	// Short with value
	k, v, ok = ParseFlag("-p=8080")
	AssertTrue(t, ok)
	AssertEqual(t, "p", k)
	AssertEqual(t, "8080", v)
}

func TestUtils_ParseFlag_ShortInvalid_Bad(t *T) {
	// Multiple chars with single dash — invalid
	_, _, ok := ParseFlag("-verbose")
	AssertFalse(t, ok)

	_, _, ok = ParseFlag("-port")
	AssertFalse(t, ok)
}

func TestUtils_ParseFlag_LongValid_Good(t *T) {
	k, v, ok := ParseFlag("--verbose")
	AssertTrue(t, ok)
	AssertEqual(t, "verbose", k)
	AssertEqual(t, "", v)

	k, v, ok = ParseFlag("--port=8080")
	AssertTrue(t, ok)
	AssertEqual(t, "port", k)
	AssertEqual(t, "8080", v)
}

func TestUtils_ParseFlag_LongInvalid_Bad(t *T) {
	// Single char with double dash — invalid; key and value come back empty.
	k, v, ok := ParseFlag("--v")
	AssertFalse(t, ok)
	AssertEqual(t, "", k)
	AssertEqual(t, "", v)

	// Any other single-rune long name is equally rejected.
	_, _, ok = ParseFlag("--x")
	AssertFalse(t, ok)

	// A bare "--" has a zero-length name and is rejected.
	k, v, ok = ParseFlag("--")
	AssertFalse(t, ok)
	AssertEqual(t, "", k)
	AssertEqual(t, "", v)

	// "--=value" has an empty name even though a value is present — invalid.
	_, _, ok = ParseFlag("--=value")
	AssertFalse(t, ok)
}

func TestUtils_ParseFlag_NotAFlag_Bad(t *T) {
	_, _, ok := ParseFlag("hello")
	AssertFalse(t, ok)

	_, _, ok = ParseFlag("")
	AssertFalse(t, ok)
}

// --- IsFlag ---

func TestUtils_IsFlag_Good(t *T) {
	AssertTrue(t, IsFlag("-v"))
	AssertTrue(t, IsFlag("--verbose"))
	AssertTrue(t, IsFlag("-"))
}

func TestUtils_IsFlag_Bad(t *T) {
	// Only a leading dash makes a flag; bare words, numbers, filenames, and
	// strings with an interior dash are all non-flags.
	AssertFalse(t, IsFlag("hello"))
	AssertFalse(t, IsFlag(""))
	AssertFalse(t, IsFlag("deploy"))
	AssertFalse(t, IsFlag("123"))
	AssertFalse(t, IsFlag("file.txt"))
	AssertFalse(t, IsFlag("a-b"))
}

// --- Arg ---

func TestUtils_Arg_String_Good(t *T) {
	r := Arg(0, "hello", 42, true)
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello", r.Value)
}

func TestUtils_Arg_Int_Good(t *T) {
	r := Arg(1, "hello", 42, true)
	AssertTrue(t, r.OK)
	AssertEqual(t, 42, r.Value)
}

func TestUtils_Arg_Bool_Good(t *T) {
	r := Arg(2, "hello", 42, true)
	AssertTrue(t, r.OK)
	AssertEqual(t, true, r.Value)
}

func TestUtils_Arg_UnsupportedType_Good(t *T) {
	r := Arg(0, 3.14)
	AssertTrue(t, r.OK)
	AssertEqual(t, 3.14, r.Value)
}

func TestUtils_Arg_OutOfBounds_Bad(t *T) {
	r := Arg(5, "only", "two")
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestUtils_Arg_NoArgs_Bad(t *T) {
	r := Arg(0)
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestUtils_Arg_ErrorDetection_Good(t *T) {
	err := NewError("fail")
	r := Arg(0, err)
	AssertTrue(t, r.OK)
	AssertEqual(t, err, r.Value)
}

// --- ArgString ---

func TestUtils_ArgString_Good(t *T) {
	AssertEqual(t, "hello", ArgString(0, "hello", 42))
	AssertEqual(t, "world", ArgString(1, "hello", "world"))
	// The last position is reachable, and an empty string stored at an
	// in-bounds index is returned faithfully (not treated as "missing").
	AssertEqual(t, "c", ArgString(2, "a", "b", "c"))
	AssertEqual(t, "", ArgString(0, ""))
}

func TestUtils_ArgString_WrongType_Bad(t *T) {
	// A non-string at the requested index yields "" for every non-string
	// kind, and only the requested index is type-checked.
	AssertEqual(t, "", ArgString(0, 42))
	AssertEqual(t, "", ArgString(0, true))
	AssertEqual(t, "", ArgString(0, 3.14))
	AssertEqual(t, "", ArgString(1, "ok", 99))
}

func TestUtils_ArgString_OutOfBounds_Bad(t *T) {
	// Past the end returns ""; so does the boundary index == len, and the
	// no-args case.
	AssertEqual(t, "", ArgString(3, "only"))
	AssertEqual(t, "", ArgString(1, "only"))
	AssertEqual(t, "", ArgString(0))
}

// --- ArgInt ---

func TestUtils_ArgInt_Good(t *T) {
	AssertEqual(t, 42, ArgInt(0, 42, "hello"))
	AssertEqual(t, 99, ArgInt(1, 0, 99))
	// Zero and negative ints stored in bounds are returned as-is.
	AssertEqual(t, 0, ArgInt(0, 0, "x"))
	AssertEqual(t, -5, ArgInt(0, -5))
	AssertEqual(t, 7, ArgInt(2, "a", "b", 7))
}

func TestUtils_ArgInt_WrongType_Bad(t *T) {
	// The extractor matches the exact `int` type only: strings, floats,
	// bools, and even int64 all fail the type assertion and yield 0.
	AssertEqual(t, 0, ArgInt(0, "not an int"))
	AssertEqual(t, 0, ArgInt(0, 3.14))
	AssertEqual(t, 0, ArgInt(0, true))
	AssertEqual(t, 0, ArgInt(0, int64(5)))
}

func TestUtils_ArgInt_OutOfBounds_Bad(t *T) {
	// Past the end returns 0; so does the boundary index == len, and the
	// no-args case.
	AssertEqual(t, 0, ArgInt(5, 1, 2))
	AssertEqual(t, 0, ArgInt(2, 1, 2))
	AssertEqual(t, 0, ArgInt(0))
}

// --- ArgBool ---

func TestUtils_ArgBool_Good(t *T) {
	AssertEqual(t, true, ArgBool(0, true, "hello"))
	AssertEqual(t, false, ArgBool(1, true, false))
	// A stored false at an in-bounds index is returned (not confused with a
	// missing arg), and the last position is reachable.
	AssertEqual(t, false, ArgBool(0, false))
	AssertEqual(t, true, ArgBool(2, 1, "x", true))
}

func TestUtils_ArgBool_WrongType_Bad(t *T) {
	// Non-bool values never coerce: a string, an int 1, and a non-bool at the
	// requested index all return false.
	AssertEqual(t, false, ArgBool(0, "not a bool"))
	AssertEqual(t, false, ArgBool(0, 1))
	AssertEqual(t, false, ArgBool(1, true, 0))
}

func TestUtils_ArgBool_OutOfBounds_Bad(t *T) {
	// Past the end returns false; so does the boundary index == len, and the
	// no-args case.
	AssertEqual(t, false, ArgBool(5, true))
	AssertEqual(t, false, ArgBool(1, true))
	AssertEqual(t, false, ArgBool(0))
}

// --- Result.New() ---

func TestUtils_Result_New_SingleArg_Good(t *T) {
	r := Result{}.New("value")
	AssertTrue(t, r.OK)
	AssertEqual(t, "value", r.Value)
}

func TestUtils_Result_New_NilError_Good(t *T) {
	r := Result{}.New("value", nil)
	AssertTrue(t, r.OK)
	AssertEqual(t, "value", r.Value)
}

func TestUtils_Result_New_WithError_Bad(t *T) {
	err := NewError("fail")
	r := Result{}.New("value", err)
	AssertFalse(t, r.OK)
	AssertEqual(t, err, r.Value)
}

func TestUtils_Arg_Good(t *T) {
	r := Arg(0, "agent", 42, true)
	AssertTrue(t, r.OK)
	AssertEqual(t, "agent", r.Value)
}

func TestUtils_Arg_Bad(t *T) {
	r := Arg(4, "agent")
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestUtils_Arg_Ugly(t *T) {
	AssertPanics(t, func() {
		_ = Arg(-1, "agent")
	})
}

func TestUtils_ArgBool_Bad(t *T) {
	// The string "true"/"false" are not bools — no parsing happens — and a
	// nil value is not a bool either.
	AssertFalse(t, ArgBool(0, "true"))
	AssertFalse(t, ArgBool(0, "false"))
	AssertFalse(t, ArgBool(0, nil))
}

func TestUtils_ArgBool_Ugly(t *T) {
	AssertPanics(t, func() {
		_ = ArgBool(-1, true)
	})
}

func TestUtils_ArgInt_Bad(t *T) {
	// A numeric-looking string is not parsed, nil is not an int, and an
	// unsigned uint is a distinct type that does not satisfy `int`.
	AssertEqual(t, 0, ArgInt(0, "8080"))
	AssertEqual(t, 0, ArgInt(0, nil))
	AssertEqual(t, 0, ArgInt(0, uint(5)))
}

func TestUtils_ArgInt_Ugly(t *T) {
	AssertPanics(t, func() {
		_ = ArgInt(-1, 8080)
	})
}

func TestUtils_ArgString_Bad(t *T) {
	// An int is not stringified, nil is not a string, and a []byte (close
	// kin of string) still fails the exact-type assertion.
	AssertEqual(t, "", ArgString(0, 8080))
	AssertEqual(t, "", ArgString(0, nil))
	AssertEqual(t, "", ArgString(0, []byte("hi")))
}

func TestUtils_ArgString_Ugly(t *T) {
	AssertPanics(t, func() {
		_ = ArgString(-1, "agent")
	})
}

func TestUtils_FilterArgs_Bad(t *T) {
	// nil yields nil, and a slice composed entirely of droppable entries
	// (empty strings + "-test." flags) filters down to nil, not []string{}.
	AssertNil(t, FilterArgs(nil))
	AssertNil(t, FilterArgs([]string{"", "-test.v", "-test.run=X"}))
	AssertNil(t, FilterArgs([]string{"", ""}))
}

func TestUtils_FilterArgs_Ugly(t *T) {
	clean := FilterArgs([]string{"", "-test.v", "-test.run=Agent", "deploy"})
	AssertEqual(t, []string{"deploy"}, clean)

	// Only empties and the "-test." prefix are dropped — ordinary flags like
	// "-v" and "--verbose" survive, and original order is preserved.
	mixed := FilterArgs([]string{"-v", "--verbose", "-test.count=1", "build"})
	AssertEqual(t, []string{"-v", "--verbose", "build"}, mixed)
}

func TestUtils_ID_Bad(t *T) {
	// Successive IDs never collide, and the format is "id-<counter>-<random>":
	// an "id-" prefix followed by a counter and a dash-separated random suffix.
	a, b := ID(), ID()
	AssertNotEqual(t, a, b)
	AssertNotEqual(t, b, ID())
	AssertTrue(t, HasPrefix(a, "id-"))
	rest := TrimPrefix(a, "id-")
	AssertNotEmpty(t, rest)
	AssertContains(t, rest, "-")
}

func TestUtils_ID_Ugly(t *T) {
	seen := make(map[string]bool)
	for i := 0; i < 64; i++ {
		id := ID()
		AssertFalse(t, seen[id])
		seen[id] = true
	}
}

func TestUtils_IsFlag_Ugly(t *T) {
	// IsFlag is a pure leading-dash test, so degenerate dash-only and
	// malformed flag-shaped strings all count as flags...
	AssertTrue(t, IsFlag("-"))
	AssertTrue(t, IsFlag("--"))
	AssertTrue(t, IsFlag("---"))
	AssertTrue(t, IsFlag("-=x"))
	// ...while a dash that is not in the leading position does not.
	AssertFalse(t, IsFlag("x-"))
}

func TestUtils_JoinPath_Good(t *T) {
	AssertEqual(t, "agent/dispatch/homelab", JoinPath("agent", "dispatch", "homelab"))
	AssertEqual(t, "a/b", JoinPath("a", "b"))
	AssertEqual(t, "solo", JoinPath("solo"))
	// Already-slashed segments are joined verbatim — JoinPath does not clean.
	AssertEqual(t, "a/b/c/d", JoinPath("a/b", "c/d"))
}

func TestUtils_JoinPath_Bad(t *T) {
	// No segments and a lone empty segment both yield the empty path.
	AssertEqual(t, "", JoinPath())
	AssertEqual(t, "", JoinPath(""))
	AssertEmpty(t, JoinPath())
	// The empty result is specific to empty/no input — a single real
	// segment is returned non-empty and verbatim.
	AssertNotEqual(t, "", JoinPath("a"))
	AssertEqual(t, "a", JoinPath("a"))
}

func TestUtils_JoinPath_Ugly(t *T) {
	// Empty segments are preserved as empty path elements (no cleaning).
	AssertEqual(t, "agent//dispatch", JoinPath("agent", "", "dispatch"))
	AssertEqual(t, "/leading", JoinPath("", "leading"))
	AssertEqual(t, "trailing/", JoinPath("trailing", ""))
}

func TestUtils_ParseFlag_Good(t *T) {
	key, value, ok := ParseFlag("--agent=codex")
	AssertTrue(t, ok)
	AssertEqual(t, "agent", key)
	AssertEqual(t, "codex", value)
}

func TestUtils_ParseFlag_Bad(t *T) {
	key, value, ok := ParseFlag("agent")
	AssertFalse(t, ok)
	AssertEqual(t, "", key)
	AssertEqual(t, "", value)
}

func TestUtils_ParseFlag_Ugly(t *T) {
	key, value, ok := ParseFlag("-")
	AssertFalse(t, ok)
	AssertEqual(t, "", key)
	AssertEqual(t, "", value)
}

func TestUtils_SanitisePath_Bad(t *T) {
	// Empty input and paths whose final element resolves to a dot segment
	// all collapse to the "invalid" sentinel.
	AssertEqual(t, "invalid", SanitisePath(""))
	AssertEqual(t, "invalid", SanitisePath("config/.."))
	AssertEqual(t, "invalid", SanitisePath("tmp/."))
}

func TestUtils_SanitisePath_Ugly(t *T) {
	// Traversal and directory prefixes are discarded — only the final, real
	// path element survives, and a trailing slash is trimmed first.
	AssertEqual(t, "passwd", SanitisePath("../../etc/passwd"))
	AssertEqual(t, "trailing", SanitisePath("trailing/"))
	AssertEqual(t, "z", SanitisePath("x/y/z"))
}

func TestUtils_ValidateName_Bad(t *T) {
	r := ValidateName("")
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "invalid name")
}

func TestUtils_ValidateName_Ugly(t *T) {
	r := ValidateName("agent/dispatch")
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "path separator")
}
