// SPDX-License-Identifier: EUPL-1.2

// Base I/O primitives for the Core framework.
//
// Re-exports stdlib io interfaces and the EOF sentinel as core types so
// consumer packages declare Reader/Writer parameters via core without
// importing io directly. Bundle size is unaffected by aliases — they
// resolve to the same underlying types — but the AX-6 import-line cost
// drops to zero for typed-parameter use.
//
// Usage:
//
//	func handle(r core.Reader, w core.Writer) error { ... }
//
//	if errors.Is(err, core.EOF) { ... }
//
//	n := core.Copy(dst, src)
//	if !n.OK { return n }
package core

import (
	"bytes"
	"io"
)

// Reader is the canonical io.Reader interface, exported as core.Reader.
//
//	var reader core.Reader = core.NewReader("agent payload")
//	r := core.Copy(core.Stdout(), reader)
//	if !r.OK { return r }
type Reader = io.Reader

// Writer is the canonical io.Writer interface, exported as core.Writer.
//
//	var writer core.Writer = core.Stdout()
//	r := core.WriteString(writer, "agent ready\n")
//	if !r.OK { return r }
type Writer = io.Writer

// Closer is the canonical io.Closer interface.
//
//	r := core.HTTPGet("https://api.lethean.example/health")
//	if !r.OK { return r }
//	var closer core.Closer = r.Value.(*core.Response).Body
//	defer closer.Close()
type Closer = io.Closer

// ReadCloser composes Reader and Closer.
//
//	r := core.HTTPGet("https://api.lethean.example/health")
//	if !r.OK { return r }
//	body := r.Value.(*core.Response).Body
//	var reader core.ReadCloser = body
//	defer reader.Close()
type ReadCloser = io.ReadCloser

// WriteCloser composes Writer and Closer.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Create("logs/agent.log")
//	if !r.OK { return r }
//	writer := r.Value.(core.WriteCloser)
//	defer writer.Close()
type WriteCloser = io.WriteCloser

// ReadWriter composes Reader and Writer.
//
//	buf := core.NewBufferString("agent payload")
//	var rw core.ReadWriter = buf
//	core.WriteString(rw, " acknowledged")
type ReadWriter = io.ReadWriter

// ReadWriteCloser composes Reader, Writer, and Closer.
//
//	a, b := core.NetPipe()
//	defer a.Close()
//	defer b.Close()
//	var rwc core.ReadWriteCloser = a
//	_ = rwc
type ReadWriteCloser = io.ReadWriteCloser

// Seeker is the canonical io.Seeker interface.
//
//	reader := core.NewReader("agent payload")
//	var seeker core.Seeker = reader
//	seeker.Seek(0, 0)
type Seeker = io.Seeker

// EOF is the canonical end-of-stream sentinel error.
//
//	if errors.Is(err, core.EOF) { /* end of stream */ }
var EOF = io.EOF

// Discard is a Writer on which all Write calls succeed without doing anything.
//
//	table := core.NewTable(core.Discard)
var Discard Writer = io.Discard

// Copy copies from src to dst until EOF on src or an error occurs.
// Returns Result wrapping the number of bytes copied (int64).
//
//	r := core.Copy(dst, src)
//	if !r.OK { return r }
//	n := r.Value.(int64)
func Copy(dst Writer, src Reader) Result {
	n, err := io.Copy(dst, src)
	if err != nil {
		return Result{err, false}
	}
	return Result{n, true}
}

// CopyN copies n bytes (or until an error) from src to dst.
//
//	r := core.CopyN(dst, src, 1024)
func CopyN(dst Writer, src Reader, n int64) Result {
	written, err := io.CopyN(dst, src, n)
	if err != nil {
		return Result{err, false}
	}
	return Result{written, true}
}

// WriteString writes the contents of s to w. Returns Result wrapping
// the number of bytes written (int).
//
//	r := core.WriteString(stdout, "hello\n")
func WriteString(w Writer, s string) Result {
	n, err := io.WriteString(w, s)
	if err != nil {
		return Result{err, false}
	}
	return Result{n, true}
}

// ReadAll reads all bytes from reader and closes it when it implements Closer.
//
//	r := core.ReadAll(core.NewReader("hello"))
//	if r.OK { core.Println(r.Value.(string)) }
func ReadAll(reader any) Result {
	rc, ok := reader.(Reader)
	if !ok {
		return Result{E("core.ReadAll", "not a reader", nil), false}
	}
	var data []byte
	var err error
	// Fast path: if the reader knows its remaining size (bytes.Reader,
	// bytes.Buffer, strings.Reader, *io.LimitedReader-of-known all
	// expose Len()), allocate the destination once at the exact size
	// instead of paying io.ReadAll's 5-10 buffer doublings (which cost
	// roughly 3x the final byte count in transient allocations).
	if sizer, hasLen := reader.(interface{ Len() int }); hasLen {
		data, err = readAllSized(rc, sizer.Len())
	} else {
		data, err = io.ReadAll(rc)
	}
	if closer, ok := reader.(Closer); ok {
		closer.Close()
	}
	if err != nil {
		return Result{err, false}
	}
	return Result{string(data), true}
}

// readAllSized reads exactly n bytes (or until EOF) into a pre-allocated
// buffer. Used by ReadAll when the source's remaining length is known.
func readAllSized(r Reader, n int) ([]byte, error) {
	if n <= 0 {
		return nil, nil
	}
	buf := make([]byte, n)
	read := 0
	for read < n {
		m, err := r.Read(buf[read:])
		read += m
		if err != nil {
			if err == io.EOF {
				return buf[:read], nil
			}
			return buf[:read], err
		}
	}
	return buf, nil
}

// Buffer is an alias for bytes.Buffer — an in-memory byte sequence with
// io.Reader/io.Writer methods. Lets consumers declare buffer-typed
// fields and locals without importing bytes.
//
//	type Sink struct {
//	    out core.Buffer
//	}
//
//	var buf core.Buffer
//	buf.WriteString("ready")
type Buffer = bytes.Buffer

// NewBuffer returns a bytes.Buffer initialised with b.
// With no input, it returns an empty bytes.Buffer.
//
//	buf := core.NewBuffer([]byte("hello"))
//	empty := core.NewBuffer()
func NewBuffer(b ...[]byte) *bytes.Buffer {
	if len(b) == 0 {
		return &bytes.Buffer{}
	}
	return bytes.NewBuffer(b[0])
}

// NewBufferString returns a bytes.Buffer initialised with s.
//
//	buf := core.NewBufferString("hello")
func NewBufferString(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}

// NewBufferReader returns a bytes.Reader over b — the bytes flavour of
// NewReader, which works on strings. Use for HTTP body construction
// and other Reader-shaped consumers.
//
//	body := core.NewBufferReader([]byte(`{"port":8080}`))
//	req, _ := core.HTTPNewRequest("POST", url, body)
func NewBufferReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}

// LimitReader returns a Reader that reads from r but stops with EOF
// after n bytes. Useful for bounding HTTP body reads at a maximum
// size to prevent memory blow-ups from oversized responses.
//
//	body := core.ReadAll(core.LimitReader(resp.Body, 4<<20))
//	if !body.OK { return body }
//	bytes := body.Value.([]byte)
func LimitReader(r Reader, n int64) Reader {
	return io.LimitReader(r, n)
}
