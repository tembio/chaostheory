package internal

import (
	"common"
	"encoding/json"
	"time"
)

// MockReceiver is a mock implementation of the Receiver interface for testing.
type MockReceiver struct {
	Events [][]byte // Pre-serialized event payloads to deliver
}

// Receive calls the handler with each mock event.
func (m *MockReceiver) Receive(handler func([]byte) error) error {
	for _, event := range m.Events {
		time.Sleep(2 * time.Second) // wait for rabbit to start

		if err := handler(event); err != nil {
			return err
		}
	}
	return nil
}

// NewMockEventReceiver creates a mock receiver with predefined bet events.
func NewMockEventReceiver() *MockReceiver {
	events := []common.BetEvent{
		{EventID: 1, EventType: "bet", UserID: 42, Amount: 100.0},
		{EventID: 2, EventType: "bet", UserID: 43, Amount: 200.0},
	}
	var payloads [][]byte
	for _, e := range events {
		b, _ := json.Marshal(e)
		payloads = append(payloads, b)
	}
	return &MockReceiver{Events: payloads}
}
