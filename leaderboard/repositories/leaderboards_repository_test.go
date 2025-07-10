package repositories

import (
	"os"
	"testing"
)

func TestSQLiteLeaderboardsRepository_UpdateAndGetAll(t *testing.T) {
	dbPath := "test_leaderboards.db"

	// Call the init_db.sh script to create the schema for the test DB
	err := runInitDBScript(dbPath)
	if err != nil {
		t.Fatalf("failed to run init_db.sh: %v", err)
	}

	repo, err := NewSQLiteLeaderboardsRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create SQLiteLeaderboardsRepository: %v", err)
	}
	defer func() {
		repo.Close()
		os.Remove(dbPath)
	}()

	// Update scores for different competitions and users
	err = repo.Update(1, 10, 100.0)
	if err != nil {
		t.Errorf("failed to update score: %v", err)
	}
	err = repo.Update(1, 20, 200.0)
	if err != nil {
		t.Errorf("failed to update score: %v", err)
	}
	err = repo.Update(2, 10, 300.0)
	if err != nil {
		t.Errorf("failed to update score: %v", err)
	}

	// GetAll retrieves all leaderboards
	scores, err := repo.GetAll()
	if err != nil {
		t.Fatalf("failed to get all scores: %v", err)
	}

	if len(scores[1]) != 2 {
		t.Errorf("expected 2 users in competition 1, got %d", len(scores[1]))
	}
	if len(scores[2]) != 1 {
		t.Errorf("expected 1 user in competition 2, got %d", len(scores[2]))
	}

	var found10, found20 bool
	for _, user := range scores[1] {
		if user.ID == 10 && user.Score == 100.0 {
			found10 = true
		}
		if user.ID == 20 && user.Score == 200.0 {
			found20 = true
		}
	}
	if !found10 {
		t.Errorf("user 10 with score 100.0 not found in competition 1")
	}
	if !found20 {
		t.Errorf("user 20 with score 200.0 not found in competition 1")
	}

	if scores[2][0].ID != 10 || scores[2][0].Score != 300.0 {
		t.Errorf("expected user 10 with score 300.0 in competition 2, got ID=%d, Score=%v", scores[2][0].ID, scores[2][0].Score)
	}
}
