package runtime

import (
	"bot-daedalus/config"
)

type Bot interface {
	HandleRequest(mf SerializedMessageFactory)
}

type DefaultBot struct {
	TokenFactory TokenFactory
}

func (b *DefaultBot) HandleRequest(mf SerializedMessageFactory) {
	cfg := config.GetScenarioConfig()
	provider, _ := GetProvider(cfg.ProviderConfig, b.TokenFactory, mf)
	s := Scenario{cfg.StateMachineConfig, provider}
	s.HandleCommand()
}
