// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

// --- Constructors ---

func TestReturn_ReturnOK_Good(t *T) {
	r := ReturnOK("brain")
	AssertTrue(t, r.OK())
	AssertEqual(t, "brain", r.Value)
}

func TestReturn_ReturnOK_Bad(t *T) {
	// The zero value is a successful Return of T's zero.
	var r Return[string]
	AssertTrue(t, r.OK())
	AssertEqual(t, "", r.Value)
}

func TestReturn_ReturnOK_Ugly(t *T) {
	// A typed nil is still success — Err is the only discriminant.
	r := ReturnOK[*Options](nil)
	AssertTrue(t, r.OK())
	AssertNil(t, r.Value)
}

func TestReturn_ReturnFail_Good(t *T) {
	cause := NewError("agent offline")
	r := ReturnFail[string](cause)
	AssertFalse(t, r.OK())
	AssertEqual(t, cause, r.Err)
}

func TestReturn_ReturnFail_Bad(t *T) {
	// A nil error makes a *successful* Return — Fail with nil is misuse
	// that degrades to OK, mirroring ReturnFrom(v, nil).
	r := ReturnFail[string](nil)
	AssertTrue(t, r.OK())
}

func TestReturn_ReturnFail_Ugly(t *T) {
	r := ReturnFail[int](NewCode("fs.notfound", "missing"))
	AssertEqual(t, 0, r.Value)
	AssertEqual(t, "fs.notfound", r.Code())
}

func TestReturn_ReturnFrom_Good(t *T) {
	r := ReturnFrom("payload", nil)
	AssertTrue(t, r.OK())
	AssertEqual(t, "payload", r.Value)
}

func TestReturn_ReturnFrom_Bad(t *T) {
	cause := NewError("read failed")
	r := ReturnFrom([]byte(nil), cause)
	AssertFalse(t, r.OK())
	AssertEqual(t, cause, r.Err)
}

func TestReturn_ReturnFrom_Ugly(t *T) {
	// A value AND an error: Err wins the discrimination, Value survives
	// for callers that inspect partial reads.
	r := ReturnFrom("partial", NewError("truncated"))
	AssertFalse(t, r.OK())
	AssertEqual(t, "partial", r.Value)
}

func TestReturn_ReturnOf_Good(t *T) {
	r := ReturnOf[string](Ok("brain"))
	AssertTrue(t, r.OK())
	AssertEqual(t, "brain", r.Value)
}

func TestReturn_ReturnOf_Bad(t *T) {
	cause := NewError("agent offline")
	r := ReturnOf[string](Fail(cause))
	AssertFalse(t, r.OK())
	AssertEqual(t, cause, r.Err)
}

func TestReturn_ReturnOf_Ugly(t *T) {
	// OK Result, wrong type: fails rather than zero-values silently.
	r := ReturnOf[int](Ok("not an int"))
	AssertFalse(t, r.OK())
	AssertError(t, r.Err)
}

func TestReturn_ReturnTry_Good(t *T) {
	r := ReturnTry(func() int { return 42 })
	AssertTrue(t, r.OK())
	AssertEqual(t, 42, r.Value)
}

func TestReturn_ReturnTry_Bad(t *T) {
	cause := NewError("parser exploded")
	r := ReturnTry(func() int { panic(cause) })
	AssertFalse(t, r.OK())
	AssertEqual(t, cause, r.Err)
}

func TestReturn_ReturnTry_Ugly(t *T) {
	// A non-error panic value is wrapped, never re-thrown.
	r := ReturnTry(func() string { panic("raw string panic") })
	AssertFalse(t, r.OK())
	AssertError(t, r.Err)
}

// --- Methods ---

func TestReturn_Return_OK_Good(t *T) {
	AssertTrue(t, ReturnOK(1).OK())
}

func TestReturn_Return_OK_Bad(t *T) {
	AssertFalse(t, ReturnFail[int](NewError("x")).OK())
}

func TestReturn_Return_OK_Ugly(t *T) {
	// OK is purely Err==nil — a zero Value does not mean failure.
	AssertTrue(t, ReturnOK(0).OK())
}

func TestReturn_Return_Or_Good(t *T) {
	AssertEqual(t, "live", ReturnOK("live").Or("fallback"))
}

func TestReturn_Return_Or_Bad(t *T) {
	AssertEqual(t, "fallback", ReturnFail[string](NewError("x")).Or("fallback"))
}

func TestReturn_Return_Or_Ugly(t *T) {
	// A successful zero value is returned, not the fallback.
	AssertEqual(t, "", ReturnOK("").Or("fallback"))
}

func TestReturn_Return_Must_Good(t *T) {
	AssertEqual(t, 42, ReturnOK(42).Must())
}

func TestReturn_Return_Must_Bad(t *T) {
	AssertPanics(t, func() {
		ReturnFail[int](NewError("agent offline")).Must()
	})
}

func TestReturn_Return_Must_Ugly(t *T) {
	// The panic value is the underlying error itself.
	AssertPanicsWithError(t, "agent offline", func() {
		ReturnFail[int](NewError("agent offline")).Must()
	})
}

func TestReturn_Return_Code_Good(t *T) {
	AssertEqual(t, "fs.notfound", ReturnFail[int](NewCode("fs.notfound", "missing")).Code())
}

func TestReturn_Return_Code_Bad(t *T) {
	AssertEqual(t, "", ReturnOK(1).Code())
}

func TestReturn_Return_Code_Ugly(t *T) {
	// E() sets Operation, not Code — Code stays empty, same as Result.Code.
	AssertEqual(t, "", ReturnFail[int](E("svc.Op", "failed", nil)).Code())
}

func TestReturn_Return_Result_Good(t *T) {
	r := ReturnOK("brain").Result()
	AssertTrue(t, r.OK)
	AssertEqual(t, "brain", r.Value)
}

func TestReturn_Return_Result_Bad(t *T) {
	cause := NewError("agent offline")
	r := ReturnFail[string](cause).Result()
	AssertFalse(t, r.OK)
	AssertEqual(t, cause, r.Err())
}

func TestReturn_Return_Result_Ugly(t *T) {
	// Round trip: Result → ReturnOf → Result preserves value and state.
	orig := Ok(42)
	back := ReturnOf[int](orig).Result()
	AssertEqual(t, orig, back)
}

func TestReturn_Return_Log_Good(t *T) {
	// No-op on success: the Return passes through untouched.
	r := ReturnOK("fine").Log("test.Op", "should not log")
	AssertTrue(t, r.OK())
	AssertEqual(t, "fine", r.Value)
}

func TestReturn_Return_Log_Bad(t *T) {
	cause := NewError("agent offline")
	r := ReturnFail[string](cause).Log("test.Op", "logged and passed through")
	AssertEqual(t, cause, r.Err)
}

func TestReturn_Return_Log_Ugly(t *T) {
	// Chainable in either state, always identity on the payload.
	r := ReturnFrom("v", nil).Log("a", "b").Log("c", "d")
	AssertEqual(t, "v", r.Value)
}
