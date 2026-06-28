// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the LSP diagnostic pipeline in lsp.go. The four
// built-in sources (ax-7 naming, result-shape, spor, test-imports)
// are registered in init(), so LSPComputeDiagnostics exercises the
// whole scan + helper chain on each call. Content is crafted to enter
// the result-shape and import-extraction paths so their helpers run.
//
// The JSON-RPC server loop (run/readMessage/dispatch/handle*) is driven
// by stdin and is exercised by the functional tests, not benchmarked
// here — it is request-rate bound, not a per-token hot path.
//
// Run:    go test -bench='BenchmarkLSP' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var (
	lspSink    []LSPDiagnostic
	lspSinkStr []string
)

func BenchmarkLSPComputeDiagnostics(b *B) {
	uri := "file:///agent/worker.go"
	content := []byte(`package worker

import "strings"

func process(in string) (string, error) {
	return strings.TrimSpace(in), nil
}

type widget struct{ name string }
`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lspSink = LSPComputeDiagnostics(uri, content)
	}
}

func BenchmarkLSPComputeDiagnostics_Test(b *B) {
	uri := "file:///agent/worker_test.go"
	content := []byte(`package worker_test

import "fmt"

func TestProcess(t *T) {
	fmt.Println("ok")
}
`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lspSink = LSPComputeDiagnostics(uri, content)
	}
}

func BenchmarkLSPDiagnosticSources(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lspSinkStr = LSPDiagnosticSources()
	}
}
