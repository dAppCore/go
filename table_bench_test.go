// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Table primitive in table.go.
// Per AX-11 — Table backs every CLI list output (services, registries,
// audits, ticket triage). Cells are aligned via text/tabwriter; the
// bench gates the alignment path + flush overhead.
//
// Run:    go test -bench='BenchmarkTable' -benchmem -run='^$' .

package core_test

import (
	"bytes"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	tableSinkTable  *Table
	tableSinkResult Result
)

// --- Constructor ---

func BenchmarkTable_NewTable(b *B) {
	var buf bytes.Buffer
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		tableSinkTable = NewTable(&buf)
	}
}

// --- Row + Flush cycle ---

func BenchmarkTable_Row_TwoCells(b *B) {
	var buf bytes.Buffer
	table := NewTable(&buf)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = table.Row("name", "value")
	}
}

func BenchmarkTable_Row_FourCells(b *B) {
	var buf bytes.Buffer
	table := NewTable(&buf)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = table.Row("service", "running", "homelab", "9000")
	}
}

func BenchmarkTable_Flush_Empty(b *B) {
	var buf bytes.Buffer
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		table := NewTable(&buf)
		tableSinkResult = table.Flush()
	}
}

func BenchmarkTable_BuildAndFlush_10Rows(b *B) {
	var buf bytes.Buffer
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		table := NewTable(&buf)
		table.Row("Name", "Status", "Host")
		for range 10 {
			table.Row("agent", "ok", "homelab")
		}
		tableSinkResult = table.Flush()
	}
}
