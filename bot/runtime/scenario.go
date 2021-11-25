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

	if cmd.GetType() == TYPE_TEXT_INPUT {
		var commands []Command
		commands, err = currentState.GetCommandListByProto(cmd)
		fmt.Println("COMMANDS RETRIEVED:")
		for _, command := range commands {
			fmt.Println(command.Debug())
		}
		fmt.Println("END CR")

		for _, c := range commands {
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
	} else {
		actualTransition, err = currentState.GetTransition(cmd)
	}

	if actualTransition == nil {
		if lastOrderCommand != nil {
			actualTransition, err = currentState.GetTransition(lastOrderCommand)
		} else if token.GetState() != "unknown" {
			_ = s.Provider.SendMarkupMessage(
				[]string{"Маркет💵", "Кошелек💠"},
				"К сожалению я не знаю такой комманды. Вы можете воспользоваться меню ниже.",
				ProviderContext{
					State:   currentState,
					Command: nil,
					Token:   token,
				},
			)

			return token
		} else {
			return token
		}
	}

	if err != nil && token.GetState() != "unknown" {
		_ = s.Provider.SendMarkupMessage(
			[]string{"Маркет💵", "Кошелек💠"},
			"Произошел сбой и я не могу произвести запрошенное действие, но вы можете воспользоваться меню ниже.",
			ProviderContext{
				State:   currentState,
				Command: nil,
				Token:   token,
			},
		)

		return token
	} else if err != nil {
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
	ActionRegistry   func(string, map[string]interface{}) Action
	CommandRegistry  func(string, string, []interface{}) Command
	Repository       DelayedTransitionRepository
	Provider         ChatProvider
	states           map[string]*State
	currentPlaceName string
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
