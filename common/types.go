package common

// User represents a user entity
type User struct {
	ID    uint
	Score float64
}

// Competition represents a competition entity
type Competition struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	ScoreRule string         `json:"score_rule"`
	StartTime string         `json:"start_time"`
	EndTime   string         `json:"end_time"`
	Rewards   map[string]int `json:"rewards"`
}

// EventType represents the type of event in the system
type EventType string

func (e EventType) String() string {
	return string(e)
}

const (
	EventTypeBet        EventType = "bet"
	EventTypeWin        EventType = "win"
	EventTypeLoss       EventType = "loss"
	EventTypeCreateUser EventType = "create_user"
)

// Event is an interface that defines the methods required for an event
type Event interface {
	GetEventID() uint
	GetEventType() EventType
}

// BetEvent represents a betting event in the system
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

func (b BetEvent) GetEventID() uint {
	return b.EventID
}

func (b BetEvent) GetEventType() EventType {
	return b.EventType
}

// UserEvent represents an event related to a user
type UserEvent struct {
	EventID   uint      `json:"event_id"`
	EventType EventType `json:"event_type"`
	UserID    uint      `json:"user_id"`
}

func (u UserEvent) GetEventID() uint {
	return u.EventID
}

func (u UserEvent) GetEventType() EventType {
	return u.EventType
}
