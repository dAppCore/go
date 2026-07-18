// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

// ExampleReturnOK wraps a value in a successful typed Return.
func ExampleReturnOK() {
	r := ReturnOK("brain")
	Println(r.OK(), r.Value)
	// Output: true brain
}

// ExampleReturnFail wraps an error in a failed typed Return.
func ExampleReturnFail() {
	r := ReturnFail[string](NewError("agent offline"))
	Println(r.OK())
	// Output: false
}

// ExampleReturnFrom adapts a stdlib (value, error) pair.
func ExampleReturnFrom() {
	parse := func(s string) (int, error) {
		if s == "" {
			return 0, NewError("empty input")
		}
		return len(s), nil
	}
	Println(ReturnFrom(parse("brain")).Or(0))
	Println(ReturnFrom(parse("")).Or(-1))
	// Output:
	// 5
	// -1
}

// ExampleReturnOf lifts an untyped Result into a typed Return.
func ExampleReturnOf() {
	r := ReturnOf[string](Ok("brain"))
	Println(r.Value)
	// Output: brain
}

// ExampleReturnTry converts a panic into a failed Return.
func ExampleReturnTry() {
	r := ReturnTry(func() int { panic(NewError("parser exploded")) })
	Println(r.OK())
	// Output: false
}

// ExampleReturn_OK reports success — Err is the only discriminant.
func ExampleReturn_OK() {
	Println(ReturnOK(0).OK())
	// Output: true
}

// ExampleReturn_Or returns the value or a typed fallback — no assertion.
func ExampleReturn_Or() {
	guest := "guest"
	Println(ReturnFail[string](NewError("no session")).Or(guest))
	// Output: guest
}

// ExampleReturn_Must unwraps for fast-fail paths.
func ExampleReturn_Must() {
	Println(ReturnOK(8080).Must())
	// Output: 8080
}

// ExampleReturn_Code exposes the stable error code, as Result.Code does.
func ExampleReturn_Code() {
	r := ReturnFail[int](NewCode("fs.notfound", "no such file"))
	Println(r.Code())
	// Output: fs.notfound
}

// ExampleReturn_Result erases the type back to the universal bus shape.
func ExampleReturn_Result() {
	r := ReturnOK("done").Result()
	Println(r.OK, r.String())
	// Output: true done
}

// ExampleReturn_Log opts a failure into the package log path and chains.
func ExampleReturn_Log() {
	r := ReturnOK("quiet").Log("svc.Op", "never logged — success is silent")
	Println(r.Value)
	// Output: quiet
}
