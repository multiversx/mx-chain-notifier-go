package config

import "github.com/multiversx/mx-chain-core-go/core"

// GeneralConfig defines the config setup based on main config file
type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
	Redis        RedisConfig
	RabbitMQ     RabbitMQConfig
	Flags        *FlagsConfig
}

// ConnectorApiConfig maps the connector configuration
type ConnectorApiConfig struct {
	Port            string
	Username        string
	Password        string
	CheckDuplicates bool
}

// RedisConfig maps the redis configuration
type RedisConfig struct {
	Url            string
	Channel        string
	MasterName     string
	SentinelUrl    string
	ConnectionType string
	TTL            uint32
}

// RabbitMQConfig maps the rabbitMQ configuration
type RabbitMQConfig struct {
	Url                     string
	EventsExchange          RabbitMQExchangeConfig
	RevertEventsExchange    RabbitMQExchangeConfig
	FinalizedEventsExchange RabbitMQExchangeConfig
	BlockTxsExchange        RabbitMQExchangeConfig
	BlockScrsExchange       RabbitMQExchangeConfig
}

// RabbitMQExchangeConfig holds the configuration for a rabbitMQ exchange
type RabbitMQExchangeConfig struct {
	Name string
	Type string
}

// FlagsConfig holds the values for CLI flags
type FlagsConfig struct {
	LogLevel          string
	SaveLogFile       bool
	GeneralConfigPath string
	WorkingDir        string
	APIType           string
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
