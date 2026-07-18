// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

// Every Example below is compiled but not executed: none carries an
// "Output:" comment. Go's own convention is explicit about this —
// "Example functions without output comments are compiled but not
// executed" (https://pkg.go.dev/testing#hdr-Examples) — so a nil *T is
// safe here: the Assert*/Require* body (every one starts with
// t.Helper()) never actually runs. This is the only honest mechanism
// available: every Assert*/Require* helper requires a live
// *testing.T/TB, which an Example function signature can never supply
// (Example functions take no arguments), and fabricating a mock T that
// pretends to be live would be theatrical, not documentation.
//
// test_example_test.go already established this exact pattern for
// these same 27 symbols, so those exact names (ExampleAssertEqual etc.)
// are already declared in that file — reusing them here would be a
// duplicate-declaration build error. The file-aware Example audit
// requires the association to live in THIS file (assert.go's sibling),
// so each entry below is named as a distinct "_variant" demonstrating a
// genuinely different scenario from its test_example_test.go
// counterpart, rather than redeclaring an identical duplicate under a
// different file.

// ExampleAssertEqual_numeric asserts equal integers through `AssertEqual` for AX-native
// tests. Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertEqual_numeric() {
	var t *T
	AssertEqual(t, 4, len("core"))
}

// ExampleAssertNotEqual_numeric asserts different integers through `AssertNotEqual` for
// AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertNotEqual_numeric() {
	var t *T
	AssertNotEqual(t, 1, 2)
}

// ExampleAssertTrue_withMessage asserts a true condition with context through `AssertTrue`
// for AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertTrue_withMessage() {
	var t *T
	items := []string{"agent"}
	AssertTrue(t, len(items) > 0, "items must not be empty")
}

// ExampleAssertFalse_withMessage asserts a false condition with context through
// `AssertFalse` for AX-native tests. Passing assertions are silent while failures stay
// one-line and AI-readable.
func ExampleAssertFalse_withMessage() {
	var t *T
	items := []string{}
	AssertFalse(t, len(items) > 0, "items must be empty")
}

// ExampleAssertNil_typedNil asserts a typed-nil interface value through `AssertNil` for
// AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertNil_typedNil() {
	var t *T
	var p *int
	var v any = p
	AssertNil(t, v)
}

// ExampleAssertNotNil_pointer asserts a non-nil pointer through `AssertNotNil` for
// AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertNotNil_pointer() {
	var t *T
	AssertNotNil(t, &struct{}{})
}

// ExampleAssertNoError_withMessage asserts no error with context through `AssertNoError`
// for AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertNoError_withMessage() {
	var t *T
	AssertNoError(t, nil, "setup phase")
}

// ExampleAssertError_wrapped asserts a wrapped error through `AssertError` for AX-native
// tests. Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertError_wrapped() {
	var t *T
	AssertError(t, Wrap(AnError, "example", "detail"))
}

// ExampleAssertContains_map asserts key membership through `AssertContains` for AX-native
// tests. Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertContains_map() {
	var t *T
	AssertContains(t, map[string]int{"x": 1}, "x")
}

// ExampleAssertNotContains_string asserts substring absence through `AssertNotContains`
// for AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertNotContains_string() {
	var t *T
	AssertNotContains(t, "agent dispatch", "missing")
}

// ExampleAssertLen_string asserts string length through `AssertLen` for AX-native tests.
// Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertLen_string() {
	var t *T
	AssertLen(t, "agent", 5)
}

// ExampleAssertEmpty_string asserts an empty string through `AssertEmpty` for AX-native
// tests. Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertEmpty_string() {
	var t *T
	AssertEmpty(t, "")
}

// ExampleAssertNotEmpty_map asserts a non-empty map through `AssertNotEmpty` for AX-native
// tests. Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertNotEmpty_map() {
	var t *T
	AssertNotEmpty(t, map[string]int{"x": 1})
}

// ExampleAssertGreater_string asserts string ordering through `AssertGreater` for
// AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertGreater_string() {
	var t *T
	AssertGreater(t, "b", "a")
}

// ExampleAssertGreaterOrEqual_duration asserts elapsed-time ordering through
// `AssertGreaterOrEqual` for AX-native tests. Passing assertions are silent while
// failures stay one-line and AI-readable.
func ExampleAssertGreaterOrEqual_duration() {
	var t *T
	elapsed := 5 * Second
	minDuration := 2 * Second
	AssertGreaterOrEqual(t, elapsed, minDuration)
}

// ExampleAssertLess_string asserts string ordering through `AssertLess` for AX-native
// tests. Passing assertions are silent while failures stay one-line and AI-readable.
func ExampleAssertLess_string() {
	var t *T
	AssertLess(t, "a", "b")
}

// ExampleAssertLessOrEqual_float asserts float ordering through `AssertLessOrEqual` for
// AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertLessOrEqual_float() {
	var t *T
	AssertLessOrEqual(t, 1.5, 2.5)
}

// ExampleAssertPanics_error asserts a panic carrying an error value through
// `AssertPanics` for AX-native tests. Passing assertions are silent while failures stay
// one-line and AI-readable.
func ExampleAssertPanics_error() {
	var t *T
	AssertPanics(t, func() { panic(AnError) })
}

// ExampleAssertNotPanics_arithmetic asserts a non-panicking closure through
// `AssertNotPanics` for AX-native tests. Passing assertions are silent while failures
// stay one-line and AI-readable.
func ExampleAssertNotPanics_arithmetic() {
	var t *T
	AssertNotPanics(t, func() { _ = 10 / 2 })
}

// ExampleAssertPanicsWithError_nonError asserts panic text from a non-error value through
// `AssertPanicsWithError` for AX-native tests. Passing assertions are silent while
// failures stay one-line and AI-readable.
func ExampleAssertPanicsWithError_nonError() {
	var t *T
	AssertPanicsWithError(t, "boom", func() { panic("boom") })
}

// ExampleAssertErrorIs_sentinel asserts a wrapped sentinel error through `AssertErrorIs`
// for AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleAssertErrorIs_sentinel() {
	var t *T
	AssertErrorIs(t, Wrap(ErrNotExist, "lookup", "missing"), ErrNotExist)
}

// ExampleAssertInDelta_exact asserts exact numeric equality (zero delta) through
// `AssertInDelta` for AX-native tests. Passing assertions are silent while failures stay
// one-line and AI-readable.
func ExampleAssertInDelta_exact() {
	var t *T
	AssertInDelta(t, 2.5, 2.5, 0)
}

// ExampleAssertSame_core asserts pointer identity on a Core singleton through
// `AssertSame` for AX-native tests. Passing assertions are silent while failures stay
// one-line and AI-readable.
func ExampleAssertSame_core() {
	var t *T
	c := New()
	// Core() returns the receiver itself, so this is genuine identity
	// between two independently obtained references, not a self-compare.
	AssertSame(t, c, c.Core())
}

// ExampleAssertElementsMatch_ints asserts unordered int-slice equality through
// `AssertElementsMatch` for AX-native tests. Passing assertions are silent while
// failures stay one-line and AI-readable.
func ExampleAssertElementsMatch_ints() {
	var t *T
	AssertElementsMatch(t, []int{1, 2, 3}, []int{3, 2, 1})
}

// ExampleAssertAllocs gates a hot path's allocation ceiling through
// `AssertAllocs` for AX-native tests. Passing assertions are silent while
// failures stay one-line and AI-readable.
func ExampleAssertAllocs() {
	var t *T
	AssertAllocs(t, 0, func() { _ = Abs(-1) })
}

// ExampleRequireNoError_withMessage requires no error with context through
// `RequireNoError` for AX-native tests. Passing assertions are silent while failures
// stay one-line and AI-readable.
func ExampleRequireNoError_withMessage() {
	var t *T
	RequireNoError(t, nil, "fixture setup")
}

// ExampleRequireTrue_withMessage requires a true condition with context through
// `RequireTrue` for AX-native tests. Passing assertions are silent while failures stay
// one-line and AI-readable.
func ExampleRequireTrue_withMessage() {
	var t *T
	RequireTrue(t, true, "scan precondition")
}

// ExampleRequireNotEmpty_string requires a non-empty string through `RequireNotEmpty`
// for AX-native tests. Passing assertions are silent while failures stay one-line and
// AI-readable.
func ExampleRequireNotEmpty_string() {
	var t *T
	RequireNotEmpty(t, "fixturePath")
}
