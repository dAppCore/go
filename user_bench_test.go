// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the user lookups in user.go.
// Per AX-11 — UserCurrent fires on every ~ expansion in Fs path
// resolution + every log identity tag. UserLookup / UserLookupID /
// UserGroupLookup hit /etc/passwd or NSS — slower; the bench gates
// the cost so a Core reroute can't quietly add lookups per-call.
//
// Run:    go test -bench='BenchmarkUser' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var userSinkResult Result

// --- Current user (cached by os/user) ---

func BenchmarkUser_Current(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		userSinkResult = UserCurrent()
	}
}

// --- Lookup by username ---

func BenchmarkUser_Lookup_Hit(b *B) {
	// "root" exists on every POSIX dev box.
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		userSinkResult = UserLookup("root")
	}
}

func BenchmarkUser_Lookup_Miss(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		userSinkResult = UserLookup("__core_no_such_user__")
	}
}

// --- Lookup by uid ---

func BenchmarkUser_LookupID_Hit(b *B) {
	// uid 0 is root on every POSIX system.
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		userSinkResult = UserLookupID("0")
	}
}

// --- Group lookup ---

func BenchmarkUser_GroupLookup_Miss(b *B) {
	// Use a guaranteed-miss name; group names vary across platforms.
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		userSinkResult = UserGroupLookup("__core_no_such_group__")
	}
}
