// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the utility helpers in utils.go.
// Per AX-11 — these are small but call-site-frequent: ID() is generated
// per task / per request / per agent event; ValidateName guards service /
// action / command registration; ParseFlag fires on every CLI invocation.
//
// Run:    go test -bench='BenchmarkUtils' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	utilsSinkString  string
	utilsSinkStrings []string
	utilsSinkResult  Result
	utilsSinkBool    bool
	utilsSinkInt     int
	utilsSinkAny     any
)

// --- ID ---

func BenchmarkUtils_ID(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkString = ID()
	}
}

// --- ValidateName ---

func BenchmarkUtils_ValidateName_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkResult = ValidateName("brain")
	}
}

func BenchmarkUtils_ValidateName_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkResult = ValidateName("")
	}
}

func BenchmarkUtils_ValidateName_Traversal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkResult = ValidateName("../escape")
	}
}

// --- SanitisePath ---

func BenchmarkUtils_SanitisePath_Safe(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkString = SanitisePath("agent.json")
	}
}

func BenchmarkUtils_SanitisePath_Traversal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkString = SanitisePath("../../etc/passwd")
	}
}

// --- JoinPath / IsFlag ---

func BenchmarkUtils_JoinPath_Three(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkString = JoinPath("deploy", "to", "homelab")
	}
}

func BenchmarkUtils_IsFlag_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkBool = IsFlag("--verbose")
	}
}

func BenchmarkUtils_IsFlag_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkBool = IsFlag("deploy")
	}
}

// --- Arg / ArgString / ArgInt / ArgBool ---

func BenchmarkUtils_Arg_String(b *B) {
	args := []any{"name", 42, true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkResult = Arg(0, args...)
	}
}

func BenchmarkUtils_ArgString_Hit(b *B) {
	args := []any{"name", 42, true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkString = ArgString(0, args...)
	}
}

func BenchmarkUtils_ArgString_OOB(b *B) {
	args := []any{"name"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkString = ArgString(5, args...)
	}
}

func BenchmarkUtils_ArgInt(b *B) {
	args := []any{"name", 42}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkInt = ArgInt(1, args...)
	}
}

func BenchmarkUtils_ArgBool(b *B) {
	args := []any{"name", 42, true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkBool = ArgBool(2, args...)
	}
}

// --- FilterArgs ---

func BenchmarkUtils_FilterArgs(b *B) {
	args := []string{"deploy", "-test.run=foo", "homelab", "", "--verbose"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utilsSinkStrings = FilterArgs(args)
	}
}

// --- ParseFlag ---

func BenchmarkUtils_ParseFlag_DoubleDashName(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k, v, ok := ParseFlag("--verbose")
		utilsSinkString = k
		_ = v
		utilsSinkBool = ok
	}
}

func BenchmarkUtils_ParseFlag_DoubleDashKeyVal(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k, v, ok := ParseFlag("--port=8080")
		utilsSinkString = k
		_ = v
		utilsSinkBool = ok
	}
}

func BenchmarkUtils_ParseFlag_SingleDash(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k, v, ok := ParseFlag("-v")
		utilsSinkString = k
		_ = v
		utilsSinkBool = ok
	}
}

func BenchmarkUtils_ParseFlag_NotAFlag(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k, v, ok := ParseFlag("hello")
		utilsSinkString = k
		_ = v
		utilsSinkBool = ok
	}
}
