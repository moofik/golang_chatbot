package runtime

import (
	"bot-daedalus/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetProvider(config config.ProviderConfig, c *gin.Context) (ChatProvider, error)  {
	if config.Name == "telegram" {
		return &TelegramProvider{config, GetMessage(c, config)}, nil
	}

	return nil, fmt.Errorf("wrong provider type passed")
}
