// SPDX-License-Identifier: EUPL-1.2

package core

func TestAssert_assertCmpFloat64_Good(t *T) {
	AssertEqual(t, -1, assertCmpFloat64(1.25, 2.5))
	AssertEqual(t, -1, assertCmpFloat64(-2.5, -1.25)) // negatives: -2.5 < -1.25
	AssertEqual(t, -1, assertCmpFloat64(0, 1e-12))    // tiny positive difference
	AssertEqual(t, 0, assertCmpFloat64(NaN(), 2.5))   // NaN unordered: neither < nor > -> 0
}
func TestAssert_assertCmpFloat64_Bad(t *T) {
	AssertEqual(t, 0, assertCmpFloat64(2.5, 2.5))
	AssertEqual(t, 0, assertCmpFloat64(0, 0))
	AssertEqual(t, 0, assertCmpFloat64(-1.5, -1.5))
	AssertEqual(t, 0, assertCmpFloat64(NaN(), NaN())) // NaN != NaN, neither branch fires -> 0
}
func TestAssert_assertCmpFloat64_Ugly(t *T) {
	AssertEqual(t, 1, assertCmpFloat64(-0.5, -1.5))
	AssertEqual(t, 1, assertCmpFloat64(2.5, 1.25))
	AssertEqual(t, 1, assertCmpFloat64(1e-12, 0)) // tiny positive vs zero -> greater
}
func TestAssert_assertCmpInt64_Good(t *T) {
	AssertEqual(t, -1, assertCmpInt64(-1, 1))
	AssertEqual(t, -1, assertCmpInt64(0, 1))            // adjacent values
	AssertEqual(t, -1, assertCmpInt64(-1<<62, 1<<62))   // wide negative-to-positive span
}
func TestAssert_assertCmpInt64_Bad(t *T) {
	AssertEqual(t, 0, assertCmpInt64(42, 42))
	AssertEqual(t, 0, assertCmpInt64(0, 0))
	AssertEqual(t, 0, assertCmpInt64(-1<<62, -1<<62)) // equal large negatives
}
func TestAssert_assertCmpInt64_Ugly(t *T) {
	AssertEqual(t, 1, assertCmpInt64(1<<62, -1<<62))
	AssertEqual(t, 1, assertCmpInt64(1, 0))  // adjacent values
	AssertEqual(t, 1, assertCmpInt64(0, -1)) // zero greater than negative
}
func TestAssert_assertCmpUint64_Good(t *T) {
	AssertEqual(t, -1, assertCmpUint64(1, 2))
	AssertEqual(t, -1, assertCmpUint64(0, 1))      // zero less than one
	AssertEqual(t, -1, assertCmpUint64(1, 1<<63))  // small vs high bit set
}
func TestAssert_assertCmpUint64_Bad(t *T) {
	AssertEqual(t, 0, assertCmpUint64(42, 42))
	AssertEqual(t, 0, assertCmpUint64(0, 0))
	AssertEqual(t, 0, assertCmpUint64(1<<63, 1<<63)) // equal high values
}
func TestAssert_assertCmpUint64_Ugly(t *T) {
	AssertEqual(t, 1, assertCmpUint64(1<<63, 1))
	AssertEqual(t, 1, assertCmpUint64(1, 0))          // one greater than zero
	AssertEqual(t, 1, assertCmpUint64(^uint64(0), 0)) // max uint64 greater than zero
}
func TestAssert_assertCompare_Good(t *T) {
	cmp, ok := assertCompare("agent", "brain")
	AssertTrue(t, ok)
	AssertEqual(t, -1, cmp) // string ordering

	// cross-kind numerics compare by value
	cmp, ok = assertCompare(int64(4), uint64(3))
	AssertTrue(t, ok)
	AssertEqual(t, 1, cmp)

	cmp, ok = assertCompare(uint(4), int(5))
	AssertTrue(t, ok)
	AssertEqual(t, -1, cmp)

	cmp, ok = assertCompare(uint(4), uint64(4))
	AssertTrue(t, ok)
	AssertEqual(t, 0, cmp)
}
func TestAssert_assertCompare_Bad(t *T) {
	cmp, ok := assertCompare(struct{ Name string }{"agent"}, struct{ Name string }{"agent"})

	AssertFalse(t, ok)
	AssertEqual(t, 0, cmp)
}
func TestAssert_assertCompare_Ugly(t *T) {
	cmp, ok := assertCompare(-1, uint(1)) // signed vs unsigned
	AssertTrue(t, ok)
	AssertEqual(t, -1, cmp)

	cmp, ok = assertCompare(uint(4), 4.5) // uint vs float
	AssertTrue(t, ok)
	AssertEqual(t, -1, cmp)

	cmp, ok = assertCompare(float64(6), uint(5)) // float vs uint
	AssertTrue(t, ok)
	AssertEqual(t, 1, cmp)

	cmp, ok = assertCompare(float32(3.5), float64(3.5)) // float32 vs float64
	AssertTrue(t, ok)
	AssertEqual(t, 0, cmp)

	cmp, ok = assertCompare("brain", "agent") // string reverse ordering
	AssertTrue(t, ok)
	AssertEqual(t, 1, cmp)
}
func TestAssert_assertContains_Good(t *T) {
	AssertTrue(t, assertContains([]string{"agent", "dispatch"}, "dispatch"))
	AssertTrue(t, assertContains([]string{"agent", "dispatch"}, "agent"))    // first element
	AssertFalse(t, assertContains([]string{"agent", "dispatch"}, "missing")) // absent element
	AssertTrue(t, assertContains([]int{1, 2, 3}, 2))                         // non-string slice via deep-equal
	AssertFalse(t, assertContains([]int{1, 2, 3}, "2"))                      // type mismatch -> deep-equal false
}
func TestAssert_assertContains_Bad(t *T) {
	AssertFalse(t, assertContains("agent dispatch", "missing"))
	AssertTrue(t, assertContains("agent dispatch", "dispatch")) // substring present
	AssertTrue(t, assertContains("agent dispatch", ""))         // empty substring always contained
	AssertFalse(t, assertContains([]string{"a"}, "missing"))    // slice miss
	AssertFalse(t, assertContains("agent", 42))                 // non-string needle vs string haystack
	AssertFalse(t, assertContains("agent", []byte("a")))        // slice-kind needle is not a string
}
func TestAssert_assertContains_Ugly(t *T) {
	AssertTrue(t, assertContains(map[string]int{"session": 1}, "session"))
	AssertFalse(t, assertContains(map[string]int{"session": 1}, "missing")) // absent key
	AssertTrue(t, assertContains(map[int]string{7: "x"}, 7))                 // int key membership
	AssertFalse(t, assertContains(42, "x"))                                  // unsupported kind -> false
}
func TestAssert_assertIsEmpty_Good(t *T) {
	AssertTrue(t, assertIsEmpty(""))
	AssertTrue(t, assertIsEmpty(nil))              // nil is empty
	AssertTrue(t, assertIsEmpty([]int{}))          // zero-length slice
	AssertTrue(t, assertIsEmpty(map[string]int{})) // zero-length map
	AssertTrue(t, assertIsEmpty(0))                // zero int via zero-value compare
	AssertFalse(t, assertIsEmpty("agent"))         // non-empty string contrast
}
func TestAssert_assertIsEmpty_Bad(t *T) {
	AssertFalse(t, assertIsEmpty("agent"))
	AssertFalse(t, assertIsEmpty([]int{1}))               // non-empty slice
	AssertFalse(t, assertIsEmpty(map[string]int{"k": 1})) // non-empty map
	AssertFalse(t, assertIsEmpty(42))                     // non-zero int
	AssertTrue(t, assertIsEmpty(0))                       // zero value is empty (contrast)
}
func TestAssert_assertIsEmpty_Ugly(t *T) {
	names := []string{}

	AssertTrue(t, assertIsEmpty(&names)) // pointer recurses to empty slice

	full := []string{"agent"}

	AssertFalse(t, assertIsEmpty(&full)) // pointer recurses to non-empty slice

	var typedNil *[]string

	AssertTrue(t, assertIsEmpty(typedNil)) // typed-nil pointer is empty

	var m map[string]int

	AssertTrue(t, assertIsEmpty(m)) // nil map is empty
}
func TestAssert_assertIsNil_Good(t *T) {
	AssertTrue(t, assertIsNil(nil))

	var p *int

	AssertTrue(t, assertIsNil(p)) // typed-nil pointer inside any

	var s []int

	AssertTrue(t, assertIsNil(s)) // nil slice

	var m map[string]int

	AssertTrue(t, assertIsNil(m)) // nil map
}
func TestAssert_assertIsNil_Bad(t *T) {
	AssertFalse(t, assertIsNil(0))
	AssertFalse(t, assertIsNil(""))      // empty string is not nil
	AssertFalse(t, assertIsNil([]int{})) // non-nil empty slice is not nil

	x := 42

	AssertFalse(t, assertIsNil(&x)) // non-nil pointer
}
func TestAssert_assertIsNil_Ugly(t *T) {
	var sessions map[string]string

	AssertTrue(t, assertIsNil(sessions)) // nil map

	sessions = map[string]string{}

	AssertFalse(t, assertIsNil(sessions)) // non-nil empty map is not nil

	var fn func()

	AssertTrue(t, assertIsNil(fn)) // nil func
}
func TestAssert_assertMsg_Good(t *T) {
	AssertEqual(t, " — agent retry", assertMsg([]string{"agent", "retry"}))
	AssertEqual(t, " — solo", assertMsg([]string{"solo"}))    // single element: no join separator
	AssertEqual(t, " — a b c", assertMsg([]string{"a", "b", "c"})) // three elements space-joined
}
func TestAssert_assertMsg_Bad(t *T) {
	AssertEqual(t, "", assertMsg(nil))
	AssertEqual(t, "", assertMsg([]string{})) // empty slice also yields ""
}
func TestAssert_assertMsg_Ugly(t *T) {
	AssertEqual(t, " — lethean degraded", assertMsg([]string{"lethean", "degraded"}))
	AssertEqual(t, " — ", assertMsg([]string{""})) // single empty string: len!=0 so prefix added, joins to ""
}

// --- assertFail: the centralised emitter for Assert* / Require* ---

// stubT records Errorf/Fatal calls without aborting the test, so the
// triplet can verify assertFail's two output formats.
type stubT struct {
	*T
	msgs  []string
	fatal bool
}

func (s *stubT) Helper() { /* no-op helper for testing.TB compatibility */ }
func (s *stubT) Error(args ...any) {
	s.msgs = append(s.msgs, Sprint(args...))
}
func (s *stubT) Errorf(format string, args ...any) {
	s.msgs = append(s.msgs, Sprintf(format, args...))
}
func (s *stubT) Fatal(args ...any) {
	s.fatal = true
	s.msgs = append(s.msgs, Sprint(args...))
}
func (s *stubT) Fatalf(format string, args ...any) {
	s.fatal = true
	s.msgs = append(s.msgs, Sprintf(format, args...))
}

func TestAssert_assertFail_Good(t *T) {
	prev := AssertVerbose
	defer func() { AssertVerbose = prev }()
	AssertVerbose = false

	st := &stubT{T: t}
	assertFail(st, false, "TestAssertion", nil, "want", 1, "got", 2)

	AssertEqual(t, false, st.fatal)
	AssertLen(t, st.msgs, 1)
	AssertContains(t, st.msgs[0], "TestAssertion")
	AssertContains(t, st.msgs[0], "want=1")
	AssertContains(t, st.msgs[0], "got=2")
}

func TestAssert_assertFail_Bad(t *T) {
	prev := AssertVerbose
	defer func() { AssertVerbose = prev }()
	AssertVerbose = false

	st := &stubT{T: t}
	assertFail(st, true, "RequireFatal", []string{"agent", "context"}, "want", "non-nil", "got", nil)

	AssertEqual(t, true, st.fatal)
	AssertContains(t, st.msgs[0], "RequireFatal")
	AssertContains(t, st.msgs[0], "agent context")
}

func TestAssert_assertFail_Ugly(t *T) {
	prev := AssertVerbose
	defer func() { AssertVerbose = prev }()
	AssertVerbose = true

	st := &stubT{T: t}
	assertFail(st, false, "VerboseAssertion", []string{"homelab", "drained"}, "want", 42, "got", 13)

	AssertEqual(t, false, st.fatal)
	AssertContains(t, st.msgs[0], "VerboseAssertion failed")
	AssertContains(t, st.msgs[0], "want: 42")
	AssertContains(t, st.msgs[0], "got: 13")
	AssertContains(t, st.msgs[0], "msg:")

	// verbose + fatal: the Require* path under verbose formatting
	st = &stubT{T: t}
	assertFail(st, true, "VerboseRequire", []string{"agent", "halted"}, "want", true, "got", false)
	AssertTrue(t, st.fatal)
	AssertContains(t, st.msgs[0], "VerboseRequire failed")
	AssertContains(t, st.msgs[0], "agent halted")
}

func assertStub(t *T) *stubT {
	t.Helper()
	return &stubT{T: t}
}

func assertOneMessage(t *T, st *stubT, want string) {
	t.Helper()
	AssertLen(t, st.msgs, 1)
	AssertContains(t, st.msgs[0], want)
}

func TestAssert_AssertEqual_Bad(t *T) {
	st := assertStub(t)
	AssertEqual(st, "expected", "actual")
	assertOneMessage(t, st, "AssertEqual")
}

func TestAssert_AssertNotEqual_Bad(t *T) {
	st := assertStub(t)
	AssertNotEqual(st, "same", "same")
	assertOneMessage(t, st, "AssertNotEqual")
}

func TestAssert_AssertTrue_Bad(t *T) {
	st := assertStub(t)
	AssertTrue(st, false)
	assertOneMessage(t, st, "AssertTrue")
}

func TestAssert_AssertFalse_Bad(t *T) {
	st := assertStub(t)
	AssertFalse(st, true)
	assertOneMessage(t, st, "AssertFalse")
}

func TestAssert_AssertNil_Bad(t *T) {
	st := assertStub(t)
	AssertNil(st, 42)
	assertOneMessage(t, st, "AssertNil")
}

func TestAssert_AssertNotNil_Bad(t *T) {
	st := assertStub(t)
	AssertNotNil(st, nil)
	assertOneMessage(t, st, "AssertNotNil")
}

func TestAssert_AssertNoError_Bad(t *T) {
	st := assertStub(t)
	AssertNoError(st, AnError)
	assertOneMessage(t, st, "AssertNoError")
}

func TestAssert_AssertError_Bad(t *T) {
	st := assertStub(t)
	AssertError(st, nil)
	assertOneMessage(t, st, "AssertError")
}

func TestAssert_AssertError_Ugly(t *T) {
	// Error present, but the required substring is missing.
	st := assertStub(t)
	AssertError(st, AnError, "missing")
	assertOneMessage(t, st, "want-substring")
}

func TestAssert_AssertContains_Bad(t *T) {
	st := assertStub(t)
	AssertContains(st, "agent", "missing")
	assertOneMessage(t, st, "AssertContains")
}

func TestAssert_AssertNotContains_Bad(t *T) {
	st := assertStub(t)
	AssertNotContains(st, "agent", "gen")
	assertOneMessage(t, st, "AssertNotContains")
}

func TestAssert_AssertLen_Bad(t *T) {
	st := assertStub(t)
	AssertLen(st, []string{"agent"}, 2)
	assertOneMessage(t, st, "AssertLen")
}

func TestAssert_AssertLen_Ugly(t *T) {
	// A kind that has no length.
	st := assertStub(t)
	AssertLen(st, 42, 1)
	assertOneMessage(t, st, "unsupported kind")
}

func TestAssert_AssertEmpty_Bad(t *T) {
	st := assertStub(t)
	AssertEmpty(st, "agent")
	assertOneMessage(t, st, "AssertEmpty")
}

func TestAssert_AssertNotEmpty_Bad(t *T) {
	st := assertStub(t)
	AssertNotEmpty(st, "")
	assertOneMessage(t, st, "AssertNotEmpty")
}

func TestAssert_AssertGreater_Bad(t *T) {
	st := assertStub(t)
	AssertGreater(st, 1, 2)
	assertOneMessage(t, st, "AssertGreater")
}

func TestAssert_AssertGreater_Ugly(t *T) {
	// Incomparable operands.
	st := assertStub(t)
	AssertGreater(st, struct{}{}, struct{}{})
	assertOneMessage(t, st, "incomparable got")
}

func TestAssert_AssertGreaterOrEqual_Bad(t *T) {
	st := assertStub(t)
	AssertGreaterOrEqual(st, 1, 2)
	assertOneMessage(t, st, "want>=")
}

func TestAssert_AssertGreaterOrEqual_Ugly(t *T) {
	// Incomparable operands.
	st := assertStub(t)
	AssertGreaterOrEqual(st, struct{}{}, struct{}{})
	assertOneMessage(t, st, "AssertGreaterOrEqual")
}

func TestAssert_AssertLess_Bad(t *T) {
	st := assertStub(t)
	AssertLess(st, 2, 1)
	assertOneMessage(t, st, "want<")
}

func TestAssert_AssertLess_Ugly(t *T) {
	// Incomparable operands.
	st := assertStub(t)
	AssertLess(st, struct{}{}, struct{}{})
	assertOneMessage(t, st, "AssertLess")
}

func TestAssert_AssertLessOrEqual_Bad(t *T) {
	st := assertStub(t)
	AssertLessOrEqual(st, 2, 1)
	assertOneMessage(t, st, "want<=")
}

func TestAssert_AssertLessOrEqual_Ugly(t *T) {
	// Incomparable operands.
	st := assertStub(t)
	AssertLessOrEqual(st, struct{}{}, struct{}{})
	assertOneMessage(t, st, "AssertLessOrEqual")
}

func TestAssert_AssertPanics_Bad(t *T) {
	st := assertStub(t)
	AssertPanics(st, func() {})
	assertOneMessage(t, st, "normal-return")
}

func TestAssert_AssertNotPanics_Bad(t *T) {
	st := assertStub(t)
	AssertNotPanics(st, func() { panic("agent failed") })
	assertOneMessage(t, st, "got panic")
}

func TestAssert_AssertPanicsWithError_Bad(t *T) {
	st := assertStub(t)
	AssertPanicsWithError(st, "agent failed", func() {})
	assertOneMessage(t, st, "normal-return")
}

func TestAssert_AssertPanicsWithError_Ugly(t *T) {
	// Panics, but with the wrong error text.
	st := assertStub(t)
	AssertPanicsWithError(st, "agent failed", func() { panic("other failure") })
	assertOneMessage(t, st, "want-substring")
}

func TestAssert_AssertErrorIs_Bad(t *T) {
	st := assertStub(t)
	AssertErrorIs(st, AnError, NewError("different"))
	assertOneMessage(t, st, "AssertErrorIs")
}

func TestAssert_AssertInDelta_Bad(t *T) {
	st := assertStub(t)
	AssertInDelta(st, 1, 2, 0.1)
	assertOneMessage(t, st, "actual-diff")
}

func TestAssert_AssertInDelta_Ugly(t *T) {
	// A NaN operand can't be compared against a delta.
	st := assertStub(t)
	AssertInDelta(st, NaN(), 1, 0.1)
	assertOneMessage(t, st, "NaN involved want")
}

func TestAssert_AssertSame_Bad(t *T) {
	left := 1
	right := 1
	st := assertStub(t)
	AssertSame(st, &left, &right)
	assertOneMessage(t, st, "AssertSame")
}

func TestAssert_AssertSame_Ugly(t *T) {
	// Non-pointer operands.
	st := assertStub(t)
	AssertSame(st, 1, 1)
	assertOneMessage(t, st, "both args must be pointers")
}

func TestAssert_AssertElementsMatch_Bad(t *T) {
	st := assertStub(t)
	AssertElementsMatch(st, []int{1}, []int{1, 2})
	assertOneMessage(t, st, "len-mismatch")

	st = assertStub(t)
	AssertElementsMatch(st, []int{1, 3}, []int{1, 2})
	assertOneMessage(t, st, "missing element")
}

func TestAssert_AssertElementsMatch_Ugly(t *T) {
	// Non-slice operands.
	st := assertStub(t)
	AssertElementsMatch(st, 1, []int{1})
	assertOneMessage(t, st, "both args must be slices")
}

func TestAssert_RequireNoError_Bad(t *T) {
	st := assertStub(t)
	RequireNoError(st, AnError)
	assertOneMessage(t, st, "RequireNoError")
	AssertTrue(t, st.fatal)
}

func TestAssert_RequireTrue_Bad(t *T) {
	st := assertStub(t)
	RequireTrue(st, false)
	assertOneMessage(t, st, "RequireTrue")
	AssertTrue(t, st.fatal)
}

func TestAssert_RequireNotEmpty_Bad(t *T) {
	st := assertStub(t)
	RequireNotEmpty(st, "")
	assertOneMessage(t, st, "RequireNotEmpty")
	AssertTrue(t, st.fatal)
}

// assertCompareMixedKinds cases folded into TestAssert_assertCompare_{Good,Ugly}.
