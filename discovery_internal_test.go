// SPDX-License-Identifier: EUPL-1.2

package core

func TestDiscovery_registerBuiltins_Good(t *T) {
	// The three builtins answer: capability map, identity, liveness.
	c := New(WithOption("name", "homelab"))

	actions := c.Action("core.actions").Run(c.Context(), NewOptions())
	info := c.Action("core.info").Run(c.Context(), NewOptions())
	health := c.Action("core.health").Run(c.Context(), NewOptions())

	AssertTrue(t, actions.OK)
	AssertContains(t, actions.Value.(map[string]any)["actions"], "core.health")
	AssertEqual(t, "homelab", info.Value.(map[string]any)["name"])
	AssertEqual(t, "ok", health.Value.(map[string]any)["status"])
}

func TestDiscovery_registerBuiltins_Bad(t *T) {
	// An unnamed app still reports identity: empty name, real OS and arch.
	c := newCore()

	info := c.Action("core.info").Run(c.Context(), NewOptions()).Value.(map[string]any)

	AssertEqual(t, "", info["name"])
	AssertEqual(t, OS(), info["os"])
	AssertEqual(t, Arch(), info["arch"])
}

func TestDiscovery_registerBuiltins_Ugly(t *T) {
	// core.actions is live: it sees actions and schemas registered after the builtins.
	c := newCore()
	a := c.Action("agent.run", func(Context, Options) Result { return Ok(nil) })
	a.Schema = NewOptions(Option{Key: "repo"})

	out := c.Action("core.actions").Run(c.Context(), NewOptions()).Value.(map[string]any)

	AssertContains(t, out["actions"], "agent.run")
	AssertEqual(t, []string{"repo"}, out["schemas"].(map[string][]string)["agent.run"])
}
