package petrinet

import "fmt"

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
	definition     Definition
	markingStorage MarkingStorage
	name           string
}

func (w *DefaultWorkflow) GetMarking(subject interface{}) (*Marking, error) {
	m, err := w.markingStorage.GetMarking(subject)

	if err != nil {
		return nil, err
	}

	if len(m.Places) == 0 {
		if len(w.definition.InitialPlaces) == 0 {
			return nil, fmt.Errorf("the Marking is empty and there is no initial place for workflow %s", w.name)
		}

		for _, place := range w.definition.InitialPlaces {
			m.Mark(place)
		}

		err := w.markingStorage.SetMarking(subject, m)

		if err != nil {
			return nil, err
		}
	}

	for name, _ := range m.Places {
		if len(w.definition.Places) == 0 {
			return nil, fmt.Errorf(
				"it seems you forgot to add places to the workflow %s",
				name,
			)
		}

		if _, ok := w.definition.Places[name]; !ok {
			return nil, fmt.Errorf(
				"place %s is not valid for workflow %s",
				name,
				w.name,
			)
		}
	}

	return m, nil
}
