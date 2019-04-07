package main

import (
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Location    string
	Online      bool
	Cost        string
	Image       string
	Link        string
	Tags        []string
	RawDate     string
	Date        time.Time
}
