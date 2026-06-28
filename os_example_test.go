package core_test

import . "dappco.re/go"

// ExampleFileMode uses file mode aliases through `FileMode` for process IO access. File
// modes and standard streams are exposed through core aliases.
func ExampleFileMode() {
	var mode FileMode = 0644
	Println((mode & ModePerm) == 0644)
	// Output: true
}

// ExampleModeDir checks directory mode aliases through `ModeDir` for process IO access.
// File modes and standard streams are exposed through core aliases.
func ExampleModeDir() {
	Println(ModeDir.IsDir())
	// Output: true
}

// ExampleStdin reads standard input through `Stdin` for process IO access. File modes and
// standard streams are exposed through core aliases.
func ExampleStdin() {
	Println(Stdin() != nil)
	// Output: true
}

// ExampleStdout reads standard output through `Stdout` for process IO access. File modes
// and standard streams are exposed through core aliases.
func ExampleStdout() {
	Println(Stdout() != nil)
	// Output: true
}

// ExampleStderr reads standard error through `Stderr` for process IO access. File modes
// and standard streams are exposed through core aliases.
func ExampleStderr() {
	Println(Stderr() != nil)
	// Output: true
}

// ExampleChmod sets a file's permission bits at the OS boundary, then
// reads them back through Stat. Chmod is the unsandboxed sibling of the
// c.Fs() permission helpers.
func ExampleChmod() {
	path := PathJoin(TempDir(), "core-example-chmod")
	WriteFile(path, []byte("#!/bin/sh\n"), 0o644)
	defer Remove(path)

	Chmod(path, 0o755)
	info := Stat(path)
	Println(info.Value.(FsFileInfo).Mode().Perm() == 0o755)
	// Output: true
}

// ExampleErrNotExist matches a failed Open against the re-exported
// sentinel without importing os.
// ExampleTempDir returns the system temporary directory through `TempDir`.
func ExampleTempDir() {
	Println(TempDir() != "")
	// Output: true
}

func ExampleErrNotExist() {
	r := Open(PathJoin(TempDir(), "core-example-definitely-missing"))
	Println(Is(r.Value.(error), ErrNotExist))
	// Output: true
}

// --- file operations (demonstration calls) ---

// ExampleReadFile reads a file's bytes through `ReadFile`.
func ExampleReadFile() {
	r := ReadFile(PathJoin(TempDir(), "config.json"))
	_ = r // r.Value is []byte on success
}

// ExampleWriteFile writes bytes to a file through `WriteFile`.
func ExampleWriteFile() {
	_ = WriteFile(PathJoin(TempDir(), "out.txt"), []byte("data"), 0o644)
}

// ExampleMkdirAll creates a directory tree through `MkdirAll`.
func ExampleMkdirAll() {
	_ = MkdirAll(PathJoin(TempDir(), "a", "b"), 0o755)
}

// ExampleMkdir creates a single directory through `Mkdir`.
func ExampleMkdir() {
	_ = Mkdir(PathJoin(TempDir(), "newdir"), 0o755)
}

// ExampleCreate creates or truncates a file through `Create`.
func ExampleCreate() {
	_ = Create(PathJoin(TempDir(), "new.txt"))
}

// ExampleOpen opens a file for reading through `Open`.
func ExampleOpen() {
	_ = Open(PathJoin(TempDir(), "data.txt"))
}

// ExampleOpenFile opens a file with flags and mode through `OpenFile`.
func ExampleOpenFile() {
	_ = OpenFile(PathJoin(TempDir(), "log.txt"), 0, 0o644)
}

// ExampleStat returns file metadata through `Stat`.
func ExampleStat() {
	_ = Stat(PathJoin(TempDir(), "f"))
}

// ExampleLstat returns link metadata without following it through `Lstat`.
func ExampleLstat() {
	_ = Lstat(PathJoin(TempDir(), "f"))
}

// ExampleRemove deletes a file or empty directory through `Remove`.
func ExampleRemove() {
	_ = Remove(PathJoin(TempDir(), "f"))
}

// ExampleRemoveAll deletes a path and any children through `RemoveAll`.
func ExampleRemoveAll() {
	_ = RemoveAll(PathJoin(TempDir(), "dir"))
}

// ExampleRename moves a file through `Rename`.
func ExampleRename() {
	_ = Rename(PathJoin(TempDir(), "a"), PathJoin(TempDir(), "b"))
}

// ExampleSymlink creates a symbolic link through `Symlink`.
func ExampleSymlink() {
	_ = Symlink(PathJoin(TempDir(), "target"), PathJoin(TempDir(), "link"))
}

// ExampleReadlink resolves a symbolic link through `Readlink`.
func ExampleReadlink() {
	_ = Readlink(PathJoin(TempDir(), "link"))
}

// ExampleMkdirTemp creates a uniquely-named temp directory through `MkdirTemp`.
func ExampleMkdirTemp() {
	_ = MkdirTemp(TempDir(), "ex-*")
}

// ExampleCreateTemp creates a uniquely-named temp file through `CreateTemp`.
func ExampleCreateTemp() {
	_ = CreateTemp(TempDir(), "ex-*")
}

// ExampleChdir changes the working directory through `Chdir`.
func ExampleChdir() {
	_ = Chdir(TempDir())
}

// ExampleHostname returns the host name through `Hostname`.
func ExampleHostname() {
	_ = Hostname() // r.Value is the host name on success
}

// ExampleExecutable returns the running binary's path through `Executable`.
func ExampleExecutable() {
	_ = Executable()
}

// ExampleGetwd returns the working directory through `Getwd`.
func ExampleGetwd() {
	_ = Getwd()
}

// ExampleUserHomeDir returns the user's home directory through `UserHomeDir`.
func ExampleUserHomeDir() {
	_ = UserHomeDir()
}

// ExampleUserConfigDir returns the user's config directory through `UserConfigDir`.
func ExampleUserConfigDir() {
	_ = UserConfigDir()
}

// ExampleUserCacheDir returns the user's cache directory through `UserCacheDir`.
func ExampleUserCacheDir() {
	_ = UserCacheDir()
}

// --- process & environment ---

// ExampleArgs returns the command-line arguments through `Args`.
func ExampleArgs() {
	Println(len(Args()) > 0)
	// Output: true
}

// ExampleGetpid returns the current process ID through `Getpid`.
func ExampleGetpid() {
	Println(Getpid() > 0)
	// Output: true
}

// ExampleGetppid returns the parent process ID through `Getppid`.
func ExampleGetppid() {
	Println(Getppid() > 0)
	// Output: true
}

// ExampleEnviron returns the environment as KEY=VALUE strings through `Environ`.
func ExampleEnviron() {
	Println(len(Environ()) > 0)
	// Output: true
}

// ExampleGetenv reads an environment variable through `Getenv`.
func ExampleGetenv() {
	Setenv("CORE_EXAMPLE", "homelab")
	Println(Getenv("CORE_EXAMPLE"))
	Unsetenv("CORE_EXAMPLE")
	// Output: homelab
}

// ExampleLookupEnv reports whether an environment variable is set through `LookupEnv`.
func ExampleLookupEnv() {
	Setenv("CORE_EXAMPLE", "x")
	_, ok := LookupEnv("CORE_EXAMPLE")
	Println(ok)
	Unsetenv("CORE_EXAMPLE")
	// Output: true
}

// ExampleIsNotExist reports whether an error means "not found" through `IsNotExist`.
func ExampleIsNotExist() {
	Println(IsNotExist(ErrNotExist))
	// Output: true
}

// ExampleIsExist reports whether an error means "already exists" through `IsExist`.
func ExampleIsExist() {
	r := Mkdir(TempDir(), 0o755) // TempDir already exists
	if !r.OK {
		_ = IsExist(r.Value.(error))
	}
}

// ExampleIsPermission reports whether an error means "permission denied" through `IsPermission`.
func ExampleIsPermission() {
	r := Open("/root/secret-core-example")
	if !r.OK {
		_ = IsPermission(r.Value.(error))
	}
}
