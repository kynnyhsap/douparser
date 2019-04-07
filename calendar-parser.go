package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

const (
	calendarPageUrl     = "https://dou.ua/calendar/page"
	eventsListSelector  = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
	titleSelector       = "h2.title"
	linkSelector        = "h2.title a"
	imageSelector       = "h2.title a img.logo"
	descriptionSelector = "p.b-typo"
	dateSelector        = "div.when-and-where span.date"
	costSelector        = "div.when-and-where span"
	locationSelector    = "div.when-and-where"
	tagsSelector        = "div.more a"
)

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
		cp.Events = append(cp.Events, cp.ParseEvent(s))
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

func (cp *CalendarParser) ParseEvent(selection *goquery.Selection) Event {
	var event Event

	title := selection.Find(titleSelector).Text()
	event.Title = strings.TrimSpace(title)

	event.Link, _ = selection.Find(linkSelector).Attr("href")

	event.Image, _ = selection.Find(imageSelector).Attr("src")

	event.ID, _ = strconv.Atoi(strings.Split(event.Link, "/")[4])
	description := selection.Find(descriptionSelector).Text()

	event.Description = strings.TrimSpace(description)

	date := selection.Find(dateSelector).Text()
	event.RawDate = strings.TrimSpace(date)
	selection.Find(dateSelector).Remove()

	cost := selection.Find(costSelector).Text()
	event.Cost = strings.TrimSpace(cost)
	selection.Find(costSelector).Remove()

	location := selection.Find(locationSelector).Text()
	event.Location = strings.TrimSpace(location)
	if event.Location == "Online" {
		event.Online = true
	}

	selection.Find(tagsSelector).Each(func(i int, tagSelection *goquery.Selection) {
		tag := tagSelection.Text()
		event.Tags = append(event.Tags, strings.TrimSpace(tag))
	})

	return event
}
