// SPDX-License-Identifier: EUPL-1.2

//go:build !windows

package core_test

import . "dappco.re/go"

// Example_openFileNoFollow refuses to open a symlinked final path component through
// `O_NOFOLLOW`. The flag carries a real symlink-refusal semantic only on
// unix — it is 0 (a no-op) on Windows, so this example is built only under
// !windows, matching os_nofollow_unix.go's own build tag.
//
// NOTE: this cannot be named ExampleO_NOFOLLOW — go vet's example-association
// parser splits an Example name on its FIRST underscore and treats the head
// as the target identifier, so ExampleO_NOFOLLOW resolves to "O" (unknown)
// and go vet — and `go test`'s automatic pre-test vet subset — hard-fails
// the whole package build with "refers to unknown identifier: O" (same
// failure mode as SHA3_256, confirmed with a standalone repro). Any
// exported identifier containing an underscore is unnameable as a
// symbol-attached Example under Go's own naming grammar; the package-level
// leading-underscore form used here is the only mechanism that stays
// vet-clean.
func Example_openFileNoFollow() {
	dir := TempDir()
	target := PathJoin(dir, "real.log")
	link := PathJoin(dir, "current.log")
	WriteFile(target, []byte("ready"), 0o600)
	Symlink(target, link)
	defer Remove(target)
	defer Remove(link)

	r := OpenFile(link, O_CREATE|O_EXCL|O_NOFOLLOW|O_WRONLY, 0o600)
	Println(r.OK)
	// Output: false
}
