// SPDX-License-Identifier: EUPL-1.2

package core

func TestContract_discoverServiceName_Good(t *T) {
	// The vanity path and the GitHub subdir-module path give the same name, so
	// moving a module path doesn't rename its service.
	AssertEqual(t, "ansible", discoverServiceName("dappco.re/go/ansible"))
	AssertEqual(t, "ansible", discoverServiceName("github.com/dappcore/go-ansible/go"))
	AssertEqual(t, "config", discoverServiceName("github.com/dappcore/config/go"))
}

func TestContract_discoverServiceName_Bad(t *T) {
	// No package path (an unnamed type): no name, so WithService reports the error.
	AssertEqual(t, "", discoverServiceName(""))
}

func TestContract_discoverServiceName_Ugly(t *T) {
	// Subpackages keep their own name; a bare or rootless "go" stays "go".
	AssertEqual(t, "marketplace", discoverServiceName("github.com/dappcore/go-scm/go/marketplace"))
	AssertEqual(t, "go", discoverServiceName("go"))
	AssertEqual(t, "go", discoverServiceName("/go"))
	AssertEqual(t, "lns", discoverServiceName("github.com/dappcore/Go-LNS/go"))
}
