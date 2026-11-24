// SPDX-License-Identifier: EUPL-1.2

//go:build windows

package core

import "syscall"

// setStdout swaps the process's stdout handle for the test and puts it back.
func setStdout(t *T, h syscall.Handle) {
	t.Helper()
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("SetStdHandle")
	prev, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	RequireNoError(t, err)
	which := syscall.STD_OUTPUT_HANDLE // -11: a variable, so it wraps to the DWORD
	stdOutput := uintptr(which)
	proc.Call(stdOutput, uintptr(h))
	t.Cleanup(func() { proc.Call(stdOutput, uintptr(prev)) })
}

// Good — under go test stdout is a pipe or file: there is output.
func TestConsole_hasStdout_Good(t *T) {
	AssertTrue(t, hasStdout())
}

// Bad — no stdout handle at all, as in a double-clicked GUI exe.
func TestConsole_hasStdout_Bad(t *T) {
	setStdout(t, 0)

	AssertFalse(t, hasStdout())
}

// Ugly — stdout is the NUL device (a character device, not a file or pipe)
// and still counts as output that's going somewhere.
func TestConsole_hasStdout_Ugly(t *T) {
	nul, err := syscall.Open("NUL", syscall.O_WRONLY, 0)
	RequireNoError(t, err)
	defer syscall.CloseHandle(nul)
	setStdout(t, nul)

	AssertTrue(t, hasStdout())
}
