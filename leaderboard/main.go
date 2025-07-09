package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	comps := []Competition{
		{
			ID:        1,
			Name:      "Monthly Challenge",
			ScoreRule: "event_type=='bet' && distributor=='evo' ? amount : 0",
			StartTime: "2023-10-01T00:00:00Z",
			EndTime:   "2023-10-31T23:59:59Z",
			Rewards:   []string{"Gold Medal", "Silver Medal", "Bronze Medal"},
		},
		{
			ID:        2,
			Name:      "Weekly Sprint",
			ScoreRule: "event_type=='bet' && game=='Poker' ? amount : 0",
			StartTime: "2023-10-01T00:00:00Z",
			EndTime:   "2023-10-31T23:59:59Z",
			Rewards:   []string{"Gold Medal", "Silver Medal", "Bronze Medal"},
		},
	}

	fmt.Println("Leaderboard service started")

	ruleEvaluator := &BetRuleEvaluator{}
	leaderboard := NewLeaderboard(ruleEvaluator)
	// Register competitions in the leaderboard
	leaderboard.RegistrerCompetition(&comps[0])
	leaderboard.RegistrerCompetition(&comps[1])

	// Initialize SQLiteScoreRepository
	dbPath := "leaderboard.db"
	scoreRepository, err := NewSQLiteScoreRepository(dbPath)
	if err != nil {
		fmt.Printf("Error initializing SQLiteScoreRepository: %v\n", err)
		return
	}
	defer scoreRepository.db.Close()

	// Restore existing scores from SQLite
	allScores, err := scoreRepository.GetAllScores()
	if err != nil {
		fmt.Printf("Error retrieving scores: %v\n", err)
		return
	}
	leaderboard.LoadScores(allScores)

	// Initialize RabbitMQ receivers
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitURL := fmt.Sprintf("amqp://guest:guest@rabbitleaderboard:%s/", rabbitPort)
	betQueue := "bet_events"

	var betReceiver *RabbitMQReceiver
	for {
		betReceiver, err = NewRabbitMQReceiver(rabbitURL, betQueue)
		if err == nil {
			break
		}
		fmt.Printf("RabbitMQ not ready, retrying in 100ms: %v\n", err)
		time.Sleep(100 * time.Millisecond)
	}
	defer betReceiver.Close()

	// receiver := NewMockBetEventReceiver()

	go func() {
		for {
			err := betReceiver.Receive(func(body []byte, acknowledgeEvent func()) error {
				err := BetEventHandler(body, leaderboard, scoreRepository, acknowledgeEvent)
				if err != nil {
					fmt.Printf("Error handling bet event: %v\n", err)
				}
				return err
			})
			if err != nil {
				fmt.Printf("Error receiving bet event: %v\n", err)
			}
		}
	}()

	stop := make(chan bool)
	for {
		<-stop // Block until an event is received
	}

	// Register competitions in the leaderboard
	// leaderboard.RegistrerCompetition(&comps[0])
	// leaderboard.RegistrerCompetition(&comps[1])

	// fmt.Println("Competitions Results:")
	// for compID, users := range leaderboard.competitionsResults {
	// 	fmt.Printf("Competition ID: %d\n", compID)
	// 	for _, user := range users {
	// 		fmt.Printf("  User ID: %d, Score: %.2f\n", user.ID, user.Score)
	// 	}
	// }

}
