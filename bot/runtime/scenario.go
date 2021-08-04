package runtime

import (
	"bot-daedalus/config"
	"fmt"
)

type Scenario struct {
	Config   config.StateMachineConfig
	Provider ChatProvider
}

func (s *Scenario) HandleCommand()  {
	// get state by token
	// run state actions

	states := GetStates(s.Config)
	cmd := s.Provider.GetCommand()
	token := s.Provider.GetToken()

	fmt.Println("Running actions")
	for _, state := range states {
		for _, action := range state.Actions {
			action.Run(s.Provider, token, state, cmd)
		}
	}
}
