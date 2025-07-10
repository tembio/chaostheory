package internal

import (
	"math/rand"
	"time"

	"common"
)

// PossibleBetValues represents the possible values for a bet event
type PossibleBetValues struct {
	Currencies  map[string]float64 `json:"currencies"`
	Games       []string           `json:"games"`
	Distributor []string           `json:"distributor"`
	Studio      []string           `json:"studio"`
}

type EventFactory struct {
	PossibleBetValues *PossibleBetValues

	eventIDCounter uint
	userCounter    uint
}

// CreateBetEvent creates a bet event with random values
func (ef *EventFactory) CreateBetEvent(userID uint, eventType common.EventType) common.BetEvent {
	possibleValues := ef.PossibleBetValues

	amount := rand.Float64()*1000 + 1
	currencies := make([]string, 0, len(possibleValues.Currencies))
	for k := range possibleValues.Currencies {
		currencies = append(currencies, k)
	}
	currency := currencies[rand.Intn(len(currencies))]
	exchangeRate := possibleValues.Currencies[currency]
	game := possibleValues.Games[rand.Intn(len(possibleValues.Games))]
	distributor := possibleValues.Distributor[rand.Intn(len(possibleValues.Distributor))]
	studio := possibleValues.Studio[rand.Intn(len(possibleValues.Studio))]
	timestamp := time.Now().UTC().Format(time.RFC3339)

	ef.eventIDCounter++
	betEvent := common.BetEvent{
		EventID:      ef.eventIDCounter,
		EventType:    eventType,
		UserID:       userID,
		Amount:       amount,
		Currency:     currency,
		ExchangeRate: exchangeRate,
		Game:         game,
		Distributor:  distributor,
		Studio:       studio,
		Timestamp:    timestamp,
	}

	return betEvent
}

// CreateUserEvent returns a UserEvent with a monotonically increasing id
func (ef *EventFactory) CreateUserEvent() common.UserEvent {
	ef.userCounter++
	ef.eventIDCounter++
	userEvent := common.UserEvent{EventID: ef.eventIDCounter, EventType: common.EventTypeCreateUser, UserID: ef.userCounter}

	return userEvent
}

// GetUserCount returns the number of users created
func (ef *EventFactory) GetUserCount() uint {
	return ef.userCounter
}

// NewEventFactory constructor for EventFactory
func NewEventFactory(possibleBetValues *PossibleBetValues) *EventFactory {
	return &EventFactory{
		PossibleBetValues: possibleBetValues,
		eventIDCounter:    0,
		userCounter:       0,
	}
}
