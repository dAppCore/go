// SPDX-License-Identifier: EUPL-1.2

// White-box benchmarks for unexported helpers a black-box (core_test)
// bench can't reach: action/task name sanitisation, sandbox path
// resolution, and the assertion failure/message formatters (driven
// through a capturing TB stub so the bench exercises the real formatting
// path without actually failing).
//
// Run:    go test -bench='Benchmark(Action_safeName|Task_safeName|Fs_path|AssertFail|AssertMsg)' -benchmem -run='^$' .

package core

// assertBenchStub embeds *B — inheriting the unexported method that makes
// it a testing.TB — and swallows Error/Fatal so assertFail runs its real
// formatting path without failing the benchmark.
type assertBenchStub struct{ *B }

func (assertBenchStub) Error(args ...any) {}
func (assertBenchStub) Fatal(args ...any) {}

var whiteboxSink string

func BenchmarkAction_safeName(b *B) {
	a := &Action{Name: "agentic.dispatch"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		whiteboxSink = a.safeName()
	}
}

func BenchmarkTask_safeName(b *B) {
	t := &Task{Name: "deploy.homelab"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		whiteboxSink = t.safeName()
	}
}

func BenchmarkFs_path(b *B) {
	m := (&Fs{}).New("/tmp")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		whiteboxSink = m.path("workspace/agent/readme.md")
	}
}

func BenchmarkAssertFail(b *B) {
	stub := assertBenchStub{b}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		assertFail(stub, false, "Bench", nil, "want", 1, "got", 2)
	}
}

func BenchmarkAssertMsg(b *B) {
	msg := []string{"context line", "detail line"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		whiteboxSink = assertMsg(msg)
	}
}
