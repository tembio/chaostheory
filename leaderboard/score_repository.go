package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// ScoreRepository defines the interface for updating user scores per competition
// This allows for different implementations (e.g., in-memory, database, etc.)
type ScoreRepository interface {
	Update(competitionID, userID uint, score float64)
}

// SQLiteScoreRepository implements ScoreRepository using a SQLite database
type SQLiteScoreRepository struct {
	db *sql.DB
}

// NewSQLiteScoreRepository opens (or creates) a SQLite DB and ensures the Scores table exists
func NewSQLiteScoreRepository(dbPath string) (*SQLiteScoreRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteScoreRepository{db: db}, nil
}

// GetAllScores retrieves all scores for all competitions. Returns a competitionsResults map.
// competitionsResults is assumed to be map[uint][]User, where User has ID and Score fields.
func (sr *SQLiteScoreRepository) GetAllScores() (map[uint][]User, error) {
	rows, err := sr.db.Query(`SELECT competition_id, user_id, score FROM Scores`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	competitionsResults := make(map[uint][]User)
	for rows.Next() {
		var competitionID, userID uint
		var score float64
		if err := rows.Scan(&competitionID, &userID, &score); err != nil {
			return nil, err
		}
		user := User{ID: userID, Score: score}
		competitionsResults[competitionID] = append(competitionsResults[competitionID], user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return competitionsResults, nil
}

// Update inserts or updates the score for a user in a competition
func (sr *SQLiteScoreRepository) Update(competitionID, userID uint, score float64) error {
	_, err := sr.db.Exec(
		`INSERT INTO Scores (competition_id, user_id, score) VALUES (?, ?, ?)
		ON CONFLICT(competition_id, user_id) DO UPDATE SET score=excluded.score;`,
		competitionID, userID, score,
	)
	if err != nil {
		return err
	}
	return nil
}
