package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	eventstore "github.com/fiatjaf/eventstore/badger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nbd-wtf/go-nostr"
)

func main() {
	// Read events from stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Parse JSON into nostr.Event
		var event nostr.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			log.Printf("Failed to parse event JSON: %v\n", err)
			continue
		}

		// Validate the event
		if ok, err := event.CheckSignature(); !ok {
			log.Printf("Invalid event signature: %v (%s)\n", err, event.ID)
			continue
		}

		// Save the event
		if err := common.GetBackend().SaveEvent(nil, &event); err != nil {
			log.Printf("%v (%s)\n", err, event.ID)
			continue
		}

		log.Printf("Imported event %s\n", event.ID)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading from stdin:", err)
	}
}
