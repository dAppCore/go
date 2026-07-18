package core_test

import (
	. "dappco.re/go"
)

func TestReflect_DeepEqual_Good(t *T) {
	left := map[string]int{"agent": 1, "health": 2}
	right := map[string]int{"health": 2, "agent": 1}
	AssertTrue(t, DeepEqual(left, right))
}

func TestReflect_DeepEqual_Bad(t *T) {
	// Same-typed slices with differing contents are not equal.
	AssertFalse(t, DeepEqual([]string{"agent"}, []string{"operator"}))
	// Differing lengths are not equal.
	AssertFalse(t, DeepEqual([]int{1, 2}, []int{1, 2, 3}))
	// Same keys, different values.
	AssertFalse(t, DeepEqual(map[string]int{"a": 1}, map[string]int{"a": 2}))
	// Distinct dynamic types are never deeply equal, even with equal magnitude.
	AssertFalse(t, DeepEqual(1, int64(1)))
	// Differing struct field values.
	AssertFalse(t, DeepEqual(reflectTagged{Name: "x"}, reflectTagged{Name: "y"}))
}

func TestReflect_DeepEqual_Ugly(t *T) {
	// A nil slice and an empty non-nil slice are NOT deeply equal.
	AssertFalse(t, DeepEqual([]string(nil), []string{}))
	// Likewise a nil map vs an empty non-nil map.
	AssertFalse(t, DeepEqual(map[string]int(nil), map[string]int{}))
	// Two untyped nils short-circuit to equal (x == y).
	AssertTrue(t, DeepEqual(nil, nil))
	// nil against a concrete value is not equal.
	AssertFalse(t, DeepEqual(nil, 0))
	// Distinct non-nil func values are never deeply equal.
	AssertFalse(t, DeepEqual(func() {}, func() {}))
}

func TestReflect_TypeOf_Good(t *T) {
	// Kind classification across the common scalar and composite kinds.
	AssertEqual(t, KindString, TypeOf("agent").Kind())
	AssertEqual(t, "string", TypeOf("agent").Name())
	AssertEqual(t, KindInt, TypeOf(42).Kind())
	AssertEqual(t, KindFloat64, TypeOf(3.14).Kind())
	AssertEqual(t, KindBool, TypeOf(true).Kind())
	AssertEqual(t, KindSlice, TypeOf([]int{}).Kind())
	AssertEqual(t, KindMap, TypeOf(map[string]int{}).Kind())
	// A struct value reports its declared name and field count.
	st := TypeOf(reflectTagged{})
	AssertEqual(t, KindStruct, st.Kind())
	AssertEqual(t, "reflectTagged", st.Name())
	AssertEqual(t, 1, st.NumField())
}

func TestReflect_TypeOf_Bad(t *T) {
	// A nil interface value yields a nil Type.
	AssertNil(t, TypeOf(nil))
	var e error
	AssertNil(t, TypeOf(e))
	// But a typed nil pointer carries its type — Type is non-nil.
	var p *int
	AssertNotNil(t, TypeOf(p))
	AssertEqual(t, KindPointer, TypeOf(p).Kind())
}

func TestReflect_TypeOf_Ugly(t *T) {
	var c *Core
	pt := TypeOf(c)
	AssertEqual(t, KindPointer, pt.Kind())
	// The pointed-to type is the Core struct.
	AssertEqual(t, KindStruct, pt.Elem().Kind())
	AssertEqual(t, "Core", pt.Elem().Name())
	// A pointer to a scalar resolves its element kind too.
	var n *int
	AssertEqual(t, KindInt, TypeOf(n).Elem().Kind())
}

func TestReflect_ValueOf_Good(t *T) {
	// Int extraction and kind.
	AssertEqual(t, int64(42), ValueOf(42).Int())
	AssertEqual(t, KindInt, ValueOf(42).Kind())
	// String, bool and float accessors mirror the concrete value.
	AssertEqual(t, "agent", ValueOf("agent").String())
	AssertTrue(t, ValueOf(true).Bool())
	AssertEqual(t, 3.5, ValueOf(3.5).Float())
	// Composite values expose length.
	AssertEqual(t, 3, ValueOf([]int{1, 2, 3}).Len())
	// Interface() round-trips the original value.
	AssertEqual(t, 42, ValueOf(42).Interface())
}

func TestReflect_ValueOf_Bad(t *T) {
	// A nil interface produces the zero Value: invalid and not valid.
	AssertEqual(t, KindInvalid, ValueOf(nil).Kind())
	AssertFalse(t, ValueOf(nil).IsValid())
	// A typed nil pointer yields a valid Value that is nil.
	var p *int
	v := ValueOf(p)
	AssertTrue(t, v.IsValid())
	AssertEqual(t, KindPointer, v.Kind())
	AssertTrue(t, v.IsNil())
}

func TestReflect_ValueOf_Ugly(t *T) {
	var c *Core
	v := ValueOf(c)
	AssertEqual(t, KindPointer, v.Kind())
	AssertTrue(t, v.IsNil())
}

func TestReflect_Zero_Good(t *T) {
	z := Zero(TypeOf(42))
	AssertEqual(t, 0, z.Interface())
	AssertEqual(t, KindInt, z.Kind())
	AssertEqual(t, int64(0), z.Int())
	// Zero of string and bool types.
	AssertEqual(t, "", Zero(TypeFor[string]()).String())
	AssertFalse(t, Zero(TypeFor[bool]()).Bool())
	// Zero of a slice type is a nil, empty slice.
	zs := Zero(TypeFor[[]int]())
	AssertTrue(t, zs.IsNil())
	AssertEqual(t, 0, zs.Len())
	// Zero of a struct type has zeroed fields.
	AssertEqual(t, "", Zero(TypeFor[reflectTagged]()).Field(0).String())
}

func TestReflect_Zero_Bad(t *T) {
	AssertPanics(t, func() {
		_ = Zero(nil)
	})
}

func TestReflect_Zero_Ugly(t *T) {
	z := Zero(TypeOf((*Core)(nil)))
	AssertEqual(t, KindPointer, z.Kind())
	AssertTrue(t, z.IsNil())
}

func TestReflect_KindUintptr_Good(t *T) {
	var p uintptr
	AssertEqual(t, KindUintptr, TypeOf(p).Kind())
	// The compile-time path agrees with the runtime path.
	AssertEqual(t, KindUintptr, TypeFor[uintptr]().Kind())
	AssertEqual(t, "uintptr", TypeFor[uintptr]().Name())
	// The Kind constant stringifies to "uintptr".
	AssertEqual(t, "uintptr", KindUintptr.String())
	// A uintptr Value reads back through Uint.
	v := ValueOf(uintptr(7))
	AssertEqual(t, KindUintptr, v.Kind())
	AssertEqual(t, uint64(7), v.Uint())
}

func TestReflect_KindUintptr_Bad(t *T) {
	// A plain int is Int, not Uintptr.
	AssertNotEqual(t, KindUintptr, TypeOf(42).Kind())
	AssertEqual(t, KindInt, TypeOf(42).Kind())
	// The unsigned kinds are still distinct from Uintptr.
	AssertNotEqual(t, KindUintptr, TypeOf(uint(1)).Kind())
	AssertEqual(t, KindUint, TypeOf(uint(1)).Kind())
	// A pointer is Pointer, not Uintptr.
	AssertNotEqual(t, KindUintptr, TypeOf((*int)(nil)).Kind())
}

func TestReflect_KindUintptr_Ugly(t *T) {
	// UnsafePointer is a distinct kind, not Uintptr.
	AssertNotEqual(t, KindUintptr, KindUnsafePointer)
	// Nor is it the same as Pointer or Uint.
	AssertNotEqual(t, KindUintptr, KindPointer)
	AssertNotEqual(t, KindUintptr, KindUint)
	// The two pointer-ish kinds stringify distinctly.
	AssertEqual(t, "uintptr", KindUintptr.String())
	AssertEqual(t, "unsafe.Pointer", KindUnsafePointer.String())
}

type reflectTagged struct {
	Name string `json:"name"`
}

func TestReflect_StructField_Good(t *T) {
	var f StructField = TypeFor[reflectTagged]().Field(0)
	AssertEqual(t, "Name", f.Name)
	AssertEqual(t, "name", f.Tag.Get("json"))
}

func TestReflect_StructField_Bad(t *T) {
	typ := TypeFor[reflectTagged]()
	// A valid index yields a StructField; out-of-range / non-struct panic.
	var _ StructField = typ.Field(0)
	AssertEqual(t, 1, typ.NumField())
	AssertPanics(t, func() { _ = typ.Field(5) })            // past the field count
	AssertPanics(t, func() { _ = typ.Field(-1) })           // negative index
	AssertPanics(t, func() { _ = TypeFor[int]().Field(0) }) // non-struct type
}

func TestReflect_StructField_Ugly(t *T) {
	var f StructField = TypeFor[reflectTagged]().Field(0)
	AssertEqual(t, KindString, f.Type.Kind())
	AssertEqual(t, "Name", f.Name)
	// Struct tags are parsed: present key returns its value, absent key empty.
	AssertEqual(t, "name", f.Tag.Get("json"))
	AssertEqual(t, "", f.Tag.Get("missing"))
	// The field is exported and not embedded.
	AssertTrue(t, f.IsExported())
	AssertFalse(t, f.Anonymous)
	// An exported field has an empty PkgPath.
	AssertEmpty(t, f.PkgPath)
}

func TestReflect_TypeFor_Good(t *T) {
	AssertEqual(t, KindString, TypeFor[string]().Kind())
	AssertEqual(t, KindInt, TypeFor[int]().Kind())
	AssertEqual(t, KindSlice, TypeFor[[]byte]().Kind())
	AssertEqual(t, KindMap, TypeFor[map[string]int]().Kind())
	AssertEqual(t, KindStruct, TypeFor[reflectTagged]().Kind())
	// A pointer type resolves its element kind.
	AssertEqual(t, KindPointer, TypeFor[*int]().Kind())
	AssertEqual(t, KindInt, TypeFor[*int]().Elem().Kind())
	// TypeFor[T] agrees with TypeOf on a value of T.
	AssertEqual(t, TypeOf(""), TypeFor[string]())
}

func TestReflect_TypeFor_Bad(t *T) {
	AssertNotEqual(t, TypeFor[int](), TypeFor[string]())
	// Differently sized integer types are distinct.
	AssertNotEqual(t, TypeFor[int](), TypeFor[int64]())
	// Element type makes slice types distinct.
	AssertNotEqual(t, TypeFor[[]int](), TypeFor[[]string]())
	// The same type is interned: repeated calls are equal.
	AssertEqual(t, TypeFor[int](), TypeFor[int]())
}

func TestReflect_TypeFor_Ugly(t *T) {
	// TypeFor[any] yields the empty interface type — zero methods.
	AssertEqual(t, KindInterface, TypeFor[any]().Kind())
	AssertEqual(t, 0, TypeFor[any]().NumMethod())
	// A non-empty interface (error) is also Interface kind but has a method.
	AssertEqual(t, KindInterface, TypeFor[error]().Kind())
	AssertEqual(t, 1, TypeFor[error]().NumMethod())
	// The empty and non-empty interface types are distinct.
	AssertNotEqual(t, TypeFor[any](), TypeFor[error]())
}

func TestReflect_NewValue_Good(t *T) {
	ptr := NewValue(TypeFor[int]())
	AssertEqual(t, KindPointer, ptr.Kind())
	AssertEqual(t, int64(0), ptr.Elem().Int())
}

func TestReflect_NewValue_Bad(t *T) {
	// A nil type panics — there is nothing to allocate.
	AssertPanics(t, func() { _ = NewValue(nil) })
	// A valid type does not panic and yields a non-nil pointer to a zero value.
	AssertNotPanics(t, func() { _ = NewValue(TypeFor[int]()) })
	ptr := NewValue(TypeFor[string]())
	AssertEqual(t, KindPointer, ptr.Kind())
	AssertFalse(t, ptr.IsNil())
	AssertEqual(t, "", ptr.Elem().String())
}

func TestReflect_NewValue_Ugly(t *T) {
	// The new value is addressable and settable through its pointer.
	ptr := NewValue(TypeFor[int]())
	ptr.Elem().SetInt(7)
	AssertEqual(t, int64(7), ptr.Elem().Int())
}

func TestReflect_MakeSlice_Good(t *T) {
	s := MakeSlice(TypeFor[[]int](), 0, 4)
	AssertEqual(t, 0, s.Len())
	AssertEqual(t, 4, s.Cap())
}

func TestReflect_MakeSlice_Bad(t *T) {
	// A non-slice type panics.
	AssertPanics(t, func() { _ = MakeSlice(TypeFor[int](), 0, 0) })
	// A map type is also not a slice.
	AssertPanics(t, func() { _ = MakeSlice(TypeFor[map[string]int](), 0, 0) })
	// len greater than cap panics.
	AssertPanics(t, func() { _ = MakeSlice(TypeFor[[]int](), 5, 2) })
	// A negative length panics.
	AssertPanics(t, func() { _ = MakeSlice(TypeFor[[]int](), -1, 0) })
}

func TestReflect_MakeSlice_Ugly(t *T) {
	s := MakeSlice(TypeFor[[]int](), 2, 2)
	AssertEqual(t, 2, s.Len())
	AssertEqual(t, int64(0), s.Index(0).Int())
}

func TestReflect_MakeMap_Good(t *T) {
	m := MakeMap(TypeFor[map[string]int]())
	AssertEqual(t, KindMap, m.Kind())
	AssertEqual(t, 0, m.Len())
}

func TestReflect_MakeMap_Bad(t *T) {
	// A scalar type is not a map.
	AssertPanics(t, func() { _ = MakeMap(TypeFor[int]()) })
	// A slice type is not a map either.
	AssertPanics(t, func() { _ = MakeMap(TypeFor[[]int]()) })
	// A genuine map type does not panic and starts empty.
	AssertNotPanics(t, func() { _ = MakeMap(TypeFor[map[int]string]()) })
	m := MakeMap(TypeFor[map[int]string]())
	AssertEqual(t, KindMap, m.Kind())
	AssertEqual(t, 0, m.Len())
}

func TestReflect_MakeMap_Ugly(t *T) {
	m := MakeMap(TypeFor[map[string]int]())
	m.SetMapIndex(ValueOf("k"), ValueOf(9))
	AssertEqual(t, 1, m.Len())
}

func TestReflect_MakeMapWithSize_Good(t *T) {
	m := MakeMapWithSize(TypeFor[map[string]int](), 16)
	AssertEqual(t, KindMap, m.Kind())
	// The size is a hint only: the map starts empty but is non-nil.
	AssertFalse(t, m.IsNil())
	AssertEqual(t, 0, m.Len())
	// The pre-sized map is immediately usable.
	m.SetMapIndex(ValueOf("k"), ValueOf(3))
	AssertEqual(t, 1, m.Len())
	AssertEqual(t, int64(3), m.MapIndex(ValueOf("k")).Int())
}

func TestReflect_MakeMapWithSize_Bad(t *T) {
	// A scalar type is not a map.
	AssertPanics(t, func() { _ = MakeMapWithSize(TypeFor[int](), 4) })
	// A slice type is not a map.
	AssertPanics(t, func() { _ = MakeMapWithSize(TypeFor[[]int](), 4) })
	// A real map type with a large hint does not panic.
	AssertNotPanics(t, func() { _ = MakeMapWithSize(TypeFor[map[string]int](), 1024) })
}

func TestReflect_MakeMapWithSize_Ugly(t *T) {
	// A zero size hint is valid; the map is still usable.
	m := MakeMapWithSize(TypeFor[map[string]int](), 0)
	m.SetMapIndex(ValueOf("k"), ValueOf(1))
	AssertEqual(t, 1, m.Len())
}

func TestReflect_CopyValue_Good(t *T) {
	src := ValueOf([]int{1, 2, 3})
	dst := MakeSlice(TypeFor[[]int](), 3, 3)
	n := CopyValue(dst, src)
	AssertEqual(t, 3, n)
	AssertEqual(t, int64(2), dst.Index(1).Int())
}

func TestReflect_CopyValue_Bad(t *T) {
	// Copying into a zero-length destination copies nothing.
	src := ValueOf([]int{1, 2, 3})
	dst := MakeSlice(TypeFor[[]int](), 0, 0)
	AssertEqual(t, 0, CopyValue(dst, src))
}

func TestReflect_CopyValue_Ugly(t *T) {
	// Copy stops at the shorter length (dst here).
	src := ValueOf([]int{1, 2, 3, 4})
	dst := MakeSlice(TypeFor[[]int](), 2, 2)
	AssertEqual(t, 2, CopyValue(dst, src))
}

func TestReflect_MakeFunc_Good(t *T) {
	// Synthesise func(int) int that doubles its argument.
	doubler := MakeFunc(TypeFor[func(int) int](), func(args []Value) []Value {
		return []Value{ValueOf(int(args[0].Int()) * 2)}
	})
	fn := doubler.Interface().(func(int) int)
	AssertEqual(t, 10, fn(5))
}

func TestReflect_MakeFunc_Bad(t *T) {
	// A non-func type panics.
	AssertPanics(t, func() {
		_ = MakeFunc(TypeFor[int](), func(args []Value) []Value { return nil })
	})
}

func TestReflect_MakeFunc_Ugly(t *T) {
	// A no-arg, no-result func is valid and callable.
	noop := MakeFunc(TypeFor[func()](), func(args []Value) []Value { return nil })
	fn := noop.Interface().(func())
	AssertNotPanics(t, fn)
}
