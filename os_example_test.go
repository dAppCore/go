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
func ExampleErrNotExist() {
	r := Open(PathJoin(TempDir(), "core-example-definitely-missing"))
	Println(Is(r.Value.(error), ErrNotExist))
	// Output: true
}
