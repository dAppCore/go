package core_test

import . "dappco.re/go"

// ExampleJSONMarshal_config initialises configuration values through `JSONMarshal` for
// configuration serialisation. Serialisation and parsing return core Results for
// configuration payloads.
func ExampleJSONMarshal_config() {
	type appConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	r := JSONMarshal(appConfig{Host: "localhost", Port: 8080})
	Println(string(r.Value.([]byte)))
	// Output: {"host":"localhost","port":8080}
}

// ExampleJSONMarshalString serialises a value to JSON text through `JSONMarshalString` for
// configuration serialisation. Serialisation and parsing return core Results for
// configuration payloads.
func ExampleJSONMarshalString() {
	type appConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	Println(JSONMarshalString(appConfig{Host: "localhost", Port: 8080}))
	// Output: {"host":"localhost","port":8080}
}

// ExampleJSONUnmarshal parses JSON bytes through `JSONUnmarshal` for configuration
// serialisation. Serialisation and parsing return core Results for configuration payloads.
func ExampleJSONUnmarshal() {
	type appConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	var cfg appConfig
	JSONUnmarshal([]byte(`{"host":"localhost","port":8080}`), &cfg)
	Println(cfg.Host, cfg.Port)
	// Output: localhost 8080
}

// ExampleJSONUnmarshalString_config initialises configuration values through
// `JSONUnmarshalString` for configuration serialisation. Serialisation and parsing return
// core Results for configuration payloads.
func ExampleJSONUnmarshalString_config() {
	type appConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	var cfg appConfig
	JSONUnmarshalString(`{"host":"localhost","port":8080}`, &cfg)
	Println(cfg.Host, cfg.Port)
	// Output: localhost 8080
}

// ExampleRawMessage defers JSON decoding through the `RawMessage` alias
// for envelope-then-payload parsing. Serialisation and parsing return
// core Results for configuration payloads.
func ExampleRawMessage() {
	type envelope struct {
		Type string     `json:"type"`
		Data RawMessage `json:"data"`
	}
	var env envelope
	JSONUnmarshal([]byte(`{"type":"ping","data":{"port":8080}}`), &env)
	Println(env.Type)
	Println(string(env.Data))
	// Output:
	// ping
	// {"port":8080}
}

// ExampleJSONNewEncoder streams two values straight to stdout as JSONL
// — one JSON object per line — without buffering the whole batch. Encode
// writes the trailing newline itself.
func ExampleJSONNewEncoder() {
	type row struct {
		Name string `json:"name"`
	}
	enc := JSONNewEncoder(Stdout())
	enc.Encode(row{Name: "a"})
	enc.Encode(row{Name: "b"})
	// Output:
	// {"name":"a"}
	// {"name":"b"}
}

// ExampleJSONNewDecoder pulls a single value from a reader, the
// streaming counterpart to JSONUnmarshal.
func ExampleJSONNewDecoder() {
	type cfg struct {
		Port int `json:"port"`
	}
	var c cfg
	JSONNewDecoder(NewReader(`{"port":8080}`)).Decode(&c)
	Println(c.Port)
	// Output: 8080
}

// ExampleJSONValid gates a payload without decoding it.
func ExampleJSONValid() {
	Println(JSONValid([]byte(`{"ok":true}`)))
	Println(JSONValid([]byte(`{"ok":`)))
	// Output:
	// true
	// false
}

// ExampleJSONMarshalIndent renders indented JSON through `JSONMarshalIndent`.
func ExampleJSONMarshalIndent() {
	r := JSONMarshalIndent(map[string]int{"count": 3}, "", "  ")
	Println(string(r.Value.([]byte)))
	// Output:
	// {
	//   "count": 3
	// }
}
