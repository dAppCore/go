// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func TestOS_FileMode_Good_Alias(t *T) {
	var mode FileMode = 0o644
	AssertEqual(t, FileMode(0o644), mode)
	AssertEqual(t, FileMode(0o644), mode.Perm())
	AssertFalse(t, mode.IsDir())
	AssertTrue(t, mode.IsRegular())
	AssertEqual(t, "-rw-r--r--", mode.String())
}

func TestOS_ModePerm_Good(t *T) {
	AssertEqual(t, FileMode(0o777), ModePerm)
	AssertEqual(t, FileMode(0), ModeDir&ModePerm)
	AssertEqual(t, FileMode(0o644), (ModeDir|0o644).Perm())
	AssertTrue(t, ModeType&ModeDir != 0)
	AssertTrue(t, ModeType&ModeSymlink != 0)
}

func TestOS_ModeDir_Good(t *T) {
	mode := ModeDir | 0o755
	AssertTrue(t, mode.IsDir())
	AssertEqual(t, FileMode(0o755), mode.Perm())
}

func TestOS_Stdin_Good_NotNil(t *T) {
	in := Stdin()
	AssertNotNil(t, in)
	f, ok := in.(*OSFile)
	AssertTrue(t, ok)
	AssertEqual(t, "/dev/stdin", f.Name())
}

func TestOS_Stdout_Good_NotNil(t *T) {
	out := Stdout()
	AssertNotNil(t, out)
	f, ok := out.(*OSFile)
	AssertTrue(t, ok)
	AssertEqual(t, "/dev/stdout", f.Name())
}

func TestOS_Stderr_Good_NotNil(t *T) {
	errStream := Stderr()
	AssertNotNil(t, errStream)
	f, ok := errStream.(*OSFile)
	AssertTrue(t, ok)
	AssertNotEmpty(t, f.Name())
}

func TestOS_Args_Good(t *T) {
	args := Args()

	AssertNotEmpty(t, args)
	AssertNotEmpty(t, args[0])
}

func TestOS_Args_Bad(t *T) {
	args := Args()

	AssertNotNil(t, args)
	AssertGreaterOrEqual(t, len(args), 1)
	AssertNotEmpty(t, args[0])
}

func TestOS_Args_Ugly(t *T) {
	first := Args()
	second := Args()

	AssertEqual(t, first[0], second[0])
}

func TestOS_Chdir_Good(t *T) {
	cwd := Getwd()
	RequireTrue(t, cwd.OK)
	defer func() { AssertTrue(t, Chdir(cwd.Value.(string)).OK) }()
	dir := t.TempDir()
	realDir := PathEvalSymlinks(dir)
	RequireTrue(t, realDir.OK)

	r := Chdir(dir)

	AssertTrue(t, r.OK)
	after := Getwd()
	AssertTrue(t, after.OK)
	AssertEqual(t, realDir.Value.(string), after.Value.(string))
}

func TestOS_Chdir_Bad(t *T) {
	before := Getwd()
	RequireTrue(t, before.OK)

	r := Chdir(Path(t.TempDir(), "missing"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	after := Getwd()
	AssertTrue(t, after.OK)
	AssertEqual(t, before.Value.(string), after.Value.(string))
}

func TestOS_Chdir_Ugly(t *T) {
	cwd := Getwd()
	RequireTrue(t, cwd.OK)

	r := Chdir(".")

	AssertTrue(t, r.OK)
	after := Getwd()
	AssertTrue(t, after.OK)
	AssertEqual(t, cwd.Value.(string), after.Value.(string))
}

func TestOS_Create_Bad(t *T) {
	path := Path(t.TempDir(), "missing", "agent.log")

	r := Create(path)

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	AssertFalse(t, Stat(path).OK)
}

func TestOS_Create_Ugly(t *T) {
	path := Path(t.TempDir(), "agent.log")
	r := Create(path)
	RequireTrue(t, r.OK)
	CloseStream(r.Value)

	second := Create(path)
	RequireTrue(t, second.OK)
	CloseStream(second.Value)
	read := ReadFile(path)
	AssertTrue(t, read.OK)
	AssertEqual(t, []byte{}, read.Value.([]byte))
}

func TestOS_DirFS_Good(t *T) {
	dir := t.TempDir()
	path := Path(dir, "agent.txt")
	AssertTrue(t, WriteFile(path, []byte("ready"), 0o644).OK)

	r := ReadFSFile(DirFS(dir), "agent.txt")

	AssertTrue(t, r.OK)
	AssertEqual(t, []byte("ready"), r.Value.([]byte))
}

func TestOS_DirFS_Bad(t *T) {
	fsys := DirFS(t.TempDir())

	r := ReadFSFile(fsys, "missing.txt")

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	escape := ReadFSFile(fsys, "../escape.txt")
	AssertFalse(t, escape.OK)
}

func TestOS_DirFS_Ugly(t *T) {
	r := ReadDir(DirFS(t.TempDir()), ".")

	AssertTrue(t, r.OK)
	AssertLen(t, r.Value.([]FsDirEntry), 0)
}

func TestOS_Environ_Good(t *T) {
	t.Setenv("CORE_AX7_AGENT", "dispatch")

	env := Environ()

	found := false
	for _, entry := range env {
		if entry == "CORE_AX7_AGENT=dispatch" {
			found = true
		}
	}
	AssertTrue(t, found)
}

func TestOS_Environ_Bad(t *T) {
	env := Environ()

	AssertNotNil(t, env)
	AssertNotEmpty(t, env)
	for _, entry := range env {
		AssertContains(t, entry, "=")
	}
}

func TestOS_Environ_Ugly(t *T) {
	t.Setenv("CORE_AX7_EMPTY", "")

	env := Environ()

	found := false
	for _, entry := range env {
		if entry == "CORE_AX7_EMPTY=" {
			found = true
		}
	}
	AssertTrue(t, found)
}

func TestOS_Getenv_Good(t *T) {
	t.Setenv("CORE_AX7_TOKEN", "session-token")
	AssertEqual(t, "session-token", Getenv("CORE_AX7_TOKEN"))

	t.Setenv("CORE_AX7_TOKEN", "rotated")
	AssertEqual(t, "rotated", Getenv("CORE_AX7_TOKEN"))
	value, ok := LookupEnv("CORE_AX7_TOKEN")
	AssertTrue(t, ok)
	AssertEqual(t, "rotated", value)
}

func TestOS_Getenv_Bad(t *T) {
	RequireTrue(t, Unsetenv("CORE_AX7_MISSING").OK)

	AssertEqual(t, "", Getenv("CORE_AX7_MISSING"))
	_, ok := LookupEnv("CORE_AX7_MISSING")
	AssertFalse(t, ok)
}

func TestOS_Getenv_Ugly(t *T) {
	t.Setenv("CORE_AX7_EMPTY", "")

	AssertEqual(t, "", Getenv("CORE_AX7_EMPTY"))
	value, ok := LookupEnv("CORE_AX7_EMPTY")
	AssertTrue(t, ok)
	AssertEqual(t, "", value)
}

func TestOS_Getpid_Good(t *T) {
	pid := Getpid()

	AssertGreater(t, pid, 0)
	AssertEqual(t, pid, Getpid())
}

func TestOS_Getpid_Bad(t *T) {
	pid := Getpid()

	AssertNotEqual(t, 0, pid)
	AssertGreater(t, pid, 0)
	AssertNotEqual(t, Getppid(), pid)
}

func TestOS_Getpid_Ugly(t *T) {
	first := Getpid()
	second := Getpid()

	AssertEqual(t, first, second)
	AssertGreater(t, first, 0)
}

func TestOS_Getppid_Good(t *T) {
	ppid := Getppid()

	AssertGreater(t, ppid, 0)
	AssertEqual(t, ppid, Getppid())
}

func TestOS_Getppid_Bad(t *T) {
	ppid := Getppid()

	AssertNotEqual(t, 0, ppid)
	AssertGreater(t, ppid, 0)
	AssertNotEqual(t, Getpid(), ppid)
}

func TestOS_Getppid_Ugly(t *T) {
	first := Getppid()
	second := Getppid()

	AssertEqual(t, first, second)
	AssertGreater(t, first, 0)
}

func TestOS_Getwd_Good(t *T) {
	r := Getwd()

	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_Getwd_Bad(t *T) {
	r := Getwd()

	AssertTrue(t, r.OK)
	wd := r.Value.(string)
	AssertNotEmpty(t, wd)
	AssertTrue(t, PathIsAbs(wd))
}

func TestOS_Getwd_Ugly(t *T) {
	cwd := Getwd()
	RequireTrue(t, cwd.OK)
	defer func() { AssertTrue(t, Chdir(cwd.Value.(string)).OK) }()
	dir := t.TempDir()
	realDir := PathEvalSymlinks(dir)
	RequireTrue(t, realDir.OK)

	AssertTrue(t, Chdir(dir).OK)
	r := Getwd()

	AssertTrue(t, r.OK)
	AssertEqual(t, realDir.Value.(string), r.Value.(string))
}

func TestOS_Hostname_Good(t *T) {
	r := Hostname()

	AssertTrue(t, r.OK)
	name := r.Value.(string)
	AssertNotEmpty(t, name)
	AssertNotContains(t, name, " ")
}

func TestOS_Hostname_Bad(t *T) {
	r := Hostname()

	AssertTrue(t, r.OK)
	name, ok := r.Value.(string)
	AssertTrue(t, ok)
	AssertNotEmpty(t, name)
}

func TestOS_Hostname_Ugly(t *T) {
	first := Hostname()
	second := Hostname()

	AssertTrue(t, first.OK)
	AssertTrue(t, second.OK)
	AssertEqual(t, first.Value.(string), second.Value.(string))
}

func TestOS_IsExist_Good(t *T) {
	r := Mkdir(t.TempDir(), 0o755)

	AssertFalse(t, r.OK)
	AssertTrue(t, IsExist(r.Value.(error)))
}

func TestOS_IsExist_Bad(t *T) {
	AssertFalse(t, IsExist(AnError))
	AssertFalse(t, IsExist(ErrNotExist))
	AssertFalse(t, IsExist(ErrPermission))
}

func TestOS_IsExist_Ugly(t *T) {
	AssertFalse(t, IsExist(nil))
	AssertTrue(t, IsExist(ErrExist))
	// The predicate only unwraps OS error types, not fmt %w wrapping —
	// unlike core.Is, which walks the whole chain.
	AssertFalse(t, IsExist(Errorf("create: %w", ErrExist)))
	AssertTrue(t, Is(Errorf("create: %w", ErrExist), ErrExist))
}

func TestOS_IsNotExist_Good(t *T) {
	r := ReadFile(Path(t.TempDir(), "missing.txt"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
}

func TestOS_IsNotExist_Bad(t *T) {
	AssertFalse(t, IsNotExist(AnError))
	AssertFalse(t, IsNotExist(ErrExist))
	AssertFalse(t, IsNotExist(ErrPermission))
}

func TestOS_IsNotExist_Ugly(t *T) {
	AssertFalse(t, IsNotExist(nil))
	AssertTrue(t, IsNotExist(ErrNotExist))
	// fmt %w wrapping defeats the predicate but not core.Is.
	AssertFalse(t, IsNotExist(Errorf("open: %w", ErrNotExist)))
	AssertTrue(t, Is(Errorf("open: %w", ErrNotExist), ErrNotExist))
}

func TestOS_IsPermission_Good(t *T) {
	AssertTrue(t, IsPermission(ErrPermissionForTest))
	AssertTrue(t, IsPermission(ErrPermission))
	AssertFalse(t, IsPermission(ErrNotExist))
}

func TestOS_IsPermission_Bad(t *T) {
	AssertFalse(t, IsPermission(AnError))
	AssertFalse(t, IsPermission(ErrExist))
	AssertFalse(t, IsPermission(ErrNotExist))
}

func TestOS_IsPermission_Ugly(t *T) {
	AssertFalse(t, IsPermission(nil))
	AssertTrue(t, IsPermission(ErrPermission))
	// fmt %w wrapping defeats the predicate but not core.Is.
	AssertFalse(t, IsPermission(Errorf("chmod: %w", ErrPermission)))
	AssertTrue(t, Is(Errorf("chmod: %w", ErrPermission), ErrPermission))
}

func TestOS_LookupEnv_Good(t *T) {
	t.Setenv("CORE_AX7_LOOKUP", "present")

	value, ok := LookupEnv("CORE_AX7_LOOKUP")

	AssertTrue(t, ok)
	AssertEqual(t, "present", value)
}

func TestOS_LookupEnv_Bad(t *T) {
	Unsetenv("CORE_AX7_LOOKUP_MISSING")

	value, ok := LookupEnv("CORE_AX7_LOOKUP_MISSING")

	AssertFalse(t, ok)
	AssertEqual(t, "", value)
}

func TestOS_LookupEnv_Ugly(t *T) {
	t.Setenv("CORE_AX7_LOOKUP_EMPTY", "")

	value, ok := LookupEnv("CORE_AX7_LOOKUP_EMPTY")

	AssertTrue(t, ok)
	AssertEqual(t, "", value)
}

func TestOS_Lstat_Good(t *T) {
	path := Path(t.TempDir(), "agent.txt")
	AssertTrue(t, WriteFile(path, []byte("ready"), 0o644).OK)

	r := Lstat(path)

	AssertTrue(t, r.OK)
	info := r.Value.(interface{ Name() string })
	AssertEqual(t, "agent.txt", info.Name())
}

func TestOS_Lstat_Bad(t *T) {
	r := Lstat(Path(t.TempDir(), "missing.txt"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	empty := Lstat("")
	AssertFalse(t, empty.OK)
}

func TestOS_Lstat_Ugly(t *T) {
	dir := t.TempDir()
	target := Path(dir, "agent.txt")
	link := Path(dir, "current")
	AssertTrue(t, WriteFile(target, []byte("ready"), 0o644).OK)
	RequireNoError(t, SymlinkForTest(target, link))

	r := Lstat(link)

	AssertTrue(t, r.OK)
	info := r.Value.(interface{ Mode() FileMode })
	AssertTrue(t, info.Mode()&ModeSymlink != 0)
}

func TestOS_Mkdir_Good(t *T) {
	path := Path(t.TempDir(), "agent")

	r := Mkdir(path, 0o755)

	AssertTrue(t, r.OK)
	AssertTrue(t, Stat(path).OK)
}

func TestOS_Mkdir_Bad(t *T) {
	dir := t.TempDir()

	r := Mkdir(dir, 0o755)

	AssertFalse(t, r.OK)
}

func TestOS_Mkdir_Ugly(t *T) {
	r := Mkdir("", 0o755)

	AssertFalse(t, r.OK)
	AssertNotNil(t, r.Value)
	AssertTrue(t, IsNotExist(r.Value.(error)))
}

func TestOS_MkdirAll_Good(t *T) {
	path := Path(t.TempDir(), "agent", "dispatch", "logs")

	r := MkdirAll(path, 0o755)

	AssertTrue(t, r.OK)
	AssertTrue(t, Stat(path).OK)
}

func TestOS_MkdirAll_Bad(t *T) {
	dir := t.TempDir()
	blocker := Path(dir, "agent")
	AssertTrue(t, WriteFile(blocker, []byte("file"), 0o644).OK)

	r := MkdirAll(Path(blocker, "dispatch"), 0o755)

	AssertFalse(t, r.OK)
}

func TestOS_MkdirAll_Ugly(t *T) {
	r := MkdirAll("", 0o755)

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
}

func TestOS_MkdirTemp_Good(t *T) {
	r := MkdirTemp("", "agent-*")
	RequireTrue(t, r.OK)
	defer RemoveAll(r.Value.(string))

	AssertTrue(t, Stat(r.Value.(string)).OK)
}

func TestOS_MkdirTemp_Bad(t *T) {
	r := MkdirTemp(Path(t.TempDir(), "missing"), "agent-*")

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
}

func TestOS_MkdirTemp_Ugly(t *T) {
	r := MkdirTemp("", "")
	RequireTrue(t, r.OK)
	defer RemoveAll(r.Value.(string))

	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_Open_Ugly(t *T) {
	path := Path(t.TempDir(), "empty.txt")
	AssertTrue(t, WriteFile(path, nil, 0o644).OK)

	r := Open(path)

	AssertTrue(t, r.OK)
	CloseStream(r.Value)
}

func TestOS_OpenFile_Good(t *T) {
	path := Path(t.TempDir(), "agent.log")

	r := OpenFile(path, O_CREATE|O_WRONLY, 0o644)
	RequireTrue(t, r.OK)
	defer CloseStream(r.Value)

	AssertTrue(t, WriteAll(r.Value, "ready").OK)
}

func TestOS_OpenFile_Bad(t *T) {
	r := OpenFile(Path(t.TempDir(), "missing.log"), O_RDONLY, 0o644)

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	wr := OpenFile(t.TempDir(), O_WRONLY, 0o644)
	AssertFalse(t, wr.OK)
}

func TestOS_OpenFile_Ugly(t *T) {
	path := Path(t.TempDir(), "agent.log")
	AssertTrue(t, WriteFile(path, []byte("ready"), 0o644).OK)

	r := OpenFile(path, O_CREATE|O_EXCL|O_WRONLY, 0o644)

	AssertFalse(t, r.OK)
}

func TestOS_ReadFile_Bad(t *T) {
	r := ReadFile(Path(t.TempDir(), "missing.txt"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	dir := ReadFile(t.TempDir())
	AssertFalse(t, dir.OK)
}

func TestOS_ReadFile_Ugly(t *T) {
	path := Path(t.TempDir(), "empty.txt")
	AssertTrue(t, WriteFile(path, nil, 0o644).OK)

	r := ReadFile(path)

	AssertTrue(t, r.OK)
	AssertEqual(t, []byte{}, r.Value.([]byte))
}

func TestOS_Remove_Ugly(t *T) {
	r := Remove(Path(t.TempDir(), "missing.txt"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	dir := t.TempDir()
	RequireTrue(t, WriteFile(Path(dir, "child"), nil, 0o644).OK)
	nonEmpty := Remove(dir)
	AssertFalse(t, nonEmpty.OK)
}

func TestOS_RemoveAll_Good(t *T) {
	dir := t.TempDir()
	path := Path(dir, "agent", "dispatch.log")
	AssertTrue(t, MkdirAll(PathDir(path), 0o755).OK)
	AssertTrue(t, WriteFile(path, []byte("ready"), 0o644).OK)

	r := RemoveAll(Path(dir, "agent"))

	AssertTrue(t, r.OK)
	AssertFalse(t, Stat(Path(dir, "agent")).OK)
}

func TestOS_RemoveAll_Bad(t *T) {
	missing := Path(t.TempDir(), "missing")

	r := RemoveAll(missing)

	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
	AssertTrue(t, RemoveAll(missing).OK)
	AssertFalse(t, Stat(missing).OK)
}

func TestOS_RemoveAll_Ugly(t *T) {
	r := RemoveAll("")

	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
	AssertTrue(t, RemoveAll(Path(t.TempDir(), "a", "b", "c")).OK)
}

func TestOS_Rename_Bad(t *T) {
	dir := t.TempDir()

	r := Rename(Path(dir, "missing.txt"), Path(dir, "agent.txt"))

	AssertFalse(t, r.OK)
}

func TestOS_Rename_Ugly(t *T) {
	dir := t.TempDir()
	oldPath := Path(dir, "agent.tmp")
	newPath := Path(dir, "agent.json")
	AssertTrue(t, WriteFile(oldPath, []byte("new"), 0o644).OK)
	AssertTrue(t, WriteFile(newPath, []byte("old"), 0o644).OK)

	r := Rename(oldPath, newPath)

	AssertTrue(t, r.OK)
	read := ReadFile(newPath)
	AssertTrue(t, read.OK)
	AssertEqual(t, []byte("new"), read.Value.([]byte))
}

func TestOS_Stat_Bad(t *T) {
	r := Stat(Path(t.TempDir(), "missing.txt"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	empty := Stat("")
	AssertFalse(t, empty.OK)
}

func TestOS_Stat_Ugly(t *T) {
	dir := t.TempDir()

	r := Stat(dir)

	AssertTrue(t, r.OK)
	info := r.Value.(interface{ IsDir() bool })
	AssertTrue(t, info.IsDir())
}

func TestOS_Stdin_Good(t *T) {
	in := Stdin()

	AssertNotNil(t, in)
	AssertEqual(t, uintptr(0), in.(*OSFile).Fd())
}

func TestOS_Stdin_Bad(t *T) {
	var reader Reader = Stdin()

	AssertNotNil(t, reader)
	_, isFile := reader.(*OSFile)
	AssertTrue(t, isFile)
}

func TestOS_Stdin_Ugly(t *T) {
	first := Stdin()
	second := Stdin()

	AssertEqual(t, first, second)
}

func TestOS_Stdout_Good(t *T) {
	out := Stdout()

	AssertNotNil(t, out)
	AssertEqual(t, uintptr(1), out.(*OSFile).Fd())
}

func TestOS_Stdout_Bad(t *T) {
	r := WriteString(Stdout(), "")

	AssertTrue(t, r.OK)
	AssertEqual(t, 0, r.Value.(int))
}

func TestOS_Stdout_Ugly(t *T) {
	first := Stdout()
	second := Stdout()

	AssertEqual(t, first, second)
}

func TestOS_Stderr_Good(t *T) {
	errStream := Stderr()

	AssertNotNil(t, errStream)
	AssertEqual(t, errStream, Stderr())
	AssertTrue(t, WriteString(errStream, "").OK)
}

func TestOS_Stderr_Bad(t *T) {
	r := WriteString(Stderr(), "")

	AssertTrue(t, r.OK)
	AssertEqual(t, 0, r.Value.(int))
}

func TestOS_Stderr_Ugly(t *T) {
	first := Stderr()
	second := Stderr()

	AssertEqual(t, first, second)
}

func TestOS_TempDir_Good(t *T) {
	dir := TempDir()

	AssertNotEmpty(t, dir)
	AssertTrue(t, Stat(dir).OK)
}

func TestOS_TempDir_Bad(t *T) {
	dir := TempDir()

	AssertNotEmpty(t, dir)
	AssertTrue(t, PathIsAbs(dir))
	AssertEqual(t, dir, TempDir())
}

func TestOS_TempDir_Ugly(t *T) {
	custom := t.TempDir()
	t.Setenv("TMPDIR", custom)

	got := TempDir()

	AssertNotEmpty(t, got)
	AssertEqual(t, custom, got)
}

func TestOS_Unsetenv_Good(t *T) {
	t.Setenv("CORE_AX7_UNSET", "value")

	r := Unsetenv("CORE_AX7_UNSET")
	_, ok := LookupEnv("CORE_AX7_UNSET")

	AssertTrue(t, r.OK)
	AssertFalse(t, ok)
}

func TestOS_Unsetenv_Bad(t *T) {
	r := Unsetenv("")

	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
	AssertTrue(t, Unsetenv("").OK)
}

func TestOS_Unsetenv_Ugly(t *T) {
	r := Unsetenv("CORE_AX7_ALREADY_MISSING")

	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
	_, ok := LookupEnv("CORE_AX7_ALREADY_MISSING")
	AssertFalse(t, ok)
}

func TestOS_UserCacheDir_Good(t *T) {
	r := UserCacheDir()

	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_UserCacheDir_Bad(t *T) {
	r := UserCacheDir()

	AssertTrue(t, r.OK)
	dir := r.Value.(string)
	AssertNotEmpty(t, dir)
	AssertTrue(t, PathIsAbs(dir))
}

func TestOS_UserCacheDir_Ugly(t *T) {
	t.Setenv("HOME", t.TempDir())

	r := UserCacheDir()

	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_UserConfigDir_Good(t *T) {
	r := UserConfigDir()

	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_UserConfigDir_Bad(t *T) {
	r := UserConfigDir()

	AssertTrue(t, r.OK)
	dir := r.Value.(string)
	AssertNotEmpty(t, dir)
	AssertTrue(t, PathIsAbs(dir))
}

func TestOS_UserConfigDir_Ugly(t *T) {
	t.Setenv("HOME", t.TempDir())

	r := UserConfigDir()

	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_UserHomeDir_Good(t *T) {
	r := UserHomeDir()

	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.(string))
}

func TestOS_UserHomeDir_Bad(t *T) {
	r := UserHomeDir()

	AssertTrue(t, r.OK)
	dir := r.Value.(string)
	AssertNotEmpty(t, dir)
	AssertTrue(t, PathIsAbs(dir))
}

func TestOS_UserHomeDir_Ugly(t *T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	r := UserHomeDir()

	AssertTrue(t, r.OK)
	AssertEqual(t, home, r.Value.(string))
}

func TestOS_WriteFile_Good(t *T) {
	path := Path(t.TempDir(), "agent.json")

	r := WriteFile(path, []byte(`{"agent":"dispatch"}`), 0o644)

	AssertTrue(t, r.OK)
	read := ReadFile(path)
	AssertTrue(t, read.OK)
	AssertEqual(t, []byte(`{"agent":"dispatch"}`), read.Value.([]byte))
}

func TestOS_WriteFile_Bad(t *T) {
	dir := t.TempDir()

	r := WriteFile(dir, []byte("not a file"), 0o644)

	AssertFalse(t, r.OK)
	AssertNotNil(t, r.Value)
	missing := WriteFile(Path(dir, "missing", "f.txt"), []byte("x"), 0o644)
	AssertFalse(t, missing.OK)
	AssertTrue(t, IsNotExist(missing.Value.(error)))
}

func TestOS_WriteFile_Ugly(t *T) {
	path := Path(t.TempDir(), "empty.txt")

	r := WriteFile(path, nil, 0o600)

	AssertTrue(t, r.OK)
	read := ReadFile(path)
	AssertTrue(t, read.OK)
	AssertEqual(t, []byte{}, read.Value.([]byte))
}

func TestOS_Chmod_Good(t *T) {
	path := Path(t.TempDir(), "bin")
	RequireTrue(t, WriteFile(path, []byte("#!/bin/sh\n"), 0o644).OK)

	r := Chmod(path, 0o755)

	AssertTrue(t, r.OK)
	info := Stat(path)
	RequireTrue(t, info.OK)
	AssertEqual(t, FileMode(0o755), info.Value.(FsFileInfo).Mode().Perm())
}

func TestOS_Chmod_Bad(t *T) {
	r := Chmod(Path(t.TempDir(), "missing"), 0o755)

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	empty := Chmod("", 0o755)
	AssertFalse(t, empty.OK)
}

func TestOS_Chmod_Ugly(t *T) {
	// Re-applying the same mode is a no-op that must still succeed.
	path := Path(t.TempDir(), "f")
	RequireTrue(t, WriteFile(path, nil, 0o600).OK)

	AssertTrue(t, Chmod(path, 0o600).OK)
	AssertTrue(t, Chmod(path, 0o600).OK)
}

func TestOS_Symlink_Good(t *T) {
	dir := t.TempDir()
	target := Path(dir, "target")
	RequireTrue(t, WriteFile(target, []byte("ready"), 0o644).OK)
	link := Path(dir, "link")

	r := Symlink(target, link)

	AssertTrue(t, r.OK)
	read := ReadFile(link)
	AssertTrue(t, read.OK)
	AssertEqual(t, []byte("ready"), read.Value.([]byte))
}

func TestOS_Symlink_Bad(t *T) {
	// A link path under a non-existent parent directory cannot be created.
	r := Symlink("target", Path(t.TempDir(), "missing", "link"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
}

func TestOS_Symlink_Ugly(t *T) {
	// Creating a link where a file already exists must fail, not clobber.
	dir := t.TempDir()
	link := Path(dir, "link")
	RequireTrue(t, WriteFile(link, nil, 0o644).OK)

	AssertFalse(t, Symlink("target", link).OK)
}

func TestOS_Readlink_Good(t *T) {
	dir := t.TempDir()
	target := Path(dir, "target")
	RequireTrue(t, WriteFile(target, nil, 0o644).OK)
	link := Path(dir, "link")
	RequireTrue(t, Symlink(target, link).OK)

	r := Readlink(link)

	AssertTrue(t, r.OK)
	AssertEqual(t, target, r.Value.(string))
}

func TestOS_Readlink_Bad(t *T) {
	r := Readlink(Path(t.TempDir(), "missing"))

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
}

func TestOS_Readlink_Ugly(t *T) {
	// A regular file is not a symlink — Readlink must report failure.
	path := Path(t.TempDir(), "regular")
	RequireTrue(t, WriteFile(path, nil, 0o644).OK)

	AssertFalse(t, Readlink(path).OK)
}

func TestOS_CreateTemp_Good(t *T) {
	r := CreateTemp(t.TempDir(), "agent-*.json")

	AssertTrue(t, r.OK)
	f := r.Value.(*OSFile)
	defer CloseStream(f)
	AssertTrue(t, HasSuffix(f.Name(), ".json"))
}

func TestOS_CreateTemp_Bad(t *T) {
	r := CreateTemp(Path(t.TempDir(), "missing"), "agent-*")

	AssertFalse(t, r.OK)
	AssertTrue(t, IsNotExist(r.Value.(error)))
	sep := CreateTemp(t.TempDir(), "bad/pattern-*")
	AssertFalse(t, sep.OK)
}

func TestOS_CreateTemp_Ugly(t *T) {
	// Two calls with the same pattern must yield distinct files.
	dir := t.TempDir()
	first := CreateTemp(dir, "x-*")
	RequireTrue(t, first.OK)
	defer CloseStream(first.Value)
	second := CreateTemp(dir, "x-*")
	RequireTrue(t, second.OK)
	defer CloseStream(second.Value)

	AssertNotEqual(t, first.Value.(*OSFile).Name(), second.Value.(*OSFile).Name())
}

func TestOS_Executable_Good(t *T) {
	r := Executable()

	AssertTrue(t, r.OK)
	AssertTrue(t, PathIsAbs(r.Value.(string)))
}

func TestOS_Executable_Bad(t *T) {
	// Executable takes no input that could be made invalid; the contract
	// is that it resolves to a real on-disk binary on supported platforms.
	r := Executable()

	AssertTrue(t, r.OK)
	path := r.Value.(string)
	AssertNotEmpty(t, path)
	AssertTrue(t, Stat(path).OK)
}

func TestOS_Executable_Ugly(t *T) {
	// Repeated calls within a process return the same path.
	first := Executable()
	second := Executable()

	RequireTrue(t, first.OK)
	RequireTrue(t, second.OK)
	AssertEqual(t, first.Value.(string), second.Value.(string))
}

func TestOS_ErrNotExist_Good(t *T) {
	r := Open(Path(t.TempDir(), "missing"))

	RequireTrue(t, !r.OK)
	AssertTrue(t, Is(r.Value.(error), ErrNotExist))
}

func TestOS_ErrNotExist_Bad(t *T) {
	// A successful open carries no error to match against the sentinel, and an
	// unrelated error does not match ErrNotExist either.
	path := Path(t.TempDir(), "present")
	RequireTrue(t, WriteFile(path, nil, 0o644).OK)
	r := Open(path)

	RequireTrue(t, r.OK)
	CloseStream(r.Value)
	AssertFalse(t, Is(AnError, ErrNotExist))
}

func TestOS_ErrNotExist_Ugly(t *T) {
	// The other sentinels are distinct from ErrNotExist.
	AssertFalse(t, Is(ErrExist, ErrNotExist))
	AssertFalse(t, Is(ErrPermission, ErrNotExist))
	AssertFalse(t, Is(ErrInvalid, ErrNotExist))
	AssertFalse(t, Is(ErrClosed, ErrNotExist))
}
