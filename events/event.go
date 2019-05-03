package events

import (
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Location    string
	Cost        string
	Image       string
	RawDate     string
	Start       time.Time
	End         time.Time
	Online      bool
	Tags        []string
}
