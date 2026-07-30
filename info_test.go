// SPDX-License-Identifier: EUPL-1.2

package core_test

import (
	. "dappco.re/go"
)

func TestInfo_Env_OS_Good(t *T) {
	v := Env("OS")
	AssertNotEmpty(t, v)
	AssertContains(t, []string{"darwin", "linux", "windows"}, v)
}

func TestInfo_Env_ARCH_Good(t *T) {
	v := Env("ARCH")
	AssertNotEmpty(t, v)
	AssertContains(t, []string{"amd64", "arm64", "386"}, v)
}

func TestInfo_Env_GO_Good(t *T) {
	v := Env("GO")
	AssertNotEmpty(t, v)
	AssertTrue(t, HasPrefix(v, "go"))
	AssertEqual(t, GoVersion(), v)
	AssertNotContains(t, v, " ")
	AssertEqual(t, v, Env("GO"))
}

func TestInfo_Env_DS_Good(t *T) {
	ds := Env("DS")
	AssertContains(t, []string{"/", "\\"}, ds)
	AssertLen(t, ds, 1)
	AssertEqual(t, string(PathSeparator), ds)
}

func TestInfo_Env_PS_Good(t *T) {
	ps := Env("PS")
	AssertContains(t, []string{":", ";"}, ps)
	AssertLen(t, ps, 1)
	AssertEqual(t, string(PathListSeparator), ps)
	AssertNotEqual(t, Env("DS"), ps)
}

func TestInfo_Env_DIR_HOME_Good(t *T) {
	home := Env("DIR_HOME")
	AssertNotEmpty(t, home)
	AssertTrue(t, PathIsAbs(home), "DIR_HOME should be absolute")
}

func TestInfo_Env_DIR_DATA_Good(t *T) {
	dataDir := Env("DIR_DATA")
	AssertNotEmpty(t, dataDir)
	AssertTrue(t, PathIsAbs(dataDir), "DIR_DATA should be absolute")
}

func TestInfo_Env_DIR_TMP_Good(t *T) {
	tmp := Env("DIR_TMP")
	AssertNotEmpty(t, tmp)
	AssertEqual(t, TempDir(), tmp)
	AssertTrue(t, PathIsAbs(tmp), "DIR_TMP should be absolute")
}

func TestInfo_Env_DIR_CONFIG_Good(t *T) {
	cfg := Env("DIR_CONFIG")
	AssertNotEmpty(t, cfg)
	AssertTrue(t, PathIsAbs(cfg), "DIR_CONFIG should be absolute")
	if r := UserConfigDir(); r.OK {
		AssertEqual(t, r.Value.(string), cfg)
	}
}

func TestInfo_Env_DIR_CACHE_Good(t *T) {
	cache := Env("DIR_CACHE")
	AssertNotEmpty(t, cache)
	AssertTrue(t, PathIsAbs(cache), "DIR_CACHE should be absolute")
	if r := UserCacheDir(); r.OK {
		AssertEqual(t, r.Value.(string), cache)
	}
}

func TestInfo_Env_HOSTNAME_Good(t *T) {
	host := Env("HOSTNAME")
	AssertNotEmpty(t, host)
	AssertNotContains(t, host, " ")
	if r := Hostname(); r.OK {
		AssertEqual(t, r.Value.(string), host)
	}
}

func TestInfo_Env_USER_Good(t *T) {
	user := Env("USER")
	AssertNotEmpty(t, user)
	AssertEqual(t, Username(), user)
	AssertEqual(t, user, Env("USER"))
}

func TestInfo_Env_PID_Good(t *T) {
	pid := Env("PID")
	AssertNotEmpty(t, pid)
	AssertEqual(t, Itoa(Getpid()), pid)
	AssertGreater(t, Getpid(), 0)
}

func TestInfo_Env_NUM_CPU_Good(t *T) {
	n := Env("NUM_CPU")
	AssertNotEmpty(t, n)
	AssertEqual(t, Itoa(NumCPU()), n)
	AssertGreaterOrEqual(t, NumCPU(), 1)
}

func TestInfo_Env_CORE_START_Good(t *T) {
	ts := Env("CORE_START")
	RequireNotEmpty(t, ts)
	r := TimeParse(TimeRFC3339, ts)
	AssertTrue(t, r.OK, "CORE_START should be valid RFC3339")
}

func TestInfo_Env_Bad_Unknown(t *T) {
	AssertEqual(t, "", Env("NOPE"))
	AssertNotContains(t, EnvKeys(), "NOPE")
	AssertNotEmpty(t, Env("OS"))
}

func TestInfo_Env_Good_CoreInstance(t *T) {
	c := New()
	AssertEqual(t, Env("OS"), c.Env("OS"))
	AssertEqual(t, Env("DIR_HOME"), c.Env("DIR_HOME"))
}

func TestInfo_EnvKeys_Good(t *T) {
	keys := EnvKeys()
	AssertNotEmpty(t, keys)
	AssertContains(t, keys, "OS")
	AssertContains(t, keys, "DIR_HOME")
	AssertContains(t, keys, "CORE_START")
}

func TestInfo_Arch_Good(t *T) {
	arch := Arch()
	AssertNotEmpty(t, arch)
	AssertEqual(t, Env("ARCH"), arch)
	AssertEqual(t, Lower(arch), arch)
	AssertNotContains(t, arch, " ")
	AssertEqual(t, arch, Arch())
}

func TestInfo_Arch_Bad(t *T) {
	arch := Arch()
	AssertEqual(t, Env("ARCH"), arch)
	AssertNotEmpty(t, arch)
	AssertNotEqual(t, OS(), arch)
}

func TestInfo_Arch_Ugly(t *T) {
	arch := Arch()
	AssertEqual(t, Lower(arch), arch)
	AssertNotContains(t, arch, " ")
	AssertNotEmpty(t, arch)
}

func TestInfo_Env_Good(t *T) {
	// Env resolves init-populated keys to their non-empty values.
	AssertNotEmpty(t, Env("OS"))
	AssertNotEmpty(t, Env("DIR_HOME"))
	AssertEqual(t, OS(), Env("OS"))
}

func TestInfo_Env_Bad(t *T) {
	AssertEqual(t, "", Env("CORE_TEST_MISSING"))
	AssertNotContains(t, EnvKeys(), "CORE_TEST_MISSING")
	AssertNotEmpty(t, Env("DIR_HOME"))
}

func TestInfo_Env_Ugly(t *T) {
	t.Setenv("CORE_TEST_SESSION", "token")
	AssertEqual(t, "token", Env("CORE_TEST_SESSION"))
	AssertEqual(t, Getenv("CORE_TEST_SESSION"), Env("CORE_TEST_SESSION"))
	AssertNotContains(t, EnvKeys(), "CORE_TEST_SESSION")
}

func TestInfo_EnvKeys_Bad(t *T) {
	keys := EnvKeys()
	AssertNotContains(t, keys, "CORE_TEST_MISSING")
	AssertNotContains(t, keys, "NOPE")
	AssertNotEmpty(t, keys)
	AssertContains(t, keys, "OS")
}

func TestInfo_EnvKeys_Ugly(t *T) {
	keys := EnvKeys()
	RequireNotEmpty(t, keys)
	keys[0] = "MUTATED"

	AssertNotContains(t, EnvKeys(), "MUTATED")
}

func TestInfo_GoVersion_Good(t *T) {
	v := GoVersion()
	AssertTrue(t, HasPrefix(v, "go"))
	AssertNotEmpty(t, v)
	AssertEqual(t, Env("GO"), v)
	AssertEqual(t, v, GoVersion())
}

func TestInfo_GoVersion_Bad(t *T) {
	AssertEqual(t, Env("GO"), GoVersion())
	AssertNotEqual(t, "go", GoVersion())
	AssertNotEmpty(t, GoVersion())
}

func TestInfo_GoVersion_Ugly(t *T) {
	v := GoVersion()
	AssertNotContains(t, v, " ")
	AssertTrue(t, HasPrefix(v, "go"))
	AssertNotEmpty(t, v)
}

func TestInfo_NumCPU_Good(t *T) {
	n := NumCPU()
	AssertGreaterOrEqual(t, n, 1)
	AssertEqual(t, Itoa(n), Env("NUM_CPU"))
	AssertEqual(t, n, NumCPU())
}

func TestInfo_NumCPU_Bad(t *T) {
	AssertEqual(t, Env("NUM_CPU"), Itoa(NumCPU()))
	AssertGreater(t, NumCPU(), 0)
	AssertNotEqual(t, "0", Env("NUM_CPU"))
}

func TestInfo_NumCPU_Ugly(t *T) {
	AssertNotPanics(t, func() {
		_ = NumCPU()
	})
}

func TestInfo_OS_Good(t *T) {
	os := OS()
	AssertNotEmpty(t, os)
	AssertContains(t, []string{"linux", "darwin", "windows", "freebsd", "openbsd", "netbsd"}, os)
}

func TestInfo_OS_Bad(t *T) {
	AssertEqual(t, Env("OS"), OS())
	AssertNotEmpty(t, OS())
	AssertNotEqual(t, Arch(), OS())
}

func TestInfo_OS_Ugly(t *T) {
	os := OS()
	AssertEqual(t, Lower(os), os)
	AssertNotContains(t, os, " ")
	AssertNotEmpty(t, os)
}

func TestInfo_StackBuf_Good(t *T) {
	stack := string(StackBuf())
	AssertContains(t, stack, "goroutine")
	AssertContains(t, stack, "StackBuf")
	AssertContains(t, stack, "TestInfo_StackBuf_Good")
}

func TestInfo_StackBuf_Bad(t *T) {
	buf := StackBuf()
	AssertNotEmpty(t, buf)
	AssertContains(t, string(buf), "goroutine")
	AssertGreater(t, len(buf), 10)
}

func TestInfo_StackBuf_Ugly(t *T) {
	first := string(StackBuf())
	AssertNotEqual(t, "", first)
	AssertContains(t, first, "goroutine")
	AssertContains(t, string(StackBuf()), "TestInfo_StackBuf_Ugly")
}

// --- Pinner ---

// Pinner keeps a Go allocation at a fixed address across a garbage collection,
// so a pointer handed to C stays valid. Its observable contract is that the
// address does not move while pinned and that the value survives — which is
// what these check, since the pin itself is invisible from Go.

func TestInfo_Pinner_Good(t *T) {
	buf := make([]byte, 8)
	buf[0] = 42

	var p Pinner
	p.Pin(&buf[0])
	defer p.Unpin()

	AssertEqual(t, byte(42), buf[0])
}

// Bad: Unpin is what releases the pin, and pinning again after it must work —
// a Pinner that could only be used once would leak for every reused buffer.
func TestInfo_Pinner_Bad(t *T) {
	buf := make([]byte, 4)

	var p Pinner
	p.Pin(&buf[0])
	p.Unpin()

	p.Pin(&buf[0])
	p.Unpin()

	AssertEqual(t, 4, len(buf))
}

// Ugly: several allocations pinned through one Pinner are all released by a
// single Unpin — the shape PinnedView relies on.
func TestInfo_Pinner_Ugly(t *T) {
	a := make([]byte, 2)
	b := make([]int32, 2)
	a[1], b[1] = 9, 9

	var p Pinner
	p.Pin(&a[0])
	p.Pin(&b[0])
	p.Unpin()

	AssertEqual(t, byte(9), a[1])
	AssertEqual(t, int32(9), b[1])
}
