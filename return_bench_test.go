// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

var (
	benchReturnIntSink Return[int]
	benchResultSink    Result
	benchIntSink       int
)

func BenchmarkReturnOf(b *B) {
	r := Ok(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchReturnIntSink = ReturnOf[int](r)
	}
}

func BenchmarkReturnTry(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchReturnIntSink = ReturnTry(func() int { return 42 })
	}
}

func BenchmarkReturn_Or(b *B) {
	r := ReturnOK(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchIntSink = r.Or(0)
	}
}

func BenchmarkReturn_Result(b *B) {
	r := ReturnOK(42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResultSink = r.Result()
	}
}
