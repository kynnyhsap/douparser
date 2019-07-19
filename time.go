package douparser

import (
	"regexp"
	"strconv"
)

type douTime struct {
	hours   int
	minutes int
}

func (dt douTime) defined() bool {
	return dt.hours > 0 && dt.minutes > 0
}

func timeFromStrings(hours, minutes string) douTime {
	h, _ := strconv.Atoi(hours)
	m, _ := strconv.Atoi(minutes)

	return douTime{h, m}
}

func parseRawTime(rawTime string) (douTime, douTime) {
	if rawTime == "" {
		return douTime{}, douTime{}
	}

	var re = regexp.MustCompile(`(\d{2}):(\d{2})(?:\D+)?(?:(\d{2}):(\d{2}))?`)
	match := re.FindStringSubmatch(rawTime)

	if match[1] == "" || match[2] == "" {
		return douTime{}, douTime{}
	}

	start := timeFromStrings(match[1], match[2])

	if match[3] == "" || match[4] == "" {
		return start, douTime{}
	}

	end := timeFromStrings(match[3], match[4])

	return start, end
}
