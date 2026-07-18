package core_test

import . "dappco.re/go"

// Example_sha3_256 hashes a byte payload with SHA3-256 for Ethereum-compatible digest
// work. SHA3 and Keccak helpers cover Ethereum-compatible digest inputs without direct
// crypto imports.
//
// NOTE: this cannot be named ExampleSHA3_256 — go vet's example-association parser
// splits an Example name on its FIRST underscore and treats the head as the target
// identifier, so ExampleSHA3_256 resolves to "SHA3" (unknown) and go vet — and `go
// test`'s automatic pre-test vet subset — hard-fails the whole package build with
// "refers to unknown identifier: SHA3". Confirmed with a standalone repro. Any exported
// identifier containing an underscore (SHA3_256, SHA3_256Hex) is unnameable as a
// symbol-attached Example under Go's own naming grammar; the package-level
// leading-underscore form used here is the only mechanism that stays vet-clean.
func Example_sha3_256() {
	sum := SHA3_256([]byte("hello"))
	Println(HexEncode(sum[:])[:16])
	// Output: 3338be694f50c5f3
}

// Example_sha3_256Hex hashes a byte payload and renders SHA3-256 hex for
// Ethereum-compatible digest work. SHA3 and Keccak helpers cover Ethereum-compatible
// digest inputs without direct crypto imports.
//
// NOTE: see Example_sha3_256 — ExampleSHA3_256Hex has the same go-vet-breaking shape.
func Example_sha3_256Hex() {
	Println(SHA3_256Hex([]byte("hello"))[:16])
	// Output: 3338be694f50c5f3
}

// ExampleKeccak256 hashes bytes with Keccak through `Keccak256` for Ethereum-compatible
// digest work. SHA3 and Keccak helpers cover Ethereum-compatible digest inputs without
// direct crypto imports.
func ExampleKeccak256() {
	sum := Keccak256([]byte("hello"))
	Println(HexEncode(sum[:])[:16])
	// Output: 1c8aff950685c2ed
}

// ExampleKeccak256Hex hashes bytes with Keccak and hex output through `Keccak256Hex` for
// Ethereum-compatible digest work. SHA3 and Keccak helpers cover Ethereum-compatible
// digest inputs without direct crypto imports.
func ExampleKeccak256Hex() {
	Println(Keccak256Hex([]byte("hello"))[:16])
	// Output: 1c8aff950685c2ed
}

// ExampleSHA3Shake128 generates SHAKE128 output through `SHA3Shake128` for
// Ethereum-compatible digest work. SHA3 and Keccak helpers cover Ethereum-compatible
// digest inputs without direct crypto imports.
func ExampleSHA3Shake128() {
	sum := SHA3Shake128([]byte("hello"), 8)
	Println(len(sum))
	Println(HexEncode(sum))
	// Output:
	// 8
	// 8eb4b6a932f28033
}

// ExampleSHA3Shake256 generates SHAKE256 output through `SHA3Shake256` for
// Ethereum-compatible digest work. SHA3 and Keccak helpers cover Ethereum-compatible
// digest inputs without direct crypto imports.
func ExampleSHA3Shake256() {
	sum := SHA3Shake256([]byte("hello"), 8)
	Println(len(sum))
	Println(HexEncode(sum))
	// Output:
	// 8
	// 1234075ae4a1e773
}
