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
	start := time.Now()
	err, eventsList := parser.ScrapCalendarEvents()

	fmt.Printf("Parsed %d events in %f seconds\n\n",
		len(eventsList), time.Since(start).Seconds())

	for _, event := range eventsList {
		fmt.Println(event.RawDate)
	}

	err, tags := parser.ScrapEventTags()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tags)

}
