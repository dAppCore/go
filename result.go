// SPDX-License-Identifier: EUPL-1.2

// Result ergonomics — methods and free functions that collapse the
// common Result-handling patterns into one-liners. Result itself is
// defined in options.go alongside Options as a Core primitive; this
// file extends it with the call-site helpers.
//
//	user := core.MustCast[*User](c.Drive().Get(opts))
//	timeout := core.HTTPGet(url).Or(defaultResp)
package core

// Ok wraps v in a successful Result. The canonical "happy path"
// constructor — replaces the awkward `Result{v, true}` literal.
//
//	return core.Ok(parsed)
func Ok(v any) Result {
	return Result{Value: v, OK: true}
}

// Fail wraps err in a failed Result. The canonical "sad path"
// constructor — pair Result.Code() / Result.Error() to introspect.
// (Named Fail rather than Err to avoid colliding with the *Err type
// in error.go.)
//
//	if err := decode(b); err != nil { return core.Fail(err) }
func Fail(err error) Result {
	return Result{Value: err, OK: false}
}

// ResultOf adapts a stdlib (value, error) pair into a Result —
// OK=false carrying err when err != nil, OK=true carrying v otherwise.
// The replacement for Result{}.New(v, err).
//
//	r := core.ResultOf(os.ReadFile(path))
//	if !r.OK { return r }
func ResultOf(v any, err error) Result {
	if err != nil {
		return Result{Value: err, OK: false}
	}
	return Result{Value: v, OK: true}
}

// Error returns the error message when the Result represents a failure,
// or "" when OK. Convenience for logging without unwrapping Value.
//
//	if !r.OK { core.Error("dispatch failed", "err", r.Error()) }
func (r Result) Error() string {
	if r.OK {
		return ""
	}
	if err, ok := r.Value.(error); ok {
		return err.Error()
	}
	if s, ok := r.Value.(string); ok {
		return s
	}
	return "unknown error"
}

// Err returns the failure as a Go error — nil when OK, the unwrapped error Value
// when it is one, otherwise the Result itself (which satisfies error via
// [Result.Error]). The idiomatic bridge from a Result to an error at an API
// boundary, replacing the hand-rolled "if !r.OK { return r.Value.(error) }"
// unwrap scattered across consumers.
//
//	if err := core.JSONUnmarshal(data, &cfg).Err(); err != nil {
//	    return err
//	}
func (r Result) Err() error {
	if r.OK {
		return nil
	}
	if err, ok := r.Value.(error); ok {
		return err
	}
	return r
}

// Code returns the stable error code from the Result's failure, or ""
// when OK or when the failure isn't a *core.Err with a Code populated.
// Codes form a flat keyspace agents grep on (e.g. "fs.notfound",
// "json.invalid", "http.timeout", "crypto.algo.unsupported"). See the
// Stable codespace section in AGENTS.md for the canonical list.
//
//	r := c.Fs().Read("/missing")
//	if r.Code() == "fs.notfound" { core.Println("first run") }
//	switch r.Code() {
//	case "http.timeout": retry()
//	case "http.refused": fallback()
//	}
func (r Result) Code() string {
	if r.OK {
		return ""
	}
	if e, ok := r.Value.(*Err); ok {
		return e.Code
	}
	return ""
}

// Must returns Value when OK; panics with the underlying error when
// not. Use for fast-fail paths — init, test setup, must-have config.
// Production request paths should check r.OK and return r.
//
//	cfg := core.JSONUnmarshal(data, &Config{}).Must().(*Config)
//	dir := core.PathAbs(".").Must().(string)
func (r Result) Must() any {
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			panic(err)
		}
		panic(r.Value)
	}
	return r.Value
}

// Or returns Value when OK, fallback otherwise. Convenience for
// optional reads where a default is acceptable.
//
//	port := core.EnvGet("PORT").Or("8080").(string)
//	timeout := core.ParseDuration(s).Or(5 * core.Second).(Duration)
func (r Result) Or(fallback any) any {
	if r.OK {
		return r.Value
	}
	return fallback
}

// --- Typed getters (the Options accessor dialect, output side) ---
//
// Each getter returns the type's zero value when the Result is failed
// or Value isn't that type — the same contract as the Options accessors
// (options.go), so the universal input and universal output speak one
// dialect. String is the single documented divergence (Stringer-first).

// String renders the Result for humans and logs — Result implements
// Stringer through it. Returns the string Value verbatim when OK holds
// a string, Sprint(Value) for any other OK value, and the Error() text
// when failed. This is the ONE deliberate divergence from the strict
// getter contract, so `%v` of any Result reads well in logs. For a
// strict typed read use Cast[string](r).
//
//	name := core.EnvGet("USER").String()
//	core.Println(c.Fs().Read("/missing").String())  // error text, not ""
func (r Result) String() string {
	if !r.OK {
		return r.Error()
	}
	if s, ok := r.Value.(string); ok {
		return s
	}
	return Sprint(r.Value)
}

// Int retrieves an int Value, 0 when failed or not an int.
// Strict — no numeric promotion, mirroring Options.Int.
//
//	port := c.QUERY(portQuery{}).Int()
func (r Result) Int() int {
	if !r.OK {
		return 0
	}
	i, _ := r.Value.(int)
	return i
}

// Bool retrieves a bool Value, false when failed or not a bool.
//
//	enabled := c.QUERY(flagQuery{}).Bool()
func (r Result) Bool() bool {
	if !r.OK {
		return false
	}
	b, _ := r.Value.(bool)
	return b
}

// Float64 retrieves a float64 Value, 0 when failed or wrong type.
// Promotes int/int64/float32 so JSON-decoded numbers work uniformly —
// the same promotion set as Options.Float64.
//
//	weight := c.QUERY(weightQuery{}).Float64()
func (r Result) Float64() float64 {
	if !r.OK {
		return 0
	}
	switch v := r.Value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	return 0
}

// Duration retrieves a Duration Value, 0 when failed or wrong type.
// Accepts a Duration directly or a string parsed via ParseDuration —
// the same contract as Options.Duration.
//
//	timeout := core.ParseDuration("30s").Duration()
func (r Result) Duration() Duration {
	if !r.OK {
		return 0
	}
	switch v := r.Value.(type) {
	case Duration:
		return v
	case string:
		if d := ParseDuration(v); d.OK {
			if dur, ok := d.Value.(Duration); ok {
				return dur
			}
		}
	}
	return 0
}

// Bytes retrieves a []byte Value, nil when failed or not []byte.
// Returns the backing slice directly — no defensive copy (hot I/O
// path; callers that mutate must copy).
//
//	body := core.HTTPGet(url).Bytes()
func (r Result) Bytes() []byte {
	if !r.OK {
		return nil
	}
	b, _ := r.Value.([]byte)
	return b
}

// Cast extracts a typed value from a Result. Returns (zero, false) when
// the Result is not OK or Value isn't assignable to T. Single
// expression replacing the (Result.OK check + type assertion) pair.
//
//	cfg, ok := core.Cast[*Config](core.JSONUnmarshal(data, &Config{}))
//	if !ok { return r }
//	if user, ok := core.Cast[*User](svc.Get(id)); ok { use(user) }
func Cast[T any](r Result) (T, bool) {
	var zero T
	if !r.OK {
		return zero, false
	}
	v, ok := r.Value.(T)
	if !ok {
		return zero, false
	}
	return v, true
}

// MustCast extracts a typed value from a Result. Panics with the
// underlying error when the Result is not OK or when Value isn't
// assignable to T. Use for fast-fail paths — init, test setup,
// must-have config — where the type-assertion is part of the contract.
//
//	cfg := core.MustCast[*Config](core.JSONUnmarshal(data, &Config{}))
//	dir := core.MustCast[string](core.PathAbs("."))
func MustCast[T any](r Result) T {
	var zero T
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			panic(err)
		}
		panic(r.Value)
	}
	v, ok := r.Value.(T)
	if !ok {
		panic(E("MustCast", Sprintf("Value is %T, not %T", r.Value, zero), nil))
	}
	return v
}

// Try runs fn and converts its outcome into a Result. A nil error or
// a returned value sets OK=true; a returned error or a panic sets
// OK=false. Bridges legacy code that panics or returns (T, error).
//
//	r := core.Try(func() any {
//	    return riskyParse(input)  // may panic
//	})
//	if !r.OK { core.Error("parse failed", "err", r.Error()) }
func Try(fn func() any) (r Result) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				r = Result{Value: err, OK: false}
				return
			}
			r = Result{Value: E("Try", "panic recovered", nil), OK: false}
		}
	}()
	v := fn()
	if err, ok := v.(error); ok && err != nil {
		return Result{Value: err, OK: false}
	}
	return Result{Value: v, OK: true}
}
