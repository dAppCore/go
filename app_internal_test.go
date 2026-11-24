// SPDX-License-Identifier: EUPL-1.2

package core

func TestApp_isExecutable_Good(t *T) {
	// Windows has no execute bit: there, executable means a listed extension.
	name := "agent"
	if OS() == "windows" {
		name += ".exe"
	}
	path := Path(t.TempDir(), name)
	RequireTrue(t, WriteFile(path, []byte("#!/bin/sh\n"), 0o755).OK)

	AssertTrue(t, isExecutable(path))
}
func TestApp_isExecutable_Bad(t *T) {
	AssertFalse(t, isExecutable(Path(t.TempDir(), "missing")))
}
func TestApp_isExecutable_Ugly(t *T) {
	AssertFalse(t, isExecutable(t.TempDir()))
}

func TestApp_isExecutableWith_Good(t *T) {
	// The Windows receipt: a 0644 file — no execute bit anywhere — IS
	// executable when its extension is listed. The old mode&0111 test
	// rejected exactly this file, which is why Find never succeeded on
	// the platform.
	path := Path(t.TempDir(), "git.exe")
	RequireTrue(t, WriteFile(path, []byte("MZ"), 0o644).OK)

	AssertTrue(t, isExecutableWith(path, ".COM;.EXE;.BAT;.CMD"))
}
func TestApp_isExecutableWith_Bad(t *T) {
	// An unlisted extension is not executable under PATHEXT semantics,
	// however runnable its mode bits claim it is.
	path := Path(t.TempDir(), "notes.txt")
	RequireTrue(t, WriteFile(path, []byte("x"), 0o755).OK)

	AssertFalse(t, isExecutableWith(path, ".COM;.EXE"))
}
func TestApp_isExecutableWith_Ugly(t *T) {
	// An empty list is the POSIX arm: mode bits decide, so a 0644 file
	// stays non-executable and a directory never qualifies either way.
	path := Path(t.TempDir(), "plain")
	RequireTrue(t, WriteFile(path, []byte("x"), 0o644).OK)

	AssertFalse(t, isExecutableWith(path, ""))
	AssertFalse(t, isExecutableWith(t.TempDir(), ".EXE"))
}

func TestApp_App_findWith_Good(t *T) {
	// Bare name "git" resolves to git.exe on a fixture PATH under
	// PATHEXT semantics — candidate generation plus the extension
	// check working together, from any host platform.
	dir := t.TempDir()
	RequireTrue(t, WriteFile(Path(dir, "git.exe"), []byte("MZ"), 0o644).OK)

	r := App{}.findWith("git", "Git", dir, ".COM;.EXE;.BAT;.CMD")
	RequireTrue(t, r.OK)
	AssertEqual(t, Path(dir, "git.exe"), r.Value.(*App).Path)
}
func TestApp_App_findWith_Bad(t *T) {
	// Nothing on the fixture PATH satisfies the name; and an empty
	// PATH is its own distinct failure.
	AssertFalse(t, App{}.findWith("git", "Git", t.TempDir(), ".EXE").OK)
	AssertFalse(t, App{}.findWith("git", "Git", "", ".EXE").OK)
}
func TestApp_App_findWith_Ugly(t *T) {
	// A forward-slash path is a direct check on every platform — Go
	// accepts "/" on Windows, so "sub/tool" must never be hunted on
	// PATH. The candidate resolves via extension there too.
	dir := t.TempDir()
	RequireTrue(t, MkdirAll(Path(dir, "sub"), 0o755).OK)
	RequireTrue(t, WriteFile(Path(dir, "sub", "tool.exe"), []byte("MZ"), 0o644).OK)

	r := App{}.findWith(Concat(dir, "/sub/tool"), "Tool", "", ".EXE")
	RequireTrue(t, r.OK)
	AssertEqual(t, Path(dir, "sub", "tool.exe"), r.Value.(*App).Path)
}

func TestApp_splitPathExt_Good(t *T) {
	exts := splitPathExt(".COM;.EXE;bat; .Cmd ;")
	AssertEqual(t, 4, len(exts))
	AssertEqual(t, ".com", exts[0])
	AssertEqual(t, ".exe", exts[1])
	AssertEqual(t, ".bat", exts[2])
	AssertEqual(t, ".cmd", exts[3])
}
func TestApp_splitPathExt_Bad(t *T) {
	AssertEqual(t, 0, len(splitPathExt("")))
}
func TestApp_executableCandidates_Good(t *T) {
	// A name already carrying a listed extension is tried as itself
	// first, then with each extension appended; a bare name only with
	// the extensions; POSIX (empty list) is the path alone.
	withExt := executableCandidates("dir/git.exe", ".COM;.EXE")
	AssertEqual(t, 3, len(withExt))
	AssertEqual(t, "dir/git.exe", withExt[0])

	bare := executableCandidates("dir/git", ".COM;.EXE")
	AssertEqual(t, 2, len(bare))
	AssertEqual(t, "dir/git.com", bare[0])
	AssertEqual(t, "dir/git.exe", bare[1])

	posix := executableCandidates("dir/git", "")
	AssertEqual(t, 1, len(posix))
	AssertEqual(t, "dir/git", posix[0])
}

func TestApp_splitPathExt_Ugly(t *T) {
	// Only separators and blanks: nothing survives, and the result is nil.
	AssertNil(t, splitPathExt(" ; ;;"))
}

func TestApp_executableCandidates_Bad(t *T) {
	// A PATHEXT of only separators is no list at all: the path alone, as on POSIX.
	got := executableCandidates("dir/git", ";;")
	AssertEqual(t, []string{"dir/git"}, got)
}

func TestApp_executableCandidates_Ugly(t *T) {
	// An unlisted extension isn't tried bare: "tool.sh" with .EXE only gets ".exe" appended.
	got := executableCandidates("dir/tool.sh", ".EXE")
	AssertEqual(t, []string{"dir/tool.sh.exe"}, got)
}

func TestApp_hasListedExt_Good(t *T) {
	// Case-insensitive against the lower-cased list splitPathExt produces.
	AssertTrue(t, hasListedExt("dir/GIT.EXE", []string{".com", ".exe"}))
}

func TestApp_hasListedExt_Bad(t *T) {
	AssertFalse(t, hasListedExt("dir/git", []string{".exe"}))
	AssertFalse(t, hasListedExt("dir/git.sh", []string{".exe"}))
}

func TestApp_hasListedExt_Ugly(t *T) {
	// Empty list, or a dot-only name: never a match.
	AssertFalse(t, hasListedExt("dir/git.exe", nil))
	AssertFalse(t, hasListedExt("dir/.", []string{"."}))
}
