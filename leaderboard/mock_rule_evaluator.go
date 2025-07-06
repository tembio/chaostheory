package main

import "common"

// MockRuleEvaluator implements RuleEvaluatorInterface for testing
// It allows custom behavior for EvaluateRules

type MockRuleEvaluator struct {
	matches       []Match
	evaluateError error
}

func (m *MockRuleEvaluator) AddRule(rule string) {
}

func (m *MockRuleEvaluator) EvaluateRules(event common.BetEvent) ([]Match, error) {
	return m.matches, m.evaluateError
}
