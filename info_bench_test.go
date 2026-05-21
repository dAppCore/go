// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the system-info primitives in info.go.
// Per AX-11 — Env is the universal env lookup: Core keys hit the
// init-populated map; unknown keys fall through to os.Getenv (a
// syscall on first hit, then cached). OS / Arch / NumCPU are inline
// runtime accessors. StackBuf is paid per crash report.
//
// Run:    go test -bench='BenchmarkInfo' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	infoSinkString  string
	infoSinkInt     int
	infoSinkBytes   []byte
	infoSinkStrings []string
)

// --- Env (hot path) ---

func BenchmarkInfo_Env_CoreKey(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = Env("DIR_HOME")
	}
}

func BenchmarkInfo_Env_OSKey(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = Env("OS")
	}
}

func BenchmarkInfo_Env_Fallthrough(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = Env("PATH")
	}
}

func BenchmarkInfo_Env_Missing(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = Env("__CORE_NO_SUCH_KEY__")
	}
}

// --- EnvKeys (registry walk) ---

func BenchmarkInfo_EnvKeys(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkStrings = EnvKeys()
	}
}

// --- Inline runtime accessors ---

func BenchmarkInfo_OS(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = OS()
	}
}

func BenchmarkInfo_Arch(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = Arch()
	}
}

func BenchmarkInfo_GoVersion(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkString = GoVersion()
	}
}

func BenchmarkInfo_NumCPU(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkInt = NumCPU()
	}
}

// --- StackBuf (crash report path) ---

func BenchmarkInfo_StackBuf(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		infoSinkBytes = StackBuf()
	}
}
