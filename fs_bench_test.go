// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the filesystem primitives in fs.go.
// Per AX-11 — Fs is on the load-bearing path for go-mlx model file
// loading (GGUF / safetensors reads), agent config persistence, log
// rotation, and every dapp.workspace operation. The wrappers add
// sandbox path validation on top of the stdlib os.* primitives so
// regressions HERE cascade across the ecosystem.
//
// Run:    go test -bench='BenchmarkFs' -benchmem -run='^$' .

package core_test

import (
	"os"
	"path/filepath"
	"testing"

	. "dappco.re/go"
)

// Sinks defeat compiler DCE on the read paths.
var (
	fsSinkResult Result
	fsSinkBool   bool
	fsSinkString string
)

// fsBenchFixture builds a *Fs over a temp dir pre-populated with files
// of various sizes. Returns the *Fs plus a cleanup func the bench can
// defer.
func fsBenchFixture(tb testing.TB, sizes map[string]int) *Fs {
	tb.Helper()
	dir, err := os.MkdirTemp("", "fsbench-*")
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() { os.RemoveAll(dir) })
	// On macOS, /var/folders is a symlink to /private/var/folders.
	// Fs.validatePath resolves symlinks during sandbox checks, so the
	// fixture must use the symlink-resolved path to avoid spurious
	// "sandbox escape detected" warnings during benches.
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}

	// Sub-directory with a few files for List/Walk benches.
	if err := os.Mkdir(filepath.Join(dir, "models"), 0o755); err != nil {
		tb.Fatal(err)
	}
	for name, size := range sizes {
		payload := make([]byte, size)
		for i := range payload {
			payload[i] = byte('a' + (i % 26))
		}
		if err := os.WriteFile(filepath.Join(dir, name), payload, 0o644); err != nil {
			tb.Fatal(err)
		}
	}
	for i, name := range []string{"qwen.gguf", "llama.gguf", "config.json"} {
		_ = i
		if err := os.WriteFile(filepath.Join(dir, "models", name), []byte("placeholder"), 0o644); err != nil {
			tb.Fatal(err)
		}
	}

	return (&Fs{}).New(dir)
}

// --- Read ---

func BenchmarkFs_Read_1KB(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 1024})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Read("f.bin")
	}
}

func BenchmarkFs_Read_64KB(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 64 * 1024})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Read("f.bin")
	}
}

func BenchmarkFs_Read_1MB(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 1024 * 1024})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Read("f.bin")
	}
}

// --- Existence checks ---

func BenchmarkFs_Exists_Hit(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 128})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Exists("f.bin")
	}
}

func BenchmarkFs_Exists_Miss(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 128})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Exists("missing.bin")
	}
}

func BenchmarkFs_IsFile(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 128})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.IsFile("f.bin")
	}
}

func BenchmarkFs_IsDir(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 128})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.IsDir("models")
	}
}

func BenchmarkFs_Stat(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 128})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Stat("f.bin")
	}
}

// --- List / TempDir ---

func BenchmarkFs_List_SmallDir(b *B) {
	fs := fsBenchFixture(b, map[string]int{"f.bin": 128})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.List("models")
	}
}

func BenchmarkFs_TempDir(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.TempDir("agent-")
	}
}

// --- Root / New ---

func BenchmarkFs_Root(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkString = fs.Root()
	}
}

func BenchmarkFs_New(b *B) {
	dir, _ := os.MkdirTemp("", "fsbench-new-")
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}
	b.Cleanup(func() { os.RemoveAll(dir) })
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = (&Fs{}).New(dir)
	}
}

// --- Open / Create / Delete ---
//
// These mutate state, so they cycle through a unique path per iteration
// and clean up between runs. The bench wall-time captures the cycle
// cost (creation + close) which is the realistic per-op floor.

func BenchmarkFs_Create(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Create("create.bin")
		if fsSinkResult.OK {
			CloseStream(fsSinkResult.Value)
		}
		fs.Delete("create.bin")
	}
}

// --- Write ---

func BenchmarkFs_Write_Small(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	content := "x"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Write("w.bin", content)
	}
}

func BenchmarkFs_Write_1KB(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	content := string(make([]byte, 1024))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Write("w.bin", content)
	}
}

func BenchmarkFs_WriteAtomic_1KB(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	content := string(make([]byte, 1024))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.WriteAtomic("w.bin", content)
	}
}

// --- Append / Stream ---

func BenchmarkFs_Append_Small(b *B) {
	fs := fsBenchFixture(b, map[string]int{"appendme.txt": 0})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Append("appendme.txt")
		if fsSinkResult.OK {
			CloseStream(fsSinkResult.Value)
		}
	}
}

func BenchmarkFs_ReadStream(b *B) {
	fs := fsBenchFixture(b, map[string]int{"stream.bin": 4096})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.ReadStream("stream.bin")
		if fsSinkResult.OK {
			CloseStream(fsSinkResult.Value)
		}
	}
}

// --- EnsureDir / DeleteAll / Rename ---

func BenchmarkFs_EnsureDir_Exists(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.EnsureDir("models")
	}
}

func BenchmarkFs_EnsureDir_Create(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Cycle the same dir so the bench captures the make-then-noop floor.
		fsSinkResult = fs.EnsureDir("ephemeral")
	}
}

func BenchmarkFs_DeleteAll(b *B) {
	fs := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fs.EnsureDir("nuked")
		fsSinkResult = fs.DeleteAll("nuked")
	}
}

func BenchmarkFs_Rename(b *B) {
	fs := fsBenchFixture(b, map[string]int{"rename-src.bin": 16})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = fs.Rename("rename-src.bin", "rename-dst.bin")
		if fsSinkResult.OK {
			// Cycle back so the next iteration has the src again.
			fs.Rename("rename-dst.bin", "rename-src.bin")
		}
	}
}

// --- NewUnrestricted ---

func BenchmarkFs_NewUnrestricted(b *B) {
	parent := fsBenchFixture(b, map[string]int{})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = parent.NewUnrestricted()
	}
}

// --- ReadDir / ReadFSFile (the embed-side helpers) ---

func BenchmarkFs_ReadDir_FS(b *B) {
	fsys := DirFS(".")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = ReadDir(fsys, ".")
	}
}

func BenchmarkFs_ReadFSFile(b *B) {
	fs := fsBenchFixture(b, map[string]int{"r.bin": 1024})
	fsys := DirFS(fs.Root())
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = ReadFSFile(fsys, "r.bin")
	}
}

func BenchmarkFs_WriteStream(b *B) {
	fs := fsBenchFixture(b, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := fs.WriteStream("s" + Itoa(i))
		if wc, ok := r.Value.(WriteCloser); ok {
			wc.Close()
		}
	}
}

func BenchmarkWriteAll(b *B) {
	fs := fsBenchFixture(b, nil)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = WriteAll(fs.WriteStream("w"+Itoa(i)).Value, "payload")
	}
}

func BenchmarkFs_WalkSeq(b *B) {
	fs := fsBenchFixture(b, map[string]int{"a.txt": 64, "b.txt": 64})
	root := fs.Root()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for entry, err := range fs.WalkSeq(root) {
			_, _ = entry, err
		}
	}
}

func BenchmarkFs_WalkSeqSkip(b *B) {
	fs := fsBenchFixture(b, map[string]int{"a.txt": 64, "b.txt": 64})
	root := fs.Root()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for entry, err := range fs.WalkSeqSkip(root, "vendor") {
			_, _ = entry, err
		}
	}
}

func BenchmarkWalkDir(b *B) {
	fs := fsBenchFixture(b, map[string]int{"a.txt": 64})
	fsys := DirFS(fs.Root())
	fn := func(path string, d FsDirEntry, err error) error { return nil }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fsSinkResult = WalkDir(fsys, ".", fn)
	}
}
