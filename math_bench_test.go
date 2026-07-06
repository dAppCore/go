// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the math + ordering primitives in math.go.
// Per AX-11 — these wrappers sit on tight loops in tokenisers
// (argmax over logits), config validation, and numeric formatters.
// Even a 1-2 ns per-call cost compounds across token-generation
// inner loops in go-mlx.
//
// Run:    go test -bench='BenchmarkMath|BenchmarkCompare|BenchmarkMin|BenchmarkMax|BenchmarkAbs|BenchmarkPow|BenchmarkFloor|BenchmarkCeil|BenchmarkRound|BenchmarkNaN' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks — prevent compiler dead-code elimination from collapsing the
// bench loop. Every benchmark stores its result here so the call must
// actually execute.
var (
	mathSinkInt    int
	mathSinkFloat  float64
	mathSinkBool   bool
	mathSinkString string
	_              = mathSinkString // touched to keep declaration
)

// --- Compare / Min / Max ---

func BenchmarkCompare_Int(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Compare(i, b.N-i)
	}
}

func BenchmarkCompare_String(b *B) {
	a, c := "alpha", "beta"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if i&1 == 0 {
			mathSinkInt = Compare(a, c)
		} else {
			mathSinkInt = Compare(c, a)
		}
	}
}

func BenchmarkMin_Int(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Min(i, b.N-i)
	}
}

func BenchmarkMax_Int(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Max(i, b.N-i)
	}
}

func BenchmarkMin_Float(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Min(float64(i), float64(b.N-i))
	}
}

func BenchmarkMax_Float(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Max(float64(i), float64(b.N-i))
	}
}

// --- Abs ---

func BenchmarkAbs_PosInt(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Abs(i)
	}
}

func BenchmarkAbs_NegInt(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Abs(-i)
	}
}

func BenchmarkAbs_Float(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Abs(-float64(i) - 0.14)
	}
}

// --- Clamp / Sign ---

func BenchmarkClamp_Int_InRange(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Clamp(50, 0, 100)
	}
}

func BenchmarkClamp_Int_BelowLo(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Clamp(-10, 0, 100)
	}
}

func BenchmarkClamp_Float(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Clamp(0.42, 0.0, 1.0)
	}
}

func BenchmarkSign_Pos(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Sign(42)
	}
}

func BenchmarkSign_Neg(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Sign(-42)
	}
}

func BenchmarkSign_Zero(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkInt = Sign(0)
	}
}

// --- NaN / IsNaN ---

func BenchmarkNaN(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = NaN()
	}
}

func BenchmarkIsNaN_Number(b *B) {
	v := 3.14
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkBool = IsNaN(v + float64(i))
	}
}

func BenchmarkIsNaN_NaN(b *B) {
	v := NaN()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkBool = IsNaN(v)
	}
}

// --- Floating-point primitives ---

func BenchmarkPow(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Pow(float64(i)+2, 2)
	}
}

func BenchmarkFloor(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Floor(float64(i) + 0.7)
	}
}

func BenchmarkCeil(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Ceil(float64(i) + 0.1)
	}
}

func BenchmarkRound(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mathSinkFloat = Round(float64(i) + 0.5)
	}
}
