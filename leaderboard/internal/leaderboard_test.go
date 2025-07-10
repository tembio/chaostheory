package internal

import (
	"common"
	"testing"
)

func TestLeaderboard_RegisterCompetition(t *testing.T) {
	comp := &common.Competition{
		ID:        1,
		Name:      "Test Competition",
		ScoreRule: "event_type=='win' ? amount : 0",
		StartTime: "2023-01-01T00:00:00Z",
		EndTime:   "2023-12-31T23:59:59Z",
		Rewards:   map[string]int{"1-2": 100, "3-5": 50, "6+": 25},
	}

	mockEval := &MockRuleEvaluator{}
	lb := NewLeaderboard(mockEval)

	// Register a valid competition
	lb.RegisterCompetition(comp)
	if _, exists := lb.rulesToCompetition[comp.ScoreRule]; !exists {
		t.Errorf("expected competition to be registered")
	}

	// Register nil competition
	lb.RegisterCompetition(nil)
	// Should not panic or add anything

	// Register competition with empty ScoreRule
	compEmpty := &common.Competition{ID: 2, Name: "No Rule"}
	lb.RegisterCompetition(compEmpty)
	if _, exists := lb.rulesToCompetition[compEmpty.ScoreRule]; exists {
		t.Errorf("should not register competition with empty ScoreRule")
	}

	// Register duplicate ScoreRule
	lb.RegisterCompetition(comp)
	// Should not add duplicate
	count := 0
	for _, id := range lb.rulesToCompetition {
		if id == comp.ID {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected only one registration for the same ScoreRule, got %d", count)
	}
}

func TestLeaderboard_Update(t *testing.T) {
	comp := &common.Competition{
		ID:        1,
		Name:      "Test Competition",
		ScoreRule: "event_type=='win' ? amount : 0",
		StartTime: "2023-01-01T00:00:00Z",
		EndTime:   "2023-12-31T23:59:59Z",
		Rewards:   map[string]int{"1-2": 100, "3-5": 50, "6+": 25},
	}

	// Normal case: evaluationResult is a valid float
	mockEval := &MockRuleEvaluator{
		Matches: []Match{{Rule: comp.ScoreRule, Result: 100.0}},
	}
	lb := NewLeaderboard(mockEval)
	lb.RegisterCompetition(comp)

	event := common.BetEvent{
		EventID:      1,
		EventType:    common.EventTypeWin,
		UserID:       42,
		Amount:       100.0,
		Currency:     "USD",
		ExchangeRate: 1.0,
		Game:         "Poker",
		Distributor:  "evo",
		Studio:       "StudioX",
		Timestamp:    "2023-10-01T12:00:00Z",
	}
	updates, err := lb.Update(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 1 {
		t.Fatalf("expected 1 update, got %d", len(updates))
	}
	if updates[0].UserID != 42 {
		t.Errorf("expected UserID 42, got %d", updates[0].UserID)
	}
	if updates[0].Score != 100.0 {
		t.Errorf("expected Score 100.0, got %f", updates[0].Score)
	}

	// Case: evaluationResult == 0
	mockEvalZero := &MockRuleEvaluator{
		Matches: []Match{{Rule: comp.ScoreRule, Result: 0}},
	}
	lbZero := NewLeaderboard(mockEvalZero)
	lbZero.RegisterCompetition(comp)
	updates, err = lbZero.Update(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("expected 0 updates for evaluationResult==0, got %d", len(updates))
	}

	// Case: event.EventType == common.EventTypeBet
	mockEvalBet := &MockRuleEvaluator{
		Matches: []Match{{Rule: comp.ScoreRule, Result: 100.0}},
	}
	lbBet := NewLeaderboard(mockEvalBet)
	lbBet.RegisterCompetition(comp)
	eventBet := event
	eventBet.EventType = common.EventTypeBet
	updates, err = lbBet.Update(eventBet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 1 {
		t.Errorf("expected 1 update for EventTypeBet, got %d", len(updates))
	}
	if updates[0].Score != eventBet.Amount {
		t.Errorf("expected Score 100 for EventTypeBet, got %f", updates[0].Score)
	}

	// Case: event.EventType == common.EventTypeLoss
	mockEvalLoss := &MockRuleEvaluator{
		Matches: []Match{{Rule: comp.ScoreRule, Result: 150.0}},
	}
	lbLoss := NewLeaderboard(mockEvalLoss)
	lbLoss.RegisterCompetition(comp)
	eventLoss := event
	eventLoss.EventType = common.EventTypeLoss
	updates, err = lbLoss.Update(eventLoss)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("expected 0 updates for EventTypeLoss, got %d", len(updates))
	}

	// Case: conversion to float64 fails
	mockEvalErr := &MockRuleEvaluator{
		Matches: []Match{{Rule: comp.ScoreRule, Result: "not a number"}},
	}
	lbErr := NewLeaderboard(mockEvalErr)
	lbErr.RegisterCompetition(comp)
	updates, _ = lbErr.Update(event)
	if len(updates) != 0 {
		t.Errorf("expected 0 updates when amount is not a float, got %d", len(updates))
	}
}

func TestLeaderboard_Load(t *testing.T) {
	lb := NewLeaderboard(&MockRuleEvaluator{})
	leaderboards := map[uint][]common.User{
		1: {
			{ID: 1, Score: 10.0},
			{ID: 2, Score: 20.0},
		},
		2: {
			{ID: 3, Score: 30.0},
		},
	}
	lb.Load(leaderboards)

	if len(lb.competitionsResults) != 2 {
		t.Errorf("expected 2 competitions loaded, got %d", len(lb.competitionsResults))
	}
	if user, ok := lb.competitionsResults[1][1]; !ok || user.Score != 10.0 {
		t.Errorf("expected user 1 with score 10.0 in competition 1")
	}
	if user, ok := lb.competitionsResults[2][3]; !ok || user.Score != 30.0 {
		t.Errorf("expected user 3 with score 30.0 in competition 2")
	}
}
