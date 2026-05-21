// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the IPC primitives in ipc.go.
// Per AX-11 — IPC is the message bus underpinning ACTION broadcast,
// QUERY first-OK-wins, and QUERYALL gather. Every service-startup /
// dispatch / progress / completed event flows through here. The bench
// gates the per-event floor in both the no-handler and N-handler shapes.
//
// Run:    go test -bench='BenchmarkIPC' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var ipcSinkResult Result

// benchMessage is a minimal Message + Query implementation for the bus.
type benchMessage struct{ id int }

// benchQuery satisfies the Query iface (any non-Message any works —
// the bus is just an interface{}-passing surface).
type benchQuery struct{ id int }

// noopHandler returns OK without doing work.
func noopHandler() func(*Core, Message) Result {
	return func(_ *Core, _ Message) Result { return Result{OK: true} }
}

// noopQuery returns OK with a small Value.
func noopQueryHandler() QueryHandler {
	return func(_ *Core, _ Query) Result { return Result{Value: 42, OK: true} }
}

// --- ACTION (broadcast) ---

func BenchmarkIPC_ACTION_NoHandlers(b *B) {
	c := New()
	msg := benchMessage{id: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ipcSinkResult = c.ACTION(msg)
	}
}

func BenchmarkIPC_ACTION_OneHandler(b *B) {
	c := New()
	c.RegisterAction(noopHandler())
	msg := benchMessage{id: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ipcSinkResult = c.ACTION(msg)
	}
}

func BenchmarkIPC_ACTION_TenHandlers(b *B) {
	c := New()
	for i := 0; i < 10; i++ {
		c.RegisterAction(noopHandler())
	}
	msg := benchMessage{id: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ipcSinkResult = c.ACTION(msg)
	}
}

// --- QUERY (first OK wins) ---

func BenchmarkIPC_QUERY_NoHandlers(b *B) {
	c := New()
	q := benchQuery{id: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ipcSinkResult = c.QUERY(q)
	}
}

func BenchmarkIPC_QUERY_OneHandler(b *B) {
	c := New()
	c.RegisterQuery(noopQueryHandler())
	q := benchQuery{id: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ipcSinkResult = c.QUERY(q)
	}
}

// --- QUERYALL ---

func BenchmarkIPC_QUERYALL_TenHandlers(b *B) {
	c := New()
	for i := 0; i < 10; i++ {
		c.RegisterQuery(noopQueryHandler())
	}
	q := benchQuery{id: 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ipcSinkResult = c.QUERYALL(q)
	}
}

// --- Register* ---

func BenchmarkIPC_RegisterAction(b *B) {
	c := New()
	handler := noopHandler()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.RegisterAction(handler)
	}
}

func BenchmarkIPC_RegisterActions_Five(b *B) {
	c := New()
	hs := []func(*Core, Message) Result{
		noopHandler(), noopHandler(), noopHandler(), noopHandler(), noopHandler(),
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.RegisterActions(hs...)
	}
}

func BenchmarkIPC_RegisterQuery(b *B) {
	c := New()
	q := noopQueryHandler()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.RegisterQuery(q)
	}
}
