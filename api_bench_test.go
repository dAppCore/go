// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the HTTP primitives in api.go.
// Per AX-11 — api.go sits on every HTTP-touching path in the
// ecosystem: consumer services, dappco.re/go-api, MCP transport,
// HTTPTestServer-driven test suites. Each public function deserves
// a perf contract so that regressions surface against a published
// floor rather than only after a downstream consumer complains.
//
// Coverage matches the verbosity of api_test.go + api_example_test.go
// — one benchmark per public function (plus shape variants where the
// behaviour depends on input size). HTTP* network-bound functions are
// benched against an in-process *httptest.Server so the numbers are
// stable on any machine.
//
// Run:    go test -bench='Benchmark.*HTTP|BenchmarkServe|BenchmarkMultipart' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

// --- Fixtures ---

var (
	apiBenchHandler = HandlerFunc(func(w ResponseWriter, r *Request) {
		_, _ = WriteToResponse(w, "ok")
	})

	apiBenchURLForm = URLValues{
		"name":  []string{"cladius"},
		"agent": []string{"go"},
	}
)

// WriteToResponse is a small adapter to write a response without bringing
// io into the bench file's surface. Mirrors what consumer handlers do.
func WriteToResponse(w ResponseWriter, body string) (int, error) {
	return w.Write([]byte(body))
}

// --- Constructors / cheap wrappers ---

func BenchmarkNewServeMux(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewServeMux()
	}
}

func BenchmarkHTTPStripPrefix(b *B) {
	h := apiBenchHandler
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTTPStripPrefix("/api", h)
	}
}

func BenchmarkHTTPStatusText_200(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTTPStatusText(200)
	}
}

func BenchmarkHTTPStatusText_404(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTTPStatusText(404)
	}
}

func BenchmarkHTTPStatusText_500(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTTPStatusText(500)
	}
}

// --- Request construction ---

func BenchmarkNewHTTPRequest_GET(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewHTTPRequest("GET", "https://forge.lthn.sh/api/health", nil)
	}
}

func BenchmarkNewHTTPRequest_POSTBody(b *B) {
	body := []byte(`{"port":8080,"host":"homelab.lan"}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewHTTPRequest("POST", "https://forge.lthn.sh/api/agent", NewBufferReader(body))
	}
}

func BenchmarkNewHTTPRequestContext(b *B) {
	ctx := Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewHTTPRequestContext(ctx, "GET", "https://forge.lthn.sh/api/health", nil)
	}
}

func BenchmarkNewHTTPTestRequest(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewHTTPTestRequest("GET", "/health", nil)
	}
}

// --- HTTPError ---

func BenchmarkHTTPError(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rec := NewHTTPTestRecorder()
		HTTPError(rec, "not found", 404)
	}
}

// --- HTTP round-trip via httptest.Server ---
//
// These exercise the full HTTPGet/HTTPPost path against an in-process
// server so the numbers reflect the wrapper cost on top of stdlib's
// net/http — stable across machines because there's no real network.

func BenchmarkHTTPGet(b *B) {
	srv := NewHTTPTestServer(apiBenchHandler)
	defer srv.Close()
	url := srv.URL
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := HTTPGet(url)
		if r.OK {
			_ = ReadAll(r.Value.(*Response).Body)
		}
	}
}

func BenchmarkHTTPPost(b *B) {
	srv := NewHTTPTestServer(apiBenchHandler)
	defer srv.Close()
	url := srv.URL
	body := []byte(`{"name":"cladius"}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := HTTPPost(url, "application/json", NewBufferReader(body))
		if r.OK {
			_ = ReadAll(r.Value.(*Response).Body)
		}
	}
}

func BenchmarkHTTPPostForm(b *B) {
	srv := NewHTTPTestServer(apiBenchHandler)
	defer srv.Close()
	url := srv.URL
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := HTTPPostForm(url, apiBenchURLForm)
		if r.OK {
			_ = ReadAll(r.Value.(*Response).Body)
		}
	}
}

// --- Test-server / recorder constructors ---

func BenchmarkNewHTTPTestServer(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv := NewHTTPTestServer(apiBenchHandler)
		srv.Close()
	}
}

func BenchmarkNewHTTPTestRecorder(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewHTTPTestRecorder()
	}
}

// --- Multipart ---

func BenchmarkNewMultipartReader(b *B) {
	body := NewBufferReader([]byte("--boundary\r\n"))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewMultipartReader(body, "boundary")
	}
}

func BenchmarkNewMultipartWriter(b *B) {
	dst := NewBuffer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst.Reset()
		_ = NewMultipartWriter(dst)
	}
}

// --- High-level API struct (Stream / Call / RegisterProtocol) ---
//
// The API.Stream/Call path needs a configured Drive + protocol factory.
// Benched here against a tiny in-memory factory so the numbers reflect
// registry lookup + dispatch, not network cost.

func BenchmarkAPI_RegisterProtocol(b *B) {
	c := New()
	factory := func(*DriveHandle) (Stream, error) { return nil, nil }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.API().RegisterProtocol("bench", factory)
	}
}

func BenchmarkAPI_Protocols(b *B) {
	c := New()
	factory := func(*DriveHandle) (Stream, error) { return nil, nil }
	c.API().RegisterProtocol("bench", factory)
	c.API().RegisterProtocol("bench2", factory)
	c.API().RegisterProtocol("bench3", factory)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.API().Protocols()
	}
}

// --- HTTP file server constructors ---
//
// HTTPListenAndServe blocks on a network bind, so we can't run it in a
// bench loop. The pure constructors HTTPFileServer + HTTPFS are cheap
// adapter-builders worth gating, and NewHTTPTestTLSServer is the TLS
// twin of NewHTTPTestServer.

func BenchmarkHTTPFS(b *B) {
	fsys := DirFS(".")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTTPFS(fsys)
	}
}

func BenchmarkHTTPFileServer(b *B) {
	root := HTTPFS(DirFS("."))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = HTTPFileServer(root)
	}
}

func BenchmarkNewHTTPTestTLSServer(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		srv := NewHTTPTestTLSServer(apiBenchHandler)
		srv.Close()
	}
}

// --- remote-endpoint surface (Drive + protocol resolution) ---
//
// A mock protocol factory stands in for a live transport so these bench
// the resolve + scheme-extraction + call dispatch path, not the network.
// HTTPListenAndServe is omitted — it blocks serving until shut down.

var apiSinkResult Result

func apiRemoteFixture() *Core {
	c := New()
	c.API().RegisterProtocol("http", mockFactory("pong"))
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "charon"},
		Option{Key: "transport", Value: "http://127.0.0.1:9101/mcp"},
	))
	return c
}

func BenchmarkAPI_Stream(b *B) {
	c := apiRemoteFixture()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := c.API().Stream("charon")
		if r.OK {
			r.Value.(Stream).Close()
		}
	}
}

func BenchmarkAPI_Call(b *B) {
	c := apiRemoteFixture()
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		apiSinkResult = c.API().Call("charon", "agentic.status", opts)
	}
}

func BenchmarkCore_RemoteAction(b *B) {
	c := apiRemoteFixture()
	ctx := Background()
	opts := NewOptions()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		apiSinkResult = c.RemoteAction("charon:agentic.status", ctx, opts)
	}
}

// HTTPListenAndServe is benched on its error path: an invalid address
// fails the bind immediately and returns the error Result, with no real
// socket. The success path blocks serving, so it can't be unit-benched.
func BenchmarkHTTPListenAndServe(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		apiSinkResult = Result{Value: HTTPListenAndServe("127.0.0.1:-1", nil)}
	}
}
