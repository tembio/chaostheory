package repositories

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"common"
)

// LeaderboardsRepository defines the interface for updating user scores per competition
// This allows for different implementations (e.g., in-memory, database, etc.)
type LeaderboardsRepository interface {
	Update(competitionID, userID uint, score float64) error
	GetAll() (map[uint][]common.User, error)
	GetTopN(competitionID uint, n int) ([]*common.User, error)
}

// SQLiteLeaderboards implements LeaderboardsRepository using a SQLite database
type SQLiteLeaderboards struct {
	db *sql.DB
}

// NewSQLiteLeaderboardsRepository opens (or creates) a SQLite DB and ensures the Leaderboards table exists
func NewSQLiteLeaderboardsRepository(dbPath string) (*SQLiteLeaderboards, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteLeaderboards{db: db}, nil
}

// GetAll retrieves all scores for all competitions. Returns a competitionsResults map.
// competitionsResults is assumed to be map[uint][]User, where User has ID and Score fields.
func (sr *SQLiteLeaderboards) GetAll() (map[uint][]common.User, error) {
	rows, err := sr.db.Query(`SELECT competition_id, user_id, score FROM Leaderboards`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	competitionsResults := make(map[uint][]common.User)
	for rows.Next() {
		var competitionID, userID uint
		var score float64
		if err := rows.Scan(&competitionID, &userID, &score); err != nil {
			return nil, err
		}
		user := common.User{ID: userID, Score: score}
		competitionsResults[competitionID] = append(competitionsResults[competitionID], user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return competitionsResults, nil
}

// Update inserts or updates the score for a user in a competition
func (sr *SQLiteLeaderboards) Update(competitionID, userID uint, score float64) error {
	_, err := sr.db.Exec(
		`INSERT INTO Leaderboards (competition_id, user_id, score) VALUES (?, ?, ?)
		ON CONFLICT(competition_id, user_id) DO UPDATE SET score=excluded.score;`,
		competitionID, userID, score,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetTopN retrieves the top N users for a given competition, ordered by greatest score
func (sr *SQLiteLeaderboards) GetTopN(competitionID uint, n int) ([]*common.User, error) {
	rows, err := sr.db.Query(`SELECT user_id, score FROM Leaderboards WHERE competition_id = ? ORDER BY score DESC LIMIT ?`, competitionID, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*common.User
	for rows.Next() {
		var user common.User
		if err := rows.Scan(&user.ID, &user.Score); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Close closes the SQLite database connection
func (sr *SQLiteLeaderboards) Close() {
	if sr.db != nil {
		sr.db.Close()
	}
}
