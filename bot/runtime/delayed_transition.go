package runtime

import "time"

type DelayedTransition struct {
	TokenId   string
	ExecuteAt time.Time
	Scenario  string
}

func (t *DelayedTransition) getScenarioName() string {
	return t.Scenario
}

func (t *DelayedTransition) getExecutionTime() time.Time {
	return t.ExecuteAt
}

func (t *DelayedTransition) getTokenId() string {
	return t.TokenId
}

type DelayedTransitionRepository interface {
	save(transition *DelayedTransition)
	getByDay(time time.Time) []*DelayedTransition
	getAll(time time.Time) []*DelayedTransition
	delete(t *DelayedTransition)
}
