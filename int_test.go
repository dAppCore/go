package core_test

import (
	. "dappco.re/go"
)

// --- Atoi ---

func TestInt_Atoi_Good(t *T) {
	r := Atoi("42")
	AssertTrue(t, r.OK)
	AssertEqual(t, 42, r.Value)
}

func TestInt_Atoi_Bad(t *T) {
	r := Atoi("not-an-int")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

func TestInt_Atoi_Ugly(t *T) {
	r := Atoi("999999999999999999999999999999")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

// --- Itoa ---

func TestInt_Itoa_Good(t *T) {
	AssertEqual(t, "42", Itoa(42))
}

func TestInt_Itoa_Bad(t *T) {
	AssertEqual(t, "-42", Itoa(-42))
}

func TestInt_Itoa_Ugly(t *T) {
	AssertEqual(t, "0", Itoa(0))
}

// --- FormatInt ---

func TestInt_FormatInt_Good(t *T) {
	AssertEqual(t, "ff", FormatInt(255, 16))
}

func TestInt_FormatInt_Bad(t *T) {
	AssertEqual(t, "-ff", FormatInt(-255, 16))
}

func TestInt_FormatInt_Ugly(t *T) {
	AssertEqual(t, "z", FormatInt(35, 36))
	AssertEqual(t, "0", FormatInt(0, 2))
}

// --- FormatUint ---

func TestInt_FormatUint_Good(t *T) {
	AssertEqual(t, "ff", FormatUint(255, 16))
}

func TestInt_FormatUint_Bad(t *T) {
	AssertEqual(t, "10", FormatUint(10, 10))
}

func TestInt_FormatUint_Ugly(t *T) {
	AssertEqual(t, "z", FormatUint(35, 36))
	AssertEqual(t, "0", FormatUint(0, 2))
}

// --- ParseInt ---

func TestInt_ParseInt_Good(t *T) {
	r := ParseInt("ff", 16, 64)
	AssertTrue(t, r.OK)
	AssertEqual(t, int64(255), r.Value)
}

func TestInt_ParseInt_Bad(t *T) {
	r := ParseInt("not-an-int", 10, 64)
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

func TestInt_ParseInt_Ugly(t *T) {
	r := ParseInt("255", 10, 8)
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

// --- AppendInt ---

func TestInt_AppendInt_Good(t *T) {
	AssertEqual(t, "42", AsString(AppendInt(nil, 42, 10)))
}

func TestInt_AppendInt_Bad(t *T) {
	AssertEqual(t, "-ff", AsString(AppendInt(nil, -255, 16)))
}

// Ugly: the point of the Append form is that it EXTENDS dst rather than
// replacing it — a caller assembling a line depends on the prefix surviving.
func TestInt_AppendInt_Ugly(t *T) {
	AssertEqual(t, "n=-7", AsString(AppendInt([]byte("n="), -7, 10)))
	AssertEqual(t, "0", AsString(AppendInt(nil, 0, 2)))
}

// --- AppendUint ---

func TestInt_AppendUint_Good(t *T) {
	AssertEqual(t, "255", AsString(AppendUint(nil, 255, 10)))
}

func TestInt_AppendUint_Bad(t *T) {
	AssertEqual(t, "ff", AsString(AppendUint(nil, 255, 16)))
}

// Ugly: appends to a non-empty dst, and the unsigned maximum — the value an
// int64-shaped formatter would have rendered negative.
func TestInt_AppendUint_Ugly(t *T) {
	AssertEqual(t, "id-1", AsString(AppendUint([]byte("id-"), 1, 10)))
	AssertEqual(t, "18446744073709551615", AsString(AppendUint(nil, ^uint64(0), 10)))
}

// --- AppendBool ---

func TestInt_AppendBool_Good(t *T) {
	AssertEqual(t, "true", AsString(AppendBool(nil, true)))
}

func TestInt_AppendBool_Bad(t *T) {
	AssertEqual(t, "false", AsString(AppendBool(nil, false)))
}

func TestInt_AppendBool_Ugly(t *T) {
	AssertEqual(t, "ok=true", AsString(AppendBool([]byte("ok="), true)))
}

// --- AppendFloat ---

func TestInt_AppendFloat_Good(t *T) {
	AssertEqual(t, "1.5", AsString(AppendFloat(nil, 1.5, 'g', -1, 64)))
}

// Bad: the two values a float field can hold that are not numbers. Built by
// division rather than a math helper because a test here may only import
// dappco.re/go, and core exposes no NaN/Inf constructor.
func TestInt_AppendFloat_Bad(t *T) {
	zero := 0.0
	AssertEqual(t, "NaN", AsString(AppendFloat(nil, zero/zero, 'g', -1, 64)))
	AssertEqual(t, "+Inf", AsString(AppendFloat(nil, 1.0/zero, 'g', -1, 64)))
}

// Ugly: 'g' drops the exponent for small magnitudes but not large ones, and a
// 32-bit bitSize rounds to the nearest float32 — both surprising but correct.
func TestInt_AppendFloat_Ugly(t *T) {
	AssertEqual(t, "1e+21", AsString(AppendFloat(nil, 1e21, 'g', -1, 64)))
	AssertEqual(t, "0.1", AsString(AppendFloat(nil, float64(float32(0.1)), 'g', -1, 32)))
}

// --- AppendQuote ---

func TestInt_AppendQuote_Good(t *T) {
	AssertEqual(t, `"hello"`, AsString(AppendQuote(nil, "hello")))
}

func TestInt_AppendQuote_Bad(t *T) {
	AssertEqual(t, `""`, AsString(AppendQuote(nil, "")))
}

// Ugly: the reason the logger quotes at all — an embedded quote, backslash or
// newline has to come back escaped or the line it is written into is corrupt.
func TestInt_AppendQuote_Ugly(t *T) {
	AssertEqual(t, `"a\"b\\c\n"`, AsString(AppendQuote(nil, "a\"b\\c\n")))
}
