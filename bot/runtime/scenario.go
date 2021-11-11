package runtime

import (
	"bot-daedalus/petrinet"
	"fmt"
	"reflect"

	"github.com/spf13/viper"
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

	if cmd.GetType() == TYPE_TEXT_INPUT {
		var transitions []*petrinet.Transition
		var commands []Command
		transitions, err = currentState.GetTransitionListByUniqueness(cmd)
		//finding valid transitions
		commands, err = currentState.GetCommandListByUniqueness(cmd)

		for _, c := range commands {
			for _, t := range transitions {
				if c.Pass(s.Provider, token, t) {
					actualTransition = t
					break
				}
			}

			if actualTransition != nil {
				break
			}
		}
	} else {
		actualTransition, err = currentState.GetTransition(cmd)
	}

	if actualTransition == nil {
		fmt.Println("actual transition not found")
		return token
	}

	if err != nil {
		// handle state error
		err := s.Provider.SendTextMessage(err.Error(), ProviderContext{
			State: currentState,
			Command: &UserInputCommand{
				Text: "",
				Metadata: &Metadata{
					Cmd:   "/system",
					Place: "noplace",
				},
			},
			Token: token,
		})

		if err != nil {
			return nil
		}

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
		newCmd, _ := newState.GetCommandByProto(&InstantTransitionCommand{Metadata: &Metadata{
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
	ActionRegistry   func(string, map[string]string) Action
	Repository       DelayedTransitionRepository
	Provider         ChatProvider
	states           map[string]*State
	currentPlaceName string
}

func (b *ScenarioBuilder) BuildScenario(path string, name string) (*Scenario, error) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	statesConfig := viper.Get("scenario.states")
	statesConfigValue := reflect.ValueOf(statesConfig)

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
		params := make(map[string]string)

		for _, kk := range innerMap.MapKeys() {
			if kk.Elem().String() == "name" {
				name = innerMap.MapIndex(kk).Elem().String()
			}

			if kk.Elem().String() == "params" {
				needParams = true
				reflectParams := innerMap.MapIndex(kk).Elem()
				for _, kk := range reflectParams.MapKeys() {
					params[kk.Elem().String()] = reflectParams.MapIndex(kk).Elem().String()
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
					newCmd := CreateCommand(cmdType, b.currentPlaceName, arguments)
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
