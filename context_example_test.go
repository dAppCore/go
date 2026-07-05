package core_test

import (
	. "dappco.re/go"
)

// ExampleBackground creates a background context through `Background` for request lifetime
// control. Cancellation and timeout lifetimes are created through the core context
// surface.
func ExampleBackground() {
	ctx := Background()
	Println(ctx.Err() == nil)
	// Output: true
}

// ExampleWithTimeout creates a timeout context through `WithTimeout` for request lifetime
// control. Cancellation and timeout lifetimes are created through the core context
// surface.
func ExampleWithTimeout() {
	ctx, cancel := WithTimeout(Background(), 50*Millisecond)
	defer cancel()
	Println(ctx.Err() == nil)
	// Output: true
}

// ExampleWithCancel creates a cancellable context through `WithCancel` for request
// lifetime control. Cancellation and timeout lifetimes are created through the core
// context surface.
func ExampleWithCancel() {
	ctx, cancel := WithCancel(Background())
	cancel()
	Println(ctx.Err() != nil)
	// Output: true
}

// ExampleTODO returns a non-nil placeholder context through `TODO`.
func ExampleTODO() {
	Println(TODO() != nil)
	// Output: true
}

// ExampleWithValue carries a request-scoped value through `WithValue`.
func ExampleWithValue() {
	ctx := WithValue(Background(), "agent", "codex")
	Println(ctx.Value("agent"))
	// Output: codex
}

// ExampleWithDeadline derives a context that cancels at a deadline through `WithDeadline`.
func ExampleWithDeadline() {
	ctx, cancel := WithDeadline(Background(), Now().Add(Hour))
	defer cancel()
	_ = ctx // cancels automatically once the deadline passes
}
