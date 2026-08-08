// SPDX-License-Identifier: EUPL-1.2

// Application identity for the Core framework.

package core

// App holds the application identity and optional GUI runtime.
//
//	app := core.App{}.New(core.NewOptions(
//	    core.Option{Key: "name", Value: "Core CLI"},
//	    core.Option{Key: "version", Value: "1.0.0"},
//	))
type App struct {
	Name        string
	Version     string
	Description string
	Filename    string
	Path        string
	Runtime     any // GUI runtime (e.g., Wails App). Nil for CLI-only.
}

// New creates an App from Options.
//
//	app := core.App{}.New(core.NewOptions(
//	    core.Option{Key: "name", Value: "myapp"},
//	    core.Option{Key: "version", Value: "1.0.0"},
//	))
func (a App) New(opts Options) App {
	if name := opts.String("name"); name != "" {
		a.Name = name
	}
	if version := opts.String("version"); version != "" {
		a.Version = version
	}
	if desc := opts.String("description"); desc != "" {
		a.Description = desc
	}
	if filename := opts.String("filename"); filename != "" {
		a.Filename = filename
	}
	return a
}

// Find locates a program on PATH and returns a Result containing the App.
// Uses core.Stat to search PATH directories — no os/exec dependency.
//
//	r := core.App{}.Find("node", "Node.js")
//	if r.OK { app := r.Value.(*App) }
func (a App) Find(filename, name string) Result {
	pathExt := Env("PATHEXT")
	if pathExt == "" && OS() == "windows" {
		// Go's own default when the variable is unset.
		pathExt = ".COM;.EXE;.BAT;.CMD"
	}
	return a.findWith(filename, name, Env("PATH"), pathExt)
}

// findWith is Find with the environment as arguments, so tests can drive
// the Windows resolution rules from any platform against a fixture PATH.
// An empty extension list means POSIX semantics; a non-empty one means
// Windows PATHEXT semantics — the list IS the platform switch.
func (a App) findWith(filename, name, pathEnv, pathExt string) Result {
	// A path rather than a bare name is checked directly. Go accepts "/"
	// as a separator on Windows too, so both spellings mean "path" —
	// matching only the platform separator sent "bin/tool" on a PATH
	// hunt there instead of probing it.
	if Contains(filename, string(PathSeparator)) || Contains(filename, "/") {
		abs := PathAbs(filename)
		if !abs.OK {
			return abs
		}
		for _, candidate := range executableCandidates(abs.Value.(string), pathExt) {
			if isExecutableWith(candidate, pathExt) {
				return Result{&App{Name: name, Filename: filename, Path: candidate}, true}
			}
		}
		return Result{E("app.Find", Concat(filename, " not found"), nil), false}
	}

	// Search PATH
	if pathEnv == "" {
		return Result{E("app.Find", "PATH is empty", nil), false}
	}
	for _, dir := range Split(pathEnv, string(PathListSeparator)) {
		for _, candidate := range executableCandidates(PathJoin(dir, filename), pathExt) {
			if !isExecutableWith(candidate, pathExt) {
				continue
			}
			abs := PathAbs(candidate)
			if !abs.OK {
				continue
			}
			return Result{&App{Name: name, Filename: filename, Path: abs.Value.(string)}, true}
		}
	}
	return Result{E("app.Find", Concat(filename, " not found on PATH"), nil), false}
}

// executableCandidates lists the filenames a program name may resolve to.
// With no extension list (POSIX) the path stands alone. With one (Windows
// PATHEXT) the name is tried as itself when it already carries a listed
// extension, then with each listed extension appended — "git" on disk is
// "git.exe", and the bare join alone would never find it.
func executableCandidates(path, pathExt string) []string {
	exts := splitPathExt(pathExt)
	if len(exts) == 0 {
		return []string{path}
	}
	out := make([]string, 0, len(exts)+1)
	if hasListedExt(path, exts) {
		out = append(out, path)
	}
	for _, ext := range exts {
		out = append(out, Concat(path, ext))
	}
	return out
}

// splitPathExt parses a PATHEXT value (";"-separated, case-insensitive)
// into normalised lower-case ".ext" entries. Empty input yields nil —
// the POSIX arm.
func splitPathExt(pathExt string) []string {
	if pathExt == "" {
		return nil
	}
	var out []string
	for _, e := range Split(pathExt, ";") {
		e = Lower(Trim(e))
		if e == "" {
			continue
		}
		if !HasPrefix(e, ".") {
			e = Concat(".", e)
		}
		out = append(out, e)
	}
	return out
}

// hasListedExt reports whether path's extension appears in exts
// (already lower-cased ".ext" entries).
func hasListedExt(path string, exts []string) bool {
	got := Lower(PathExt(path))
	if got == "" {
		return false
	}
	for _, e := range exts {
		if got == e {
			return true
		}
	}
	return false
}

// isExecutable checks if a path exists and is executable on this platform.
func isExecutable(path string) bool {
	pathExt := Env("PATHEXT")
	if pathExt == "" && OS() == "windows" {
		pathExt = ".COM;.EXE;.BAT;.CMD"
	}
	return isExecutableWith(path, pathExt)
}

// isExecutableWith is isExecutable with the extension list as an argument.
// POSIX (empty list) asks the mode bits. Windows has no execute bit —
// Stat synthesises 0666, or 0444 for a read-only file — so the old
// mode&0111 test rejected EVERY file on the platform, git.exe included,
// and Find could never succeed there. With a list present the answer is
// "a regular file whose extension is listed", and the mode is not asked.
func isExecutableWith(path, pathExt string) bool {
	r := Stat(path)
	if !r.OK {
		return false
	}
	info := r.Value.(interface {
		IsDir() bool
		Mode() FileMode
	})
	if info.IsDir() {
		return false
	}
	exts := splitPathExt(pathExt)
	if len(exts) == 0 {
		// Regular file with at least one execute bit
		return info.Mode()&0111 != 0
	}
	return hasListedExt(path, exts)
}
