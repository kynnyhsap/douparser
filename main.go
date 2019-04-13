package main

import (
	"dou-parser/parser"
	"fmt"
	"log"
	"time"
)

func main() {
	var calendarParser parser.EventsParser

	start := time.Now()

	err := calendarParser.ParseAll()
	if err != nil {
		log.Fatal(err) // =(
	}

	fmt.Printf("Parsed %d events in %f seconds\n", len(calendarParser.Events), time.Since(start).Seconds())

	for i, event := range calendarParser.Events {
		if event.End.IsZero() {
			fmt.Printf("%d(%s): %s \n", i+1, event.RawDate, event.Start)
		} else {
			fmt.Printf("%d(%s): %s - %s \n", i+1, event.RawDate, event.Start, event.End)
		}
	}

}
