// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Service primitive in service.go.
// Per AX-11 — Service registration is one-shot per process, but the
// lookup path (Service / ServiceFor / MustServiceFor) is on every
// runtime call that resolves a registered dependency. The Wails RPC
// dispatcher uses Services() to enumerate the bound surface on boot.
//
// Run:    go test -bench='BenchmarkService' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	svcSinkResult  Result
	svcSinkBool    bool
	svcSinkStrings []string
	svcSinkAny     any
)

type benchService struct {
	id int
}

// serviceFixture builds a Core with one named service ready to look up.
func serviceFixture() (*Core, *benchService) {
	c := New()
	inst := &benchService{id: 42}
	c.RegisterService("bench.svc", inst)
	return c, inst
}

// --- Service register / lookup ---

func BenchmarkService_Register_Fresh(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := New()
		svcSinkResult = c.Service("bench.svc", Service{Name: "bench.svc"})
	}
}

func BenchmarkService_RegisterService_Fresh(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := New()
		svcSinkResult = c.RegisterService("bench.svc", &benchService{id: 1})
	}
}

func BenchmarkService_Lookup_Hit_Instance(b *B) {
	c, _ := serviceFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		svcSinkResult = c.Service("bench.svc")
	}
}

func BenchmarkService_Lookup_Miss(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		svcSinkResult = c.Service("noexist.svc")
	}
}

// --- Typed lookups ---

func BenchmarkService_ServiceFor_Hit(b *B) {
	c, _ := serviceFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, svcSinkBool = ServiceFor[*benchService](c, "bench.svc")
	}
}

func BenchmarkService_ServiceFor_TypeMiss(b *B) {
	c, _ := serviceFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, svcSinkBool = ServiceFor[*Core](c, "bench.svc")
	}
}

func BenchmarkService_MustServiceFor_Hit(b *B) {
	c, _ := serviceFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		svcSinkAny = MustServiceFor[*benchService](c, "bench.svc")
	}
}

// --- Listing ---

func BenchmarkService_Services_Empty(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		svcSinkStrings = c.Services()
	}
}

func BenchmarkService_Services_TenRegistered(b *B) {
	c := New()
	for i := 0; i < 10; i++ {
		c.RegisterService(Sprintf("bench.svc.%d", i), &benchService{id: i})
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		svcSinkStrings = c.Services()
	}
}
