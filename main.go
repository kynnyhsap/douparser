package main

import (
	"dou-parser/events"
	"dou-parser/parser"
	"fmt"
	"log"
	"time"
)

func printEvents(list []events.Event) {
	for i, event := range list {
		fmt.Printf(`
%d. 🆔 %d 🔴 %s
	💰 %s
	⚓ %s
	📅 %s (%s  %s)
	🗒️ %s
	🏷️ %v
--------------------------------------------------------------------------------------------------------
		`, i, event.ID,
			event.Title,
			event.Cost,
			event.Location,
			event.RawDate,
			event.Start,
			event.End,
			event.ShortDescription,
			event.Tags)
	}
}

func main() {
	var p parser.EventsParser

	start := time.Now()
	err := p.ParseAll()
	fmt.Printf("Parsed %d events in %f seconds\n\n", len(p.Events), time.Since(start).Seconds())

	for _, event := range p.Events {
		fmt.Println(event.RawDate)
	}
	//printEvents(p.Events)

	err = p.ParseTags()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p.Tags)

}
