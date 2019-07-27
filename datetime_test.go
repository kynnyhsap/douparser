package douparser

import (
	"testing"
	"time"
)

func isDouDateTimeEqual(a, b douDateTime) bool {
	return a.year == b.year &&
		a.month == b.month &&
		a.day == b.day &&
		a.hours == b.hours &&
		a.minutes == b.minutes
}

func TestDateDefined(t *testing.T) {
	f := func(dt douDateTime, defined bool) {
		got := dt.dateDefined()

		if got != defined {
			t.Errorf(`dateDefined() for %s failed
				expected: %t
				recieved: %t
			`, dt, defined, got)
		}
	}

	f(douDateTime{year: 2012, month: time.September, day: 13}, true)
	f(douDateTime{year: 2012, month: time.September, day: 0}, false)
	f(douDateTime{year: 0, month: time.September, day: 13}, false)
	f(douDateTime{year: 2012, month: 0, day: 13}, false)
}

func TestTimeDefined(t *testing.T) {
	f := func(dt douDateTime, defined bool) {
		got := dt.timeDefined()

		if got != defined {
			t.Errorf(`timeDefined() for %s failed
				expected: %t
				recieved: %t
			`, dt, defined, got)
		}
	}

	f(douDateTime{hours: 13, minutes: 41}, true)
	f(douDateTime{hours: 12, minutes: 0}, true)
	f(douDateTime{hours: 0, minutes: 0}, false)
	f(douDateTime{hours: 0, minutes: 32}, false)
	f(douDateTime{hours: -1, minutes: -2}, false)
}

func TestParseDateString(t *testing.T) {
	f := func(date string, year int, month time.Month, day int) {
		dt := douDateTime{}

		dt.parseDateString(date)

		if dt.year != year {
			t.Errorf(`parseDateString(%s) failed with year
				expected: %d
				recieved: %d
			`, date, year, dt.year)
		}

		if dt.month != month {
			t.Errorf(`parseDateString(%s) failed with month
				expected: %d
				recieved: %d
			`, date, month, dt.month)
		}

		if dt.day != day {
			t.Errorf(`parseDateString(%s) failed with day
				expected: %d
				recieved: %d
			`, date, day, dt.day)
		}
	}

	f("", 0, 0, 0)
	f("=)", 0, 0, 0)
	f("23 мая", 0, time.May, 23)
	f("12 сентября", 0, time.September, 12)
	f("12 сентября", 0, time.September, 12)
	f("9 марта 2013", 2013, time.March, 9)
	f("14 липня (п'ятниця)", 0, time.July, 14)
	f("14 декабря (Четверг)", 0, time.December, 14)
	f("14 december (Friday)", 0, time.December, 14)
}

func TestParseTimeString(t *testing.T) {
	f := func(time string, hours, minutes int) {
		dt := douDateTime{}

		dt.parseTimeString(time)

		if dt.hours != hours {
			t.Errorf(`parseTimeString(%s) failed with hours
				expected: %d
				recieved: %d
			`, time, hours, dt.hours)
		}

		if dt.minutes != minutes {
			t.Errorf(`parseTimeString(%s) failed with minutes
				expected: %d
				recieved: %d
			`, time, minutes, dt.minutes)
		}
	}

	f("", 0, 0)
	f("00:00", 0, 0)
	f("30:0", 0, 0)
	f("1:02", 0, 0)
	f("3:2", 0, 0)
	f("16:53", 16, 53)
	f("16:53", 16, 53)
	f("18:00", 18, 0)
	f("00:05", 0, 5)
}

func TestParseRawDates(t *testing.T) {
	f := func(dates string, start, end douDateTime) {
		var date1, date2 douDateTime

		parseRawDates(dates, &date1, &date2)

		if !isDouDateTimeEqual(start, date1) {
			t.Errorf(`parseRawDates(%s) failed with start DOU calendar date
				expected: %s
				recieved: %s
			`, dates, start, date1)
		}

		if !isDouDateTimeEqual(end, date2) {
			t.Errorf(`parseRawDates(%s) failed with end DOU calendar date
				expected: %s
				recieved: %s
			`, dates, end, date2)
		}
	}

	f("23 октября",
		douDateTime{year: 2019, month: time.October, day: 23},
		douDateTime{})
	f("23 июля 2016",
		douDateTime{year: 2016, month: time.July, day: 23},
		douDateTime{})
	f("26 — 29 ноября",
		douDateTime{year: 2019, month: time.November, day: 26},
		douDateTime{year: 2019, month: time.November, day: 29})
	f("26 сентября — 29 марта",
		douDateTime{year: 2019, month: time.September, day: 26},
		douDateTime{year: 2020, month: time.March, day: 29})
}

func TestParseRawClockTimes(t *testing.T) {
	f := func(times string, start, end douDateTime) {
		var time1, time2 douDateTime

		parseRawClockTimes(times, &time1, &time2)

		if !isDouDateTimeEqual(start, time1) {
			t.Errorf(`parseRawClockTimes(%s) failed with start DOU clock time
				expected: %s
				recieved: %s
			`, times, start, time1)
		}

		if !isDouDateTimeEqual(end, time2) {
			t.Errorf(`parseRawClockTimes(%s) failed with end DOU clock time
				expected: %s
				recieved: %s
			`, times, end, time2)
		}
	}

	f("18:00",
		douDateTime{hours: 18, minutes: 0},
		douDateTime{})
	f("18:40",
		douDateTime{hours: 18, minutes: 40},
		douDateTime{})
	f("18:00 - 21:20",
		douDateTime{hours: 18, minutes: 0},
		douDateTime{hours: 21, minutes: 20})
	f("18:00 -------kek------- 21:30",
		douDateTime{hours: 18, minutes: 0},
		douDateTime{hours: 21, minutes: 30})
}

func TestResolveTimes(t *testing.T) {}
