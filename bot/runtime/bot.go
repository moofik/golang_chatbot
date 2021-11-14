package runtime

import (
	"bot-daedalus/config"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Bot interface {
	HandleRequest(mf SerializedMessageFactory)
}

type DefaultBot struct {
	ScenarioPath    string
	ScenarioName    string
	TokenFactory    TokenFactory
	TokenRepository TokenRepository
	ActionRegistry  func(string, map[string]string) Action
	CommandRegistry func(string, string, []interface{}) Command
}

func (b *DefaultBot) HandleRequest(mf SerializedMessageFactory) {
	cfg := config.GetScenarioConfig(b.ScenarioPath, b.ScenarioName)
	provider, _ := GetProvider(cfg.ProviderConfig, cfg.Name, b.TokenFactory, mf, b.TokenRepository)

	sbuilder := ScenarioBuilder{
		ActionRegistry:  b.ActionRegistry,
		CommandRegistry: b.CommandRegistry,
		Repository:      nil,
		Provider:        provider,
		states:          nil,
	}

	s, err := sbuilder.BuildScenario(b.ScenarioPath, b.ScenarioName)

	if err != nil {
		panic(err)
	}

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
