#!/bin/sh
# This script creates the tables in the SQLite database if they does not exist.
# also prloads some initial data for competitions.


DB_PATH=${1:-leaderboard.db}
GENERATE_TEST_DATA=${2:-1} # 1 = generate test data, 0 = skip

# Create the SQLite database file if it does not exist
if [ ! -f "$DB_PATH" ]; then
    sqlite3 "$DB_PATH" ".databases"
fi

sqlite3 "$DB_PATH" <<EOF
CREATE TABLE IF NOT EXISTS Leaderboards (
    competition_id INTEGER,
    user_id INTEGER,
    score REAL,
    PRIMARY KEY (competition_id, user_id)
);

CREATE TABLE IF NOT EXISTS BetEvents (
    event_id INTEGER PRIMARY KEY,
    user_id INTEGER,
    amount REAL
);

CREATE TABLE IF NOT EXISTS Competitions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    scorerule TEXT,
    starttime TEXT,
    endtime TEXT,
    rewards TEXT
);
EOF

if [ "$GENERATE_TEST_DATA" != "noTestData" ]; then
sqlite3 "$DB_PATH" <<EOF
INSERT INTO Competitions (id, name, scorerule, starttime, endtime, rewards) VALUES (
    1,
    'Monthly Challenge',
    'event_type==''bet'' && distributor==''evo'' ? amount : 0',
    '2023-10-01T00:00:00Z',
    '2023-10-31T23:59:59Z',
    '{"1-2":100,"3-5":50,"6+":25}'
)
ON CONFLICT(id) DO NOTHING;

INSERT INTO Competitions (id, name, scorerule, starttime, endtime, rewards) VALUES (
    2,
    'Weekly Sprint',
    'event_type==''bet'' && game==''Poker'' ? amount : 0',
    '2023-10-01T00:00:00Z',
    '2023-10-31T23:59:59Z',
    '{"1":300,"2-5":40,"6+":20}'
)
ON CONFLICT(id) DO NOTHING;
EOF
fi


