// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Data primitive in data.go.
// Per AX-11 — Data is the embedded-asset registry every agent persona,
// every prompt template, every CLI banner loads from. Read hot path is
// ReadString — agent boot reads dozens; chat-turn fast-path reads one.
// The bench gates the resolve→read pipeline.
//
// Run:    go test -bench='BenchmarkData' -benchmem -run='^$' .

package core_test

import (
	"testing/fstest"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	dataSinkResult  Result
	dataSinkStrings []string
)

// dataFixture mounts a small in-memory FS as a Data entry.
// Mirrors how consumer packages declare their embedded.FS at boot.
func dataFixture() *Core {
	c := New()
	fsys := fstest.MapFS{
		"prompts/coding.md":         &fstest.MapFile{Data: []byte("# Coding\n\n```go\nfunc main() {}\n```\n")},
		"prompts/review.md":         &fstest.MapFile{Data: []byte("# Review checklist\n- correctness\n- tests\n- docs\n")},
		"prompts/triage.md":         &fstest.MapFile{Data: []byte("# Triage\n\nSurface the smallest reproducer.\n")},
		"flow/deploy/homelab.yaml":  &fstest.MapFile{Data: []byte("target: homelab\nbranch: dev\n")},
		"flow/deploy/de1.yaml":      &fstest.MapFile{Data: []byte("target: de1\nbranch: main\n")},
		"persona/code/cladius.md":   &fstest.MapFile{Data: []byte("persona: cladius\n")},
		"persona/code/hephaestus.md": &fstest.MapFile{Data: []byte("persona: hephaestus\n")},
	}
	c.Data().New(NewOptions(
		Option{Key: "name", Value: "brain"},
		Option{Key: "source", Value: FS(fsys)},
		Option{Key: "path", Value: "."},
	))
	return c
}

// --- Read paths ---

func BenchmarkData_ReadFile_Hit(b *B) {
	c := dataFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkResult = c.Data().ReadFile("brain/prompts/coding.md")
	}
}

func BenchmarkData_ReadFile_Miss(b *B) {
	c := dataFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkResult = c.Data().ReadFile("brain/noexist.md")
	}
}

func BenchmarkData_ReadString_Hit(b *B) {
	c := dataFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkResult = c.Data().ReadString("brain/prompts/coding.md")
	}
}

// --- List paths ---

func BenchmarkData_List(b *B) {
	c := dataFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkResult = c.Data().List("brain/prompts")
	}
}

func BenchmarkData_ListNames(b *B) {
	c := dataFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkResult = c.Data().ListNames("brain/prompts")
	}
}

// --- Mounts registry ---

func BenchmarkData_Mounts(b *B) {
	c := dataFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkStrings = c.Data().Mounts()
	}
}

// --- New (mount registration) ---

func BenchmarkData_New(b *B) {
	fsys := fstest.MapFS{
		"a.md": &fstest.MapFile{Data: []byte("a")},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := New()
		dataSinkResult = c.Data().New(NewOptions(
			Option{Key: "name", Value: "bench"},
			Option{Key: "source", Value: FS(fsys)},
			Option{Key: "path", Value: "."},
		))
	}
}

func BenchmarkData_Extract(b *B) {
	fsys := fstest.MapFS{
		"tmpl/a.md": &fstest.MapFile{Data: []byte("content")},
	}
	c := New()
	c.Data().New(NewOptions(
		Option{Key: "name", Value: "bench"},
		Option{Key: "source", Value: FS(fsys)},
		Option{Key: "path", Value: "."},
	))
	base := b.TempDir()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dataSinkResult = c.Data().Extract("tmpl", PathJoin(base, Itoa(i)), nil)
	}
}
