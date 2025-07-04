package main

import (
	"encoding/json"
	"fmt"
	"os"

	"common"
)

// Config holds the configuration for event generation
type Config struct {
	EventRateMean        uint              `json:"eventRateMean"`
	EventRateStd         uint              `json:"eventRateStd"`
	Interval             uint              `json:"interval"`
	MaxNumberOfUsers     uint              `json:"maxNumberOfUsers"`
	DefaultNumberOfUsers uint              `json:"defaultNumberOfUsers"`
	PossibleBetValues    PossibleBetValues `json:"possibleBetValues"`
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	eventFactory := NewEventFactory(&config.PossibleBetValues)
	eventGenerator := NewEventGenerator(config, eventFactory)

	// Read RabbitMQ port from environment variable, default to 5672 if not set
	rabbitPort := os.Getenv("RABBIT_PORT")
	if rabbitPort == "" {
		rabbitPort = "5673"
	}
	rabbitURL := fmt.Sprintf("amqp://guest:guest@localhost:%s/", rabbitPort)
	userQueue := "user_events"
	betQueue := "bet_events"

	userSender, err := NewRabbitMQSender(rabbitURL, userQueue)
	if err != nil {
		fmt.Printf("Error creating user event sender: %v\n", err)
		return
	}
	defer userSender.Close()

	betSender, err := NewRabbitMQSender(rabbitURL, betQueue)
	if err != nil {
		fmt.Printf("Error creating bet event sender: %v\n", err)
		return
	}
	defer betSender.Close()

	var Func = func(events []common.Event) {
		for _, event := range events {
			switch event.GetEventType() {
			case common.EventTypeCreateUser:
				if err := userSender.Send(event); err != nil {
					fmt.Printf("Failed to send user event: %v\n", err)
				}
			default:
				if err := betSender.Send(event); err != nil {
					fmt.Printf("Failed to send bet event: %v\n", err)
				}
			}
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
