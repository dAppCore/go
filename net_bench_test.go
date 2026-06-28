// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the net primitives in net.go.
// Per AX-11 — most of net.go does real network I/O (Dial / Listen)
// which is IO-bound and untestable in a deterministic bench. The
// pure-string parsers (ParseIP / ParseCIDR) and NetPipe (no kernel
// network stack) are the parts a unit bench can usefully gate.
//
// Run:    go test -bench='BenchmarkNet' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	netSinkIP     IP
	netSinkResult Result
)

// --- ParseIP ---

func BenchmarkNet_ParseIP_v4(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		netSinkIP = ParseIP("192.168.1.1")
	}
}

func BenchmarkNet_ParseIP_v6(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		netSinkIP = ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	}
}

func BenchmarkNet_ParseIP_v6Compact(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		netSinkIP = ParseIP("::1")
	}
}

// --- ParseCIDR ---

func BenchmarkNet_ParseCIDR_v4(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		netSinkResult = ParseCIDR("192.168.1.0/24")
	}
}

func BenchmarkNet_ParseCIDR_v6(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		netSinkResult = ParseCIDR("2001:db8::/32")
	}
}

// --- NetPipe ---

func BenchmarkNet_NetPipe(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, c := NetPipe()
		a.Close()
		c.Close()
	}
}

func BenchmarkNetListen(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NetListen("tcp", "127.0.0.1:0")
		if r.OK {
			r.Value.(Listener).Close()
		}
	}
}

func BenchmarkNetListenPacket(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NetListenPacket("udp", "127.0.0.1:0")
		if r.OK {
			r.Value.(PacketConn).Close()
		}
	}
}

// netDialFixture starts a loopback listener with an accept-and-close loop
// so the Dial benches have a server that drains the backlog on full runs.
func netDialFixture(b *B) string {
	ln := NetListen("tcp", "127.0.0.1:0").Value.(Listener)
	b.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().String()
}

func BenchmarkNetDial(b *B) {
	addr := netDialFixture(b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NetDial("tcp", addr)
		if r.OK {
			r.Value.(Conn).Close()
		}
	}
}

func BenchmarkNetDialTimeout(b *B) {
	addr := netDialFixture(b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NetDialTimeout("tcp", addr, Second)
		if r.OK {
			r.Value.(Conn).Close()
		}
	}
}
