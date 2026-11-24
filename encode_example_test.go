package core_test

import . "github.com/dappcore/go"

// ExampleHexEncode encodes bytes as hex through `HexEncode` for token and payload
// encoding. Byte and string encoders return predictable wrapper values for tokens and
// payloads.
func ExampleHexEncode() {
	Println(HexEncode([]byte("hello")))
	// Output: 68656c6c6f
}

// ExampleHexDecode decodes hex text through `HexDecode` for token and payload encoding.
// Byte and string encoders return predictable wrapper values for tokens and payloads.
func ExampleHexDecode() {
	r := HexDecode("68656c6c6f")
	Println(string(r.Value.([]byte)))
	// Output: hello
}

// ExampleBase64Encode encodes bytes as base64 through `Base64Encode` for token and payload
// encoding. Byte and string encoders return predictable wrapper values for tokens and
// payloads.
func ExampleBase64Encode() {
	Println(Base64Encode([]byte("hello")))
	// Output: aGVsbG8=
}

// ExampleBase64Decode decodes base64 text through `Base64Decode` for token and payload
// encoding. Byte and string encoders return predictable wrapper values for tokens and
// payloads.
func ExampleBase64Decode() {
	r := Base64Decode("aGVsbG8=")
	Println(string(r.Value.([]byte)))
	// Output: hello
}

// ExampleBase64URLEncode encodes bytes for URL-safe base64 through `Base64URLEncode` for
// token and payload encoding. Byte and string encoders return predictable wrapper values
// for tokens and payloads.
func ExampleBase64URLEncode() {
	Println(Base64URLEncode([]byte("hello?")))
	// Output: aGVsbG8_
}

// ExampleBase64URLDecode decodes URL-safe base64 text through `Base64URLDecode` for token
// and payload encoding. Byte and string encoders return predictable wrapper values for
// tokens and payloads.
func ExampleBase64URLDecode() {
	r := Base64URLDecode("aGVsbG8_")
	Println(string(r.Value.([]byte)))
	// Output: hello?
}

// ExampleLittleEndianUint16 reads a packet's opcode.
func ExampleLittleEndianUint16() {
	Println(LittleEndianUint16([]byte{0x34, 0x12}) == 0x1234)
	// Output: true
}

// ExampleLittleEndianUint32 reads a length field.
func ExampleLittleEndianUint32() {
	Println(LittleEndianUint32([]byte{10, 0, 0, 0}))
	// Output: 10
}

// ExampleBigEndianUint16 reads a port in network order.
func ExampleBigEndianUint16() {
	Println(BigEndianUint16([]byte{0x1f, 0x90}))
	// Output: 8080
}

// ExampleAppendLittleEndianUint16 builds a packet header.
func ExampleAppendLittleEndianUint16() {
	Println(AppendLittleEndianUint16(nil, 0x1234))
	// Output: [52 18]
}

// ExampleAppendLittleEndianUint32 builds a length field.
func ExampleAppendLittleEndianUint32() {
	Println(AppendLittleEndianUint32(nil, 10))
	// Output: [10 0 0 0]
}
