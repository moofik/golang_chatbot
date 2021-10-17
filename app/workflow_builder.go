package app

import (
	"bot-daedalus/bot/runtime"
	"bot-daedalus/petrinet"
)

func BuildWorkflow(states []*runtime.State) petrinet.Workflow {
	d, _ := buildDefinition(states)

	return &petrinet.DefaultWorkflow{
		Definition: d,
		Name:       "test",
	}
}

func buildDefinition(states []*runtime.State) (*petrinet.Definition, error) {
	d, _ := petrinet.CreateDefinition(nil, nil, nil)

	for _, state := range states {
		d.AddPlace(state.Name)

		for _, transition := range state.Transitions {
			err := d.AddTransition(transition)
			if err != nil {
				return nil, err
			}
		}
	}

	return d, nil
}
