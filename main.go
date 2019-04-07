package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

const eventsListSelector = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
const calendarPageUrl = "https://dou.ua/calendar/page"

func getPageUrl(page int) string {
	return fmt.Sprintf("%s-%d/", calendarPageUrl, page)
}

func parseCalendarPage(page int) ([]Event, error) {
	var events []Event

	res, err := http.Get(getPageUrl(page))
	if err != nil {
		return events, err
	}

	if res.StatusCode == 404 {
		return events, fmt.Errorf("404")
	} else if res.StatusCode != 200 {
		return events, fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return events, err
	}

	err = res.Body.Close()
	if err != nil {
		return events, err
	}

	doc.Find(eventsListSelector).Each(func(i int, s *goquery.Selection) {
		var event Event
		event.ParseQuickInfo(s)
		events = append(events, event)
	})

	return events, nil
}

func parseCalendar() []Event {
	var events []Event

	for i := 0; ; i++ {
		parsedEvents, err := parseCalendarPage(i)
		if err != nil {
			if err.Error() == "404" {
				break
			} else {
				fmt.Printf("ERROR in `parseCalendarPage(%d)`: %s", i, err.Error())
			}
		}

		events = append(events, parsedEvents...)
	}

	return events
}

func main() {
	start := time.Now()
	events := parseCalendar()
	fmt.Printf("Parsed %d events in %f seconds\n", len(events), time.Since(start).Seconds())

	//for key, event := range events {
	//	fmt.Println(key, event)
	//}
}
