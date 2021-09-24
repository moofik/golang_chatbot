package app

import (
	"bot-daedalus/config"
	"bot-daedalus/petrinet"
)

//skeleton
func BuildWorkflow(config.StateMachineConfig) petrinet.Workflow {
	return &petrinet.DefaultWorkflow{
		Definition: &petrinet.Definition{nil, nil, nil},
		Name:       "test",
	}
}
