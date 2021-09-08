package runtime

import (
	"bot-daedalus/config"
	"fmt"
)

func GetProvider(config config.ProviderConfig, tf TokenFactory, mf SerializedMessageFactory) (ChatProvider, error) {
	if config.Name == "telegram" {
		return &TelegramProvider{
			tf,
			mf,
			config,
			mf.GetSerializedMessage(config),
		}, nil
	}

	return nil, fmt.Errorf("wrong provider type passed")
}
