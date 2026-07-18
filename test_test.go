package core_test

import (
	. "dappco.re/go"
)

func TestTest_AssertCLI_Bad(t *T) {
	c := New()
	AssertCLI(t, c, CLITest{
		Name:   "process capability refused",
		Cmd:    "agent-dispatch",
		WantOK: false,
	})
}

func TestTest_AssertCLI_Ugly(t *T) {
	c := New()
	var captured Options
	c.Action("process.run", func(_ Context, opts Options) Result {
		captured = opts
		return Result{Value: "homelab health check ok\n", OK: true}
	})

	dir := t.TempDir()
	AssertCLI(t, c, CLITest{
		Name:     "env and directory dispatch",
		Cmd:      "agent-dispatch",
		Args:     []string{"health", "--site=homelab"},
		Dir:      dir,
		Env:      []string{"CORE_AGENT=codex"},
		WantOK:   true,
		Contains: "health check ok",
	})

	AssertEqual(t, dir, captured.String("dir"))
	env, ok := captured.Get("env").Value.([]string)
	AssertTrue(t, ok)
	AssertContains(t, env, "CORE_AGENT=codex")
}

func TestTest_AssertCLIs_Bad(t *T) {
	AssertCLIs(t, New(), []CLITest{
		{Name: "missing process runner", Cmd: "agent-dispatch", WantOK: false},
	})
}

func TestTest_AssertCLIs_Ugly(t *T) {
	// nil and empty case slices are both no-op: zero sub-tests, zero failures.
	AssertCLIs(t, New(), nil)
	AssertCLIs(t, New(), []CLITest{})
	// Cases with no Name fall back to Cmd for the sub-test label; with no
	// process.run service registered the dispatch fails, so WantOK:false holds.
	AssertCLIs(t, New(), []CLITest{
		{Cmd: "agent-dispatch", WantOK: false},
		{Cmd: "agent-deploy", Args: []string{"--site=homelab"}, WantOK: false},
	})
}

func TestTest_AssertContains_Good(t *T) {
	// String haystack: substring match anywhere in the string.
	AssertContains(t, "agent dispatch completed", "dispatch")
	AssertContains(t, "agent dispatch completed", "agent")
	AssertContains(t, "agent dispatch completed", "completed")
	AssertContains(t, "agent dispatch completed", "agent dispatch completed")
	// Absent substring is the false branch of the predicate.
	AssertNotContains(t, "agent dispatch completed", "rollback")
	// A non-string needle against a string haystack never matches.
	AssertNotContains(t, "agent dispatch completed", 5)
}

func TestTest_AssertContains_Bad(t *T) {
	// Map haystack: membership is tested against the KEYS, not the values.
	AssertContains(t, map[string]int{"session": 1}, "session")
	AssertContains(t, map[string]int{"session": 1, "agent": 2}, "agent")
	AssertNotContains(t, map[string]int{"session": 1}, "missing")
	// The map value (1) is not a key, so it must not be reported as present.
	AssertNotContains(t, map[string]int{"session": 1}, 1)
	// Non-string key types are matched by deep equality.
	AssertContains(t, map[int]string{7: "x"}, 7)
	AssertNotContains(t, map[int]string{7: "x"}, 8)
}

func TestTest_AssertContains_Ugly(t *T) {
	// strings.Contains("", "") is true: every string contains the empty string.
	AssertContains(t, "", "")
	AssertContains(t, "anything", "")
	// The empty string contains nothing non-empty.
	AssertNotContains(t, "", "x")
	// Slice membership finds an empty-string element.
	AssertContains(t, []string{"", "a"}, "")
	AssertNotContains(t, []string{"a", "b"}, "")
}

func TestTest_AssertElementsMatch_Good(t *T) {
	AssertElementsMatch(t, []string{"agent", "health", "deploy"}, []string{"deploy", "agent", "health"})
	// Identical order still matches.
	AssertElementsMatch(t, []int{1, 2, 3}, []int{1, 2, 3})
	// Fully reversed.
	AssertElementsMatch(t, []int{1, 2, 3}, []int{3, 2, 1})
	// Single element.
	AssertElementsMatch(t, []string{"solo"}, []string{"solo"})
}

func TestTest_AssertElementsMatch_Bad(t *T) {
	// Duplicate elements must match by count, not just by set membership.
	AssertElementsMatch(t, []int{1, 1, 2}, []int{1, 2, 1})
	AssertElementsMatch(t, []string{"a", "a", "a"}, []string{"a", "a", "a"})
	AssertElementsMatch(t, []string{"a", "a", "b"}, []string{"b", "a", "a"})
}

func TestTest_AssertElementsMatch_Ugly(t *T) {
	// Two empty slices match trivially.
	AssertElementsMatch(t, []string{}, []string{})
	// A nil slice and an empty slice are both length-0, so they match.
	AssertElementsMatch(t, []int(nil), []int{})
	// Arrays are supported alongside slices.
	AssertElementsMatch(t, [2]int{1, 2}, [2]int{2, 1})
}

func TestTest_AssertEmpty_Good(t *T) {
	AssertEmpty(t, []string{})
	AssertEmpty(t, "")
	AssertEmpty(t, map[string]int{})
	AssertEmpty(t, make([]int, 0))
	// Inverse: a populated slice is not empty.
	AssertNotEmpty(t, []string{"agent"})
}

func TestTest_AssertEmpty_Bad(t *T) {
	// Untyped nil is empty.
	AssertEmpty(t, nil)
	// Typed nil collections report empty without panicking.
	var s []string
	AssertEmpty(t, s)
	var m map[string]int
	AssertEmpty(t, m)
	var p *int
	AssertEmpty(t, p)
}

func TestTest_AssertEmpty_Ugly(t *T) {
	// Non-container kinds are empty when equal to their zero value.
	AssertEmpty(t, 0)
	AssertEmpty(t, false)
	AssertEmpty(t, 0.0)
	AssertEmpty(t, struct{ Agent string }{})
	// Inverse: non-zero scalars are not empty.
	AssertNotEmpty(t, 1)
	AssertNotEmpty(t, struct{ Agent string }{Agent: "codex"})
}

func TestTest_AssertEqual_Good(t *T) {
	AssertEqual(t, "agent dispatch", "agent dispatch")
	// Computed values compare equal.
	AssertEqual(t, 6, 2*3)
	// Deep equality across slices and order-independent maps.
	AssertEqual(t, []int{1, 2, 3}, []int{1, 2, 3})
	AssertEqual(t, map[string]int{"a": 1, "b": 2}, map[string]int{"b": 2, "a": 1})
	AssertNotEqual(t, "agent", "deploy")
}

func TestTest_AssertEqual_Bad(t *T) {
	// A plain NewError carries only its message.
	AssertEqual(t, NewError("refused").Error(), "refused")
	// E prefixes the operation: "op: msg".
	AssertEqual(t, E("agent.Dispatch", "session refused", nil).Error(), "agent.Dispatch: session refused")
	// Wrap nests the cause after the message.
	AssertEqual(t, Wrap(AnError, "homelab.Health", "down").Error(), "homelab.Health: down: core test sentinel error")
	AssertNotEqual(t, NewError("a").Error(), NewError("b").Error())
}

func TestTest_AssertEqual_Ugly(t *T) {
	var left []string
	var right []string
	AssertEqual(t, left, right)
}

func TestTest_AssertError_Good(t *T) {
	AssertError(t, AnError)
	// Optional substrings must each appear in err.Error().
	AssertError(t, AnError, "sentinel")
	AssertError(t, AnError, "core test sentinel error")
	AssertError(t, NewError("boom"), "boom")
}

func TestTest_AssertError_Bad(t *T) {
	AssertError(t, E("agent.Dispatch", "session token refused", nil), "session token refused")
	// Multiple substrings checked in a single call: operation + message.
	AssertError(t, E("agent.Dispatch", "session token refused", nil), "agent.Dispatch", "session token refused")
	AssertError(t, E("homelab.Deploy", "rollback", nil), "homelab.Deploy: rollback")
}

func TestTest_AssertError_Ugly(t *T) {
	AssertError(t, Wrap(AnError, "homelab.Health", "nested failure"), "nested failure")
	// The wrapped operation and the underlying cause both surface in the string.
	AssertError(t, Wrap(AnError, "homelab.Health", "nested failure"), "homelab.Health")
	AssertError(t, Wrap(AnError, "homelab.Health", "nested failure"), "core test sentinel error")
}

func TestTest_AssertErrorIs_Good(t *T) {
	AssertErrorIs(t, Wrap(AnError, "agent.Dispatch", "failed"), AnError)
	// Double-wrapped chains still unwrap to the sentinel.
	AssertErrorIs(t, Wrap(Wrap(AnError, "a", "b"), "c", "d"), AnError)
	// An error is always itself.
	sentinel := NewError("self")
	AssertErrorIs(t, sentinel, sentinel)
}

func TestTest_AssertErrorIs_Bad(t *T) {
	first := NewError("first")
	joined := ErrorJoin(first, AnError)
	// Both joined members are reachable via Is.
	AssertErrorIs(t, joined, AnError)
	AssertErrorIs(t, joined, first)
}

func TestTest_AssertErrorIs_Ugly(t *T) {
	// errors.Is(nil, nil) is true.
	AssertErrorIs(t, nil, nil)
	// Wrap(nil, ...) returns nil, so it also Is nil.
	AssertErrorIs(t, Wrap(nil, "agent.Skip", "no-op"), nil)
	sentinel := NewError("only")
	AssertErrorIs(t, sentinel, sentinel)
}

func TestTest_AssertFalse_Good(t *T) {
	AssertFalse(t, 1 == 2)
	AssertFalse(t, Contains("agent", "xyz"))
	AssertFalse(t, len("") > 0)
}

func TestTest_AssertFalse_Bad(t *T) {
	// The zero-value Result is not OK.
	AssertFalse(t, Result{}.OK)
	AssertFalse(t, Result{OK: false}.OK)
	AssertTrue(t, Result{OK: true}.OK)
}

func TestTest_AssertFalse_Ugly(t *T) {
	var token *string
	AssertFalse(t, token != nil)
	AssertTrue(t, token == nil)
	// After assignment the pointer is non-nil.
	v := "tok"
	token = &v
	AssertFalse(t, token == nil)
	AssertTrue(t, token != nil)
}

func TestTest_AssertGreater_Good(t *T) {
	AssertGreater(t, 8, 2)
	AssertGreater(t, 2.5, 2.4)
	// Mixed int/float comparison.
	AssertGreater(t, 3, 2.9)
	AssertLess(t, 2, 8)
}

func TestTest_AssertGreater_Bad(t *T) {
	// Strings compare lexicographically.
	AssertGreater(t, "session-token-b", "session-token-a")
	AssertGreater(t, "abc", "abb")
	AssertLess(t, "a", "b")
}

func TestTest_AssertGreater_Ugly(t *T) {
	// Mixed-sign: uint(1) > int(-1).
	AssertGreater(t, uint(1), -1)
	AssertGreater(t, 0, -1)
	AssertGreater(t, uint(5), uint(3))
	AssertGreater(t, int8(2), int8(-3))
}

func TestTest_AssertGreaterOrEqual_Good(t *T) {
	AssertGreaterOrEqual(t, 3, 3)
	AssertGreaterOrEqual(t, 4, 3)
	AssertGreaterOrEqual(t, 3.0, 3.0)
	AssertGreaterOrEqual(t, 5, 4.5)
}

func TestTest_AssertGreaterOrEqual_Bad(t *T) {
	// Equal strings satisfy >=.
	AssertGreaterOrEqual(t, "deploy", "deploy")
	AssertGreaterOrEqual(t, "deploz", "deploy")
	AssertGreaterOrEqual(t, "b", "a")
}

func TestTest_AssertGreaterOrEqual_Ugly(t *T) {
	AssertGreaterOrEqual(t, 0.5, 0)
	AssertGreaterOrEqual(t, 0, 0)
	AssertGreaterOrEqual(t, uint(1), -1)
	AssertGreaterOrEqual(t, 0.0, 0)
}

func TestTest_AssertInDelta_Good(t *T) {
	AssertInDelta(t, 1.00, 1.01, 0.02)
	// Exact equality with zero delta.
	AssertInDelta(t, 5.0, 5.0, 0.0)
	// Boundary: |want-got| == delta passes (fail only when strictly greater).
	AssertInDelta(t, 1.0, 1.5, 0.5)
	AssertInDelta(t, 100.0, 100.05, 0.1)
}

func TestTest_AssertInDelta_Bad(t *T) {
	AssertInDelta(t, 10, 10, 0)
	AssertInDelta(t, 0.0, 0.0, 0.0)
	AssertInDelta(t, -3.0, -3.0, 0.0)
}

func TestTest_AssertInDelta_Ugly(t *T) {
	AssertInDelta(t, -0.1, -0.1001, 0.001)
	AssertInDelta(t, -1.0, -1.2, 0.5)
	// Values straddling zero, within delta.
	AssertInDelta(t, -0.0001, 0.0001, 0.001)
}

func TestTest_AssertLen_Good(t *T) {
	AssertLen(t, []string{"agent", "dispatch"}, 2)
	// String length is measured in bytes.
	AssertLen(t, "agent", 5)
	AssertLen(t, []int{}, 0)
	AssertLen(t, [3]int{}, 3)
}

func TestTest_AssertLen_Bad(t *T) {
	AssertLen(t, map[string]bool{"refused": true}, 1)
	AssertLen(t, map[string]int{"a": 1, "b": 2, "c": 3}, 3)
	AssertLen(t, map[string]int{}, 0)
}

func TestTest_AssertLen_Ugly(t *T) {
	// An empty buffered channel has length 0.
	AssertLen(t, make(chan string, 3), 0)
	// A channel's length reflects its buffered, not-yet-received items.
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	AssertLen(t, ch, 2)
	AssertLen(t, "", 0)
}

func TestTest_AssertLess_Good(t *T) {
	AssertLess(t, 2, 8)
	AssertLess(t, 2.4, 2.5)
	AssertLess(t, 2.9, 3)
	AssertGreater(t, 8, 2)
}

func TestTest_AssertLess_Bad(t *T) {
	AssertLess(t, "agent-a", "agent-b")
	AssertLess(t, "a", "b")
	AssertGreater(t, "b", "a")
}

func TestTest_AssertLess_Ugly(t *T) {
	// Mixed-sign: int(-1) < uint(1).
	AssertLess(t, -1, uint(1))
	AssertLess(t, -1, 0)
	AssertLess(t, uint(3), uint(5))
}

func TestTest_AssertLessOrEqual_Good(t *T) {
	AssertLessOrEqual(t, 3, 3)
	AssertLessOrEqual(t, 2, 3)
	AssertLessOrEqual(t, 3.0, 3.0)
	AssertLessOrEqual(t, 4.5, 5)
}

func TestTest_AssertLessOrEqual_Bad(t *T) {
	AssertLessOrEqual(t, "deploy", "deploy")
	AssertLessOrEqual(t, "deploy", "deploz")
	AssertLessOrEqual(t, "a", "b")
}

func TestTest_AssertLessOrEqual_Ugly(t *T) {
	AssertLessOrEqual(t, 0, 0.5)
	AssertLessOrEqual(t, 0, 0)
	AssertLessOrEqual(t, -1, uint(1))
}

func TestTest_AssertNil_Good(t *T) {
	AssertNil(t, nil)
	// Typed-nil values of every nilable kind report nil.
	var p *int
	AssertNil(t, p)
	var m map[string]int
	AssertNil(t, m)
	var f func()
	AssertNil(t, f)
}

func TestTest_AssertNil_Bad(t *T) {
	var session *string
	AssertNil(t, session)
	var s []int
	AssertNil(t, s)
	var e error
	AssertNil(t, e)
	// Inverse: an addressable value yields a non-nil pointer.
	x := 5
	AssertNotNil(t, &x)
}

func TestTest_AssertNil_Ugly(t *T) {
	var meta map[string]string
	AssertNil(t, meta)
	var ch chan int
	AssertNil(t, ch)
	// An initialised (but empty) map is NOT nil.
	AssertNotNil(t, map[string]string{})
}

func TestTest_AssertNoError_Good(t *T) {
	AssertNoError(t, nil)
	// Wrap(nil, ...) collapses to nil.
	AssertNoError(t, Wrap(nil, "agent.Skip", "no-op"))
	// ErrorJoin of nothing is nil.
	AssertNoError(t, ErrorJoin())
}

func TestTest_AssertNoError_Bad(t *T) {
	// WrapCode with no cause and no code returns nil.
	AssertNoError(t, WrapCode(nil, "", "op", "m"))
	// Joining only nils yields nil.
	AssertNoError(t, ErrorJoin(nil, nil))
	// Inverse: a real error is reported.
	AssertError(t, AnError)
}

func TestTest_AssertNoError_Ugly(t *T) {
	var err error
	AssertNoError(t, err)
	AssertNoError(t, error(nil))
	AssertNoError(t, ErrorJoin(nil))
}

func TestTest_AssertNotContains_Good(t *T) {
	AssertNotContains(t, []string{"agent", "deploy"}, "rollback")
	// Membership is case-sensitive.
	AssertNotContains(t, []string{"agent"}, "AGENT")
	AssertNotContains(t, []int{1, 2, 3}, 4)
	// Inverse: a present element is found.
	AssertContains(t, []string{"agent", "deploy"}, "agent")
}

func TestTest_AssertNotContains_Bad(t *T) {
	AssertNotContains(t, map[string]int{"session": 1}, "missing")
	// Key match is case-sensitive.
	AssertNotContains(t, map[string]int{"session": 1}, "Session")
	AssertNotContains(t, map[int]string{1: "a"}, 2)
	// Inverse: a present key is found.
	AssertContains(t, map[string]int{"session": 1}, "session")
}

func TestTest_AssertNotContains_Ugly(t *T) {
	// A nil haystack has Invalid kind and matches nothing.
	AssertNotContains(t, nil, "agent")
	// Unsupported haystack kinds (int) never report membership.
	AssertNotContains(t, 42, "agent")
	// Non-empty needle is absent from an empty string.
	AssertNotContains(t, "", "x")
	// A non-string needle never matches a string haystack.
	AssertNotContains(t, "agent", 5)
}

func TestTest_AssertNotEmpty_Good(t *T) {
	AssertNotEmpty(t, "agent")
	AssertNotEmpty(t, []int{1})
	AssertNotEmpty(t, map[string]int{"a": 1})
	AssertNotEmpty(t, 1)
	// Inverse: the empty string is empty.
	AssertEmpty(t, "")
}

func TestTest_AssertNotEmpty_Bad(t *T) {
	AssertNotEmpty(t, []string{"refused"})
	// A length-1 array is non-empty even when its element is the zero value.
	AssertNotEmpty(t, [1]int{0})
	// A true bool differs from its zero value.
	AssertNotEmpty(t, true)
	AssertNotEmpty(t, "false")
}

func TestTest_AssertNotEmpty_Ugly(t *T) {
	AssertNotEmpty(t, struct{ Agent string }{Agent: "codex"})
	// Inverse: the zero-value struct is empty.
	AssertEmpty(t, struct{ Agent string }{})
	AssertNotEmpty(t, struct{ N int }{N: 1})
}

func TestTest_AssertNotEqual_Good(t *T) {
	AssertNotEqual(t, "agent-a", "agent-b")
	AssertNotEqual(t, 1, 2)
	AssertNotEqual(t, []int{1}, []int{2})
	// Inverse: equal values are equal.
	AssertEqual(t, "x", "x")
}

func TestTest_AssertNotEqual_Bad(t *T) {
	// Different dynamic types are never deeply equal.
	AssertNotEqual(t, 1, "1")
	AssertNotEqual(t, int64(1), int32(1))
	AssertNotEqual(t, 1.0, 1)
}

func TestTest_AssertNotEqual_Ugly(t *T) {
	// A nil slice and an empty slice are not deeply equal.
	AssertNotEqual(t, []string(nil), []string{})
	// Likewise a nil map versus an empty map.
	AssertNotEqual(t, map[string]int(nil), map[string]int{})
	// Inverse: two empty non-nil slices are equal.
	AssertEqual(t, []string{}, []string{})
}

func TestTest_AssertNotNil_Good(t *T) {
	AssertNotNil(t, "agent")
	AssertNotNil(t, 42)
	x := 1
	AssertNotNil(t, &x)
	AssertNotNil(t, []int{})
}

func TestTest_AssertNotNil_Bad(t *T) {
	// Empty-but-not-nil values are still non-nil.
	AssertNotNil(t, "")
	AssertNotNil(t, 0)
	AssertNotNil(t, false)
	AssertNotNil(t, []int{})
}

func TestTest_AssertNotNil_Ugly(t *T) {
	AssertNotNil(t, struct{}{})
	AssertNotNil(t, [0]int{})
	AssertNotNil(t, map[string]int{})
	AssertNotNil(t, make(chan int))
}

func TestTest_AssertNotPanics_Good(t *T) {
	AssertNotPanics(t, func() { /* no-op closure verifies non-panicking behaviour */ })
	AssertNotPanics(t, func() {
		err := E("agent.Run", "boom", AnError)
		AssertErrorIs(t, err, AnError)
	})
	AssertNotPanics(t, func() {
		items := []string{"agent", "deploy"}
		AssertContains(t, items, "agent")
	})
	AssertNotPanics(t, func() {
		AssertTrue(t, Contains("agent dispatch", "dispatch"))
	})
}

func TestTest_AssertNotPanics_Bad(t *T) {
	AssertNotPanics(t, func() {
		_ = Result{}
	})
}

func TestTest_AssertNotPanics_Ugly(t *T) {
	AssertNotPanics(t, func() {
		var opts Options
		_ = opts.Len()
	})
}

func TestTest_AssertPanics_Good(t *T) {
	AssertPanics(t, func() { panic("agent halted") })
	AssertPanics(t, func() { panic(NewError("session expired")) })
	AssertPanics(t, func() {
		var items []int
		_ = items[3]
	})
	AssertPanics(t, func() {
		var m map[string]int
		m["x"] = 1
	})
}

func TestTest_AssertPanics_Bad(t *T) {
	AssertPanics(t, func() { panic(AnError) })
	AssertPanics(t, func() { panic(42) })
	AssertPanics(t, func() {
		var x any = "string"
		_ = x.(int)
	})
}

func TestTest_AssertPanics_Ugly(t *T) {
	AssertPanics(t, func() {
		var values []string
		_ = values[1]
	})
}

func TestTest_AssertPanicsWithError_Good(t *T) {
	AssertPanicsWithError(t, "session token expired", func() {
		panic(NewError("session token expired"))
	})
}

func TestTest_AssertPanicsWithError_Bad(t *T) {
	AssertPanicsWithError(t, "entitlement denied", func() {
		panic("entitlement denied")
	})
}

func TestTest_AssertPanicsWithError_Ugly(t *T) {
	AssertPanicsWithError(t, "", func() {
		panic("")
	})
}

func TestTest_AssertSame_Good(t *T) {
	c := New()
	AssertSame(t, c, c.Core())
	// Subsystem accessors return stable singleton pointers.
	AssertSame(t, c.Fs(), c.Fs())
	AssertSame(t, c.Options(), c.Options())
	x := "agent"
	p := &x
	q := &x
	// Two independently taken addresses of the same variable are the
	// same pointer — a real identity check, not a self-comparison.
	AssertSame(t, p, q)
}

func TestTest_AssertSame_Bad(t *T) {
	log := Default()
	// A previously obtained reference matches a freshly obtained one —
	// Default() hands back the same stored logger each call.
	AssertSame(t, log, Default())
	AssertSame(t, Default(), Default())
	c := New()
	opts := c.Options()
	AssertSame(t, opts, c.Options())
}

func TestTest_AssertSame_Ugly(t *T) {
	var left *Core
	var right *Core
	AssertSame(t, left, right)
}

func TestTest_AssertTrue_Good(t *T) {
	AssertTrue(t, 1 < 2)
	AssertTrue(t, Contains("agent dispatch", "dispatch"))
	AssertTrue(t, Result{OK: true}.OK)
	AssertTrue(t, len([]string{"agent"}) == 1)
}

func TestTest_AssertTrue_Bad(t *T) {
	// The zero-value Result is not OK, so its negation is true.
	AssertTrue(t, !Result{}.OK)
	AssertTrue(t, !Contains("agent", "deploy"))
	AssertFalse(t, Result{}.OK)
}

func TestTest_AssertTrue_Ugly(t *T) {
	AssertTrue(t, len([]string{}) == 0)
	AssertTrue(t, len([]string{"a"}) == 1)
	AssertTrue(t, len("") == 0)
}

func TestTest_RequireNoError_Good(t *T) {
	RequireNoError(t, nil)
	RequireNoError(t, Wrap(nil, "agent.Skip", "no-op"))
	RequireNoError(t, ErrorJoin(nil, nil))
	// The test continues past the requires and can still assert.
	AssertNoError(t, ErrorJoin())
}

func TestTest_RequireNoError_Bad(t *T) {
	var err error
	RequireNoError(t, err)
	// WrapCode with no cause and no code is nil.
	RequireNoError(t, WrapCode(nil, "", "op", "m"))
	AssertNil(t, WrapCode(nil, "", "op", "m"))
}

func TestTest_RequireNoError_Ugly(t *T) {
	RequireNoError(t, nil)
	RequireNoError(t, error(nil))
	var e error
	RequireNoError(t, e)
	AssertNoError(t, e)
}

func TestTest_RequireNotEmpty_Good(t *T) {
	RequireNotEmpty(t, "agent")
	RequireNotEmpty(t, []string{"deploy"})
	RequireNotEmpty(t, map[string]int{"session": 1})
	RequireNotEmpty(t, 42)
}

func TestTest_RequireNotEmpty_Bad(t *T) {
	RequireNotEmpty(t, []string{"refused"})
	// A length-1 array is non-empty regardless of its element value.
	RequireNotEmpty(t, [1]int{0})
	RequireNotEmpty(t, true)
}

func TestTest_RequireNotEmpty_Ugly(t *T) {
	RequireNotEmpty(t, struct{ Agent string }{Agent: "codex"})
	// A pointer to a non-empty value is non-empty (dereferenced).
	s := "x"
	RequireNotEmpty(t, &s)
	RequireNotEmpty(t, struct{ N int }{N: 1})
}

func TestTest_RequireTrue_Good(t *T) {
	RequireTrue(t, 1 < 2)
	RequireTrue(t, Contains("agent dispatch", "agent"))
	RequireTrue(t, len([]int{1, 2, 3}) == 3)
	AssertTrue(t, Contains("deploy homelab", "homelab"))
}

func TestTest_RequireTrue_Bad(t *T) {
	RequireTrue(t, !Result{}.OK)
	RequireTrue(t, !Contains("agent", "deploy"))
	RequireTrue(t, Result{OK: true}.OK)
}

func TestTest_RequireTrue_Ugly(t *T) {
	RequireTrue(t, len([]string{}) == 0)
	RequireTrue(t, len("") == 0)
	RequireTrue(t, len([]int{1, 2}) == 2)
}

func TestTest_AllocsPerRun_Good(t *T) {
	AssertEqual(t, 0.0, AllocsPerRun(100, func() {}))
}

func TestTest_AllocsPerRun_Bad(t *T) {
	// An allocating body reports at least one alloc per run.
	AssertTrue(t, AllocsPerRun(100, func() { testAllocSink = make([]byte, 1024) }) >= 1)
}

func TestTest_AllocsPerRun_Ugly(t *T) {
	// Misuse-panic contract (stdlib): zero runs divides by zero. Same
	// documented-panic idiom as RandPick.
	AssertPanics(t, func() { AllocsPerRun(0, func() {}) })
}

var testAllocSink []byte
