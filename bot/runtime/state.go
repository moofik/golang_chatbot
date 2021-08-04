package runtime

import (
	"bot-daedalus/config"
)

type State struct {
	Name        string
	Actions     []Action
	Transitions map[string]string
}

func GetStates(c config.StateMachineConfig) []*State {
	actions := []Action{
		&SendTextMessage{"Привет {{.FirstName}} {{.LastName}}. Я нахожусь в режиме простоя. Сейчас я не могу тебе ничем помочь."},
	}

	return []*State{
		&State{
			"mock",
			actions,
			map[string]string{},
		},
	}
}
