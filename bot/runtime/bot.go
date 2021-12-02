package runtime

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Bot interface {
	HandleRequest(mf SerializedMessageFactory)
}

type DefaultBot struct {
	ScenarioPath      string
	ScenarioName      string
	TokenFactory      TokenFactory
	TokenRepository   TokenRepository
	ActionRegistry    func(string, map[string]interface{}) Action
	CommandRegistry   func(string, string, []interface{}) Command
	StateErrorHandler func(p ChatProvider, ctx ProviderContext)
}

func (b *DefaultBot) GetBaseActors(mf SerializedMessageFactory) (ScenarioConfig, ChatProvider, *Scenario) {
	cfg := GetScenarioConfig(b.ScenarioPath, b.ScenarioName)
	provider, _ := GetProvider(cfg.ProviderConfig, cfg.Name, b.TokenFactory, mf, b.TokenRepository)

	sbuilder := ScenarioBuilder{
		ActionRegistry:    b.ActionRegistry,
		CommandRegistry:   b.CommandRegistry,
		StateErrorHandler: b.StateErrorHandler,
		Repository:        nil,
		Provider:          provider,
		states:            nil,
	}

	s, err := sbuilder.BuildScenario(cfg)

	if err != nil {
		panic(err)
	}

	return cfg, provider, s
}

func (b *DefaultBot) HandleRequest(mf SerializedMessageFactory) {
	_, _, s := b.GetBaseActors(mf)
	token := s.Provider.GetToken()
	currentState := s.GetCurrentState(token)
	cmd := s.Provider.GetCommand(currentState)
	token = s.HandleCommand(cmd, currentState, token)
	b.TokenRepository.Persist(token)
}

func (b *DefaultBot) LogRequest(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(body))

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}
