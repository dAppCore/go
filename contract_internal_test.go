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

func TestContract_newCore_Good(t *T) {
	// Every subsystem is wired, with the builtins registered and nothing else.
	c := newCore()

	AssertNotNil(t, c.app)
	AssertNotNil(t, c.fs)
	AssertNotNil(t, c.config)
	AssertNotNil(t, c.ipc)
	AssertNotNil(t, c.api)
	AssertTrue(t, c.api.core == c)
	AssertTrue(t, c.Action("core.health").Exists())
	AssertEqual(t, 0, len(c.Services()))
}

func TestContract_newCore_Bad(t *T) {
	// No options applied and no lock: registration is open and the context live.
	c := newCore()

	AssertNil(t, c.options)
	AssertFalse(t, c.services.Locked())
	AssertNil(t, c.Context().Err())
}

func TestContract_newCore_Ugly(t *T) {
	// Two cores share nothing: separate registries and a boot time each.
	a, b := newCore(), newCore()
	a.Action("only.a", func(Context, Options) Result { return Ok(nil) })

	AssertFalse(t, b.Action("only.a").Exists())
	AssertFalse(t, a.ipc == b.ipc)
	AssertFalse(t, a.bootTime.IsZero())
}
