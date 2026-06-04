// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the I18n primitive in i18n.go.
// Per AX-11 — c.I18n().Translate fires on every CLI command-line
// description render, every UI label boot, every error message that
// uses the message-id pattern. Default (no translator) returns the key
// — the bench harness gates the fast-path floor. AddLocales / Locales
// are boot-only but worth gating.
//
// Run:    go test -bench='BenchmarkI18n' -benchmem -run='^$' .

package core_test

import (
	"testing/fstest"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE.
var (
	i18nSinkResult  Result
	i18nSinkString  string
	i18nSinkStrings []string
)

// benchTranslator is a fake Translator for the with-translator paths.
type benchTranslator struct{ lang string }

func (t *benchTranslator) Translate(messageID string, args ...any) Result {
	return Result{messageID + ":translated", true}
}

func (t *benchTranslator) SetLanguage(lang string) Result { t.lang = lang; return Ok(nil) }
func (t *benchTranslator) Language() string               { return t.lang }
func (t *benchTranslator) AvailableLanguages() []string   { return []string{"en", "en-GB", "de"} }

// --- Translate (default + with translator) ---

func BenchmarkI18n_Translate_Default(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().Translate("cmd.deploy.description")
	}
}

func BenchmarkI18n_Translate_WithTranslator(b *B) {
	c := New()
	c.I18n().SetTranslator(&benchTranslator{lang: "en-GB"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().Translate("cmd.deploy.description")
	}
}

func BenchmarkI18n_Translate_WithArgs(b *B) {
	c := New()
	c.I18n().SetTranslator(&benchTranslator{lang: "en-GB"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().Translate("cmd.deploy.message", "homelab", 9000)
	}
}

// --- Language accessors ---

func BenchmarkI18n_Language_Default(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkString = c.I18n().Language()
	}
}

func BenchmarkI18n_SetLanguage_NoTranslator(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().SetLanguage("en-GB")
	}
}

func BenchmarkI18n_SetLanguage_WithTranslator(b *B) {
	c := New()
	c.I18n().SetTranslator(&benchTranslator{lang: "en"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().SetLanguage("en-GB")
	}
}

func BenchmarkI18n_AvailableLanguages_Default(b *B) {
	c := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkStrings = c.I18n().AvailableLanguages()
	}
}

func BenchmarkI18n_AvailableLanguages_WithTranslator(b *B) {
	c := New()
	c.I18n().SetTranslator(&benchTranslator{lang: "en"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkStrings = c.I18n().AvailableLanguages()
	}
}

// --- Locale collection ---

func BenchmarkI18n_AddLocales(b *B) {
	c := New()
	r := Mount(FS(fstest.MapFS{
		"en.yaml": &fstest.MapFile{Data: []byte("hello: Hello")},
	}), ".")
	emb := r.Value.(*Embed)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.I18n().AddLocales(emb)
	}
}

func BenchmarkI18n_Locales(b *B) {
	c := New()
	r := Mount(FS(fstest.MapFS{
		"en.yaml": &fstest.MapFile{Data: []byte("hello: Hello")},
	}), ".")
	c.I18n().AddLocales(r.Value.(*Embed))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().Locales()
	}
}

// --- Translator setter + getter ---

func BenchmarkI18n_SetTranslator(b *B) {
	c := New()
	tr := &benchTranslator{lang: "en"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.I18n().SetTranslator(tr)
	}
}

func BenchmarkI18n_Translator_Hit(b *B) {
	c := New()
	c.I18n().SetTranslator(&benchTranslator{lang: "en"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		i18nSinkResult = c.I18n().Translator()
	}
}
