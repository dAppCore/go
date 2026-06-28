// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for lsp.go. White-box (package core) so the JSON-RPC server
// loop — dispatch, the document handlers, readMessage/writeMessage,
// respond/notify, publishDiagnostics — can be driven directly. These are
// unexported and editor-rate (fire per keystroke), not reachable from a
// black-box bench. Output goes to Discard so the numbers are the protocol
// machinery, not the terminal.
//
// Run:    go test -bench='BenchmarkLSP|BenchmarkLsp' -benchmem -run='^$' .

package core

var (
	lspSink    []LSPDiagnostic
	lspSinkStr []string
	lspSinkRes Result
)

func lspBenchServer() *lspServer {
	return &lspServer{out: Discard, documents: NewRegistry[[]byte]()}
}

// --- diagnostic pipeline (exported) ---

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

// --- JSON-RPC server loop (white-box) ---

func BenchmarkLspServer_dispatch_Initialize(b *B) {
	srv := lspBenchServer()
	raw := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.dispatch(raw)
	}
}

func BenchmarkLspServer_dispatch_DidOpen(b *B) {
	srv := lspBenchServer()
	raw := []byte(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///x.go","text":"package x\n"}}}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.dispatch(raw) // → handleDocumentSync → lspExtractDocument → publishDiagnostics → notify → writeMessage
	}
}

func BenchmarkLspServer_dispatch_DidChange(b *B) {
	srv := lspBenchServer()
	raw := []byte(`{"jsonrpc":"2.0","method":"textDocument/didChange","params":{"textDocument":{"uri":"file:///x.go"},"contentChanges":[{"text":"package x\n"}]}}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.dispatch(raw) // → handleDocumentChange → lspExtractDocumentChange
	}
}

func BenchmarkLspServer_dispatch_DidClose(b *B) {
	srv := lspBenchServer()
	srv.dispatch([]byte(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///x.go","text":"package x\n"}}}`))
	raw := []byte(`{"jsonrpc":"2.0","method":"textDocument/didClose","params":{"textDocument":{"uri":"file:///x.go"}}}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.dispatch(raw) // → handleDocumentClose
	}
}

func BenchmarkLspServer_readMessage(b *B) {
	frame := "Content-Length: 17\r\n\r\n{\"jsonrpc\":\"2.0\"}"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv := lspBenchServer()
		srv.in = NewBufReader(NewReader(frame))
		lspSinkRes = srv.readMessage()
	}
}

func BenchmarkLspServer_writeMessage(b *B) {
	srv := lspBenchServer()
	payload := lspMessage{JSONRPC: "2.0", Method: "textDocument/publishDiagnostics"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lspSinkRes = srv.writeMessage(payload)
	}
}

func BenchmarkLspServer_respond(b *B) {
	srv := lspBenchServer()
	id := 1
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.respond(&id, nil, nil)
	}
}

func BenchmarkLspServer_notify(b *B) {
	srv := lspBenchServer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.notify("textDocument/publishDiagnostics", nil)
	}
}

func BenchmarkLspServer_publishDiagnostics(b *B) {
	srv := lspBenchServer()
	content := []byte("package x\n")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv.publishDiagnostics("file:///x.go", content)
	}
}

// LSPServe with an already-cancelled context: the run loop checks ctx
// before reading stdin, so it returns immediately — covers the entry +
// shutdown path without blocking on input.
func BenchmarkLSPServe(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx, cancel := WithCancel(Background())
		cancel()
		lspSinkRes = LSPServe(ctx)
	}
}
