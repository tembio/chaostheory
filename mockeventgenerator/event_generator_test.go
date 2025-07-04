package main

import (
	"testing"
	"time"

	"common"
)

func getTestConfig() *Config {
	return &Config{
		EventRateMean:        3,
		EventRateStd:         1,
		Interval:             1,
		MaxNumberOfUsers:     5,
		DefaultNumberOfUsers: 2,
		PossibleBetValues:    *getTestPossibleBetValues(),
	}
}

func TestGenerateRandomEvents_UserCreation(t *testing.T) {
	factory := NewEventFactory(getTestPossibleBetValues())
	gen := NewEventGenerator(getTestConfig(), factory)

	// Should create a user if none exist
	events := gen.GenerateRandomEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(common.UserEvent); !ok {
		t.Errorf("expected UserEvent, got %T", events[0])
	}
}

func TestGenerateRandomEvents_BetEvents(t *testing.T) {
	factory := NewEventFactory(getTestPossibleBetValues())
	gen := NewEventGenerator(getTestConfig(), factory)

	// Create some users
	for i := 0; i < 3; i++ {
		factory.CreateUserEvent()
	}

	foundBet := false
	foundWinOrLoss := false
	for i := 0; i < 20; i++ {
		events := gen.GenerateRandomEvents()
		if len(events) == 2 {
			if be, ok := events[0].(common.BetEvent); ok && be.EventType == common.EventTypeBet {
				foundBet = true
			}
			if be, ok := events[1].(common.BetEvent); ok && (be.EventType == common.EventTypeWin || be.EventType == common.EventTypeLoss) {
				foundWinOrLoss = true
			}
		}
	}

	if !foundBet {
		t.Error("expected to find at least one BetEvent with EventTypeBet")
	}
	if !foundWinOrLoss {
		t.Error("expected to find at least one BetEvent with EventTypeWin or EventTypeLoss")
	}
}

func TestRunEventGeneration_ProducesEvents(t *testing.T) {
	factory := NewEventFactory(getTestPossibleBetValues())
	gen := NewEventGenerator(getTestConfig(), factory)
	ch := make(chan []common.Event, 2)

	go func() {
		gen.RunEventGeneration(func(events []common.Event) {
			ch <- events
		})
	}()

	time.Sleep(10 * time.Millisecond)

	select {
	case events := <-ch:
		// It's valid for there to be zero or more events, so just check we received a slice
		if events == nil {
			t.Error("expected a slice of events, got nil")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for events")
	}
}
