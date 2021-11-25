package runtime

import (
	"fmt"
)

func GetProvider(config ProviderConfig, scenarioName string, tf TokenFactory, mf SerializedMessageFactory, tr TokenRepository) (ChatProvider, error) {
	if config.Name == "telegram" {
		return &TelegramProvider{
			tf,
			scenarioName,
			mf,
			config,
			mf.GetSerializedMessage(config),
			tr,
		}, nil
	}

	return nil, fmt.Errorf("wrong provider type passed")
}
