// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the entitlement primitives in entitlement.go.
// Per AX-11 — c.Entitled is on every Action.Run (see action.go entitlement
// gate). Default checker is permissive (one closure call); installed
// checkers may do real work (workspace lookup, quota fetch). Both must
// not regress the per-action gate cost. NearLimit + UsagePercent are
// inspection helpers consumers call on the returned Entitlement.
//
// Run:    go test -bench='BenchmarkEntitle' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	entitleSinkEnt   Entitlement
	entitleSinkBool  bool
	entitleSinkFloat float64
)

// --- Entitled (default checker, permissive) ---

func BenchmarkEntitle_Entitled_NoQty(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkEnt = c.Entitled("bench.action")
	}
}

func BenchmarkEntitle_Entitled_Quantity(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkEnt = c.Entitled("bench.action", 3)
	}
}

// --- Entitled (custom checker) ---

func BenchmarkEntitle_Entitled_CustomChecker(b *B) {
	c := New()
	checker := func(action string, qty int, ctx Context) Entitlement {
		return Entitlement{Allowed: true, Limit: 100, Used: 25, Remaining: 75}
	}
	c.SetEntitlementChecker(checker)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkEnt = c.Entitled("bench.action", 1)
	}
}

// --- Entitled denied path (triggers Security log) ---

func BenchmarkEntitle_Entitled_Denied(b *B) {
	quietDefault(b)
	c := New()
	c.SetEntitlementChecker(func(_ string, _ int, _ Context) Entitlement {
		return Entitlement{Allowed: false, Reason: "bench-denied"}
	})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkEnt = c.Entitled("bench.action")
	}
}

// --- Entitlement helpers ---

func BenchmarkEntitle_NearLimit_Below(b *B) {
	e := Entitlement{Allowed: true, Limit: 100, Used: 25}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkBool = e.NearLimit(0.8)
	}
}

func BenchmarkEntitle_NearLimit_Above(b *B) {
	e := Entitlement{Allowed: true, Limit: 100, Used: 95}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkBool = e.NearLimit(0.8)
	}
}

func BenchmarkEntitle_NearLimit_Unlimited(b *B) {
	e := Entitlement{Allowed: true, Unlimited: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkBool = e.NearLimit(0.8)
	}
}

func BenchmarkEntitle_UsagePercent(b *B) {
	e := Entitlement{Allowed: true, Limit: 100, Used: 75}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entitleSinkFloat = e.UsagePercent()
	}
}

// --- Setters + RecordUsage ---

func BenchmarkEntitle_SetEntitlementChecker(b *B) {
	c := New()
	checker := EntitlementChecker(func(_ string, _ int, _ Context) Entitlement {
		return Entitlement{Allowed: true, Unlimited: true}
	})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.SetEntitlementChecker(checker)
	}
}

func BenchmarkEntitle_SetUsageRecorder(b *B) {
	c := New()
	rec := UsageRecorder(func(_ string, _ int, _ Context) {})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.SetUsageRecorder(rec)
	}
}

func BenchmarkEntitle_RecordUsage_Nil(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.RecordUsage("bench.action", 1)
	}
}

func BenchmarkEntitle_RecordUsage_Recorder(b *B) {
	c := New()
	c.SetUsageRecorder(func(_ string, _ int, _ Context) {})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.RecordUsage("bench.action", 1)
	}
}
