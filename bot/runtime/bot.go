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
}

func (b *DefaultBot) HandleRequest(mf *DefaultSerializedMessageFactory) {
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

func (b *DefaultBot) LogRequest(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(body))

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}
