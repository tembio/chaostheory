package main

import (
	"common"
	"testing"
)

func TestLeaderboard_Update(t *testing.T) {
	comp := &Competition{
		ID:        1,
		Name:      "Test Competition",
		ScoreRule: "event_type=='win' ? amount : 0",
		StartTime: "2023-01-01T00:00:00Z",
		EndTime:   "2023-12-31T23:59:59Z",
		Rewards:   []string{"Gold"},
	}

	// Normal case: evaluationResult is a valid float
	mockEval := &MockRuleEvaluator{
		matches: []Match{{Rule: comp.ScoreRule, Result: 100.0}},
	}
	lb := NewLeaderboard(mockEval)
	lb.RegistrerCompetition(comp)

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
		matches: []Match{{Rule: comp.ScoreRule, Result: 0}},
	}
	lbZero := NewLeaderboard(mockEvalZero)
	lbZero.RegistrerCompetition(comp)
	updates, err = lbZero.Update(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("expected 0 updates for evaluationResult==0, got %d", len(updates))
	}

	// Case: event.EventType == common.EventTypeBet
	mockEvalBet := &MockRuleEvaluator{
		matches: []Match{{Rule: comp.ScoreRule, Result: 100.0}},
	}
	lbBet := NewLeaderboard(mockEvalBet)
	lbBet.RegistrerCompetition(comp)
	eventBet := event
	eventBet.EventType = common.EventTypeBet
	updates, err = lbBet.Update(eventBet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("expected 0 updates for EventTypeBet, got %d", len(updates))
	}

	// Case: event.EventType == common.EventTypeLoss
	mockEvalLoss := &MockRuleEvaluator{
		matches: []Match{{Rule: comp.ScoreRule, Result: 150.0}},
	}
	lbLoss := NewLeaderboard(mockEvalLoss)
	lbLoss.RegistrerCompetition(comp)
	eventLoss := event
	eventLoss.EventType = common.EventTypeLoss
	eventLoss.Amount = 150.0
	updates, err = lbLoss.Update(eventLoss)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 1 {
		t.Fatalf("expected 1 update for loss, got %d", len(updates))
	}
	if updates[0].Score != -150.0 {
		t.Errorf("expected Score -150.0 for loss, got %f", updates[0].Score)
	}

	// Case: conversion to float64 fails
	mockEvalErr := &MockRuleEvaluator{
		matches: []Match{{Rule: comp.ScoreRule, Result: "not a number"}},
	}
	lbErr := NewLeaderboard(mockEvalErr)
	lbErr.RegistrerCompetition(comp)
	_, err = lbErr.Update(event)
	if err == nil {
		t.Error("expected error for failed float64 conversion, got nil")
	}
}
