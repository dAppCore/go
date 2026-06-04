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
	AssertFalse(t, DeepEqual([]string{"agent"}, []string{"operator"}))
}

func TestReflect_DeepEqual_Ugly(t *T) {
	AssertFalse(t, DeepEqual([]string(nil), []string{}))
}

func TestReflect_TypeOf_Good(t *T) {
	AssertEqual(t, KindString, TypeOf("agent").Kind())
}

func TestReflect_TypeOf_Bad(t *T) {
	AssertNil(t, TypeOf(nil))
}

func TestReflect_TypeOf_Ugly(t *T) {
	var c *Core
	AssertEqual(t, KindPointer, TypeOf(c).Kind())
}

func TestReflect_ValueOf_Good(t *T) {
	AssertEqual(t, int64(42), ValueOf(42).Int())
}

func TestReflect_ValueOf_Bad(t *T) {
	AssertEqual(t, KindInvalid, ValueOf(nil).Kind())
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
}

func TestReflect_KindUintptr_Bad(t *T) {
	AssertNotEqual(t, KindUintptr, TypeOf(42).Kind())
}

func TestReflect_KindUintptr_Ugly(t *T) {
	// UnsafePointer is a distinct kind, not Uintptr.
	AssertNotEqual(t, KindUintptr, KindUnsafePointer)
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
	// Indexing past the field count panics.
	AssertPanics(t, func() { _ = TypeFor[reflectTagged]().Field(5) })
}

func TestReflect_StructField_Ugly(t *T) {
	AssertEqual(t, KindString, TypeFor[reflectTagged]().Field(0).Type.Kind())
}

func TestReflect_TypeFor_Good(t *T) {
	AssertEqual(t, KindString, TypeFor[string]().Kind())
}

func TestReflect_TypeFor_Bad(t *T) {
	AssertNotEqual(t, TypeFor[int](), TypeFor[string]())
}

func TestReflect_TypeFor_Ugly(t *T) {
	// TypeFor[any] yields a nil-interface element type.
	AssertEqual(t, KindInterface, TypeFor[any]().Kind())
}

func TestReflect_NewValue_Good(t *T) {
	ptr := NewValue(TypeFor[int]())
	AssertEqual(t, KindPointer, ptr.Kind())
	AssertEqual(t, int64(0), ptr.Elem().Int())
}

func TestReflect_NewValue_Bad(t *T) {
	AssertPanics(t, func() { _ = NewValue(nil) })
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
	AssertPanics(t, func() { _ = MakeMap(TypeFor[int]()) })
}

func TestReflect_MakeMap_Ugly(t *T) {
	m := MakeMap(TypeFor[map[string]int]())
	m.SetMapIndex(ValueOf("k"), ValueOf(9))
	AssertEqual(t, 1, m.Len())
}

func TestReflect_MakeMapWithSize_Good(t *T) {
	m := MakeMapWithSize(TypeFor[map[string]int](), 16)
	AssertEqual(t, 0, m.Len())
}

func TestReflect_MakeMapWithSize_Bad(t *T) {
	AssertPanics(t, func() { _ = MakeMapWithSize(TypeFor[int](), 4) })
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
