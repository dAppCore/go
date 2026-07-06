package core_test

import (
	. "dappco.re/go"
)

// --- NewOptions ---

func TestOptions_NewOptions_Good(t *T) {
	opts := NewOptions(
		Option{Key: "name", Value: "brain"},
		Option{Key: "port", Value: 8080},
	)
	AssertEqual(t, 2, opts.Len())
}

func TestOptions_NewOptions_Empty_Good(t *T) {
	opts := NewOptions()
	AssertEqual(t, 0, opts.Len())
	AssertFalse(t, opts.Has("anything"))
}

// --- Options.Set ---

func TestOptions_Set_Good(t *T) {
	opts := NewOptions()
	opts.Set("name", "brain")
	AssertEqual(t, "brain", opts.String("name"))
}

func TestOptions_Set_Update_Good(t *T) {
	opts := NewOptions(Option{Key: "name", Value: "old"})
	opts.Set("name", "new")
	AssertEqual(t, "new", opts.String("name"))
	AssertEqual(t, 1, opts.Len())
}

// --- Options.Get ---

func TestOptions_Get_Good(t *T) {
	opts := NewOptions(
		Option{Key: "name", Value: "brain"},
		Option{Key: "port", Value: 8080},
	)
	r := opts.Get("name")
	AssertTrue(t, r.OK)
	AssertEqual(t, "brain", r.Value)
}

func TestOptions_Get_Bad(t *T) {
	opts := NewOptions(Option{Key: "name", Value: "brain"})
	r := opts.Get("missing")
	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

// --- Options.Has ---

func TestOptions_Has_Good(t *T) {
	opts := NewOptions(Option{Key: "debug", Value: true})
	AssertTrue(t, opts.Has("debug"))
	AssertFalse(t, opts.Has("missing"))
}

// --- Options.String ---

func TestOptions_String_Good(t *T) {
	opts := NewOptions(
		Option{Key: "name", Value: "brain"},
		Option{Key: "host", Value: "localhost"},
		Option{Key: "empty", Value: ""},
	)
	AssertEqual(t, "brain", opts.String("name"))
	AssertEqual(t, "localhost", opts.String("host"))
	// A present key whose value is the empty string returns "".
	AssertEqual(t, "", opts.String("empty"))
	AssertTrue(t, opts.Has("empty"))
	// A missing key also returns "" — indistinguishable via String alone.
	AssertEqual(t, "", opts.String("missing"))
	AssertFalse(t, opts.Has("missing"))
}

func TestOptions_String_Bad(t *T) {
	opts := NewOptions(Option{Key: "port", Value: 8080})
	AssertEqual(t, "", opts.String("port"))
	AssertEqual(t, "", opts.String("missing"))
}

// --- Options.Int ---

func TestOptions_Int_Good(t *T) {
	opts := NewOptions(
		Option{Key: "port", Value: 8080},
		Option{Key: "zero", Value: 0},
		Option{Key: "neg", Value: -42},
	)
	AssertEqual(t, 8080, opts.Int("port"))
	// A present key holding 0 is returned as 0, and the key still exists.
	AssertEqual(t, 0, opts.Int("zero"))
	AssertTrue(t, opts.Has("zero"))
	AssertEqual(t, -42, opts.Int("neg"))
	// Missing key yields the zero default.
	AssertEqual(t, 0, opts.Int("missing"))
}

func TestOptions_Int_Bad(t *T) {
	opts := NewOptions(Option{Key: "name", Value: "brain"})
	AssertEqual(t, 0, opts.Int("name"))
	AssertEqual(t, 0, opts.Int("missing"))
}

// --- Options.Bool ---

func TestOptions_Bool_Good(t *T) {
	opts := NewOptions(
		Option{Key: "debug", Value: true},
		Option{Key: "verbose", Value: false},
	)
	AssertTrue(t, opts.Bool("debug"))
	// A present key holding false returns false, and the key still exists.
	AssertFalse(t, opts.Bool("verbose"))
	AssertTrue(t, opts.Has("verbose"))
	// Missing key yields the false default.
	AssertFalse(t, opts.Bool("missing"))
}

func TestOptions_Bool_Bad(t *T) {
	opts := NewOptions(Option{Key: "name", Value: "brain"})
	AssertFalse(t, opts.Bool("name"))
	AssertFalse(t, opts.Bool("missing"))
}

// --- Options.Items ---

func TestOptions_Items_Good(t *T) {
	opts := NewOptions(Option{Key: "a", Value: 1}, Option{Key: "b", Value: 2})
	items := opts.Items()
	AssertLen(t, items, 2)
}

// --- Options with typed struct ---

func TestOptions_Get_Struct_Good(t *T) {
	type BrainConfig struct {
		Name       string
		OllamaURL  string
		Collection string
	}
	cfg := BrainConfig{Name: "brain", OllamaURL: "http://localhost:11434", Collection: "openbrain"}
	opts := NewOptions(Option{Key: "config", Value: cfg})

	r := opts.Get("config")
	AssertTrue(t, r.OK)
	bc, ok := r.Value.(BrainConfig)
	AssertTrue(t, ok)
	AssertEqual(t, "brain", bc.Name)
	AssertEqual(t, "http://localhost:11434", bc.OllamaURL)
}

// --- Result ---

func TestOptions_Result_New_Good(t *T) {
	// Single non-error arg: value is stored and OK flips true.
	r := Result{}.New("value")
	AssertEqual(t, "value", r.Value)
	AssertTrue(t, r.OK)

	// A non-string, non-error value passes through unchanged.
	num := Result{}.New(42)
	AssertEqual(t, 42, num.Value)
	AssertTrue(t, num.OK)

	// Two-arg (value, nil-error) form: value wins, OK true.
	var noErr error
	pair := Result{}.New("file", noErr)
	AssertEqual(t, "file", pair.Value)
	AssertTrue(t, pair.OK)
}

func TestOptions_Result_New_Error_Bad(t *T) {
	err := E("test", "failed", nil)
	r := Result{}.New(err)
	AssertFalse(t, r.OK)
	AssertEqual(t, err, r.Value)
}

// --- WithOption ---

func TestOptions_WithOption_Good(t *T) {
	c := New(
		WithOption("name", "myapp"),
		WithOption("port", 8080),
	)
	AssertEqual(t, "myapp", c.App().Name)
	AssertEqual(t, 8080, c.Options().Int("port"))
}

func TestOptions_NewOptions_Bad(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"}, Option{Key: "agent", Value: "hades"})

	AssertEqual(t, 2, opts.Len())
	AssertEqual(t, "codex", opts.String("agent"))
}

func TestOptions_NewOptions_Ugly(t *T) {
	items := []Option{{Key: "agent", Value: "codex"}}
	opts := NewOptions(items...)
	items[0].Value = "hades"

	AssertEqual(t, "codex", opts.String("agent"))
}

func TestOptions_Options_Bool_Good(t *T) {
	opts := NewOptions(
		Option{Key: "enabled", Value: true},
		Option{Key: "disabled", Value: false},
	)

	AssertTrue(t, opts.Bool("enabled"))
	AssertFalse(t, opts.Bool("disabled"))
	AssertTrue(t, opts.Has("disabled"))
}

func TestOptions_Options_Bool_Bad(t *T) {
	opts := NewOptions(
		Option{Key: "enabled", Value: "true"},
		Option{Key: "count", Value: 1},
	)

	// Bool does no coercion: the string "true" is the wrong type -> false.
	AssertFalse(t, opts.Bool("enabled"))
	// The key still exists even though Bool returns the false default.
	AssertTrue(t, opts.Has("enabled"))
	// An int value is also not a bool.
	AssertFalse(t, opts.Bool("count"))
	// Missing key returns false too.
	AssertFalse(t, opts.Bool("missing"))
}

func TestOptions_Options_Bool_Ugly(t *T) {
	opts := NewOptions(Option{Key: "", Value: true})

	// The empty string is a legitimate key.
	AssertTrue(t, opts.Bool(""))
	AssertTrue(t, opts.Has(""))
	// A non-empty lookup still misses and returns the false default.
	AssertFalse(t, opts.Bool("enabled"))
	AssertFalse(t, opts.Has("enabled"))
}

func TestOptions_Options_Get_Good(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"})
	r := opts.Get("agent")

	AssertTrue(t, r.OK)
	AssertEqual(t, "codex", r.Value)
}

func TestOptions_Options_Get_Bad(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"})
	r := opts.Get("missing")

	AssertFalse(t, r.OK)
	AssertNil(t, r.Value)
}

func TestOptions_Options_Get_Ugly(t *T) {
	opts := NewOptions(Option{Key: "", Value: "empty-key"})
	r := opts.Get("")

	AssertTrue(t, r.OK)
	AssertEqual(t, "empty-key", r.Value)
}

func TestOptions_Options_Has_Good(t *T) {
	opts := NewOptions(
		Option{Key: "agent", Value: "codex"},
		Option{Key: "nilval", Value: nil},
	)

	AssertTrue(t, opts.Has("agent"))
	// Has reports key presence, not value-ness: a nil value still counts.
	AssertTrue(t, opts.Has("nilval"))
	// And the matching Get confirms the key resolves with OK true.
	AssertTrue(t, opts.Get("agent").OK)
}

func TestOptions_Options_Has_Bad(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"})

	AssertFalse(t, opts.Has("missing"))
	// Lookups are exact, case-sensitive string matches.
	AssertFalse(t, opts.Has("Agent"))
	AssertFalse(t, opts.Has("agent "))
	// The corresponding Get also reports OK false for a miss.
	AssertFalse(t, opts.Get("missing").OK)
}

func TestOptions_Options_Has_Ugly(t *T) {
	opts := NewOptions(Option{Key: "", Value: "empty-key"})

	AssertTrue(t, opts.Has(""))
	// Only the empty key exists; the value text is not itself a key.
	AssertFalse(t, opts.Has("empty-key"))
}

func TestOptions_Options_Int_Good(t *T) {
	opts := NewOptions(
		Option{Key: "port", Value: 8080},
		Option{Key: "zero", Value: 0},
		Option{Key: "neg", Value: -42},
	)

	AssertEqual(t, 8080, opts.Int("port"))
	// A present 0 is returned as 0; the key still exists.
	AssertEqual(t, 0, opts.Int("zero"))
	AssertTrue(t, opts.Has("zero"))
	AssertEqual(t, -42, opts.Int("neg"))
}

func TestOptions_Options_Int_Bad(t *T) {
	opts := NewOptions(
		Option{Key: "port", Value: "8080"},
		Option{Key: "big", Value: int64(8080)},
	)

	// String is not int — no parsing, returns the zero default.
	AssertEqual(t, 0, opts.Int("port"))
	// int64 does not satisfy the int type assertion either.
	AssertEqual(t, 0, opts.Int("big"))
	// Missing key also yields 0.
	AssertEqual(t, 0, opts.Int("missing"))
}

func TestOptions_Options_Int_Ugly(t *T) {
	opts := NewOptions(
		Option{Key: "", Value: -1},
		Option{Key: "max", Value: 2147483647},
	)

	// Negative value under the empty key resolves correctly.
	AssertEqual(t, -1, opts.Int(""))
	// Large positive int round-trips unchanged.
	AssertEqual(t, 2147483647, opts.Int("max"))
}

func TestOptions_Options_Items_Good(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"}, Option{Key: "region", Value: "homelab"})

	items := opts.Items()
	AssertLen(t, items, 2)
	// Insertion order is preserved.
	AssertEqual(t, "agent", items[0].Key)
	AssertEqual(t, "codex", items[0].Value)
	AssertEqual(t, "region", items[1].Key)
	AssertEqual(t, "homelab", items[1].Value)
	AssertEqual(t, []Option{{Key: "agent", Value: "codex"}, {Key: "region", Value: "homelab"}}, items)
}

func TestOptions_Options_Items_Bad(t *T) {
	opts := NewOptions()

	items := opts.Items()
	AssertEmpty(t, items)
	AssertLen(t, items, 0)
	// Items always allocates a slice — never returns nil, even when empty.
	AssertNotNil(t, items)
}

func TestOptions_Options_Items_Ugly(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"})
	items := opts.Items()
	items[0].Value = "hades"

	AssertEqual(t, "codex", opts.String("agent"))
}

func TestOptions_Options_Len_Good(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"}, Option{Key: "debug", Value: true})

	AssertEqual(t, 2, opts.Len())
	// Setting a new key grows Len.
	opts.Set("region", "homelab")
	AssertEqual(t, 3, opts.Len())
	// Updating an existing key is in-place — Len is unchanged.
	opts.Set("agent", "hades")
	AssertEqual(t, 3, opts.Len())
}

func TestOptions_Options_Len_Bad(t *T) {
	opts := NewOptions()

	AssertEqual(t, 0, opts.Len())
	// A miss lookup is a read-only no-op — Len stays 0.
	AssertFalse(t, opts.Has("x"))
	AssertEqual(t, 0, opts.Len())
	// The first Set grows it to 1.
	opts.Set("first", 1)
	AssertEqual(t, 1, opts.Len())
}

func TestOptions_Options_Len_Ugly(t *T) {
	opts := NewOptions(Option{Key: "agent", Value: "codex"})
	opts.Set("agent", "hades")

	AssertEqual(t, 1, opts.Len())
}

func TestOptions_Options_Set_Good(t *T) {
	opts := NewOptions()
	opts.Set("agent", "codex")

	AssertEqual(t, "codex", opts.String("agent"))
}

func TestOptions_Options_Set_Bad(t *T) {
	var opts *Options

	// Set has a pointer receiver and dereferences the nil pointer -> panic.
	AssertPanics(t, func() { opts.Set("agent", "codex") })

	// A valid Options does not panic and the write takes effect.
	valid := NewOptions()
	AssertNotPanics(t, func() { valid.Set("agent", "codex") })
	AssertEqual(t, "codex", valid.String("agent"))
	AssertEqual(t, 1, valid.Len())
}

func TestOptions_Options_Set_Ugly(t *T) {
	opts := NewOptions(Option{Key: "", Value: "before"})
	opts.Set("", "after")

	AssertEqual(t, 1, opts.Len())
	AssertEqual(t, "after", opts.String(""))
}

func TestOptions_Options_String_Good(t *T) {
	opts := NewOptions(
		Option{Key: "agent", Value: "codex"},
		Option{Key: "empty", Value: ""},
	)

	AssertEqual(t, "codex", opts.String("agent"))
	// A present empty-string value returns "" and the key still exists.
	AssertEqual(t, "", opts.String("empty"))
	AssertTrue(t, opts.Has("empty"))
	// Missing key returns the "" default.
	AssertEqual(t, "", opts.String("missing"))
}

func TestOptions_Options_String_Bad(t *T) {
	opts := NewOptions(
		Option{Key: "port", Value: 8080},
		Option{Key: "flag", Value: true},
	)

	// Wrong type yields "" — no fmt-style stringification of an int.
	AssertEqual(t, "", opts.String("port"))
	// The key exists; only the typed accessor returns the default.
	AssertTrue(t, opts.Has("port"))
	// A bool value is also not a string.
	AssertEqual(t, "", opts.String("flag"))
	// Missing key returns "" too.
	AssertEqual(t, "", opts.String("missing"))
}

func TestOptions_Options_String_Ugly(t *T) {
	opts := NewOptions(Option{Key: "", Value: "empty-key"})

	// The empty string is a usable key.
	AssertEqual(t, "empty-key", opts.String(""))
	// The value text is not itself a key — that lookup misses.
	AssertEqual(t, "", opts.String("empty-key"))
	AssertFalse(t, opts.Has("empty-key"))
}

// --- Options.Float64 ---

func TestOptions_Options_Float64_Good(t *T) {
	opts := NewOptions(
		Option{Key: "weight", Value: 0.75},
		Option{Key: "zero", Value: 0.0},
		Option{Key: "neg", Value: -2.5},
	)
	AssertEqual(t, 0.75, opts.Float64("weight"))
	// A present 0.0 returns 0; the key still exists.
	AssertEqual(t, 0.0, opts.Float64("zero"))
	AssertTrue(t, opts.Has("zero"))
	AssertEqual(t, -2.5, opts.Float64("neg"))
	// Missing key returns the 0 default.
	AssertEqual(t, 0.0, opts.Float64("missing"))
}

func TestOptions_Options_Float64_Bad(t *T) {
	opts := NewOptions(Option{Key: "name", Value: "brain"})
	AssertEqual(t, 0.0, opts.Float64("name"))
	AssertEqual(t, 0.0, opts.Float64("missing"))
}

func TestOptions_Options_Float64_Ugly(t *T) {
	// int and float32 promote to float64 — JSON-decoded numbers work.
	opts := NewOptions(
		Option{Key: "i", Value: 42},
		Option{Key: "i64", Value: int64(99)},
		Option{Key: "f32", Value: float32(1.5)},
	)
	AssertEqual(t, 42.0, opts.Float64("i"))
	AssertEqual(t, 99.0, opts.Float64("i64"))
	AssertEqual(t, 1.5, opts.Float64("f32"))
}

// --- Options.Duration ---

func TestOptions_Options_Duration_Good(t *T) {
	opts := NewOptions(
		Option{Key: "timeout", Value: 5 * Second},
		Option{Key: "zero", Value: Duration(0)},
		Option{Key: "ms", Value: 250 * Millisecond},
	)
	AssertEqual(t, 5*Second, opts.Duration("timeout"))
	// A present zero Duration returns 0; the key still exists.
	AssertEqual(t, Duration(0), opts.Duration("zero"))
	AssertTrue(t, opts.Has("zero"))
	AssertEqual(t, 250*Millisecond, opts.Duration("ms"))
	// Missing key returns the 0 default.
	AssertEqual(t, Duration(0), opts.Duration("missing"))
}

func TestOptions_Options_Duration_Bad(t *T) {
	opts := NewOptions(Option{Key: "name", Value: 12345})
	AssertEqual(t, Duration(0), opts.Duration("name"))
	AssertEqual(t, Duration(0), opts.Duration("missing"))
}

func TestOptions_Options_Duration_Ugly(t *T) {
	// String values parse via ParseDuration; bad strings return 0.
	opts := NewOptions(
		Option{Key: "good", Value: "250ms"},
		Option{Key: "bad", Value: "not-a-duration"},
	)
	AssertEqual(t, 250*Millisecond, opts.Duration("good"))
	AssertEqual(t, Duration(0), opts.Duration("bad"))
}

func TestOptions_Result_New_Bad(t *T) {
	r := Result{}.New(AnError)

	AssertFalse(t, r.OK)
	AssertEqual(t, AnError, r.Value)
}

func TestOptions_Result_New_Ugly(t *T) {
	r := Result{Value: "existing", OK: true}.New()

	AssertTrue(t, r.OK)
	AssertEqual(t, "existing", r.Value)
}
