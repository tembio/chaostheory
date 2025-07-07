module leaderboard

go 1.22.4

require (
	common v0.0.0
	github.com/expr-lang/expr v1.17.5
)

require github.com/mattn/go-sqlite3 v1.14.28 // indirect

replace common => ../common
