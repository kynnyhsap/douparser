package douparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/k3a/html2text"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	calendarURL = "https://dou.ua/calendar"
	archiveURL  = "https://dou.ua/calendar/archive"

	// css selectors
	tagsSelector           = "body > div > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > div.page-head > h1 > select:nth-child(3) > option"
	eventCardsListSelector = "body > div.g-page > div.l-content.m-content > div > div.col70.m-cola > div > div > div.col50.m-cola > article"
	eventCellSelector      = "body > div > div.l-content.m-content > div.l-content-wrap > div.cell.g-right-shadowed.mobtab-maincol"
)

func eventPageURL(id int) string {
	return calendarURL + "/" + strconv.Itoa(id)
}

func calendarPageURL(page int) string {
	return fmt.Sprintf("%s/page-%d/", calendarURL, page)
}

func scrapEvent(eventID int) (Event, error) {
	res, err := http.Get(eventPageURL(eventID))
	if err != nil {
		return Event{}, err
	}

	if res.StatusCode != 200 {
		return Event{}, fmt.Errorf("unhandled response status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Event{}, err
	}

	event := parseEvent(doc.Find(eventCellSelector))

	return event, nil
}

func newDateFromDouDateAndTime(d douDate, t douTime) time.Time {
	return time.Date(d.year, d.month, d.day, t.hours, t.minutes, 0, 0, locale)
}

func combineDouTimeAndDate(sTime, eTime douTime, sDate, eDate douDate) (time.Time, time.Time) {
	if !sDate.defined() {
		return time.Time{}, time.Time{}
	}

	start := newDateFromDouDateAndTime(sDate, sTime)

	if eDate.defined() {
		return start, newDateFromDouDateAndTime(eDate, sTime)
	}

	if eTime.defined() {
		return start, newDateFromDouDateAndTime(sDate, eTime)
	}

	return start, time.Time{}
}

func parseEvent(s *goquery.Selection) Event {
	var event Event

	event.Title = strings.TrimSpace(s.Find(".page-head h1").Text())

	event.Image, _ = s.Find(".event-info img.event-info-logo").Attr("src")

	htmlDescription, err := s.Find("article.b-typo").Html()
	if err == nil {
		event.FullDescription = strings.TrimSpace(html2text.HTML2Text(htmlDescription))
	}

	var startTime, endTime douTime
	var startDate, endDate douDate

	s.Find(".event-info .event-info-row").Each(func(i int, infoRow *goquery.Selection) {
		infoType := strings.TrimSpace(infoRow.Find(".dt").Text())
		d := strings.TrimSpace(infoRow.Find(".dd").Text())

		switch infoType {
		case "Відбудеться", "Пройдет", "Date":
			event.RawDate = d
			startDate, endDate = parseRawDate(d)
		case "Початок", "Начало", "Time":
			event.RawTime = d
			startTime, endTime = parseRawTime(d)
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

	event.Start, event.End = combineDouTimeAndDate(startTime, endTime, startDate, endDate)

	s.Find(".b-post-tags a").Each(func(i int, tagLink *goquery.Selection) {
		event.Tags = append(event.Tags, tagLink.Text())
	})

	return event
}

func scrapCalendarPage(page int) ([]Event, error) {
	events := make([]Event, 0)

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

	doc.Find(eventCardsListSelector).Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Find("h2.title a").Attr("href")
		id, _ := strconv.Atoi(strings.Split(link, "/")[4])

		event, err := scrapEvent(id)
		if err != nil {
			return
		}

		event.ID = id
		event.ShortDescription = strings.TrimSpace(selection.Find("p.b-typo").Text())

		events = append(events, event)
	})

	return events, nil
}

func CalendarEvents() ([]Event, error) {
	var events []Event

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
