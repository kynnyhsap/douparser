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
	f("12 сентября", douDate{12, time.September, 0})
	f("23 мая", douDate{23, time.May, 0})
	f("9 марта 2013", douDate{9, time.March, 2013})
	f("14 декабря (Четверг)", douDate{14, time.December, 0})
	f("14 december (Friday)", douDate{14, time.December, 0})
	f("14 липня (п'ятниця)", douDate{14, time.July, 0})
}

func TestParseRawDate(t *testing.T) {
	f := func(rawDate string, expectedStart, expectedEnd douDate) {
		receivedStart, receivedEnd := parseRawDate(rawDate)

		if !isDouDatesEqual(expectedStart, receivedStart) {
			t.Error("For", rawDate, " want Start date: ", expectedStart, ", got: ", receivedStart)
		}

		if !isDouDatesEqual(expectedEnd, receivedEnd) {
			t.Error("For", rawDate, " want End date: ", expectedEnd, ", got: ", receivedEnd)
		}
	}

	f("23 октября", douDate{23, time.October, 2019}, douDate{})
	f("23 июля 2016", douDate{23, time.July, 2016}, douDate{})
	f("26 — 29 ноября", douDate{26, time.November, 2019}, douDate{29, time.November, 2019})
	f("26 сентября — 29 марта", douDate{26, time.September, 2019}, douDate{29, time.March, 2020})
}
