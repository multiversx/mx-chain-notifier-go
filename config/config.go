package config

import "github.com/multiversx/mx-chain-core-go/core"

// Configs holds all configs
type Configs struct {
	GeneralConfig   GeneralConfig
	ApiRoutesConfig APIRoutesConfig
	Flags           FlagsConfig
}

// GeneralConfig defines the config setup based on main config file
type GeneralConfig struct {
	ConnectorApi ConnectorApiConfig
	Redis        RedisConfig
	RabbitMQ     RabbitMQConfig
}

// ConnectorApiConfig maps the connector configuration
type ConnectorApiConfig struct {
	Port            string
	Username        string
	Password        string
	CheckDuplicates bool
}

// APIRoutesConfig holds the configuration related to Rest API routes
type APIRoutesConfig struct {
	APIPackages map[string]APIPackageConfig
}

// APIPackageConfig holds the configuration for the routes of each package
type APIPackageConfig struct {
	Routes []RouteConfig
}

// RouteConfig holds the configuration for a single route
type RouteConfig struct {
	Name string
	Open bool
	Auth bool
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
	BlockEventsExchange     RabbitMQExchangeConfig
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
	APIConfigPath     string
	WorkingDir        string
	APIType           string
	RestApiInterface  string
}

// LoadGeneralConfig return a GeneralConfig instance by reading the provided toml file
func LoadGeneralConfig(filePath string) (*GeneralConfig, error) {
	cfg := &GeneralConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

// LoadAPIConfig return a APIRoutesConfig instance by reading the provided toml file
func LoadAPIConfig(filePath string) (*APIRoutesConfig, error) {
	cfg := &APIRoutesConfig{}
	err := core.LoadTomlFile(cfg, filePath)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
