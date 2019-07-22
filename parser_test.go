package douparser

import (
	"testing"
	"time"
)

func TestCombineDouTimeAndDate(t *testing.T) {
	f := func(sd, ed douDate, st, et douTime, expectedStart, expectedEnd time.Time) {
		t.Helper()
		receivedStart, receivedEnd := combineDouTimeAndDate(st, et, sd, ed)

		if !receivedStart.Equal(expectedStart) {
			t.Errorf(`For start %s %s
	want: %s
	got: %s`, sd, st, expectedStart, receivedStart)
		}

		if !receivedEnd.Equal(expectedEnd) {
			t.Errorf(`For end %s %s
	want: %s
	got: %s`, ed, et, expectedEnd, receivedEnd)
		}
	}

	f(douDate{}, douDate{}, douTime{}, douTime{}, time.Time{}, time.Time{})

	f(douDate{13, time.September, 2019},
		douDate{},
		douTime{},
		douTime{},
		time.Date(2019, time.September, 13, 0, 0, 0, 0, locale),
		time.Time{})

	f(douDate{13, time.September, 2019},
		douDate{16, time.September, 2019},
		douTime{},
		douTime{},
		time.Date(2019, time.September, 13, 0, 0, 0, 0, locale),
		time.Date(2019, time.September, 16, 0, 0, 0, 0, locale))

	f(douDate{13, time.September, 2019},
		douDate{},
		douTime{18, 30},
		douTime{},
		time.Date(2019, time.September, 13, 18, 30, 0, 0, locale),
		time.Time{})

	f(douDate{13, time.September, 2019},
		douDate{},
		douTime{18, 30},
		douTime{21, 40},
		time.Date(2019, time.September, 13, 18, 30, 0, 0, locale),
		time.Date(2019, time.September, 13, 21, 40, 0, 0, locale))

	f(douDate{13, time.September, 2019},
		douDate{16, time.September, 2019},
		douTime{18, 30},
		douTime{21, 40},
		time.Date(2019, time.September, 13, 18, 30, 0, 0, locale),
		time.Date(2019, time.September, 16, 18, 30, 0, 0, locale))
}
