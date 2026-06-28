// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func TestPath_Relative(t *T) {
	home := Env("DIR_HOME")

	ds := Env("DS")
	AssertEqual(t, home+ds+"Code"+ds+".core", Path("Code", ".core"))
}

func TestPath_Absolute(t *T) {
	ds := Env("DS")
	AssertEqual(t, "/tmp"+ds+"workspace", Path("/tmp", "workspace"))
	// An absolute first segment bypasses the DIR_HOME anchor and flows
	// through PathJoin, which cleans while joining: .. resolves and
	// redundant separators collapse.
	AssertEqual(t, "/tmp", Path("/tmp", "workspace", ".."))
	AssertEqual(t, "/a"+ds+"b", Path("/a//b"))
	AssertEqual(t, "/etc", Path("/tmp/../etc"))
	AssertEqual(t, "/usr"+ds+"local", Path("/usr/local"))
	AssertTrue(t, PathIsAbs(Path("/tmp", "workspace")))
}

func TestPath_Empty(t *T) {
	home := Env("DIR_HOME")

	AssertEqual(t, home, Path())
	// With no segments Path returns the home anchor verbatim: non-empty,
	// absolute, and with no trailing separator.
	AssertNotEmpty(t, Path())
	AssertTrue(t, PathIsAbs(Path()))
	AssertFalse(t, HasSuffix(Path(), Env("DS")))
}

func TestPath_Cleans(t *T) {
	home := Env("DIR_HOME")
	ds := Env("DS")

	AssertEqual(t, home+ds+"Code", Path("Code", "sub", ".."))
	// .. unwinds segment-by-segment, even back to the bare home anchor.
	AssertEqual(t, home, Path("a", "b", "..", ".."))
	// A lone "." segment is dropped during cleaning.
	AssertEqual(t, home+ds+"a"+ds+"b", Path("a", ".", "b"))
	// Mixed redundant elements collapse to a single clean path.
	AssertEqual(t, home+ds+"x", Path("x", "y", "z", "..", ".."))
}

func TestPath_CleanDoubleSlash(t *T) {
	ds := Env("DS")
	AssertEqual(t, ds+"tmp"+ds+"file", Path("/tmp//file"))
	// Repeated and leading-doubled separators all collapse to one.
	AssertEqual(t, ds+"a"+ds+"b"+ds+"c", Path("/a///b//c"))
	AssertEqual(t, ds+"x"+ds+"y", Path("/x/", "/y"))
}

func TestPath_PathBase(t *T) {
	AssertEqual(t, "core", PathBase("/Users/snider/Code/core"))
	AssertEqual(t, "homelab", PathBase("deploy/to/homelab"))
	// A trailing separator is trimmed before the last element is taken.
	AssertEqual(t, "b", PathBase("a/b/"))
	// A bare name with no separator is returned unchanged.
	AssertEqual(t, "file.go", PathBase("file.go"))
	// A leading-separator single element keeps just the name.
	AssertEqual(t, "a", PathBase("/a"))
}

func TestPath_PathBase_Root(t *T) {
	// Trimming "/" off "/" leaves "", which PathBase reports as the
	// separator itself.
	AssertEqual(t, "/", PathBase("/"))
	AssertEqual(t, string(PathSeparator), PathBase(string(PathSeparator)))
	// An absolute single element with a trailing slash trims back to the name.
	AssertEqual(t, "a", PathBase("/a/"))
}

func TestPath_PathBase_Empty(t *T) {
	AssertEqual(t, ".", PathBase(""))
	// "." has no separator, so it is returned verbatim rather than
	// re-mapped to the empty-string special case.
	AssertEqual(t, ".", PathBase("."))
	// The parent marker has no separator either and is returned as-is.
	AssertEqual(t, "..", PathBase(".."))
}

func TestPath_PathDir(t *T) {
	AssertEqual(t, "/Users/snider/Code", PathDir("/Users/snider/Code/core"))
	// Drops only the final element; intermediate elements are preserved.
	AssertEqual(t, "a/b", PathDir("a/b/c"))
	AssertEqual(t, "/srv/dappcore", PathDir("/srv/dappcore/agent.json"))
}

func TestPath_PathDir_Root(t *T) {
	// The only separator is at index 0, so the dir collapses to "/".
	AssertEqual(t, "/", PathDir("/file"))
	AssertEqual(t, "/", PathDir("/agent.json"))
	// Collapse to root only happens for a single element; deeper paths keep
	// their leading directory.
	AssertEqual(t, "/a", PathDir("/a/b"))
}

func TestPath_PathDir_NoDir(t *T) {
	// No separator at all means there is no directory component.
	AssertEqual(t, ".", PathDir("file.go"))
	AssertEqual(t, ".", PathDir("agent"))
	AssertEqual(t, ".", PathDir(""))
}

func TestPath_PathExt(t *T) {
	AssertEqual(t, ".go", PathExt("main.go"))
	AssertEqual(t, "", PathExt("Makefile"))
	AssertEqual(t, ".gz", PathExt("archive.tar.gz"))
}

func TestPath_EnvConsistency(t *T) {
	// Path() with no args is exactly the DIR_HOME anchor, and the DS env
	// key is the separator Path uses to assemble anchored paths.
	AssertEqual(t, Env("DIR_HOME"), Path())
	ds := Env("DS")
	AssertEqual(t, string(PathSeparator), ds)
	AssertEqual(t, Env("DIR_HOME")+ds+"Code", Path("Code"))
}

func TestPath_PathGlob_Good(t *T) {
	dir := t.TempDir()
	f := (&Fs{}).New("/")
	f.Write(Path(dir, "a.txt"), "a")
	f.Write(Path(dir, "b.txt"), "b")
	f.Write(Path(dir, "c.log"), "c")

	matches := PathGlob(Path(dir, "*.txt"))
	AssertLen(t, matches, 2)
}

func TestPath_PathGlob_NoMatch(t *T) {
	// A syntactically valid pattern pointing at a missing tree yields no
	// matches but is not an error.
	matches := PathGlob("/nonexistent/pattern-*.xyz")
	AssertEmpty(t, matches)
	// A real, populated directory whose pattern matches nothing also
	// returns an empty slice.
	dir := t.TempDir()
	f := (&Fs{}).New("/")
	f.Write(Path(dir, "agent.log"), "x")
	AssertEmpty(t, PathGlob(Path(dir, "*.txt")))
}

func TestPath_PathIsAbs_Good(t *T) {
	AssertTrue(t, PathIsAbs("/tmp"))
	AssertTrue(t, PathIsAbs("/"))
	AssertFalse(t, PathIsAbs("relative"))
	AssertFalse(t, PathIsAbs(""))
}

func TestPath_CleanPath_Good(t *T) {
	AssertEqual(t, "/a/b", CleanPath("/a//b", "/"))
	AssertEqual(t, "/a/c", CleanPath("/a/b/../c", "/"))
	AssertEqual(t, "/", CleanPath("/", "/"))
	AssertEqual(t, ".", CleanPath("", "/"))
}

func TestPath_PathDir_TrailingSlash(t *T) {
	// PathDir does not pre-trim a trailing separator, so the final empty
	// element is dropped, leaving the path without its trailing slash.
	result := PathDir("/Users/snider/Code/")
	AssertEqual(t, "/Users/snider/Code", result)
	AssertEqual(t, "a/b", PathDir("a/b/"))
}

// --- PathRel ---

func TestPath_PathRel_Good_Descendant(t *T) {
	r := PathRel("/var/lib/foo", "/var/lib/foo/bar/baz")
	AssertTrue(t, r.OK)
	AssertEqual(t, "bar/baz", r.Value.(string))
}

func TestPath_PathRel_Good_Sibling(t *T) {
	r := PathRel("/a", "/b")
	AssertTrue(t, r.OK)
	AssertEqual(t, "../b", r.Value.(string))
}

func TestPath_PathRel_Good_Identical(t *T) {
	r := PathRel("/x/y", "/x/y")
	AssertTrue(t, r.OK)
	AssertEqual(t, ".", r.Value.(string))
}

func TestPath_PathRel_Bad_MixedAbsRel(t *T) {
	// filepath.Rel rejects when one is absolute and the other relative,
	// in EITHER direction. On failure, OK is false and Value carries the
	// underlying error rather than a string.
	r := PathRel("/abs/path", "rel/path")
	AssertFalse(t, r.OK)
	_, isStr := r.Value.(string)
	AssertFalse(t, isStr)
	AssertNotNil(t, r.Value)
	// Reverse pairing (relative base, absolute target) is also rejected.
	r2 := PathRel("rel/path", "/abs/path")
	AssertFalse(t, r2.OK)
}

// --- PathAbs ---

func TestPath_PathAbs_Good_AlreadyAbsolute(t *T) {
	r := PathAbs("/already/absolute")
	AssertTrue(t, r.OK)
	AssertEqual(t, "/already/absolute", r.Value.(string))
}

func TestPath_PathAbs_Good_Relative(t *T) {
	r := PathAbs("./relative/path")
	AssertTrue(t, r.OK)
	abs := r.Value.(string)
	// Resolved against cwd — must end in the relative tail and be absolute.
	AssertTrue(t, len(abs) > 0 && abs[0] == '/', "PathAbs must return absolute: %s", abs)
	AssertContains(t, abs, "relative/path")
}

// --- AX-7 canonical triplets ---

func TestPath_Path_Good(t *T) {
	home := Env("DIR_HOME")
	ds := Env("DS")
	AssertEqual(t, home+ds+"Code"+ds+"core", Path("Code", "core"))
}

func TestPath_Path_Bad(t *T) {
	// Degenerate input (no segments) still yields a usable absolute anchor
	// rather than an empty string or "."
	AssertEqual(t, Env("DIR_HOME"), Path())
	AssertTrue(t, PathIsAbs(Path()))
	AssertNotEmpty(t, Path())
	AssertNotEqual(t, ".", Path())
}

func TestPath_Path_Ugly(t *T) {
	home := Env("DIR_HOME")
	ds := Env("DS")
	AssertEqual(t, home+ds+"agent", Path("workspace", "..", "agent"))
}

func TestPath_PathJoin_Good(t *T) {
	AssertEqual(t, "deploy/to/homelab", PathToSlash(PathJoin("deploy", "to", "homelab")))
	// Unlike Path, PathJoin preserves a relative result (no home anchor).
	AssertFalse(t, PathIsAbs(PathJoin("deploy", "to", "homelab")))
	// An absolute first segment is kept absolute.
	AssertEqual(t, "/a/b", PathToSlash(PathJoin("/a", "b")))
	// Empty interior segments are skipped.
	AssertEqual(t, "a/b", PathToSlash(PathJoin("a", "", "b")))
}

func TestPath_PathJoin_Bad(t *T) {
	// No segments, or only empty segments, both join to the empty string.
	AssertEqual(t, "", PathJoin())
	AssertEqual(t, "", PathJoin(""))
	AssertEqual(t, "", PathJoin("", ""))
}

func TestPath_PathJoin_Ugly(t *T) {
	// .. is resolved during the join.
	AssertEqual(t, "agent", PathJoin("deploy", "..", "agent"))
	// .. may unwind past all named segments into a parent reference.
	AssertEqual(t, "..", PathJoin("a", "..", ".."))
	// Redundant separators within segments collapse.
	AssertEqual(t, "a/b", PathToSlash(PathJoin("a/", "/b")))
}

func TestPath_PathBase_Good(t *T) {
	AssertEqual(t, "agent.json", PathBase("/srv/dappcore/agent.json"))
	// Relative paths take their last element too.
	AssertEqual(t, "homelab", PathBase("deploy/to/homelab"))
	// A trailing separator is trimmed before the base is read.
	AssertEqual(t, "core", PathBase("/Users/snider/Code/core/"))
}

func TestPath_PathBase_Bad(t *T) {
	// The empty string maps to the current-directory marker.
	AssertEqual(t, ".", PathBase(""))
	// "." has no separator and is returned as-is.
	AssertEqual(t, ".", PathBase("."))
	// A bare name with no path structure is returned unchanged.
	AssertEqual(t, "a", PathBase("a"))
}

func TestPath_PathBase_Ugly(t *T) {
	// The bare separator trims to "", which reports back as the separator.
	AssertEqual(t, string(PathSeparator), PathBase(string(PathSeparator)))
	// A name with no directory component is returned verbatim.
	AssertEqual(t, "file.go", PathBase("file.go"))
	// A trailing separator on a multi-element path is trimmed first.
	AssertEqual(t, "b", PathBase("a/b/"))
}

func TestPath_PathDir_Good(t *T) {
	AssertEqual(t, "/srv/dappcore", PathDir("/srv/dappcore/agent.json"))
	// Deeper nesting keeps every element except the last.
	AssertEqual(t, "/srv/dappcore/agent", PathDir("/srv/dappcore/agent/config.json"))
	// Relative input keeps its relative directory.
	AssertEqual(t, "deploy/to", PathDir("deploy/to/homelab"))
}

func TestPath_PathDir_Bad(t *T) {
	// No separator: there is no directory part.
	AssertEqual(t, ".", PathDir("agent.json"))
	AssertEqual(t, ".", PathDir(""))
	AssertEqual(t, ".", PathDir("agent"))
}

func TestPath_PathDir_Ugly(t *T) {
	// The only separator sits at index 0, so the dir collapses to it.
	AssertEqual(t, string(PathSeparator), PathDir(string(PathSeparator)+"agent.json"))
	// A trailing separator drops the empty final element.
	AssertEqual(t, "a/b", PathDir("a/b/"))
	// A single relative element with a trailing slash keeps the element as
	// the dir (the custom impl differs from stdlib filepath.Dir here).
	AssertEqual(t, "a", PathDir("a/"))
}

func TestPath_PathExt_Good(t *T) {
	// Only the final dotted suffix is the extension.
	AssertEqual(t, ".json", PathExt("agent.config.json"))
	AssertEqual(t, ".go", PathExt("main.go"))
	// The extension is read from the base element, ignoring directories.
	AssertEqual(t, ".txt", PathExt("/srv/dappcore/notes.txt"))
}

func TestPath_PathExt_Bad(t *T) {
	// No dot anywhere in the base means no extension.
	AssertEqual(t, "", PathExt("Makefile"))
	AssertEqual(t, "", PathExt(""))
	// A dot in a directory element does not count as the file's extension.
	AssertEqual(t, "", PathExt("dir.d/file"))
}

func TestPath_PathExt_Ugly(t *T) {
	// A leading dot is the start of the name, not an extension separator
	// (lastIndex is 0, and i<=0 is treated as "no extension").
	AssertEqual(t, "", PathExt(".env"))
	AssertEqual(t, "", PathExt(".gitignore"))
	// But a dotfile that ALSO has a later dot does have an extension.
	AssertEqual(t, ".local", PathExt(".env.local"))
}

func TestPath_PathIsAbs_Bad(t *T) {
	// Plain relative paths and dot-relative prefixes are not absolute.
	AssertFalse(t, PathIsAbs("tmp/agent"))
	AssertFalse(t, PathIsAbs("./tmp"))
	AssertFalse(t, PathIsAbs("../tmp"))
	AssertFalse(t, PathIsAbs("agent"))
}

func TestPath_PathIsAbs_Ugly(t *T) {
	// Windows drive prefixes count as absolute with either slash flavour.
	AssertTrue(t, PathIsAbs(`C:\agent\workspace`))
	AssertTrue(t, PathIsAbs("C:/agent/workspace"))
	// A drive letter WITHOUT a following separator is not absolute.
	AssertFalse(t, PathIsAbs("C:agent"))
	// The empty string is never absolute.
	AssertFalse(t, PathIsAbs(""))
}

func TestPath_CleanPath_Bad(t *T) {
	// The empty path is the degenerate input; it cleans to ".".
	AssertEqual(t, ".", CleanPath("", "/"))
	// A path that fully unwinds via .. also reduces to a relative marker.
	AssertEqual(t, "..", CleanPath("foo/../..", "/"))
	// Balanced descent and ascent collapses all the way to ".".
	AssertEqual(t, ".", CleanPath("a/b/../..", "/"))
}

func TestPath_CleanPath_Ugly(t *T) {
	// .. cannot escape an absolute root; the leading climbs are discarded.
	AssertEqual(t, "/agent", CleanPath("/../../agent", "/"))
	// A non-native separator exercises the manual Split/Join algorithm
	// rather than the filepath.Clean fast path.
	AssertEqual(t, "a|c", CleanPath("a|b|..|c", "|"))
	// The rooted-with-.. rule also holds under a custom separator: a leading
	// climb out of the root is discarded.
	AssertEqual(t, "|a", CleanPath("|..|a", "|"))
}

func TestPath_PathGlob_Bad(t *T) {
	// A malformed pattern makes filepath.Glob error; PathGlob swallows it
	// and returns no matches rather than panicking.
	AssertEmpty(t, PathGlob("["))
	AssertEmpty(t, PathGlob("[a-"))
	AssertEmpty(t, PathGlob("a[b"))
}

func TestPath_PathGlob_Ugly(t *T) {
	dir := t.TempDir()
	f := (&Fs{}).New("/")
	f.Write(Path(dir, "agent one.txt"), "one")
	matches := PathGlob(Path(dir, "agent*.txt"))
	AssertLen(t, matches, 1)
	AssertContains(t, matches[0], "agent one.txt")
}

func TestPath_PathMatch_Good(t *T) {
	r := PathMatch("agent.?", "agent.a")
	AssertTrue(t, r.OK)
	AssertTrue(t, r.Value.(bool))
}

func TestPath_PathMatch_Bad(t *T) {
	// A malformed pattern surfaces as OK=false with the error in Value.
	r := PathMatch("[", "agent")
	AssertFalse(t, r.OK)
	AssertNotNil(t, r.Value)
	_, isBool := r.Value.(bool)
	AssertFalse(t, isBool)
	// A well-formed pattern that simply does not match is NOT an error:
	// OK stays true and the boolean result is false.
	r2 := PathMatch("*.go", "main.py")
	AssertTrue(t, r2.OK)
	AssertFalse(t, r2.Value.(bool))
}

func TestPath_PathMatch_Ugly(t *T) {
	r := PathMatch("agent.*", "agent.")
	AssertTrue(t, r.OK)
	AssertTrue(t, r.Value.(bool))
}

func TestPath_PathRel_Good(t *T) {
	r := PathRel("/srv/dappcore", "/srv/dappcore/agent/config.json")
	AssertTrue(t, r.OK)
	AssertEqual(t, "agent/config.json", r.Value.(string))
}

func TestPath_PathRel_Bad(t *T) {
	// Absolute base with a relative target cannot be related.
	r := PathRel("/srv/dappcore", "agent/config.json")
	AssertFalse(t, r.OK)
	AssertNotNil(t, r.Value)
	// Relative base with an absolute target is equally rejected.
	r2 := PathRel("srv/dappcore", "/agent/config.json")
	AssertFalse(t, r2.OK)
}

func TestPath_PathRel_Ugly(t *T) {
	r := PathRel("/srv/dappcore", "/srv/dappcore")
	AssertTrue(t, r.OK)
	AssertEqual(t, ".", r.Value.(string))
}

func TestPath_PathAbs_Good(t *T) {
	r := PathAbs(".")
	AssertTrue(t, r.OK)
	AssertTrue(t, PathIsAbs(r.Value.(string)))
}

func TestPath_PathAbs_Bad(t *T) {
	r := PathAbs("agent/workspace")
	AssertTrue(t, r.OK, "PathAbs is infallible for ordinary relative input")
	AssertTrue(t, PathIsAbs(r.Value.(string)))
}

func TestPath_PathAbs_Ugly(t *T) {
	r := PathAbs("")
	AssertTrue(t, r.OK)
	AssertTrue(t, PathIsAbs(r.Value.(string)))
}

func TestPath_PathEvalSymlinks_Good(t *T) {
	dir := t.TempDir()
	r := PathEvalSymlinks(dir)
	AssertTrue(t, r.OK)
	AssertTrue(t, PathIsAbs(r.Value.(string)))
	AssertEqual(t, PathBase(dir), PathBase(r.Value.(string)))
}

func TestPath_PathEvalSymlinks_Bad(t *T) {
	// A path that does not exist cannot be resolved.
	r := PathEvalSymlinks(Path(t.TempDir(), "missing"))
	AssertFalse(t, r.OK)
	AssertNotNil(t, r.Value)
	// A non-directory parent in the path is also unresolvable.
	dir := t.TempDir()
	file := Path(dir, "agent.txt")
	(&Fs{}).New("/").Write(file, "x")
	r2 := PathEvalSymlinks(Path(file, "child"))
	AssertFalse(t, r2.OK)
}

func TestPath_PathEvalSymlinks_Ugly(t *T) {
	dir := t.TempDir()
	file := Path(dir, "agent.txt")
	(&Fs{}).New("/").Write(file, "ready")
	r := PathEvalSymlinks(file)
	AssertTrue(t, r.OK)
	AssertTrue(t, PathIsAbs(r.Value.(string)))
	AssertEqual(t, "agent.txt", PathBase(r.Value.(string)))
}

func TestPath_PathToSlash_Good(t *T) {
	slashed := PathToSlash(PathJoin("deploy", "to", "homelab"))
	AssertEqual(t, "deploy/to/homelab", slashed)
	// The result is always forward-slash separated, never the Windows form.
	AssertContains(t, slashed, "/")
	AssertFalse(t, Contains(slashed, `\`))
}

func TestPath_PathToSlash_Bad(t *T) {
	// Nothing to convert: a bare element and the empty string pass through.
	AssertEqual(t, "agent", PathToSlash("agent"))
	AssertEqual(t, "", PathToSlash(""))
	// The current-directory marker is likewise untouched.
	AssertEqual(t, ".", PathToSlash("."))
}

func TestPath_PathToSlash_Ugly(t *T) {
	// Already-slashed input is returned verbatim and is idempotent.
	AssertEqual(t, "agent/dispatch", PathToSlash("agent/dispatch"))
	AssertEqual(t, "agent/dispatch", PathToSlash(PathToSlash("agent/dispatch")))
	// An absolute slash path is preserved leading slash and all.
	AssertEqual(t, "/a/b/c", PathToSlash("/a/b/c"))
}

func TestPath_PathWalk_Good(t *T) {
	dir := t.TempDir()
	f := (&Fs{}).New("/")
	f.Write(Path(dir, "agent.txt"), "ready")
	visited := 0
	r := PathWalk(dir, func(path string, _ FsFileInfo, err error) error {
		AssertNoError(t, err)
		if PathBase(path) == "agent.txt" {
			visited++
		}
		return nil
	})
	AssertTrue(t, r.OK)
	AssertEqual(t, 1, visited)
}

func TestPath_PathWalk_Bad(t *T) {
	r := PathWalk(Path(t.TempDir(), "missing"), func(path string, _ FsFileInfo, err error) error {
		return err
	})
	AssertFalse(t, r.OK)
}

func TestPath_PathWalk_Ugly(t *T) {
	dir := t.TempDir()
	f := (&Fs{}).New("/")
	f.EnsureDir(Path(dir, "skip"))
	f.Write(Path(dir, "skip/hidden.txt"), "hidden")
	visitedHidden := false
	r := PathWalk(dir, func(path string, info FsFileInfo, err error) error {
		AssertNoError(t, err)
		if info.IsDir() && PathBase(path) == "skip" {
			return PathSkipDir
		}
		if PathBase(path) == "hidden.txt" {
			visitedHidden = true
		}
		return nil
	})
	AssertTrue(t, r.OK)
	AssertFalse(t, visitedHidden)
}

func TestPath_PathWalkDir_Good(t *T) {
	dir := t.TempDir()
	(&Fs{}).New("/").Write(Path(dir, "agent.txt"), "ready")
	visited := 0
	r := PathWalkDir(dir, func(path string, d FsDirEntry, err error) error {
		AssertNoError(t, err)
		if !d.IsDir() && PathBase(path) == "agent.txt" {
			visited++
		}
		return nil
	})
	AssertTrue(t, r.OK)
	AssertEqual(t, 1, visited)
}

func TestPath_PathWalkDir_Bad(t *T) {
	r := PathWalkDir(Path(t.TempDir(), "missing"), func(path string, d FsDirEntry, err error) error {
		return err
	})
	AssertFalse(t, r.OK)
}

func TestPath_PathWalkDir_Ugly(t *T) {
	dir := t.TempDir()
	f := (&Fs{}).New("/")
	f.Write(Path(dir, "first.txt"), "first")
	f.Write(Path(dir, "second.txt"), "second")
	visited := 0
	r := PathWalkDir(dir, func(path string, d FsDirEntry, err error) error {
		AssertNoError(t, err)
		visited++
		if !d.IsDir() {
			return PathSkipAll
		}
		return nil
	})
	AssertTrue(t, r.OK)
	AssertLess(t, visited, 4)
}

func TestPath_PathChangeExt_Good(t *T) {
	// newExt is accepted with or without a leading dot.
	AssertEqual(t, "agent.yaml", PathChangeExt("agent.json", "yaml"))
	AssertEqual(t, "agent.yaml", PathChangeExt("agent.json", ".yaml"))
	// Only the final extension is replaced; earlier dots are preserved.
	AssertEqual(t, "archive.tar.zip", PathChangeExt("archive.tar.gz", "zip"))
}

func TestPath_PathChangeExt_Bad(t *T) {
	// With no existing extension, newExt is simply appended.
	AssertEqual(t, "README.md", PathChangeExt("README", ".md"))
	AssertEqual(t, "README.md", PathChangeExt("README", "md"))
	// A dotfile has no detected extension, so newExt is appended whole.
	AssertEqual(t, ".env.txt", PathChangeExt(".env", ".txt"))
}

func TestPath_PathChangeExt_Ugly(t *T) {
	// An empty newExt strips the existing extension entirely.
	AssertEqual(t, "agent", PathChangeExt("agent.json", ""))
	// Empty newExt on a path with no extension is a no-op.
	AssertEqual(t, "README", PathChangeExt("README", ""))
	// Only the FINAL extension is stripped; earlier dots survive.
	AssertEqual(t, "archive.tar", PathChangeExt("archive.tar.gz", ""))
}
