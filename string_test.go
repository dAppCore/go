package core_test

import (
	. "dappco.re/go"
)

// --- String Operations ---

func TestString_HasPrefix_Good(t *T) {
	AssertTrue(t, HasPrefix("--verbose", "--"))
	AssertTrue(t, HasPrefix("-v", "-"))
	AssertFalse(t, HasPrefix("hello", "-"))
}

func TestString_HasSuffix_Good(t *T) {
	AssertTrue(t, HasSuffix("test.go", ".go"))
	AssertFalse(t, HasSuffix("test.go", ".py"))
	AssertTrue(t, HasSuffix("archive.tar.gz", ".gz"))
	AssertTrue(t, HasSuffix("README", "README")) // suffix equals whole string
}

func TestString_TrimPrefix_Good(t *T) {
	AssertEqual(t, "verbose", TrimPrefix("--verbose", "--"))
	AssertEqual(t, "hello", TrimPrefix("hello", "--"))
	AssertEqual(t, "-verbose", TrimPrefix("--verbose", "-")) // removes only one occurrence
	AssertEqual(t, "", TrimPrefix("--", "--"))               // exact match drains to empty
}

func TestString_TrimSuffix_Good(t *T) {
	AssertEqual(t, "test", TrimSuffix("test.go", ".go"))
	AssertEqual(t, "test.go", TrimSuffix("test.go", ".py"))
	AssertEqual(t, "archive.tar", TrimSuffix("archive.tar.gz", ".gz"))
	AssertEqual(t, "", TrimSuffix(".go", ".go")) // exact match drains to empty
}

func TestString_Contains_Good(t *T) {
	AssertTrue(t, Contains("hello world", "world"))
	AssertFalse(t, Contains("hello world", "mars"))
	AssertTrue(t, Contains("hello", "hello")) // full-string match
	AssertTrue(t, Contains("café", "é"))      // multibyte substring
}

func TestString_Split_Good(t *T) {
	AssertEqual(t, []string{"a", "b", "c"}, Split("a/b/c", "/"))
	AssertEqual(t, []string{"a", "", "b"}, Split("a//b", "/"))         // empties between adjacent seps kept
	AssertEqual(t, []string{"", "a", ""}, Split("/a/", "/"))           // leading/trailing empties kept
	AssertEqual(t, []string{"noseparator"}, Split("noseparator", "/")) // sep absent → single element
}

func TestString_SplitN_Good(t *T) {
	AssertEqual(t, []string{"key", "value=extra"}, SplitN("key=value=extra", "=", 2))
	AssertEqual(t, []string{"a", "b", "c"}, SplitN("a=b=c", "=", -1)) // n<0 → unlimited
	AssertEqual(t, []string{"a=b=c"}, SplitN("a=b=c", "=", 1))        // n=1 → whole string, no split
}

func TestString_Join_Good(t *T) {
	AssertEqual(t, "a/b/c", Join("/", "a", "b", "c"))
	AssertEqual(t, "a", Join("/", "a"))             // single part → no separator
	AssertEqual(t, "a..b", Join(".", "a", "", "b")) // empty part still separated
	AssertEqual(t, "cmd.deploy.description", Join(".", "cmd", "deploy", "description"))
}

func TestString_Replace_Good(t *T) {
	AssertEqual(t, "deploy.to.homelab", Replace("deploy/to/homelab", "/", "."))
	AssertEqual(t, "a-b-c", Replace("a_b_c", "_", "-")) // all occurrences replaced
	AssertEqual(t, "bbb", Replace("aaa", "a", "b"))     // each char rewritten
	AssertEqual(t, "", Replace("aaa", "a", ""))         // replace-with-empty deletes matches
}

func TestString_Lower_Good(t *T) {
	AssertEqual(t, "hello", Lower("HELLO"))
	AssertEqual(t, "hello", Lower("hello"))   // already lower → unchanged
	AssertEqual(t, "go123!", Lower("GO123!")) // digits and punctuation preserved
	AssertEqual(t, "café", Lower("CAFÉ"))     // non-ASCII path: É → é
}

func TestString_Upper_Good(t *T) {
	AssertEqual(t, "HELLO", Upper("hello"))
	AssertEqual(t, "HELLO", Upper("HELLO"))   // already upper → unchanged
	AssertEqual(t, "GO123!", Upper("go123!")) // digits and punctuation preserved
	AssertEqual(t, "CAFÉ", Upper("café"))     // non-ASCII path: é → É
}

func TestString_Trim_Good(t *T) {
	AssertEqual(t, "hello", Trim("  hello  "))
	AssertEqual(t, "hello world", Trim("  hello world  ")) // interior whitespace kept
	AssertEqual(t, "hello", Trim("\t\n hello \r\n"))       // tabs/newlines/returns trimmed
	AssertEqual(t, "", Trim("   "))                        // all whitespace → empty
}

func TestString_TrimCutset_Good(t *T) {
	AssertEqual(t, "path", TrimCutset("//path//", "/"))
	AssertEqual(t, "name", TrimCutset("[name]", "[]"))
	AssertEqual(t, "core", TrimCutset("xxcorexx", "x"))       // repeated cut char both ends
	AssertEqual(t, "mid/path", TrimCutset("/mid/path/", "/")) // interior cutset chars preserved
}

func TestString_TrimCutset_Bad(t *T) {
	AssertEqual(t, "", TrimCutset("////", "/"))
	AssertEqual(t, "hello", TrimCutset("hello", ""))
	AssertEqual(t, "hello", TrimCutset("hello", "xyz")) // none of cutset present → unchanged
	AssertEqual(t, "", TrimCutset("", "/"))             // empty input → empty
}

func TestString_TrimCutset_Ugly(t *T) {
	AssertEqual(t, "abc", TrimCutset("xyzabcxyz", "xyz"))
	AssertEqual(t, "abc", TrimCutset("zyxabczxy", "xyz")) // cutset is a set; order irrelevant
	AssertEqual(t, "", TrimCutset("xyzyx", "xyz"))        // whole string made of cutset chars
}

func TestString_TrimLeft_Good(t *T) {
	AssertEqual(t, "path", TrimLeft("///path", "/"))
	AssertEqual(t, "verbose", TrimLeft("---verbose", "-"))
	AssertEqual(t, "core//", TrimLeft("//core//", "/")) // only leading trimmed, trailing kept
	AssertEqual(t, "abc", TrimLeft("xyzabc", "xyz"))    // cutset acts as a set
}

func TestString_TrimLeft_Bad(t *T) {
	AssertEqual(t, "", TrimLeft("////", "/"))
	AssertEqual(t, "path/", TrimLeft("path/", "/"))
	AssertEqual(t, "hello", TrimLeft("hello", "")) // empty cutset → unchanged
	AssertEqual(t, "", TrimLeft("", "/"))          // empty input → empty
}

func TestString_TrimLeft_Ugly(t *T) {
	AssertEqual(t, "/path", TrimLeft(" \t /path", " \t"))
	AssertEqual(t, "path  ", TrimLeft("  path  ", " ")) // trailing spaces preserved
	AssertEqual(t, "bc", TrimLeft("aaabc", "a"))        // run of one cut char
}

func TestString_TrimRight_Good(t *T) {
	AssertEqual(t, "path", TrimRight("path///", "/"))
	AssertEqual(t, "hello", TrimRight("hello!!!", "!"))
	AssertEqual(t, "//core", TrimRight("//core//", "/")) // only trailing trimmed, leading kept
	AssertEqual(t, "abc", TrimRight("abcxyz", "xyz"))    // cutset acts as a set
}

func TestString_TrimRight_Bad(t *T) {
	AssertEqual(t, "", TrimRight("////", "/"))
	AssertEqual(t, "/path", TrimRight("/path", "/"))
	AssertEqual(t, "hello", TrimRight("hello", "")) // empty cutset → unchanged
	AssertEqual(t, "", TrimRight("", "/"))          // empty input → empty
}

func TestString_TrimRight_Ugly(t *T) {
	AssertEqual(t, "path", TrimRight("path \t ", " \t"))
	AssertEqual(t, "  path", TrimRight("  path  ", " ")) // leading spaces preserved
	AssertEqual(t, "ab", TrimRight("abccc", "c"))        // run of one cut char
}

func TestString_Index_Good(t *T) {
	AssertEqual(t, 3, Index("key=value", "="))
	AssertEqual(t, 0, Index("hello", "h"))
	AssertEqual(t, 2, Index("hello", "llo"))  // multi-byte substring offset
	AssertEqual(t, 5, Index("café!", "!"))    // byte offset counts multibyte é as 2 bytes
	AssertEqual(t, 0, Index("abcabc", "abc")) // returns first occurrence, not last
}

func TestString_Index_Bad(t *T) {
	AssertEqual(t, -1, Index("nothing here", "?"))
	AssertEqual(t, -1, Index("", "x"))
	AssertEqual(t, -1, Index("go", "golang")) // sep longer than s
	AssertEqual(t, -1, Index("Hello", "h"))   // case sensitive
}

func TestString_Index_Ugly(t *T) {
	AssertEqual(t, 0, Index("hello", ""))
	AssertEqual(t, 0, Index("", ""))    // both empty → 0
	AssertEqual(t, 0, Index("abc", "")) // empty sep matches at the front
}

func TestString_Clone_Good(t *T) {
	AssertEqual(t, "agent", Clone("agent"))
	AssertEqual(t, "café", Clone("café")) // multi-byte preserved
}

func TestString_Clone_Bad(t *T) {
	AssertEqual(t, "", Clone("")) // empty clones to empty
}

func TestString_Clone_Ugly(t *T) {
	// The clone detaches from a large parent's backing array — equal
	// content, independent storage.
	parent := Repeat("x", 4096)
	sub := Clone(parent[:8])
	AssertEqual(t, "xxxxxxxx", sub)
	AssertEqual(t, 8, len(sub))
}

func TestString_LastIndex_Good(t *T) {
	AssertEqual(t, 3, LastIndex("a.b.c", "."))       // last of two dots
	AssertEqual(t, 7, LastIndex("foo/bar/baz", "/")) // last separator
	AssertEqual(t, 0, LastIndex("hello", "h"))
}

func TestString_LastIndex_Bad(t *T) {
	AssertEqual(t, -1, LastIndex("agent", "x"))   // absent
	AssertEqual(t, -1, LastIndex("go", "golang")) // substr longer than s
}

func TestString_LastIndex_Ugly(t *T) {
	AssertEqual(t, 2, LastIndex("aaa", "a"))  // repeated → last position
	AssertEqual(t, 5, LastIndex("agent", "")) // empty substr → len(s)
	AssertEqual(t, 0, LastIndex("", ""))      // both empty → 0
}

func TestString_Builder_Good(t *T) {
	var b Builder
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	AssertEqual(t, "hello world", b.String())
}

func TestString_Builder_Bad(t *T) {
	var b Builder
	AssertEqual(t, "", b.String())
	AssertEqual(t, 0, b.Len())
	n, err := b.WriteString("") // writing nothing is a no-op
	AssertNoError(t, err)
	AssertEqual(t, 0, n)
	AssertEqual(t, "", b.String())
}

func TestString_Builder_Ugly(t *T) {
	type Sink struct {
		out Builder
	}
	s := Sink{}
	s.out.WriteString("ok")
	AssertEqual(t, "ok", s.out.String())
}

func TestString_RuneCount_Good(t *T) {
	AssertEqual(t, 5, RuneCount("hello"))
	AssertEqual(t, 1, RuneCount("🔥"))
	AssertEqual(t, 0, RuneCount(""))
}

func TestString_Concat_Good(t *T) {
	AssertEqual(t, "agent.dispatch.ready", Concat("agent.", "dispatch", ".ready"))
	AssertEqual(t, "https://host/api/v1", Concat("https://", "host", "/api/v1"))
	AssertEqual(t, "abc", Concat("a", "b", "c"))
	AssertEqual(t, "single", Concat("single")) // single part returned verbatim
}

func TestString_Concat_Bad(t *T) {
	AssertEqual(t, "", Concat())
	AssertEqual(t, "", Concat("", "", "")) // all-empty parts → empty
	AssertEqual(t, "", Concat(""))         // single empty part
}

func TestString_Concat_Ugly(t *T) {
	AssertEqual(t, "token=", Concat("token", "=", "")) // trailing empty contributes nothing
	AssertEqual(t, "=value", Concat("", "=", "value")) // leading empty contributes nothing
	AssertEqual(t, "ab", Concat("a", "", "b"))         // interior empty collapses
}

func TestString_Contains_Bad(t *T) {
	AssertFalse(t, Contains("agent dispatch", "homelab"))
	AssertFalse(t, Contains("Hello", "hello")) // case sensitive
	AssertFalse(t, Contains("go", "golang"))   // substr longer than s
	AssertFalse(t, Contains("", "x"))          // empty s, non-empty substr
}

func TestString_Contains_Ugly(t *T) {
	AssertTrue(t, Contains("", ""))         // both empty
	AssertTrue(t, Contains("anything", "")) // empty substr is always contained
	AssertTrue(t, Contains("🔥", ""))        // empty substr contained even in multibyte s
}

func TestString_HasPrefix_Bad(t *T) {
	AssertFalse(t, HasPrefix("agent.dispatch", "task."))
	AssertFalse(t, HasPrefix("Agent", "agent")) // case sensitive
	AssertFalse(t, HasPrefix("go", "golang"))   // prefix longer than s
	AssertFalse(t, HasPrefix("", "x"))          // empty s, non-empty prefix
}

func TestString_HasPrefix_Ugly(t *T) {
	AssertTrue(t, HasPrefix("agent.dispatch", "")) // empty prefix
	AssertTrue(t, HasPrefix("", ""))               // both empty
	AssertTrue(t, HasPrefix("exact", "exact"))     // prefix equals whole string
}

func TestString_HasSuffix_Bad(t *T) {
	AssertFalse(t, HasSuffix("agent.yaml", ".json"))
	AssertFalse(t, HasSuffix("AGENT.YAML", ".yaml")) // case sensitive
	AssertFalse(t, HasSuffix("go", "golang"))        // suffix longer than s
	AssertFalse(t, HasSuffix("", "x"))               // empty s, non-empty suffix
}

func TestString_HasSuffix_Ugly(t *T) {
	AssertTrue(t, HasSuffix("agent.yaml", ""))           // empty suffix
	AssertTrue(t, HasSuffix("", ""))                     // both empty
	AssertTrue(t, HasSuffix(".gitignore", ".gitignore")) // suffix equals whole string
}

func TestString_Join_Bad(t *T) {
	AssertEqual(t, "", Join("/"))
	AssertEqual(t, "", Join(""))      // no parts, empty sep
	AssertEqual(t, "", Join("/", "")) // single empty part → empty
}

func TestString_Join_Ugly(t *T) {
	AssertEqual(t, "agentdispatchready", Join("", "agent", "dispatch", "ready"))
	AssertEqual(t, "a-b", Join("-", "a", "b"))
	AssertEqual(t, "//", Join("/", "", "", "")) // three empty parts → two separators
}

func TestString_Lower_Bad(t *T) {
	AssertEqual(t, "agent-01", Lower("agent-01")) // already lowercase
	AssertEqual(t, "123-456", Lower("123-456"))   // no letters → unchanged
	AssertEqual(t, "snake_case", Lower("snake_case"))
}

func TestString_Lower_Ugly(t *T) {
	AssertEqual(t, "", Lower(""))   // empty input
	AssertEqual(t, "ñ", Lower("Ñ")) // Latin-1 supplement: Ñ → ñ
	AssertEqual(t, "é", Lower("é")) // already-lowercase non-ASCII unchanged
}

func TestString_NewBuilder_Good(t *T) {
	b := NewBuilder()
	n, err := b.WriteString("agent")

	AssertNoError(t, err)
	AssertEqual(t, 5, n)
	AssertEqual(t, "agent", b.String())
}

func TestString_NewBuilder_Bad(t *T) {
	b := NewBuilder()

	AssertEqual(t, "", b.String())
	AssertEqual(t, 0, b.Len())
}

func TestString_NewBuilder_Ugly(t *T) {
	b := NewBuilder()
	_, err := b.WriteString("session")
	RequireNoError(t, err)
	b.Reset()

	AssertEqual(t, "", b.String())
	AssertEqual(t, 0, b.Len())
}

func TestString_NewReader_Good(t *T) {
	reader := NewReader("agent")
	buf := make([]byte, 5)
	n, err := reader.Read(buf)

	AssertNoError(t, err)
	AssertEqual(t, 5, n)
	AssertEqual(t, "agent", string(buf))
}

func TestString_NewReader_Bad(t *T) {
	reader := NewReader("")
	buf := make([]byte, 1)
	n, err := reader.Read(buf)

	AssertEqual(t, 0, n)
	AssertErrorIs(t, err, EOF)
}

func TestString_NewReader_Ugly(t *T) {
	reader := NewReader("agent dispatch")
	offset, err := reader.Seek(6, 0)
	RequireNoError(t, err)
	buf := make([]byte, 8)
	n, err := reader.Read(buf)

	AssertEqual(t, int64(6), offset)
	AssertNoError(t, err)
	AssertEqual(t, 8, n)
	AssertEqual(t, "dispatch", string(buf))
}

func TestString_Replace_Bad(t *T) {
	AssertEqual(t, "agent/dispatch", Replace("agent/dispatch", ".", "/")) // old absent → unchanged
	AssertEqual(t, "hello", Replace("hello", "z", "Z"))                   // old absent → unchanged
	AssertEqual(t, "", Replace("", "a", "b"))                             // empty s → empty
}

func TestString_Replace_Ugly(t *T) {
	AssertEqual(t, ".a.g.e.n.t.", Replace("agent", "", ".")) // empty old inserts around every rune
	AssertEqual(t, "-a-b-", Replace("ab", "", "-"))          // including both ends
	AssertEqual(t, "x", Replace("", "", "x"))                // empty s, empty old → single insert
}

func TestString_RuneCount_Bad(t *T) {
	AssertEqual(t, 14, RuneCount("agent dispatch")) // space counts as a rune
	AssertNotEqual(t, len("agent"), RuneCount("agent dispatch"))
	AssertEqual(t, 4, RuneCount("café"))              // 4 runes...
	AssertNotEqual(t, len("café"), RuneCount("café")) // ...but 5 bytes
}

func TestString_RuneCount_Ugly(t *T) {
	AssertEqual(t, 2, RuneCount(string([]byte{0xff, 'a'})))        // invalid byte → 1 rune, then 'a'
	AssertEqual(t, 3, RuneCount(string([]byte{0xff, 0xfe, 0xfd}))) // each invalid byte counts once
	AssertEqual(t, 0, RuneCount(""))                               // empty → 0
}

func TestString_Split_Bad(t *T) {
	AssertEqual(t, []string{"agent.dispatch"}, Split("agent.dispatch", "/")) // sep absent → single element
	AssertEqual(t, []string{"", ""}, Split("/", "/"))                        // lone sep → two empties
	AssertLen(t, Split("a,b,c,d", ","), 4)
}

func TestString_Split_Ugly(t *T) {
	AssertEqual(t, []string{"a", "b", "c"}, Split("abc", ""))       // empty sep splits per rune
	AssertEqual(t, []string{"c", "a", "f", "é"}, Split("café", "")) // splits on runes, not bytes
	AssertLen(t, Split("", ""), 0)                                  // empty s and empty sep → empty slice
}

func TestString_SplitN_Bad(t *T) {
	AssertNil(t, SplitN("agent=dispatch", "=", 0)) // n=0 → nil
	AssertLen(t, SplitN("agent=dispatch", "=", 0), 0)
	AssertEqual(t, []string{"agent", "dispatch"}, SplitN("agent=dispatch", "=", 5)) // n>parts → all parts
}

func TestString_SplitN_Ugly(t *T) {
	AssertEqual(t, []string{"agent", "dispatch", "ready"}, SplitN("agent=dispatch=ready", "=", -1))
	AssertEqual(t, []string{"a", "b", "c"}, SplitN("abc", "", -1))          // empty sep → per rune
	AssertEqual(t, []string{"noseparator"}, SplitN("noseparator", "=", -1)) // sep absent → single element
}

func TestString_Trim_Bad(t *T) {
	AssertEqual(t, "agent", Trim("agent")) // no whitespace → unchanged
	AssertEqual(t, "a b", Trim("a b"))     // interior whitespace kept
	AssertEqual(t, "", Trim(""))           // empty → empty
}

func TestString_Trim_Ugly(t *T) {
	AssertEqual(t, "agent", Trim("\n\tagent\r\n"))
	AssertEqual(t, "agent", Trim(" agent ")) // NBSP is Unicode whitespace, trimmed
	AssertEqual(t, "", Trim("\t\n\r "))      // only whitespace → empty
}

func TestString_TrimPrefix_Bad(t *T) {
	AssertEqual(t, "agent.dispatch", TrimPrefix("agent.dispatch", "task.")) // prefix absent
	AssertEqual(t, "Agent", TrimPrefix("Agent", "agent"))                   // case sensitive
	AssertEqual(t, "go", TrimPrefix("go", "golang"))                        // prefix longer than s
}

func TestString_TrimPrefix_Ugly(t *T) {
	AssertEqual(t, "agent.dispatch", TrimPrefix("agent.dispatch", "")) // empty prefix → unchanged
	AssertEqual(t, "", TrimPrefix("", ""))                             // both empty
	AssertEqual(t, "", TrimPrefix("", "x"))                            // empty s with prefix → empty
}

func TestString_TrimSuffix_Bad(t *T) {
	AssertEqual(t, "agent.yaml", TrimSuffix("agent.yaml", ".json")) // suffix absent
	AssertEqual(t, "AGENT.YAML", TrimSuffix("AGENT.YAML", ".yaml")) // case sensitive
	AssertEqual(t, "go", TrimSuffix("go", "golang"))                // suffix longer than s
}

func TestString_TrimSuffix_Ugly(t *T) {
	AssertEqual(t, "agent.yaml", TrimSuffix("agent.yaml", "")) // empty suffix → unchanged
	AssertEqual(t, "", TrimSuffix("", ""))                     // both empty
	AssertEqual(t, "", TrimSuffix("", "x"))                    // empty s with suffix → empty
}

func TestString_Upper_Bad(t *T) {
	AssertEqual(t, "AGENT-01", Upper("AGENT-01")) // already uppercase
	AssertEqual(t, "123-456", Upper("123-456"))   // no letters → unchanged
	AssertEqual(t, "SNAKE_CASE", Upper("SNAKE_CASE"))
}

func TestString_Upper_Ugly(t *T) {
	AssertEqual(t, "", Upper(""))   // empty input
	AssertEqual(t, "Ñ", Upper("ñ")) // Latin-1 supplement: ñ → Ñ
	AssertEqual(t, "É", Upper("É")) // already-uppercase non-ASCII unchanged
}

func TestString_HTMLEscape_Good(t *T) {
	AssertEqual(
		t,
		"&lt;p title=&#34;agent &amp; dispatch&#34;&gt;ready&lt;/p&gt;",
		HTMLEscape(`<p title="agent & dispatch">ready</p>`),
	)
}

func TestString_HTMLEscape_Bad(t *T) {
	AssertEqual(t, "", HTMLEscape(""))                             // empty → empty
	AssertEqual(t, "plain text 123", HTMLEscape("plain text 123")) // no special chars → unchanged
	AssertEqual(t, "a &amp; b", HTMLEscape("a & b"))               // ampersand escaped
}

func TestString_HTMLEscape_Ugly(t *T) {
	AssertEqual(t, "&#34;&amp;&#39;&lt;&gt;\x00", HTMLEscape("\"&'<>\x00")) // all five escaped, NUL passthrough
	AssertEqual(t, "&amp;amp;", HTMLEscape("&amp;"))                        // re-escapes an already-escaped entity
	AssertEqual(t, "café &lt;3", HTMLEscape("café <3"))                     // multibyte preserved, < escaped
}

func TestString_HTMLUnescape_Good(t *T) {
	AssertEqual(
		t,
		`<p title="agent & dispatch">ready</p>`,
		HTMLUnescape("&lt;p title=&#34;agent &amp; dispatch&#34;&gt;ready&lt;/p&gt;"),
	)
}

func TestString_HTMLUnescape_Bad(t *T) {
	AssertEqual(t, "", HTMLUnescape(""))                     // empty → empty
	AssertEqual(t, "plain text", HTMLUnescape("plain text")) // no entities → unchanged
	AssertEqual(t, "a & b", HTMLUnescape("a &amp; b"))       // entity decoded
}

func TestString_HTMLUnescape_Ugly(t *T) {
	AssertEqual(t, "agent &unknown; dispatch", HTMLUnescape("agent &unknown; dispatch"))     // unknown entity left as-is
	AssertEqual(t, `<p title="x & y">`, HTMLUnescape("&lt;p title=&#34;x &amp; y&#34;&gt;")) // round-trips escapes
	AssertEqual(t, "A", HTMLUnescape("&#65;"))                                               // numeric entity decoded
}

func TestString_IndexAny_Good(t *T) {
	AssertEqual(t, 1, IndexAny("a/b\\c", "/\\"))
	AssertEqual(t, 0, IndexAny("@host", "@:")) // match at the front
	AssertEqual(t, 4, IndexAny("user@host", "@"))
	AssertEqual(t, 3, IndexAny("café", "é")) // byte offset of a multibyte rune
}

func TestString_IndexAny_Bad(t *T) {
	AssertEqual(t, -1, IndexAny("abc", "/\\"))
	AssertEqual(t, -1, IndexAny("abc", "XYZ")) // none present
	AssertEqual(t, -1, IndexAny("ABC", "abc")) // case sensitive
}

func TestString_IndexAny_Ugly(t *T) {
	AssertEqual(t, -1, IndexAny("", "/"))
	AssertEqual(t, -1, IndexAny("abc", "")) // empty chars never matches
	AssertEqual(t, -1, IndexAny("", ""))    // both empty
}

func TestString_ContainsAny_Good(t *T) {
	AssertTrue(t, ContainsAny("user@host", "@:"))
	AssertTrue(t, ContainsAny("path/to", "/\\")) // one of the set present
	AssertTrue(t, ContainsAny("café", "é"))      // multibyte rune in set
}

func TestString_ContainsAny_Bad(t *T) {
	AssertFalse(t, ContainsAny("userhost", "@:"))
	AssertFalse(t, ContainsAny("ABC", "abc"))          // case sensitive
	AssertFalse(t, ContainsAny("plain", "0123456789")) // no digits present
}

func TestString_ContainsAny_Ugly(t *T) {
	AssertFalse(t, ContainsAny("", "@"))
	AssertFalse(t, ContainsAny("abc", "")) // empty chars never matches
	AssertFalse(t, ContainsAny("", ""))    // both empty
}

func TestString_ContainsRune_Good(t *T) {
	AssertTrue(t, ContainsRune("café", 'é'))
	AssertTrue(t, ContainsRune("hello", 'h')) // first rune
	AssertTrue(t, ContainsRune("hello", 'o')) // last rune
	AssertTrue(t, ContainsRune("a/b", '/'))
}

func TestString_ContainsRune_Bad(t *T) {
	AssertFalse(t, ContainsRune("cafe", 'é'))  // rune absent
	AssertFalse(t, ContainsRune("hello", 'H')) // case sensitive
	AssertFalse(t, ContainsRune("hello", 'z'))
}

func TestString_ContainsRune_Ugly(t *T) {
	AssertFalse(t, ContainsRune("", 'a'))    // empty s
	AssertFalse(t, ContainsRune("abc", 0))   // NUL absent
	AssertTrue(t, ContainsRune("a\x00b", 0)) // NUL present and detected
}

func TestString_Count_Good(t *T) {
	AssertEqual(t, 2, Count("a.b.c", "."))
	AssertEqual(t, 3, Count("aaa", "a"))   // single-char count
	AssertEqual(t, 2, Count("aaaa", "aa")) // non-overlapping matches
	AssertEqual(t, 1, Count("hello", "ll"))
}

func TestString_Count_Bad(t *T) {
	AssertEqual(t, 0, Count("abc", "."))
	AssertEqual(t, 0, Count("", "x"))       // empty s, non-empty substr
	AssertEqual(t, 0, Count("abc", "abcd")) // substr longer than s
}

func TestString_Count_Ugly(t *T) {
	// Empty substr counts code-point boundaries: 1 + RuneCount.
	AssertEqual(t, 4, Count("abc", ""))
	AssertEqual(t, 1, Count("", ""))     // empty s, empty substr → 1
	AssertEqual(t, 5, Count("café", "")) // 1 + 4 runes (boundaries, not bytes)
}

func TestString_EqualFold_Good(t *T) {
	AssertTrue(t, EqualFold("Bearer", "bearer"))
	AssertTrue(t, EqualFold("GoLang", "golang"))
	AssertTrue(t, EqualFold("CAFÉ", "café")) // Unicode case-folding
	AssertTrue(t, EqualFold("ABC", "ABC"))   // identical inputs
}

func TestString_EqualFold_Bad(t *T) {
	AssertFalse(t, EqualFold("Bearer", "Basic"))
	AssertFalse(t, EqualFold("go", "golang")) // differing lengths
	AssertFalse(t, EqualFold("café", "cafe")) // accent differs
}

func TestString_EqualFold_Ugly(t *T) {
	AssertTrue(t, EqualFold("", ""))   // both empty
	AssertFalse(t, EqualFold("", "x")) // one empty
	AssertTrue(t, EqualFold("K", "k")) // Kelvin sign folds to 'k'
}

func TestString_Repeat_Good(t *T) {
	AssertEqual(t, "========", Repeat("=", 8))
	AssertEqual(t, "ababab", Repeat("ab", 3)) // multi-char unit
	AssertEqual(t, "x", Repeat("x", 1))       // count 1
	AssertEqual(t, "abcabc", Repeat("abc", 2))
}

func TestString_Repeat_Bad(t *T) {
	AssertEqual(t, "", Repeat("x", 0))  // count 0 → empty
	AssertEqual(t, "", Repeat("", 100)) // empty unit → empty regardless of count
	AssertEqual(t, "", Repeat("", 0))
}

func TestString_Repeat_Ugly(t *T) {
	// Negative count panics per stdlib contract.
	AssertPanics(t, func() { Repeat("x", -1) })
	AssertPanics(t, func() { Repeat("ab", -5) })  // any negative panics
	AssertNotPanics(t, func() { Repeat("x", 0) }) // zero is safe
}

func TestString_Fields_Good(t *T) {
	AssertEqual(t, []string{"go", "test", "./..."}, Fields("  go   test ./... "))
	AssertEqual(t, []string{"a", "b", "c"}, Fields("a b c"))           // single spaces
	AssertEqual(t, []string{"one"}, Fields("   one   "))               // surrounding whitespace stripped
	AssertEqual(t, []string{"tab", "newline"}, Fields("tab\tnewline")) // any whitespace run splits
}

func TestString_Fields_Bad(t *T) {
	AssertLen(t, Fields("   "), 0)
	AssertLen(t, Fields("\t\n\r "), 0) // mixed whitespace → empty
	AssertEmpty(t, Fields("   "))
}

func TestString_Fields_Ugly(t *T) {
	AssertLen(t, Fields(""), 0)
	AssertEmpty(t, Fields(""))
	AssertEqual(t, []string{"x"}, Fields(" x ")) // NBSP is Unicode whitespace
}

func TestString_Cut_Good(t *T) {
	before, after, found := Cut("port=8080", "=")
	AssertTrue(t, found)
	AssertEqual(t, "port", before)
	AssertEqual(t, "8080", after)
}

func TestString_Cut_Bad(t *T) {
	before, after, found := Cut("noseparator", "=")
	AssertFalse(t, found)
	AssertEqual(t, "noseparator", before)
	AssertEqual(t, "", after)
}

func TestString_Cut_Ugly(t *T) {
	// Empty sep cuts before the first byte.
	before, after, found := Cut("abc", "")
	AssertTrue(t, found)
	AssertEqual(t, "", before)
	AssertEqual(t, "abc", after)
}

func TestString_CutPrefix_Good(t *T) {
	rest, found := CutPrefix("--verbose", "--")
	AssertTrue(t, found)
	AssertEqual(t, "verbose", rest)
}

func TestString_CutPrefix_Bad(t *T) {
	rest, found := CutPrefix("verbose", "--")
	AssertFalse(t, found)
	AssertEqual(t, "verbose", rest)
}

func TestString_CutPrefix_Ugly(t *T) {
	// Empty prefix always cuts, leaving s unchanged.
	rest, found := CutPrefix("abc", "")
	AssertTrue(t, found)
	AssertEqual(t, "abc", rest)
}

func TestString_CutSuffix_Good(t *T) {
	base, found := CutSuffix("main.go", ".go")
	AssertTrue(t, found)
	AssertEqual(t, "main", base)
}

func TestString_CutSuffix_Bad(t *T) {
	base, found := CutSuffix("main.rs", ".go")
	AssertFalse(t, found)
	AssertEqual(t, "main.rs", base)
}

func TestString_CutSuffix_Ugly(t *T) {
	base, found := CutSuffix("abc", "")
	AssertTrue(t, found)
	AssertEqual(t, "abc", base)
}

func TestString_StringReader_Good(t *T) {
	var r *StringReader = NewReader("payload")
	AssertEqual(t, 7, r.Len())
	AssertEqual(t, int64(7), r.Size())
	buf := make([]byte, 3)
	n, err := r.Read(buf)
	AssertNoError(t, err)
	AssertEqual(t, 3, n)
	AssertEqual(t, "pay", string(buf))
	AssertEqual(t, 4, r.Len()) // Len reports the unread remainder
}
