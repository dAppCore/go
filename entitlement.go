// SPDX-License-Identifier: EUPL-1.2

// Permission primitive for the Core framework.
// Entitlement answers "can [subject] do [action] with [quantity]?"
// Default: everything permitted (trusted conclave).
// With go-entitlements: checks workspace packages, features, usage, boosts.
// With commerce-matrix: checks entity hierarchy, lock cascade.
//
// Usage:
//
//	e := c.Entitled("process.run")           // boolean gate
//	e := c.Entitled("social.accounts", 3)    // quantity check
//	if e.Allowed { proceed() }
//	if e.NearLimit(0.8) { showUpgradePrompt() }
//
// Registration:
//
//	c.SetEntitlementChecker(myChecker)
//	c.SetUsageRecorder(myRecorder)
package core

// Entitlement is the result of a permission check.
// Carries context for both boolean gates (Allowed) and usage limits (Limit/Used/Remaining).
//
//	e := c.Entitled("social.accounts", 3)
//	e.Allowed     // true
//	e.Limit       // 5
//	e.Used        // 2
//	e.Remaining   // 3
//	e.NearLimit(0.8) // false
type Entitlement struct {
	Allowed   bool   // permission granted
	Unlimited bool   // no cap (agency tier, admin, trusted conclave)
	Limit     int    // total allowed (0 = boolean gate)
	Used      int    // current consumption
	Remaining int    // Limit - Used
	Reason    string // denial reason — for UI and audit logging
}

// NearLimit returns true if usage exceeds the threshold percentage.
//
//	if e.NearLimit(0.8) { showUpgradePrompt() }
func (e Entitlement) NearLimit(threshold float64) bool {
	if e.Unlimited || e.Limit == 0 {
		return false
	}
	return float64(e.Used)/float64(e.Limit) >= threshold
}

// UsagePercent returns current usage as a percentage of the limit.
//
//	pct := e.UsagePercent() // 75.0
func (e Entitlement) UsagePercent() float64 {
	if e.Limit == 0 {
		return 0
	}
	return float64(e.Used) / float64(e.Limit) * 100
}

// EntitlementChecker answers "can [subject] do [action] with [quantity]?"
// Subject comes from context (workspace, entity, user — consumer's concern).
//
//	checker := func(action string, quantity int, ctx Context) core.Entitlement {
//	    return core.Entitlement{Allowed: action == "process.run", Limit: 10, Used: quantity}
//	}
//	core.New().SetEntitlementChecker(core.EntitlementChecker(checker))
type EntitlementChecker func(action string, quantity int, ctx Context) Entitlement

// UsageRecorder records consumption after a gated action succeeds.
// Consumer packages provide the implementation (database, cache, etc).
//
//	recorder := func(action string, quantity int, ctx Context) {
//	    core.Info("usage recorded", "action", action, "quantity", quantity)
//	}
//	core.New().SetUsageRecorder(core.UsageRecorder(recorder))
type UsageRecorder func(action string, quantity int, ctx Context)

// defaultChecker — trusted conclave, everything permitted.
func defaultChecker(_ string, _ int, _ Context) Entitlement {
	return Entitlement{Allowed: true, Unlimited: true}
}

// Entitled checks if an action is permitted in the current context.
// Default: always returns Allowed=true, Unlimited=true.
// Denials are logged via core.Security().
//
//	e := c.Entitled("process.run")
//	e := c.Entitled("social.accounts", 3)
func (c *Core) Entitled(action string, quantity ...int) Entitlement {
	qty := 1
	if len(quantity) > 0 {
		qty = quantity[0]
	}

	e := c.entitlementChecker(action, qty, c.Context())

	if !e.Allowed {
		Security("entitlement.denied", "action", action, "quantity", qty, "reason", e.Reason)
	}

	return e
}

// SetEntitlementChecker replaces the default (permissive) checker.
// Called by go-entitlements or commerce-matrix during OnStartup.
//
//	func (s *EntitlementService) OnStartup(ctx Context) core.Result {
//	    s.Core().SetEntitlementChecker(s.check)
//	    return core.Result{OK: true}
//	}
func (c *Core) SetEntitlementChecker(checker EntitlementChecker) {
	c.entitlementChecker = checker
}

// RecordUsage records consumption after a gated action succeeds.
// Delegates to the registered UsageRecorder. No-op if none registered.
//
//	e := c.Entitled("ai.credits", 10)
//	if e.Allowed {
//	    doWork()
//	    c.RecordUsage("ai.credits", 10)
//	}
func (c *Core) RecordUsage(action string, quantity ...int) {
	if c.usageRecorder == nil {
		return
	}
	qty := 1
	if len(quantity) > 0 {
		qty = quantity[0]
	}
	c.usageRecorder(action, qty, c.Context())
}

// SetUsageRecorder registers a usage tracking function.
// Called by go-entitlements during OnStartup.
//
//	c := core.New()
//	c.SetUsageRecorder(func(action string, quantity int, ctx Context) {
//	    core.Info("usage recorded", "action", action, "quantity", quantity)
//	})
func (c *Core) SetUsageRecorder(recorder UsageRecorder) {
	c.usageRecorder = recorder
}

// --- Policy: declarative entitlements + quotas ---

// Policy is a declarative entitlement rulebook: exact action names or
// "prefix.*" globs mapping to "allow", "deny", or an integer quota.
// Checker gates at Action.Run; Recorder decrements consumption through
// the metering loop — wire both and quotas enforce themselves. Unruled
// actions default to allowed (rules restrict, absence permits).
//
//	p := core.NewPolicy(core.NewOptions(
//	    core.Option{Key: "agentic.*", Value: "allow"},
//	    core.Option{Key: "admin.purge", Value: "deny"},
//	    core.Option{Key: "ai.credits", Value: 100},
//	))
//	c.SetEntitlementChecker(p.Checker())
//	c.SetUsageRecorder(p.Recorder())
type Policy struct {
	rules Options
	mu    Mutex
	used  map[string]int // consumption per quota rule key
}

// NewPolicy builds a Policy from a rule set. Rule values: "allow",
// "deny", or an int quota. Exact keys beat globs; longer globs beat
// shorter ones.
//
//	p := core.NewPolicy(core.NewOptions(core.Option{Key: "admin.*", Value: "deny"}))
func NewPolicy(rules Options) *Policy {
	return &Policy{rules: rules, used: map[string]int{}}
}

// match finds the most specific rule for an action: exact match first,
// then the longest matching "prefix.*" glob.
func (p *Policy) match(action string) (string, any, bool) {
	if r := p.rules.Get(action); r.OK {
		return action, r.Value, true
	}
	bestKey := ""
	var bestVal any
	for _, opt := range p.rules.Items() {
		if !HasSuffix(opt.Key, ".*") {
			continue
		}
		if HasPrefix(action, TrimSuffix(opt.Key, "*")) && len(opt.Key) > len(bestKey) {
			bestKey, bestVal = opt.Key, opt.Value
		}
	}
	if bestKey == "" {
		return "", nil, false
	}
	return bestKey, bestVal, true
}

// Checker returns the EntitlementChecker enforcing this policy.
//
//	c.SetEntitlementChecker(p.Checker())
func (p *Policy) Checker() EntitlementChecker {
	return func(action string, quantity int, _ Context) Entitlement {
		key, val, ok := p.match(action)
		if !ok {
			return Entitlement{Allowed: true, Unlimited: true}
		}
		switch v := val.(type) {
		case string:
			if v == "deny" {
				return Entitlement{Allowed: false, Reason: Concat("policy: ", key, " denied")}
			}
			return Entitlement{Allowed: true, Unlimited: true}
		case int:
			if quantity <= 0 {
				quantity = 1
			}
			p.mu.Lock()
			used := p.used[key]
			p.mu.Unlock()
			e := Entitlement{
				Allowed:   v-used >= quantity,
				Limit:     v,
				Used:      used,
				Remaining: v - used,
			}
			if !e.Allowed {
				e.Reason = Concat("policy: quota exhausted for ", key)
			}
			return e
		}
		return Entitlement{Allowed: true, Unlimited: true}
	}
}

// Recorder returns the UsageRecorder that decrements quotas — pair it
// with Checker via SetUsageRecorder so successful gated actions consume.
//
//	c.SetUsageRecorder(p.Recorder())
func (p *Policy) Recorder() UsageRecorder {
	return func(action string, quantity int, _ Context) {
		key, val, ok := p.match(action)
		if !ok {
			return
		}
		if _, isQuota := val.(int); !isQuota {
			return
		}
		if quantity <= 0 {
			quantity = 1
		}
		p.mu.Lock()
		p.used[key] += quantity
		p.mu.Unlock()
	}
}
