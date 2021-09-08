package runtime

import (
	"bot-daedalus/config"
)

type Bot interface {
	HandleRequest(mf SerializedMessageFactory)
}

type DefaultBot struct {
	ScenarioPath string
	ScenarioName string
	TokenFactory TokenFactory
}

func (b *DefaultBot) HandleRequest(mf SerializedMessageFactory) {
	cfg := config.GetScenarioConfig(b.ScenarioPath, b.ScenarioName)
	provider, _ := GetProvider(cfg.ProviderConfig, b.TokenFactory, mf)
	s := Scenario{cfg.StateMachineConfig, provider}
	s.HandleCommand()
}
