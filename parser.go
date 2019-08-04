package douparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/k3a/html2text"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	eventsPerPage = 20
	calendarURL   = "https://dou.ua/calendar"
	archiveURL    = "https://dou.ua/calendar/archive"

	// css selectors
	tagsSelector           = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option"
	eventCardsListSelector = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
	eventCellSelector      = "body > div > div.l-content.m-content > div.l-content-wrap > div.cell.g-right-shadowed.mobtab-maincol"
)

var selectors = map[string]string{
	"tags":      "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option",
	"eventCard": "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article",
	"eventFull": "body > div > div.l-content.m-content > div.l-content-wrap > div.cell.g-right-shadowed.mobtab-maincol",
}

func fetchEventDocument(eventID int) (*goquery.Document, error) {
	res, err := http.Get(fmt.Sprintf("%s/%d", calendarURL, eventID))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

func scrapEvent(eventID int) (Event, error) {
	doc, err := fetchEventDocument(eventID)
	if err != nil {
		return Event{}, err
	}

	event := parseEvent(doc)
	event.ID = eventID

	return event, nil
}

func parseEvent(eventDocument *goquery.Document) Event {
	var event Event

	s := eventDocument.Find(eventCellSelector)

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

func totalPages() (int, error) {
	const selector = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.b-paging > span:nth-last-child(2) > a"

	res, err := http.Get(calendarURL)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, err
	}

	err = res.Body.Close()
	if err != nil {
		return 0, err
	}

	total, err := strconv.Atoi(doc.Find(selector).Text())
	if err != nil {
		return 0, err
	}

	return total, err
}

func fetchPageDocument(page int) (*goquery.Document, error) {
	res, err := http.Get(fmt.Sprintf("%s/page-%d/", calendarURL, page))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("404")
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return doc, err
	}

	return doc, err
}

func parsePageEvents(pageDocument *goquery.Document) []Event {
	events := make([]Event, 0, eventsPerPage)

	pageDocument.Find("body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article").Each(func(i int, selection *goquery.Selection) {
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

func Events() []Event {
	total, _ := totalPages()
	events := make([]Event, 0, total*eventsPerPage)

	//var wg sync.WaitGroup

	for i := 0; i < total; i++ {
		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//
		//}()
		doc, err := fetchPageDocument(i)
		if err != nil {
			log.Print(err)
		}

		parsed := parsePageEvents(doc)
		events = append(events, parsed...)

		//<-time.After(200 * time.Millisecond)
	}

	//wg.Wait()

	return events
}

func Tags() ([]string, error) {
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

	doc.Find(tagsSelector).Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tags = append(tags, tag)
	})

	return tags, nil
}
