// SPDX-License-Identifier: EUPL-1.2

//go:build !windows

package core_test

import . "dappco.re/go"

// O_NOFOLLOW carries a symlink-refusal open semantic only on unix (it is
// 0 on Windows), so its Good/Bad/Ugly triplet lives in a !windows test
// file alongside the unix declaration in os_nofollow_unix.go.

// Good — opening a regular (non-symlink) final component with O_NOFOLLOW
// succeeds exactly as a normal create-or-write would.
func TestOs_OpenFile_NOFOLLOW_Good(t *T) {
	path := Path(t.TempDir(), "agent.log")

	r := OpenFile(path, O_CREATE|O_EXCL|O_NOFOLLOW|O_WRONLY, 0o600)
	RequireTrue(t, r.OK)
	defer CloseStream(r.Value)

	AssertTrue(t, WriteAll(r.Value, "ready").OK)
}

// Bad — the final component is a symlink, so O_NOFOLLOW refuses the open
// and the Result reports the failure.
func TestOs_OpenFile_NOFOLLOW_Bad(t *T) {
	dir := t.TempDir()
	target := Path(dir, "real.log")
	link := Path(dir, "current.log")
	AssertTrue(t, WriteFile(target, []byte("ready"), 0o600).OK)
	RequireNoError(t, SymlinkForTest(target, link))

	r := OpenFile(link, O_CREATE|O_EXCL|O_NOFOLLOW|O_WRONLY, 0o600)

	AssertFalse(t, r.OK)
}

// Ugly — a dangling symlink (target does not exist) opened for read with
// O_NOFOLLOW still fails: the refusal fires on the link itself, not the
// resolved destination.
func TestOs_OpenFile_NOFOLLOW_Ugly(t *T) {
	dir := t.TempDir()
	link := Path(dir, "dangling.log")
	RequireNoError(t, SymlinkForTest(Path(dir, "nowhere.log"), link))

	r := OpenFile(link, O_RDONLY|O_NOFOLLOW, 0o600)

	AssertFalse(t, r.OK)
}
