package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/petrinet"
	"fmt"
)

type State struct {
	Name              string
	Actions           []Action
	TransitionStorage *TransitionStorage
}

func (s *State) GetTransition(command command.Command) (*petrinet.Transition, StateError) {
	transition := s.TransitionStorage.FindTransitionByCommand(command)
	if transition == nil {
		return nil, fmt.Errorf("transition not found for state %s and command %s", s.Name, command.Debug())
	}
	return transition, nil
}

func (s *State) Execute(token TokenProxy, provider ChatProvider, command command.Command) ActionError {
	for _, action := range s.Actions {
		err := action.Run(provider, token, s, command)

		if err != nil {
			return err
		}
	}

	return nil
}

type StateError interface {
	error
}

type GenericStateError struct {
	innerError error
}

func (m *GenericStateError) Error() string {
	return m.innerError.Error()
}
