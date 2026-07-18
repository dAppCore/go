// SPDX-License-Identifier: EUPL-1.2

// Built-in introspection actions — the capability map made queryable.
// Every Core answers these three (entitlement-gated at Run like any
// other action), so an agent can walk a mesh of Core apps discovering
// what each one can do:
//
//	r := c.API("codex").Discover()                       // remote peer
//	r = c.Action("core.actions").Run(ctx, core.NewOptions()) // local
//
// With the colon law, discovery composes: register a Drive handle for
// a peer and c.Action("peer:core.actions") resolves transparently.
package core

// registerBuiltins wires the core.* introspection actions during
// construction — before any conclave seal, so they ride every bundle.
func registerBuiltins(c *Core) {
	// core.actions — the capability map: every registered action name,
	// plus declared Schema keys where an action published them.
	c.Action("core.actions", func(_ Context, _ Options) Result {
		names := c.Actions()
		schemas := make(map[string][]string)
		for _, name := range names {
			a := c.Action(name)
			if a.Schema.Len() > 0 {
				keys := make([]string, 0, a.Schema.Len())
				for _, opt := range a.Schema.Items() {
					keys = append(keys, opt.Key)
				}
				schemas[name] = keys
			}
		}
		return Ok(map[string]any{"actions": names, "schemas": schemas})
	})

	// core.info — application identity and enabled features.
	c.Action("core.info", func(_ Context, _ Options) Result {
		return Ok(map[string]any{
			"name":     c.app.Name,
			"version":  c.app.Version,
			"os":       OS(),
			"arch":     Arch(),
			"features": c.config.featureFlags().Names(),
		})
	})

	// core.health — liveness probe with process uptime.
	c.Action("core.health", func(_ Context, _ Options) Result {
		return Ok(map[string]any{
			"status": "ok",
			"uptime": Since(c.bootTime).String(),
		})
	})
}
