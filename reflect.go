// SPDX-License-Identifier: EUPL-1.2

// Reflection primitives — re-exports of Go's standard reflect package as
// Core types, so consumers never need to write `import "reflect"`.
// Reflection is a low-level tool; reach for it only when generic
// constraints and type switches genuinely cannot express the operation.
//
//	if core.TypeOf(v).Kind() == core.KindStruct { ... }
package core

import "reflect"

// Type is the runtime type descriptor returned by TypeOf.
//
//	t := core.TypeOf(opts)
//	name := t.Name()
type Type = reflect.Type

// Value is the runtime value handle returned by ValueOf.
//
//	v := core.ValueOf(42)
//	n := v.Int()
type Value = reflect.Value

// Kind classifies a Type into one of the basic Go kinds (Bool, Int,
// String, Struct, Map, Slice, etc.).
//
//	if core.TypeOf(opts).Kind() == core.KindStruct { ... }
type Kind = reflect.Kind

// Common Kind constants — re-exported from reflect so consumers can
// compare without importing reflect.
//
//	switch core.TypeOf(v).Kind() {
//	case core.KindString:
//	case core.KindInt, core.KindInt64:
//	}
const (
	KindInvalid       = reflect.Invalid
	KindBool          = reflect.Bool
	KindInt           = reflect.Int
	KindInt8          = reflect.Int8
	KindInt16         = reflect.Int16
	KindInt32         = reflect.Int32
	KindInt64         = reflect.Int64
	KindUint          = reflect.Uint
	KindUint8         = reflect.Uint8
	KindUint16        = reflect.Uint16
	KindUint32        = reflect.Uint32
	KindUint64        = reflect.Uint64
	KindFloat32       = reflect.Float32
	KindFloat64       = reflect.Float64
	KindComplex64     = reflect.Complex64
	KindComplex128    = reflect.Complex128
	KindArray         = reflect.Array
	KindChan          = reflect.Chan
	KindFunc          = reflect.Func
	KindInterface     = reflect.Interface
	KindMap           = reflect.Map
	KindPointer       = reflect.Pointer
	KindSlice         = reflect.Slice
	KindString        = reflect.String
	KindStruct        = reflect.Struct
	KindUintptr       = reflect.Uintptr
	KindUnsafePointer = reflect.UnsafePointer
)

// StructField is an alias for reflect.StructField — one field's
// descriptor (Name, Type, Tag, offset) returned by Type.Field during
// struct introspection.
//
//	f := core.TypeOf(opts).Field(0)
//	tag := f.Tag.Get("json")
type StructField = reflect.StructField

// TypeOf returns the runtime type of v. Returns nil if v is a nil
// interface value.
//
//	t := core.TypeOf(opts)
//	if t.Kind() == core.KindStruct { ... }
func TypeOf(v any) Type {
	return reflect.TypeOf(v)
}

// ValueOf returns a Value initialised to the concrete value stored in v.
// Returns the zero Value if v is a nil interface value.
//
//	val := core.ValueOf(42)
//	n := val.Int()  // 42
func ValueOf(v any) Value {
	return reflect.ValueOf(v)
}

// DeepEqual reports whether x and y are deeply equal — recursively
// comparing slices, maps, structs, etc. Differs from == in that it
// follows pointers and compares unexported fields.
//
//	core.DeepEqual([]int{1, 2, 3}, []int{1, 2, 3})  // true
//	core.DeepEqual(map[string]int{"a": 1}, map[string]int{"a": 1})  // true
func DeepEqual(x, y any) bool {
	return reflect.DeepEqual(x, y)
}

// Zero returns a Value representing the zero value for type t. Used
// for default-value comparison without an explicit zero literal.
//
//	t := core.TypeOf((*MyStruct)(nil))
//	zero := core.Zero(t.Elem()).Interface()
func Zero(t Type) Value {
	return reflect.Zero(t)
}

// TypeFor returns the Type for the compile-time type T. The type-safe
// replacement for TypeOf((*T)(nil)).Elem() — no nil-pointer dance, no
// runtime value needed.
//
//	t := core.TypeFor[MyStruct]()
//	if t.Kind() == core.KindStruct { ... }
func TypeFor[T any]() Type {
	return reflect.TypeFor[T]()
}

// NewValue returns a Value representing a pointer to a new zero value
// of type t — the reflective equivalent of new(T). Named NewValue (not
// New) because core.New is the framework constructor.
//
//	ptr := core.NewValue(core.TypeFor[Config]())  // *Config, zeroed
//	cfg := ptr.Elem().Interface().(Config)
func NewValue(t Type) Value {
	return reflect.New(t)
}

// MakeSlice returns a Value representing a new slice of element type's
// slice t with the given length and capacity. t must have Kind Slice;
// callers control that, so this stays infallible (panics on a non-slice
// type, matching the stdlib contract).
//
//	s := core.MakeSlice(core.TypeFor[[]int](), 0, 8)
func MakeSlice(t Type, len, cap int) Value {
	return reflect.MakeSlice(t, len, cap)
}

// MakeMap returns a Value representing a new empty map of map type t.
//
//	m := core.MakeMap(core.TypeFor[map[string]int]())
func MakeMap(t Type) Value {
	return reflect.MakeMap(t)
}

// MakeMapWithSize returns a new empty map of type t pre-sized for about
// n entries — the reflective make(map, n) hint.
//
//	m := core.MakeMapWithSize(core.TypeFor[map[string]int](), 64)
func MakeMapWithSize(t Type, n int) Value {
	return reflect.MakeMapWithSize(t, n)
}

// CopyValue copies the contents of src into dst until dst is full or
// src is exhausted, returning the number of elements copied. Both must
// be slices (or dst an array) with assignable element types. Named
// CopyValue (not Copy) because core.Copy is the io stream copier.
//
//	n := core.CopyValue(dstVal, srcVal)
func CopyValue(dst, src Value) int {
	return reflect.Copy(dst, src)
}

// MakeFunc returns a new function Value of type t whose body calls fn
// with the in-arguments and returns fn's results. Used to synthesise
// functions (proxies, generic adapters) at runtime — reach for it only
// when a closure over a concrete signature genuinely cannot.
//
//	fn := core.MakeFunc(t, func(args []core.Value) []core.Value { ... })
func MakeFunc(t Type, fn func(args []Value) (results []Value)) Value {
	return reflect.MakeFunc(t, fn)
}
