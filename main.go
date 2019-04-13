package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	var parser EventsParser

	start := time.Now()

	err := parser.ParseAll()
	if err != nil {
		log.Fatal(err) // =(
	}

	fmt.Printf("Parsed %d events in %f seconds\n", len(parser.Events), time.Since(start).Seconds())

	for i, event := range parser.Events {
		if event.End.IsZero() {
			fmt.Printf("%d(%s): %s \n", i+1, event.RawDate, event.Start)
		} else {
			fmt.Printf("%d(%s): %s - %s \n", i+1, event.RawDate, event.Start, event.End)
		}
	}

}
