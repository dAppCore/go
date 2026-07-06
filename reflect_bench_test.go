// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the reflection primitives in reflect.go.
// Per AX-11 — TypeOf / ValueOf / DeepEqual / Zero are on every
// serialiser entry point, every contract registry lookup, every test
// AssertEqual comparison. Reflect is slow by nature; the bench harness
// locks in the contract floor so a Core reroute can't silently slow it.
//
// Run:    go test -bench='BenchmarkReflect' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	reflectSinkType  Type
	reflectSinkValue Value
	reflectSinkBool  bool
	reflectSinkKind  Kind
	reflectSinkInt   int
)

// Fixtures
type reflectStruct struct {
	Name string
	Port int
}

var (
	reflectFixSmall  = 42
	reflectFixString = "homelab"
	reflectFixStruct = reflectStruct{Name: "agent", Port: 9000}
	reflectFixSlice  = []int{1, 2, 3, 4, 5, 6, 7, 8}
	reflectFixMap    = map[string]int{"a": 1, "b": 2, "c": 3}
)

// --- TypeOf / ValueOf ---

func BenchmarkReflect_TypeOf_Int(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkType = TypeOf(reflectFixSmall)
	}
}

func BenchmarkReflect_TypeOf_Struct(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkType = TypeOf(reflectFixStruct)
	}
}

func BenchmarkReflect_ValueOf_Int(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = ValueOf(reflectFixSmall)
	}
}

func BenchmarkReflect_ValueOf_Struct(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = ValueOf(reflectFixStruct)
	}
}

// --- Kind inspection ---

func BenchmarkReflect_Kind(b *B) {
	t := TypeOf(reflectFixStruct)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkKind = t.Kind()
	}
}

// --- DeepEqual ---

func BenchmarkReflect_DeepEqual_Ints(b *B) {
	a, c := 42, 42
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkBool = DeepEqual(a, c)
	}
}

func BenchmarkReflect_DeepEqual_Slices(b *B) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	c := []int{1, 2, 3, 4, 5, 6, 7, 8}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkBool = DeepEqual(a, c)
	}
}

func BenchmarkReflect_DeepEqual_Map(b *B) {
	a := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map[string]int{"a": 1, "b": 2, "c": 3}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkBool = DeepEqual(a, c)
	}
}

func BenchmarkReflect_DeepEqual_Mismatch(b *B) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	c := []int{1, 2, 3, 4, 5, 6, 7, 9}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkBool = DeepEqual(a, c)
	}
}

// --- Zero ---

func BenchmarkReflect_Zero_Struct(b *B) {
	t := TypeOf(reflectFixStruct)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = Zero(t)
	}
}

// --- TypeFor / NewValue / MakeSlice / MakeMap / CopyValue ---

func BenchmarkReflect_TypeFor_Struct(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkType = TypeFor[reflectStruct]()
	}
}

func BenchmarkReflect_NewValue_Struct(b *B) {
	t := TypeFor[reflectStruct]()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = NewValue(t)
	}
}

func BenchmarkReflect_MakeSlice_Int(b *B) {
	t := TypeFor[[]int]()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = MakeSlice(t, 0, 8)
	}
}

func BenchmarkReflect_MakeMap_StringInt(b *B) {
	t := TypeFor[map[string]int]()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = MakeMap(t)
	}
}

func BenchmarkReflect_CopyValue_Int(b *B) {
	src := ValueOf(reflectFixSlice)
	dst := MakeSlice(TypeFor[[]int](), len(reflectFixSlice), len(reflectFixSlice))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkInt = CopyValue(dst, src)
	}
}

func BenchmarkReflect_MakeMapWithSize(b *B) {
	t := TypeFor[map[string]int]()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = MakeMapWithSize(t, 8)
	}
}

func BenchmarkReflect_MakeFunc(b *B) {
	t := TypeFor[func()]()
	fn := func(args []Value) []Value { return nil }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reflectSinkValue = MakeFunc(t, fn)
	}
}
