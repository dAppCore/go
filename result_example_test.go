package core_test

import (
	. "dappco.re/go"
)

// ExampleResult_Must unwraps a successful Result through `Result.Must` for Result-based
// control flow. Success, fallback, casting, and error inspection all use the same Result
// shape.
func ExampleResult_Must() {
	v := (Result{Value: 42, OK: true}).Must()
	Println(v)
	// Output: 42
}

// ExampleResult_Or falls back from a failed Result through `Result.Or` for Result-based
// control flow. Success, fallback, casting, and error inspection all use the same Result
// shape.
func ExampleResult_Or() {
	v := (Result{OK: false}).Or("fallback")
	Println(v)
	// Output: fallback
}

// ExampleResult_Error writes or renders an error through `Result.Error` for Result-based
// control flow. Success, fallback, casting, and error inspection all use the same Result
// shape.
func ExampleResult_Error() {
	r := Result{Value: NewError("bad config"), OK: false}
	Println(r.Error())
	// Output: bad config
}

// ExampleResult_Code reads a Result error code through `Result.Code` for Result-based
// control flow. Success, fallback, casting, and error inspection all use the same Result
// shape.
func ExampleResult_Code() {
	r := Result{Value: NewCode("fs.notfound", "missing file"), OK: false}
	Println(r.Code())
	// Output: fs.notfound
}

// ExampleCast casts a Result value through `Cast` for Result-based control flow. Success,
// fallback, casting, and error inspection all use the same Result shape.
func ExampleCast() {
	r := Result{Value: "hello", OK: true}
	s, ok := Cast[string](r)
	Println(s, ok)
	// Output: hello true
}

// ExampleTry converts panic-prone work into a Result through `Try` for Result-based
// control flow. Success, fallback, casting, and error inspection all use the same Result
// shape.
func ExampleTry() {
	r := Try(func() any {
		return 42
	})
	Println(r.OK, r.Value)
	// Output: true 42
}

// ExampleOk wraps a value in a successful Result through `Ok` for Result-based control
// flow. Success, fallback, casting, and error inspection all use the same Result shape.
func ExampleOk() {
	r := Ok(42)
	Println(r.OK, r.Value)
	// Output: true 42
}

// ExampleFail wraps an error in a failed Result through `Fail` for Result-based control
// flow. Success, fallback, casting, and error inspection all use the same Result shape.
func ExampleFail() {
	r := Fail(NewError("boom"))
	Println(r.OK)
	Println(r.Error())
	// Output:
	// false
	// boom
}

// ExampleResultOf adapts a (value, error) pair through `ResultOf` for Result-based control
// flow. Success, fallback, casting, and error inspection all use the same Result shape.
func ExampleResultOf() {
	r := ResultOf("data", nil)
	Println(r.OK, r.Value)
	// Output: true data
}

// ExampleMustCast unwraps and type-asserts a Result through `MustCast` for Result-based
// control flow. Success, fallback, casting, and error inspection all use the same Result shape.
func ExampleMustCast() {
	n := MustCast[int](Ok(42))
	Println(n)
	// Output: 42
}

// ExampleResult_String shows the Stringer-first contract: the value when
// OK, the error text when failed — %v of any Result reads well in logs.
func ExampleResult_String() {
	Println(Ok("brain").String())
	Println(Fail(NewError("agent offline")).String())
	// Output:
	// brain
	// agent offline
}

// ExampleResult_Int reads a typed int with the Options accessor contract.
func ExampleResult_Int() {
	Println(Ok(8080).Int())
	// Output: 8080
}

// ExampleResult_Bool reads a typed bool, false on failure or wrong type.
func ExampleResult_Bool() {
	Println(Ok(true).Bool())
	// Output: true
}

// ExampleResult_Float64 promotes int/int64/float32, as Options.Float64 does.
func ExampleResult_Float64() {
	Println(Ok(3).Float64())
	// Output: 3
}

// ExampleResult_Duration accepts a Duration or a ParseDuration string.
func ExampleResult_Duration() {
	Println(Ok("30s").Duration())
	// Output: 30s
}

// ExampleResult_Bytes reads a []byte Value directly — no copy, no assert.
func ExampleResult_Bytes() {
	Println(string(Ok([]byte("payload")).Bytes()))
	// Output: payload
}

// ExampleResult_Err returns the failure as an error — nil when OK, and
// never nil when failed.
func ExampleResult_Err() {
	Println(Ok("fine").Err() == nil)
	Println(Fail(NewError("agent offline")).Err() != nil)
	// Output:
	// true
	// true
}
