package dates

import (
	"strconv"
	"strings"
	"time"
)

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

var locale, _ = time.LoadLocation("Local")

func getYear(month time.Month) int {
	now := time.Now()

	if month < now.Month() {
		return now.Year() + 1
	}

	return now.Year()
}

func newDate(month time.Month, day int) time.Time {
	return time.Date(getYear(month), month, day, 9, 0, 0, 0, locale)
}

func Parse(raw string) [2]time.Time {
	var dates [2]time.Time

	rawDates := strings.Split(raw, " — ")

	if len(rawDates) < 2 {
		meta := strings.Split(strings.TrimSpace(rawDates[0]), " ")

		month := months[meta[1]]
		day, _ := strconv.Atoi(meta[0])

		dates[0] = newDate(month, day)
	} else {
		meta1 := strings.Split(strings.TrimSpace(rawDates[0]), " ")
		meta2 := strings.Split(strings.TrimSpace(rawDates[1]), " ")

		if len(meta1) < 2 {
			day1, _ := strconv.Atoi(meta1[0])
			day2, _ := strconv.Atoi(meta2[0])
			month := months[meta2[1]]

			dates[0] = newDate(month, day1)
			dates[1] = newDate(month, day2)
		} else {
			day1, _ := strconv.Atoi(meta1[0])
			day2, _ := strconv.Atoi(meta2[0])
			month1 := months[meta1[1]]
			month2 := months[meta2[1]]

			dates[0] = newDate(month1, day1)
			dates[1] = newDate(month2, day2)
		}
	}

	return dates
}
