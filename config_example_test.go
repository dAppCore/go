package core_test

import (
	. "dappco.re/go"
)

// ExampleConfig_Set sets a value through `Config.Set` for service configuration. Callers
// write untyped settings and read back typed values or enabled features.
func ExampleConfig_Set() {
	c := New()
	c.Config().Set("database.host", "localhost")
	c.Config().Set("database.port", 5432)

	Println(c.Config().String("database.host"))
	Println(c.Config().Int("database.port"))
	// Output:
	// localhost
	// 5432
}

// ExampleNewConfigVar_config initialises configuration values through `NewConfigVar` for
// service configuration. Callers write untyped settings and read back typed values or
// enabled features.
func ExampleNewConfigVar_config() {
	v := NewConfigVar("enabled")
	Println(v.Get())
	Println(v.IsSet())
	// Output:
	// enabled
	// true
}

// ExampleConfigVar_Set sets a value through `ConfigVar.Set` for service configuration.
// Callers write untyped settings and read back typed values or enabled features.
func ExampleConfigVar_Set() {
	var v ConfigVar[string]
	v.Set("enabled")
	Println(v.Get())
	Println(v.IsSet())
	// Output:
	// enabled
	// true
}

// ExampleConfigVar_Unset clears a value through `ConfigVar.Unset` for service
// configuration. Callers write untyped settings and read back typed values or enabled
// features.
func ExampleConfigVar_Unset() {
	v := NewConfigVar("enabled")
	v.Unset()
	Println(v.Get())
	Println(v.IsSet())
	// Output:
	//
	// false
}

// ExampleConfig_New creates an empty configuration with no enabled feature flags. Callers
// write untyped settings and read back typed values or enabled features.
func ExampleConfig_New() {
	cfg := (&Config{}).New()
	Println(cfg.EnabledFeatures())
	// Output: []
}

// ExampleConfigOptions groups settings and feature flags through `ConfigOptions` for
// service configuration. Callers write untyped settings and read back typed values or
// enabled features.
func ExampleConfigOptions() {
	opts := ConfigOptions{
		Settings: map[string]any{"host": "localhost"},
		Features: map[string]bool{"debug": true},
	}
	Println(opts.Settings["host"])
	Println(opts.Features["debug"])
	// Output:
	// localhost
	// true
}

// ExampleConfig_Get retrieves a value through `Config.Get` for service configuration.
// Callers write untyped settings and read back typed values or enabled features.
func ExampleConfig_Get() {
	cfg := (&Config{}).New()
	cfg.Set("host", "localhost")
	Println(cfg.Get("host").Value)
	Println(cfg.Get("missing").OK)
	// Output:
	// localhost
	// false
}

// ExampleConfig_String renders `Config.String` as a stable string for service
// configuration. Callers write untyped settings and read back typed values or enabled
// features.
func ExampleConfig_String() {
	cfg := (&Config{}).New()
	cfg.Set("host", "localhost")
	Println(cfg.String("host"))
	// Output: localhost
}

// ExampleConfig_Int reads an integer through `Config.Int` for service configuration.
// Callers write untyped settings and read back typed values or enabled features.
func ExampleConfig_Int() {
	cfg := (&Config{}).New()
	cfg.Set("port", 8080)
	Println(cfg.Int("port"))
	// Output: 8080
}

// ExampleConfig_Bool reads a boolean through `Config.Bool` for service configuration.
// Callers write untyped settings and read back typed values or enabled features.
func ExampleConfig_Bool() {
	cfg := (&Config{}).New()
	cfg.Set("debug", true)
	Println(cfg.Bool("debug"))
	// Output: true
}

// ExampleConfigGet reads a typed config value through `ConfigGet` for service
// configuration. Callers write untyped settings and read back typed values or enabled
// features.
func ExampleConfigGet() {
	cfg := (&Config{}).New()
	cfg.Set("host", "localhost")
	Println(ConfigGet[string](cfg, "host"))
	// Output: localhost
}

// ExampleConfig_Enable turns on named feature flags for service configuration. Callers
// write untyped settings and read back typed values or enabled features.
func ExampleConfig_Enable() {
	c := New()
	c.Config().Enable("dark-mode")
	c.Config().Enable("beta-features")

	Println(c.Config().Enabled("dark-mode"))
	features := c.Config().EnabledFeatures()
	SliceSort(features)
	Println(features)
	// Output:
	// true
	// [beta-features dark-mode]
}

// ExampleConfig_Disable turns off a previously enabled feature flag for service
// configuration. Callers write untyped settings and read back typed values or enabled
// features.
func ExampleConfig_Disable() {
	c := New()
	c.Config().Enable("debug")
	c.Config().Disable("debug")
	Println(c.Config().Enabled("debug"))
	// Output: false
}

// ExampleConfig_Enabled checks whether a feature is enabled through `Config.Enabled` for
// service configuration. Callers write untyped settings and read back typed values or
// enabled features.
func ExampleConfig_Enabled() {
	c := New()
	c.Config().Enable("debug")
	Println(c.Config().Enabled("debug"))
	// Output: true
}

// ExampleConfig_EnabledFeatures lists enabled feature flags through
// `Config.EnabledFeatures` for service configuration. Callers write untyped settings and
// read back typed values or enabled features.
func ExampleConfig_EnabledFeatures() {
	c := New()
	c.Config().Enable("beta")
	c.Config().Enable("debug")
	features := c.Config().EnabledFeatures()
	SliceSort(features)
	Println(features)
	// Output: [beta debug]
}

// ExampleConfigVar toggles a standalone config variable through `ConfigVar` for service
// configuration. Callers write untyped settings and read back typed values or enabled
// features.
func ExampleConfigVar() {
	v := NewConfigVar(42)
	Println(v.Get(), v.IsSet())

	v.Unset()
	Println(v.Get(), v.IsSet())
	// Output:
	// 42 true
	// 0 false
}

// ExampleConfig_Group scopes config keys under a group prefix through
// `Config.Group` — the view reads and writes "<group>.<key>" in the shared store.
func ExampleConfig_Group() {
	c := New()
	db := c.Config("database")
	db.Set("host", "localhost")

	Println(c.Config().String("database.host"))
	Println(db.String("host"))
	// Output:
	// localhost
	// localhost
}

// ExampleConfigVar_Get reads a typed config var through `ConfigVar.Get`.
func ExampleConfigVar_Get() {
	v := NewConfigVar("https://api.lthn.ai")
	Println(v.Get())
	// Output: https://api.lthn.ai
}

// ExampleConfigVar_IsSet reports whether a config var holds a value through `ConfigVar.IsSet`.
func ExampleConfigVar_IsSet() {
	v := NewConfigVar("ready")
	Println(v.IsSet())
	// Output: true
}

// ExampleFeature_Name returns the feature's name through `Feature.Name`.
func ExampleFeature_Name() {
	Println(New().Feature("dark-mode").Name())
	// Output: dark-mode
}

// ExampleFeature_Enable activates a feature through `Feature.Enable`.
func ExampleFeature_Enable() {
	c := New()
	c.Feature("beta").Enable()
	Println(c.Feature("beta").Enabled())
	// Output: true
}

// ExampleFeature_Disable deactivates a feature through `Feature.Disable`.
func ExampleFeature_Disable() {
	c := New()
	f := c.Feature("beta")
	f.Enable()
	f.Disable()
	Println(f.Enabled())
	// Output: false
}

// ExampleFeature_Enabled reports whether a feature is active through `Feature.Enabled`.
func ExampleFeature_Enabled() {
	Println(New().Feature("unset").Enabled())
	// Output: false
}

// ExampleConfig_Load merges a JSON file into the store with dotted keys.
func ExampleConfig_Load() {
	c := New()
	r := c.Config().Load("/nonexistent/config.json")
	Println(r.OK)
	// Output: false
}

// ExampleConfig_FromEnv imports prefixed environment variables:
// MYAPP_DATABASE_HOST becomes "database.host".
func ExampleConfig_FromEnv() {
	c := New()
	r := c.Config().FromEnv("NO_SUCH_PREFIX_XYZ_")
	Println(r.Int())
	// Output: 0
}
