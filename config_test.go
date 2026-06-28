package core_test

import (
	. "dappco.re/go"
)

// --- Config ---

func TestConfig_Get_Bad(t *T) {
	c := New()
	r := c.Config().Get("missing")
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

// --- Feature Flags ---

func TestConfig_Features_Disable_Good(t *T) {
	c := New()
	c.Config().Enable("feature")
	AssertTrue(t, c.Config().Enabled("feature"))

	c.Config().Disable("feature")
	AssertFalse(t, c.Config().Enabled("feature"))
}

func TestConfig_Features_CaseSensitive(t *T) {
	c := New()
	c.Config().Enable("Feature")
	AssertTrue(t, c.Config().Enabled("Feature"))
	AssertFalse(t, c.Config().Enabled("feature"))
}

func TestConfig_EnabledFeatures_Good(t *T) {
	c := New()
	c.Config().Enable("a")
	c.Config().Enable("b")
	c.Config().Enable("c")
	c.Config().Disable("b")

	features := c.Config().EnabledFeatures()
	AssertContains(t, features, "a")
	AssertContains(t, features, "c")
	AssertNotContains(t, features, "b")
}

// --- Group (G1) ---

func TestConfig_Group_Good(t *T) {
	c := New()
	c.Config().Group("database").Set("host", "localhost")

	// Stored under the prefixed key, readable via root, .Group(), or shortcut.
	AssertEqual(t, "localhost", c.Config().String("database.host"))
	AssertEqual(t, "localhost", c.Config().Group("database").String("host"))
	AssertEqual(t, "localhost", c.Config("database").String("host"))
}

func TestConfig_Group_Bad(t *T) {
	c := New()
	// An empty group name is the root — no prefix applied.
	c.Config().Group("").Set("host", "localhost")
	AssertEqual(t, "localhost", c.Config().String("host"))
}

func TestConfig_Group_Ugly(t *T) {
	c := New()
	// Nested groups compose into a dotted prefix.
	c.Config("a").Group("b").Set("k", "v")
	AssertEqual(t, "v", c.Config().String("a.b.k"))
	AssertEqual(t, "v", c.Config("a").String("b.k"))
}

// --- Feature (G1) ---

func TestConfig_Feature_Name_Good(t *T) {
	AssertEqual(t, "dark-mode", New().Feature("dark-mode").Name())
}

func TestConfig_Feature_Name_Bad(t *T) {
	AssertEqual(t, "", New().Feature("").Name())
}

func TestConfig_Feature_Name_Ugly(t *T) {
	// The handle reports its raw name verbatim, dots and all, set or not.
	AssertEqual(t, "ns.sub.flag", New().Feature("ns.sub.flag").Name())
}

func TestConfig_Feature_Enable_Good(t *T) {
	c := New()
	c.Feature("dark-mode").Enable()
	AssertTrue(t, c.Feature("dark-mode").Enabled())
	AssertTrue(t, c.Config().Enabled("dark-mode")) // same backing store
}

func TestConfig_Feature_Enable_Bad(t *T) {
	c := New()
	// Enabling one feature leaves an unrelated one off.
	c.Feature("a").Enable()
	AssertFalse(t, c.Feature("b").Enabled())
}

func TestConfig_Feature_Enable_Ugly(t *T) {
	c := New()
	// Enable is idempotent — re-enabling keeps it on, listed once.
	c.Feature("beta").Enable()
	c.Feature("beta").Enable()
	AssertTrue(t, c.Feature("beta").Enabled())
	AssertContains(t, c.Config().EnabledFeatures(), "beta")
}

func TestConfig_Feature_Disable_Good(t *T) {
	c := New()
	c.Feature("beta").Enable()
	c.Feature("beta").Disable()
	AssertFalse(t, c.Feature("beta").Enabled())
}

func TestConfig_Feature_Disable_Bad(t *T) {
	c := New()
	// Disabling a never-set feature keeps it off and absent from the set.
	c.Feature("never-set").Disable()
	AssertFalse(t, c.Feature("never-set").Enabled())
	AssertNotContains(t, c.Config().EnabledFeatures(), "never-set")
}

func TestConfig_Feature_Disable_Ugly(t *T) {
	c := New()
	// Double-disable is safe and stays off.
	c.Feature("beta").Enable()
	c.Feature("beta").Disable()
	c.Feature("beta").Disable()
	AssertFalse(t, c.Feature("beta").Enabled())
}

func TestConfig_Feature_Enabled_Good(t *T) {
	c := New()
	c.Feature("dark-mode").Enable()
	AssertTrue(t, c.Feature("dark-mode").Enabled())
}

func TestConfig_Feature_Enabled_Bad(t *T) {
	// An unset feature reads as disabled.
	AssertFalse(t, New().Feature("never-set").Enabled())
}

func TestConfig_Feature_Enabled_Ugly(t *T) {
	c := New()
	// Feature handles are NOT auto-namespaced: a grouped flag is invisible
	// to the ungrouped handle of the same leaf name.
	c.Config("ui").Enable("dark")
	AssertTrue(t, c.Config().Enabled("ui.dark"))
	AssertFalse(t, c.Feature("dark").Enabled())
}

// --- AX-7 canonical triplets ---

func TestConfig_Config_New_Good(t *T) {
	cfg := (&Config{}).New()
	AssertNotNil(t, cfg)
	AssertFalse(t, cfg.Get("agent.host").OK)
	AssertEmpty(t, cfg.EnabledFeatures())
}

func TestConfig_Config_New_Bad(t *T) {
	cfg := (&Config{}).New()
	AssertEqual(t, "", cfg.String("missing"))
	AssertEqual(t, 0, cfg.Int("missing"))
	AssertFalse(t, cfg.Bool("missing"))
}

func TestConfig_Config_New_Ugly(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.host", "homelab.lthn.sh")
	cfg = cfg.New()
	AssertFalse(t, cfg.Get("agent.host").OK)
}

func TestConfig_Config_Set_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.host", "homelab.lthn.sh")
	r := cfg.Get("agent.host")
	AssertTrue(t, r.OK)
	AssertEqual(t, "homelab.lthn.sh", r.Value)
}

func TestConfig_Config_Set_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("", "root-setting")
	r := cfg.Get("")
	AssertTrue(t, r.OK)
	AssertEqual(t, "root-setting", r.Value)
}

func TestConfig_Config_Set_Ugly(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.optional", nil)
	r := cfg.Get("agent.optional")
	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
}

func TestConfig_Config_Get_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.mode", "review")
	r := cfg.Get("agent.mode")
	AssertTrue(t, r.OK)
	AssertEqual(t, "review", r.Value)
}

func TestConfig_Config_Get_Bad(t *T) {
	cfg := (&Config{}).New()
	r := cfg.Get("agent.mode")
	AssertFalse(t, r.OK)
}

func TestConfig_Config_Get_Ugly(t *T) {
	var cfg Config
	r := cfg.Get("agent.mode")
	AssertFalse(t, r.OK)
}

func TestConfig_Config_String_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.host", "homelab.lthn.sh")
	AssertEqual(t, "homelab.lthn.sh", cfg.String("agent.host"))
}

func TestConfig_Config_String_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.port", 9101)
	AssertEqual(t, "", cfg.String("agent.port"))
}

func TestConfig_Config_String_Ugly(t *T) {
	var cfg Config
	// Zero-value Config: missing key yields the empty string, no panic.
	AssertEqual(t, "", cfg.String("agent.host"))
	// Set lazily initialises the store, so a later String read sees the value.
	cfg.Set("agent.host", "homelab.lthn.sh")
	AssertEqual(t, "homelab.lthn.sh", cfg.String("agent.host"))
	// A non-string value reads back as the empty string (type mismatch).
	cfg.Set("agent.port", 9101)
	AssertEqual(t, "", cfg.String("agent.port"))
}

func TestConfig_Config_Int_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.port", 9101)
	AssertEqual(t, 9101, cfg.Int("agent.port"))
}

func TestConfig_Config_Int_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.port", "9101")
	AssertEqual(t, 0, cfg.Int("agent.port"))
}

func TestConfig_Config_Int_Ugly(t *T) {
	var cfg Config
	// Zero-value Config: missing key yields 0, no panic.
	AssertEqual(t, 0, cfg.Int("agent.port"))
	// Set lazily initialises the store, so a later Int read sees the value.
	cfg.Set("agent.port", 9101)
	AssertEqual(t, 9101, cfg.Int("agent.port"))
	// A non-int value reads back as 0 (type mismatch).
	cfg.Set("agent.host", "homelab.lthn.sh")
	AssertEqual(t, 0, cfg.Int("agent.host"))
}

func TestConfig_Config_Bool_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.enabled", true)
	AssertTrue(t, cfg.Bool("agent.enabled"))
}

func TestConfig_Config_Bool_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("agent.enabled", "true")
	AssertFalse(t, cfg.Bool("agent.enabled"))
}

func TestConfig_Config_Bool_Ugly(t *T) {
	var cfg Config
	// Zero-value Config: missing key yields false, no panic.
	AssertFalse(t, cfg.Bool("agent.enabled"))
	// Set lazily initialises the store, so a later Bool read sees the value.
	cfg.Set("agent.enabled", true)
	AssertTrue(t, cfg.Bool("agent.enabled"))
	// A non-bool value reads back as false (type mismatch).
	cfg.Set("agent.mode", "true")
	AssertFalse(t, cfg.Bool("agent.mode"))
}

func TestConfig_ConfigGet_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("session.token", "lethean-session")
	AssertEqual(t, "lethean-session", ConfigGet[string](cfg, "session.token"))
}

func TestConfig_ConfigGet_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Set("session.ttl", "60")
	AssertEqual(t, 0, ConfigGet[int](cfg, "session.ttl"))
}

func TestConfig_ConfigGet_Ugly(t *T) {
	AssertPanics(t, func() {
		ConfigGet[string](nil, "session.token")
	})
}

func TestConfig_Config_Enable_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Enable("agent.dispatch")
	AssertTrue(t, cfg.Enabled("agent.dispatch"))
}

func TestConfig_Config_Enable_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Enable("")
	AssertTrue(t, cfg.Enabled(""))
}

func TestConfig_Config_Enable_Ugly(t *T) {
	var cfg Config
	cfg.Enable("agent.dispatch")
	AssertTrue(t, cfg.Enabled("agent.dispatch"))
}

func TestConfig_Config_Disable_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Enable("agent.dispatch")
	cfg.Disable("agent.dispatch")
	AssertFalse(t, cfg.Enabled("agent.dispatch"))
}

func TestConfig_Config_Disable_Bad(t *T) {
	cfg := (&Config{}).New()
	cfg.Disable("missing.feature")
	AssertFalse(t, cfg.Enabled("missing.feature"))
}

func TestConfig_Config_Disable_Ugly(t *T) {
	var cfg Config
	cfg.Disable("")
	AssertFalse(t, cfg.Enabled(""))
}

func TestConfig_Config_Enabled_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Enable("agent.dispatch")
	AssertTrue(t, cfg.Enabled("agent.dispatch"))
}

func TestConfig_Config_Enabled_Bad(t *T) {
	cfg := (&Config{}).New()
	// Never-set feature is not enabled.
	AssertFalse(t, cfg.Enabled("agent.dispatch"))
	// Explicitly disabled feature is not enabled.
	cfg.Disable("agent.review")
	AssertFalse(t, cfg.Enabled("agent.review"))
	// Enabled-then-disabled reverts to not enabled.
	cfg.Enable("agent.dispatch")
	cfg.Disable("agent.dispatch")
	AssertFalse(t, cfg.Enabled("agent.dispatch"))
}

func TestConfig_Config_Enabled_Ugly(t *T) {
	var cfg Config
	// Zero-value Config: querying an unset feature is false, no panic.
	AssertFalse(t, cfg.Enabled("agent.dispatch"))
	// Enable lazily initialises the feature store on a zero-value Config.
	cfg.Enable("agent.dispatch")
	AssertTrue(t, cfg.Enabled("agent.dispatch"))
	// Disable flips it back off.
	cfg.Disable("agent.dispatch")
	AssertFalse(t, cfg.Enabled("agent.dispatch"))
}

func TestConfig_Config_EnabledFeatures_Good(t *T) {
	cfg := (&Config{}).New()
	cfg.Enable("agent.dispatch")
	cfg.Enable("agent.review")
	features := cfg.EnabledFeatures()
	AssertContains(t, features, "agent.dispatch")
	AssertContains(t, features, "agent.review")
}

func TestConfig_Config_EnabledFeatures_Bad(t *T) {
	cfg := (&Config{}).New()
	// Nothing enabled yet.
	AssertEmpty(t, cfg.EnabledFeatures())
	// Disabling features (without enabling) leaves the active set empty.
	cfg.Disable("agent.dispatch")
	cfg.Disable("agent.review")
	AssertEmpty(t, cfg.EnabledFeatures())
	// Enabling then disabling the same feature drops it from the active set.
	cfg.Enable("agent.dispatch")
	cfg.Disable("agent.dispatch")
	AssertEmpty(t, cfg.EnabledFeatures())
}

func TestConfig_Config_EnabledFeatures_Ugly(t *T) {
	cfg := (&Config{}).New()
	cfg.Enable("agent.dispatch")
	cfg.Disable("agent.review")
	features := cfg.EnabledFeatures()
	AssertContains(t, features, "agent.dispatch")
	AssertNotContains(t, features, "agent.review")
}

func TestConfig_NewConfigVar_Good(t *T) {
	v := NewConfigVar("homelab.lthn.sh")
	AssertTrue(t, v.IsSet())
	AssertEqual(t, "homelab.lthn.sh", v.Get())
}

func TestConfig_NewConfigVar_Bad(t *T) {
	v := NewConfigVar("")
	AssertTrue(t, v.IsSet())
	AssertEqual(t, "", v.Get())
}

func TestConfig_NewConfigVar_Ugly(t *T) {
	v := NewConfigVar[*App](nil)
	AssertTrue(t, v.IsSet())
	AssertNil(t, v.Get())
}

func TestConfig_ConfigVar_Get_Good(t *T) {
	v := NewConfigVar("codex")
	AssertEqual(t, "codex", v.Get())
	// Get reflects the latest Set value.
	v.Set("virgil")
	AssertEqual(t, "virgil", v.Get())
	// Get is non-destructive: repeated reads return the same value and keep it set.
	AssertEqual(t, "virgil", v.Get())
	AssertTrue(t, v.IsSet())
}

func TestConfig_ConfigVar_Get_Bad(t *T) {
	// Zero-value ConfigVar returns the type's zero and reports unset.
	var v ConfigVar[string]
	AssertEqual(t, "", v.Get())
	AssertFalse(t, v.IsSet())
	// A numeric ConfigVar zero-value is 0, also unset.
	var n ConfigVar[int]
	AssertEqual(t, 0, n.Get())
	AssertFalse(t, n.IsSet())
}

func TestConfig_ConfigVar_Get_Ugly(t *T) {
	v := NewConfigVar("session-token")
	v.Unset()
	AssertEqual(t, "", v.Get())
}

func TestConfig_ConfigVar_Set_Good(t *T) {
	var v ConfigVar[string]
	v.Set("codex")
	AssertTrue(t, v.IsSet())
	AssertEqual(t, "codex", v.Get())
}

func TestConfig_ConfigVar_Set_Bad(t *T) {
	var v ConfigVar[int]
	v.Set(0)
	AssertTrue(t, v.IsSet())
	AssertEqual(t, 0, v.Get())
}

func TestConfig_ConfigVar_Set_Ugly(t *T) {
	var v ConfigVar[*App]
	v.Set(nil)
	AssertTrue(t, v.IsSet())
	AssertNil(t, v.Get())
}

func TestConfig_ConfigVar_IsSet_Good(t *T) {
	v := NewConfigVar(true)
	AssertTrue(t, v.IsSet())
	// NewConfigVar marks the var set even when the value is the zero value.
	z := NewConfigVar(false)
	AssertTrue(t, z.IsSet())
	AssertFalse(t, z.Get())
	// Set on an existing var keeps it set.
	v.Set(false)
	AssertTrue(t, v.IsSet())
}

func TestConfig_ConfigVar_IsSet_Bad(t *T) {
	// Zero-value ConfigVar is not set.
	var v ConfigVar[bool]
	AssertFalse(t, v.IsSet())
	// Setting it (even to the zero value) marks it set.
	v.Set(false)
	AssertTrue(t, v.IsSet())
	// Unset reverts it to not-set.
	v.Unset()
	AssertFalse(t, v.IsSet())
}

func TestConfig_ConfigVar_IsSet_Ugly(t *T) {
	v := NewConfigVar(false)
	AssertTrue(t, v.IsSet(), "set false differs from unset")
	v.Unset()
	AssertFalse(t, v.IsSet())
}

func TestConfig_ConfigVar_Unset_Good(t *T) {
	v := NewConfigVar("session-token")
	v.Unset()
	AssertFalse(t, v.IsSet())
	AssertEqual(t, "", v.Get())
}

func TestConfig_ConfigVar_Unset_Bad(t *T) {
	var v ConfigVar[string]
	v.Unset()
	AssertFalse(t, v.IsSet())
	AssertEqual(t, "", v.Get())
}

func TestConfig_ConfigVar_Unset_Ugly(t *T) {
	v := NewConfigVar(&App{Name: "agent"})
	v.Unset()
	AssertFalse(t, v.IsSet())
	AssertNil(t, v.Get())
}
