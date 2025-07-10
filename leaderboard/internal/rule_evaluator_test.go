package internal

import (
	"common"
	"testing"
)

func TestBetRuleEvaluator_AddRuleAndEvaluateRules(t *testing.T) {
	eval := &BetRuleEvaluator{}
	eval.AddRule("event_type == 'bet' ? amount : 0")
	eval.AddRule("game == 'Poker' ? 10 : 0")

	event := common.BetEvent{
		EventID:      1,
		EventType:    common.EventTypeBet,
		UserID:       42,
		Amount:       123.45,
		Currency:     "USD",
		ExchangeRate: 1.0,
		Game:         "Poker",
		Distributor:  "evo",
		Studio:       "StudioX",
		Timestamp:    "2023-10-01T12:00:00Z",
	}

	matches, err := eval.EvaluateRules(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].Rule != "event_type == 'bet' ? amount : 0" {
		t.Errorf("unexpected rule: %s", matches[0].Rule)
	}
	if matches[0].Result != 123.45 {
		t.Errorf("expected result 123.45, got %v", matches[0].Result)
	}
	if matches[1].Rule != "game == 'Poker' ? 10 : 0" {
		t.Errorf("unexpected rule: %s", matches[1].Rule)
	}

	v, ok := matches[1].Result.(int)
	if !ok || float32(v) != 10.0 {
		t.Errorf("expected result 10.0 as float64, got %T %v", matches[1].Result, matches[1].Result)
	}
}

func TestBetRuleEvaluator_EvaluateRules_Error(t *testing.T) {
	eval := &BetRuleEvaluator{}
	eval.AddRule("not a valid expr")
	event := common.BetEvent{}
	_, err := eval.EvaluateRules(event)
	if err == nil {
		t.Error("expected error for invalid rule expression")
	}
}
