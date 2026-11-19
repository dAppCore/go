// SPDX-License-Identifier: EUPL-1.2

package core

func TestLock_Core_sealConclave_Good(t *T) {
	// With the service lock on, every capability registry freezes.
	c := New(WithServiceLock())

	c.sealConclave()

	AssertTrue(t, c.ipc.actions.Locked())
	AssertTrue(t, c.ipc.tasks.Locked())
	AssertTrue(t, c.commands.Registry.Locked())
	AssertTrue(t, c.api.protocols.Locked())
	AssertTrue(t, c.drive.Registry.Locked())
	AssertTrue(t, c.data.Registry.Locked())
	AssertTrue(t, c.config.featureFlags().Sealed())
}

func TestLock_Core_sealConclave_Bad(t *T) {
	// Without the service lock, nothing is frozen: the conclave stays open.
	c := New()

	c.sealConclave()

	AssertFalse(t, c.ipc.actions.Locked())
	AssertFalse(t, c.commands.Registry.Locked())
	AssertFalse(t, c.config.featureFlags().Sealed())
	AssertTrue(t, c.Command("late", Command{Action: func(Options) Result { return Ok(nil) }}).OK)
}

func TestLock_Core_sealConclave_Ugly(t *T) {
	// Feature flags are sealed, not locked: existing flags still toggle, new
	// ones can't appear; runtime Config settings stay writable.
	c := New(WithServiceLock())
	c.Config().Enable("beta")

	c.sealConclave()
	c.Config().Disable("beta")
	c.Config().Enable("gamma")
	c.Config().Set("runtime.key", "still writable")

	AssertFalse(t, c.Config().Enabled("beta"))
	AssertFalse(t, c.Config().Enabled("gamma"))
	AssertEqual(t, "still writable", c.Config().String("runtime.key"))
	// Registrations after the seal are refused out loud, not dropped while reporting OK.
	AssertFalse(t, c.Command("late", Command{Action: func(Options) Result { return Ok(nil) }}).OK)
	AssertFalse(t, c.Command("late").OK)
	AssertFalse(t, c.Drive().New(NewOptions(Option{Key: "name", Value: "late"}, Option{Key: "transport", Value: "http://x"})).OK)
	AssertFalse(t, c.Data().MountDir("late", t.TempDir()).OK)
}
