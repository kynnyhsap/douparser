package douparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
)

const (
	calendarURL = "https://dou.ua/calendar"
	archiveURL  = calendarURL + "/archive"
)

func buildEventsPageUrl(page int, fromArchive bool) string {
	if fromArchive {
		return archiveURL + "/" + strconv.Itoa(page) + "/"
	}

	return calendarURL + "/page-" + strconv.Itoa(page) + "/"
}

func fetchEventDocument(eventID int) (*goquery.Document, error) {
	url := calendarURL + "/" + strconv.Itoa(eventID) + "/"
	res, err := http.Get(url)
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

func fetchPageDocument(page int, fromArchive bool) (*goquery.Document, error) {
	res, err := http.Get(buildEventsPageUrl(page, fromArchive))
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
