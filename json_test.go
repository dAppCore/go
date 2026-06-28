package core_test

import (
	. "dappco.re/go"
)

type testJSON struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// --- JSONMarshal ---

func TestJson_JSONMarshal_Good(t *T) {
	r := JSONMarshal(testJSON{Name: "brain", Port: 8080})
	AssertTrue(t, r.OK)
	AssertContains(t, string(r.Value.([]byte)), `"name":"brain"`)
}

func TestJson_JSONMarshal_Bad_Unmarshalable(t *T) {
	r := JSONMarshal(make(chan int))
	AssertFalse(t, r.OK)
}

// --- JSONMarshalString ---

func TestJson_JSONMarshalString_Good(t *T) {
	s := JSONMarshalString(testJSON{Name: "x", Port: 1})
	AssertContains(t, s, `"name":"x"`)
}

func TestJson_JSONMarshalString_Ugly_Fallback(t *T) {
	s := JSONMarshalString(make(chan int))
	AssertEqual(t, "{}", s)
}

// --- JSONUnmarshal ---

func TestJson_JSONUnmarshal_Good(t *T) {
	var target testJSON
	r := JSONUnmarshal([]byte(`{"name":"brain","port":8080}`), &target)
	AssertTrue(t, r.OK)
	AssertEqual(t, "brain", target.Name)
	AssertEqual(t, 8080, target.Port)
}

func TestJson_JSONUnmarshal_Bad_Invalid(t *T) {
	var target testJSON
	r := JSONUnmarshal([]byte(`not json`), &target)
	AssertFalse(t, r.OK)
}

// --- JSONUnmarshalString ---

func TestJson_JSONUnmarshalString_Good(t *T) {
	var target testJSON
	r := JSONUnmarshalString(`{"name":"x","port":1}`, &target)
	AssertTrue(t, r.OK)
	AssertEqual(t, "x", target.Name)
}

func TestJson_JSONMarshal_Bad(t *T) {
	r := JSONMarshal(make(chan int))

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestJson_JSONMarshal_Ugly(t *T) {
	r := JSONMarshal(nil)

	AssertTrue(t, r.OK)
	AssertEqual(t, "null", string(r.Value.([]byte)))
}

func TestJson_JSONMarshalIndent_Good(t *T) {
	r := JSONMarshalIndent(testJSON{Name: "codex", Port: 8080}, "", "  ")

	AssertTrue(t, r.OK)
	AssertContains(t, string(r.Value.([]byte)), "\n  \"name\": \"codex\"")
}

func TestJson_JSONMarshalIndent_Bad(t *T) {
	r := JSONMarshalIndent(make(chan int), "", "  ")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestJson_JSONMarshalIndent_Ugly(t *T) {
	r := JSONMarshalIndent(testJSON{Name: "codex", Port: 8080}, ">>", "\t")

	AssertTrue(t, r.OK)
	AssertContains(t, string(r.Value.([]byte)), "\n>>\t\"name\": \"codex\"")
}

func TestJson_JSONMarshalString_Bad(t *T) {
	AssertEqual(t, "{}", JSONMarshalString(make(chan int)))
}

func TestJson_JSONMarshalString_Ugly(t *T) {
	AssertEqual(t, "null", JSONMarshalString(nil))
}

func TestJson_JSONUnmarshal_Bad(t *T) {
	var target testJSON
	r := JSONUnmarshal([]byte(`not json`), &target)

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestJson_JSONUnmarshal_Ugly(t *T) {
	var target testJSON
	r := JSONUnmarshal([]byte(`{"name":"codex","extra":true}`), &target)

	AssertTrue(t, r.OK)
	AssertEqual(t, "codex", target.Name)
	AssertEqual(t, 0, target.Port)
}

func TestJson_JSONUnmarshalString_Bad(t *T) {
	var target testJSON
	r := JSONUnmarshalString(`not json`, &target)

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestJson_JSONUnmarshalString_Ugly(t *T) {
	target := testJSON{Name: "codex", Port: 8080}
	r := JSONUnmarshalString(`{}`, &target)

	AssertTrue(t, r.OK)
	AssertEqual(t, "codex", target.Name)
	AssertEqual(t, 8080, target.Port)
}

func TestJson_RawMessage_Good(t *T) {
	type envelope struct {
		Type string     `json:"type"`
		Data RawMessage `json:"data"`
	}
	var env envelope
	r := JSONUnmarshal([]byte(`{"type":"ping","data":{"port":8080}}`), &env)

	AssertTrue(t, r.OK)
	AssertEqual(t, "ping", env.Type)
	AssertEqual(t, `{"port":8080}`, string(env.Data))
}

func TestJson_RawMessage_Bad(t *T) {
	var raw RawMessage
	r := JSONUnmarshal([]byte(`{"port":8080}`), &raw)

	AssertTrue(t, r.OK)
	AssertEqual(t, `{"port":8080}`, string(raw))
}

func TestJson_RawMessage_Ugly(t *T) {
	raw := RawMessage(`{"a":1}`)
	r := JSONMarshal(raw)

	AssertTrue(t, r.OK)
	AssertEqual(t, `{"a":1}`, string(r.Value.([]byte)))
}

func TestJson_JSONNewEncoder_Good(t *T) {
	b := NewBuilder()
	err := JSONNewEncoder(b).Encode(testJSON{Name: "stream", Port: 80})

	AssertNoError(t, err)
	AssertContains(t, b.String(), `"name":"stream"`)
}

func TestJson_JSONNewEncoder_Bad(t *T) {
	b := NewBuilder()
	err := JSONNewEncoder(b).Encode(make(chan int))

	AssertError(t, err)
}

func TestJson_JSONNewEncoder_Ugly(t *T) {
	// Encode appends a trailing newline and can be called repeatedly to
	// stream multiple values (JSONL).
	b := NewBuilder()
	enc := JSONNewEncoder(b)
	AssertNoError(t, enc.Encode(testJSON{Name: "a", Port: 1}))
	AssertNoError(t, enc.Encode(testJSON{Name: "b", Port: 2}))
	AssertEqual(t, 2, Count(b.String(), "\n"))
}

func TestJson_JSONNewDecoder_Good(t *T) {
	var target testJSON
	err := JSONNewDecoder(NewReader(`{"name":"stream","port":80}`)).Decode(&target)

	AssertNoError(t, err)
	AssertEqual(t, "stream", target.Name)
	AssertEqual(t, 80, target.Port)
}

func TestJson_JSONNewDecoder_Bad(t *T) {
	var target testJSON
	err := JSONNewDecoder(NewReader(`not json`)).Decode(&target)

	AssertError(t, err)
}

func TestJson_JSONNewDecoder_Ugly(t *T) {
	// More()/Decode loop over a concatenated stream of two objects.
	dec := JSONNewDecoder(NewReader(`{"name":"a","port":1}{"name":"b","port":2}`))
	var names []string
	for dec.More() {
		var tj testJSON
		AssertNoError(t, dec.Decode(&tj))
		names = append(names, tj.Name)
	}
	AssertEqual(t, []string{"a", "b"}, names)
}

func TestJson_JSONNumber_Good(t *T) {
	// UseNumber keeps the literal so a large int survives without float
	// rounding.
	dec := JSONNewDecoder(NewReader(`{"big":9007199254740993}`))
	dec.UseNumber()
	var m map[string]JSONNumber
	AssertNoError(t, dec.Decode(&m))
	v := m["big"].Int64
	n, err := v()
	AssertNoError(t, err)
	AssertEqual(t, int64(9007199254740993), n)
}

func TestJson_JSONNumber_Bad(t *T) {
	var n JSONNumber = "not-a-number"
	_, err := n.Int64()
	AssertError(t, err)
}

func TestJson_JSONNumber_Ugly(t *T) {
	var n JSONNumber = "3.14"
	AssertEqual(t, "3.14", n.String())
}

func TestJson_JSONDelim_Good(t *T) {
	dec := JSONNewDecoder(NewReader(`["x"]`))
	tok, err := dec.Token()
	AssertNoError(t, err)
	d, ok := tok.(JSONDelim)
	AssertTrue(t, ok)
	AssertEqual(t, "[", d.String())
}

func TestJson_JSONDelim_Bad(t *T) {
	// A string token is not a JSONDelim.
	dec := JSONNewDecoder(NewReader(`"plain"`))
	tok, err := dec.Token()
	AssertNoError(t, err)
	_, ok := tok.(JSONDelim)
	AssertFalse(t, ok)
}

func TestJson_JSONDelim_Ugly(t *T) {
	AssertEqual(t, "}", JSONDelim('}').String())
}

func TestJson_JSONValid_Good(t *T) {
	AssertTrue(t, JSONValid([]byte(`{"ok":true}`)))
}

func TestJson_JSONValid_Bad(t *T) {
	AssertFalse(t, JSONValid([]byte(`{"ok":`)))
}

func TestJson_JSONValid_Ugly(t *T) {
	// Empty input is not valid JSON.
	AssertFalse(t, JSONValid(nil))
}

func TestJson_JSONUnmarshaler_Good(t *T) {
	// RawMessage satisfies both interfaces and round-trips JSON verbatim.
	var raw RawMessage
	var um JSONUnmarshaler = &raw
	RequireNoError(t, um.UnmarshalJSON([]byte(`{"k":1}`)))
	var m JSONMarshaler = raw
	out, err := m.MarshalJSON()
	AssertNoError(t, err)
	AssertEqual(t, `{"k":1}`, string(out))
}
