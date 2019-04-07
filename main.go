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
}
