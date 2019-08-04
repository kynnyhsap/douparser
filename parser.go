package douparser

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

const (
	eventsPerPage = 20
)

func scrapEvent(eventID int) (Event, error) {
	doc, err := fetchEventDocument(eventID)
	if err != nil {
		return Event{}, err
	}

	event := parseEvent(doc)
	event.ID = eventID

	return event, nil
}

func totalPages(fromArchive bool) (int, error) {
	const selector = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.b-paging > span:nth-last-child(2) > a"

	doc, err := fetchPageDocument(1, fromArchive)
	if err != nil {
		return 0, err
	}

	total, err := strconv.Atoi(doc.Find(selector).Text())
	if err != nil {
		return 0, err
	}

	return total, err
}

func Events() ([]Event, error) {
	total, _ := totalPages(false)
	events := make([]Event, 0, total*eventsPerPage)

	for i := 1; i <= total; i++ {
		doc, err := fetchPageDocument(i, false)
		if err != nil {
			return events, err
		}

		parsed := parsePageEvents(doc)
		events = append(events, parsed...)
	}

	return events, nil
}

func Tags() ([]string, error) {
	tags := make([]string, 0)

	doc, err := fetchPageDocument(1, true)
	if err != nil {
		return tags, err
	}

	doc.Find(selectors["tags"]).Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tags = append(tags, tag)
	})

	return tags, nil
}
