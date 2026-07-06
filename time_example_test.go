package core_test

import . "dappco.re/go"

// ExampleNow reads the current time through `Now` for health-check timing. Durations,
// parsing, and timestamps use core time wrappers for service code.
func ExampleNow() {
	Println(!Now().IsZero())
	// Output: true
}

// ExampleUnixNow reads the current Unix timestamp through `UnixNow` for health-check
// timing. Durations, parsing, and timestamps use core time wrappers for service code.
func ExampleUnixNow() {
	Println(UnixNow() > 0)
	// Output: true
}

// ExampleSleep pauses execution through `Sleep` for health-check timing. Durations,
// parsing, and timestamps use core time wrappers for service code.
func ExampleSleep() {
	Sleep(0)
	Println("awake")
	// Output: awake
}

// ExampleSince measures elapsed time through `Since` for health-check timing. Durations,
// parsing, and timestamps use core time wrappers for service code.
func ExampleSince() {
	Println(Since(UnixTime(0)) > 0)
	// Output: true
}

// ExampleUntil measures time until a deadline through `Until` for health-check timing.
// Durations, parsing, and timestamps use core time wrappers for service code.
func ExampleUntil() {
	Println(Until(UnixTime(32503680000)) > 0)
	// Output: true
}

// ExampleParseDuration parses duration text through `ParseDuration` for health-check
// timing. Durations, parsing, and timestamps use core time wrappers for service code.
func ExampleParseDuration() {
	r := ParseDuration("250ms")
	Println(r.Value)
	// Output: 250ms
}

// ExampleTimeFormat formats a timestamp through `TimeFormat` for health-check timing.
// Durations, parsing, and timestamps use core time wrappers for service code.
func ExampleTimeFormat() {
	t := UnixTime(1714262400)
	Println(TimeFormat(t, TimeDateOnly))
	// Output: 2024-04-28
}

// ExampleTimeParse parses a timestamp through `TimeParse` for health-check timing.
// Durations, parsing, and timestamps use core time wrappers for service code.
func ExampleTimeParse() {
	r := TimeParse(TimeRFC3339, "2026-04-28T07:00:00Z")
	Println(Contains(Sprint(r.Value), "2026-04-28 07:00:00"))
	// Output: true
}

// ExampleUnixTime builds a timestamp from seconds through `UnixTime` for health-check
// timing. Durations, parsing, and timestamps use core time wrappers for service code.
func ExampleUnixTime() {
	Println(Contains(Sprint(UnixTime(0)), "1970-01-01"))
	// Output: true
}

// ExampleUnix builds a timestamp from seconds + nanoseconds through `Unix` for health-check
// timing. Mirrors the stdlib two-arg signature; UnixTime is the sec-only shorthand.
func ExampleUnix() {
	Println(Contains(Sprint(Unix(0, 0)), "1970-01-01"))
	// Output: true
}

// ExampleUnixMilli builds a timestamp from milliseconds through `UnixMilli` for parsing
// JSON timestamps and other ms-resolution APIs without importing time directly.
func ExampleUnixMilli() {
	Println(Contains(Sprint(UnixMilli(0)), "1970-01-01"))
	// Output: true
}

// ExampleAfter shows the channel-returning timer used in select for
// timeouts. Pair with a body select-case to bound any blocking read.
func ExampleAfter() {
	ch := make(chan string, 1)
	ch <- "msg"
	select {
	case msg := <-ch:
		Println(msg)
	case <-After(100 * Millisecond):
		Println("timeout")
	}
	// Output: msg
}

// ExampleNewTicker fires periodic ticks for poll loops. Caller MUST Stop
// the ticker to release the underlying timer; the example reads one
// tick + stops to keep the example test bounded.
func ExampleNewTicker() {
	ticker := NewTicker(10 * Millisecond)
	defer ticker.Stop()
	<-ticker.C
	Println("tick")
	// Output: tick
}

// ExampleDate constructs a Time from calendar fields through `Date`.
func ExampleDate() {
	t := Date(2024, January, 15, 10, 30, 0, 0, UTC)
	Println(t.Year(), t.Month(), t.Day())
	// Output: 2024 January 15
}

// ExampleTick delivers periodic ticks on a channel through `Tick`.
func ExampleTick() {
	ch := Tick(Hour)
	_ = ch // a long interval; the channel fires once per period
}

// ExampleNewTimer fires once after a delay through `NewTimer`.
func ExampleNewTimer() {
	timer := NewTimer(Hour)
	defer timer.Stop()
}

// ExampleAfterFunc runs a function after a delay through `AfterFunc`.
func ExampleAfterFunc() {
	timer := AfterFunc(Hour, func() {})
	defer timer.Stop()
}
