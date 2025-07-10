package internal

import (
	"common"
	"fmt"
)

type rulesToCompetitionID map[string]uint        // map[rule]competition
type usersIDToUser map[uint]*common.User         // map[userID]User
type scoresPerCompetition map[uint]usersIDToUser // map[competitionID]map[userID]User

type UpdatedData struct {
	CompetitionID uint
	UserID        uint
	Score         float64
}

type LeaderboardInterface interface {
	// Update processes a bet event and returns updated scores for users in competitions
	Update(event common.BetEvent) ([]*UpdatedData, error)
	// RegisterCompetition registers a competition with its score rule
	RegisterCompetition(comp *common.Competition)
}

// RuleEvaluator abstracts rule evaluation for Leaderboard
type RuleEvaluator interface {
	AddRule(rule string)
	EvaluateRules(event common.BetEvent) ([]Match, error)
}

type Leaderboard struct {
	ruleEvaluator       RuleEvaluator
	rulesToCompetition  rulesToCompetitionID
	competitionsResults scoresPerCompetition
}

// NewLeaderboard creates and returns a new Leaderboard instance
func NewLeaderboard(evaluator RuleEvaluator) *Leaderboard {
	return &Leaderboard{
		ruleEvaluator:       evaluator,
		rulesToCompetition:  rulesToCompetitionID{},
		competitionsResults: scoresPerCompetition{},
	}
}

// RegisterCompetition adds a new competition to the leaderboard
// If the competition's score rule is empty or already registered, it skips registration.
func (lb *Leaderboard) RegisterCompetition(comp *common.Competition) {
	if comp == nil || comp.ScoreRule == "" {
		// TODO: use logs
		fmt.Printf("Skipping registration of competition due to empty ScoreRule\n")
		return
	}
	if _, exists := lb.rulesToCompetition[comp.ScoreRule]; exists {
		// TODO: use logs
		fmt.Printf("Competition with ScoreRule '%s' already registered for competition ID %d\n", comp.ScoreRule, comp.ID)
		return
	}
	lb.ruleEvaluator.AddRule(comp.ScoreRule)
	lb.rulesToCompetition[comp.ScoreRule] = comp.ID
}

// Update updates the leaderboard with the results of a bet event
func (lb *Leaderboard) Update(event common.BetEvent) ([]*UpdatedData, error) {
	var updates []*UpdatedData

	if event.EventType == common.EventTypeLoss {
		return nil, nil // Skip loss events, only process bets and wins
	}

	matches, err := lb.ruleEvaluator.EvaluateRules(event)
	if err != nil {
		// TODO: use logs
		return nil, fmt.Errorf("error evaluating rules: %w", err)
	}

	for _, match := range matches {
		amount, err := toFloat64(match.Result)
		if err != nil {
			// TODO: use logs
			fmt.Printf("Event %d: Error converting output to float64: %v\n", event.EventID, err)
			continue // Skip this match if conversion fails
		}

		const epsilon = 1e-9
		if amount < epsilon && amount > -epsilon {
			continue // Skip rules that evaluate to 0
		}

		competitionID := lb.rulesToCompetition[match.Rule]
		// Initialize the map for the competition if it doesn't exist
		if _, exists := lb.competitionsResults[competitionID]; !exists {
			lb.competitionsResults[competitionID] = usersIDToUser{}
		}

		amount = toUSD(amount, event.ExchangeRate)

		// Update the user's score in the competition
		if user, exists := lb.competitionsResults[competitionID][event.UserID]; exists {
			user.Score += amount
			lb.competitionsResults[competitionID][event.UserID] = user
		} else {
			lb.competitionsResults[competitionID][event.UserID] = &common.User{
				ID:    event.UserID,
				Score: amount,
			}
		}

		updatedUser := lb.competitionsResults[competitionID][event.UserID]
		updates = append(updates, &UpdatedData{
			CompetitionID: competitionID,
			UserID:        updatedUser.ID,
			Score:         updatedUser.Score,
		})
	}

	return updates, nil
}

// Load populates the Leaderboard data with the provided leaderboards
func (lb *Leaderboard) Load(leaderboards map[uint][]common.User) {
	if lb.competitionsResults == nil {
		lb.competitionsResults = scoresPerCompetition{}
	}
	for compID, users := range leaderboards {
		if _, ok := lb.competitionsResults[compID]; !ok {
			lb.competitionsResults[compID] = usersIDToUser{}
		}
		for _, user := range users {
			lb.competitionsResults[compID][user.ID] = &user
		}
	}
}

// toFloat64 safely converts an interface{} to float64, handling int, int64, and float64
func toFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}

// toUSD converts an amount in any currency to USD using the exchange rate
func toUSD(amount, exchangeRate float64) float64 {
	return amount * exchangeRate
}
