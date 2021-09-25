package petrinet

import "fmt"

type Transition struct {
	Name string
	From []string
	To   []string
}

type TransitionError interface {
	error
	GetBlockerList() *BlockerList
	GetTransitionName() string
}

type GenericTransitionError struct {
	BlockerList    *BlockerList
	TransitionName string
	innerError     error
}

func (m *GenericTransitionError) Error() string {
	return m.innerError.Error()
}

func (m *GenericTransitionError) GetBlockerList() *BlockerList {
	return m.BlockerList
}

func (m *GenericTransitionError) GetTransitionName() string {
	return m.TransitionName
}

type NotDefinedTransitionError struct {
	BlockerList    *BlockerList
	TransitionName string
	WorkflowName   string
}

func (m *NotDefinedTransitionError) Error() string {
	return fmt.Sprintf("transition %s is not defined for workflow %s", m.TransitionName, m.WorkflowName)
}

func (m *NotDefinedTransitionError) GetBlockerList() *BlockerList {
	return m.BlockerList
}

func (m *NotDefinedTransitionError) GetTransitionName() string {
	return m.TransitionName
}

type NotEnabledTransitionError struct {
	BlockerList    *BlockerList
	TransitionName string
	WorkflowName   string
}

func (m *NotEnabledTransitionError) Error() string {
	return fmt.Sprintf("transition %s is not enabled for workflow %s", m.TransitionName, m.WorkflowName)
}

func (m *NotEnabledTransitionError) GetBlockerList() *BlockerList {
	return m.BlockerList
}

func (m *NotEnabledTransitionError) GetTransitionName() string {
	return m.TransitionName
}
