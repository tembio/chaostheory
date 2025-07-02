package main

import (
	"testing"
	"time"
)

func getTestPossibleBetValues() *PossibleBetValues {
	return &PossibleBetValues{
		Currencies:  map[string]float64{"USD": 1.0, "EUR": 0.9},
		Games:       []string{"Poker", "Blackjack"},
		Distributor: []string{"DistA", "DistB"},
		Studio:      []string{"StudioX", "StudioY"},
	}
}

func TestCreateBetEvent(t *testing.T) {
	factory := NewEventFactory(getTestPossibleBetValues())
	userID := uint(42)
	eventType := EventTypeBet
	betEvent := factory.CreateBetEvent(userID, eventType)

	if betEvent.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, betEvent.UserID)
	}
	if betEvent.EventType != eventType {
		t.Errorf("expected EventType %v, got %v", eventType, betEvent.EventType)
	}
	if betEvent.Amount < 1.0 || betEvent.Amount > 1001.0 {
		t.Errorf("amount %f out of range", betEvent.Amount)
	}
	if _, ok := factory.PossibleBetValues.Currencies[betEvent.Currency]; !ok {
		t.Errorf("unexpected currency: %s", betEvent.Currency)
	}
	if betEvent.ExchangeRate != factory.PossibleBetValues.Currencies[betEvent.Currency] {
		t.Errorf("exchange rate mismatch for currency %s", betEvent.Currency)
	}
	foundGame := false
	for _, g := range factory.PossibleBetValues.Games {
		if betEvent.Game == g {
			foundGame = true
			break
		}
	}
	if !foundGame {
		t.Errorf("unexpected game: %s", betEvent.Game)
	}
	foundDist := false
	for _, d := range factory.PossibleBetValues.Distributor {
		if betEvent.Distributor == d {
			foundDist = true
			break
		}
	}
	if !foundDist {
		t.Errorf("unexpected distributor: %s", betEvent.Distributor)
	}
	foundStudio := false
	for _, s := range factory.PossibleBetValues.Studio {
		if betEvent.Studio == s {
			foundStudio = true
			break
		}
	}
	if !foundStudio {
		t.Errorf("unexpected studio: %s", betEvent.Studio)
	}
	if _, err := time.Parse(time.RFC3339, betEvent.Timestamp); err != nil {
		t.Errorf("invalid timestamp: %s", betEvent.Timestamp)
	}
}

func TestCreateUserEvent(t *testing.T) {
	factory := NewEventFactory(getTestPossibleBetValues())
	userEvent1 := factory.CreateUserEvent()
	userEvent2 := factory.CreateUserEvent()

	if userEvent1.EventID == 0 {
		t.Errorf("unexpected userEvent1: %+v", userEvent1)
	}
	if userEvent1.EventType != EventTypeCreateUser {
		t.Errorf("expected EventType %v, got %v", EventTypeCreateUser, userEvent1.EventType)
	}
	if userEvent2.EventID != userEvent1.EventID+1 {
		t.Errorf("unexpected userEvent2: %+v", userEvent2)
	}
	if userEvent2.UserID != userEvent1.UserID+1 {
		t.Errorf("unexpected userEvent2: %+v", userEvent2)
	}
	if userEvent2.EventType != EventTypeCreateUser {
		t.Errorf("expected EventType %v, got %v", EventTypeCreateUser, userEvent2.EventType)
	}
}

func TestGetUserCount(t *testing.T) {
	factory := NewEventFactory(getTestPossibleBetValues())
	if factory.GetUserCount() != 0 {
		t.Errorf("expected 0 users, got %d", factory.GetUserCount())
	}
	factory.CreateUserEvent()
	factory.CreateUserEvent()
	if factory.GetUserCount() != 2 {
		t.Errorf("expected 2 users, got %d", factory.GetUserCount())
	}
}
