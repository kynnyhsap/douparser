package douparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

const (
	calendarURL = "https://dou.ua/calendar"
	archiveURL  = "https://dou.ua/calendar/archive"

	// css selectors
	allTagsSelector     = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option"
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

func singleEventURL(id int) string {
	return calendarURL + "/" + strconv.Itoa(id)
}

func eventsPageURL(page int) string {
	return fmt.Sprintf("%s/page-%d/", calendarURL, page)
}

func scrapPage(page int) ([]DouEvent, error) {
	events := make([]DouEvent, 0)

	res, err := http.Get(eventsPageURL(page))
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

	doc.Find(eventsListSelector).Each(func(i int, selection *goquery.Selection) {
		event := parseEvent(selection)
		// scrap single event (image url, full description, address)
		events = append(events, event)
	})

	return events, nil
}

func parseEvent(selection *goquery.Selection) DouEvent {
	var event DouEvent

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
	event.Start, event.End = parseRawDate(event.RawDate)

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

func ParseCalendarEvents() ([]DouEvent, error) {
	var events []DouEvent

	for page := 0; ; page++ {
		parsedEventsFromPage, err := scrapPage(page)

		events = append(events, parsedEventsFromPage...)

		if err != nil {
			if err.Error() == "404" {
				break
			} else {
				return events, err
			}
		}
	}

	return events, nil
}

func ParseEventTags() ([]string, error) {
	tags := make([]string, 0)

	res, err := http.Get(archiveURL)
	if err != nil {
		return tags, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return tags, err
	}

	err = res.Body.Close()
	if err != nil {
		return tags, err
	}

	doc.Find(allTagsSelector).Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tags = append(tags, tag)
	})

	return tags, nil
}
