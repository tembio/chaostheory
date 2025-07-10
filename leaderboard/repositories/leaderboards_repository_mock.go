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
	ReturnErr   error
	GetTopNFunc func(competitionID uint, n int) ([]*common.User, error)
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
