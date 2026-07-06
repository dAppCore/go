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

// ExampleZero builds the zero Value of a type through `Zero`.
func ExampleZero() {
	Println(Zero(TypeFor[int]()).Int())
	// Output: 0
}

// ExampleMakeSlice allocates a typed slice Value through `MakeSlice`.
func ExampleMakeSlice() {
	s := MakeSlice(TypeFor[[]int](), 0, 4)
	Println(s.Cap())
	// Output: 4
}

// ExampleMakeMap allocates a typed map Value through `MakeMap`.
func ExampleMakeMap() {
	Println(MakeMap(TypeFor[map[string]int]()).Len())
	// Output: 0
}

// ExampleMakeMapWithSize allocates a sized map Value through `MakeMapWithSize`.
func ExampleMakeMapWithSize() {
	Println(MakeMapWithSize(TypeFor[map[string]int](), 8).Len())
	// Output: 0
}

// ExampleCopyValue copies between slice Values through `CopyValue`.
func ExampleCopyValue() {
	dst := MakeSlice(TypeFor[[]int](), 3, 3)
	Println(CopyValue(dst, dst))
	// Output: 3
}

// ExampleMakeFunc builds a callable Value at runtime through `MakeFunc`.
func ExampleMakeFunc() {
	fn := MakeFunc(TypeFor[func()](), func(args []Value) []Value { return nil })
	_ = fn // an invokable reflect.Value wrapping the synthesised function
}
