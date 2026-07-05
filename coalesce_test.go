// SPDX-License-Identifier: EUPL-1.2

package core_test

import (
	. "dappco.re/go"
)

func TestCoalesce_Coalesce_Good(t *T) {
	AssertEqual(t, "third", Coalesce("", "", "third"))
}

func TestCoalesce_Coalesce_Bad(t *T) {
	AssertEqual(t, "", Coalesce("", ""))
}

func TestCoalesce_Coalesce_Ugly(t *T) {
	// first non-zero wins even when later values are also non-zero
	AssertEqual(t, 7, Coalesce(0, 0, 7, 9))
}

func TestCoalesce_FirstPositive_Good(t *T) {
	AssertEqual(t, 16, FirstPositive(0, 0, 16))
}

func TestCoalesce_FirstPositive_Bad(t *T) {
	AssertEqual(t, 0, FirstPositive(0, 0))
}

func TestCoalesce_FirstPositive_Ugly(t *T) {
	// negatives are skipped (unlike Coalesce, which treats a negative as non-zero)
	AssertEqual(t, 5, FirstPositive(-1, 0, 5))
}
