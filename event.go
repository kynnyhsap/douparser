package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	RawLocation string
	Online      bool
	Cost        string
	Image       string
	Link        string
	Tags        []string
	RawDate     string
	Date        time.Time
}

func (e *Event) ParseQuickInfo(selection *goquery.Selection) {
	title := selection.Find("h2.title").Text()
	e.Title = strings.TrimSpace(title)

	link, hrefExists := selection.Find("h2.title a").Attr("href")
	if !hrefExists {
		log.Fatal("attribute href does not exits")
	}

	e.Link = link

	image, srcExists := selection.Find("h2.title a img.logo").Attr("src")
	if !srcExists {
		log.Fatal("attribute src does not exits")
	}

	e.Image = image

	id, err := strconv.Atoi(strings.Split(e.Link, "/")[4])
	if err != nil {
		log.Fatal(err)
	}

	e.ID = id

	description := selection.Find("p.b-typo").Text()
	e.Description = strings.TrimSpace(description)

	date := selection.Find("div.when-and-where span.date").Text()
	e.RawDate = strings.TrimSpace(date)
	selection.Find("div.when-and-where span.date").Remove()

	cost := selection.Find("div.when-and-where span").Text()
	e.Cost = strings.TrimSpace(cost)
	selection.Find("div.when-and-where span").Remove()

	location := selection.Find("div.when-and-where").Text()
	e.RawLocation = strings.TrimSpace(location)

	if e.RawLocation == "Online" {
		e.Online = true
	}

	selection.Find("div.more a").Each(func(i int, tagSelection *goquery.Selection) {
		tag := tagSelection.Text()

		e.Tags = append(e.Tags, strings.TrimSpace(tag))
	})
}
