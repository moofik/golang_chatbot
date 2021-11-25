package runtime

import (
	"fmt"
	"github.com/spf13/viper"
)

type ProviderConfig struct {
	Name  string
	Token string
}

type ScenarioConfig struct {
	Name           string
	ProviderConfig ProviderConfig
	States         interface{}
}

func GetScenarioConfig(path string, name string) ScenarioConfig {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	providerConfig := ProviderConfig{
		Name:  viper.GetString("scenario.provider.name"),
		Token: viper.GetString("scenario.provider.token"),
	}

	return ScenarioConfig{
		Name:           viper.GetString("scenario.name"),
		ProviderConfig: providerConfig,
		States:         viper.Get("scenario.states"),
	}
}
