package petrinet

import "fmt"

type Transition struct {
	Name string
	From []string
	To   []string
}

type NotEnabledTransitionError struct {
	BlockerList    *BlockerList
	TransitionName string
}

func (m *NotEnabledTransitionError) Error() string {
	return fmt.Sprintf("transition %s is not enabled for workflow", m.TransitionName)
}
