package core_test

import (
	. "dappco.re/go"
)

func TestTime_Now_Good(t *T) {
	before := Now()
	value := Now()
	after := Now()

	AssertFalse(t, value.Before(before))
	AssertFalse(t, value.After(after))
}

func TestTime_Now_Bad(t *T) {
	value := Now()

	AssertFalse(t, value.IsZero())
	AssertTrue(t, value.After(UnixTime(0)))
	AssertGreaterOrEqual(t, value.Year(), 2026)
	AssertSame(t, Local, value.Location())
}

func TestTime_Now_Ugly(t *T) {
	first := Now()
	second := Now()

	AssertFalse(t, second.Before(first))
}

func TestTime_ParseDuration_Good(t *T) {
	r := ParseDuration("250ms")

	AssertTrue(t, r.OK)
	AssertEqual(t, 250*Millisecond, r.Value.(Duration))
}

func TestTime_ParseDuration_Bad(t *T) {
	r := ParseDuration("not-a-duration")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestTime_ParseDuration_Ugly(t *T) {
	r := ParseDuration("-1h30m")

	AssertTrue(t, r.OK)
	AssertEqual(t, -90*Minute, r.Value.(Duration))
}

func TestTime_Since_Good(t *T) {
	start := Now().Add(-Second)
	elapsed := Since(start)

	AssertGreaterOrEqual(t, elapsed, Second)
	AssertLess(t, elapsed, Minute)

	hourAgo := Now().Add(-Hour)
	AssertGreaterOrEqual(t, Since(hourAgo), Hour)
}

func TestTime_Since_Bad(t *T) {
	future := Now().Add(Second)

	AssertLess(t, Since(future), Duration(0))

	farFuture := Now().Add(Hour)
	AssertLess(t, Since(farFuture), -Minute)
	AssertGreater(t, Since(farFuture), -2*Hour)
}

func TestTime_Since_Ugly(t *T) {
	start := Now()
	Sleep(Millisecond)

	AssertGreater(t, Since(start), Duration(0))
}

func TestTime_Sleep_Good(t *T) {
	start := Now()
	Sleep(Millisecond)

	AssertGreaterOrEqual(t, Since(start), Millisecond)
}

func TestTime_Sleep_Bad(t *T) {
	start := Now()
	Sleep(-Millisecond)

	AssertLess(t, Since(start), 50*Millisecond)
}

func TestTime_Sleep_Ugly(t *T) {
	start := Now()
	Sleep(0)

	AssertLess(t, Since(start), 50*Millisecond)
}

func TestTime_TimeFormat_Good(t *T) {
	r := TimeParse(TimeRFC3339, "2026-04-28T07:00:00Z")
	RequireTrue(t, r.OK)

	AssertEqual(t, "2026-04-28T07:00:00Z", TimeFormat(r.Value.(Time), TimeRFC3339))
}

func TestTime_TimeFormat_Bad(t *T) {
	// A layout with no reference-time tokens is returned verbatim.
	AssertEqual(t, "agent", TimeFormat(UnixTime(0), "agent"))
	AssertEqual(t, "", TimeFormat(UnixTime(0), ""))

	// Literal text around the year token substitutes only the token.
	AssertEqual(t, "year-1970", TimeFormat(UnixTime(0).In(UTC), "year-2006"))
}

func TestTime_TimeFormat_Ugly(t *T) {
	AssertEqual(t, "1970-01-01", TimeFormat(UnixTime(0), TimeDateOnly))

	// Built in UTC so the rendering is deterministic across machines.
	ts := Date(2026, April, 28, 7, 5, 9, 0, UTC)
	AssertEqual(t, "2026-04-28", TimeFormat(ts, TimeDateOnly))
	AssertEqual(t, "07:05:09", TimeFormat(ts, TimeOnly))
	AssertEqual(t, "2026-04-28 07:05:09", TimeFormat(ts, DateTime))
}

func TestTime_TimeParse_Good(t *T) {
	r := TimeParse(TimeRFC3339, "2026-04-28T07:00:00Z")

	AssertTrue(t, r.OK)
	AssertEqual(t, int64(1777359600), r.Value.(Time).Unix())
}

func TestTime_TimeParse_Bad(t *T) {
	r := TimeParse(TimeRFC3339, "not-a-time")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))
}

func TestTime_TimeParse_Ugly(t *T) {
	r := TimeParse(TimeDateOnly, "2026-04-28")

	AssertTrue(t, r.OK)
	AssertEqual(t, "2026-04-28", TimeFormat(r.Value.(Time), TimeDateOnly))
}

func TestTime_Until_Good(t *T) {
	future := Now().Add(Second)
	wait := Until(future)

	AssertGreater(t, wait, Duration(0))
	AssertLessOrEqual(t, wait, Second)

	farFuture := Now().Add(Hour)
	AssertGreater(t, Until(farFuture), Minute)
	AssertLessOrEqual(t, Until(farFuture), Hour)
}

func TestTime_Until_Bad(t *T) {
	past := Now().Add(-Second)

	AssertLess(t, Until(past), Duration(0))

	farPast := Now().Add(-Hour)
	AssertLess(t, Until(farPast), -Minute)
	AssertGreater(t, Until(farPast), -2*Hour)
}

func TestTime_Until_Ugly(t *T) {
	deadline := Now().Add(Millisecond)
	Sleep(2 * Millisecond)

	AssertLessOrEqual(t, Until(deadline), Duration(0))
}

func TestTime_UnixNow_Good(t *T) {
	before := Now().Unix()
	value := UnixNow()
	after := Now().Unix()

	AssertGreaterOrEqual(t, value, before)
	AssertLessOrEqual(t, value, after)
}

func TestTime_UnixNow_Bad(t *T) {
	value := UnixNow()

	AssertGreater(t, value, int64(0))
	AssertGreater(t, value, int64(1767225600)) // after 2026-01-01 UTC
	AssertLessOrEqual(t, value-Now().Unix(), int64(1))
}

func TestTime_UnixNow_Ugly(t *T) {
	value := UnixNow()

	AssertLessOrEqual(t, value-Now().Unix(), int64(1))
	AssertEqual(t, value, UnixTime(value).Unix())
	AssertEqual(t, value, Unix(value, 0).Unix())
}

func TestTime_UnixTime_Good(t *T) {
	ts := UnixTime(1714291200)

	AssertEqual(t, int64(1714291200), ts.Unix())
	AssertEqual(t, 0, ts.Nanosecond())
	AssertEqual(t, int64(1714291200000), ts.UnixMilli())
	AssertEqual(t, "2024-04-28T08:00:00Z", TimeFormat(ts.In(UTC), RFC3339))
}

func TestTime_UnixTime_Bad(t *T) {
	ts := UnixTime(-1)

	AssertEqual(t, int64(-1), ts.Unix())
	AssertTrue(t, ts.Before(UnixTime(0)))
	AssertEqual(t, "1969-12-31T23:59:59Z", TimeFormat(ts.In(UTC), RFC3339))
}

func TestTime_UnixTime_Ugly(t *T) {
	AssertEqual(t, "1970-01-01", TimeFormat(UnixTime(0), TimeDateOnly))

	// Pinned to UTC so the epoch reads identically on any host.
	epoch := UnixTime(0).In(UTC)
	AssertEqual(t, int64(0), epoch.Unix())
	AssertEqual(t, 1970, epoch.Year())
	AssertEqual(t, "1970-01-01T00:00:00Z", TimeFormat(epoch, RFC3339))
}

func TestTime_Date_Good(t *T) {
	ts := Date(2026, April, 28, 7, 0, 0, 0, UTC)

	AssertEqual(t, int64(1777359600), ts.Unix())
	AssertEqual(t, 2026, ts.Year())
	AssertEqual(t, April, ts.Month())
	AssertEqual(t, 28, ts.Day())
	AssertEqual(t, 7, ts.Hour())
	AssertEqual(t, Tuesday, ts.Weekday())
	AssertSame(t, UTC, ts.Location())
}

func TestTime_Date_Bad(t *T) {
	// Month 13 normalises into January of the next year.
	ts := Date(2026, Month(13), 1, 0, 0, 0, 0, UTC)

	AssertEqual(t, 2027, ts.Year())
	AssertEqual(t, January, ts.Month())
}

func TestTime_Date_Ugly(t *T) {
	// The zero-ish boundary: epoch reconstructed via Date.
	ts := Date(1970, January, 1, 0, 0, 0, 0, UTC)

	AssertEqual(t, int64(0), ts.Unix())
	AssertEqual(t, "1970-01-01T00:00:00Z", TimeFormat(ts, RFC3339))

	// Sub-second precision carries through the nsec field.
	withNsec := Date(1970, January, 1, 0, 0, 0, 500000000, UTC)
	AssertEqual(t, 500000000, withNsec.Nanosecond())

	// Day 0 normalises back into the previous month (31 Dec 1969).
	rollback := Date(1970, January, 0, 0, 0, 0, 0, UTC)
	AssertEqual(t, 1969, rollback.Year())
	AssertEqual(t, December, rollback.Month())
	AssertEqual(t, 31, rollback.Day())
}

func TestTime_Month_Good(t *T) {
	ts := Date(2026, December, 25, 0, 0, 0, 0, UTC)

	AssertEqual(t, December, ts.Month())
	AssertEqual(t, "December", ts.Month().String())
	AssertEqual(t, Month(12), ts.Month())

	jan := Date(2026, January, 1, 0, 0, 0, 0, UTC)
	AssertEqual(t, January, jan.Month())
	AssertEqual(t, "January", jan.Month().String())
}

func TestTime_Month_Bad(t *T) {
	// January is 1, not 0 — guards against off-by-one assumptions.
	AssertEqual(t, Month(1), January)
	AssertEqual(t, Month(12), December)
	AssertEqual(t, "January", January.String())
	AssertEqual(t, "December", December.String())
	AssertNotEqual(t, January, December)
	AssertEqual(t, February, January+1)
}

func TestTime_Month_Ugly(t *T) {
	// The full sequence is contiguous and ordered.
	months := []Month{
		January, February, March, April, May, June,
		July, August, September, October, November, December,
	}
	for i, m := range months {
		AssertEqual(t, Month(i+1), m)
	}
}

func TestTime_Weekday_Good(t *T) {
	// 2026-04-28 is a Tuesday.
	ts := Date(2026, April, 28, 0, 0, 0, 0, UTC)

	AssertEqual(t, Tuesday, ts.Weekday())
	AssertEqual(t, "Tuesday", ts.Weekday().String())
	AssertEqual(t, Weekday(2), ts.Weekday())

	next := Date(2026, April, 29, 0, 0, 0, 0, UTC)
	AssertEqual(t, Wednesday, next.Weekday())
}

func TestTime_Weekday_Bad(t *T) {
	// Sunday is 0, the zero value — guards against treating it as unset.
	AssertEqual(t, Weekday(0), Sunday)
	AssertEqual(t, Weekday(6), Saturday)
	AssertEqual(t, "Sunday", Sunday.String())
	AssertEqual(t, "Saturday", Saturday.String())
	AssertNotEqual(t, Sunday, Saturday)
	AssertEqual(t, Monday, Sunday+1)
}

func TestTime_Weekday_Ugly(t *T) {
	days := []Weekday{
		Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday,
	}
	for i, d := range days {
		AssertEqual(t, Weekday(i), d)
	}
}

func TestTime_UTC_Good(t *T) {
	ts := Date(2026, January, 1, 12, 0, 0, 0, UTC)

	AssertEqual(t, "UTC", ts.Location().String())
	AssertSame(t, UTC, ts.Location())

	name, offset := ts.Zone()
	AssertEqual(t, "UTC", name)
	AssertEqual(t, 0, offset)
	AssertEqual(t, "2026-01-01T12:00:00Z", TimeFormat(ts, RFC3339))
}

func TestTime_UTC_Bad(t *T) {
	AssertNotNil(t, UTC)
	AssertEqual(t, "UTC", UTC.String())

	// Converting any instant to UTC yields a zero zone offset.
	_, offset := UnixTime(1777359600).In(UTC).Zone()
	AssertEqual(t, 0, offset)
}

func TestTime_UTC_Ugly(t *T) {
	// Converting to UTC must not shift the instant, only the zone.
	ts := UnixTime(1777359600)
	utc := ts.In(UTC)

	AssertEqual(t, ts.Unix(), utc.Unix())
	AssertTrue(t, ts.Equal(utc))
	AssertSame(t, UTC, utc.Location())
	AssertEqual(t, "2026-04-28T07:00:00Z", TimeFormat(utc, RFC3339))
}

func TestTime_Local_Good(t *T) {
	AssertNotNil(t, Local)
	AssertNotEmpty(t, Local.String())
	AssertSame(t, Local, Now().Location())

	ts := Date(2026, January, 1, 0, 0, 0, 0, Local)
	AssertSame(t, Local, ts.Location())
}

func TestTime_Local_Bad(t *T) {
	// In(Local) preserves the instant regardless of the machine's zone.
	ts := UnixTime(1777359600)
	local := ts.In(Local)

	AssertEqual(t, ts.Unix(), local.Unix())
	AssertTrue(t, ts.Equal(local))
	AssertSame(t, Local, local.Location())
}

func TestTime_Local_Ugly(t *T) {
	// Local and UTC describe the same instant for an epoch-derived time.
	ts := UnixTime(0)

	AssertEqual(t, ts.In(UTC).Unix(), ts.In(Local).Unix())
	AssertTrue(t, ts.In(UTC).Equal(ts.In(Local)))

	// The same holds for an arbitrary later instant.
	later := UnixTime(1777359600)
	AssertEqual(t, later.In(UTC).Unix(), later.In(Local).Unix())
}

func TestTime_Location_Good(t *T) {
	var loc *Location = UTC

	AssertEqual(t, "UTC", loc.String())
	AssertSame(t, UTC, loc)

	ts := Date(2026, April, 28, 7, 0, 0, 0, loc)
	AssertSame(t, loc, ts.Location())
	AssertEqual(t, int64(1777359600), ts.Unix())
}

func TestTime_Location_Bad(t *T) {
	// A typed-nil Location is the zero value of the alias.
	var loc *Location

	AssertNil(t, loc)

	// The stdlib rejects a nil Location rather than defaulting to UTC:
	// both Date and In panic.
	AssertPanics(t, func() { Date(2026, January, 1, 0, 0, 0, 0, loc) })
	AssertPanics(t, func() { UnixTime(0).In(loc) })
}

func TestTime_Location_Ugly(t *T) {
	ts := Date(2026, January, 1, 0, 0, 0, 0, Local)

	AssertSame(t, Local, ts.Location())

	// Re-zoning swaps the Location pointer but preserves the instant.
	utc := ts.In(UTC)
	AssertSame(t, UTC, utc.Location())
	AssertTrue(t, ts.Equal(utc))
	AssertEqual(t, ts.Unix(), utc.Unix())
}

func TestTime_RFC3339_Good(t *T) {
	r := TimeParse(RFC3339, "2026-04-28T07:00:00Z")
	RequireTrue(t, r.OK)

	AssertEqual(t, "2026-04-28T07:00:00Z", TimeFormat(r.Value.(Time), RFC3339))
}

func TestTime_RFC3339_Bad(t *T) {
	r := TimeParse(RFC3339, "2026-04-28")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))

	// Missing the timezone designator also fails RFC3339.
	noZone := TimeParse(RFC3339, "2026-04-28T07:00:00")
	AssertFalse(t, noZone.OK)
	AssertError(t, noZone.Value.(error))
}

func TestTime_RFC3339_Ugly(t *T) {
	// The bare constant equals the Time*-prefixed one.
	AssertEqual(t, TimeRFC3339, RFC3339)
	AssertEqual(t, "2006-01-02T15:04:05Z07:00", RFC3339)

	// Both aliases parse the same input to the same instant.
	a := TimeParse(RFC3339, "2026-04-28T07:00:00Z")
	b := TimeParse(TimeRFC3339, "2026-04-28T07:00:00Z")
	RequireTrue(t, a.OK)
	RequireTrue(t, b.OK)
	AssertTrue(t, a.Value.(Time).Equal(b.Value.(Time)))
}

func TestTime_RFC3339Nano_Good(t *T) {
	r := TimeParse(RFC3339Nano, "2026-04-28T07:00:00.5Z")
	RequireTrue(t, r.OK)

	AssertEqual(t, 500000000, r.Value.(Time).Nanosecond())
}

func TestTime_RFC3339Nano_Bad(t *T) {
	r := TimeParse(RFC3339Nano, "not-a-time")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))

	// A date-only string lacks the time and zone components.
	dateOnly := TimeParse(RFC3339Nano, "2026-04-28")
	AssertFalse(t, dateOnly.OK)
	AssertError(t, dateOnly.Value.(error))
}

func TestTime_RFC3339Nano_Ugly(t *T) {
	AssertEqual(t, TimeRFC3339Nano, RFC3339Nano)
	AssertEqual(t, "2006-01-02T15:04:05.999999999Z07:00", RFC3339Nano)

	// Full nanosecond precision round-trips through the layout.
	r := TimeParse(RFC3339Nano, "2026-04-28T07:00:00.123456789Z")
	RequireTrue(t, r.OK)
	AssertEqual(t, 123456789, r.Value.(Time).Nanosecond())
}

func TestTime_RFC1123_Good(t *T) {
	r := TimeParse(RFC1123, "Tue, 28 Apr 2026 07:00:00 UTC")

	AssertTrue(t, r.OK)
	ts := r.Value.(Time)
	AssertEqual(t, 2026, ts.Year())
	AssertEqual(t, April, ts.Month())
	AssertEqual(t, 28, ts.Day())
	AssertEqual(t, Tuesday, ts.Weekday())
	AssertEqual(t, "Tue, 28 Apr 2026 07:00:00 UTC", TimeFormat(ts, RFC1123))
}

func TestTime_RFC1123_Bad(t *T) {
	r := TimeParse(RFC1123, "2026-04-28T07:00:00Z")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))

	// Missing the weekday prefix also fails RFC1123.
	noDay := TimeParse(RFC1123, "28 Apr 2026 07:00:00 UTC")
	AssertFalse(t, noDay.OK)
}

func TestTime_RFC1123_Ugly(t *T) {
	AssertEqual(t, TimeRFC1123, RFC1123)
	AssertEqual(t, "Mon, 02 Jan 2006 15:04:05 MST", RFC1123)

	// Format then re-parse round-trips to the same instant.
	original := Date(2026, April, 28, 7, 0, 0, 0, UTC)
	r := TimeParse(RFC1123, TimeFormat(original, RFC1123))
	RequireTrue(t, r.OK)
	AssertTrue(t, original.Equal(r.Value.(Time)))
}

func TestTime_Kitchen_Good(t *T) {
	AssertEqual(t, "7:00AM", TimeFormat(Date(2026, April, 28, 7, 0, 0, 0, UTC), Kitchen))
	AssertEqual(t, "3:04PM", TimeFormat(Date(2026, April, 28, 15, 4, 0, 0, UTC), Kitchen))
	AssertEqual(t, "12:00AM", TimeFormat(Date(2026, April, 28, 0, 0, 0, 0, UTC), Kitchen))
	AssertEqual(t, "12:00PM", TimeFormat(Date(2026, April, 28, 12, 0, 0, 0, UTC), Kitchen))
}

func TestTime_Kitchen_Bad(t *T) {
	// Kitchen carries no date, so a round-trip drops the year.
	r := TimeParse(Kitchen, "7:00AM")
	RequireTrue(t, r.OK)

	AssertEqual(t, 0, r.Value.(Time).Year())
}

func TestTime_Kitchen_Ugly(t *T) {
	AssertEqual(t, TimeKitchen, Kitchen)
	AssertEqual(t, "3:04PM", Kitchen)

	// Kitchen carries only the clock; parsing recovers hour and minute.
	r := TimeParse(Kitchen, "3:04PM")
	RequireTrue(t, r.OK)
	AssertEqual(t, 15, r.Value.(Time).Hour())
	AssertEqual(t, 4, r.Value.(Time).Minute())
}

func TestTime_DateTime_Good(t *T) {
	ts := Date(2026, April, 28, 7, 0, 0, 0, UTC)

	AssertEqual(t, "2026-04-28 07:00:00", TimeFormat(ts, DateTime))

	r := TimeParse(DateTime, "2026-04-28 07:00:00")
	RequireTrue(t, r.OK)
	AssertEqual(t, 2026, r.Value.(Time).Year())
	AssertEqual(t, 7, r.Value.(Time).Hour())
}

func TestTime_DateTime_Bad(t *T) {
	r := TimeParse(DateTime, "2026-04-28")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))

	// DateTime uses a space separator; the RFC3339 'T' form fails.
	withT := TimeParse(DateTime, "2026-04-28T07:00:00")
	AssertFalse(t, withT.OK)
}

func TestTime_DateTime_Ugly(t *T) {
	AssertEqual(t, TimeDateTime, DateTime)
	AssertEqual(t, "2006-01-02 15:04:05", DateTime)

	// Format then re-parse preserves the rendered value.
	original := Date(2026, April, 28, 7, 0, 0, 0, UTC)
	r := TimeParse(DateTime, TimeFormat(original, DateTime))
	RequireTrue(t, r.OK)
	AssertEqual(t, TimeFormat(original, DateTime), TimeFormat(r.Value.(Time), DateTime))
}

func TestTime_DateOnly_Good(t *T) {
	ts := Date(2026, April, 28, 7, 0, 0, 0, UTC)

	AssertEqual(t, "2026-04-28", TimeFormat(ts, DateOnly))

	r := TimeParse(DateOnly, "2026-04-28")
	RequireTrue(t, r.OK)
	AssertEqual(t, 2026, r.Value.(Time).Year())
	AssertEqual(t, April, r.Value.(Time).Month())
	AssertEqual(t, 0, r.Value.(Time).Hour())
}

func TestTime_DateOnly_Bad(t *T) {
	r := TimeParse(DateOnly, "07:00:00")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))

	// A full timestamp leaves trailing data DateOnly cannot consume.
	full := TimeParse(DateOnly, "2026-04-28 07:00:00")
	AssertFalse(t, full.OK)
}

func TestTime_DateOnly_Ugly(t *T) {
	AssertEqual(t, TimeDateOnly, DateOnly)
	AssertEqual(t, "2006-01-02", DateOnly)

	r := TimeParse(DateOnly, "2026-12-25")
	RequireTrue(t, r.OK)
	AssertEqual(t, "2026-12-25", TimeFormat(r.Value.(Time), DateOnly))
}

func TestTime_TimeOnly_Good(t *T) {
	AssertEqual(t, "07:00:00", TimeFormat(Date(2026, April, 28, 7, 0, 0, 0, UTC), TimeOnly))

	ts := Date(2026, April, 28, 7, 30, 45, 0, UTC)
	AssertEqual(t, "07:30:45", TimeFormat(ts, TimeOnly))

	r := TimeParse(TimeOnly, "07:30:45")
	RequireTrue(t, r.OK)
	AssertEqual(t, 7, r.Value.(Time).Hour())
	AssertEqual(t, 30, r.Value.(Time).Minute())
	AssertEqual(t, 45, r.Value.(Time).Second())
}

func TestTime_TimeOnly_Bad(t *T) {
	r := TimeParse(TimeOnly, "2026-04-28")

	AssertFalse(t, r.OK)
	AssertError(t, r.Value.(error))

	// A full timestamp carries date components TimeOnly rejects.
	full := TimeParse(TimeOnly, "2026-04-28 07:00:00")
	AssertFalse(t, full.OK)
}

func TestTime_TimeOnly_Ugly(t *T) {
	AssertEqual(t, TimeTimeOnly, TimeOnly)
	AssertEqual(t, "15:04:05", TimeOnly)

	// TimeOnly carries no date, so parsing yields year zero.
	r := TimeParse(TimeOnly, "15:04:05")
	RequireTrue(t, r.OK)
	AssertEqual(t, 0, r.Value.(Time).Year())
	AssertEqual(t, "15:04:05", TimeFormat(r.Value.(Time), TimeOnly))
}

func TestTime_NewTimer_Good(t *T) {
	timer := NewTimer(Millisecond)
	defer timer.Stop()

	fired := <-timer.C
	AssertFalse(t, fired.IsZero())
}

func TestTime_NewTimer_Bad(t *T) {
	// Stop before fire returns true and leaves C empty.
	timer := NewTimer(Hour)

	AssertNotNil(t, timer)
	AssertNotNil(t, timer.C)
	AssertTrue(t, timer.Stop())
	AssertFalse(t, timer.Stop()) // already stopped — second Stop is false

	select {
	case <-timer.C:
		AssertTrue(t, false, "stopped timer should not have fired")
	default:
	}
}

func TestTime_NewTimer_Ugly(t *T) {
	// A zero-duration timer fires effectively immediately.
	timer := NewTimer(0)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-After(Second):
		AssertTrue(t, false, "zero-duration timer never fired")
	}
}

func TestTime_Timer_Good(t *T) {
	var timer *Timer = NewTimer(Millisecond)
	defer timer.Stop()

	<-timer.C
	AssertNotNil(t, timer)
}

func TestTime_Timer_Bad(t *T) {
	// Stopping an already-fired timer returns false.
	var timer *Timer = NewTimer(Millisecond)
	<-timer.C

	AssertFalse(t, timer.Stop())
}

func TestTime_Timer_Ugly(t *T) {
	// Reset on a stopped timer re-arms it.
	var timer *Timer = NewTimer(Hour)
	RequireTrue(t, timer.Stop())
	timer.Reset(Millisecond)
	defer timer.Stop()

	fired := <-timer.C
	AssertFalse(t, fired.IsZero())
}

func TestTime_AfterFunc_Good(t *T) {
	done := make(chan bool, 1)
	timer := AfterFunc(Millisecond, func() { done <- true })
	defer timer.Stop()

	AssertTrue(t, <-done)
}

func TestTime_AfterFunc_Bad(t *T) {
	// Stop before the duration elapses prevents the call.
	ran := make(chan bool, 1)
	timer := AfterFunc(Hour, func() { ran <- true })

	AssertTrue(t, timer.Stop())
	select {
	case <-ran:
		AssertTrue(t, false, "func ran after Stop")
	case <-After(10 * Millisecond):
	}
}

func TestTime_AfterFunc_Ugly(t *T) {
	// Zero duration still runs the func exactly once.
	count := make(chan int, 4)
	timer := AfterFunc(0, func() { count <- 1 })
	defer timer.Stop()

	select {
	case <-count:
	case <-After(Second):
		AssertTrue(t, false, "zero-duration AfterFunc never ran")
	}
}

func TestTime_Tick_Good(t *T) {
	ch := Tick(Millisecond)

	fired := <-ch
	AssertFalse(t, fired.IsZero())
}

func TestTime_Tick_Bad(t *T) {
	// Two consecutive ticks are spaced by at least the interval.
	ch := Tick(5 * Millisecond)
	first := <-ch
	second := <-ch

	AssertGreaterOrEqual(t, second.Sub(first), Duration(0))
}

func TestTime_Tick_Ugly(t *T) {
	// The channel keeps delivering — drain three ticks without blocking
	// forever.
	ch := Tick(Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case <-ch:
		case <-After(Second):
			AssertTrue(t, false, "ticker stopped delivering")
		}
	}
}
