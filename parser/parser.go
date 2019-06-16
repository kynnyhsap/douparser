package parser

import (
	"dou-parser/dates"
	"dou-parser/events"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

type EventsParser struct {
	FromArchive bool
	Events      []events.Event
	Tags        []string
}

const (
	calendarURL         = "https://dou.ua/calendar"
	archiveURL          = "https://dou.ua/calendar/archive"
	eventsListSelector  = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
	titleSelector       = "h2.title"
	linkSelector        = "h2.title a"
	imageSelector       = "h2.title a img.logo"
	descriptionSelector = "p.b-typo"
	dateSelector        = "div.when-and-where span.date"
	costSelector        = "div.when-and-where span"
	locationSelector    = "div.when-and-where"
	tagsSelector        = "div.more a"

	allTagsSelector = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option"
)

func singleEventURL(id int) string {
	return calendarURL + "/" + strconv.Itoa(id)
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
		event := p.ParseEvent(s)

		// scrap single event (image url, full description, address)

		p.Events = append(p.Events, event)
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

func (p *EventsParser) ParseEvent(selection *goquery.Selection) events.Event {
	var event events.Event

	title := selection.Find(titleSelector).Text()
	event.Title = strings.TrimSpace(title)

	event.Image, _ = selection.Find(imageSelector).Attr("src")

	link, _ := selection.Find(linkSelector).Attr("href")
	event.ID, _ = strconv.Atoi(strings.Split(link, "/")[4])
	description := selection.Find(descriptionSelector).Text()

	event.ShortDescription = strings.TrimSpace(description)

	date := selection.Find(dateSelector).Text()
	selection.Find(dateSelector).Remove()
	event.RawDate = strings.TrimSpace(date)
	parsedDates := dates.Parse(event.RawDate)
	event.Start = parsedDates[0]
	event.End = parsedDates[1]

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

func (p *EventsParser) ParseTags() error {
	res, err := http.Get(archiveURL)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	doc.Find(allTagsSelector).Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		p.Tags = append(p.Tags, tag)
	})

	return nil
}
