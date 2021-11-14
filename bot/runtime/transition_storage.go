package runtime

import (
	"bot-daedalus/petrinet"
	"sort"
)

type TransitionAndCommandElement struct {
	Transition *petrinet.Transition
	Command    Command
}

type TransitionStorage struct {
	index         int
	transitionMap map[string]*TransitionAndCommandElement
}

func (ts *TransitionStorage) FindTransitionByCommand(command Command) *petrinet.Transition {
	for _, element := range ts.transitionMap {
		if element.Command.ToHash() == command.ToHash() {
			return element.Transition
		}
	}

	return nil
}

func (ts *TransitionStorage) FindTransitionByUniqueness(command Command) *petrinet.Transition {
	for _, element := range ts.transitionMap {
		if element.Command.ToUniquenessHash() == command.ToUniquenessHash() {
			return element.Transition
		}
	}

	return nil
}

func (ts *TransitionStorage) FindTransitionListByUniqueness(command Command) []*petrinet.Transition {
	var tss []*petrinet.Transition

	for _, element := range ts.transitionMap {
		if element.Command.ToUniquenessHash() == command.ToUniquenessHash() {
			tss = append(tss, element.Transition)
		}
	}

	return tss
}

func (ts *TransitionStorage) FindTransitionByProto(command Command) *petrinet.Transition {
	for _, element := range ts.transitionMap {
		if element.Command.ToProtoHash() == command.ToProtoHash() {
			return element.Transition
		}
	}

	return nil
}

func (ts *TransitionStorage) FindCommandByProto(command Command) Command {
	for _, element := range ts.transitionMap {
		if element.Command.ToProtoHash() == command.ToProtoHash() {
			return element.Command
		}
	}

	return nil
}

func (ts *TransitionStorage) FindCommandByUniqueness(command Command) Command {
	for _, element := range ts.transitionMap {
		if element.Command.ToUniquenessHash() == command.ToUniquenessHash() {
			return element.Command
		}
	}

	return nil
}

func (ts *TransitionStorage) FindCommand(command Command) Command {
	for _, element := range ts.transitionMap {
		if element.Command.ToHash() == command.ToHash() {
			return element.Command
		}
	}

	return nil
}

func (ts *TransitionStorage) FindCommandListByProto(command Command) []Command {
	var elements []Command

	for _, element := range ts.transitionMap {
		if element.Command.ToProtoHash() == command.ToProtoHash() {
			elements = append(elements, element.Command)
		}
	}

	return elements
}

func (ts *TransitionStorage) Add(c Command, t *petrinet.Transition) {
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

func (ts *TransitionStorage) AllButtonCommands() []Command {
	commandsMap := map[string]Command{}
	var commandsSlice []Command

	for _, element := range ts.transitionMap {
		if element.Command.GetMetadata().Cmd == "button" {
			commandsMap[element.Command.GetMetadata().Uniqueness] = element.Command
		}
	}

	keys := make([]string, 0)
	for k, _ := range commandsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		commandsSlice = append(commandsSlice, commandsMap[k])
	}

	return commandsSlice
}

func (ts *TransitionStorage) FindCommandListByUniqueness(c Command) []Command {
	var commandsSlice []Command

	for _, element := range ts.transitionMap {
		if element.Command.ToUniquenessHash() == c.ToUniquenessHash() {
			commandsSlice = append(commandsSlice, element.Command)
		}
	}

	return commandsSlice
}
