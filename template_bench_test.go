// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the text-template primitives in template.go.
// Per AX-11 — templates back every CLI banner, every dispatch brief,
// every onboarding doc. ParseTemplate is one-shot per template; the
// hot path is ExecuteTemplate firing on each render. The bench gates
// both surfaces.
//
// Run:    go test -bench='BenchmarkTemplate' -benchmem -run='^$' .

package core_test

import (
	"bytes"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	templateSinkTemplate *Template
	templateSinkResult   Result
)

// Fixtures
var (
	templateSimpleText = "Hello {{.Name}}!"
	templateRichText   = `Agent {{.Name}} ({{.Role}}) on host {{.Host}}:{{.Port}}
{{range .Tags}}  - {{.}}
{{end}}Status: {{.Status}}`
)

type templateData struct {
	Name, Role, Host string
	Port             int
	Tags             []string
	Status           string
}

var templateRichData = templateData{
	Name: "cladius", Role: "project-leader",
	Host: "homelab.lan", Port: 9000,
	Tags: []string{"alive", "sandbox", "verified"}, Status: "ready",
}

// --- NewTemplate ---

func BenchmarkTemplate_NewTemplate(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		templateSinkTemplate = NewTemplate("bench")
	}
}

// --- Parse ---

func BenchmarkTemplate_Parse_Simple(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		templateSinkResult = ParseTemplate("bench", templateSimpleText)
	}
}

func BenchmarkTemplate_Parse_Rich(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		templateSinkResult = ParseTemplate("bench", templateRichText)
	}
}

// --- Execute (hot path) ---

func BenchmarkTemplate_Execute_Simple(b *B) {
	r := ParseTemplate("bench", templateSimpleText)
	tmpl := r.Value.(*Template)
	var buf bytes.Buffer
	data := struct{ Name string }{Name: "Snider"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		templateSinkResult = ExecuteTemplate(tmpl, &buf, data)
	}
}

func BenchmarkTemplate_Execute_Rich(b *B) {
	r := ParseTemplate("bench", templateRichText)
	tmpl := r.Value.(*Template)
	var buf bytes.Buffer
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		templateSinkResult = ExecuteTemplate(tmpl, &buf, templateRichData)
	}
}

// --- ParseFS ---

func BenchmarkTemplate_ParseFS(b *B) {
	// DirFS(".") on the package dir gives us a known file (one of the
	// bench_test.go files) we can ParseFS against; we expect it to fail
	// to parse as a template (text won't be a valid template surface).
	// That's fine — the bench captures the open + read + parse-attempt
	// path, which is what consumers pay.
	fsys := DirFS(".")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		templateSinkResult = ParseTemplateFS(fsys, "doc.go")
	}
}

func BenchmarkTemplate_ParseFiles(b *B) {
	path := PathJoin(b.TempDir(), "greeting.tmpl")
	WriteFile(path, []byte("hello {{.Name}}"), 0o644)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		templateSinkResult = ParseTemplateFiles(path)
	}
}
