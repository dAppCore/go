// SPDX-License-Identifier: EUPL-1.2

package core

import "os"

// AttachParentConsole and useForStdio swap the process's stdout, which only
// internal code can watch and put back; these tests live here for that.

// Bad — os.Stdout has been pointed at a closed file. AttachParentConsole
// reads the process's real stdout handle, not os.Stdout, so it still sees
// go test's pipe, attaches nothing and leaves os.Stdout alone.
func TestConsole_AttachParentConsole_Bad(t *T) {
	closed, err := os.CreateTemp(t.TempDir(), "closed")
	RequireNoError(t, err)
	RequireNoError(t, closed.Close())
	prev := os.Stdout
	os.Stdout = closed
	defer func() { os.Stdout = prev }()

	r := AttachParentConsole()

	RequireTrue(t, r.OK)
	AssertEqual(t, false, r.Value)
	AssertTrue(t, os.Stdout == closed)
}

// Good — an open file becomes stdout and stderr, and comes back as the
// Writer for the logger to follow.
func TestOs_useForStdio_Good(t *T) {
	f, err := os.CreateTemp(t.TempDir(), "console")
	RequireNoError(t, err)
	defer f.Close()
	prevOut, prevErr := os.Stdout, os.Stderr
	defer func() { os.Stdout, os.Stderr = prevOut, prevErr }()

	w, ok := useForStdio(f)

	RequireTrue(t, ok)
	AssertTrue(t, os.Stdout == f && os.Stderr == f)
	AssertTrue(t, w == Writer(f))
}

// Bad — something that isn't a file is refused and stdout is untouched.
func TestOs_useForStdio_Bad(t *T) {
	prev := os.Stdout

	w, ok := useForStdio(NewBuffer())

	AssertFalse(t, ok)
	AssertNil(t, w)
	AssertTrue(t, os.Stdout == prev)
}

// Ugly — nil (OpenFile's Value on a failure) is refused the same way.
func TestOs_useForStdio_Ugly(t *T) {
	prev := os.Stdout

	_, ok := useForStdio(nil)

	AssertFalse(t, ok)
	AssertTrue(t, os.Stdout == prev)
}
