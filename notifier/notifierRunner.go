package notifier

import (
	"fmt"
	"os"
	"os/signal"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/api/gin"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/facade"
	"github.com/ElrondNetwork/notifier-go/factory"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
)

var log = logger.GetOrCreate("notifierRunner")

// TODO: evaluate adding to configs
const maxLockerConRetries = 10

type notifierRunner struct {
	configs *config.GeneralConfig
}

// NewNotifierRunner create a new notifierRunner instance
func NewNotifierRunner(cfgs *config.GeneralConfig) (*notifierRunner, error) {
	if cfgs == nil {
		return nil, fmt.Errorf("nil configs provided")
	}

	return &notifierRunner{
		configs: cfgs,
	}, nil
}

// Start will trigger the notifier service
func (nr *notifierRunner) Start() error {
	lockService, err := factory.CreateLockService(nr.configs.Flags.APIType, nr.configs)
	if err != nil {
		return err
	}

	publisher, err := factory.CreatePublisher(nr.configs.Flags.APIType, nr.configs)
	if err != nil {
		return err
	}

	hub, err := factory.CreateHub(nr.configs.Flags.APIType, nr.configs.ConnectorApi.HubType)
	if err != nil {
		return err
	}

	argsEventsHandler := factory.ArgsEventsHandlerFactory{
		APIConfig:           nr.configs.ConnectorApi,
		Locker:              lockService,
		MqPublisher:         publisher,
		HubPublisher:        hub,
		APIType:             nr.configs.Flags.APIType,
		MaxLockerConRetries: maxLockerConRetries,
	}
	eventsHandler, err := factory.CreateEventsHandler(argsEventsHandler)
	if err != nil {
		return err
	}

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
		APIConfig:     nr.configs.ConnectorApi,
		Hub:           hub,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return err
	}

	webServerArgs := gin.ArgsWebServerHandler{
		Facade: facade,
		Config: nr.configs.ConnectorApi,
		Type:   nr.configs.Flags.APIType,
	}
	webServer, err := gin.NewWebServerHandler(webServerArgs)
	if err != nil {
		return err
	}

	startHandlers(hub, publisher)

	err = webServer.Run()
	if err != nil {
		return err
	}

	err = waitForGracefulShutdown(webServer, publisher, hub)
	if err != nil {
		return err
	}
	log.Debug("closing eventNotifier proxy...")

	return nil
}

func startHandlers(hub dispatcher.Hub, publisher rabbitmq.PublisherService) {
	hub.Run()
	publisher.Run()
}

func waitForGracefulShutdown(
	server shared.WebServerHandler,
	publisher rabbitmq.PublisherService,
	hub dispatcher.Hub,
) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	err := server.Close()
	if err != nil {
		return err
	}

	err = publisher.Close()
	if err != nil {
		return err
	}

	err = hub.Close()
	if err != nil {
		return err
	}

	return nil
}
