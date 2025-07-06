package main

import (
	"fmt"

	"common"
)

func main() {
	comps := []Competition{
		{
			ID:        1,
			Name:      "Monthly Challenge",
			ScoreRule: "event_type=='win' && distributor=='evo' ? amount : 0",
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

	event := common.BetEvent{
		EventID:      1,
		EventType:    "win",
		UserID:       42,
		Amount:       100.0,
		Currency:     "USD",
		ExchangeRate: 1.0,
		Game:         "Poker",
		Distributor:  "evo",
		Studio:       "StudioX",
		Timestamp:    "2023-10-01T12:00:00Z",
	}

	ruleEvaluator := &BetRuleEvaluator{}
	leaderboard := NewLeaderboard(ruleEvaluator)

	// Register competitions in the leaderboard
	leaderboard.RegistrerCompetition(&comps[0])
	leaderboard.RegistrerCompetition(&comps[1])

	UpdatedData, err := leaderboard.Update(event)
	if err != nil {
		fmt.Printf("Error updating leaderboard: %v\n", err)
	}

	for _, update := range UpdatedData {
		fmt.Printf("Competition ID: %d, User ID: %d, Score: %.2f\n",
			update.CompetitionID, update.UserID, update.Score)
	}

	fmt.Println("Competitions Results:")
	for compID, users := range leaderboard.competitionsResults {
		fmt.Printf("Competition ID: %d\n", compID)
		for _, user := range users {
			fmt.Printf("  User ID: %d, Score: %.2f\n", user.ID, user.Score)
		}
	}

}
