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
	DateRaw     string
	DateStart   time.Time
	DateEnd     time.Time
}
