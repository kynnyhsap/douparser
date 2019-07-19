package douparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/k3a/html2text"
	"net/http"
	"strconv"
	"strings"
)

const (
	calendarURL = "https://dou.ua/calendar"
	archiveURL  = "https://dou.ua/calendar/archive"

	// css selectors
	allTagsSelector    = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option"
	eventsListSelector = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
)

func eventPageURL(id int) string {
	return calendarURL + "/" + strconv.Itoa(id)
}

func calendarPageURL(page int) string {
	return fmt.Sprintf("%s/page-%d/", calendarURL, page)
}

func scrapEvent(eventID int) (DouEvent, error) {
	res, err := http.Get(eventPageURL(eventID))
	if err != nil {
		return DouEvent{}, err
	}

	if res.StatusCode != 200 {
		return DouEvent{}, fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return DouEvent{}, err
	}

	s := doc.Find("body > div > div.l-content.m-content > div.l-content-wrap > div.cell.g-right-shadowed.mobtab-maincol")

	return parseEvent(s), nil
}

func parseEvent(s *goquery.Selection) DouEvent {
	var event DouEvent

	event.Title = strings.TrimSpace(s.Find(".page-head h1").Text())
	event.Image, _ = s.Find(".event-info img.event-info-logo").Attr("src")

	htmlDescription, err := s.Find("article.b-typo").Html()
	if err == nil {
		event.FullDescription = strings.TrimSpace(html2text.HTML2Text(htmlDescription))
	}

	s.Find(".event-info .event-info-row").Each(func(i int, row *goquery.Selection) {
		t := strings.TrimSpace(row.Find(".dt").Text())
		d := strings.TrimSpace(row.Find(".dd").Text())

		switch t {
		case "Відбудеться":
			event.RawDate = d
			event.Start, event.End = parseRawDate(d)
		case "Початок": // parse time
			break
		case "Місце":
			if d == "Online" {
				event.Online = true
				break
			}
			event.Location = d
		case "Вартість":
			event.Cost = d
		}
		//...
	})

	s.Find(".b-post-tags a").Each(func(i int, tagLink *goquery.Selection) {
		event.Tags = append(event.Tags, tagLink.Text())
	})

	return event
}

func scrapCalendarPage(page int) ([]DouEvent, error) {
	events := make([]DouEvent, 0)

	res, err := http.Get(calendarPageURL(page))
	if err != nil {
		return events, err
	}

	if res.StatusCode == 404 {
		return events, fmt.Errorf("404")
	}
	if res.StatusCode != 200 {
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
		link, _ := selection.Find("h2.title a").Attr("href")
		id, _ := strconv.Atoi(strings.Split(link, "/")[4])

		event, err := scrapEvent(id)
		if err != nil {
			return
		}

		event.ID = id
		event.ShortDescription = selection.Find("p.b-typo").Text()

		events = append(events, event)
	})

	return events, nil
}

func ScrapCalendarEvents() ([]DouEvent, error) {
	var events []DouEvent

	for page := 0; ; page++ {
		scrappedEvents, err := scrapCalendarPage(page)

		events = append(events, scrappedEvents...)

		if err != nil {
			if err.Error() == "404" {
				break
			}

			return events, err
		}
	}

	return events, nil
}

func ScrapEventTags() ([]string, error) {
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
