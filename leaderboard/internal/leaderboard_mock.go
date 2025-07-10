package internal

import "common"

// MockLeaderboard implements LeaderboardInterface for testing
type MockLeaderboard struct {
	UpdateCalled bool
	UpdateData   []common.BetEvent
	ReturnData   []*UpdatedData
	ReturnErr    error
}

// Update simulates the Update method of LeaderboardInterface
func (m *MockLeaderboard) Update(event common.BetEvent) ([]*UpdatedData, error) {
	m.UpdateCalled = true
	m.UpdateData = append(m.UpdateData, event)
	return m.ReturnData, m.ReturnErr
}

// RegisterCompetition simulates the RegisterCompetition method of LeaderboardInterface
func (m *MockLeaderboard) RegisterCompetition(comp *common.Competition) {
	// No-op for mock
}
