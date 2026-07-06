// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Registry primitive in registry.go.
// Per AX-11 — Registry[T] underpins service / action / drive / command
// / data / config / i18n / contract registries. Every Set / Get / Has /
// Names / Each / List call across the framework lands here. Lock / Seal
// are state ops only fired at boot but worth gating against regressions.
//
// Run:    go test -bench='BenchmarkRegistry' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	registrySinkResult  Result
	registrySinkBool    bool
	registrySinkInt     int
	registrySinkStrings []string
	registrySinkVals    []int
)

// registryFixture builds a Registry[int] pre-populated with n entries.
func registryFixture(n int) *Registry[int] {
	r := NewRegistry[int]()
	for i := 0; i < n; i++ {
		r.Set(Sprintf("entry.%d", i), i)
	}
	return r
}

// --- Constructor ---

func BenchmarkRegistry_NewRegistry(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewRegistry[int]()
	}
}

// --- Set / Get / Has ---

func BenchmarkRegistry_Set_Fresh(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewRegistry[int]()
		registrySinkResult = r.Set("k", 42)
	}
}

func BenchmarkRegistry_Get_Hit(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkResult = r.Get("entry.50")
	}
}

func BenchmarkRegistry_Get_Miss(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkResult = r.Get("noexist")
	}
}

func BenchmarkRegistry_Has_Hit(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkBool = r.Has("entry.50")
	}
}

func BenchmarkRegistry_Has_Miss(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkBool = r.Has("noexist")
	}
}

// --- Names / List / Each / Len ---

func BenchmarkRegistry_Names_100(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkStrings = r.Names()
	}
}

func BenchmarkRegistry_List_PrefixHit(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkVals = r.List("entry.*")
	}
}

func BenchmarkRegistry_List_NarrowHit(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkVals = r.List("entry.9*")
	}
}

func BenchmarkRegistry_Each_100(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sum := 0
		r.Each(func(_ string, v int) {
			sum += v
		})
		registrySinkInt = sum
	}
}

func BenchmarkRegistry_Len(b *B) {
	r := registryFixture(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkInt = r.Len()
	}
}

// --- Delete ---

func BenchmarkRegistry_Delete(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := registryFixture(10)
		registrySinkResult = r.Delete("entry.5")
	}
}

// --- Enable / Disable / Disabled ---

func BenchmarkRegistry_Disable(b *B) {
	r := registryFixture(10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Disable("entry.5")
	}
}

func BenchmarkRegistry_Enable(b *B) {
	r := registryFixture(10)
	r.Disable("entry.5")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Enable("entry.5")
	}
}

func BenchmarkRegistry_Disabled_Hit(b *B) {
	r := registryFixture(10)
	r.Disable("entry.5")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkBool = r.Disabled("entry.5")
	}
}

// --- Lock / Seal state ---

func BenchmarkRegistry_Lock(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewRegistry[int]()
		r.Lock()
	}
}

func BenchmarkRegistry_Locked(b *B) {
	r := NewRegistry[int]()
	r.Lock()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkBool = r.Locked()
	}
}

func BenchmarkRegistry_Seal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewRegistry[int]()
		r.Seal()
	}
}

func BenchmarkRegistry_Sealed(b *B) {
	r := NewRegistry[int]()
	r.Seal()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		registrySinkBool = r.Sealed()
	}
}

func BenchmarkRegistry_Open(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewRegistry[int]()
		r.Lock()
		r.Open()
	}
}
