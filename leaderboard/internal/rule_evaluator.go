package internal

import (
	"github.com/expr-lang/expr"

	"common"
)

type BetRuleEvaluator struct {
	rules []string
}

type Match struct {
	Rule   string
	Result any
}

// AddRule appends a rule string to the RuleEvaluator
func (re *BetRuleEvaluator) AddRule(rule string) {
	re.rules = append(re.rules, rule)
}

// EvaluateRules evaluates an event against a list of rules and returns matches
func (evaluator *BetRuleEvaluator) EvaluateRules(event common.BetEvent) ([]Match, error) {
	betEventEnv := map[string]any{
		"event_id":      event.EventID,
		"event_type":    event.EventType.String(),
		"user_id":       event.UserID,
		"amount":        event.Amount,
		"currency":      event.Currency,
		"exchange_rate": event.ExchangeRate,
		"game":          event.Game,
		"distributor":   event.Distributor,
		"studio":        event.Studio,
		"timestamp":     event.Timestamp,
	}

	// TODO: cache the compiled programs for each rule to avoid recompilation

	var matches []Match
	for _, rule := range evaluator.rules {
		program, err := expr.Compile(rule, expr.Env(betEventEnv))
		if err != nil {
			return nil, err
		}
		output, err := expr.Run(program, betEventEnv)
		if err != nil {
			return nil, err
		}

		matches = append(matches, Match{
			Rule:   rule,
			Result: output,
		})
	}
	return matches, nil
}
