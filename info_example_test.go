package core_test

import (
	. "dappco.re/go"
)

// ExampleSysInfo declares runtime system information through `SysInfo` for runtime
// diagnostics. Runtime and environment facts are exposed as stable diagnostic values.
func ExampleSysInfo() {
	var info *SysInfo
	_ = info
}

// ExampleEnv reads environment access through `Env` for runtime diagnostics. Runtime and
// environment facts are exposed as stable diagnostic values.
func ExampleEnv() {
	Println(Env("OS"))   // e.g. "darwin"
	Println(Env("ARCH")) // e.g. "arm64"
}

// ExampleEnvKeys lists environment keys through `EnvKeys` for runtime diagnostics. Runtime
// and environment facts are exposed as stable diagnostic values.
func ExampleEnvKeys() {
	keys := EnvKeys()
	Println(len(keys) > 0)
	// Output: true
}

// ExampleOS returns the host operating system through `OS` (e.g. "darwin", "linux").
func ExampleOS() {
	Println(OS())
}

// ExampleArch returns the host CPU architecture through `Arch` (e.g. "arm64", "amd64").
func ExampleArch() {
	Println(Arch())
}

// ExampleGoVersion returns the Go runtime version through `GoVersion`.
func ExampleGoVersion() {
	Println(GoVersion())
}

// ExampleNumCPU returns the number of logical CPUs through `NumCPU`.
func ExampleNumCPU() {
	Println(NumCPU() > 0)
	// Output: true
}

// ExampleStackBuf captures the current goroutine stack through `StackBuf`.
func ExampleStackBuf() {
	Println(len(StackBuf()) > 0)
	// Output: true
}
