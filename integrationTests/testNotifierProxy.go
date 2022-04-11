package integrationTests

import (
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
	"github.com/ElrondNetwork/notifier-go/facade"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/ElrondNetwork/notifier-go/process"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
	"github.com/ElrondNetwork/notifier-go/redis"
)

type testNotifier struct {
	Facade         FacadeHandler
	Publisher      PublisherHandler
	WSHandler      dispatcher.WSHandler
	RedisClient    *mocks.RedisClientMock
	RabbitMQClient *mocks.RabbitClientMock
}

// NewTestNotifierWithWS will create a notifier instance for websockets flow
func NewTestNotifierWithWS(cfg *config.GeneralConfig) (*testNotifier, error) {
	redisClient := mocks.NewRedisClientMock()
	locker, err := redis.NewRedlockWrapper(redisClient)
	if err != nil {
		return nil, err
	}

	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	publisher, err := hub.NewCommonHub(args)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Config:              cfg.ConnectorApi,
		Locker:              locker,
		Publisher:           publisher,
		MaxLockerConRetries: 2,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	upgrader, err := ws.NewWSUpgraderWrapper(1024, 1024)
	if err != nil {
		return nil, err
	}
	wsHandlerArgs := ws.ArgsWebSocketProcessor{
		Hub:      publisher,
		Upgrader: upgrader,
	}
	wsHandler, err := ws.NewWebSocketProcessor(wsHandlerArgs)
	if err != nil {
		return nil, err
	}

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
		APIConfig:     cfg.ConnectorApi,
		WSHandler:     wsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return nil, err
	}

	return &testNotifier{
		Facade:         facade,
		Publisher:      publisher,
		WSHandler:      wsHandler,
		RedisClient:    redisClient,
		RabbitMQClient: mocks.NewRabbitClientMock(),
	}, nil
}

// NewTestNotifierWithRabbitMq will create a notifier instance with rabbitmq
func NewTestNotifierWithRabbitMq(cfg *config.GeneralConfig) (*testNotifier, error) {
	redisClient := mocks.NewRedisClientMock()
	locker, err := redis.NewRedlockWrapper(redisClient)
	if err != nil {
		return nil, err
	}

	rabbitmqMock := mocks.NewRabbitClientMock()
	publisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client: rabbitmqMock,
		Config: cfg.RabbitMQ,
	}
	publisher, err := rabbitmq.NewRabbitMqPublisher(publisherArgs)

	argsEventsHandler := process.ArgsEventsHandler{
		Config:              cfg.ConnectorApi,
		Locker:              locker,
		Publisher:           publisher,
		MaxLockerConRetries: 2,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	wsHandler := &disabled.WSHandler{}
	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
		APIConfig:     cfg.ConnectorApi,
		WSHandler:     wsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return nil, err
	}

	return &testNotifier{
		Facade:         facade,
		Publisher:      publisher,
		WSHandler:      wsHandler,
		RedisClient:    redisClient,
		RabbitMQClient: rabbitmqMock,
	}, nil
}

func GetDefaultConfigs() *config.GeneralConfig {
	return &config.GeneralConfig{
		ConnectorApi: config.ConnectorApiConfig{
			Port:            "8081",
			Username:        "user",
			Password:        "pass",
			CheckDuplicates: false,
		},
		Redis: config.RedisConfig{
			Url:            "redis://localhost:6379",
			Channel:        "pub-sub",
			MasterName:     "mymaster",
			SentinelUrl:    "localhost:26379",
			ConnectionType: "sentinel",
		},
		RabbitMQ: config.RabbitMQConfig{
			Url:                     "amqp://guest:guest@localhost:5672",
			EventsExchange:          "all_events",
			RevertEventsExchange:    "revert_events",
			FinalizedEventsExchange: "finalized_events",
		},
		Flags: &config.FlagsConfig{
			LogLevel:          "*:INFO",
			SaveLogFile:       false,
			GeneralConfigPath: "./config/config.toml",
			WorkingDir:        "",
			APIType:           "notifier",
		},
	}
}
