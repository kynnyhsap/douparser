package douparser

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func Tags() ([]string, error) {
	tags := make([]string, 0)

	doc, err := fetchPageDocument(1, true)
	if err != nil {
		return tags, err
	}

	doc.Find(selectors["tags"]).Each(func(i int, s *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(s.Text()))
	})

	return tags, nil
}
