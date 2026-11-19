// SPDX-License-Identifier: EUPL-1.2

package core

func TestEntitlement_defaultChecker_Good(t *T) {
	entitlement := defaultChecker("process.run", 1, Background())

	AssertTrue(t, entitlement.Allowed)
	AssertTrue(t, entitlement.Unlimited)
}
func TestEntitlement_defaultChecker_Bad(t *T) {
	entitlement := defaultChecker("process.delete", -1, Background())

	AssertTrue(t, entitlement.Allowed)
	AssertTrue(t, entitlement.Unlimited)
}
func TestEntitlement_defaultChecker_Ugly(t *T) {
	entitlement := defaultChecker("", 0, nil)

	AssertTrue(t, entitlement.Allowed)
	AssertTrue(t, entitlement.Unlimited)
}

func TestEntitlement_Policy_match_Good(t *T) {
	// An exact rule wins over any glob.
	p := NewPolicy(NewOptions(Option{Key: "social.*", Value: false}, Option{Key: "social.post", Value: true}))

	key, val, ok := p.match("social.post")

	AssertTrue(t, ok)
	AssertEqual(t, "social.post", key)
	AssertEqual(t, true, val)
}

func TestEntitlement_Policy_match_Bad(t *T) {
	// No rule covers the action: nothing matched.
	p := NewPolicy(NewOptions(Option{Key: "social.*", Value: true}))

	key, val, ok := p.match("agent.run")

	AssertFalse(t, ok)
	AssertEqual(t, "", key)
	AssertNil(t, val)
}

func TestEntitlement_Policy_match_Ugly(t *T) {
	// The longest matching glob wins; a glob needs its dot ("socialx" isn't "social.*").
	p := NewPolicy(NewOptions(
		Option{Key: "social.*", Value: 1},
		Option{Key: "social.accounts.*", Value: 2},
	))

	key, val, ok := p.match("social.accounts.add")
	AssertTrue(t, ok)
	AssertEqual(t, "social.accounts.*", key)
	AssertEqual(t, 2, val)

	_, _, ok = p.match("socialx")
	AssertFalse(t, ok)
}
