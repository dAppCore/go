// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the SQL primitives in sql.go.
// Per AX-11 — sql.go is overwhelmingly type aliases (DB, Tx, Rows,
// Stmt, Row, SQLResult, NullX). The three callable funcs are: SQLOpen
// (driver-bind path), SQLDrivers (diagnostic list), and SQLIsNoRows
// (the canonical not-found check). core/go itself imports no SQL
// drivers — SQLOpen with an unregistered driver returns the no-driver
// error path, which is the floor consumers pay when a driver is
// missing.
//
// Run:    go test -bench='BenchmarkSQL' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	sqlSinkResult  Result
	sqlSinkStrings []string
	sqlSinkBool    bool
)

// --- SQLOpen ---

func BenchmarkSQL_Open_NoDriver(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sqlSinkResult = SQLOpen("nodriver", "anything")
	}
}

// --- SQLDrivers ---

func BenchmarkSQL_Drivers(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sqlSinkStrings = SQLDrivers()
	}
}

// --- SQLIsNoRows ---

func BenchmarkSQL_IsNoRows_Hit(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sqlSinkBool = SQLIsNoRows(ErrNoRows)
	}
}

func BenchmarkSQL_IsNoRows_Miss(b *B) {
	err := NewError("different error")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sqlSinkBool = SQLIsNoRows(err)
	}
}

func BenchmarkSQL_IsNoRows_Wrapped(b *B) {
	wrapped := Wrap(ErrNoRows, "bench.SQL", "wrapped not-found")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sqlSinkBool = SQLIsNoRows(wrapped)
	}
}
