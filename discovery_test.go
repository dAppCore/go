// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func TestDiscovery_CoreActions_Good(t *T) {
	c := New()
	a := c.Action("agent.deploy", func(Context, Options) Result { return Ok(nil) })
	a.Schema = NewOptions(Option{Key: "target", Value: "required"})
	r := c.Action("core.actions").Run(Background(), NewOptions())
	AssertTrue(t, r.OK)
	m := r.Value.(map[string]any)
	AssertContains(t, m["actions"].([]string), "agent.deploy")
	AssertEqual(t, []string{"target"}, m["schemas"].(map[string][]string)["agent.deploy"])
}

func TestDiscovery_CoreInfo_Good(t *T) {
	c := New(WithOption("name", "meshnode"))
	r := c.Action("core.info").Run(Background(), NewOptions())
	AssertTrue(t, r.OK)
	AssertEqual(t, "meshnode", r.Value.(map[string]any)["name"].(string))
}

func TestDiscovery_CoreHealth_Good(t *T) {
	r := New().Action("core.health").Run(Background(), NewOptions())
	AssertTrue(t, r.OK)
	AssertEqual(t, "ok", r.Value.(map[string]any)["status"].(string))
}

func TestDiscovery_CoreActions_Bad_EntitlementGated(t *T) {
	// The built-ins pass the same gate as every action.
	c := New()
	c.SetEntitlementChecker(func(action string, _ int, _ Context) Entitlement {
		return Entitlement{Allowed: action != "core.actions", Reason: "introspection denied"}
	})
	AssertFalse(t, c.Action("core.actions").Run(Background(), NewOptions()).OK)
	AssertTrue(t, c.Action("core.health").Run(Background(), NewOptions()).OK)
}

func TestDiscovery_CoreActions_Ugly_RideEveryBundle(t *T) {
	// Sealed conclave: the built-ins were registered pre-seal and survive.
	c := New(WithServiceLock())
	AssertTrue(t, c.ServiceStartup(Background(), nil).OK)
	AssertTrue(t, c.Action("core.actions").Run(Background(), NewOptions()).OK)
}
