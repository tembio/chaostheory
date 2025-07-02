package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	eventFactory := NewEventFactory(&config.PossibleBetValues)
	eventGenerator := NewEventGenerator(config, eventFactory)

	var Func = func(events []Event) {
		for _, event := range events {
			fmt.Printf("Bet Events: %d, %s\n", event.GetEventID(), event.GetEventType())
		}
	}

	eventGenerator.RunEventGeneration(Func)
}

// LoadConfig reads the EventGenerator config from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
