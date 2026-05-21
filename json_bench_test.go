// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the JSON primitives in json.go.
// Per AX-11 — JSON marshal/unmarshal sits on every config-load,
// API-response, persisted-state, IPC-message, and metric-emit path.
// Even a 5-10% allocation drop here compounds across the ecosystem.
//
// The bench surface deliberately covers small/medium/large payload
// shapes and the convenience String variants so we can see whether
// the []byte(s) and string(data) round-trips at the API edges are
// material vs the underlying json.Marshal/Unmarshal work.
//
// Run:    go test -bench='BenchmarkJSON' -benchmem -run='^$' .
// Filter: go test -bench='BenchmarkJSONMarshal' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// --- Fixtures ---

type benchJSONSmall struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type benchJSONMedium struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Tags      []string          `json:"tags"`
	Score     float64           `json:"score"`
	Active    bool              `json:"active"`
	Counters  map[string]int    `json:"counters"`
	UpdatedAt int64             `json:"updated_at"`
	Owner     benchJSONSmall    `json:"owner"`
	Labels    map[string]string `json:"labels"`
}

type benchJSONLarge struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Status    string             `json:"status"`
	Items     []benchJSONMedium  `json:"items"`
	History   []string           `json:"history"`
	Meta      map[string]string  `json:"meta"`
	Stats     map[string]float64 `json:"stats"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
}

var (
	fixtureSmall = benchJSONSmall{Port: 8080, Host: "homelab.lan"}

	fixtureMedium = benchJSONMedium{
		ID:     "agt-7421",
		Name:   "cladius",
		Tags:   []string{"agent", "go", "tokeniser", "production"},
		Score:  0.8732,
		Active: true,
		Counters: map[string]int{
			"requests": 12384,
			"errors":   3,
			"warnings": 17,
		},
		UpdatedAt: 1715763600,
		Owner:     benchJSONSmall{Port: 9000, Host: "forge.lthn.sh"},
		Labels: map[string]string{
			"env":    "production",
			"region": "homelab",
			"tier":   "primary",
		},
	}

	fixtureLarge = benchJSONLarge{
		ID:     "batch-2026-05-21",
		Name:   "AX-11 propagation sweep",
		Status: "running",
		Items: []benchJSONMedium{
			fixtureMedium, fixtureMedium, fixtureMedium,
			fixtureMedium, fixtureMedium, fixtureMedium,
			fixtureMedium, fixtureMedium,
		},
		History: []string{
			"normalize-cache", "scratch-pool", "byte-elision",
			"pre-warm", "ascii-fast-path", "concat-grow",
		},
		Meta: map[string]string{
			"author":   "cladius",
			"reviewer": "snider",
			"trigger":  "AX-11 sweep",
		},
		Stats: map[string]float64{
			"speedup":      2.21,
			"alloc_drop":   0.95,
			"scores_per_s": 15898,
		},
		CreatedAt: 1715000000,
		UpdatedAt: 1715763600,
	}

	// Pre-marshalled fixtures for Unmarshal benches.
	jsonSmall, _  = marshalToBytes(fixtureSmall)
	jsonMedium, _ = marshalToBytes(fixtureMedium)
	jsonLarge, _  = marshalToBytes(fixtureLarge)

	jsonSmallStr  = string(jsonSmall)
	jsonMediumStr = string(jsonMedium)
)

func marshalToBytes(v any) ([]byte, error) {
	r := JSONMarshal(v)
	if !r.OK {
		return nil, r.Value.(error)
	}
	return r.Value.([]byte), nil
}

// --- Marshal ---

func BenchmarkJSONMarshal_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = JSONMarshal(fixtureSmall)
	}
}

func BenchmarkJSONMarshal_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = JSONMarshal(fixtureMedium)
	}
}

func BenchmarkJSONMarshal_Large(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = JSONMarshal(fixtureLarge)
	}
}

func BenchmarkJSONMarshalIndent_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = JSONMarshalIndent(fixtureMedium, "", "  ")
	}
}

func BenchmarkJSONMarshalString_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = JSONMarshalString(fixtureMedium)
	}
}

// --- Unmarshal ---

func BenchmarkJSONUnmarshal_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v benchJSONSmall
		_ = JSONUnmarshal(jsonSmall, &v)
	}
}

func BenchmarkJSONUnmarshal_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v benchJSONMedium
		_ = JSONUnmarshal(jsonMedium, &v)
	}
}

func BenchmarkJSONUnmarshal_Large(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v benchJSONLarge
		_ = JSONUnmarshal(jsonLarge, &v)
	}
}

func BenchmarkJSONUnmarshalString_Small(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v benchJSONSmall
		_ = JSONUnmarshalString(jsonSmallStr, &v)
	}
}

func BenchmarkJSONUnmarshalString_Medium(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v benchJSONMedium
		_ = JSONUnmarshalString(jsonMediumStr, &v)
	}
}

// --- Unmarshal into generic targets ---

func BenchmarkJSONUnmarshal_IntoMap(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v map[string]any
		_ = JSONUnmarshal(jsonMedium, &v)
	}
}

func BenchmarkJSONUnmarshal_IntoRawMessage(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var v RawMessage
		_ = JSONUnmarshal(jsonMedium, &v)
	}
}
