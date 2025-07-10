package internal

// Config holds the configuration for event generation
type Config struct {
	EventRateMean        uint              `json:"eventRateMean"`
	EventRateStd         uint              `json:"eventRateStd"`
	Interval             uint              `json:"interval"`
	MaxNumberOfUsers     uint              `json:"maxNumberOfUsers"`
	DefaultNumberOfUsers uint              `json:"defaultNumberOfUsers"`
	PossibleBetValues    PossibleBetValues `json:"possibleBetValues"`
}
