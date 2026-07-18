package core_test

import (
	. "dappco.re/go"
)

func TestResult_Result_Error_Good(t *T) {
	r := Result{Value: NewError("agent dispatch failed"), OK: false}
	AssertEqual(t, "agent dispatch failed", r.Error())
}

func TestResult_Result_Error_Bad(t *T) {
	r := Result{Value: "session token refused", OK: false}
	AssertEqual(t, "session token refused", r.Error())
}

func TestResult_Result_Error_Ugly(t *T) {
	r := Result{OK: false}
	AssertEqual(t, "unknown error", r.Error())
}

func TestResult_Result_Err_Good(t *T) {
	AssertTrue(t, Result{Value: "ready", OK: true}.Err() == nil)
}

func TestResult_Result_Err_Bad(t *T) {
	err := NewError("dispatch failed")
	AssertEqual(t, err, Result{Value: err, OK: false}.Err())
}

func TestResult_Result_Err_Ugly(t *T) {
	// non-error failure value: Err() returns the Result itself as an error
	r := Result{Value: "session refused", OK: false}
	AssertEqual(t, "session refused", r.Err().Error())
}

func TestResult_Result_Code_Good(t *T) {
	r := Result{Value: NewCode("agent.refused", "dispatch refused"), OK: false}
	AssertEqual(t, "agent.refused", r.Code())
}

func TestResult_Result_Code_Bad(t *T) {
	r := Result{Value: NewError("plain failure"), OK: false}
	AssertEqual(t, "", r.Code())
}

func TestResult_Result_Code_Ugly(t *T) {
	r := Result{Value: NewCode("agent.refused", "dispatch refused"), OK: true}
	AssertEqual(t, "", r.Code())
}

func TestResult_Result_Must_Good(t *T) {
	r := Result{Value: "agent-ready", OK: true}
	AssertEqual(t, "agent-ready", r.Must())
}

func TestResult_Result_Must_Bad(t *T) {
	r := Result{Value: NewError("session token expired"), OK: false}
	AssertPanicsWithError(t, "session token expired", func() {
		_ = r.Must()
	})
}

func TestResult_Result_Must_Ugly(t *T) {
	r := Result{Value: "panic text", OK: false}
	AssertPanicsWithError(t, "panic text", func() {
		_ = r.Must()
	})
}

func TestResult_Result_Or_Good(t *T) {
	r := Result{Value: "primary agent", OK: true}
	AssertEqual(t, "primary agent", r.Or("fallback agent"))
}

func TestResult_Result_Or_Bad(t *T) {
	r := Result{Value: NewError("missing agent"), OK: false}
	AssertEqual(t, "fallback agent", r.Or("fallback agent"))
}

func TestResult_Result_Or_Ugly(t *T) {
	r := Result{Value: nil, OK: true}
	AssertNil(t, r.Or("fallback agent"))
}

func TestResult_Cast_Good(t *T) {
	value, ok := Cast[string](Result{Value: "codex", OK: true})
	AssertTrue(t, ok)
	AssertEqual(t, "codex", value)
}

func TestResult_Cast_Bad(t *T) {
	value, ok := Cast[string](Result{Value: "codex", OK: false})
	AssertFalse(t, ok)
	AssertEqual(t, "", value)
}

func TestResult_Cast_Ugly(t *T) {
	value, ok := Cast[int](Result{Value: "codex", OK: true})
	AssertFalse(t, ok)
	AssertEqual(t, 0, value)
}

func TestResult_Try_Good(t *T) {
	r := Try(func() any { return "dispatch-complete" })
	AssertTrue(t, r.OK)
	AssertEqual(t, "dispatch-complete", r.Value)
}

func TestResult_Try_Bad(t *T) {
	r := Try(func() any { return NewError("dispatch refused") })
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "dispatch refused")
}

func TestResult_Try_Ugly(t *T) {
	r := Try(func() any { panic("worker panic") })
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error), "panic recovered")
}

// --- Ok ---

func TestResult_Ok_Good(t *T) {
	r := Ok(42)
	AssertTrue(t, r.OK)
	AssertEqual(t, 42, r.Value)
}

func TestResult_Ok_Bad(t *T) {
	// Ok with nil value — still OK=true; consumer can opt to treat nil as a sentinel.
	r := Ok(nil)
	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
}

func TestResult_Ok_Ugly(t *T) {
	// Wrapping an error in Ok still produces OK=true (caller's choice).
	r := Ok(NewError("warning"))
	AssertTrue(t, r.OK)
	AssertNotNil(t, r.Value)
}

// --- Err ---

func TestResult_Fail_Good(t *T) {
	r := Fail(NewError("dispatch failed"))
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "dispatch failed")
}

func TestResult_Fail_Bad(t *T) {
	// Err with nil — produces OK=false, Value nil.
	r := Fail(nil)
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestResult_Fail_Ugly(t *T) {
	// Err with a coded error — Code() pulls through.
	r := Fail(NewCode("net.timeout", "homelab unreachable"))
	AssertFalse(t, r.OK)
	AssertEqual(t, "net.timeout", r.Code())
}

// --- ResultOf ---

func TestResult_ResultOf_Good(t *T) {
	r := ResultOf("payload", nil)
	AssertTrue(t, r.OK)
	AssertEqual(t, "payload", r.Value)
}

func TestResult_ResultOf_Bad(t *T) {
	r := ResultOf("payload", NewError("network"))
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "network")
}

func TestResult_ResultOf_Ugly(t *T) {
	// (nil, nil) is treated as success.
	r := ResultOf(nil, nil)
	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
}

// --- MustCast ---

func TestResult_MustCast_Good(t *T) {
	r := Ok("agent.dispatch")
	got := MustCast[string](r)
	AssertEqual(t, "agent.dispatch", got)
}

func TestResult_MustCast_Bad(t *T) {
	// Failed Result panics with the underlying error.
	AssertPanics(t, func() {
		_ = MustCast[string](Fail(NewError("boom")))
	})
}

func TestResult_MustCast_Ugly(t *T) {
	// Type mismatch panics with a descriptive error.
	AssertPanicsWithError(t, "not", func() {
		_ = MustCast[*int](Ok("string-not-pointer"))
	})
}

// --- Typed getters (the Options accessor dialect, output side) ---

func TestResult_Result_String_Good(t *T) {
	AssertEqual(t, "brain", Ok("brain").String())
}

func TestResult_Result_String_Bad(t *T) {
	AssertEqual(t, "agent offline", Fail(NewError("agent offline")).String())
}

func TestResult_Result_String_Ugly(t *T) {
	// The documented Stringer divergence: non-string OK values render.
	AssertEqual(t, "42", Ok(42).String())
}

func TestResult_Result_Int_Good(t *T) {
	AssertEqual(t, 8080, Ok(8080).Int())
}

func TestResult_Result_Int_Bad(t *T) {
	AssertEqual(t, 0, Fail(NewError("no port")).Int())
}

func TestResult_Result_Int_Ugly(t *T) {
	// Strict — an int64 does not promote to int.
	AssertEqual(t, 0, Ok(int64(9)).Int())
}

func TestResult_Result_Bool_Good(t *T) {
	AssertTrue(t, Ok(true).Bool())
}

func TestResult_Result_Bool_Bad(t *T) {
	AssertFalse(t, Fail(NewError("no flag")).Bool())
}

func TestResult_Result_Bool_Ugly(t *T) {
	// Strict — a "true" string never coerces.
	AssertFalse(t, Ok("true").Bool())
}

func TestResult_Result_Float64_Good(t *T) {
	AssertEqual(t, 1.5, Ok(1.5).Float64())
}

func TestResult_Result_Float64_Bad(t *T) {
	AssertEqual(t, 0.0, Fail(NewError("no weight")).Float64())
}

func TestResult_Result_Float64_Ugly(t *T) {
	// Promotes int/int64/float32 — the Options.Float64 contract.
	AssertEqual(t, 3.0, Ok(3).Float64())
	AssertEqual(t, float64(float32(2.5)), Ok(float32(2.5)).Float64())
}

func TestResult_Result_Duration_Good(t *T) {
	AssertEqual(t, 5*Second, Ok(5*Second).Duration())
}

func TestResult_Result_Duration_Bad(t *T) {
	AssertEqual(t, Duration(0), Fail(NewError("no timeout")).Duration())
}

func TestResult_Result_Duration_Ugly(t *T) {
	// A string Value parses via ParseDuration — the Options contract.
	AssertEqual(t, 30*Second, Ok("30s").Duration())
}

func TestResult_Result_Bytes_Good(t *T) {
	AssertEqual(t, "payload", string(Ok([]byte("payload")).Bytes()))
}

func TestResult_Result_Bytes_Bad(t *T) {
	AssertNil(t, Fail(NewError("no body")).Bytes())
}

func TestResult_Result_Bytes_Ugly(t *T) {
	// Strict — a string is not []byte.
	AssertNil(t, Ok("text").Bytes())
}

// (Err triplet lives beside Error above — the senior contract from the
// dedup-primitives lane: non-error failure Values return the Result
// itself, which satisfies error via Result.Error.)
