package douparser

import (
	"testing"
)

func isDouTimesEqual(t1, t2 douTime) bool {
	return t1.minutes == t2.minutes && t1.hours == t2.hours
}

func TestParseRawTime(t *testing.T) {
	f := func(rawTime string, expectedStart, expectedEnd douTime) {
		receivedStart, receivedEnd := parseRawTime(rawTime)

		if !isDouTimesEqual(expectedStart, receivedStart) {
			t.Error("For", rawTime, " want Start time: ", expectedStart, ", got: ", receivedStart)
		}

		if !isDouTimesEqual(expectedEnd, receivedEnd) {
			t.Error("For", rawTime, " want End time: ", expectedEnd, ", got: ", receivedEnd)
		}
	}

	f("", douTime{}, douTime{})
	f("18:31", douTime{18, 31}, douTime{})
	f("14:00 - 16:30", douTime{14, 0}, douTime{16, 30})
	f("18:01 =bad=text= 21:34", douTime{18, 1}, douTime{21, 34})
}
