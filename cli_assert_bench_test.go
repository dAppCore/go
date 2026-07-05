// SPDX-License-Identifier: EUPL-1.2

// Benchmark for the CLI assertion helper in cli_assert.go. *B satisfies
// testing.TB, so AssertCLI benches against a fake process.run handler on
// the success path (no failure fired). AssertCLIs takes *testing.T, not
// TB, so it cannot be driven from a benchmark.
//
// Run:    go test -bench='BenchmarkAssertCLI' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

func BenchmarkAssertCLI(b *B) {
	c := New()
	c.Action("process.run", func(ctx Context, opts Options) Result {
		return Result{Value: "go version go1.26.0\n", OK: true}
	})
	tc := CLITest{Cmd: "go", Args: []string{"version"}, WantOK: true, Contains: "go1.26"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		AssertCLI(b, c, tc)
	}
}
