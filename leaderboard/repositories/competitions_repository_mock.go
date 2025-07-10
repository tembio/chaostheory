package repositories

import "common"

// MockCompetitions is a mock implementation of CompetitionsRepository for testing
type MockCompetitions struct {
	LastCreated *common.Competition
	LastID      uint
	CreateErr   error
}

// Create inserts a new competition and returns the ID
func (m *MockCompetitions) Create(c *common.Competition) (uint, error) {
	m.LastCreated = c
	m.LastID++
	return m.LastID, m.CreateErr
}

// GetAll retrieves all competitions, returning an empty slice and nil error
func (m *MockCompetitions) GetAll() ([]*common.Competition, error) { return nil, nil }

// Close is a no-op for the mock implementation
func (m *MockCompetitions) Close() {}
