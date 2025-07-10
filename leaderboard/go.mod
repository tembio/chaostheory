module leaderboard

go 1.22.4

require (
	common v0.0.0
	github.com/expr-lang/expr v1.17.5
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/rabbitmq/amqp091-go v1.10.0
)

require github.com/gorilla/mux v1.8.1

replace common => ../common
