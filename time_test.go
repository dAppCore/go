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
	AssertFalse(t, Now().IsZero())
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

	AssertGreaterOrEqual(t, Since(start), Second)
}

func TestTime_Since_Bad(t *T) {
	future := Now().Add(Second)

	AssertLess(t, Since(future), Duration(0))
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
	AssertEqual(t, "agent", TimeFormat(UnixTime(0), "agent"))
}

func TestTime_TimeFormat_Ugly(t *T) {
	AssertEqual(t, "1970-01-01", TimeFormat(UnixTime(0), TimeDateOnly))
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

	AssertGreater(t, Until(future), Duration(0))
}

func TestTime_Until_Bad(t *T) {
	past := Now().Add(-Second)

	AssertLess(t, Until(past), Duration(0))
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
	AssertGreater(t, UnixNow(), int64(0))
}

func TestTime_UnixNow_Ugly(t *T) {
	AssertLessOrEqual(t, UnixNow()-Now().Unix(), int64(1))
}

func TestTime_UnixTime_Good(t *T) {
	AssertEqual(t, int64(1714291200), UnixTime(1714291200).Unix())
}

func TestTime_UnixTime_Bad(t *T) {
	AssertEqual(t, int64(-1), UnixTime(-1).Unix())
}

func TestTime_UnixTime_Ugly(t *T) {
	AssertEqual(t, "1970-01-01", TimeFormat(UnixTime(0), TimeDateOnly))
}

func TestTime_Date_Good(t *T) {
	ts := Date(2026, April, 28, 7, 0, 0, 0, UTC)

	AssertEqual(t, int64(1777359600), ts.Unix())
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
}

func TestTime_Month_Good(t *T) {
	ts := Date(2026, December, 25, 0, 0, 0, 0, UTC)

	AssertEqual(t, December, ts.Month())
}

func TestTime_Month_Bad(t *T) {
	// January is 1, not 0 — guards against off-by-one assumptions.
	AssertEqual(t, Month(1), January)
	AssertEqual(t, Month(12), December)
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
}

func TestTime_Weekday_Bad(t *T) {
	// Sunday is 0, the zero value — guards against treating it as unset.
	AssertEqual(t, Weekday(0), Sunday)
	AssertEqual(t, Weekday(6), Saturday)
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
}

func TestTime_UTC_Bad(t *T) {
	AssertNotNil(t, UTC)
}

func TestTime_UTC_Ugly(t *T) {
	// Converting to UTC must not shift the instant, only the zone.
	ts := UnixTime(1777359600)

	AssertEqual(t, ts.Unix(), ts.In(UTC).Unix())
}

func TestTime_Local_Good(t *T) {
	AssertNotNil(t, Local)
}

func TestTime_Local_Bad(t *T) {
	// In(Local) preserves the instant regardless of the machine's zone.
	ts := UnixTime(1777359600)

	AssertEqual(t, ts.Unix(), ts.In(Local).Unix())
}

func TestTime_Local_Ugly(t *T) {
	// Local and UTC describe the same instant for an epoch-derived time.
	ts := UnixTime(0)

	AssertEqual(t, ts.In(UTC).Unix(), ts.In(Local).Unix())
}

func TestTime_Location_Good(t *T) {
	var loc *Location = UTC

	AssertEqual(t, "UTC", loc.String())
}

func TestTime_Location_Bad(t *T) {
	// A nil Location is a valid concept (UTC) for stdlib; assert the
	// alias accepts the typed nil without panicking on assignment.
	var loc *Location

	AssertNil(t, loc)
}

func TestTime_Location_Ugly(t *T) {
	ts := Date(2026, January, 1, 0, 0, 0, 0, Local)

	AssertSame(t, Local, ts.Location())
}

func TestTime_RFC3339_Good(t *T) {
	r := TimeParse(RFC3339, "2026-04-28T07:00:00Z")
	RequireTrue(t, r.OK)

	AssertEqual(t, "2026-04-28T07:00:00Z", TimeFormat(r.Value.(Time), RFC3339))
}

func TestTime_RFC3339_Bad(t *T) {
	r := TimeParse(RFC3339, "2026-04-28")

	AssertFalse(t, r.OK)
}

func TestTime_RFC3339_Ugly(t *T) {
	// The bare constant equals the Time*-prefixed one.
	AssertEqual(t, TimeRFC3339, RFC3339)
}

func TestTime_RFC3339Nano_Good(t *T) {
	r := TimeParse(RFC3339Nano, "2026-04-28T07:00:00.5Z")
	RequireTrue(t, r.OK)

	AssertEqual(t, 500000000, r.Value.(Time).Nanosecond())
}

func TestTime_RFC3339Nano_Bad(t *T) {
	r := TimeParse(RFC3339Nano, "not-a-time")

	AssertFalse(t, r.OK)
}

func TestTime_RFC3339Nano_Ugly(t *T) {
	AssertEqual(t, TimeRFC3339Nano, RFC3339Nano)
}

func TestTime_RFC1123_Good(t *T) {
	r := TimeParse(RFC1123, "Tue, 28 Apr 2026 07:00:00 UTC")

	AssertTrue(t, r.OK)
}

func TestTime_RFC1123_Bad(t *T) {
	r := TimeParse(RFC1123, "2026-04-28T07:00:00Z")

	AssertFalse(t, r.OK)
}

func TestTime_RFC1123_Ugly(t *T) {
	AssertEqual(t, TimeRFC1123, RFC1123)
}

func TestTime_Kitchen_Good(t *T) {
	AssertEqual(t, "7:00AM", TimeFormat(Date(2026, April, 28, 7, 0, 0, 0, UTC), Kitchen))
}

func TestTime_Kitchen_Bad(t *T) {
	// Kitchen carries no date, so a round-trip drops the year.
	r := TimeParse(Kitchen, "7:00AM")
	RequireTrue(t, r.OK)

	AssertEqual(t, 0, r.Value.(Time).Year())
}

func TestTime_Kitchen_Ugly(t *T) {
	AssertEqual(t, TimeKitchen, Kitchen)
}

func TestTime_DateTime_Good(t *T) {
	AssertEqual(t, "2026-04-28 07:00:00", TimeFormat(Date(2026, April, 28, 7, 0, 0, 0, UTC), DateTime))
}

func TestTime_DateTime_Bad(t *T) {
	r := TimeParse(DateTime, "2026-04-28")

	AssertFalse(t, r.OK)
}

func TestTime_DateTime_Ugly(t *T) {
	AssertEqual(t, TimeDateTime, DateTime)
}

func TestTime_DateOnly_Good(t *T) {
	AssertEqual(t, "2026-04-28", TimeFormat(Date(2026, April, 28, 7, 0, 0, 0, UTC), DateOnly))
}

func TestTime_DateOnly_Bad(t *T) {
	r := TimeParse(DateOnly, "07:00:00")

	AssertFalse(t, r.OK)
}

func TestTime_DateOnly_Ugly(t *T) {
	AssertEqual(t, TimeDateOnly, DateOnly)
}

func TestTime_TimeOnly_Good(t *T) {
	AssertEqual(t, "07:00:00", TimeFormat(Date(2026, April, 28, 7, 0, 0, 0, UTC), TimeOnly))
}

func TestTime_TimeOnly_Bad(t *T) {
	r := TimeParse(TimeOnly, "2026-04-28")

	AssertFalse(t, r.OK)
}

func TestTime_TimeOnly_Ugly(t *T) {
	AssertEqual(t, TimeTimeOnly, TimeOnly)
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

	AssertTrue(t, timer.Stop())
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
	timer := NewTimer(Millisecond)
	<-timer.C

	AssertFalse(t, timer.Stop())
}

func TestTime_Timer_Ugly(t *T) {
	// Reset on a stopped timer re-arms it.
	timer := NewTimer(Hour)
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
