package handlers

import (
	"common"
	"encoding/json"
	"fmt"
	"net/http"

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

func (lh *LeaderboardsHandler) GetLeaderboards(w http.ResponseWriter, r *http.Request) {

}
