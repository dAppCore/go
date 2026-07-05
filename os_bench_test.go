// SPDX-License-Identifier: EUPL-1.2

// Benchmarks for the os.go wrappers — filesystem syscalls, process
// info, and environment access. The I/O benchmarks run against a
// per-benchmark temp dir; the env/info wrappers measure the syscall +
// Result-boxing overhead callers pay on every workspace setup.
//
// Exit is intentionally not benchmarked — it terminates the process.
//
// Run:    go test -bench='Benchmark(Stdin|Stdout|Mkdir|Lstat|Chmod|Symlink|Readlink|CreateTemp|Is(NotExist|Exist|Permission)|Args|Executable|Getppid|Chdir|Environ|Setenv|Unsetenv|LookupEnv)' -benchmem -run='^$' .

package core_test

import (
	. "dappco.re/go"
)

var (
	osSinkStr   string
	osSinkBool  bool
	osSinkInt   int
	osSinkSlice []string
	osSinkAny   any
)

func BenchmarkStdin(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkAny = Stdin()
	}
}

func BenchmarkStdout(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkAny = Stdout()
	}
}

func BenchmarkMkdir(b *B) {
	dir := b.TempDir()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Mkdir(PathJoin(dir, Itoa(i)), 0o755)
	}
}

func BenchmarkLstat(b *B) {
	dir := b.TempDir()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkAny = Lstat(dir)
	}
}

func BenchmarkChmod(b *B) {
	f := PathJoin(b.TempDir(), "f")
	WriteFile(f, []byte("x"), 0o644)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Chmod(f, 0o644)
	}
}

func BenchmarkSymlink(b *B) {
	dir := b.TempDir()
	target := PathJoin(dir, "target")
	WriteFile(target, []byte("x"), 0o644)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Symlink(target, PathJoin(dir, "link"+Itoa(i)))
	}
}

func BenchmarkReadlink(b *B) {
	dir := b.TempDir()
	target := PathJoin(dir, "target")
	link := PathJoin(dir, "link")
	WriteFile(target, []byte("x"), 0o644)
	Symlink(target, link)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkAny = Readlink(link)
	}
}

func BenchmarkCreateTemp(b *B) {
	dir := b.TempDir()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkAny = CreateTemp(dir, "bench-*")
	}
}

func BenchmarkIsNotExist(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkBool = IsNotExist(ErrNotExist)
	}
}

func BenchmarkIsExist(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkBool = IsExist(AnError)
	}
}

func BenchmarkIsPermission(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkBool = IsPermission(AnError)
	}
}

func BenchmarkArgs(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkSlice = Args()
	}
}

func BenchmarkExecutable(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkAny = Executable()
	}
}

func BenchmarkGetppid(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkInt = Getppid()
	}
}

func BenchmarkChdir(b *B) {
	cwd := MustCast[string](Getwd())
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Chdir(cwd)
	}
}

func BenchmarkEnviron(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkSlice = Environ()
	}
}

func BenchmarkSetenv(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Setenv("CORE_BENCH", "value")
	}
	Unsetenv("CORE_BENCH")
}

func BenchmarkUnsetenv(b *B) {
	Setenv("CORE_BENCH", "value")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Unsetenv("CORE_BENCH")
	}
}

func BenchmarkLookupEnv(b *B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		osSinkStr, osSinkBool = LookupEnv("PATH")
	}
}
