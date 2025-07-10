package repositories

import (
	"common"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// runInitDBScript runs the leaderboard/init_db.sh script for a given db path
func runInitDBScript(dbPath string) error {
	script := "../init_db.sh"
	cmd := exec.Command("bash", script, dbPath, "noTestData")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("init_db.sh failed: %v, output: %s", err, string(output))
	}
	return nil
}

func TestSQLiteCompetitionsRepository(t *testing.T) {
	dbPath := "test_competitions.db"
	os.Remove(dbPath)

	// Call the init_db.sh script to create the schema for the test DB
	err := runInitDBScript(dbPath)
	if err != nil {
		t.Fatalf("failed to run init_db.sh: %v", err)
	}

	repo, err := NewSQLiteCompetitionsRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	defer func() {
		repo.Close()
		os.Remove(dbPath)
	}()

	rewards := map[string]int{"1": 100, "2": 50}
	comp := &common.Competition{
		Name:      "SQLite Competition",
		ScoreRule: "event_type=='bet' ? amount : 0",
		StartTime: "2025-07-10T00:00:00Z",
		EndTime:   "2025-07-11T00:00:00Z",
		Rewards:   rewards,
	}

	// Create a competition
	id, err := repo.Create(comp)
	if err != nil {
		t.Fatalf("failed to create competition: %v", err)
	}
	if id == 0 {
		t.Errorf("expected non-zero ID, got %d", id)
	}

	// Get all competitions
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("failed to get all: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 competition, got %d", len(all))
	}
	got := all[0]
	if got.Name != comp.Name || got.ScoreRule != comp.ScoreRule || got.StartTime != comp.StartTime || got.EndTime != comp.EndTime {
		t.Errorf("competition fields mismatch: got %+v, want %+v", got, comp)
	}
	if len(got.Rewards) != len(rewards) {
		t.Errorf("rewards length mismatch: got %+v, want %+v", got.Rewards, rewards)
	}
	for k, v := range rewards {
		if got.Rewards[k] != v {
			t.Errorf("rewards value mismatch for key %s: got %d, want %d", k, got.Rewards[k], v)
		}
	}
}
