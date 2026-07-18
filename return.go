// SPDX-License-Identifier: EUPL-1.2

// Return[T] — the typed twin of Result for leaf functions with one
// concrete return type. Flat: a direct Value OR an error, with Err as
// the discriminant. Result stays the universal bus type; Return[T]
// gives compile-time typing at call boundaries so consumers never
// type-assert.
//
// The law (see PLAN-v0.12.0.md W1-3):
//
//  1. Anything crossing a type-erasing boundary — IPC, registries,
//     Actions, Tasks, RegistryOf — returns Result. No exceptions.
//
//  2. Leaf functions with one concrete struct-shaped return may
//     declare Return[T].
//
//  3. Scalars use Result and its typed getters, never lifted into
//     Return. One symbol never offers both shapes.
//
//     func UserByName(name string) core.Return[*User] {
//     return core.ReturnFrom(lookup(name))
//     }
//     u := UserByName("snider").Log("svc.Auth", "lookup failed").Or(guest)
package core

// Return is the typed single-value return: Value OR Err. A nil Err
// means success. The zero value is a successful Return carrying T's
// zero value.
//
//	r := core.ReturnOK(&Config{})
//	if r.OK() { cfg := r.Value }
type Return[T any] struct {
	Value T
	Err   error
}

// ReturnOK wraps v in a successful Return — the typed Ok.
//
//	return core.ReturnOK(parsed)
func ReturnOK[T any](v T) Return[T] {
	return Return[T]{Value: v}
}

// ReturnFail wraps err in a failed Return — the typed Fail.
//
//	return core.ReturnFail[*Config](core.E("config.Load", "no file", nil))
func ReturnFail[T any](err error) Return[T] {
	return Return[T]{Err: err}
}

// ReturnFrom adapts a stdlib (value, error) pair — the typed ResultOf.
//
//	return core.ReturnFrom(os.ReadFile(path))
func ReturnFrom[T any](v T, err error) Return[T] {
	return Return[T]{Value: v, Err: err}
}

// ReturnOf lifts a Result into a typed Return. A failed Result carries
// its error across (never nil — see Result.Err); an OK Result whose
// Value isn't T fails with operation "core.ReturnOf".
//
//	r := core.ReturnOf[*User](c.QUERY(userQuery{ID: id}))
func ReturnOf[T any](r Result) Return[T] {
	if !r.OK {
		return Return[T]{Err: r.Err()}
	}
	v, ok := r.Value.(T)
	if !ok {
		var zero T
		return Return[T]{Err: E("core.ReturnOf", Sprintf("Value is %T, not %T", r.Value, zero), nil)}
	}
	return Return[T]{Value: v}
}

// ReturnTry runs fn and converts a panic into a failed Return — the
// typed Try. Bridges legacy code that panics.
//
//	r := core.ReturnTry(func() *Config { return riskyParse(input) })
//	if !r.OK() { return r }
func ReturnTry[T any](fn func() T) (r Return[T]) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				r = Return[T]{Err: err}
				return
			}
			r = Return[T]{Err: E("core.ReturnTry", Sprint("panic recovered: ", rec), nil)}
		}
	}()
	return Return[T]{Value: fn()}
}

// OK reports success — Err is nil.
//
//	if !r.OK() { return r }
func (r Return[T]) OK() bool {
	return r.Err == nil
}

// Or returns Value on success, fallback on failure — typed, no
// assertion ever.
//
//	u := UserByName(name).Or(guest)
func (r Return[T]) Or(fallback T) T {
	if r.Err != nil {
		return fallback
	}
	return r.Value
}

// Must returns Value on success and panics with Err on failure. For
// fast-fail paths — init, test setup, must-have config.
//
//	cfg := core.ReturnOf[*Config](r).Must()
func (r Return[T]) Must() T {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.Value
}

// Code returns the stable error code when Err is a *core.Err with a
// Code populated, "" otherwise — the same contract as Result.Code.
//
//	if r.Code() == "fs.notfound" { firstRun() }
func (r Return[T]) Code() string {
	if e, ok := r.Err.(*Err); ok {
		return e.Code
	}
	return ""
}

// Result erases the type back to the universal bus shape: a failed
// Return carries Err as the failure Value, a successful one carries
// Value.
//
//	c.ACTION(TaskDone{Result: r.Result()})
func (r Return[T]) Result() Result {
	if r.Err != nil {
		return Result{Value: r.Err, OK: false}
	}
	return Result{Value: r.Value, OK: true}
}

// Log reports the failure through the package log path and returns the
// Return unchanged for chaining — logging handled by Core, opted into
// per call site. No-op on success, so expected-failure branches stay
// silent unless the caller asks.
//
//	u := UserByName(n).Log("svc.Auth", "lookup failed").Or(guest)
func (r Return[T]) Log(op, msg string) Return[T] {
	if r.Err != nil {
		Error(msg, "op", op, "err", r.Err)
	}
	return r
}
