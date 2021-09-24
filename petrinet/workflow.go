package petrinet

import "fmt"

type Workflow interface {
	// GetMarking is used to get current workflow state
	GetMarking(subject interface{}) (*Marking, error)
	// CanFire checks whether we can fire a transition
	CanFire(subject interface{}, transition string) (bool, error)
	// Fire fires a transition
	Fire(subject interface{}, transition string) (*Marking, error)
	// GetTransitionBlockerList is used to get blockers for transition for current subject's state
	GetTransitionBlockerList(subject interface{}, marking *Marking, transition *Transition) BlockerList
	// GetDefinition gets workflow Definition
	GetDefinition() *Definition
}

type DefaultWorkflow struct {
	Definition     *Definition
	MarkingStorage *MarkingStorage
	Name           string
}

func (w *DefaultWorkflow) GetMarking(subject interface{}) (*Marking, error) {
	marking, err := w.MarkingStorage.GetMarking(subject)

	if err != nil {
		return nil, err
	}

	if len(marking.Places) == 0 {
		if len(w.Definition.InitialPlaces) == 0 {
			return nil, fmt.Errorf("the Marking is empty and there is no initial place for workflow %s", w.Name)
		}

		for _, place := range w.Definition.InitialPlaces {
			marking.Mark(place)
		}

		err := w.MarkingStorage.SetMarking(subject, marking)

		if err != nil {
			return nil, err
		}
	}

	for name := range marking.Places {
		if len(w.Definition.Places) == 0 {
			return nil, fmt.Errorf(
				"it seems you forgot to add places to the workflow %s",
				name,
			)
		}

		if _, ok := w.Definition.Places[name]; !ok {
			return nil, fmt.Errorf(
				"place %s is not valid for workflow %s",
				name,
				w.Name,
			)
		}
	}

	return marking, nil
}

func (w *DefaultWorkflow) CanFire(subject interface{}, transition string) (bool, error) {
	marking, err := w.GetMarking(subject)

	if err != nil {
		return false, err
	}

	for _, t := range w.Definition.Transitions {
		if transition != t.Name {
			continue
		}

		blockerList := w.GetTransitionBlockerList(subject, marking, t)

		if blockerList.empty() {
			return true, nil
		}
	}

	return false, nil
}

func (w *DefaultWorkflow) GetTransitionBlockerList(subject interface{}, marking *Marking, transition *Transition) BlockerList {
	for _, place := range transition.From {
		if !marking.Has(place) {
			return BlockerList{blockers: []*Blocker{createNotEnabledBlocker()}}
		}
	}

	return BlockerList{}
}

func (w *DefaultWorkflow) Fire(subject interface{}, transition string) (*Marking, error) {
	marking, err := w.GetMarking(subject)

	if err != nil {
		return nil, err
	}

	transitionExist := true
	var approvedTransitions []*Transition
	var blockerList *BlockerList

	for _, t := range w.Definition.Transitions {
		if t.Name != transition {
			continue
		}

		transitionExist = true
		blockers := w.GetTransitionBlockerList(subject, marking, t)

		if blockers.empty() {
			approvedTransitions = append(approvedTransitions, t)
			continue
		}

		if !blockers.has(CODE_NOT_ENABLED) {
			*blockerList = blockers
		}
	}

	if !transitionExist {
		return nil, &NotEnabledTransitionError{blockerList, transition}
	}

	for _, t := range approvedTransitions {
		for _, place := range t.From {
			err := marking.Unmark(place)
			if err != nil {
				return nil, err
			}
		}

		for _, place := range t.To {
			marking.Mark(place)
		}

		err := w.MarkingStorage.SetMarking(subject, marking)
		if err != nil {
			return nil, err
		}
	}

	return marking, nil
}

func (w *DefaultWorkflow) GetDefinition() *Definition {
	return w.Definition
}
