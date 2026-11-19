// SPDX-License-Identifier: EUPL-1.2

package core

func TestConfig_ConfigOptions_init_Good(t *T) {
	var opts ConfigOptions

	opts.init()

	AssertNotNil(t, opts.Settings)
	AssertNotNil(t, opts.Features)
}
func TestConfig_ConfigOptions_init_Bad(t *T) {
	opts := ConfigOptions{Settings: map[string]any{"agent": "codex"}}

	opts.init()

	AssertEqual(t, "codex", opts.Settings["agent"])
	AssertNotNil(t, opts.Features)
}
func TestConfig_ConfigOptions_init_Ugly(t *T) {
	opts := ConfigOptions{}

	opts.init()
	opts.Settings["site"] = "homelab"
	opts.init()

	AssertEqual(t, "homelab", opts.Settings["site"])
}

func TestConfig_Config_featureFlags_Good(t *T) {
	// A zero Config seeds the Registry once from the declarative input.
	cfg := &Config{ConfigOptions: &ConfigOptions{Features: map[string]bool{"beta": true}}}

	ff := cfg.featureFlags()

	AssertTrue(t, ff.Get("beta").OK)
	AssertTrue(t, cfg.Enabled("beta"))
}

func TestConfig_Config_featureFlags_Bad(t *T) {
	// No ConfigOptions at all: an empty Registry, no panic.
	cfg := &Config{}

	AssertEqual(t, 0, cfg.featureFlags().Len())
}

func TestConfig_Config_featureFlags_Ugly(t *T) {
	// Later edits to the declarative map don't leak in: seeding happens once.
	opts := &ConfigOptions{Features: map[string]bool{}}
	cfg := &Config{ConfigOptions: opts}
	first := cfg.featureFlags()
	opts.Features["late"] = true

	AssertTrue(t, cfg.featureFlags() == first)
	AssertFalse(t, cfg.Enabled("late"))
}

func TestConfig_Config_lock_Good(t *T) {
	cfg := (&Config{}).New()

	AssertTrue(t, cfg.lock() == &cfg.mu)
}

func TestConfig_Config_lock_Bad(t *T) {
	// A Group view never guards with its own mutex: it shares the root's.
	cfg := (&Config{}).New()
	db := cfg.Group("database")

	AssertTrue(t, db.lock() == &cfg.mu)
	AssertFalse(t, db.lock() == &db.mu)
}

func TestConfig_Config_lock_Ugly(t *T) {
	// Nested groups still resolve to the one root lock.
	cfg := (&Config{}).New()

	AssertTrue(t, cfg.Group("a").Group("b").lock() == &cfg.mu)
}

func TestConfig_Config_key_Good(t *T) {
	cfg := (&Config{}).New()

	AssertEqual(t, "database.host", cfg.Group("database").key("host"))
}

func TestConfig_Config_key_Bad(t *T) {
	// No prefix: the key passes through untouched.
	cfg := (&Config{}).New()

	AssertEqual(t, "host", cfg.key("host"))
}

func TestConfig_Config_key_Ugly(t *T) {
	// Nested prefixes compose; an empty key still gets the prefix.
	cfg := (&Config{}).New()

	AssertEqual(t, "a.b.host", cfg.Group("a").Group("b").key("host"))
	AssertEqual(t, "a.", cfg.Group("a").key(""))
}

func TestConfig_Config_mergeMap_Good(t *T) {
	cfg := (&Config{}).New()

	cfg.mergeMap("", map[string]any{"name": "agent", "port": 8080})

	AssertEqual(t, "agent", cfg.String("name"))
	AssertEqual(t, 8080, cfg.Int("port"))
}

func TestConfig_Config_mergeMap_Bad(t *T) {
	// Nil and empty maps are no-ops.
	cfg := (&Config{}).New()

	cfg.mergeMap("", nil)
	cfg.mergeMap("db", map[string]any{})

	AssertFalse(t, cfg.Get("db").OK)
}

func TestConfig_Config_mergeMap_Ugly(t *T) {
	// Nested maps flatten to dotted keys under the prefix.
	cfg := (&Config{}).New()

	cfg.mergeMap("app", map[string]any{"database": map[string]any{"host": "homelab", "pool": map[string]any{"size": 4}}})

	AssertEqual(t, "homelab", cfg.String("app.database.host"))
	AssertEqual(t, 4, cfg.Int("app.database.pool.size"))
}
