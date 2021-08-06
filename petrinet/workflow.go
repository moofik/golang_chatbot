package petrinet

type Workflow interface {
	// GetMarking is used to get current workflow state
	GetMarking(subject interface{}) (Marking, error)
	// CanFire checks whether we can fire a transition
	CanFire(subject interface{}, transition string) bool
	// Fire fires a transition
	Fire(subject interface{}, transition string)
	// GetAllowedTransitions is used to get allowed transitions for current subject's state
	GetAllowedTransitions(subject interface{})
	// GetDefinition gets workflow definition
	GetDefinition()
}

type DefaultWorkflow struct {
	definition Definition
}
