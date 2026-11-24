// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "github.com/dappcore/go"

// ExampleAttachParentConsole runs first in main for a Windows GUI exe, so
// `myapp.exe search frieren` typed into PowerShell prints its results. Output
// that's already redirected is kept, so here nothing is attached.
func ExampleAttachParentConsole() {
	r := AttachParentConsole()
	Println(r.OK, r.Value)
	// Output: true false
}
