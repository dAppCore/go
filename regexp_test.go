// SPDX-License-Identifier: EUPL-1.2

package core_test

import (
	. "dappco.re/go"
)

func TestRegexp_Regex_Good_Compile(t *T) {
	r := Regex(`\d+`)
	AssertTrue(t, r.OK)
	AssertNotNil(t, r.Value)
}

func TestRegexp_Regex_Bad_InvalidPattern(t *T) {
	r := Regex(`[`)
	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestRegexp_MatchString_Good(t *T) {
	rx := Regex(`\d+`).Value.(*Regexp)
	AssertTrue(t, rx.MatchString("foo123bar"))
	AssertFalse(t, rx.MatchString("no digits"))
}

func TestRegexp_FindString_Good(t *T) {
	rx := Regex(`\d+`).Value.(*Regexp)
	AssertEqual(t, "42", rx.FindString("hello 42 world"))
	AssertEqual(t, "", rx.FindString("nothing here"))
}

func TestRegexp_FindAllString_Good(t *T) {
	rx := Regex(`\d+`).Value.(*Regexp)
	AssertEqual(t, []string{"1", "2", "3"}, rx.FindAllString("a1 b2 c3", -1))
	AssertEqual(t, []string{"1", "2"}, rx.FindAllString("a1 b2 c3", 2))
}

func TestRegexp_FindStringSubmatch_Good(t *T) {
	rx := Regex(`(\w+)=(\d+)`).Value.(*Regexp)
	got := rx.FindStringSubmatch("count=42")
	AssertEqual(t, []string{"count=42", "count", "42"}, got)
}

func TestRegexp_ReplaceAllString_Good(t *T) {
	rx := Regex(`\d`).Value.(*Regexp)
	AssertEqual(t, "aXbX", rx.ReplaceAllString("a1b2", "X"))
	// Consecutive matches each get replaced.
	AssertEqual(t, "aXXb", rx.ReplaceAllString("a12b", "X"))
	// No match leaves the input untouched.
	AssertEqual(t, "abc", rx.ReplaceAllString("abc", "X"))
	// Replacing with the empty string deletes every match.
	AssertEqual(t, "", rx.ReplaceAllString("123", ""))
	// Empty input yields empty output.
	AssertEqual(t, "", rx.ReplaceAllString("", "X"))
}

func TestRegexp_Split_Good(t *T) {
	rx := Regex(`,+`).Value.(*Regexp)
	// Runs of commas collapse into a single separator.
	AssertEqual(t, []string{"a", "b", "c"}, rx.Split("a,b,,c", -1))
	// n>0 caps the number of substrings; the remainder is left intact.
	AssertEqual(t, []string{"a", "b,,c"}, rx.Split("a,b,,c", 2))
	// No separator present returns the whole string as one element.
	AssertEqual(t, []string{"abc"}, rx.Split("abc", -1))
}

func TestRegexp_String_Good(t *T) {
	rx := Regex(`abc`).Value.(*Regexp)
	AssertEqual(t, "abc", rx.String())
	// String returns the exact source pattern, metacharacters included.
	AssertEqual(t, `\d+`, Regex(`\d+`).Value.(*Regexp).String())
	AssertEqual(t, `^(a|b)$`, Regex(`^(a|b)$`).Value.(*Regexp).String())
}

func TestRegexp_Regex_Good(t *T) {
	r := Regex(`agent-[0-9]+`)

	AssertTrue(t, r.OK)
	AssertEqual(t, `agent-[0-9]+`, r.Value.(*Regexp).String())
}

func TestRegexp_Regex_Bad(t *T) {
	r := Regex(`[agent`)

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestRegexp_Regex_Ugly(t *T) {
	r := Regex("")

	AssertTrue(t, r.OK)
	AssertTrue(t, r.Value.(*Regexp).MatchString("agent"))
}

func TestRegexp_Regexp_FindAllString_Good(t *T) {
	rx := Regex(`agent-[0-9]+`).Value.(*Regexp)

	all := rx.FindAllString("agent-1 agent-2", -1)
	AssertEqual(t, []string{"agent-1", "agent-2"}, all)
	AssertLen(t, all, 2)
	// n>0 caps the number of matches returned, leftmost first.
	AssertEqual(t, []string{"agent-1"}, rx.FindAllString("agent-1 agent-2", 1))
	// A single match still comes back as a one-element slice.
	AssertEqual(t, []string{"agent-99"}, rx.FindAllString("agent-99", -1))
}

func TestRegexp_Regexp_FindAllString_Bad(t *T) {
	rx := Regex(`agent-[0-9]+`).Value.(*Regexp)

	// No match anywhere returns nil, not an empty slice.
	AssertNil(t, rx.FindAllString("no agents", -1))
	AssertEmpty(t, rx.FindAllString("no agents", -1))
	// A partial match (prefix present, digits missing) is still no match.
	AssertNil(t, rx.FindAllString("agent-", -1))
	// Empty input never matches.
	AssertNil(t, rx.FindAllString("", -1))
}

func TestRegexp_Regexp_FindAllString_Ugly(t *T) {
	rx := Regex(`agent-[0-9]+`).Value.(*Regexp)

	// n==0 requests zero matches, so nil comes back even when matches exist.
	AssertNil(t, rx.FindAllString("agent-1 agent-2", 0))
	// Contrast: n==1 on the same input does return a match, proving the
	// nil above is driven by n, not by the absence of matches.
	AssertEqual(t, []string{"agent-1"}, rx.FindAllString("agent-1 agent-2", 1))
}

func TestRegexp_Regexp_FindString_Good(t *T) {
	rx := Regex(`agent-[0-9]+`).Value.(*Regexp)

	AssertEqual(t, "agent-7", rx.FindString("dispatch agent-7 ready"))
	// FindString returns the leftmost match when several are present.
	AssertEqual(t, "agent-1", rx.FindString("agent-1 agent-2"))
	// A match spanning the whole input is returned verbatim.
	AssertEqual(t, "agent-42", rx.FindString("agent-42"))
}

func TestRegexp_Regexp_FindString_Bad(t *T) {
	rx := Regex(`agent-[0-9]+`).Value.(*Regexp)

	AssertEqual(t, "", rx.FindString("dispatch ready"))
	// A partial match (no trailing digits) is not a match.
	AssertEqual(t, "", rx.FindString("agent-"))
	// Empty input yields the empty no-match sentinel.
	AssertEqual(t, "", rx.FindString(""))
	AssertEmpty(t, rx.FindString("dispatch ready"))
}

func TestRegexp_Regexp_FindString_Ugly(t *T) {
	rx := Regex(``).Value.(*Regexp)

	// An empty pattern matches the empty string at position 0, so every
	// FindString returns "" — indistinguishable from a no-match string-wise.
	AssertEqual(t, "", rx.FindString("dispatch"))
	AssertEqual(t, "", rx.FindString(""))
	// The empty pattern genuinely matches (it is not a failed match).
	AssertTrue(t, rx.MatchString("dispatch"))
}

func TestRegexp_Regexp_FindStringSubmatch_Good(t *T) {
	rx := Regex(`agent=([a-z]+)`).Value.(*Regexp)

	got := rx.FindStringSubmatch("agent=codex")
	AssertEqual(t, []string{"agent=codex", "codex"}, got)
	// Index 0 is the full match, index 1 the first capture group.
	AssertLen(t, got, 2)
	// The full match is located within surrounding text.
	AssertEqual(t, []string{"agent=claude", "claude"}, rx.FindStringSubmatch("run agent=claude now"))
}

func TestRegexp_Regexp_FindStringSubmatch_Bad(t *T) {
	rx := Regex(`agent=([a-z]+)`).Value.(*Regexp)

	// Digits do not satisfy [a-z], so there is no match at all.
	AssertNil(t, rx.FindStringSubmatch("agent=42"))
	// Group can't match empty (+ requires one char), so "agent=" fails.
	AssertNil(t, rx.FindStringSubmatch("agent="))
	// Empty input never matches.
	AssertNil(t, rx.FindStringSubmatch(""))
}

func TestRegexp_Regexp_FindStringSubmatch_Ugly(t *T) {
	rx := Regex(`(agent)?`).Value.(*Regexp)

	// The optional group does not participate, but the regex still matches
	// the empty string at position 0: full match "" and empty group "".
	AssertEqual(t, []string{"", ""}, rx.FindStringSubmatch(""))
	// Non-matching text still matches empty at the start, group absent.
	AssertEqual(t, []string{"", ""}, rx.FindStringSubmatch("xyz"))
	// When the optional group is present it is captured greedily.
	AssertEqual(t, []string{"agent", "agent"}, rx.FindStringSubmatch("agent"))
}

func TestRegexp_Regexp_MatchString_Good(t *T) {
	rx := Regex(`^agent\.`).Value.(*Regexp)

	AssertTrue(t, rx.MatchString("agent.dispatch"))
	// Any suffix after the anchored prefix still matches.
	AssertTrue(t, rx.MatchString("agent.task"))
	// The ^ anchor means the prefix must be at the very start.
	AssertFalse(t, rx.MatchString("x agent.dispatch"))
}

func TestRegexp_Regexp_MatchString_Bad(t *T) {
	rx := Regex(`^agent\.`).Value.(*Regexp)

	AssertFalse(t, rx.MatchString("task.dispatch"))
	// Empty input cannot satisfy the prefix.
	AssertFalse(t, rx.MatchString(""))
	// The escaped dot is literal: a different char after "agent" fails.
	AssertFalse(t, rx.MatchString("agentXdispatch"))
}

func TestRegexp_Regexp_MatchString_Ugly(t *T) {
	rx := Regex(`^$`).Value.(*Regexp)

	// ^$ matches only the fully empty string.
	AssertTrue(t, rx.MatchString(""))
	AssertFalse(t, rx.MatchString("x"))
	// A single space is non-empty, so it does not match.
	AssertFalse(t, rx.MatchString(" "))
}

func TestRegexp_Regexp_ReplaceAllString_Good(t *T) {
	rx := Regex(`/+`).Value.(*Regexp)

	// A run of slashes collapses to a single replacement.
	AssertEqual(t, "agent.dispatch.ready", rx.ReplaceAllString("agent/dispatch//ready", "."))
	// A lone slash is replaced too.
	AssertEqual(t, "a.b", rx.ReplaceAllString("a/b", "."))
	// No slash present leaves the input unchanged.
	AssertEqual(t, "agent", rx.ReplaceAllString("agent", "."))
}

func TestRegexp_Regexp_ReplaceAllString_Bad(t *T) {
	rx := Regex(`missing`).Value.(*Regexp)

	// No match: the input is returned verbatim.
	AssertEqual(t, "agent.dispatch", rx.ReplaceAllString("agent.dispatch", "task"))
	// Empty input yields empty output regardless of replacement.
	AssertEqual(t, "", rx.ReplaceAllString("", "task"))
	// Unrelated input is likewise untouched.
	AssertEqual(t, "nothing here", rx.ReplaceAllString("nothing here", "task"))
}

func TestRegexp_Regexp_ReplaceAllString_Ugly(t *T) {
	rx := Regex(`agent-([0-9]+)`).Value.(*Regexp)

	// $1 expands to the first capture group.
	AssertEqual(t, "id=42", rx.ReplaceAllString("agent-42", "id=$1"))
	// ${1} brace syntax references the same group.
	AssertEqual(t, "7", rx.ReplaceAllString("agent-7", "${1}"))
	// Every match is expanded independently.
	AssertEqual(t, "[1] [2]", rx.ReplaceAllString("agent-1 agent-2", "[$1]"))
}

func TestRegexp_Regexp_Split_Good(t *T) {
	rx := Regex(`\s+`).Value.(*Regexp)

	AssertEqual(t, []string{"agent", "dispatch", "ready"}, rx.Split("agent dispatch ready", -1))
	// A run of whitespace counts as one separator.
	AssertEqual(t, []string{"agent", "dispatch"}, rx.Split("agent   dispatch", -1))
	// n>0 caps the substrings, leaving the tail unsplit.
	AssertEqual(t, []string{"a", "b c"}, rx.Split("a b c", 2))
	// No separator returns the whole string as a single element.
	AssertEqual(t, []string{"agent"}, rx.Split("agent", -1))
}

func TestRegexp_Regexp_Split_Bad(t *T) {
	rx := Regex(`,`).Value.(*Regexp)

	// n==0 requests zero substrings, so nil comes back.
	AssertNil(t, rx.Split("agent,dispatch", 0))
	AssertNil(t, rx.Split("", 0))
	// Contrast: n<0 on the same input splits normally, proving n drove the nil.
	AssertEqual(t, []string{"agent", "dispatch"}, rx.Split("agent,dispatch", -1))
}

func TestRegexp_Regexp_Split_Ugly(t *T) {
	rx := Regex(`,`).Value.(*Regexp)

	// Leading and trailing separators produce empty leading/trailing elements.
	AssertEqual(t, []string{"", "agent", ""}, rx.Split(",agent,", -1))
	// Consecutive separators yield an empty element between them.
	AssertEqual(t, []string{"a", "", "b"}, rx.Split("a,,b", -1))
	// A string that is nothing but the separator splits into two empties.
	AssertEqual(t, []string{"", ""}, rx.Split(",", -1))
}

func TestRegexp_Regexp_String_Good(t *T) {
	rx := Regex(`agent-[0-9]+`).Value.(*Regexp)

	AssertEqual(t, `agent-[0-9]+`, rx.String())
	// String round-trips the source pattern unchanged, anchors included.
	AssertEqual(t, `^(a|b)+$`, Regex(`^(a|b)+$`).Value.(*Regexp).String())
}

func TestRegexp_Regexp_String_Bad(t *T) {
	rx := Regex(`\.`).Value.(*Regexp)

	// Escape sequences are preserved verbatim in the source.
	AssertEqual(t, `\.`, rx.String())
	AssertEqual(t, `a\sb`, Regex(`a\sb`).Value.(*Regexp).String())
	AssertEqual(t, `\d{3}`, Regex(`\d{3}`).Value.(*Regexp).String())
}

func TestRegexp_Regexp_String_Ugly(t *T) {
	rx := Regex("").Value.(*Regexp)

	// The empty pattern reports an empty source string...
	AssertEqual(t, "", rx.String())
	// ...yet is a valid regex that matches the empty string.
	AssertTrue(t, rx.MatchString(""))
	// A whitespace-only pattern is preserved, not trimmed.
	AssertEqual(t, " ", Regex(" ").Value.(*Regexp).String())
}
