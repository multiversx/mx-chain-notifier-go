package config

import "github.com/ElrondNetwork/elrond-go/core"

const configPath = "./config/config.toml"

type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
}

type ConnectorApiConfig struct {
	Port         string
	HubType      string
	DispatchType string
}

func LoadConfig() (*GeneralConfig, error) {
	cfg := &GeneralConfig{}
	err := core.LoadTomlFile(cfg, configPath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
