package config

import "github.com/ElrondNetwork/elrond-go-core/core"

type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
	PubSub       PubSubConfig
}

type ConnectorApiConfig struct {
	ApiType      string
	Port         string
	HubType      string
	DispatchType string
	Username     string
	Password     string
}

type PubSubConfig struct {
	Url     string
	Channel string
}

// LoadConfig return a GeneralConfig instance by reading the provided toml file
func LoadConfig(filePath string) (*GeneralConfig, error) {
	cfg := &GeneralConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
