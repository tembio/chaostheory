package internal

import "common"

// MockRuleEvaluator implements RuleEvaluatorInterface for testing
type MockRuleEvaluator struct {
	Matches       []Match
	EvaluateError error
}

// AddRule is a no-op for the mock implementation
func (m *MockRuleEvaluator) AddRule(rule string) {
}

// EvaluateRules simulates rule evaluation by returning predefined matches
func (m *MockRuleEvaluator) EvaluateRules(event common.BetEvent) ([]Match, error) {
	return m.Matches, m.EvaluateError
}
