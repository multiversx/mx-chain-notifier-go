package config

import "github.com/ElrondNetwork/elrond-go-core/core"

type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
	PubSub       PubSubConfig
	RabbitMQ     RabbitMQConfig
}

type ConnectorApiConfig struct {
	Port            string
	HubType         string
	DispatchType    string
	Username        string
	Password        string
	CheckDuplicates bool
}

type PubSubConfig struct {
	Url         string
	Channel     string
	MasterName  string
	SentinelUrl string
}

type RabbitMQConfig struct {
	Url                     string
	EventsExchange          string
	RevertEventsExchange    string
	FinalizedEventsExchange string
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
