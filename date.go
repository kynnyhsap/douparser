package douparser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var locale, _ = time.LoadLocation("Local")

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
	"січень":   time.January,
	"лютий":    time.February,
	"березень": time.March,
	"квітень":  time.April,
	"травень":  time.May,
	"червень":  time.June,
	"липень":   time.July,
	"серпень":  time.August,
	"вересень": time.September,
	"жовтень":  time.October,
	"листопад": time.November,
	"грудень":  time.December,
}

type douDate struct {
	day   int
	month time.Month
	year  int
}

func (dd douDate) toTimeDate() time.Time {
	if !dd.yearDefined() && !dd.monthDefined() && !dd.dayDefined() {
		return time.Time{}
	}

	return time.Date(dd.year, dd.month, dd.day, 0, 0, 0, 0, locale)
}

func (dd douDate) monthDefined() bool {
	return dd.month > 0
}

func (dd douDate) yearDefined() bool {
	return dd.year > 0
}

func (dd douDate) dayDefined() bool {
	return dd.day > 0
}

func newDouDateFromString(s string) douDate {
	var re = regexp.MustCompile(`(?:\s*)(\d{1,2})(?:(?:\s+)([а-яА-яa-zA-Z]+))?(?:(?:\s+)(?:(\d{4})|(?:\(.+\))))?`)

	match := re.FindStringSubmatch(s)

	date := douDate{}

	date.day, _ = strconv.Atoi(match[1])
	if len(match[2]) != 0 {
		date.month = defineMonth(match[2])
	}
	if len(match[3]) != 0 {
		date.year, _ = strconv.Atoi(match[3])
	}

	return date
}

func defineMonth(m string) time.Month {
	m = strings.ToLower(m)

	return months[m]
}

func defineYearByMonth(month time.Month) int {
	now := time.Now()

	if now.Month() > month {
		return now.Year() + 1
	}

	return now.Year()
}

func parseRawDate(rawDate string) (time.Time, time.Time) {
	var start douDate
	var end douDate

	dates := strings.Split(rawDate, "—")

	start = newDouDateFromString(dates[0])

	if len(dates) == 2 {
		end = newDouDateFromString(dates[1])

		if !end.yearDefined() {
			end.year = defineYearByMonth(end.month)
		}

		if !start.monthDefined() {
			start.month = end.month
		}

		if !start.yearDefined() {
			if start.monthDefined() {
				start.year = defineYearByMonth(start.month)
			} else {
				start.year = end.year
			}
		}
	} else {
		if !start.yearDefined() {
			start.year = defineYearByMonth(start.month)
		}
	}

	return start.toTimeDate(), end.toTimeDate()
}
