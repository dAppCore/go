// SPDX-License-Identifier: EUPL-1.2

//go:build !windows

package core

// AttachParentConsole is a no-op outside Windows: a terminal's output is
// already there. See console_windows.go.
//
//	core.AttachParentConsole() // Result{Value: false, OK: true}
func AttachParentConsole() Result {
	return Result{Value: false, OK: true}
}
