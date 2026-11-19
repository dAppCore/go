// SPDX-License-Identifier: EUPL-1.2

package core

import "os"

// Println writes straight to os.Stdout, so measuring it means swapping the
// process's stdout for /dev/null — only internal code can (the bench lives
// here, exempt from the test-imports rule, for that one reason). The bench
// then measures the formatting + write cost, not terminal scroll.
func BenchmarkPrintln(b *B) {
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Skip("no /dev/null")
	}
	old := os.Stdout
	os.Stdout = devnull
	defer func() {
		os.Stdout = old
		devnull.Close()
	}()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Println("agent ready", 42)
	}
}
