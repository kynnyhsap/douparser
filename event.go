package douparser

import "time"

type Event struct {
	ID               int       `json:"id"`
	Online           bool      `json:"online"`
	Tags             []string  `json:"tags"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	RawDate          string    `json:"raw_date"`
	RawTime          string    `json:"raw_time"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	FullDescription  string    `json:"full_description"`
	Location         string    `json:"location"`
	Cost             string    `json:"cost"`
	Image            string    `json:"image"`
}
