// SPDX-License-Identifier: EUPL-1.2

// JSON helpers for the Core framework.
// Wraps encoding/json so consumers don't import it directly.
// Same guardrail pattern as string.go wraps strings.
//
// Usage:
//
//	data := core.JSONMarshal(myStruct)
//	if data.OK { json := data.Value.([]byte) }
//
//	r := core.JSONUnmarshal(jsonBytes, &target)
//	if !r.OK { /* handle error */ }
package core

import "encoding/json"

// RawMessage is an alias for json.RawMessage — a deferred-decode JSON
// fragment. Lets HTTP / MCP / IPC handlers accept JSON parameters
// without committing to a concrete struct, then decode when the shape
// is known.
//
//	func HandleBridgeCall(args core.RawMessage) core.Result {
//	    var req DeployRequest
//	    if r := core.JSONUnmarshal(args, &req); !r.OK { return r }
//	    return core.Ok(req)
//	}
type RawMessage = json.RawMessage

// JSONNumber is an alias for json.Number — a JSON numeric literal kept
// as its original string so callers choose int vs float decoding (and
// avoid float64 precision loss on large integers). Pairs with the
// UseNumber() option on JSONDecoder.
//
//	var n core.JSONNumber = "9007199254740993"
//	r := n.Int64()  // exact, no float rounding
type JSONNumber = json.Number

// JSONDelim is an alias for json.Delim — one of the four structural
// tokens ( [ ] { } ) returned by JSONDecoder.Token during streaming
// token walks.
//
//	tok, _ := dec.Token()
//	if d, ok := tok.(core.JSONDelim); ok && d == '[' { /* array start */ }
type JSONDelim = json.Delim

// JSONEncoder is an alias for json.Encoder — the streaming writer that
// emits JSON values to an io.Writer one at a time. Construct via
// JSONNewEncoder.
//
//	var enc *core.JSONEncoder = core.JSONNewEncoder(w)
type JSONEncoder = json.Encoder

// JSONDecoder is an alias for json.Decoder — the streaming reader that
// pulls JSON values from an io.Reader. Construct via JSONNewDecoder.
//
//	var dec *core.JSONDecoder = core.JSONNewDecoder(r)
type JSONDecoder = json.Decoder

// JSONMarshaler is an alias for json.Marshaler — the interface a type
// implements to control its own JSON encoding.
//
//	func (id AgentID) MarshalJSON() ([]byte, error) { ... }
//	var _ core.JSONMarshaler = AgentID{}
type JSONMarshaler = json.Marshaler

// JSONUnmarshaler is an alias for json.Unmarshaler — the interface a
// type implements to control its own JSON decoding.
//
//	func (id *AgentID) UnmarshalJSON(b []byte) error { ... }
//	var _ core.JSONUnmarshaler = (*AgentID)(nil)
type JSONUnmarshaler = json.Unmarshaler

// JSONNewEncoder returns a streaming JSON encoder writing to w. Prefer
// it over JSONMarshal when emitting many values to a stream (HTTP
// response, JSONL file) — it writes incrementally without buffering
// every value into one []byte.
//
//	enc := core.JSONNewEncoder(core.Stdout())
//	for _, row := range rows { enc.Encode(row) }  // one JSON value per line
func JSONNewEncoder(w Writer) *JSONEncoder {
	return json.NewEncoder(w)
}

// JSONNewDecoder returns a streaming JSON decoder reading from r. Prefer
// it over JSONUnmarshal when consuming a stream of values or a body of
// unknown length — it decodes one value at a time and can switch to
// JSONNumber mode via dec.UseNumber().
//
//	dec := core.JSONNewDecoder(resp.Body)
//	var msg Message
//	for dec.More() { if err := dec.Decode(&msg); err != nil { break } }
func JSONNewDecoder(r Reader) *JSONDecoder {
	return json.NewDecoder(r)
}

// JSONValid reports whether data is a well-formed JSON encoding. Use it
// to gate a payload before storing or forwarding it without paying a
// full Unmarshal into a throwaway target.
//
//	if !core.JSONValid(body) { return core.NewError("malformed JSON body") }
func JSONValid(data []byte) bool {
	return json.Valid(data)
}

// JSONMarshal serialises a value to JSON bytes.
//
//	r := core.JSONMarshal(myStruct)
//	if r.OK { data := r.Value.([]byte) }
func JSONMarshal(v any) Result {
	data, err := json.Marshal(v)
	if err != nil {
		return Result{err, false}
	}
	return Result{data, true}
}

// JSONMarshalIndent serialises a value to indented JSON bytes.
//
//	r := core.JSONMarshalIndent(report, "", "  ")
//	if r.OK { data := r.Value.([]byte) }
func JSONMarshalIndent(v any, prefix, indent string) Result {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return Result{err, false}
	}
	return Result{data, true}
}

// JSONMarshalString serialises a value to a JSON string.
//
//	s := core.JSONMarshalString(myStruct)
func JSONMarshalString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	// json.Marshal returns a freshly-allocated []byte we own
	// exclusively — AsString skips the copy.
	return AsString(data)
}

// JSONUnmarshal deserialises JSON bytes into a target.
//
//	var cfg Config
//	r := core.JSONUnmarshal(data, &cfg)
func JSONUnmarshal(data []byte, target any) Result {
	if err := json.Unmarshal(data, target); err != nil {
		return Result{err, false}
	}
	return Result{OK: true}
}

// JSONUnmarshalString deserialises a JSON string into a target.
//
//	var cfg Config
//	r := core.JSONUnmarshalString(`{"port":8080}`, &cfg)
//
// Zero-copy: json.Unmarshal treats its input as read-only and does
// not alias the buffer into the unmarshalled values (Strings are
// copied via SetString), so AsBytes is safe here. Saves one alloc
// per call — load-bearing on JSONL hot paths (one call per dataset
// row, thousands per training run).
func JSONUnmarshalString(s string, target any) Result {
	return JSONUnmarshal(AsBytes(s), target)
}
