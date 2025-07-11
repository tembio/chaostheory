package handlers

import (
	"common"
	"encoding/json"
	"fmt"
	"leaderboard/internal"

	"leaderboard/repositories"
)

type BetEventHandler struct {
	leaderboardsRepo repositories.LeaderboardsRepository
	leaderboard      internal.LeaderboardInterface
	websocketHandler *WebsocketHandler
}

type UserEventHandler struct {
	leaderboardsRepo repositories.LeaderboardsRepository
}

func NewBetEventHandler(repo repositories.LeaderboardsRepository, leaderboard internal.LeaderboardInterface, websocketHandler *WebsocketHandler) *BetEventHandler {
	return &BetEventHandler{
		leaderboardsRepo: repo,
		leaderboard:      leaderboard,
		websocketHandler: websocketHandler,
	}
}

func (beh *BetEventHandler) Handle(body []byte) error {
	var betEvent common.BetEvent
	if err := json.Unmarshal(body, &betEvent); err == nil {
		fmt.Printf("Received bet event: %+v\n", betEvent)

		updatedData, err := beh.leaderboard.Update(betEvent)
		if err != nil {
			println("Error updating leaderboard:", err)
			return fmt.Errorf("error updating leaderboard: %v", err)
		}

		for _, update := range updatedData {
			if err := beh.leaderboardsRepo.Update(update.CompetitionID, update.UserID, update.Score); err != nil {
				println("Error storing score in SQLite:", err)
				return fmt.Errorf("error storing score in SQLite: %v", err)
			}
		}

		go sendCompetitionsUpdatesToWebsocket(beh.websocketHandler, beh.leaderboardsRepo, updatedData)

		return nil

	} else {
		return fmt.Errorf("error unmarshalling bet event: %v", err)
	}
}

func NewUserEventHandler(repo repositories.LeaderboardsRepository) *UserEventHandler {
	return &UserEventHandler{
		leaderboardsRepo: repo,
	}
}

func (ueh *UserEventHandler) UserEventHandler(body []byte) error {
	var userEvent common.UserEvent
	if err := json.Unmarshal(body, &userEvent); err == nil {
		// EXTRA FUNCTIONALITY: Store user in SQLite
		return nil

	} else {
		return fmt.Errorf("error unmarshalling user event: %v", err)
	}
}

func sendCompetitionsUpdatesToWebsocket(handler *WebsocketHandler, leaderboardsRepo repositories.LeaderboardsRepository, updates []*internal.UpdatedData) {
	if handler == nil {
		fmt.Println("WebSocket handler is not initialized")

		return // If no WebSocket handler, skip sending updates
	}

	if leaderboardsRepo == nil {
		fmt.Println("Leaderboards repository is not initialized")
		return
	}

	for _, updatedCompetition := range updates {
		competitionID := updatedCompetition.CompetitionID

		updates, err := leaderboardsRepo.GetTopN(competitionID, 10)
		if err != nil {
			fmt.Printf("Error retrieving top N users for competition %d: %v\n", competitionID, err)
			return
		}

		message := struct {
			CompetitionID uint
			Users         []*common.User
		}{
			CompetitionID: competitionID,
			Users:         updates,
		}

		if err := handler.SendMessage(message); err != nil {
			fmt.Printf("Error sending WebSocket message: %v\n", err)
		}
	}
}
