package config

import "github.com/ElrondNetwork/elrond-go/core"

type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
}

type ConnectorApiConfig struct {
	Port         string
	HubType      string
	DispatchType string
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
