package main

import (
	"math"
	"math/rand"
	"time"
)

// EventGenerator is in charge of generating events based on the configuration provided
type EventGenerator struct {
	EventRateMean        uint
	EventRateStd         uint
	Interval             uint
	MaxNumberOfUsers     uint
	DefaultNumberOfUsers uint

	eventFactory          *EventFactory
	sentEventsPerInterval uint
}

// CreateUserEvent constructor for EventGenerator
func NewEventGenerator(config *Config, eventFactory *EventFactory) *EventGenerator {
	return &EventGenerator{
		EventRateMean:        config.EventRateMean,
		EventRateStd:         config.EventRateStd,
		Interval:             config.Interval,
		MaxNumberOfUsers:     config.MaxNumberOfUsers,
		DefaultNumberOfUsers: config.DefaultNumberOfUsers,

		eventFactory:          eventFactory,
		sentEventsPerInterval: 0,
	}
}

// RunEventGeneration runs GenerateEvents at intervals, producing a variable number of events per interval
// The number of events is normally distributed around the mean and standard deviation provided in the config
// The sendFunc is called with the generated events
func (eg *EventGenerator) RunEventGeneration(sendFunc func([]Event)) {
	mean := float64(eg.EventRateMean)
	std := float64(eg.EventRateStd)

	for {
		numEvents := int(math.Round(rand.NormFloat64()*std + mean))
		if numEvents < 1 {
			numEvents = 0
		}
		var allEvents []Event
		for i := 0; i < numEvents; i++ {
			events := eg.GenerateRandomEvents()
			allEvents = append(allEvents, events...)
		}

		sendFunc(allEvents)

		time.Sleep(time.Duration(eg.Interval) * time.Millisecond)
	}

}

// GenerateEvents generates bet and user events, with a 20% chance of creating a new user event
// If no users have been created yet, it will create at least one new user event
func (eg *EventGenerator) GenerateRandomEvents() []Event {
	numCreatedUsers := eg.eventFactory.GetUserCount()

	// If no users have been created yet, create at least one new user event
	if numCreatedUsers == 0 {
		return []Event{eg.eventFactory.CreateUserEvent()}
	}

	// 20% chance to create a new user event
	if numCreatedUsers < eg.MaxNumberOfUsers && rand.Intn(5) == 0 {
		return []Event{eg.eventFactory.CreateUserEvent()}
	}

	// Pick a random user
	randUserID := rand.Intn(int(numCreatedUsers)) // user IDs start at 0

	return eg.createPairOfBetEvents(uint(randUserID))
}

// createPairOfBetEvents creates a random bet event and a matching win or loss event for the given user_id
func (eg *EventGenerator) createPairOfBetEvents(userID uint) []Event {
	var events []Event
	betEvent := eg.eventFactory.CreateBetEvent(userID, EventTypeBet)
	events = append(events, betEvent)

	winOrLossTypes := []EventType{EventTypeWin, EventTypeLoss}
	winOrLossType := winOrLossTypes[rand.Intn(len(winOrLossTypes))]
	betResultEvent := eg.eventFactory.CreateBetEvent(userID, winOrLossType)

	betResultEvent.Amount = rand.Float64()*1000 + 1
	betResultEvent.Currency = betEvent.Currency
	betResultEvent.ExchangeRate = betEvent.ExchangeRate
	betResultEvent.Game = betEvent.Game
	betResultEvent.Distributor = betEvent.Distributor
	betResultEvent.Studio = betEvent.Studio

	events = append(events, betResultEvent)

	return events
}
