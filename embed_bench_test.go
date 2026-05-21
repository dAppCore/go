// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the embed primitives in embed.go.
// Per AX-11 — Embed is the substrate behind c.Data() and every package
// that ships brain prompts, CLI templates, agent flows, or onboarding
// docs. Mount() / ReadFile() / ReadString() fire on every persona
// boot; Sub() / List() drive the directory walkers; AddAsset /
// GetAsset back the in-memory asset registry used by generated packs.
//
// Run:    go test -bench='BenchmarkEmbed' -benchmem -run='^$' .

package core_test

import (
	"testing/fstest"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	embedSinkResult Result
	embedSinkEmbed  *Embed
	embedSinkFS     FS
	embedSinkString string
)

// embedFixture returns a small FS suitable for the per-op read paths.
func embedFixture() FS {
	return FS(fstest.MapFS{
		"prompts/code.md":         &fstest.MapFile{Data: []byte("# Code prompt body\nLong enough to bench against.\n")},
		"prompts/review.md":       &fstest.MapFile{Data: []byte("# Review prompt\nChecklist content here.\n")},
		"flow/deploy/homelab.yml": &fstest.MapFile{Data: []byte("target: homelab\n")},
		"flow/deploy/de1.yml":     &fstest.MapFile{Data: []byte("target: de1\n")},
		"agent/persona/code.md":   &fstest.MapFile{Data: []byte("persona body\n")},
	})
}

// --- Mount ---

func BenchmarkEmbed_Mount_Root(b *B) {
	fsys := embedFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = Mount(fsys, ".")
	}
}

func BenchmarkEmbed_Mount_Sub(b *B) {
	fsys := embedFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = Mount(fsys, "prompts")
	}
}

// --- Embed read paths ---

func BenchmarkEmbed_ReadFile(b *B) {
	r := Mount(embedFixture(), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = emb.ReadFile("prompts/code.md")
	}
}

func BenchmarkEmbed_ReadString(b *B) {
	r := Mount(embedFixture(), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = emb.ReadString("prompts/code.md")
	}
}

func BenchmarkEmbed_Open(b *B) {
	r := Mount(embedFixture(), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = emb.Open("prompts/code.md")
		if embedSinkResult.OK {
			CloseStream(embedSinkResult.Value)
		}
	}
}

func BenchmarkEmbed_ReadDir(b *B) {
	r := Mount(embedFixture(), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = emb.ReadDir("prompts")
	}
}

// --- Sub + FS accessors ---

func BenchmarkEmbed_Sub(b *B) {
	r := Mount(embedFixture(), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = emb.Sub("prompts")
	}
}

func BenchmarkEmbed_FS(b *B) {
	r := Mount(embedFixture(), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkFS = emb.FS()
	}
}

func BenchmarkEmbed_BaseDirectory(b *B) {
	r := Mount(embedFixture(), "prompts")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkString = emb.BaseDirectory()
	}
}

// --- In-memory asset registry (AddAsset / GetAsset) ---

func BenchmarkEmbed_AddAsset(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AddAsset("bench", "key", "value")
	}
}

func BenchmarkEmbed_GetAsset_Hit(b *B) {
	AddAsset("bench-hit", "key", "value")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = GetAsset("bench-hit", "key")
	}
}

func BenchmarkEmbed_GetAsset_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = GetAsset("bench-miss", "noexist")
	}
}

func BenchmarkEmbed_GetAssetBytes_Hit(b *B) {
	AddAsset("bench-bytes", "key", "value")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		embedSinkResult = GetAssetBytes("bench-bytes", "key")
	}
}
