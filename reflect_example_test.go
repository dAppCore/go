package core_test

import (
	. "dappco.re/go"
)

// ExampleTypeOf reads a type through `TypeOf` for runtime inspection. Reflection stays
// behind a narrow core surface for rare inspection code.
func ExampleTypeOf() {
	t := TypeOf(42)
	Println(t.Kind())
	// Output: int
}

// ExampleValueOf reads a value wrapper through `ValueOf` for runtime inspection.
// Reflection stays behind a narrow core surface for rare inspection code.
func ExampleValueOf() {
	v := ValueOf("hello")
	Println(v.String())
	// Output: hello
}

// ExampleDeepEqual compares nested values through `DeepEqual` for runtime inspection.
// Reflection stays behind a narrow core surface for rare inspection code.
func ExampleDeepEqual() {
	Println(DeepEqual([]int{1, 2, 3}, []int{1, 2, 3}))
	// Output: true
}

// ExampleKind reads a reflection kind through `Kind` for runtime inspection. Reflection
// stays behind a narrow core surface for rare inspection code.
func ExampleKind() {
	k := TypeOf("x").Kind()
	Println(k == KindString)
	// Output: true
}

// ExampleTypeFor names a type at compile time through `TypeFor`, with no
// nil-pointer dance. Reflection stays behind a narrow core surface for
// rare inspection code.
func ExampleTypeFor() {
	Println(TypeFor[string]().Kind())
	// Output: string
}

// ExampleNewValue allocates a zeroed value of a reflected type through
// `NewValue` — the reflective new(T).
func ExampleNewValue() {
	ptr := NewValue(TypeFor[int]())
	ptr.Elem().SetInt(7)
	Println(ptr.Elem().Int())
	// Output: 7
}
