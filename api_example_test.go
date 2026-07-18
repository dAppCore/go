package core_test

import . "dappco.re/go"

type exampleStream struct {
	response []byte
	sent     []byte
}

func (s *exampleStream) Send(data []byte) error {
	s.sent = data
	return nil
}

func (s *exampleStream) Receive() ([]byte, error) {
	return s.response, nil
}

func (s *exampleStream) Close() error {
	return nil
}

// ExampleAPI_RegisterProtocol registers a transport protocol through
// `API.RegisterProtocol` for a Lethean drive integration. Transport details stay behind
// the API wrapper while callers exchange drives, streams, and Results.
func ExampleAPI_RegisterProtocol() {
	c := New()
	c.API().RegisterProtocol("http", func(h *DriveHandle) (Stream, error) {
		return &exampleStream{response: []byte("pong")}, nil
	})
	Println(c.API().Protocols())
	// Output: [http]
}

// ExampleAPI_Stream opens a stream through `API.Stream` for a Lethean drive integration.
// Transport details stay behind the API wrapper while callers exchange drives, streams,
// and Results.
func ExampleAPI_Stream() {
	c := New()
	c.API().RegisterProtocol("http", func(h *DriveHandle) (Stream, error) {
		return &exampleStream{response: []byte(Concat("connected to ", h.Name))}, nil
	})
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "charon"},
		Option{Key: "transport", Value: "http://10.69.69.165:9101"},
	))

	r := c.API().Stream("charon")
	if r.OK {
		stream := r.Value.(Stream)
		resp, _ := stream.Receive()
		Println(string(resp))
		stream.Close()
	}
	// Output: connected to charon
}

// ExampleAPI_Call calls a remote method through `API.Call` for a Lethean drive
// integration. Transport details stay behind the API wrapper while callers exchange
// drives, streams, and Results.
func ExampleAPI_Call() {
	c := New()
	c.API().RegisterProtocol("http", func(_ *DriveHandle) (Stream, error) {
		return &exampleStream{response: []byte(`{"ok":true}`)}, nil
	})
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "charon"},
		Option{Key: "transport", Value: "http://10.69.69.165:9101"},
	))

	r := c.API().Call("charon", "agent.status", NewOptions())
	Println(r.Value)
	// Output: {"ok":true}
}

// ExampleAPI_Protocols lists transport protocols through `API.Protocols` for a Lethean
// drive integration. Transport details stay behind the API wrapper while callers exchange
// drives, streams, and Results.
func ExampleAPI_Protocols() {
	c := New()
	c.API().RegisterProtocol("http", func(_ *DriveHandle) (Stream, error) {
		return &exampleStream{}, nil
	})
	Println(c.API().Protocols())
	// Output: [http]
}

// ExampleCore_RemoteAction resolves an action name locally before a remote drive prefix is
// needed. Transport details stay behind the API wrapper while callers exchange drives,
// streams, and Results.
func ExampleCore_RemoteAction() {
	c := New()
	// Local action
	c.Action("agent.status", func(_ Context, _ Options) Result {
		return Result{Value: "running", OK: true}
	})

	// No colon — resolves locally
	r := c.RemoteAction("agent.status", Background(), NewOptions())
	Println(r.Value)
	// Output: running
}

// ExampleHTTPGet fetches a local health endpoint through the core HTTP client wrapper for
// a Lethean drive integration. Transport details stay behind the API wrapper while callers
// exchange drives, streams, and Results.
func ExampleHTTPGet() {
	srv := NewHTTPTestServer(HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "ok")
	}))
	defer srv.Close()

	r := HTTPGet(srv.URL)
	defer r.Value.(*Response).Body.Close()
	body := ReadAll(r.Value.(*Response).Body)
	Println(body.Value)
	// Output: ok
}

// ExampleHTTPPost sends a reader-backed payload through the core HTTP client wrapper for a
// Lethean drive integration. Transport details stay behind the API wrapper while callers
// exchange drives, streams, and Results.
func ExampleHTTPPost() {
	srv := NewHTTPTestServer(HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "created")
	}))
	defer srv.Close()

	r := HTTPPost(srv.URL, "text/plain", NewReader("payload"))
	defer r.Value.(*Response).Body.Close()
	body := ReadAll(r.Value.(*Response).Body)
	Println(body.Value)
	// Output: created
}

// ExampleHTTPPostForm submits form data through the core HTTP client wrapper for a Lethean
// drive integration. Transport details stay behind the API wrapper while callers exchange
// drives, streams, and Results.
func ExampleHTTPPostForm() {
	srv := NewHTTPTestServer(HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "submitted")
	}))
	defer srv.Close()

	r := HTTPPostForm(srv.URL, nil)
	defer r.Value.(*Response).Body.Close()
	body := ReadAll(r.Value.(*Response).Body)
	Println(body.Value)
	// Output: submitted
}

// ExampleNewHTTPRequest builds a POST request with a payload for a deployment endpoint.
// Transport details stay behind the API wrapper while callers exchange drives, streams,
// and Results.
func ExampleNewHTTPRequest() {
	r := NewHTTPRequest("POST", "https://example.com/deploy", NewReader("payload"))
	req := r.Value.(*Request)
	Println(req.Method)
	Println(req.URL.Path)
	// Output:
	// POST
	// /deploy
}

// ExampleNewHTTPRequestContext builds a request bound to the active Core context for a
// status endpoint. Transport details stay behind the API wrapper while callers exchange
// drives, streams, and Results.
func ExampleNewHTTPRequestContext() {
	ctx := New().Context()
	r := NewHTTPRequestContext(ctx, "GET", "https://example.com/status", nil)
	req := r.Value.(*Request)
	Println(req.Context() == ctx)
	Println(req.Method)
	// Output:
	// true
	// GET
}

// ExampleHTTPStatusText reads status reason text through `HTTPStatusText` for a Lethean
// drive integration. Transport details stay behind the API wrapper while callers exchange
// drives, streams, and Results.
func ExampleHTTPStatusText() {
	Println(HTTPStatusText(201))
	// Output: Created
}

// ExampleNewMultipartWriter creates multipart form data through `NewMultipartWriter` for a
// Lethean drive integration. Transport details stay behind the API wrapper while callers
// exchange drives, streams, and Results.
func ExampleNewMultipartWriter() {
	buf := NewBuffer()
	writer := NewMultipartWriter(buf)
	writer.WriteField("name", "codex")
	boundary := writer.Boundary()
	writer.Close()

	Println(boundary != "")
	Println(Contains(buf.String(), "codex"))
	// Output:
	// true
	// true
}

// ExampleNewMultipartReader reads multipart form data through `NewMultipartReader` for a
// Lethean drive integration. Transport details stay behind the API wrapper while callers
// exchange drives, streams, and Results.
func ExampleNewMultipartReader() {
	buf := NewBuffer()
	writer := NewMultipartWriter(buf)
	writer.WriteField("name", "codex")
	boundary := writer.Boundary()
	writer.Close()

	reader := NewMultipartReader(buf, boundary)
	part, _ := reader.NextPart()
	data := ReadAll(part)
	Println(part.FormName())
	Println(data.Value)
	// Output:
	// name
	// codex
}

// ExampleNewHTTPTestServer creates an HTTP fixture server for handler validation.
// Transport details stay behind the API wrapper while callers exchange drives, streams,
// and Results.
func ExampleNewHTTPTestServer() {
	srv := NewHTTPTestServer(HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "ok")
	}))
	defer srv.Close()

	Println(HasPrefix(srv.URL, "http://"))
	// Output: true
}

// ExampleNewHTTPTestTLSServer creates a TLS fixture server for HTTPS handler validation.
// Transport details stay behind the API wrapper while callers exchange drives, streams,
// and Results.
func ExampleNewHTTPTestTLSServer() {
	srv := NewHTTPTestTLSServer(HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "ok")
	}))
	defer srv.Close()

	Println(HasPrefix(srv.URL, "https://"))
	// Output: true
}

// ExampleNewHTTPTestRecorder records response status for handler validation. Transport
// details stay behind the API wrapper while callers exchange drives, streams, and Results.
func ExampleNewHTTPTestRecorder() {
	rec := NewHTTPTestRecorder()
	rec.WriteHeader(202)
	Println(rec.Code)
	// Output: 202
}

// ExampleNewHTTPTestRequest builds a request fixture for a status route. Transport details
// stay behind the API wrapper while callers exchange drives, streams, and Results.
func ExampleNewHTTPTestRequest() {
	req := NewHTTPTestRequest("GET", "/status", nil)
	Println(req.Method)
	Println(req.URL.Path)
	// Output:
	// GET
	// /status
}

// ExampleMethodGet shows the HTTP method constants exposed by core. They are
// the canonical method strings re-exported from net/http so consumers build
// requests without importing net/http directly.
func ExampleMethodGet() {
	Println(MethodGet)
	Println(MethodPost)
	Println(MethodDelete)
	// Output:
	// GET
	// POST
	// DELETE
}

// ExampleStatusOK shows the HTTP status constants exposed by core. The
// most-reached codes (success, client error, server error) are listed
// here; the full set follows the http.StatusXxx naming pattern.
func ExampleStatusOK() {
	Println(StatusOK)
	Println(StatusBadRequest)
	Println(StatusInternalServerError)
	// Output:
	// 200
	// 400
	// 500
}

// ExampleFlusher demonstrates type-asserting a ResponseWriter to Flusher
// for streaming responses. The handler pushes one chunk and flushes
// immediately — the canonical Server-Sent Events / chunked-transfer shape.
func ExampleFlusher() {
	handler := HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "tick")
		if f, ok := w.(Flusher); ok {
			f.Flush()
		}
	})
	srv := NewHTTPTestServer(handler)
	defer srv.Close()

	r := HTTPGet(srv.URL)
	defer r.Value.(*Response).Body.Close()
	body := ReadAll(r.Value.(*Response).Body)
	Println(body.Value)
	// Output: tick
}

// ExampleHTTPFileSystem documents the type alias for HTTP file servers.
// Construct via HTTPFS to wrap a Lethean FS, or pass a stdlib http.Dir.
func ExampleHTTPFileSystem() {
	var _ HTTPFileSystem = HTTPFS(DirFS("/tmp"))
	Println("aliased")
	// Output: aliased
}

// ExampleDefaultHTTPClient uses the package-level default *HTTPClient for a
// one-off request that doesn't justify a dedicated client.
func ExampleDefaultHTTPClient() {
	srv := NewHTTPTestServer(HandlerFunc(func(w ResponseWriter, _ *Request) {
		WriteString(w, "ok")
	}))
	defer srv.Close()

	req := NewHTTPRequest(MethodGet, srv.URL, nil).Value.(*Request)
	resp, _ := DefaultHTTPClient.Do(req)
	defer resp.Body.Close()
	body := ReadAll(resp.Body)
	Println(body.Value)
	// Output: ok
}

// ExampleErrHTTPServerClosed checks for the graceful-shutdown sentinel.
// A real HTTPServer returns this from ListenAndServe after Shutdown is
// called; here the example demonstrates the equality-check shape via Is.
func ExampleErrHTTPServerClosed() {
	err := ErrHTTPServerClosed
	Println(Is(err, ErrHTTPServerClosed))
	// Output: true
}

// ExampleNewServeMux composes a small mux + serves a request via the
// test server. Same shape consumers use to register agent endpoints
// without importing net/http directly.
func ExampleNewServeMux() {
	mux := NewServeMux()
	mux.HandleFunc("/health", func(w ResponseWriter, _ *Request) {
		WriteString(w, "ok")
	})
	srv := NewHTTPTestServer(mux)
	defer srv.Close()

	r := HTTPGet(srv.URL + "/health")
	defer r.Value.(*Response).Body.Close()
	body := ReadAll(r.Value.(*Response).Body)
	Println(body.Value)
	// Output: ok
}

// ExampleHTTPStripPrefix mounts an inner handler under /api/v1/* by
// stripping the prefix before delegation. Used to compose sub-apps
// under a versioned path.
func ExampleHTTPStripPrefix() {
	inner := HandlerFunc(func(w ResponseWriter, r *Request) {
		WriteString(w, r.URL.Path)
	})
	mux := NewServeMux()
	mux.Handle("/api/v1/", HTTPStripPrefix("/api/v1", inner))
	srv := NewHTTPTestServer(mux)
	defer srv.Close()

	r := HTTPGet(srv.URL + "/api/v1/users")
	defer r.Value.(*Response).Body.Close()
	body := ReadAll(r.Value.(*Response).Body)
	Println(body.Value)
	// Output: /users
}

// ExampleHTTPListenAndServe documents the function signature without
// actually starting a server (would block forever). In production,
// pair with a Shutdown call on signal.received via a goroutine.
func ExampleHTTPListenAndServe() {
	mux := NewServeMux()
	_ = HTTPListenAndServe // documented; not invoked
	_ = mux
}

// ExampleHTTPFileServer constructs a Handler that serves a file tree
// under any mux path. Pair with HTTPFS to serve a Lethean FS or
// embed.FS root.
func ExampleHTTPFileServer() {
	mux := NewServeMux()
	mux.Handle("/static/", HTTPStripPrefix("/static/", HTTPFileServer(HTTPFS(DirFS("/tmp")))))
	_ = mux
}

// ExampleHTTPFS converts a Lethean FS to an HTTPFileSystem suitable for
// HTTPFileServer. Combined with embed.FS, the pattern serves static
// assets without a network read.
func ExampleHTTPFS() {
	hfs := HTTPFS(DirFS("/tmp"))
	_ = hfs
}

// ExampleHTTPError writes a plain-text error body with the given status
// code. The handler exits after this call; the Recorder captures the
// response for assertion in tests.
// ExampleCore_API returns the HTTP API subsystem through `Core.API`.
func ExampleCore_API() {
	Println(New().API() != nil)
	// Output: true
}

func ExampleHTTPError() {
	rec := NewHTTPTestRecorder()
	HTTPError(rec, "missing field", StatusBadRequest)
	Println(rec.Code)
	Println(rec.Body.String())
	// Output:
	// 400
	// missing field
}
