// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func ExampleLSPDiagnostic() {
	d := LSPDiagnostic{
		Range:    LSPRange{Start: LSPPosition{Line: 4}, End: LSPPosition{Line: 4, Character: 80}},
		Severity: LSPSeverityWarning,
		Source:   "test-imports",
		Message:  "imports 'context' — use core.Background()",
	}
	Println(d.Source)
	Println(d.Message)
	Println(d.Range.End.Character)
	// Output:
	// test-imports
	// imports 'context' — use core.Background()
	// 80
}

func ExampleLSPComputeDiagnostics() {
	// Register a diagnostic source, then compute diagnostics for a document
	// without starting the full server. Output is omitted: registered sources
	// are global process state, so the result depends on what else the process
	// has registered.
	LSPRegisterDiagnostic("example-todo", func(uri string, content []byte) []LSPDiagnostic {
		if Contains(string(content), "TODO") {
			return []LSPDiagnostic{{Message: "contains TODO", Source: "example-todo"}}
		}
		return nil
	})
	_ = LSPComputeDiagnostics("file:///x.go", []byte("// TODO\n"))
	_ = LSPDiagnosticSources()
}

// ExampleLSPRegisterDiagnostic registers a diagnostic source through `LSPRegisterDiagnostic`.
func ExampleLSPRegisterDiagnostic() {
	LSPRegisterDiagnostic("example-linter", func(uri string, content []byte) []LSPDiagnostic {
		return nil
	})
	// Registered sources are global process state; no deterministic output.
}

// ExampleLSPDiagnosticSources lists registered diagnostic sources through `LSPDiagnosticSources`.
func ExampleLSPDiagnosticSources() {
	_ = LSPDiagnosticSources() // the names of every registered diagnostic source
}

// ExampleLSPServe runs the language-server loop over stdio through `LSPServe`.
func ExampleLSPServe() {
	if false {
		LSPServe(Background()) // blocks serving LSP over stdin/stdout until ctx is cancelled
	}
}
