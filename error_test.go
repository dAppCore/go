package core_test

import (
	. "dappco.re/go"
)

// --- Error Creation ---

func TestError_E_Good(t *T) {
	err := E("user.Save", "failed to save", nil)
	AssertError(t, err)
	AssertContains(t, err.Error(), "user.Save")
	AssertContains(t, err.Error(), "failed to save")
}

func TestError_E_WithCause_Good(t *T) {
	cause := NewError("connection refused")
	err := E("db.Connect", "database unavailable", cause)
	AssertErrorIs(t, err, cause)
}

func TestError_Wrap_Good(t *T) {
	cause := NewError("timeout")
	err := Wrap(cause, "api.Call", "request failed")
	AssertError(t, err)
	AssertErrorIs(t, err, cause)
}

func TestError_Wrap_Nil_Good(t *T) {
	err := Wrap(nil, "api.Call", "request failed")
	AssertNil(t, err)
	// op/msg are irrelevant when the wrapped error is nil.
	AssertNil(t, Wrap(nil, "", ""))
	AssertNil(t, Wrap(nil, "op.only", ""))
	AssertNil(t, Wrap(nil, "", "msg only"))
}

func TestError_WrapCode_Good(t *T) {
	cause := NewError("invalid email")
	err := WrapCode(cause, "VALIDATION_ERROR", "user.Validate", "bad input")
	AssertError(t, err)
	AssertEqual(t, "VALIDATION_ERROR", ErrorCode(err))
}

func TestError_NewCode_Good(t *T) {
	err := NewCode("NOT_FOUND", "resource not found")
	AssertError(t, err)
	AssertEqual(t, "NOT_FOUND", ErrorCode(err))
}

// --- Error Introspection ---

func TestError_Operation_Good(t *T) {
	err := E("brain.Recall", "search failed", nil)
	AssertEqual(t, "brain.Recall", Operation(err))
	// Wrapping reports the outermost operation, not the cause's.
	AssertEqual(t, "agent.Dispatch", Operation(Wrap(err, "agent.Dispatch", "failed")))
	AssertEqual(t, "user.Validate", Operation(WrapCode(nil, "CODE", "user.Validate", "bad")))
	AssertEqual(t, "direct.Op", Operation(&Err{Operation: "direct.Op"}))
}

func TestError_Operation_Bad(t *T) {
	err := NewError("plain error")
	AssertEqual(t, "", Operation(err))
	// Foreign error types and op-less *Err both yield "".
	AssertEqual(t, "", Operation(&plainErr{msg: "external"}))
	AssertEqual(t, "", Operation(&Err{Message: "no op set"}))
	AssertEqual(t, "", Operation(E("", "msg", nil)))
}

func TestError_ErrorMessage_Good(t *T) {
	err := E("op", "the message", nil)
	AssertEqual(t, "the message", ErrorMessage(err))
	// Returns the *Err.Message field (not the formatted Error() string).
	AssertEqual(t, "outer msg", ErrorMessage(Wrap(err, "op2", "outer msg")))
	AssertEqual(t, "missing", ErrorMessage(NewCode("CODE", "missing")))
	AssertEqual(t, "direct", ErrorMessage(&Err{Operation: "x", Message: "direct"}))
}

func TestError_ErrorMessage_Plain(t *T) {
	err := NewError("plain")
	AssertEqual(t, "plain", ErrorMessage(err))
	// NewError produces an *Err, so the raw Message is returned verbatim.
	AssertEqual(t, "", ErrorMessage(NewError("")))
	AssertEqual(t, "line1\nline2", ErrorMessage(NewError("line1\nline2")))
	AssertEqual(t, "with spaces", ErrorMessage(&Err{Message: "with spaces"}))
}

func TestError_ErrorMessage_Nil(t *T) {
	AssertEqual(t, "", ErrorMessage(nil))
	var e error
	AssertEqual(t, "", ErrorMessage(e))
	// The nil guard must short-circuit before the err.Error() fallback,
	// which would otherwise nil-deref.
	AssertNotPanics(t, func() { ErrorMessage(nil) })
}

func TestError_Root_Good(t *T) {
	root := NewError("root cause")
	wrapped := Wrap(root, "layer1", "first wrap")
	double := Wrap(wrapped, "layer2", "second wrap")
	AssertEqual(t, root, Root(double))
}

func TestError_Root_Nil(t *T) {
	AssertNil(t, Root(nil))
	var e error
	AssertNil(t, Root(e))
	AssertNotPanics(t, func() { Root(nil) })
}

func TestError_StackTrace_Good(t *T) {
	err := Wrap(E("inner", "cause", nil), "outer", "wrapper")
	stack := StackTrace(err)
	AssertLen(t, stack, 2)
	AssertEqual(t, "outer", stack[0])
	AssertEqual(t, "inner", stack[1])
}

func TestError_FormatStackTrace_Good(t *T) {
	err := Wrap(E("a", "x", nil), "b", "y")
	formatted := FormatStackTrace(err)
	AssertEqual(t, "b -> a", formatted)
}

// --- ErrorLog ---

func TestError_ErrorLog_Good(t *T) {
	c := New()
	cause := NewError("boom")
	var el *ErrorLog = c.Log()
	r := el.Error(cause, "test.Operation", "something broke")
	AssertFalse(t, r.OK)
	AssertErrorIs(t, r.Value.(error), cause)
}

func TestError_ErrorLog_Nil_Good(t *T) {
	c := New()
	var el *ErrorLog = c.Log()
	r := el.Error(nil, "test.Operation", "no error")
	AssertTrue(t, r.OK)
}

func TestError_ErrorLog_Warn_Good(t *T) {
	c := New()
	cause := NewError("warning")
	r := c.Log().Warn(cause, "test.Operation", "heads up")
	AssertFalse(t, r.OK)
}

func TestError_ErrorLog_Must_Ugly(t *T) {
	c := New()
	AssertPanics(t, func() {
		c.Log().Must(NewError("fatal"), "test.Operation", "must fail")
	})
}

func TestError_ErrorLog_Must_Nil_Good(t *T) {
	c := New()
	AssertNotPanics(t, func() {
		c.Log().Must(nil, "test.Operation", "no error")
	})
}

// --- ErrorPanic ---

func TestError_ErrorPanic_Recover_Good(t *T) {
	c := New()
	// Should not panic — Recover catches it
	AssertNotPanics(t, func() {
		defer c.Error().Recover()
		panic("test panic")
	})
}

func TestError_ErrorPanic_SafeGo_Good(t *T) {
	c := New()
	done := make(chan bool, 1)
	c.Error().SafeGo(func() {
		done <- true
	})
	AssertTrue(t, <-done)
}

func TestError_ErrorPanic_SafeGo_Panic_Good(t *T) {
	c := New()
	done := make(chan bool, 1)
	c.Error().SafeGo(func() {
		defer func() { done <- true }()
		panic("caught by SafeGo")
	})
	// SafeGo recovers — goroutine completes without crashing the process
	<-done
}

// --- Standard Library Wrappers ---

func TestError_Is_Good(t *T) {
	target := NewError("target")
	wrapped := Wrap(target, "op", "msg")
	AssertTrue(t, Is(wrapped, target))
}

func TestError_As_Good(t *T) {
	err := E("op", "msg", nil)
	var e *Err
	AssertTrue(t, As(err, &e))
	AssertEqual(t, "op", e.Operation)
}

func TestError_NewError_Good(t *T) {
	err := NewError("simple error")
	AssertEqual(t, "simple error", err.Error())
	// Documented contract: NewError returns an *Err so introspection works.
	AssertEqual(t, "simple error", ErrorMessage(err))
	AssertEqual(t, "", ErrorCode(err))
	AssertEqual(t, "", Operation(err))
	var e *Err
	AssertTrue(t, As(err, &e))
	AssertEqual(t, "simple error", e.Message)
}

func TestError_ErrorJoin_Good(t *T) {
	e1 := NewError("first")
	e2 := NewError("second")
	joined := ErrorJoin(e1, e2)
	AssertErrorIs(t, joined, e1)
	AssertErrorIs(t, joined, e2)
}

// --- ErrorPanic Crash Reports ---

func TestError_ErrorPanic_Reports_Good(t *T) {
	dir := t.TempDir()
	path := Path(dir, "crashes.json")

	// Create ErrorPanic with file output
	c := New()
	// Access internals via a crash that writes to file
	// Since ErrorPanic fields are unexported, we test via Recover
	_ = c
	_ = path
	// Crash reporting needs ErrorPanic configured with filePath — tested indirectly
}

// --- ErrorPanic Crash File ---

func TestError_ErrorPanic_CrashFile_Good(t *T) {
	dir := t.TempDir()
	path := Path(dir, "crashes.json")

	// Create Core, trigger a panic through SafeGo, check crash file
	// ErrorPanic.filePath is unexported — but we can test via the package-level
	// error handling that writes crash reports

	// For now, test that Reports handles missing file gracefully
	c := New()
	r := c.Error().Reports(5)
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
	_ = path
}

// --- Error formatting branches ---

func TestError_Err_Error_WithCode_Good(t *T) {
	err := WrapCode(NewError("bad"), "INVALID", "validate", "input failed")
	AssertContains(t, err.Error(), "[INVALID]")
	AssertContains(t, err.Error(), "validate")
	AssertContains(t, err.Error(), "bad")
}

func TestError_Err_Error_CodeNoCause_Good(t *T) {
	err := NewCode("NOT_FOUND", "resource missing")
	AssertContains(t, err.Error(), "[NOT_FOUND]")
	AssertContains(t, err.Error(), "resource missing")
}

func TestError_Err_Error_NoOp_Good(t *T) {
	err := &Err{Message: "bare error"}
	AssertEqual(t, "bare error", err.Error())
	// With no Operation there is never an "op: " prefix, across every branch.
	AssertEqual(t, "bare [CODE]", (&Err{Message: "bare", Code: "CODE"}).Error())
	AssertEqual(t, "bare: cause", (&Err{Message: "bare", Cause: NewError("cause")}).Error())
	AssertEqual(t, "bare [CODE]: cause", (&Err{Message: "bare", Code: "CODE", Cause: NewError("cause")}).Error())
}

func TestError_WrapCode_NilErr_EmptyCode_Good(t *T) {
	err := WrapCode(nil, "", "op", "msg")
	AssertNil(t, err)
	// Nil is returned ONLY when both err is nil AND code is empty.
	AssertNil(t, WrapCode(nil, "", "", ""))
	AssertNotNil(t, WrapCode(nil, "CODE", "op", "msg"))
	AssertNotNil(t, WrapCode(NewError("x"), "", "op", "msg"))
}

func TestError_Wrap_PreservesCode_Good(t *T) {
	inner := WrapCode(NewError("root"), "AUTH_FAIL", "auth", "denied")
	outer := Wrap(inner, "handler", "request failed")
	AssertEqual(t, "AUTH_FAIL", ErrorCode(outer))
}

func TestError_LogWarn_Nil_Good(t *T) {
	c := New()
	r := c.LogWarn(nil, "op", "msg")
	AssertTrue(t, r.OK)
}

func TestError_LogError_Nil_Good(t *T) {
	c := New()
	r := c.LogError(nil, "op", "msg")
	AssertTrue(t, r.OK)
}

func TestError_AllOperations_Good(t *T) {
	err := Wrap(E("agent.Token", "expired", nil), "agent.Dispatch", "failed")
	var ops []string
	for op := range AllOperations(err) {
		ops = append(ops, op)
	}
	AssertEqual(t, []string{"agent.Dispatch", "agent.Token"}, ops)
}

func TestError_AllOperations_Bad(t *T) {
	var ops []string
	for op := range AllOperations(NewError("plain failure")) {
		ops = append(ops, op)
	}
	AssertEmpty(t, ops)
}

func TestError_AllOperations_Ugly(t *T) {
	var ops []string
	for op := range AllOperations(nil) {
		ops = append(ops, op)
	}
	AssertEmpty(t, ops)
}

func TestError_AllOperations_Break_Bad(t *T) {
	err := Wrap(E("agent.Token", "expired", nil), "agent.Dispatch", "failed")
	var ops []string
	for op := range AllOperations(err) {
		ops = append(ops, op)
		break
	}

	AssertEqual(t, []string{"agent.Dispatch"}, ops)
}

type plainErr struct{ msg string }

func (e *plainErr) Error() string { return e.msg }

func TestError_As_Bad(t *T) {
	// A non-*Err error never matches the *Err target.
	var structured *Err
	AssertFalse(t, As(&plainErr{msg: "plain failure"}, &structured))
	AssertNil(t, structured)
}

func TestError_As_Ugly(t *T) {
	var structured *Err
	AssertFalse(t, As(nil, &structured))
	AssertNil(t, structured)
}

func TestError_E_Bad(t *T) {
	err := E("", "", nil)
	AssertError(t, err)
	AssertEqual(t, "", err.Error())
}

func TestError_E_Ugly(t *T) {
	err := E("agent.Dispatch", "", AnError)
	AssertContains(t, err.Error(), "agent.Dispatch")
	AssertErrorIs(t, err, AnError)
}

func TestError_Err_Error_Good(t *T) {
	err := &Err{Operation: "agent.Dispatch", Message: "failed", Cause: AnError, Code: "agent.failed"}
	AssertContains(t, err.Error(), "agent.Dispatch")
	AssertContains(t, err.Error(), "[agent.failed]")
	AssertContains(t, err.Error(), AnError.Error())
}

func TestError_Err_Error_Bad(t *T) {
	err := &Err{}
	AssertEqual(t, "", err.Error())
	// Degenerate single-field forms still render the documented shape.
	AssertEqual(t, "agent.Run: ", (&Err{Operation: "agent.Run"}).Error())
	AssertEqual(t, " [CODE]", (&Err{Code: "CODE"}).Error())
	AssertEqual(t, ": c", (&Err{Cause: NewError("c")}).Error())
}

func TestError_Err_Error_Ugly(t *T) {
	err := &Err{Message: "session refused", Code: "session.refused"}
	AssertEqual(t, "session refused [session.refused]", err.Error())
	// All four fields set: op prefix, message, [code], then cause.
	full := &Err{Operation: "agent.Dispatch", Message: "boom", Code: "X", Cause: NewError("root")}
	AssertEqual(t, "agent.Dispatch: boom [X]: root", full.Error())
	// Embedded newlines/tabs/brackets pass through unescaped.
	weird := &Err{Operation: "op\n", Message: "msg\t", Code: "[nested]"}
	AssertEqual(t, "op\n: msg\t [[nested]]", weird.Error())
}

func TestError_Err_Unwrap_Good(t *T) {
	err := &Err{Cause: AnError}
	AssertEqual(t, AnError, err.Unwrap())
	// Unwrap returns the exact same cause pointer and drives errors.Is.
	AssertSame(t, AnError, err.Unwrap())
	AssertErrorIs(t, err, AnError)
	nested := &Err{Cause: E("inner", "x", AnError)}
	AssertErrorIs(t, nested.Unwrap(), AnError)
}

func TestError_Err_Unwrap_Bad(t *T) {
	err := &Err{}
	AssertNil(t, err.Unwrap())
	// A causeless Err unwraps to nil regardless of its other fields...
	AssertNil(t, (&Err{Operation: "op", Message: "msg", Code: "C"}).Unwrap())
	// ...and is therefore its own root.
	AssertEqual(t, err, Root(err))
}

func TestError_Err_Unwrap_Ugly(t *T) {
	root := NewCode("agent.refused", "dispatch refused")
	err := &Err{Cause: Wrap(root, "agent.Dispatch", "failed")}
	AssertErrorIs(t, err.Unwrap(), root)
}

func TestError_ErrorCode_Good(t *T) {
	err := NewCode("agent.refused", "dispatch refused")
	AssertEqual(t, "agent.refused", ErrorCode(err))
	AssertEqual(t, "VALIDATION", ErrorCode(WrapCode(AnError, "VALIDATION", "op", "msg")))
	// Wrap preserves the inner code through the chain.
	AssertEqual(t, "agent.refused", ErrorCode(Wrap(err, "outer", "wrapped")))
	AssertEqual(t, "DIRECT", ErrorCode(&Err{Code: "DIRECT"}))
}

func TestError_ErrorCode_Bad(t *T) {
	AssertEqual(t, "", ErrorCode(NewError("plain failure")))
	AssertEqual(t, "", ErrorCode(&plainErr{msg: "external"}))
	AssertEqual(t, "", ErrorCode(E("op", "msg", nil)))
	AssertEqual(t, "", ErrorCode(&Err{Message: "no code"}))
}

func TestError_ErrorCode_Ugly(t *T) {
	AssertEqual(t, "", ErrorCode(nil))
	var e error
	AssertEqual(t, "", ErrorCode(e))
	AssertNotPanics(t, func() { ErrorCode(nil) })
}

func TestError_ErrorJoin_Bad(t *T) {
	AssertNil(t, ErrorJoin(nil, nil))
	AssertNil(t, ErrorJoin())
	AssertNil(t, ErrorJoin(nil))
	AssertNil(t, ErrorJoin(nil, nil, nil))
}

func TestError_ErrorJoin_Ugly(t *T) {
	joined := ErrorJoin(nil, AnError)
	AssertErrorIs(t, joined, AnError)
	// nil entries are dropped; every non-nil member stays matchable.
	other := NewError("second")
	multi := ErrorJoin(nil, AnError, nil, other)
	AssertErrorIs(t, multi, AnError)
	AssertErrorIs(t, multi, other)
	AssertContains(t, multi.Error(), AnError.Error())
}

func TestError_ErrorLog_Error_Good(t *T) {
	r := New().Log().Error(AnError, "agent.Dispatch", "failed")
	AssertFalse(t, r.OK)
	AssertErrorIs(t, r.Value.(error), AnError)
}

func TestError_ErrorLog_Error_Bad(t *T) {
	r := New().Log().Error(nil, "agent.Dispatch", "no failure")
	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
	// A bare ErrorLog (no Core/logger) also short-circuits on nil err.
	r2 := (&ErrorLog{}).Error(nil, "op", "msg")
	AssertTrue(t, r2.OK)
	AssertNil(t, r2.Value)
}

func TestError_ErrorLog_Error_Ugly(t *T) {
	r := (&ErrorLog{}).Error(AnError, "agent.Dispatch", "failed")
	AssertFalse(t, r.OK)
	AssertErrorIs(t, r.Value.(error), AnError)
}

func TestError_ErrorLog_Must_Good(t *T) {
	AssertNotPanics(t, func() {
		New().Log().Must(nil, "agent.Dispatch", "ready")
	})
}

func TestError_ErrorLog_Must_Bad(t *T) {
	AssertPanicsWithError(t, "dispatch failed", func() {
		New().Log().Must(AnError, "agent.Dispatch", "dispatch failed")
	})
}

func TestError_ErrorLog_Warn_Bad(t *T) {
	r := New().Log().Warn(nil, "agent.Dispatch", "no warning")
	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
	// A bare ErrorLog (no Core/logger) also short-circuits on nil err.
	r2 := (&ErrorLog{}).Warn(nil, "op", "msg")
	AssertTrue(t, r2.OK)
	AssertNil(t, r2.Value)
}

func TestError_ErrorLog_Warn_Ugly(t *T) {
	r := (&ErrorLog{}).Warn(AnError, "agent.Dispatch", "degraded")
	AssertFalse(t, r.OK)
	AssertErrorIs(t, r.Value.(error), AnError)
}

func TestError_ErrorMessage_Bad(t *T) {
	AssertEqual(t, "plain failure", ErrorMessage(&plainErr{msg: "plain failure"}))
	// Non-*Err falls back to err.Error(), including the empty-message case.
	AssertEqual(t, "", ErrorMessage(&plainErr{msg: ""}))
	// A join of foreign errors (no *Err in the tree) returns its full string.
	joined := ErrorJoin(&plainErr{msg: "a"}, &plainErr{msg: "b"})
	AssertEqual(t, joined.Error(), ErrorMessage(joined))
}

func TestError_ErrorMessage_Ugly(t *T) {
	AssertEqual(t, "", ErrorMessage(nil))
	// Weird-but-valid messages survive verbatim.
	AssertEqual(t, "first\nsecond", ErrorMessage(E("op", "first\nsecond", nil)))
	AssertEqual(t, "  spaced  ", ErrorMessage(&Err{Message: "  spaced  "}))
	// Foreign error: full Error() string is returned, brackets and all.
	AssertEqual(t, "msg [CODE]", ErrorMessage(&plainErr{msg: "msg [CODE]"}))
}

func TestError_ErrorPanic_Recover_Bad(t *T) {
	AssertNotPanics(t, func() {
		New().Error().Recover()
	})
}

func TestError_ErrorPanic_Recover_Ugly(t *T) {
	var h *ErrorPanic
	AssertNotPanics(t, func() {
		h.Recover()
	})
}

func TestError_ErrorPanic_Reports_Bad(t *T) {
	r := New().Error().Reports(1)
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestError_ErrorPanic_Reports_Ugly(t *T) {
	r := New().Error().Reports(0)
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestError_ErrorPanic_SafeGo_Bad(t *T) {
	done := make(chan bool, 1)
	New().Error().SafeGo(func() {
		defer func() { done <- true }()
		panic("agent worker failed")
	})
	AssertTrue(t, <-done)
}

func TestError_ErrorPanic_SafeGo_Ugly(t *T) {
	done := make(chan bool, 1)
	New().Error().SafeGo(func() {
		done <- true
	})
	AssertTrue(t, <-done)
}

func TestError_FormatStackTrace_Bad(t *T) {
	AssertEqual(t, "", FormatStackTrace(NewError("plain failure")))
	// No operational context anywhere in the chain -> empty trace.
	AssertEqual(t, "", FormatStackTrace(&plainErr{msg: "external"}))
	AssertEqual(t, "", FormatStackTrace(E("", "no op", nil)))
	AssertEqual(t, "", FormatStackTrace(&Err{Message: "bare"}))
}

func TestError_FormatStackTrace_Ugly(t *T) {
	AssertEqual(t, "", FormatStackTrace(nil))
	var e error
	AssertEqual(t, "", FormatStackTrace(e))
	AssertNotPanics(t, func() { FormatStackTrace(nil) })
}

func TestError_Is_Bad(t *T) {
	AssertFalse(t, Is(NewError("left"), NewError("right")))
	AssertFalse(t, Is(NewError("left"), AnError))
	// A non-nil error never matches a nil target, even through a wrap.
	AssertFalse(t, Is(nil, AnError))
	AssertFalse(t, Is(Wrap(NewError("root"), "op", "msg"), AnError))
}

func TestError_Is_Ugly(t *T) {
	AssertTrue(t, Is(nil, nil))
	AssertTrue(t, Is(AnError, AnError))
	// A non-nil error is not "is nil"; a wrapped sentinel still matches.
	AssertFalse(t, Is(NewError("x"), nil))
	AssertTrue(t, Is(Wrap(AnError, "op", "msg"), AnError))
}

func TestError_NewCode_Bad(t *T) {
	err := NewCode("", "dispatch refused")
	AssertEqual(t, "dispatch refused", err.Error())
	AssertEqual(t, "", ErrorCode(err))
}

func TestError_NewCode_Ugly(t *T) {
	err := NewCode("", "")
	AssertEqual(t, "", err.Error())
	AssertEqual(t, "", ErrorCode(err))
	// A code with an empty message still renders the bracketed code.
	AssertEqual(t, " [CODE]", NewCode("CODE", "").Error())
	AssertEqual(t, "CODE", ErrorCode(NewCode("CODE", "")))
}

func TestError_NewError_Bad(t *T) {
	err := NewError("")
	AssertEqual(t, "", err.Error())
	// Empty text yields empty values from every introspection helper.
	AssertEqual(t, "", ErrorMessage(err))
	AssertEqual(t, "", Operation(err))
	AssertEqual(t, "", ErrorCode(err))
}

func TestError_NewError_Ugly(t *T) {
	err := NewError("session\nrefused")
	AssertContains(t, err.Error(), "session\nrefused")
	// Newlines and unicode are preserved exactly, not escaped or trimmed.
	AssertEqual(t, "session\nrefused", err.Error())
	AssertEqual(t, "session\nrefused", ErrorMessage(err))
	AssertEqual(t, "café ☕", NewError("café ☕").Error())
}

func TestError_Operation_Ugly(t *T) {
	AssertEqual(t, "", Operation(nil))
	var e error
	AssertEqual(t, "", Operation(e))
	AssertNotPanics(t, func() { Operation(nil) })
}

func TestError_Root_Bad(t *T) {
	err := NewError("plain failure")
	AssertEqual(t, err, Root(err))
	// An unwrappable error is its own root (same pointer), foreign or not.
	AssertSame(t, err, Root(err))
	foreign := &plainErr{msg: "external"}
	AssertSame(t, foreign, Root(foreign))
	noCause := E("op", "msg", nil)
	AssertSame(t, noCause, Root(noCause))
}

func TestError_Root_Ugly(t *T) {
	AssertNil(t, Root(nil))
	// A pathologically deep wrap chain resolves to the original root.
	root := NewError("deep root")
	deep := Wrap(Wrap(Wrap(root, "a", "1"), "b", "2"), "c", "3")
	AssertSame(t, root, Root(deep))
	AssertEqual(t, "deep root", Root(deep).Error())
}

func TestError_StackTrace_Bad(t *T) {
	AssertEmpty(t, StackTrace(NewError("plain failure")))
	// No operations anywhere -> empty slice, regardless of error shape.
	AssertEmpty(t, StackTrace(&plainErr{msg: "external"}))
	AssertEmpty(t, StackTrace(E("", "no op", nil)))
	AssertEmpty(t, StackTrace(&Err{Message: "bare"}))
}

func TestError_StackTrace_Ugly(t *T) {
	AssertEmpty(t, StackTrace(nil))
	var e error
	AssertEmpty(t, StackTrace(e))
	// Deep chain: operations come out outermost-first.
	deep := Wrap(Wrap(E("inner", "x", nil), "mid", "y"), "outer", "z")
	AssertEqual(t, []string{"outer", "mid", "inner"}, StackTrace(deep))
}

func TestError_Wrap_Bad(t *T) {
	AssertNil(t, Wrap(nil, "agent.Dispatch", "failed"))
	AssertNil(t, Wrap(nil, "", ""))
	// A real cause always wraps; a code-less cause yields a code-less wrap.
	AssertNotNil(t, Wrap(AnError, "op", "msg"))
	AssertEqual(t, "", ErrorCode(Wrap(NewError("x"), "op", "msg")))
}

func TestError_Wrap_Ugly(t *T) {
	inner := NewCode("agent.refused", "dispatch refused")
	err := Wrap(inner, "agent.Dispatch", "failed")
	AssertEqual(t, "agent.refused", ErrorCode(err))
}

func TestError_WrapCode_Bad(t *T) {
	AssertNil(t, WrapCode(nil, "", "agent.Dispatch", "failed"))
	AssertNil(t, WrapCode(nil, "", "", ""))
	// A code with no cause still produces an error carrying op + code.
	coded := WrapCode(nil, "CODE", "op", "msg")
	AssertEqual(t, "CODE", ErrorCode(coded))
	AssertEqual(t, "op", Operation(coded))
}

func TestError_WrapCode_Ugly(t *T) {
	err := WrapCode(nil, "agent.refused", "agent.Dispatch", "failed")
	AssertError(t, err)
	AssertEqual(t, "agent.refused", ErrorCode(err))
	AssertNil(t, Root(err).(*Err).Cause)
}
