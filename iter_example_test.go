// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func ExamplePull() {
	var seq Seq[int] = func(yield func(int) bool) {
		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	}

	next, stop := Pull(seq)
	defer stop()
	for {
		v, ok := next()
		if !ok {
			break
		}
		Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExamplePull2() {
	var seq Seq2[string, int] = func(yield func(string, int) bool) {
		_ = yield("a", 1) && yield("b", 2)
	}

	next, stop := Pull2(seq)
	defer stop()
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		Println(k, v)
	}
	// Output:
	// a 1
	// b 2
}
