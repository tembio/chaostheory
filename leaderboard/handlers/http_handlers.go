package handlers

import (
	"common"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"leaderboard/repositories"
)

// CompetitionsHandler holds dependencies for competition handlers
type CompetitionsHandler struct {
	competitionsRepo repositories.CompetitionsRepository
}

// NewCompetitionsHandler creates a new CompetitionHandler instance
func NewCompetitionsHandler(repo repositories.CompetitionsRepository) *CompetitionsHandler {
	return &CompetitionsHandler{
		competitionsRepo: repo,
	}
}

// CreateCompetition creates a new competition
func (ch *CompetitionsHandler) CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition common.Competition

	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON"))
		return
	}

	id, err := ch.competitionsRepo.Create(&competition)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to create competition: %v", err)))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

// LeaderboardsHandler holds dependencies for leaderboard handlers
type LeaderboardsHandler struct {
	leaderboardsRepo repositories.LeaderboardsRepository
}

// NewLeaderboardsHandler creates a new LeaderboardHandler instance
func NewLeaderboardsHandler(repo repositories.LeaderboardsRepository) *LeaderboardsHandler {
	return &LeaderboardsHandler{
		leaderboardsRepo: repo,
	}
}

// GetLeaderboardByID retrieves the top N users for a given competition ID
func (lh *LeaderboardsHandler) GetLeaderboardByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	countStr := r.URL.Query().Get("count")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid leaderboard id"))
		return
	}
	count := 10 // default
	if countStr != "" {
		c, err := strconv.Atoi(countStr)
		if err == nil && c > 0 {
			count = c
		}
	}

	users, err := lh.leaderboardsRepo.GetTopN(uint(id), count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get leaderboard: %v", err)))
		return
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("leaderboard with id %d not found or has no users", id)))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
