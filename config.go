// SPDX-License-Identifier: EUPL-1.2

// Settings, feature flags, and typed configuration for the Core framework.

package core

import ()

// ConfigVar is a variable that can be set, unset, and queried for its state.
//
//	host := core.NewConfigVar("homelab.lthn.sh")
//	if host.IsSet() { core.Println(host.Get()) }
type ConfigVar[T any] struct {
	val T
	set bool
}

// Get returns the current value.
//
//	val := v.Get()
func (v *ConfigVar[T]) Get() T { return v.val }

// Set sets the value and marks it as explicitly set.
//
//	v.Set(true)
func (v *ConfigVar[T]) Set(val T) { v.val = val; v.set = true }

// IsSet returns true if the value was explicitly set (distinguishes "set to false" from "never set").
//
//	if v.IsSet() { /* explicitly configured */ }
func (v *ConfigVar[T]) IsSet() bool { return v.set }

// Unset resets to zero value and marks as not set.
//
//	v.Unset()
//	v.IsSet()  // false
func (v *ConfigVar[T]) Unset() {
	v.set = false
	var zero T
	v.val = zero
}

// NewConfigVar creates a ConfigVar with an initial value marked as set.
//
//	debug := core.NewConfigVar(true)
func NewConfigVar[T any](val T) ConfigVar[T] {
	return ConfigVar[T]{val: val, set: true}
}

// ConfigOptions holds configuration data.
//
//	opts := core.ConfigOptions{
//	    Settings: map[string]any{"config.host": "homelab.lthn.sh"},
//	    Features: map[string]bool{"agentic": true},
//	}
//	_ = opts
type ConfigOptions struct {
	Settings map[string]any
	Features map[string]bool
}

func (o *ConfigOptions) init() {
	if o.Settings == nil {
		o.Settings = make(map[string]any)
	}
	if o.Features == nil {
		o.Features = make(map[string]bool)
	}
}

// Config holds configuration settings and feature flags.
//
//	cfg := (&core.Config{}).New()
//	cfg.Set("config.host", "homelab.lthn.sh")
//	core.Println(cfg.String("config.host"))
type Config struct {
	*ConfigOptions
	features *Registry[bool] // live feature-flag store; ConfigOptions.Features is declarative init input
	root     *Config         // non-nil for a Group view; the store + lock live on root
	prefix   string          // group key prefix ("" for the root Config)
	mu       RWMutex
}

// New initialises a Config with empty settings and features.
//
//	cfg := (&core.Config{}).New()
func (e *Config) New() *Config {
	e.ConfigOptions = &ConfigOptions{}
	e.ConfigOptions.init()
	e.features = NewRegistry[bool]()
	return e
}

// featureFlags returns the live feature-flag Registry, lazily creating it and
// seeding once from the declarative ConfigOptions.Features input. The Registry
// is the source of truth for Enable/Disable/Enabled; ConfigOptions.Features is
// init-only input and is not mutated by later Enable/Disable calls.
func (e *Config) featureFlags() *Registry[bool] {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.features == nil {
		e.features = NewRegistry[bool]()
		if e.ConfigOptions != nil && e.ConfigOptions.Features != nil {
			for k, v := range e.ConfigOptions.Features {
				e.features.Set(k, v)
			}
		}
	}
	return e.features
}

// lock returns the mutex guarding the underlying store. A Group view shares its
// root's lock so concurrent access to the shared Settings map stays serialised.
func (e *Config) lock() *RWMutex {
	if e.root != nil {
		return &e.root.mu
	}
	return &e.mu
}

// key applies the group prefix to a raw key ("" prefix → key unchanged).
func (e *Config) key(k string) string {
	if e.prefix == "" {
		return k
	}
	return Concat(e.prefix, ".", k)
}

// Group returns a Config view scoped under name: Set/Get/String/Int/Bool and the
// feature methods operate on "<name>.<key>" in the shared underlying store.
// Nested groups compose; the root Config's ungrouped keys are unaffected.
//
//	db := c.Config("database")
//	db.Set("host", "localhost")          // stores "database.host"
//	c.Config().String("database.host")   // "localhost"
func (e *Config) Group(name string) *Config {
	root := e
	if e.root != nil {
		root = e.root
	}
	return &Config{
		ConfigOptions: root.ConfigOptions,
		features:      root.featureFlags(),
		root:          root,
		prefix:        e.key(name),
	}
}

// Set stores a configuration value by key.
//
//	cfg := (&core.Config{}).New()
//	cfg.Set("config.host", "homelab.lthn.sh")
func (e *Config) Set(key string, val any) {
	mu := e.lock()
	mu.Lock()
	if e.ConfigOptions == nil {
		e.ConfigOptions = &ConfigOptions{}
	}
	e.ConfigOptions.init()
	e.Settings[e.key(key)] = val
	mu.Unlock()
}

// Get retrieves a configuration value by key.
//
//	cfg := (&core.Config{}).New()
//	cfg.Set("config.host", "homelab.lthn.sh")
//	r := cfg.Get("config.host")
//	if r.OK { core.Println(r.Value.(string)) }
func (e *Config) Get(key string) Result {
	mu := e.lock()
	mu.RLock()
	defer mu.RUnlock()
	if e.ConfigOptions == nil || e.Settings == nil {
		return Result{}
	}
	val, ok := e.Settings[e.key(key)]
	if !ok {
		return Result{}
	}
	return Result{val, true}
}

// String retrieves a string config value (empty string if missing).
//
//	host := c.Config().String("database.host")
func (e *Config) String(key string) string { return ConfigGet[string](e, key) }

// Int retrieves an int config value (0 if missing).
//
//	port := c.Config().Int("database.port")
func (e *Config) Int(key string) int { return ConfigGet[int](e, key) }

// Bool retrieves a bool config value (false if missing).
//
//	debug := c.Config().Bool("debug")
func (e *Config) Bool(key string) bool { return ConfigGet[bool](e, key) }

// ConfigGet retrieves a typed configuration value.
//
//	cfg := (&core.Config{}).New()
//	cfg.Set("config.host", "homelab.lthn.sh")
//	host := core.ConfigGet[string](cfg, "config.host")
//	core.Println(host)
func ConfigGet[T any](e *Config, key string) T {
	r := e.Get(key)
	if !r.OK {
		var zero T
		return zero
	}
	typed, _ := r.Value.(T)
	return typed
}

// --- Feature Flags ---

// Enable activates a feature flag.
//
//	c.Config().Enable("dark-mode")
func (e *Config) Enable(feature string) {
	e.featureFlags().Set(e.key(feature), true)
}

// Disable deactivates a feature flag.
//
//	c.Config().Disable("dark-mode")
func (e *Config) Disable(feature string) {
	e.featureFlags().Set(e.key(feature), false)
}

// Enabled returns true if a feature flag is active.
//
//	if c.Config().Enabled("dark-mode") { ... }
func (e *Config) Enabled(feature string) bool {
	r := e.featureFlags().Get(e.key(feature))
	return r.OK && r.Value.(bool)
}

// Feature is a keyed handle to a single feature flag, bound to a Config.
//
//	f := c.Feature("dark-mode")
//	f.Enable()
//	if f.Enabled() { core.Println("on") }
type Feature struct {
	cfg  *Config
	name string
}

// Name returns the feature's name.
//
//	c.Feature("dark-mode").Name() // "dark-mode"
func (f Feature) Name() string { return f.name }

// Enable activates the feature.
//
//	c.Feature("dark-mode").Enable()
func (f Feature) Enable() { f.cfg.Enable(f.name) }

// Disable deactivates the feature.
//
//	c.Feature("dark-mode").Disable()
func (f Feature) Disable() { f.cfg.Disable(f.name) }

// Enabled reports whether the feature is active.
//
//	if c.Feature("dark-mode").Enabled() { core.Println("on") }
func (f Feature) Enabled() bool { return f.cfg.Enabled(f.name) }

// EnabledFeatures returns all active feature flag names.
//
//	features := c.Config().EnabledFeatures()
func (e *Config) EnabledFeatures() []string {
	ff := e.featureFlags()
	result := make([]string, 0, ff.Len())
	ff.Each(func(name string, on bool) {
		if on {
			result = append(result, name)
		}
	})
	if len(result) == 0 {
		return nil // byte-identical to the old var-nil behaviour
	}
	return result
}
