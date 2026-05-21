// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the Core primitive in core.go.
// Per AX-11 — *Core sits at the centre of every consumer. Each accessor
// gates a subsystem the rest of the framework defers through (Fs,
// Config, Drive, Data, App, IPC, etc.) and is read inline on the hot
// path. New() is one-shot per process but worth gating against
// regressions in the constructor surface.
//
// Run:    go test -bench='BenchmarkCore' -benchmem -run='^$' .

package core_test

import (
	"io"

	. "dappco.re/go"
)

// quietDefault redirects the default logger to io.Discard for the duration
// of a sub-test. Used by LogError / LogWarn benches that otherwise spam
// stderr with formatted log lines.
func quietDefault(b *B) {
	prev := Default()
	SetDefault(NewLog(LogOptions{Level: LevelError, Output: io.Discard}))
	b.Cleanup(func() { SetDefault(prev) })
}

// Sinks defeat compiler DCE.
var (
	coreSinkCore     *Core
	coreSinkOptions  *Options
	coreSinkApp      *App
	coreSinkData     *Data
	coreSinkDrive    *Drive
	coreSinkFs       *Fs
	coreSinkConfig   *Config
	coreSinkErrPanic *ErrorPanic
	coreSinkErrLog   *ErrorLog
	coreSinkIpc      *Ipc
	coreSinkI18n     *I18n
	coreSinkEnv      string
	coreSinkCtx      Context
	coreSinkRegistry *Registry[any]
	coreSinkResult   Result
)

// --- Constructor ---

func BenchmarkCore_New(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkCore = New()
	}
}

// --- Accessors (the inline hot path) ---

func BenchmarkCore_Options(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkOptions = c.Options()
	}
}

func BenchmarkCore_App(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkApp = c.App()
	}
}

func BenchmarkCore_Data(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkData = c.Data()
	}
}

func BenchmarkCore_Drive(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkDrive = c.Drive()
	}
}

func BenchmarkCore_Fs(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkFs = c.Fs()
	}
}

func BenchmarkCore_Config(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkConfig = c.Config()
	}
}

func BenchmarkCore_Error(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkErrPanic = c.Error()
	}
}

func BenchmarkCore_Log(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkErrLog = c.Log()
	}
}

func BenchmarkCore_IPC(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkIpc = c.IPC()
	}
}

func BenchmarkCore_I18n(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkI18n = c.I18n()
	}
}

func BenchmarkCore_Env(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkEnv = c.Env("DIR_HOME")
	}
}

func BenchmarkCore_Context(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkCtx = c.Context()
	}
}

func BenchmarkCore_Core_Self(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkCore = c.Core()
	}
}

// --- IPC convenience aliases ---

func BenchmarkCore_ACTION_NoHandlers(b *B) {
	c := New()
	msg := ActionServiceStartup{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkResult = c.ACTION(msg)
	}
}

// --- LogError / LogWarn ---

func BenchmarkCore_LogError(b *B) {
	quietDefault(b)
	c := New()
	err := NewError("bench.LogError")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkResult = c.LogError(err, "bench.op", "synthetic")
	}
}

func BenchmarkCore_LogWarn(b *B) {
	quietDefault(b)
	c := New()
	err := NewError("bench.LogWarn")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkResult = c.LogWarn(err, "bench.op", "synthetic")
	}
}

// --- RegistryOf ---

func BenchmarkCore_RegistryOf_Services(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkRegistry = c.RegistryOf("services")
	}
}

func BenchmarkCore_RegistryOf_Actions(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkRegistry = c.RegistryOf("actions")
	}
}

func BenchmarkCore_RegistryOf_Unknown(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		coreSinkRegistry = c.RegistryOf("noexist")
	}
}
