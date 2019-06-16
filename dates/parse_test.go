package dates

import (
	"testing"
	"time"
)

func TestGetYear(t *testing.T) {
	now := time.Now()

	got := getYear(now.Month())

	if got != now.Year() {
		t.Errorf("Now year should be %d", now.Year())
	}
}

// Test helper: Creates needed date from yyyy-mm-dd
func d(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, locale)
}

func TestParse(t *testing.T) {
	tables := []struct {
		raw string
		d1  time.Time
		d2  time.Time
	}{
		{
			"14 апреля",
			d(2019, time.April, 14),
			time.Time{},
		},
		{
			"7 декабря",
			d(2019, time.December, 7),
			time.Time{},
		},
		{
			"1 марта",
			d(2020, time.March, 1),
			time.Time{},
		},
		{
			"20—21 июля",
			d(2019, time.July, 20),
			d(2019, time.July, 21),
		},
		{
			"    20      — 21 июля",
			d(2019, time.July, 20),
			d(2019, time.July, 21),
		},
		{
			" 27 мая     —  2 сентября  ",
			d(2019, time.May, 27),
			d(2019, time.September, 2),
		},
		{
			"16 апреля — 28 мая",
			d(2019, time.April, 16),
			d(2019, time.May, 28),
		},
		{
			"21 декабря 2018",
			d(2018, time.December, 21),
			time.Time{},
		},
		{
			"3 декабря 2013",
			d(2013, time.December, 3),
			time.Time{},
		},
		//{
		//	"24 декабря — 18 марта",
		//	d(2018, time.December, 24),
		//	d(2019, time.March, 18),
		//},
		//{
		//	"14 — 15 декабря 2018",
		//	d(2018, time.December, 14),
		//	d(2018, time.December, 15),
		//},
	}

	for _, table := range tables {
		got := Parse(table.raw)

		first := got[0]
		if !first.Equal(table.d1) {
			t.Error("Want: ", table.d1, "\nGot: ", first)
		}

		second := got[1]
		if table.d2.IsZero() {
			if !second.IsZero() {
				t.Error("Want: zero date\nGot: ", second)
			}
		} else {
			if second.IsZero() {
				t.Error("Want: ", table.d2, "\nGot: zero date")
			}
		}
	}
}
