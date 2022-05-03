package runtime

import (
	"bot-daedalus/petrinet"
	"fmt"
	"sort"
)

type State struct {
	Name              string
	Actions           []Action
	TransitionStorage *TransitionStorage
}

func (s *State) GetTransition(command Command) (*petrinet.Transition, StateError) {
	transition := s.TransitionStorage.FindTransitionByCommand(command)
	if transition == nil {
		return nil, fmt.Errorf("transition not found for state %s and command %s", s.Name, command.Debug())
	}
	return transition, nil
}

func (s *State) GetTransitionByUniqueness(command Command) (*petrinet.Transition, StateError) {
	transition := s.TransitionStorage.FindTransitionByUniqueness(command)
	if transition == nil {
		return nil, fmt.Errorf("transition not found for state %s and command %s", s.Name, command.Debug())
	}
	return transition, nil
}

func (s *State) GetTransitionListByUniqueness(command Command) ([]*petrinet.Transition, StateError) {
	transitions := s.TransitionStorage.FindTransitionListByUniqueness(command)
	if len(transitions) == 0 {
		return nil, fmt.Errorf("transition not found for state %s and command %s", s.Name, command.Debug())
	}
	return transitions, nil
}

func (s *State) GetCommandByProto(command Command) (Command, StateError) {
	return s.TransitionStorage.FindCommandByProto(command), nil
}

func (s *State) GetCommandByUniqueness(command Command) (Command, StateError) {
	return s.TransitionStorage.FindCommandByUniqueness(command), nil
}

func (s *State) GetCommandListByUniqueness(command Command) ([]Command, StateError) {
	return s.TransitionStorage.FindCommandListByUniqueness(command), nil
}

func (s *State) GetCommand(command Command) (Command, StateError) {
	return s.TransitionStorage.FindCommand(command), nil
}

func (s *State) GetCommandListByProto(command Command) ([]Command, StateError) {
	return s.TransitionStorage.FindCommandListByProto(command), nil
}

func (s *State) GetCommandList() ([]Command, StateError) {
	return s.TransitionStorage.FindCommandList(), nil
}

func (s *State) Execute(token TokenProxy, provider ChatProvider, command Command, prevState *State) ActionError {
	actions := map[string]Action{}

	for _, action := range s.Actions {
		actions[action.GetName()] = action
	}

	keys := make([]string, 0)
	for k, _ := range actions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		err := actions[k].Run(provider, token, s, prevState, command)

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
