// SPDX-License-Identifier: EUPL-1.2

//go:build windows

// The console for a Windows GUI exe. An app built -H windowsgui (no console
// window behind its UI) has no stdout when someone runs one of its commands
// from cmd or PowerShell, so the output goes nowhere. AttachParentConsole
// borrows the console of the terminal that started it. The unix side is a
// no-op in console_other.go.

package core

import (
	"log"
	"syscall"
)

// AttachParentConsole gives a Windows GUI exe the console of the terminal
// that started it, so a command's output shows. Call it first in main, before
// anything prints, and only when the app was asked to run a command.
//
//	if len(os.Args) > 1 {
//		core.AttachParentConsole()
//	}
//
// Output that's already going somewhere (redirected to a file, piped, or run
// over SSH) is left alone. A double-clicked exe has no parent console and
// stays as it was. Result.Value is true when a console was attached.
func AttachParentConsole() Result {
	if hasStdout() {
		return Result{Value: false, OK: true}
	}
	const attachParentProcess = ^uintptr(0) // (DWORD)-1
	if r, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("AttachConsole").Call(attachParentProcess); r == 0 {
		return Result{Value: false, OK: true}
	}
	r := OpenFile("CONOUT$", O_RDWR, 0)
	if !r.OK {
		return Result{E("core.AttachParentConsole", "open the console", r.Err()), false}
	}
	w, ok := useForStdio(r.Value)
	if !ok {
		return Result{E("core.AttachParentConsole", "the console isn't a file", nil), false}
	}
	log.SetOutput(w)
	return Result{Value: true, OK: true}
}

// hasStdout reports whether stdout is already a real file, pipe or console.
func hasStdout() bool {
	h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil || h == syscall.InvalidHandle || h == 0 {
		return false
	}
	t, err := syscall.GetFileType(h)
	return err == nil && t != syscall.FILE_TYPE_UNKNOWN
}
