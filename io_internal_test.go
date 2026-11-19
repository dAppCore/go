// SPDX-License-Identifier: EUPL-1.2

package core

// failingReader hands out its bytes, then fails with err.
type failingReader struct {
	data []byte
	err  error
}

func (r *failingReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, r.err
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

func TestIo_readAllSized_Good(t *T) {
	// Exactly n bytes: one sized buffer, no growth.
	r := readAllSized(NewReader("homelab"), 7)

	AssertTrue(t, r.OK)
	AssertEqual(t, "homelab", string(r.Value.([]byte)))
}

func TestIo_readAllSized_Bad(t *T) {
	// A read error is a failed Result carrying the io.read.failed code.
	r := readAllSized(&failingReader{data: []byte("ho"), err: NewError("disk gone")}, 7)

	AssertFalse(t, r.OK)
	AssertEqual(t, "io.read.failed", r.Code())
}

func TestIo_readAllSized_Ugly(t *T) {
	// The size hint was too big: EOF early returns what was read, trimmed.
	short := readAllSized(NewReader("ho"), 7)
	AssertTrue(t, short.OK)
	AssertEqual(t, "ho", string(short.Value.([]byte)))

	// n <= 0 reads nothing and yields a nil slice.
	none := readAllSized(NewReader("homelab"), 0)
	AssertTrue(t, none.OK)
	AssertNil(t, none.Value.([]byte))
}
