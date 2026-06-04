// SPDX-License-Identifier: EUPL-1.2

// Time helpers for the Core framework.

package core

import "time"

// Now returns the current local time.
//
//	started := core.Now()
func Now() time.Time {
	return time.Now()
}

// UnixNow returns the current Unix timestamp in seconds.
//
//	ts := core.UnixNow()
func UnixNow() int64 {
	return Now().Unix()
}

// Sleep pauses the current goroutine for at least d.
//
//	core.Sleep(50 * time.Millisecond)
func Sleep(d time.Duration) {
	time.Sleep(d)
}

// Since returns the time elapsed since t.
//
//	elapsed := core.Since(started)
func Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Until returns the duration until t.
//
//	wait := core.Until(deadline)
func Until(t time.Time) time.Duration {
	return time.Until(t)
}

// ParseDuration parses a duration string and returns a Result containing
// time.Duration.
//
//	r := core.ParseDuration("250ms")
//	if r.OK { timeout := r.Value.(time.Duration) }
func ParseDuration(s string) Result {
	d, err := time.ParseDuration(s)
	if err != nil {
		return Result{err, false}
	}
	return Result{d, true}
}

// Duration is a time span — alias of time.Duration so consumers can
// pass timeouts and intervals without importing the time package.
//
//	timeout := 5 * core.Second
//	ctx, cancel := core.WithTimeout(core.Background(), timeout)
type Duration = time.Duration

// Time is a moment — alias of time.Time so consumers can take
// timestamps without importing the time package.
//
//	deadline := core.Now().Add(2 * core.Minute)
type Time = time.Time

// Common Duration units. Multiply by an integer to build a Duration.
//
//	timeout := 5 * core.Second
//	pause   := 250 * core.Millisecond
const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
)

// Common time format constants. Layouts compatible with time.Format.
//
//	stamp := core.TimeFormat(core.Now(), core.TimeRFC3339)
const (
	TimeRFC3339     = time.RFC3339
	TimeRFC3339Nano = time.RFC3339Nano
	TimeRFC1123     = time.RFC1123
	TimeRFC822      = time.RFC822
	TimeKitchen     = time.Kitchen  // "3:04PM"
	TimeStamp       = time.Stamp    // "Jan _2 15:04:05"
	TimeDateTime    = time.DateTime // "2006-01-02 15:04:05"
	TimeDateOnly    = time.DateOnly // "2006-01-02"
	TimeTimeOnly    = time.TimeOnly // "15:04:05"
)

// Bare-name layout constants matching the stdlib spelling, for callers
// that mirror time.Format usage directly. Equivalent to the Time*-prefixed
// set above; both are kept so existing TimeRFC3339 references and new
// RFC3339 references resolve to the same layout.
//
//	s := core.TimeFormat(core.Now(), core.RFC3339)
const (
	RFC3339     = time.RFC3339
	RFC3339Nano = time.RFC3339Nano
	RFC1123     = time.RFC1123
	Kitchen     = time.Kitchen  // "3:04PM"
	DateTime    = time.DateTime // "2006-01-02 15:04:05"
	DateOnly    = time.DateOnly // "2006-01-02"
	TimeOnly    = time.TimeOnly // "15:04:05"
)

// TimeFormat formats t as a string using the given layout. Layout
// constants are exported as TimeRFC3339, TimeDateTime, etc.
//
//	s := core.TimeFormat(core.Now(), core.TimeRFC3339)
func TimeFormat(t time.Time, layout string) string {
	return t.Format(layout)
}

// TimeParse parses value into a time.Time using the given layout.
// Returns Result wrapping time.Time on success or the parse error.
//
//	r := core.TimeParse(core.TimeRFC3339, "2026-04-28T07:00:00Z")
//	if r.OK { ts := r.Value.(time.Time) }
func TimeParse(layout, value string) Result {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Result{err, false}
	}
	return Result{t, true}
}

// UnixTime returns the time corresponding to the given Unix timestamp
// (seconds since 1970-01-01 UTC).
//
//	ts := core.UnixTime(1714291200)
func UnixTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// Unix returns the Time at sec seconds + nsec nanoseconds since the Unix
// epoch. Mirrors the stdlib signature for callers that need sub-second
// precision; UnixTime is the sec-only shorthand.
//
//	ts := core.Unix(sec, nsec)
func Unix(sec, nsec int64) Time {
	return time.Unix(sec, nsec)
}

// UnixMilli returns the Time at the given milliseconds since the Unix
// epoch. Useful for parsing JSON timestamps and other ms-resolution APIs.
//
//	ts := core.UnixMilli(jsonField)
func UnixMilli(ms int64) Time {
	return time.UnixMilli(ms)
}

// After returns a channel that delivers the current time after duration d.
// Use in select for timeouts:
//
//	select {
//	case msg := <-ch:
//	    handle(msg)
//	case <-core.After(2 * core.Second):
//	    return core.E("timeout", "no message", nil)
//	}
func After(d Duration) <-chan Time {
	return time.After(d)
}

// Ticker delivers Time values at regular intervals on its C channel.
// Stop the ticker with Stop() to release resources.
type Ticker = time.Ticker

// NewTicker returns a new Ticker that fires every duration d on its C
// channel. The caller MUST call Stop() to release the underlying timer.
//
//	ticker := core.NewTicker(30 * core.Second)
//	defer ticker.Stop()
//	for range ticker.C {
//	    poll()
//	}
func NewTicker(d Duration) *Ticker {
	return time.NewTicker(d)
}

// Tick is a convenience wrapper for NewTicker that returns only the
// channel. The underlying Ticker is never recovered, so Tick leaks and
// is only safe for tickers that live for the lifetime of the program.
// Prefer NewTicker + Stop for anything shorter-lived.
//
//	for range core.Tick(time.Minute) { poll() }
func Tick(d Duration) <-chan Time {
	return time.Tick(d)
}

// Timer fires once on its C channel after a duration elapses.
// Stop the timer with Stop() to release resources if it has not fired.
type Timer = time.Timer

// NewTimer returns a new Timer that sends the current time on its C
// channel after at least duration d.
//
//	timer := core.NewTimer(5 * core.Second)
//	defer timer.Stop()
//	<-timer.C
func NewTimer(d Duration) *Timer {
	return time.NewTimer(d)
}

// AfterFunc waits for the duration to elapse and then calls f in its own
// goroutine. It returns a Timer that can be used to cancel the call with
// its Stop method.
//
//	timer := core.AfterFunc(2*core.Second, func() { cleanup() })
//	defer timer.Stop()
func AfterFunc(d Duration, f func()) *Timer {
	return time.AfterFunc(d, f)
}

// Location maps instants to the zone in use at that time — alias of
// time.Location so consumers can pass zones without importing the time
// package. Use the UTC and Local package variables for the common cases.
//
//	ts := core.Date(2026, core.January, 1, 0, 0, 0, 0, core.UTC)
type Location = time.Location

// UTC is the Coordinated Universal Time zone.
//
//	stamp := core.Now().In(core.UTC)
var UTC = time.UTC

// Local is the system's local time zone.
//
//	stamp := core.Now().In(core.Local)
var Local = time.Local

// Month specifies a month of the year (January = 1, ...) — alias of
// time.Month so callers can build dates without importing the time
// package.
//
//	m := core.January
type Month = time.Month

// Months of the year, for use with Date and time formatting.
//
//	ts := core.Date(2026, core.December, 25, 0, 0, 0, 0, core.UTC)
const (
	January   = time.January
	February  = time.February
	March     = time.March
	April     = time.April
	May       = time.May
	June      = time.June
	July      = time.July
	August    = time.August
	September = time.September
	October   = time.October
	November  = time.November
	December  = time.December
)

// Weekday specifies a day of the week (Sunday = 0, ...) — alias of
// time.Weekday so callers can switch on Now().Weekday() without
// importing the time package.
//
//	if core.Now().Weekday() == core.Sunday { rest() }
type Weekday = time.Weekday

// Days of the week, for use with Time.Weekday comparisons.
//
//	weekend := day == core.Saturday || day == core.Sunday
const (
	Sunday    = time.Sunday
	Monday    = time.Monday
	Tuesday   = time.Tuesday
	Wednesday = time.Wednesday
	Thursday  = time.Thursday
	Friday    = time.Friday
	Saturday  = time.Saturday
)

// Date returns the Time corresponding to the given calendar fields in the
// given Location. Out-of-range values are normalised (e.g. month 13
// rolls into the next year), mirroring the stdlib.
//
//	ts := core.Date(2026, core.April, 28, 7, 0, 0, 0, core.UTC)
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time {
	return time.Date(year, month, day, hour, min, sec, nsec, loc)
}
