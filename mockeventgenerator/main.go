package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitURL := fmt.Sprintf("amqp://guest:guest@rabbitleaderboard:%s/", rabbitPort)
	userQueue := "user_events"
	betQueue := "bet_events"

	var userSender *RabbitMQSender
	for {
		userSender, err = NewRabbitMQSender(rabbitURL, userQueue)
		if err == nil {
			break
		}
		fmt.Printf("RabbitMQ not ready, retrying in 100ms: %v\n", err)
		time.Sleep(100 * time.Millisecond)
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
			eventType := event.GetEventType()
			switch eventType {
			case common.EventTypeCreateUser:
				if err := userSender.Send(event, eventType.String()); err != nil {
					fmt.Printf("Failed to send user event: %v\n", err)
				}
			default:
				if err := betSender.Send(event, eventType.String()); err != nil {
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
