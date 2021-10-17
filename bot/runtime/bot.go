package runtime

import (
	"bot-daedalus/config"
)

type Bot interface {
	HandleRequest(mf SerializedMessageFactory)
}

type DefaultBot struct {
	ScenarioPath    string
	ScenarioName    string
	TokenFactory    TokenFactory
	TokenRepository TokenRepository
}

func (b *DefaultBot) HandleRequest(mf SerializedMessageFactory) {
	cfg := config.GetScenarioConfig(b.ScenarioPath, b.ScenarioName)
	provider, _ := GetProvider(cfg.ProviderConfig, cfg.Name, b.TokenFactory, mf)

	sbuilder := ScenarioBuilder{
		ActionRegistry: nil,
		Repository:     nil,
		Provider:       provider,
		states:         nil,
	}
	s, err := sbuilder.BuildScenario(b.ScenarioPath, b.ScenarioName)

	if err != nil {
		panic(err)
	}

	token := s.HandleCommand()
	b.TokenRepository.Persist(token)
}
