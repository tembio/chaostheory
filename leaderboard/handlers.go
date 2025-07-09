package main

import (
	"common"
	"encoding/json"
	"fmt"
)

func BetEventHandler(body []byte, leaderboard *Leaderboard, scoreRepository *SQLiteScoreRepository, acknowledgeEvent func()) error {
	fmt.Printf("Handling bet event: %s\n", body)

	var betEvent common.BetEvent
	if err := json.Unmarshal(body, &betEvent); err == nil {
		UpdatedData, err := leaderboard.Update(betEvent)
		if err != nil {
			return fmt.Errorf("error updating leaderboard: %v", err)
		}

		for _, update := range UpdatedData {
			err := scoreRepository.Update(update.CompetitionID, update.UserID, update.Score)
			if err != nil {
				return fmt.Errorf("error storing score in SQLite: %v", err)
			}
		}
		acknowledgeEvent() // Acknowledge the message after processing
		return nil

	} else {
		return fmt.Errorf("error unmarshalling bet event: %v", err)
	}
}

func UserEventHandler(body []byte, scoreRepository *SQLiteScoreRepository, acknowledgeEvent func()) error {
	fmt.Printf("Handling user event: %s\n", body)

	var userEvent common.UserEvent
	if err := json.Unmarshal(body, &userEvent); err == nil {
		// Store user
		// err := userRepository.Save(...)
		// if err != nil {
		// 	return fmt.Errorf("error storing user in SQLite: %v\n", err)
		// }
		acknowledgeEvent() // Acknowledge the message after processing
		return nil

	} else {
		return fmt.Errorf("error unmarshalling user event: %v", err)
	}
}
