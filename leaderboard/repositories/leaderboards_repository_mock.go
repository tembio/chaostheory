package repositories

import "common"

// MockLeaderboardsRepo is a mock implementation of LeaderboardsRepository for testing purposes
type MockLeaderboardsRepo struct {
	Updates []struct {
		compID, userID uint
		score          float64
	}
	AllUsers    map[uint][]common.User
	TopNUsers   []*common.User
	GetTopNFunc func(competitionID uint, n int) ([]*common.User, error)
	BetEvents   map[uint]bool

	ReturnErr           error
	StoreBetEventErr    error
	StoreBetEventCalled bool
	LastStoredBetEvent  *common.BetEvent
}

// Update appends the update to the mock's updates slice and returns the configured error
func (m *MockLeaderboardsRepo) Update(competitionID, userID uint, score float64) error {
	m.Updates = append(m.Updates, struct {
		compID, userID uint
		score          float64
	}{competitionID, userID, score})
	return m.ReturnErr
}

// GetAll returns an empty map and nil error for the mock implementation
func (m *MockLeaderboardsRepo) GetAll() (map[uint][]common.User, error) {
	return m.AllUsers, m.ReturnErr
}

// GetTopN returns an empty slice and nil error for the mock implementation
func (m *MockLeaderboardsRepo) GetTopN(competitionID uint, n int) ([]*common.User, error) {
	if m.GetTopNFunc != nil {
		return m.GetTopNFunc(competitionID, n)
	}
	return m.TopNUsers, m.ReturnErr
}

// HasBetEvent returns true if the eventID is in the BetEvents map
func (m *MockLeaderboardsRepo) HasBetEvent(eventID uint) (bool, error) {
	if m.BetEvents == nil {
		return false, nil
	}
	return m.BetEvents[eventID], nil
}

// StoreBetEvent records the eventID in the BetEvents map and returns configured error
func (m *MockLeaderboardsRepo) StoreBetEvent(event *common.BetEvent) error {
	m.StoreBetEventCalled = true
	m.LastStoredBetEvent = event
	if m.BetEvents == nil {
		m.BetEvents = make(map[uint]bool)
	}
	m.BetEvents[event.EventID] = true
	return m.StoreBetEventErr
}
