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

func TestParse(t *testing.T) {
	suits := []struct {
		raw string

		day1   int
		year1  int
		month1 time.Month

		day2   int
		year2  int
		month2 time.Month
	}{
		{"14 апреля", 14, 2019, time.April, 0, 0, 0},
		{"7 декабря", 7, 2019, time.December, 0, 0, 0},
		{"1 марта", 1, 2020, time.March, 0, 0, 0},
		{"20—21 июля", 20, 2019, time.July, 21, 2019, time.July},
		{"20  —       21 июля", 20, 2019, time.July, 21, 2019, time.July},
		{"    27 мая     — \n  27 сентября  ", 27, 2019, time.May, 27, 2019, time.September},
	}

	for _, s := range suits {
		got := Parse(s.raw)

		first := got[0]
		if first.Year() != s.year1 || first.Month() != s.month1 || first.Day() != s.day1 {
			t.Errorf("Want: %d-%d-%d\nGot: %d-%d-%d", s.year1, s.month1, s.day1, first.Year(), first.Month(), first.Day())
		}

		second := got[1]
		if s.day2 == 0 || s.month2 == 0 || s.year2 == 0 {
			if !second.IsZero() {
				t.Errorf("Expected: zero date\nGot: %d-%d-%d", second.Year(), first.Month(), first.Day())
			}
		} else {
			if second.IsZero() {
				t.Errorf("Expected: %d-%d-%d\nGot: zero date", s.year2, s.month2, s.day2)
			}
		}
	}
}
