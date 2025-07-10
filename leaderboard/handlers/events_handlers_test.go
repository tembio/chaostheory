package handlers

import (
	"common"
	"encoding/json"
	"errors"
	"leaderboard/internal"
	"leaderboard/repositories"
	"testing"
)

func TestBetEventHandler_Success(t *testing.T) {
	mockLB := &internal.MockLeaderboard{
		ReturnData: []*internal.UpdatedData{{CompetitionID: 1, UserID: 2, Score: 100}},
	}
	mockRepo := &repositories.MockLeaderboardsRepo{}
	ack := false
	ackFn := func() { ack = true }
	betEvent := common.BetEvent{EventID: 1, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)
	err := BetEventHandler(body, mockLB, mockRepo, ackFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ack {
		t.Error("expected acknowledgeEvent to be called")
	}
	if !mockLB.UpdateCalled {
		t.Error("expected leaderboard.Update to be called")
	}
	if len(mockRepo.Updates) != 1 {
		t.Errorf("expected 1 repo update, got %d", len(mockRepo.Updates))
	}
}

func TestBetEventHandler_UnmarshalError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{}
	mockRepo := &repositories.MockLeaderboardsRepo{}
	ackFn := func() {}
	body := []byte("notjson")
	err := BetEventHandler(body, mockLB, mockRepo, ackFn)
	if err == nil || err.Error() == "" {
		t.Error("expected error for invalid JSON")
	}
}

func TestBetEventHandler_LeaderboardUpdateError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{ReturnErr: errors.New("update error")}
	mockRepo := &repositories.MockLeaderboardsRepo{}
	ackFn := func() {}
	betEvent := common.BetEvent{EventID: 1, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)
	err := BetEventHandler(body, mockLB, mockRepo, ackFn)
	if err == nil || err.Error() == "" {
		t.Error("expected error from leaderboard.Update")
	}
}

func TestBetEventHandler_RepoUpdateError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{ReturnData: []*internal.UpdatedData{{CompetitionID: 1, UserID: 2, Score: 100}}}
	mockRepo := &repositories.MockLeaderboardsRepo{ReturnErr: errors.New("repo error")}
	ackFn := func() {}
	betEvent := common.BetEvent{EventID: 1, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)
	err := BetEventHandler(body, mockLB, mockRepo, ackFn)
	if err == nil || err.Error() == "" {
		t.Error("expected error from repo.Update")
	}
}

// You can add similar tests for UserEventHandler as needed.
