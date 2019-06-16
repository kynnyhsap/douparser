package events

import (
	"time"
)

type Event struct {
	ID               int
	Online           bool
	Tags             []string
	Start            time.Time
	End              time.Time
	RawDate          string
	Title            string
	ShortDescription string
	FullDescription  string
	Location         string
	Cost             string
	Image            string
}
