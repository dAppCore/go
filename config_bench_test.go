// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Config primitives in config.go.
// Per AX-11 — c.Config() backs every settings read across the
// framework: host names, ports, paths, feature gates. String/Int/Bool
// fire per-call from CLI command handlers, Service.OnStart wiring,
// and Wails-RPC bindings. Feature flag Enabled() runs on every gated
// code path. ConfigVar[T] is the per-package field used by services
// to declare typed config slots.
//
// Run:    go test -bench='BenchmarkConfig' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	configSinkConfig  *Config
	configSinkString  string
	configSinkInt     int
	configSinkBool    bool
	configSinkResult  Result
	configSinkStrings []string
)

// configFixture returns a Config pre-populated with typical settings.
func configFixture() *Config {
	c := (&Config{}).New()
	c.Set("database.host", "homelab.lthn.sh")
	c.Set("database.port", 5432)
	c.Set("api.timeout", "30s")
	c.Set("debug", true)
	c.Set("weight", 0.42)
	c.Enable("dark-mode")
	c.Enable("experimental-fleet")
	return c
}

// --- Construction ---

func BenchmarkConfig_New(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkConfig = (&Config{}).New()
	}
}

// --- Settings (Set / Get) ---

func BenchmarkConfig_Set_Fresh(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cfg := (&Config{}).New()
		cfg.Set("key", "value")
	}
}

func BenchmarkConfig_Set_Update(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cfg.Set("database.host", "homelab.lthn.sh")
	}
}

func BenchmarkConfig_Get_Hit(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkResult = cfg.Get("database.host")
	}
}

func BenchmarkConfig_Get_Miss(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkResult = cfg.Get("noexist")
	}
}

// --- Typed accessors ---

func BenchmarkConfig_String_Hit(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkString = cfg.String("database.host")
	}
}

func BenchmarkConfig_Int_Hit(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkInt = cfg.Int("database.port")
	}
}

func BenchmarkConfig_Bool_Hit(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkBool = cfg.Bool("debug")
	}
}

func BenchmarkConfig_ConfigGet_Float64(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ConfigGet[float64](cfg, "weight")
	}
}

// --- Feature flags ---

func BenchmarkConfig_Enable(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cfg.Enable("dark-mode")
	}
}

func BenchmarkConfig_Disable(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cfg.Disable("dark-mode")
	}
}

func BenchmarkConfig_Enabled_Hit(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkBool = cfg.Enabled("dark-mode")
	}
}

func BenchmarkConfig_Enabled_Miss(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkBool = cfg.Enabled("noexist")
	}
}

func BenchmarkConfig_EnabledFeatures(b *B) {
	cfg := configFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkStrings = cfg.EnabledFeatures()
	}
}

// --- ConfigVar[T] (per-package typed slot) ---

func BenchmarkConfig_NewConfigVar(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewConfigVar[int](42)
	}
}

func BenchmarkConfig_ConfigVar_Get(b *B) {
	v := NewConfigVar[int](42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkInt = v.Get()
	}
}

func BenchmarkConfig_ConfigVar_Set(b *B) {
	v := NewConfigVar[int](42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v.Set(43)
	}
}

func BenchmarkConfig_ConfigVar_IsSet(b *B) {
	v := NewConfigVar[int](42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		configSinkBool = v.IsSet()
	}
}

func BenchmarkConfig_ConfigVar_Unset(b *B) {
	v := NewConfigVar[int](42)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v.Unset()
	}
}
