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
	betEvent := common.BetEvent{EventID: 1, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)

	beh := &BetEventHandler{
		leaderboardsRepo: mockRepo,
		leaderboard:      mockLB,
		websocketHandler: nil,
	}
	err := beh.Handle(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mockRepo.StoreBetEventCalled {
		t.Error("expected StoreBetEvent to be called before leaderboard update")
	}
	if !mockLB.UpdateCalled {
		t.Error("expected leaderboard.Update to be called")
	}
	if len(mockRepo.Updates) != 1 {
		t.Errorf("expected 1 repo update, got %d", len(mockRepo.Updates))
	}
}

func TestBetEventHandler_Idempotency(t *testing.T) {
	mockLB := &internal.MockLeaderboard{}
	mockRepo := &repositories.MockLeaderboardsRepo{BetEvents: map[uint]bool{42: true}}
	betEvent := common.BetEvent{EventID: 42, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)

	beh := &BetEventHandler{
		leaderboardsRepo: mockRepo,
		leaderboard:      mockLB,
		websocketHandler: nil,
	}
	err := beh.Handle(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.StoreBetEventCalled {
		t.Error("StoreBetEvent should not be called for already processed event")
	}
	if mockLB.UpdateCalled {
		t.Error("leaderboard.Update should not be called for already processed event")
	}
}

func TestBetEventHandler_StoreBetEventError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{}
	mockRepo := &repositories.MockLeaderboardsRepo{StoreBetEventErr: errors.New("store error")}
	betEvent := common.BetEvent{EventID: 99, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)

	beh := &BetEventHandler{
		leaderboardsRepo: mockRepo,
		leaderboard:      mockLB,
		websocketHandler: nil,
	}
	err := beh.Handle(body)
	if err == nil || err.Error() == "" {
		t.Error("expected error from StoreBetEvent")
	}
	if !mockRepo.StoreBetEventCalled {
		t.Error("expected StoreBetEvent to be called")
	}
	if mockLB.UpdateCalled {
		t.Error("leaderboard.Update should not be called if StoreBetEvent fails")
	}
}

func TestBetEventHandler_UnmarshalError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{}
	mockRepo := &repositories.MockLeaderboardsRepo{}
	body := []byte("notjson")

	beh := &BetEventHandler{
		leaderboardsRepo: mockRepo,
		leaderboard:      mockLB,
		websocketHandler: nil,
	}
	err := beh.Handle(body)
	if err == nil || err.Error() == "" {
		t.Error("expected error for invalid JSON")
	}
}

func TestBetEventHandler_LeaderboardUpdateError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{ReturnErr: errors.New("update error")}
	mockRepo := &repositories.MockLeaderboardsRepo{}
	betEvent := common.BetEvent{EventID: 1, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)

	beh := &BetEventHandler{
		leaderboardsRepo: mockRepo,
		leaderboard:      mockLB,
		websocketHandler: nil,
	}
	err := beh.Handle(body)
	if err == nil || err.Error() == "" {
		t.Error("expected error from leaderboard.Update")
	}
}

func TestBetEventHandler_RepoUpdateError(t *testing.T) {
	mockLB := &internal.MockLeaderboard{ReturnData: []*internal.UpdatedData{{CompetitionID: 1, UserID: 2, Score: 100}}}
	mockRepo := &repositories.MockLeaderboardsRepo{ReturnErr: errors.New("repo error")}
	betEvent := common.BetEvent{EventID: 1, EventType: common.EventTypeBet, UserID: 2, Amount: 100}
	body, _ := json.Marshal(betEvent)

	beh := &BetEventHandler{
		leaderboardsRepo: mockRepo,
		leaderboard:      mockLB,
		websocketHandler: nil,
	}
	err := beh.Handle(body)
	if err == nil || err.Error() == "" {
		t.Error("expected error from repo.Update")
	}
}
