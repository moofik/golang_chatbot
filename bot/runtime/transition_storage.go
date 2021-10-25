package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/petrinet"
	"sort"
)

type TransitionAndCommandElement struct {
	Transition *petrinet.Transition
	Command    command.Command
}

type TransitionStorage struct {
	index         int
	transitionMap map[string]*TransitionAndCommandElement
}

func (ts *TransitionStorage) FindTransitionByCommand(command command.Command) *petrinet.Transition {
	for _, element := range ts.transitionMap {
		if element.Command.ToHash() == command.ToHash() {
			return element.Transition
		}
	}

	return nil
}

func (ts *TransitionStorage) FindTransitionByProto(command command.Command) *petrinet.Transition {
	for _, element := range ts.transitionMap {
		if element.Command.ToProtoHash() == command.ToProtoHash() {
			return element.Transition
		}
	}

	return nil
}

func (ts *TransitionStorage) FindCommandByProto(command command.Command) command.Command {
	for _, element := range ts.transitionMap {
		if element.Command.ToProtoHash() == command.ToProtoHash() {
			return element.Command
		}
	}

	return nil
}

func (ts *TransitionStorage) Add(c command.Command, t *petrinet.Transition) {
	if ts.transitionMap == nil {
		ts.transitionMap = map[string]*TransitionAndCommandElement{}
	}

	ts.transitionMap[c.ToHash()] = &TransitionAndCommandElement{
		Transition: t,
		Command:    c,
	}
}

func (ts *TransitionStorage) Next() (*petrinet.Transition, bool) {
	keys := make([]string, 0, len(ts.transitionMap))

	for k := range ts.transitionMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	has := false
	if ts.index < len(keys) {
		has = true
	}

	if has {
		if t, ok := ts.transitionMap[keys[ts.index]]; ok {
			ts.index++
			return t.Transition, true
		}
	}

	return nil, false
}

func (ts *TransitionStorage) Current() *petrinet.Transition {
	keys := make([]string, 0, len(ts.transitionMap))

	for k := range ts.transitionMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if ts.index != 0 {
		if t, ok := ts.transitionMap[keys[ts.index]]; ok {
			return t.Transition
		}
	}

	return nil
}

func (ts *TransitionStorage) Empty() bool {
	return len(ts.transitionMap) == 0
}

func (ts *TransitionStorage) Count() int {
	return len(ts.transitionMap)
}

func (ts *TransitionStorage) AllButtonCommands() []command.Command {
	var commandsSlice []command.Command

	for _, element := range ts.transitionMap {
		if element.Command.GetMetadata().Cmd == "button" {
			commandsSlice = append(commandsSlice, element.Command)
		}
	}

	return commandsSlice
}
