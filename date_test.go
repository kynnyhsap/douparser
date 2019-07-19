package douparser

import (
	"testing"
	"time"
)

func isDouDatesEqual(d1, d2 douDate) bool {
	return d1.year == d2.year &&
		d1.month == d2.month &&
		d1.day == d2.day
}

func TestNewDouDateFromString(t *testing.T) {
	f := func(rawDate string, expectedDate douDate) {
		receivedDate := newDouDateFromString(rawDate)

		if !isDouDatesEqual(receivedDate, expectedDate) {
			t.Error("For", rawDate, " want: ", expectedDate, ", got: ", receivedDate)
		}
	}

	f("12 сентября", douDate{12, time.September, 0})
	f("23 мая", douDate{23, time.May, 0})
	f("9 марта 2013", douDate{9, time.March, 2013})
	f("14 декабря (Четверг)", douDate{14, time.December, 0})
	f("14 december (Friday)", douDate{14, time.December, 0})
	f("14 грудень (п'ятниця)", douDate{14, time.December, 0})
}

func TestParseRawDate(t *testing.T) {
	f := func(rawDate string, expectedStart, expectedEnd time.Time) {
		receivedStart, receivedEnd := parseRawDate(rawDate)

		if !receivedStart.Equal(expectedStart) {
			t.Error("For", rawDate, " want Start date: ", expectedStart, ", got: ", receivedStart)
		}

		if !receivedEnd.Equal(expectedEnd) {
			t.Error("For", rawDate, " want End date: ", expectedEnd, ", got: ", receivedEnd)
		}
	}

	f("23 октября",
		time.Date(2019, time.October, 23, 0, 0, 0, 0, locale),
		time.Time{})

	f("23 июля 2016",
		time.Date(2016, time.July, 23, 0, 0, 0, 0, locale),
		time.Time{})

	f("26 — 29 ноября",
		time.Date(2019, time.November, 26, 0, 0, 0, 0, locale),
		time.Date(2019, time.November, 29, 0, 0, 0, 0, locale))

	f("26 сентября — 29 марта",
		time.Date(2019, time.September, 26, 0, 0, 0, 0, locale),
		time.Date(2020, time.March, 29, 0, 0, 0, 0, locale))
}
