package core_test

import (
	. "github.com/dappcore/go"
)

// --- HexEncode ---

func TestEncode_HexEncode_Good(t *T) {
	AssertEqual(t, "68656c6c6f", HexEncode([]byte("hello")))
}

func TestEncode_HexEncode_Bad(t *T) {
	AssertEqual(t, "", HexEncode(nil))
	AssertEqual(t, HexEncode(nil), HexEncode([]byte{}))
}

func TestEncode_HexEncode_Ugly(t *T) {
	src := []byte{0x00, 0x0f, 0x10, 0xff}
	encoded := HexEncode(src)

	src[0] = 0xff

	AssertEqual(t, "000f10ff", encoded)
	AssertNotEqual(t, encoded, HexEncode(src))
}

// --- HexDecode ---

func TestEncode_HexDecode_Good(t *T) {
	r := HexDecode("68656c6c6f")
	AssertTrue(t, r.OK)
	AssertEqual(t, []byte("hello"), r.Value)
}

func TestEncode_HexDecode_Bad(t *T) {
	r := HexDecode("not-hex")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

func TestEncode_HexDecode_Ugly(t *T) {
	r := HexDecode("abc")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

// --- Base64Encode ---

func TestEncode_Base64Encode_Good(t *T) {
	AssertEqual(t, "aGVsbG8=", Base64Encode([]byte("hello")))
}

func TestEncode_Base64Encode_Bad(t *T) {
	AssertEqual(t, "", Base64Encode(nil))
	AssertEqual(t, Base64Encode(nil), Base64Encode([]byte{}))
}

func TestEncode_Base64Encode_Ugly(t *T) {
	src := []byte{0xfb, 0xff, 0xff}
	encoded := Base64Encode(src)

	src[0] = 0x00

	AssertEqual(t, "+///", encoded)
	AssertNotEqual(t, encoded, Base64Encode(src))
}

// --- Base64Decode ---

func TestEncode_Base64Decode_Good(t *T) {
	r := Base64Decode("aGVsbG8=")
	AssertTrue(t, r.OK)
	AssertEqual(t, []byte("hello"), r.Value)
}

func TestEncode_Base64Decode_Bad(t *T) {
	r := Base64Decode("not-base64")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

func TestEncode_Base64Decode_Ugly(t *T) {
	r := Base64Decode("aGVsbG8===")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

// --- Base64URLEncode ---

func TestEncode_Base64URLEncode_Good(t *T) {
	AssertEqual(t, "aGVsbG8=", Base64URLEncode([]byte("hello")))
}

func TestEncode_Base64URLEncode_Bad(t *T) {
	AssertEqual(t, "", Base64URLEncode(nil))
	AssertEqual(t, Base64URLEncode(nil), Base64URLEncode([]byte{}))
}

func TestEncode_Base64URLEncode_Ugly(t *T) {
	src := []byte{0xfb, 0xff, 0xff}
	encoded := Base64URLEncode(src)

	src[0] = 0x00

	AssertEqual(t, "-___", encoded)
	AssertNotEqual(t, encoded, Base64URLEncode(src))
}

// --- Base64URLDecode ---

func TestEncode_Base64URLDecode_Good(t *T) {
	r := Base64URLDecode("aGVsbG8=")
	AssertTrue(t, r.OK)
	AssertEqual(t, []byte("hello"), r.Value)
}

func TestEncode_Base64URLDecode_Bad(t *T) {
	r := Base64URLDecode("not+url")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

func TestEncode_Base64URLDecode_Ugly(t *T) {
	r := Base64URLDecode("aGVsbG8===")
	AssertFalse(t, r.OK)
	_, ok := r.Value.(error)
	AssertTrue(t, ok)
}

// --- BigEndianUint64 ---

func TestEncode_BigEndianUint64_Good(t *T) {
	AssertEqual(t, uint64(1), BigEndianUint64([]byte{0, 0, 0, 0, 0, 0, 0, 1}))
}

// Bad: byte order is the whole point — the same bytes read big-endian and
// little-endian must not agree, or a caller has picked the wrong one.
func TestEncode_BigEndianUint64_Bad(t *T) {
	b := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	AssertEqual(t, uint64(1)<<56, BigEndianUint64(b))
	AssertFalse(t, BigEndianUint64(b) == LittleEndianUint64(b))
}

// Ugly: only the first eight bytes are read, so a longer slice is not an error
// and the tail is ignored.
func TestEncode_BigEndianUint64_Ugly(t *T) {
	AssertEqual(t, ^uint64(0), BigEndianUint64([]byte{255, 255, 255, 255, 255, 255, 255, 255, 9, 9}))
}

// --- LittleEndianUint64 ---

func TestEncode_LittleEndianUint64_Good(t *T) {
	AssertEqual(t, uint64(1), LittleEndianUint64([]byte{1, 0, 0, 0, 0, 0, 0, 0}))
}

func TestEncode_LittleEndianUint64_Bad(t *T) {
	AssertEqual(t, uint64(1)<<56, LittleEndianUint64([]byte{0, 0, 0, 0, 0, 0, 0, 1}))
}

// Ugly: round-trips with the write side, which is how Keccak absorbs and
// squeezes its state.
func TestEncode_LittleEndianUint64_Ugly(t *T) {
	buf := make([]byte, 8)
	PutLittleEndianUint64(buf, 0xDEADBEEFCAFEBABE)
	AssertEqual(t, uint64(0xDEADBEEFCAFEBABE), LittleEndianUint64(buf))
}

// --- PutLittleEndianUint64 ---

func TestEncode_PutLittleEndianUint64_Good(t *T) {
	buf := make([]byte, 8)
	PutLittleEndianUint64(buf, 1)
	AssertEqual(t, byte(1), buf[0])
	AssertEqual(t, byte(0), buf[7])
}

// Bad: writes exactly eight bytes and nothing beyond them, so a caller can
// write into the middle of a larger buffer without clobbering its neighbours.
func TestEncode_PutLittleEndianUint64_Bad(t *T) {
	buf := make([]byte, 10)
	buf[8], buf[9] = 7, 7
	PutLittleEndianUint64(buf[:8], ^uint64(0))
	AssertEqual(t, byte(7), buf[8])
	AssertEqual(t, byte(7), buf[9])
}

func TestEncode_PutLittleEndianUint64_Ugly(t *T) {
	buf := make([]byte, 16)
	PutLittleEndianUint64(buf[8:], 0x0102030405060708)
	AssertEqual(t, uint64(0), LittleEndianUint64(buf[:8]))
	AssertEqual(t, uint64(0x0102030405060708), LittleEndianUint64(buf[8:]))
}

// --- HexAppendEncode ---

func TestEncode_HexAppendEncode_Good(t *T) {
	AssertEqual(t, "0aff", AsString(HexAppendEncode(nil, []byte{0x0a, 0xff})))
}

// Bad: an empty src appends nothing at all — dst comes back untouched rather
// than gaining a separator or a padding byte.
func TestEncode_HexAppendEncode_Bad(t *T) {
	AssertEqual(t, "id-", AsString(HexAppendEncode([]byte("id-"), nil)))
}

// Ugly: the reason this exists over HexEncode — it extends dst in place, so a
// caller building an identifier stays at one allocation.
func TestEncode_HexAppendEncode_Ugly(t *T) {
	AssertEqual(t, "id-42-0aff", AsString(HexAppendEncode([]byte("id-42-"), []byte{0x0a, 0xff})))
	AssertEqual(t, HexEncode([]byte{0x0a, 0xff}), AsString(HexAppendEncode(nil, []byte{0x0a, 0xff})))
}

func TestEncode_LittleEndianUint16_Good(t *T) {
	AssertEqual(t, uint16(0x0201), LittleEndianUint16([]byte{1, 2}))
}

func TestEncode_LittleEndianUint16_Bad(t *T) {
	defer func() { AssertNotNil(t, recover(), "short input panics, as the stdlib does") }()
	LittleEndianUint16([]byte{1})
}

func TestEncode_LittleEndianUint16_Ugly(t *T) {
	// Extra bytes are ignored.
	AssertEqual(t, uint16(0xffff), LittleEndianUint16([]byte{0xff, 0xff, 9, 9}))
}

func TestEncode_LittleEndianUint32_Good(t *T) {
	AssertEqual(t, uint32(0x04030201), LittleEndianUint32([]byte{1, 2, 3, 4}))
}

func TestEncode_LittleEndianUint32_Bad(t *T) {
	defer func() { AssertNotNil(t, recover(), "short input panics, as the stdlib does") }()
	LittleEndianUint32([]byte{1, 2, 3})
}

func TestEncode_LittleEndianUint32_Ugly(t *T) {
	AssertEqual(t, uint32(0), LittleEndianUint32(make([]byte, 4)))
}

func TestEncode_BigEndianUint16_Good(t *T) {
	AssertEqual(t, uint16(0x0102), BigEndianUint16([]byte{1, 2}))
}

func TestEncode_BigEndianUint16_Bad(t *T) {
	defer func() { AssertNotNil(t, recover(), "short input panics, as the stdlib does") }()
	BigEndianUint16(nil)
}

func TestEncode_BigEndianUint16_Ugly(t *T) {
	AssertEqual(t, uint16(0xff00), BigEndianUint16([]byte{0xff, 0, 7}))
}

func TestEncode_AppendLittleEndianUint16_Good(t *T) {
	AssertEqual(t, []byte{9, 0x01, 0x02}, AppendLittleEndianUint16([]byte{9}, 0x0201))
}

func TestEncode_AppendLittleEndianUint16_Bad(t *T) {
	// A nil slice is a fine starting point.
	AssertEqual(t, []byte{0, 0}, AppendLittleEndianUint16(nil, 0))
}

func TestEncode_AppendLittleEndianUint16_Ugly(t *T) {
	AssertEqual(t, uint16(0xbeef), LittleEndianUint16(AppendLittleEndianUint16(nil, 0xbeef)))
}

func TestEncode_AppendLittleEndianUint32_Good(t *T) {
	AssertEqual(t, []byte{1, 2, 3, 4}, AppendLittleEndianUint32(nil, 0x04030201))
}

func TestEncode_AppendLittleEndianUint32_Bad(t *T) {
	AssertEqual(t, []byte{0, 0, 0, 0}, AppendLittleEndianUint32(nil, 0))
}

func TestEncode_AppendLittleEndianUint32_Ugly(t *T) {
	AssertEqual(t, uint32(0xdeadbeef), LittleEndianUint32(AppendLittleEndianUint32(nil, 0xdeadbeef)))
}
