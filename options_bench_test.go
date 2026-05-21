// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Options primitive in options.go.
// Per AX-11 — Options is the universal input type. NewOptions is called
// once per Core operation; Get / String / Int / Bool / Float64 /
// Duration fire on every Action handler that reads its input. Slice-
// backed (not map) — linear-scan tradeoff is fine for the small N
// option counts actually seen at call sites.
//
// Run:    go test -bench='BenchmarkOptions' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	optsSinkOptions  Options
	optsSinkResult   Result
	optsSinkBool     bool
	optsSinkString   string
	optsSinkInt      int
	optsSinkFloat    float64
	optsSinkDuration Duration
	optsSinkItems    []Option
)

// Fixtures
var (
	optsBenchSmall = NewOptions(
		Option{Key: "name", Value: "agent"},
		Option{Key: "port", Value: 9000},
	)
	optsBenchTyped = NewOptions(
		Option{Key: "host", Value: "homelab.lan"},
		Option{Key: "port", Value: 9000},
		Option{Key: "debug", Value: true},
		Option{Key: "weight", Value: 0.42},
		Option{Key: "timeout", Value: 5 * Second},
	)
)

// --- Constructors ---

func BenchmarkOptions_NewOptions_Empty(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkOptions = NewOptions()
	}
}

func BenchmarkOptions_NewOptions_Two(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkOptions = NewOptions(
			Option{Key: "name", Value: "agent"},
			Option{Key: "port", Value: 9000},
		)
	}
}

func BenchmarkOptions_NewOptions_Five(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkOptions = NewOptions(
			Option{Key: "host", Value: "homelab.lan"},
			Option{Key: "port", Value: 9000},
			Option{Key: "debug", Value: true},
			Option{Key: "weight", Value: 0.42},
			Option{Key: "timeout", Value: 5 * Second},
		)
	}
}

// --- Set ---

func BenchmarkOptions_Set_Insert(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		o := NewOptions()
		o.Set("k", "v")
	}
}

func BenchmarkOptions_Set_Update(b *B) {
	o := NewOptions(Option{Key: "k", Value: "v"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		o.Set("k", "v2")
	}
}

// --- Get / Has ---

func BenchmarkOptions_Get_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkResult = optsBenchSmall.Get("port")
	}
}

func BenchmarkOptions_Get_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkResult = optsBenchSmall.Get("noexist")
	}
}

func BenchmarkOptions_Has_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkBool = optsBenchSmall.Has("port")
	}
}

// --- Typed accessors ---

func BenchmarkOptions_String_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkString = optsBenchTyped.String("host")
	}
}

func BenchmarkOptions_String_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkString = optsBenchTyped.String("noexist")
	}
}

func BenchmarkOptions_Int_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkInt = optsBenchTyped.Int("port")
	}
}

func BenchmarkOptions_Bool_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkBool = optsBenchTyped.Bool("debug")
	}
}

func BenchmarkOptions_Float64_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkFloat = optsBenchTyped.Float64("weight")
	}
}

func BenchmarkOptions_Duration_Direct(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkDuration = optsBenchTyped.Duration("timeout")
	}
}

func BenchmarkOptions_Duration_FromString(b *B) {
	o := NewOptions(Option{Key: "t", Value: "5s"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkDuration = o.Duration("t")
	}
}

// --- Meta ---

func BenchmarkOptions_Len(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkInt = optsBenchTyped.Len()
	}
}

func BenchmarkOptions_Items(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkItems = optsBenchTyped.Items()
	}
}

// --- Result.New (deprecated but still on call sites until v0.10.0) ---

func BenchmarkOptions_Result_New_VarOnly(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkResult = Result{}.New(42)
	}
}

func BenchmarkOptions_Result_New_PairNoErr(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		optsSinkResult = Result{}.New("data", nil)
	}
}
