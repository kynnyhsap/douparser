package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/k3a/html2text"
	"strconv"
	"strings"
)

var selectors = map[string]string{
	"tags":      "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option",
	"eventCard": "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article",
	"eventFull": "body > div > div.l-content.m-content > div.l-content-wrap > div.cell.g-right-shadowed.mobtab-maincol",
}

func parseEvent(eventDocument *goquery.Document) Event {
	s := eventDocument.Find(selectors["eventFull"])

	var event Event

	event.Title = strings.TrimSpace(s.Find(".page-head h1").Text())
	event.Image, _ = s.Find(".event-info img.event-info-logo").Attr("src")
	htmlDescription, err := s.Find("article.b-typo").Html()
	if err == nil {
		event.FullDescription = strings.TrimSpace(html2text.HTML2Text(htmlDescription))
	}

	s.Find(".event-info .event-info-row").Each(func(i int, infoRow *goquery.Selection) {
		infoType := strings.TrimSpace(infoRow.Find(".dt").Text())
		d := strings.TrimSpace(infoRow.Find(".dd").Text())

		switch infoType {
		case "Відбудеться", "Пройдет", "Date":
			event.RawDate = d
		case "Початок", "Начало", "Time":
			event.RawTime = d
		case "Місце", "Место", "Place":
			if d == "Online" {
				event.Online = true
				break
			}
			event.Location = d
		case "Вартість", "Стоимость", "Price":
			event.Cost = d
		}
	})

	event.Start, event.End = parseEventTime(event.RawDate, event.RawTime)

	s.Find(".b-post-tags a").Each(func(i int, tagLink *goquery.Selection) {
		event.Tags = append(event.Tags, tagLink.Text())
	})

	return event
}

func parsePageEvents(pageDocument *goquery.Document) []Event {
	events := make([]Event, 0, eventsPerPage)

	pageDocument.Find(selectors["eventCard"]).Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Find("h2.title a").Attr("href")
		id, _ := strconv.Atoi(strings.Split(link, "/")[4])

		event, err := scrapEvent(id)
		if err != nil {
			return
		}

		event.ShortDescription = strings.TrimSpace(selection.Find("p.b-typo").Text())

		events = append(events, event)
	})

	return events
}
