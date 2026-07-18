// Sandboxed local filesystem I/O for the Core framework.
package core

import (
	"io/fs"
)

// Fs is a sandboxed local filesystem backend.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.EnsureDir("logs")
//	if !r.OK { return r }
type Fs struct {
	root            string
	rootResolved    string // symlink-resolved root, computed once at New()
	rootResolvedSep string // rootResolved + PathSeparator — prefix used by validatePath

	// validated caches successful validatePath results keyed by the
	// caller-supplied path. Fs operations re-validate per call by
	// default; storing the result here lets repeat lookups skip the
	// EvalSymlinks syscall chain (per-component Lstat, intermediate
	// string builds) which dominated the per-call cost.
	//
	// Security shape: the cache is unconditionally safe. First call
	// performs the full sandbox-escape check; subsequent calls return
	// the same already-validated path. This is a TOCTOU IMPROVEMENT
	// over the per-call path — an attacker who races a symlink swap
	// between validate and use can no longer redirect already-resolved
	// paths. Cache miss on any newly-presented path still does the
	// full check.
	validated SyncMap // string -> string
}

// FS is a generic filesystem accepted by Mount and Extract.
//
//	fsys := core.DirFS("templates")
//	r := core.Mount(fsys, ".")
type FS = fs.FS

// FsFile is a file opened from an FS.
//
//	r := emb.Open("README.md")
//	if r.OK { file := r.Value.(core.FsFile); _ = file }
type FsFile = fs.File

// FsFileInfo describes a file returned by Stat and PathWalk.
//
//	info := stat.Value.(core.FsFileInfo)
//	_ = info
type FsFileInfo = fs.FileInfo

// FsDirEntry is a directory entry returned by filesystem walkers.
//
//	r := emb.ReadDir(".")
//	if r.OK { entries := r.Value.([]core.FsDirEntry); _ = entries }
type FsDirEntry = fs.DirEntry

// WalkDirFunc visits one path during a filesystem walk.
//
//	fn := func(path string, d core.FsDirEntry, err error) error { return err }
//	_ = fn
type WalkDirFunc = fs.WalkDirFunc

// New initialises an Fs with the given root directory.
// Root "/" means unrestricted access. Empty root defaults to "/".
//
//	fs := (&core.Fs{}).New("/")
func (m *Fs) New(root string) *Fs {
	if root == "" {
		root = "/"
	}
	m.root = root
	// Resolve root symlinks ONCE at construction. validatePath uses
	// the resolved form for the sandbox-escape check; without this
	// every call paid an EvalSymlinks per path component to detect
	// (e.g.) /var → /private/var on macOS. With it, only the per-call
	// path itself needs resolving.
	m.rootResolved = root
	if root != "/" {
		if r := PathEvalSymlinks(root); r.OK {
			m.rootResolved = r.Value.(string)
		}
	}
	// Pre-compute rootResolved + sep so validatePath's HasPrefix check
	// is a single load — Concat at call time would alloc per Fs op.
	m.rootResolvedSep = m.rootResolved + string(PathSeparator)
	return m
}

// NewUnrestricted returns a new Fs with root "/", granting full filesystem access.
// Use this instead of unsafe.Pointer to bypass the sandbox.
//
//	fs := c.Fs().NewUnrestricted()
//	fs.Read("/etc/hostname")  // works — no sandbox
func (m *Fs) NewUnrestricted() *Fs {
	return (&Fs{}).New("/")
}

// Root returns the sandbox root path.
//
//	root := c.Fs().Root()  // e.g. "/home/agent/.core"
func (m *Fs) Root() string {
	if m.root == "" {
		return "/"
	}
	return m.root
}

// path sanitises and returns the full path.
// Absolute paths are sandboxed under root (unless root is "/").
// Empty root defaults to "/" — the zero value of Fs is usable.
func (m *Fs) path(p string) string {
	root := m.root
	if root == "" {
		root = "/"
	}
	if p == "" {
		return root
	}

	// If the path is relative and the medium is rooted at "/",
	// treat it as relative to the current working directory.
	// This makes io.Local behave more like the standard 'os' package.
	if root == "/" && !PathIsAbs(p) {
		cwd := ""
		if r := Getwd(); r.OK {
			cwd = r.Value.(string)
		}
		return PathJoin(cwd, p)
	}

	// Use a leading slash to resolve all .. and . internally
	// before joining with the root. This is a standard way to sandbox paths.
	clean := CleanPath("/"+p, string(PathSeparator))

	// If root is "/", allow absolute paths through
	if root == "/" {
		return clean
	}

	// Strip leading "/" so Join works correctly with root
	return PathJoin(root, clean[1:])
}

// validatePath ensures the path is within the sandbox, following symlinks if they exist.
//
// Resolves symlinks at the deepest existing prefix of the joined path
// (one syscall for fully-existing paths, walks up only for paths with
// not-yet-created suffixes). Compares the resolved form against
// m.rootResolved (cached once at New()). The old per-component walk
// did O(N) syscalls + allocations even for paths that already existed;
// this is O(1) for the common case.
//
// Security invariants preserved:
//   - CleanPath strips ".." escape attempts before path construction.
//   - PathEvalSymlinks chases the full symlink chain, so any link at
//     any depth pointing outside root is caught by the PathRel check.
//   - rootResolved is computed once at construction; symlinks added
//     later cannot retroactively change the sandbox.
func (m *Fs) validatePath(p string) Result {
	root := m.rootResolved
	if root == "" {
		root = m.root
		if root == "" {
			root = "/"
		}
	}
	if root == "/" {
		return Result{m.path(p), true}
	}

	// Cached fast path — same input path resolves to the same
	// caller-form output every time within an Fs lifetime. The cache
	// holds only OK results; rejected paths fall through to the full
	// check on every retry (cheap because rare).
	if v, ok := m.validated.Load(p); ok {
		return Result{v.(string), true}
	}

	// Build the full candidate path under the resolved root.
	clean := CleanPath("/"+p, string(PathSeparator))
	full := PathJoin(root, clean[1:])

	resolved := evalDeepestSymlinks(full)

	// Sandbox check via string-prefix instead of PathRel — equivalent
	// boundary (any escape produces a resolved path that fails both
	// `== root` and `HasPrefix(root + sep)`), zero allocation. The
	// PathRel approach allocated 3-5 intermediate strings per call.
	if resolved != root && !HasPrefix(resolved, m.rootResolvedSep) {
		username := "unknown"
		if r := UserCurrent(); r.OK {
			username = r.Value.(*User).Username
		}
		Print(Stderr(), "[%s] SECURITY sandbox escape detected root=%s path=%s attempted=%s user=%s",
			Now().Format(TimeRFC3339), m.root, p, resolved, username)
		return Result{E("fs.validatePath", Concat("sandbox escape: ", p, " resolves outside ", m.root), nil), false}
	}
	// Translate back to the caller's root form (the user-visible
	// contract is "paths under m.root"). Internally we used rootResolved
	// for the sandbox check; externally callers want paths matching
	// what they passed to New().
	var out string
	if m.root == root {
		out = resolved
	} else if resolved == root {
		out = m.root
	} else {
		// Strip the rootResolved prefix + sep to get the tail under root.
		out = PathJoin(m.root, resolved[len(m.rootResolvedSep):])
	}
	m.validated.Store(p, out)
	return Result{out, true}
}

// evalDeepestSymlinks resolves symlinks at the deepest existing prefix
// of p. For paths that fully exist, this is one EvalSymlinks call. For
// paths whose suffix doesn't yet exist (e.g. about to be Create'd),
// walks up parent-by-parent until a resolvable prefix is found, then
// appends the unresolved tail. Returns p unchanged if no prefix
// resolves (extremely unlikely — / is always resolvable).
func evalDeepestSymlinks(p string) string {
	if r := PathEvalSymlinks(p); r.OK {
		return r.Value.(string)
	}
	parent := PathDir(p)
	if parent == p || parent == "." {
		return p
	}
	return PathJoin(evalDeepestSymlinks(parent), PathBase(p))
}

// Read returns file contents as string.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Read("config/agent.json")
//	if r.OK { core.Println(r.Value.(string)) }
func (m *Fs) Read(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	r := ReadFile(vp.Value.(string))
	if !r.OK {
		return r
	}
	// ReadFile returns a freshly-allocated []byte that becomes
	// unreachable after this conversion — AsString skips the copy.
	// Same safety contract as ReadAll's fast path.
	return Result{AsString(r.Value.([]byte)), true}
}

// Write saves content to file, creating parent directories as needed.
// Files are created with mode 0644. For sensitive files (keys, secrets),
// use WriteMode with 0600.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Write("config/agent.json", `{"host":"homelab.lthn.sh"}`)
//	if !r.OK { return r }
func (m *Fs) Write(p, content string) Result {
	return m.WriteMode(p, content, 0644)
}

// WriteMode saves content to file with explicit permissions.
// Use 0600 for sensitive files (encryption output, private keys, auth hashes).
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.WriteMode("secrets/token", "lethean-token", 0o600)
//	if !r.OK { return r }
func (m *Fs) WriteMode(p, content string, mode FileMode) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	full := vp.Value.(string)
	if r := MkdirAll(PathDir(full), 0755); !r.OK {
		return r
	}
	if r := WriteFile(full, AsBytes(content), mode); !r.OK {
		return r
	}
	return Result{OK: true}
}

// TempDir creates a temporary directory, returning its path in Result.Value.
// The caller is responsible for cleanup via fs.DeleteAll(). Returns
// Result{OK: false} carrying the error if the directory could not be created.
//
//	r := fs.TempDir("agent-workspace")
//	if !r.OK {
//		return r
//	}
//	dir := r.Value.(string)
//	defer fs.DeleteAll(dir)
func (m *Fs) TempDir(prefix string) Result {
	r := MkdirTemp("", prefix)
	if !r.OK {
		cause, _ := r.Value.(error)
		return Result{E("fs.TempDir", Concat("could not create temp dir with prefix \"", prefix, "\""), cause), false}
	}
	return r
}

// ReadDir reads a directory from fsys.
//
//	r := core.ReadDir(core.DirFS("templates"), ".")
func ReadDir(fsys FS, name string) Result {
	return Result{}.New(fs.ReadDir(fsys, name))
}

// ReadFSFile reads a file from fsys.
//
//	r := core.ReadFSFile(core.DirFS("templates"), "README.md")
func ReadFSFile(fsys FS, name string) Result {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return Result{err, false}
	}
	return Result{data, true}
}

// Sub returns an FS rooted at dir inside fsys.
//
//	r := core.Sub(core.DirFS("templates"), "agent")
func Sub(fsys FS, dir string) Result {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return Result{err, false}
	}
	return Result{sub, true}
}

// WalkDir walks fsys from root, calling fn for each file or directory.
// Returns Result{OK: true} on a complete walk, or Result{OK: false}
// carrying the error (the first fn error, or a traversal failure).
//
//	r := core.WalkDir(core.DirFS("templates"), ".", fn)
//	if !r.OK {
//		return r
//	}
func WalkDir(fsys FS, root string, fn WalkDirFunc) Result {
	if err := fs.WalkDir(fsys, root, fn); err != nil {
		return Result{E("fs.WalkDir", Concat("walk \"", root, "\" failed"), err), false}
	}
	return Result{OK: true}
}

// WriteAtomic writes content by writing to a temp file then renaming.
// Rename is atomic on POSIX — concurrent readers never see a partial file.
// Use this for status files, config, or any file read from multiple goroutines.
//
//	r := fs.WriteAtomic("/status.json", jsonData)
func (m *Fs) WriteAtomic(p, content string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	full := vp.Value.(string)
	if r := MkdirAll(PathDir(full), 0755); !r.OK {
		return r
	}

	tmp := full + ".tmp." + shortRand()
	if r := WriteFile(tmp, AsBytes(content), 0644); !r.OK {
		return r
	}
	if r := Rename(tmp, full); !r.OK {
		Remove(tmp)
		return r
	}
	return Result{OK: true}
}

// EnsureDir creates directory if it doesn't exist.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.EnsureDir("logs/agent")
//	if !r.OK { return r }
func (m *Fs) EnsureDir(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	if r := MkdirAll(vp.Value.(string), 0755); !r.OK {
		return r
	}
	return Result{OK: true}
}

// IsDir reports whether path is a directory via Result.OK. The false path
// carries the reason — a validation/stat error, or "not a directory" when the
// path exists but is something else — so callers can tell "absent" from "wrong
// kind".
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	if fsys.IsDir("logs").OK { core.Println("logs ready") }
func (m *Fs) IsDir(p string) Result {
	if p == "" {
		return Result{E("fs.IsDir", "empty path", nil), false}
	}
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	r := Stat(vp.Value.(string))
	if !r.OK {
		return r
	}
	if !r.Value.(interface{ IsDir() bool }).IsDir() {
		return Result{E("fs.IsDir", Concat("\"", p, "\" is not a directory"), nil), false}
	}
	return Result{OK: true}
}

// IsFile reports whether path is a regular file via Result.OK. The false path
// carries the reason — a validation/stat error, or "not a regular file" when
// the path exists but is a directory or special file.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	if fsys.IsFile("config/agent.json").OK { core.Println("config ready") }
func (m *Fs) IsFile(p string) Result {
	if p == "" {
		return Result{E("fs.IsFile", "empty path", nil), false}
	}
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	r := Stat(vp.Value.(string))
	if !r.OK {
		return r
	}
	if !r.Value.(interface{ Mode() FileMode }).Mode().IsRegular() {
		return Result{E("fs.IsFile", Concat("\"", p, "\" is not a regular file"), nil), false}
	}
	return Result{OK: true}
}

// Exists reports whether path exists via Result.OK. A non-existent path is a
// valid answer (OK=false), not a swallowed error; a path-validation failure is
// returned as its own Result.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	if fsys.Exists("config/agent.json").OK { core.Println("config present") }
func (m *Fs) Exists(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	return Result{OK: Stat(vp.Value.(string)).OK}
}

// List returns directory entries.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.List("config")
//	if !r.OK { return r }
func (m *Fs) List(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	return ReadDir(DirFS(vp.Value.(string)), ".")
}

// Stat returns file info.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Stat("config/agent.json")
//	if !r.OK { return r }
func (m *Fs) Stat(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	return Stat(vp.Value.(string))
}

// Open opens the named file for reading.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.ReadStream("config/agent.json")
//	if !r.OK { return r }
//	defer r.Value.(core.ReadCloser).Close()
func (m *Fs) Open(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	return Open(vp.Value.(string))
}

// Create creates or truncates the named file.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.WriteStream("logs/agent.log")
//	if !r.OK { return r }
//	defer r.Value.(core.WriteCloser).Close()
func (m *Fs) Create(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	full := vp.Value.(string)
	if r := MkdirAll(PathDir(full), 0755); !r.OK {
		return r
	}
	return Create(full)
}

// Append opens the named file for appending, creating it if it doesn't exist.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Append("logs/agent.log")
//	if !r.OK { return r }
//	defer r.Value.(core.WriteCloser).Close()
func (m *Fs) Append(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	full := vp.Value.(string)
	if r := MkdirAll(PathDir(full), 0755); !r.OK {
		return r
	}
	return OpenFile(full, O_APPEND|O_CREATE|O_WRONLY, 0644)
}

// ReadStream returns a reader for the file content.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.ReadStream("config/agent.json")
//	if !r.OK { return r }
//	defer r.Value.(core.ReadCloser).Close()
func (m *Fs) ReadStream(path string) Result {
	return m.Open(path)
}

// WriteStream returns a writer for the file content.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.WriteStream("logs/agent.log")
//	if !r.OK { return r }
//	defer r.Value.(core.WriteCloser).Close()
func (m *Fs) WriteStream(path string) Result {
	return m.Create(path)
}

// WriteAll writes content to a writer and closes it if it implements Closer.
//
//	r := fs.WriteStream(path)
//	core.WriteAll(r.Value, "content")
func WriteAll(writer any, content string) Result {
	wc, ok := writer.(Writer)
	if !ok {
		return Result{E("core.WriteAll", "not a writer", nil), false}
	}
	_, err := wc.Write(AsBytes(content))
	if closer, ok := writer.(Closer); ok {
		closer.Close()
	}
	if err != nil {
		return Result{err, false}
	}
	return Result{OK: true}
}

// CloseStream closes any value that implements Closer.
//
//	core.CloseStream(r.Value)
func CloseStream(v any) {
	if closer, ok := v.(Closer); ok {
		closer.Close()
	}
}

// Delete removes a file or empty directory.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Delete("logs/old-agent.log")
//	if !r.OK { return r }
func (m *Fs) Delete(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	full := vp.Value.(string)
	if full == "/" || full == Getenv("HOME") {
		return Result{E("fs.Delete", Concat("refusing to delete protected path: ", full), nil), false}
	}
	if r := Remove(full); !r.OK {
		return r
	}
	return Result{OK: true}
}

// DeleteAll removes a file or directory recursively.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.DeleteAll("tmp/session-42")
//	if !r.OK { return r }
func (m *Fs) DeleteAll(p string) Result {
	vp := m.validatePath(p)
	if !vp.OK {
		return vp
	}
	full := vp.Value.(string)
	if full == "/" || full == Getenv("HOME") {
		return Result{E("fs.DeleteAll", Concat("refusing to delete protected path: ", full), nil), false}
	}
	if r := RemoveAll(full); !r.OK {
		return r
	}
	return Result{OK: true}
}

// Rename moves a file or directory.
//
//	fsys := (&core.Fs{}).New("/tmp/agent-workspace")
//	r := fsys.Rename("config/agent.tmp", "config/agent.json")
//	if !r.OK { return r }
func (m *Fs) Rename(oldPath, newPath string) Result {
	oldVp := m.validatePath(oldPath)
	if !oldVp.OK {
		return oldVp
	}
	newVp := m.validatePath(newPath)
	if !newVp.OK {
		return newVp
	}
	return Rename(oldVp.Value.(string), newVp.Value.(string))
}

// FsEntry is a directory entry yielded by WalkSeq and WalkSeqSkip.
// Path is relative to the walk root.
//
//	entry := core.FsEntry{Path: "config/agent.json", Name: "agent.json", IsDir: false, Mode: 0o644}
//	core.Println(entry.Path)
type FsEntry struct {
	Path  string // relative to walk root, OS-native separator
	Name  string // basename
	IsDir bool
	Mode  fs.FileMode
}

// WalkSeq walks the directory tree rooted at root within the Fs sandbox,
// yielding every entry depth-first. Iteration stops on caller break.
//
// Symlinks are not followed: validatePath rejects any symlink that resolves
// outside the sandbox before descent, and the underlying walker uses
// PathWalkDir which does not traverse into symlinked directories.
//
//	for entry, err := range c.Fs().WalkSeq("./") {
//		if err != nil { break }
//		if !entry.IsDir { /* file */ }
//	}
func (m *Fs) WalkSeq(root string) Seq2[FsEntry, error] {
	return m.walkSeq(root, nil)
}

// WalkSeqSkip walks like WalkSeq but skips any directory whose basename
// appears in skipNames (e.g. "vendor", "node_modules", ".git"). Skipped
// directories are not descended into; their contents are never yielded.
// The walk root itself is never skipped, even if its basename matches.
//
//	for entry, err := range c.Fs().WalkSeqSkip("./", "vendor", "node_modules", ".git") {
//		if err != nil { break }
//		if !entry.IsDir { /* file */ }
//	}
func (m *Fs) WalkSeqSkip(root string, skipNames ...string) Seq2[FsEntry, error] {
	skip := make(map[string]struct{}, len(skipNames))
	for _, name := range skipNames {
		if name != "" {
			skip[name] = struct{}{}
		}
	}
	return m.walkSeq(root, skip)
}

// walkSeq is the shared implementation behind WalkSeq and WalkSeqSkip.
func (m *Fs) walkSeq(root string, skip map[string]struct{}) Seq2[FsEntry, error] {
	return func(yield func(FsEntry, error) bool) {
		vp := m.validatePath(root)
		if !vp.OK {
			err, _ := vp.Value.(error)
			if err == nil {
				err = E("fs.WalkSeq", "invalid walk root", nil)
			}
			yield(FsEntry{}, err)
			return
		}
		fullRoot, _ := vp.Value.(string)
		if fullRoot == "" {
			yield(FsEntry{}, E("fs.WalkSeq", "validatePath returned empty root", nil))
			return
		}
		stop := false
		werr := PathWalkDir(fullRoot, func(path string, d FsDirEntry, walkErr error) error {
			if stop {
				return PathSkipAll
			}
			if walkErr != nil {
				if !yield(FsEntry{}, walkErr) {
					stop = true
					return PathSkipAll
				}
				return nil
			}
			if d.IsDir() && skip != nil && path != fullRoot {
				if _, ok := skip[d.Name()]; ok {
					return PathSkipDir
				}
			}
			relResult := PathRel(fullRoot, path)
			rel := ""
			if !relResult.OK {
				rel = path
			} else {
				rel = relResult.Value.(string)
			}
			info, _ := d.Info()
			mode := fs.FileMode(0)
			if info != nil {
				mode = info.Mode()
			}
			entry := FsEntry{
				Path:  rel,
				Name:  d.Name(),
				IsDir: d.IsDir(),
				Mode:  mode,
			}
			if !yield(entry, nil) {
				stop = true
				return PathSkipAll
			}
			return nil
		})
		// Per-entry errors already flowed through the yield above; a
		// failed Result here is the walker itself breaking.
		if !werr.OK {
			Debug("fs.WalkSeq: walk aborted", "err", werr.Error())
		}
	}
}
