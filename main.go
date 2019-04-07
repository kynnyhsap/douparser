package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	var parser CalendarParser

	start := time.Now()

	err := parser.ParseFullCalendar()
	if err != nil {
		log.Fatal(err) // =(
	}

	fmt.Printf("Parsed %d events in %f seconds\n", len(parser.Events), time.Since(start).Seconds())
}
