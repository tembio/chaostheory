package main

// Config holds the configuration for event generation
type Config struct {
	EventRateMean        uint              `json:"eventRateMean"`
	EventRateStd         uint              `json:"eventRateStd"`
	Interval             uint              `json:"interval"`
	MaxNumberOfUsers     uint              `json:"maxNumberOfUsers"`
	DefaultNumberOfUsers uint              `json:"defaultNumberOfUsers"`
	PossibleBetValues    PossibleBetValues `json:"possibleBetValues"`
}

// PossibleBetValues represents the possible values for a bet event
type PossibleBetValues struct {
	Currencies  map[string]float64 `json:"currencies"`
	Games       []string           `json:"games"`
	Distributor []string           `json:"distributor"`
	Studio      []string           `json:"studio"`
}

// EventType represents the type of event
type EventType string

const (
	EventTypeBet        EventType = "bet"
	EventTypeWin        EventType = "win"
	EventTypeLoss       EventType = "loss"
	EventTypeCreateUser EventType = "create_user"
)

type Event interface {
	// GetEventID returns the event ID
	GetEventID() uint
	// GetEventType returns the type of the event
	GetEventType() EventType
}

// BetEvent represents a bet event
type BetEvent struct {
	EventID      uint      `json:"event_id"`
	EventType    EventType `json:"event_type"`
	UserID       uint      `json:"user_id"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	ExchangeRate float64   `json:"exchange_rate"`
	Game         string    `json:"game"`
	Distributor  string    `json:"distributor"`
	Studio       string    `json:"studio"`
	Timestamp    string    `json:"timestamp"`
}

// UserEvent represents a user event with a monotonically increasing id
type UserEvent struct {
	EventID   uint      `json:"event_id"`
	EventType EventType `json:"event_type"`
	UserID    uint      `json:"user_id"`
}

// BetEvent implements Event interface
func (b BetEvent) GetEventID() uint {
	return b.EventID
}

func (b BetEvent) GetEventType() EventType {
	return b.EventType
}

// UserEvent implements Event interface
func (u UserEvent) GetEventID() uint {
	return u.EventID
}

func (u UserEvent) GetEventType() EventType {
	return u.EventType
}
