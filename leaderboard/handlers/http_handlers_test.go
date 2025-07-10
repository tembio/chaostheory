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

	"github.com/gorilla/mux"
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

func TestGetLeaderboardByID_NotFound(t *testing.T) {
	repo := &repositories.MockLeaderboardsRepo{}
	repo.GetTopNFunc = func(competitionID uint, n int) ([]*common.User, error) {
		return []*common.User{}, nil // Simulate leaderboard not found
	}
	h := NewLeaderboardsHandler(repo)
	r := mux.NewRouter()
	r.HandleFunc("/leaderboards/{id}", h.GetLeaderboardByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/leaderboards/123", nil)
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", resp.StatusCode)
	}
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	expected := "leaderboard with id 123 not found or has no users"
	if string(body) == "" {
		t.Errorf("expected error message in body, got empty string")
	} else if string(body) != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, string(body))
	}
}

func TestGetLeaderboardByID_Success(t *testing.T) {
	repo := &repositories.MockLeaderboardsRepo{}
	h := NewLeaderboardsHandler(repo)

	users := []*common.User{{ID: 1, Score: 100}, {ID: 2, Score: 90}}
	repo.TopNUsers = users

	r := mux.NewRouter()
	r.HandleFunc("/leaderboards/{id}", h.GetLeaderboardByID)
	req := httptest.NewRequest("GET", "/leaderboards/1?count=2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	var got []common.User
	json.NewDecoder(resp.Body).Decode(&got)
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Errorf("unexpected users: %+v", got)
	}
}

func TestGetLeaderboardByID_BadID(t *testing.T) {
	repo := &repositories.MockLeaderboardsRepo{}
	h := NewLeaderboardsHandler(repo)
	r := mux.NewRouter()
	r.HandleFunc("/leaderboards/{id}", h.GetLeaderboardByID)
	req := httptest.NewRequest("GET", "/leaderboards/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetLeaderboardByID_RepoError(t *testing.T) {
	repo := &repositories.MockLeaderboardsRepo{}
	repo.ReturnErr = errTest
	h := NewLeaderboardsHandler(repo)
	r := mux.NewRouter()
	r.HandleFunc("/leaderboards/{id}", h.GetLeaderboardByID)
	req := httptest.NewRequest("GET", "/leaderboards/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

var errTest = &mockError{"repo error"}

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }
