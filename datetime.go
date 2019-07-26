package douparser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var months = map[string]time.Month{
	// English months
	"january":   time.January,
	"february":  time.February,
	"match":     time.March,
	"april":     time.April,
	"may":       time.May,
	"june":      time.June,
	"july":      time.July,
	"august":    time.August,
	"september": time.September,
	"october":   time.October,
	"november":  time.November,
	"december":  time.December,
	// Russian months
	"января":   time.January,
	"февраля":  time.February,
	"марта":    time.March,
	"апреля":   time.April,
	"мая":      time.May,
	"июня":     time.June,
	"июля":     time.July,
	"августа":  time.August,
	"сентября": time.September,
	"октября":  time.October,
	"ноября":   time.November,
	"декабря":  time.December,
	// Ukrainian months
	"січня":     time.January,
	"лютого":    time.February,
	"березня":   time.March,
	"квітня":    time.April,
	"травня":    time.May,
	"червня":    time.June,
	"липня":     time.July,
	"серпня":    time.August,
	"вересня":   time.September,
	"жовтня":    time.October,
	"листопада": time.November,
	"грудня":    time.December,
}

func defineMonth(month string) time.Month {
	m := strings.ToLower(month)

	return months[m]
}

func defineYearByMonth(month time.Month) int {
	now := time.Now()

	if now.Month() > month {
		return now.Year() + 1
	}

	return now.Year()
}

// douDateTime - defines and help manipulate with DOU dates and time for Event
type douDateTime struct {
	year    int
	month   time.Month
	day     int
	hours   int
	minutes int
}

func (ddt douDateTime) hasYear() bool {
	return ddt.year > 0
}

func (ddt douDateTime) hasMonth() bool {
	return ddt.month > 0
}

func (ddt douDateTime) hasDay() bool {
	return ddt.day > 0
}

func (ddt douDateTime) dateDefined() bool {
	return ddt.year > 0 && ddt.month > 0 && ddt.day > 0
}

func (ddt douDateTime) timeDefined() bool {
	return ddt.hours > 0 && ddt.minutes >= 0
}

func (ddt douDateTime) toStdTime() time.Time {
	if !ddt.timeDefined() && !ddt.dateDefined() {
		return time.Time{}
	}

	return time.Date(ddt.year, ddt.month, ddt.day, ddt.hours, ddt.minutes, 0, 0, time.UTC)
}

func (ddt *douDateTime) parseTimeString(time string) {
	var re = regexp.MustCompile(`(?:\D*)(\d{2}):(\d{2})(?:\D*)`)
	match := re.FindStringSubmatch(time)

	if match[1] == "" || match[2] == "" {
		return
	}

	ddt.hours, _ = strconv.Atoi(match[1])
	ddt.minutes, _ = strconv.Atoi(match[2])
}

func (ddt *douDateTime) parseDateString(date string) {
	var re = regexp.MustCompile(`(?:\s*)(\d{1,2})(?:(?:\s+)([а-яА-яa-zA-Z]+))?(?:(?:\s+)(?:(\d{4})|(?:\(.+\))))?`) // TODO: extract to global scope

	match := re.FindStringSubmatch(date)

	ddt.day, _ = strconv.Atoi(match[1])

	if match[2] != "" {
		ddt.month = defineMonth(match[2])
	}

	if match[3] != "" {
		ddt.year, _ = strconv.Atoi(match[3])
	}
}

func parseRawDates(raw string, start, end *douDateTime) {
	dates := strings.Split(raw, "—")

	start.parseDateString(dates[0])

	if len(dates) == 2 {
		end.parseDateString(dates[1])

		if !end.hasYear() {
			end.year = defineYearByMonth(end.month)
		}

		if !start.hasMonth() {
			start.month = end.month
		}

		if !start.hasYear() {
			if start.hasMonth() {
				start.year = defineYearByMonth(start.month)
			} else {
				start.year = end.year
			}
		}
	} else {
		if !start.hasYear() {
			start.year = defineYearByMonth(start.month)
		}
	}
}

func parseRawTime(raw string, start, end *douDateTime) {
	if raw == "" {
		return
	}

	var re = regexp.MustCompile(`(?:\D*)(\d{2}:\d{2})(?:\D+)?(\d{2}:\d{2})?`)
	match := re.FindStringSubmatch(raw)

	if match[1] != "" {
		start.parseTimeString(match[1])
	}

	if match[2] != "" {
		end.parseTimeString(match[2])
	}
}

func resolveTimes(start, end *douDateTime) {
	if !start.dateDefined() {
		return
	}

	if end.dateDefined() {
		end.hours = start.hours
		end.minutes = start.minutes
		return
	}

	if end.timeDefined() {
		end.year = start.year
		end.month = start.month
		end.day = start.day

		return
	}
}

func getEventTime(rawDates, rawTimes string) (time.Time, time.Time) {
	var start, end douDateTime

	parseRawDates(rawDates, &start, &end)
	parseRawTime(rawTimes, &start, &end)
	resolveTimes(&start, &end)

	return start.toStdTime(), end.toStdTime()
}
