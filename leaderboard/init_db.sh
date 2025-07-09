#!/bin/sh
# This script creates the Scores table in the SQLite database if it does not exist.

DB_PATH=${1:-leaderboard.db}

# Create the SQLite database file if it does not exist
if [ ! -f "$DB_PATH" ]; then
    sqlite3 "$DB_PATH" ".databases"
fi

sqlite3 "$DB_PATH" <<EOF
CREATE TABLE IF NOT EXISTS Scores (
    competition_id INTEGER,
    user_id INTEGER,
    score REAL,
    PRIMARY KEY (competition_id, user_id)
);
EOF
