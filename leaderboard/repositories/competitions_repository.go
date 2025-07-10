package repositories

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"

	"common"
)

// CompetitionsRepository defines the interface for managing competitions
// This allows for different implementations (e.g., in-memory, database, etc.)
type CompetitionsRepository interface {
	Create(competition *common.Competition) (uint, error)
	GetAll() ([]*common.Competition, error)
	Close()
}

// SQLiteCompetitions implements CompetitionsRepository using SQLite
type SQLiteCompetitions struct {
	db *sql.DB
}

// NewSQLiteCompetitionsRepository opens (or creates) a SQLite DB
func NewSQLiteCompetitionsRepository(dbPath string) (*SQLiteCompetitions, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteCompetitions{db: db}, nil
}

// Create inserts a new competition and returns the ID
func (r *SQLiteCompetitions) Create(competition *common.Competition) (uint, error) {
	// Serialize Rewards map to JSON
	rewardsJSON, err := json.Marshal(competition.Rewards)
	if err != nil {
		return 0, err
	}
	res, err := r.db.Exec(`INSERT INTO Competitions (name, scorerule, starttime, endtime, rewards) VALUES (?, ?, ?, ?, ?)`,
		competition.Name, competition.ScoreRule, competition.StartTime, competition.EndTime, string(rewardsJSON))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GetAll retrieves all competitions, including all fields and deserializes Rewards
func (r *SQLiteCompetitions) GetAll() ([]*common.Competition, error) {
	rows, err := r.db.Query(`SELECT id, name, scorerule, starttime, endtime, rewards FROM Competitions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var competitions []*common.Competition
	for rows.Next() {
		var c common.Competition
		var rewardsJSON string
		if err := rows.Scan(&c.ID, &c.Name, &c.ScoreRule, &c.StartTime, &c.EndTime, &rewardsJSON); err != nil {
			return nil, err
		}
		if rewardsJSON != "" {
			if err := json.Unmarshal([]byte(rewardsJSON), &c.Rewards); err != nil {
				return nil, err
			}
		}
		competitions = append(competitions, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return competitions, nil
}

// Close closes the SQLite database connection
func (r *SQLiteCompetitions) Close() {
	if r.db != nil {
		r.db.Close()
	}
}
