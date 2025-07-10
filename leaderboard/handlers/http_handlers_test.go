package handlers

import (
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"leaderboard/repositories"
)

func TestCreateCompetitionHandler(t *testing.T) {
	repo := &repositories.MockCompetitions{}
	h := NewCompetitionsHandler(repo)
	competition := common.Competition{
		Name:      "Test Comp",
		ScoreRule: "rule",
		StartTime: "2025-07-10T00:00:00Z",
		EndTime:   "2025-07-11T00:00:00Z",
		Rewards:   map[string]int{"1": 100},
	}
	body, _ := json.Marshal(competition)
	req := httptest.NewRequest("POST", "/competitions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateCompetition(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte("id")) {
		t.Errorf("expected response to contain id, got %s", string(respBody))
	}
}

func TestCreateCompetitionHandler_BadJSON(t *testing.T) {
	repo := &repositories.MockCompetitions{}
	h := NewCompetitionsHandler(repo)
	req := httptest.NewRequest("POST", "/competitions", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()
	h.CreateCompetition(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateCompetitionHandler_CreateError(t *testing.T) {
	repo := &repositories.MockCompetitions{CreateErr: fmt.Errorf("create error")}
	h := NewCompetitionsHandler(repo)
	competition := common.Competition{Name: "Test Comp"}
	body, _ := json.Marshal(competition)
	req := httptest.NewRequest("POST", "/competitions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateCompetition(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}
