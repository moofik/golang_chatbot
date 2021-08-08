package runtime

import (
	"bot-daedalus/config"
	"github.com/gin-gonic/gin"
)

type Bot interface {
	HandleRequest(c *gin.Context)
}

type DefaultBot struct {
	message Message
}

func (b *DefaultBot) HandleRequest(c *gin.Context) {
	cfg := config.GetScenarioConfig()
	provider, _ := GetProvider(cfg.ProviderConfig, c)
	s := Scenario{cfg.StateMachineConfig, provider}
	s.HandleCommand()
}
