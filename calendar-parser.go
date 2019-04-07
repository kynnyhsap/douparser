package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

const eventsListSelector = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
const calendarPageUrl = "https://dou.ua/calendar/page"

type CalendarParser struct {
	Events []Event
}

func (cp *CalendarParser) pageURL(page int) string {
	return fmt.Sprintf("%s-%d/", calendarPageUrl, page)
}

func (cp *CalendarParser) ParsePage(page int) error {
	res, err := http.Get(cp.pageURL(page))
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		return fmt.Errorf("404")
	} else if res.StatusCode != 200 {
		return fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	doc.Find(eventsListSelector).Each(func(i int, s *goquery.Selection) {
		var event Event
		event.ParseQuickInfo(s)

		cp.Events = append(cp.Events, event)
	})

	return nil
}

func (cp *CalendarParser) ParseFullCalendar() error {
	for page := 0; ; page++ {
		err := cp.ParsePage(page)

		if err != nil {
			if err.Error() == "404" {
				break
			} else {
				//fmt.Printf("ERROR in `ParsePage(%d)` method: %s", i, err.Error())
				return err
			}
		}
	}

	return nil
}
