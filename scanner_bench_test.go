// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the buffered-reader / scanner primitives in scanner.go.
// Per AX-11 — NewBufReader is the byte-stream JSON-RPC frame parser
// in lsp.go; NewLineScanner is the line-oriented reader for log
// tailers and config parsers. Both fire on consumer hot paths whenever
// stdin / RPC / file streams are processed.
//
// Run:    go test -bench='BenchmarkScanner' -benchmem -run='^$' .

package core_test

import (
	"bufio"
	"bytes"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	scannerSinkReader  *BufReader
	scannerSinkScanner *bufio.Scanner
	scannerSinkString  string
	scannerSinkInt     int
)

// scannerCorpus returns a 1024-line fixture (~32KB) for the scan loop.
func scannerCorpus() []byte {
	var buf bytes.Buffer
	for range 1024 {
		buf.WriteString("the quick brown fox jumps over the lazy dog\n")
	}
	return buf.Bytes()
}

// --- NewBufReader ---

func BenchmarkScanner_NewBufReader(b *B) {
	corpus := scannerCorpus()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scannerSinkReader = NewBufReader(bytes.NewReader(corpus))
	}
}

func BenchmarkScanner_BufReader_ReadString_1Line(b *B) {
	corpus := scannerCorpus()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewBufReader(bytes.NewReader(corpus))
		scannerSinkString, _ = r.ReadString('\n')
	}
}

// --- NewLineScanner ---

func BenchmarkScanner_NewLineScanner(b *B) {
	corpus := scannerCorpus()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scannerSinkScanner = NewLineScanner(bytes.NewReader(corpus))
	}
}

func BenchmarkScanner_NewLineScannerWithSize(b *B) {
	corpus := scannerCorpus()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scannerSinkScanner = NewLineScannerWithSize(bytes.NewReader(corpus), 2*1024*1024)
	}
}

func BenchmarkScanner_LineScanner_Scan_1024(b *B) {
	corpus := scannerCorpus()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := NewLineScanner(bytes.NewReader(corpus))
		count := 0
		for s.Scan() {
			count++
		}
		scannerSinkInt = count
	}
}
