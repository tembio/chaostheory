package main

import (
	"fmt"
	"leaderboard/handlers"
	"leaderboard/internal"
	"leaderboard/repositories"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"common"
)

func main() {
	// Initialize SQLiteScoreRepository
	leaderboardsRepo, competitionsRepo, err := initialiseRepositories()
	if err != nil {
		fmt.Printf("Error initializing repositories: %v\n", err)
		return
	}
	defer leaderboardsRepo.Close()
	defer competitionsRepo.Close()

	defaultRuleEvaluator := &internal.BetRuleEvaluator{}
	leaderboard := internal.NewLeaderboard(defaultRuleEvaluator)

	// Load existing data from DB
	if err := loadLeaderBoardDataFromDB(leaderboard, leaderboardsRepo, competitionsRepo); err != nil {
		fmt.Printf("Error loading leaderboard data from DB: %v\n", err)
		return
	}

	///////// HTTP server setup /////////
	leaderboardsHandler := handlers.NewLeaderboardsHandler(leaderboardsRepo)
	competitionsHandler := handlers.NewCompetitionsHandler(competitionsRepo)
	websocketHandler := handlers.NewWebsocketHandler()

	r := mux.NewRouter()
	r.Handle("/leaderboards/{id}", http.HandlerFunc(leaderboardsHandler.GetLeaderboardByID)).Methods("GET")
	r.Handle("/competitions", authMiddleware(http.HandlerFunc(competitionsHandler.CreateCompetition))).Methods("POST")
	r.HandleFunc("/ws", http.HandlerFunc(websocketHandler.WebsocketHandler))

	go func() {
		// Start the HTTP server
		fmt.Println("Leaderboard API server listening on :8080")
		http.ListenAndServe(":8080", r)
	}()

	///////// RabbitMQ setup /////////
	eventHandler := handlers.NewBetEventHandler(leaderboardsRepo, leaderboard, websocketHandler)

	go func() {
		rabbitPort := os.Getenv("RABBITMQ_PORT")
		rabbitURL := fmt.Sprintf("amqp://guest:guest@rabbitleaderboard:%s/", rabbitPort)
		betQueue := "bet_events"

		var betReceiver *internal.RabbitMQReceiver
		for {
			betReceiver, err = internal.NewRabbitMQReceiver(rabbitURL, betQueue)
			if err == nil {
				break
			}
			fmt.Printf("RabbitMQ not ready, retrying in 200ms: %v\n", err)
			time.Sleep(200 * time.Millisecond)
		}
		defer betReceiver.Close()

		for {
			err := betReceiver.Receive(func(body []byte, acknowledgeEventFunc func()) error {
				if err := eventHandler.Handle(body); err != nil {
					return fmt.Errorf("error handling bet event: %v", err)
				}
				acknowledgeEventFunc()
				return nil
			})
			if err != nil {
				fmt.Printf("Error receiving bet event: %v\n", err)
			}
		}
	}()

	// Block forever so main does not exit while goroutines are running
	select {}
}

// authMiddleware is a simple middleware to check for a hardcoded Authorization header
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer secrettoken" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func initialiseRepositories() (*repositories.SQLiteLeaderboards, *repositories.SQLiteCompetitions, error) {
	dbPath := "db/leaderboard.db"
	leaderboardsRepo, err := repositories.NewSQLiteLeaderboardsRepository(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing SQLiteLeaderboardsRepository: %v", err)
	}

	competitionsRepo, err := repositories.NewSQLiteCompetitionsRepository(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing SQLiteCompetitionsRepository: %v", err)
	}
	return leaderboardsRepo, competitionsRepo, nil
}

func loadLeaderBoardDataFromDB(lb *internal.Leaderboard, leaderboardsRepo *repositories.SQLiteLeaderboards, competitionsRepo *repositories.SQLiteCompetitions) error {
	lbFromDB, err := leaderboardsRepo.GetAll()
	if err != nil {
		return fmt.Errorf("error retrieving leaderboards: %v", err)
	}
	lb.Load(lbFromDB)

	competitions, err := loadCompetitions(competitionsRepo)
	if err != nil {
		return fmt.Errorf("error retrieving competitions: %v", err)
	}
	for _, comp := range competitions {
		lb.RegisterCompetition(comp)
	}
	return nil
}

// loadCompetitions loads all competitions from the repository and returns them as a slice in memory
func loadCompetitions(repo repositories.CompetitionsRepository) ([]*common.Competition, error) {
	competitions, err := repo.GetAll()
	if err != nil {
		return nil, err
	}
	return competitions, nil
}
