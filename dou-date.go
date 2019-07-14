package douparser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var locale, _ = time.LoadLocation("Local")

type DouDate struct {
	day   int
	month time.Month
	year  int
}

func (dd DouDate) toTimeDate() time.Time {
	if dd.isYearDefined() && dd.isMonthDefined() && dd.isDayDefined() {
		return time.Time{}
	}

	return time.Date(dd.year, dd.month, dd.day, 0, 0, 0, 0, locale)
}

func (dd DouDate) isMonthDefined() bool {
	return dd.month > 0
}

func (dd DouDate) isYearDefined() bool {
	return dd.year > 0
}

func (dd DouDate) isDayDefined() bool {
	return dd.day > 0
}

func newDouDateFromString(s string) DouDate {
	re := regexp.MustCompile(`(?:\s*)(\d{1,2})(?:(?:\s+)([а-я]+))?(?:(?:\s+)(\d{4}))?`)

	match := re.FindStringSubmatch(s)

	date := DouDate{}

	date.day, _ = strconv.Atoi(match[1])
	if len(match[2]) != 0 {
		date.month = getMonth(match[2])
	}
	if len(match[3]) != 0 {
		date.year, _ = strconv.Atoi(match[3])
	}

	return date
}

func getMonth(m string) time.Month {
	var months = map[string]time.Month{
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
	}

	return months[m]
}

func getYearByMonth(month time.Month) int {
	now := time.Now()

	currentYear := now.Year()
	currentMonth := now.Month()

	if month < currentMonth {
		return currentYear + 1
	}

	return currentYear
}

func parseRawDate(rawDateString string) (time.Time, time.Time) {
	var start DouDate
	var end DouDate

	rawDates := strings.Split(rawDateString, "—")
	start = newDouDateFromString(rawDates[0])

	hasEndDate := len(rawDates) == 2

	if hasEndDate {
		end = newDouDateFromString(rawDates[1])

		if !end.isYearDefined() {
			end.year = getYearByMonth(end.month)
		}

		if !start.isMonthDefined() {
			start.month = end.month
		}

		if !start.isYearDefined() {
			start.year = end.year
		}
	} else {
		if !start.isYearDefined() {
			start.year = getYearByMonth(start.month)
		}
	}

	return start.toTimeDate(), end.toTimeDate()
}
