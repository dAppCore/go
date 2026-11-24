// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "github.com/dappcore/go"

// Good — under go test stdout is already a pipe or file, so nothing is
// attached and stdout is left exactly as it was (on every OS).
func TestConsole_AttachParentConsole_Good(t *T) {
	before := Stdout()

	r := AttachParentConsole()

	RequireTrue(t, r.OK)
	AssertEqual(t, false, r.Value)
	AssertTrue(t, Stdout() == before)
}

// Ugly — calling it again is harmless: still OK, still nothing attached.
func TestConsole_AttachParentConsole_Ugly(t *T) {
	AttachParentConsole()

	r := AttachParentConsole()

	RequireTrue(t, r.OK)
	AssertEqual(t, false, r.Value)
}
