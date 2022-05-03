package runtime

import (
	"bot-daedalus/petrinet"
	"fmt"
	"reflect"
)

type Scenario struct {
	petrinet.Workflow
	Provider ChatProvider
	States   map[string]*State
	DelayedTransitionRepository
	SettingsRepository SettingsRepository
	StateErrorHandler  func(p ChatProvider, ctx ProviderContext)
}

func (s *Scenario) GetCurrentState(token TokenProxy) *State {
	marking, err := s.Workflow.GetMarking(token)

	if err != nil {
		panic(err)
	}
	for place, state := range s.States {
		for markingPlace := range marking.Places {
			if markingPlace == place {
				return state
			}
		}
	}

	return nil
}

func SeparateRecongnizeInputs(cmds []Command) ([]Command, []Command) {
	riCmds := []Command{}
	otherCmds := []Command{}
	riCmd := RecognizeInputCommand{}

	for _, cmd := range cmds {
		if cmd.GetType() == riCmd.GetType() {
			riCmds = append(riCmds, cmd)
		} else {
			otherCmds = append(otherCmds, cmd)
		}
	}

	return riCmds, otherCmds
}

func (s *Scenario) HandleCommand(cmd Command, currentState *State, token TokenProxy) TokenProxy {
	// get state by token
	// run state actions
	if cmd == nil {
		_ = fmt.Sprintf("command not found for token %d and scenario %s", token.GetChatId(), s.Provider.GetScenarioName())
		return token
	}

	var actualTransition *petrinet.Transition
	var err error
	var lastOrderCommand Command
	setting := s.SettingsRepository.FindByScenarioName(s.Provider.GetScenarioName())

	if setting != nil && setting.IsOffline() == true {
		fmt.Println("BOT IS OFFLINE")
		return token
	}

	if cmd.GetType() == TYPE_TEXT_INPUT || cmd.GetType() == TYPE_PENDING {
		var commands []Command
		commands, err = currentState.GetCommandListByProto(cmd)
		riCmds, otherCmds := SeparateRecongnizeInputs(commands)

		for _, c := range riCmds {
			if c.GetType() == TYPE_TEXT_INPUT {
				lastOrderCommand = c
				continue
			}

			if ok, _ := c.Pass(s.Provider, cmd, token); ok {
				actualTransition, _ = currentState.GetTransitionByUniqueness(c)
			}

			if actualTransition != nil {
				break
			}
		}

		if actualTransition == nil {
			for _, c := range otherCmds {
				if c.GetType() == TYPE_TEXT_INPUT {
					lastOrderCommand = c
					continue
				}

				if ok, _ := c.Pass(s.Provider, cmd, token); ok {
					actualTransition, _ = currentState.GetTransitionByUniqueness(c)
				}

				if actualTransition != nil {
					break
				}
			}
		}
	} else {
		if actualTransition == nil {
			//fmt.Println("DBG pre 14 cmd: %v\n", cmd)
			//fmt.Println("DBG 14 cmd debg: %s\n", cmd.Debug())
			//
			//cmds, _ := currentState.GetCommandList()
			//fmt.Println("----------------")
			//for _, i2 := range cmds {
			//	fmt.Println(i2.Debug())
			//}
			//fmt.Println("----------------")
			actualTransition, err = currentState.GetTransition(cmd)
		}

		if actualTransition == nil {
			var commands []Command

			riPlaceholder := &RecognizeInputCommand{
				Marker:   "",
				Metadata: &CommandMetadata{Cmd: "recognize_input"},
			}

			commands, err = currentState.GetCommandListByProto(riPlaceholder)

			for _, c := range commands {
				if ok, _ := c.Pass(s.Provider, cmd, token); ok {
					actualTransition, _ = currentState.GetTransitionByUniqueness(c)
					if actualTransition == nil {
						fmt.Println("BUT TRANS IS NULL")
					}
				}

				if actualTransition != nil {
					break
				}
			}
		}
	}

	if actualTransition == nil {
		if lastOrderCommand != nil {
			actualTransition, err = currentState.GetTransition(lastOrderCommand)
		} else if s.StateErrorHandler != nil {
			fmt.Println("DBG 17")
			//s.StateErrorHandler(s.Provider, ProviderContext{
			//	State:   currentState,
			//	Command: nil,
			//	Token:   token,
			//})
			return token
		} else {
			return token
		}
	}

	if err != nil {
		fmt.Println(err.Error())
		return token
	}

	can, err := s.Workflow.CanFire(token, actualTransition.Name)
	var newState *State

	if can {
		newState = s.States[actualTransition.To[0]]
		err := newState.Execute(token, s.Provider, cmd, currentState)
		if err != nil {
			panic(err.Error())
			//handle action error
		}
		_, err = s.Workflow.Fire(token, actualTransition.Name)
		if err != nil {
			panic(err.Error())
			//handle state error
		}
	}

	if err != nil {
		panic(err.Error())
	}

	if !can {
		fmt.Println("Can not move further! Prohibited transition.")
	}

	if newState != nil {
		newCmd, _ := newState.GetCommandByProto(&InstantTransitionCommand{Metadata: &CommandMetadata{
			Cmd:        "instant",
			Place:      "",
			Uniqueness: "",
		}})

		if newCmd != nil {
			token = s.HandleCommand(newCmd, newState, token)
		}
	}

	return token
}

type ScenarioBuilder struct {
	ActionRegistry     func(string, map[string]interface{}) Action
	CommandRegistry    func(string, string, []interface{}) Command
	StateErrorHandler  func(p ChatProvider, ctx ProviderContext)
	Repository         DelayedTransitionRepository
	SettingsRepository SettingsRepository
	Provider           ChatProvider
	states             map[string]*State
	currentPlaceName   string
}

func (b *ScenarioBuilder) BuildScenario(config ScenarioConfig) (*Scenario, error) {
	statesConfigValue := reflect.ValueOf(config.States)
	b.walk(statesConfigValue)

	markingStorage := &petrinet.MarkingStorage{
		SingleState:  true,
		MarkingField: "State",
	}

	definition := &petrinet.Definition{
		Transitions:   nil,
		InitialPlaces: nil,
		Places:        nil,
	}

	for _, state := range b.states {
		definition.AddPlace(state.Name)
	}

	for _, state := range b.states {
		for i := 0; i < state.TransitionStorage.Count(); i++ {
			transition, _ := state.TransitionStorage.Next()
			err := definition.AddTransition(transition)
			if err != nil {
				return nil, err
			}
		}
	}

	workflow := &petrinet.DefaultWorkflow{
		Definition:     definition,
		MarkingStorage: markingStorage,
		Name:           "BotWorkflow",
	}

	return &Scenario{
		workflow,
		b.Provider,
		b.states,
		b.Repository,
		b.SettingsRepository,
		b.StateErrorHandler,
	}, nil
}

func (b *ScenarioBuilder) walk(v reflect.Value) {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	var ts *TransitionStorage
	var actions []Action
	b.currentPlaceName = ""

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			b.walk(v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Elem().String() == "name" {
				b.currentPlaceName = v.MapIndex(k).Elem().String()
			}
		}

		if b.currentPlaceName == "" {
			panic("state without name detected")
		}

		for _, k := range v.MapKeys() {
			if k.Elem().String() == "actions" {
				actions = b.walkActions(v.MapIndex(k))
			}

			if k.Elem().String() == "transitions" {
				ts = b.walkTransitions(v.MapIndex(k))
			}
		}
	default:
		// handle other types
	}

	if b.currentPlaceName != "" && ts != nil {
		if b.states == nil {
			b.states = map[string]*State{}
		}

		b.states[b.currentPlaceName] = &State{
			Name:              b.currentPlaceName,
			Actions:           actions,
			TransitionStorage: ts,
		}
	}
}

func (b *ScenarioBuilder) walkActions(v reflect.Value) []Action {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	var actions []Action

	for _, k := range v.MapKeys() {
		name := ""
		innerMap := v.MapIndex(k).Elem()
		needParams := false
		params := make(map[string]interface{})

		for _, kk := range innerMap.MapKeys() {
			if kk.Elem().String() == "name" {
				name = innerMap.MapIndex(kk).Elem().String()
			}

			if kk.Elem().String() == "params" {
				needParams = true
				reflectParams := innerMap.MapIndex(kk).Elem()
				for _, kk := range reflectParams.MapKeys() {
					params[kk.Elem().String()] = reflectParams.MapIndex(kk).Elem().Interface()
				}
			}
		}

		if (len(params) > 0 || !needParams) && name != "" {
			actions = append(actions, CreateAction(name, params, b.ActionRegistry))
		}
	}

	return actions
}

func (b *ScenarioBuilder) walkTransitions(v reflect.Value) *TransitionStorage {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	ts := &TransitionStorage{
		index:         0,
		transitionMap: nil,
	}

	// iterate over transitions
	for i := 0; i < v.Len(); i++ {
		tr := v.Index(i).Elem()
		stateTo := ""
		name := ""
		var commands []Command

		for _, kk := range tr.MapKeys() {
			if kk.Elem().String() == "state_to" {
				stateTo = tr.MapIndex(kk).Elem().String()
			}

			if kk.Elem().String() == "name" {
				name = tr.MapIndex(kk).Elem().String()
			}
		}

		for _, kk := range tr.MapKeys() {
			if kk.Elem().String() == "command" {
				cmdMap := tr.MapIndex(kk).Elem()
				cmdType := ""
				var arguments []interface{} = nil

				for _, val := range cmdMap.MapKeys() {
					if val.Elem().String() == "type" {
						cmdType = cmdMap.MapIndex(val).Elem().String()
					}

					if val.Elem().String() == "arguments" {
						arguments = cmdMap.MapIndex(val).Elem().Interface().([]interface{})
					}
				}

				if cmdType != "" {
					newCmd := CreateCommand(cmdType, b.currentPlaceName, arguments, b.CommandRegistry)
					commands = append(commands, newCmd)
				}
			}
		}

		if name != "" && stateTo != "" {
			t := &petrinet.Transition{
				Name: name,
				From: []string{b.currentPlaceName},
				To:   []string{stateTo},
			}

			for _, c := range commands {
				ts.Add(c, t)
			}
		}
	}

	return ts
}
