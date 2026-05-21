// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the path primitives in path.go.
// Per AX-11 — path ops sit on every model-file resolution, config-
// location compute, file walker, embed traversal. Most are filepath
// passthroughs but the bench harness gates the contract.
//
// Run:    go test -bench='BenchmarkPath' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	pathSinkString  string
	pathSinkBool    bool
	pathSinkStrings []string
	pathSinkResult  Result
)

// --- Path / Join ---

func BenchmarkPath_Two(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = Path("/home/agent", "config.json")
	}
}

func BenchmarkPath_Four(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = Path("/home", "agent", "models", "qwen.gguf")
	}
}

func BenchmarkPathJoin_Four(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = PathJoin("/home", "agent", "models", "qwen.gguf")
	}
}

// --- Base / Dir / Ext ---

func BenchmarkPathBase(b *B) {
	p := "/home/agent/models/qwen3.6-q4.gguf"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = PathBase(p)
	}
}

func BenchmarkPathDir(b *B) {
	p := "/home/agent/models/qwen3.6-q4.gguf"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = PathDir(p)
	}
}

func BenchmarkPathExt(b *B) {
	p := "/home/agent/models/qwen3.6-q4.gguf"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = PathExt(p)
	}
}

// --- Predicates ---

func BenchmarkPathIsAbs_Abs(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkBool = PathIsAbs("/absolute/path")
	}
}

func BenchmarkPathIsAbs_Rel(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkBool = PathIsAbs("relative/path")
	}
}

// --- Clean ---

func BenchmarkCleanPath_NoChange(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = CleanPath("/already/clean/path", "/")
	}
}

func BenchmarkCleanPath_DotDot(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = CleanPath("/home/agent/../models/./qwen", "/")
	}
}

// --- Match / Rel / ChangeExt ---

func BenchmarkPathMatch_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkResult = PathMatch("*.gguf", "qwen3.6-q4.gguf")
	}
}

func BenchmarkPathRel(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkResult = PathRel("/home/agent", "/home/agent/models/qwen.gguf")
	}
}

func BenchmarkPathChangeExt(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = PathChangeExt("/home/agent/qwen.gguf", ".safetensors")
	}
}

func BenchmarkPathToSlash(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pathSinkString = PathToSlash("/home/agent/models")
	}
}
