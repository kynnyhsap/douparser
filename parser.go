package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EventsParser struct {
	FromArchive bool
	Events      []Event
}

const (
	calendarURL         = "https://dou.ua/calendar"
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

func correctYear(month time.Month) int {
	now := time.Now()

	if month < now.Month() {
		return now.Year() + 1
	}

	return now.Year()
}

func date(m time.Month, d int) time.Time {
	return time.Date(correctYear(m), m, d, 9, 0, 0, 0, locale)
}

func parseRawDate(raw string) []time.Time {
	var dates []time.Time

	rawDates := strings.Split(raw, " — ")
	if len(rawDates) < 2 {
		meta := strings.Split(strings.TrimSpace(rawDates[0]), " ")

		month := months[meta[1]]
		day, _ := strconv.Atoi(meta[0])

		d := date(month, day)

		dates = append(dates, d)
	} else {
		meta1 := strings.Split(strings.TrimSpace(rawDates[0]), " ")
		meta2 := strings.Split(strings.TrimSpace(rawDates[1]), " ")

		if len(meta1) < 2 {
			day1, _ := strconv.Atoi(meta1[0])
			day2, _ := strconv.Atoi(meta2[0])
			month := months[meta2[1]]

			d1 := date(month, day1)
			d2 := date(month, day2)

			dates = append(dates, d1, d2)
		} else {
			day1, _ := strconv.Atoi(meta1[0])
			day2, _ := strconv.Atoi(meta2[0])
			month1 := months[meta1[1]]
			month2 := months[meta2[1]]

			d1 := date(month1, day1)
			d2 := date(month2, day2)

			dates = append(dates, d1, d2)
		}
	}

	return dates
}

func (p *EventsParser) pageURL(page int) string {
	if p.FromArchive {
		return fmt.Sprintf("%s/archive/%d/", calendarURL, page)
	}

	return fmt.Sprintf("%s/page-%d/", calendarURL, page)
}

func (p *EventsParser) ParsePage(page int) error {
	res, err := http.Get(p.pageURL(page))
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
		p.Events = append(p.Events, p.ParseEvent(s))
	})

	return nil
}

func (p *EventsParser) ParseAll() error {
	for page := 0; ; page++ {
		err := p.ParsePage(page)

		if err != nil {
			if err.Error() == "404" {
				break
			} else {
				return err
			}
		}
	}

	return nil
}

func (p *EventsParser) ParseEvent(selection *goquery.Selection) Event {
	var event Event

	title := selection.Find(titleSelector).Text()
	event.Title = strings.TrimSpace(title)

	event.Link, _ = selection.Find(linkSelector).Attr("href")

	event.Image, _ = selection.Find(imageSelector).Attr("src")

	event.ID, _ = strconv.Atoi(strings.Split(event.Link, "/")[4])
	description := selection.Find(descriptionSelector).Text()

	event.Description = strings.TrimSpace(description)

	date := selection.Find(dateSelector).Text()
	selection.Find(dateSelector).Remove()
	event.DateRaw = strings.TrimSpace(date)
	dates := parseRawDate(event.DateRaw)
	event.DateStart = dates[0]
	if len(dates) == 2 {
		event.DateEnd = dates[1]
	} else {
		event.DateEnd = time.Date(1990, 1, 1, 1, 1, 1, 1, locale)
	}

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
