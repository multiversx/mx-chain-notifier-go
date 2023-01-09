package notifier

import (
	"os"
	"os/signal"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/gin"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/facade"
	"github.com/multiversx/mx-chain-notifier-go/factory"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
)

var log = logger.GetOrCreate("notifierRunner")

type notifierRunner struct {
	configs *config.GeneralConfig
}

// NewNotifierRunner create a new notifierRunner instance
func NewNotifierRunner(cfgs *config.GeneralConfig) (*notifierRunner, error) {
	if cfgs == nil {
		return nil, ErrNilConfigs
	}

	return &notifierRunner{
		configs: cfgs,
	}, nil
}

// Start will trigger the notifier service
func (nr *notifierRunner) Start() error {
	lockService, err := factory.CreateLockService(nr.configs.ConnectorApi.CheckDuplicates, nr.configs.Redis)
	if err != nil {
		return err
	}

	publisher, err := factory.CreatePublisher(nr.configs.Flags.APIType, nr.configs)
	if err != nil {
		return err
	}

	hub, err := factory.CreateHub(nr.configs.Flags.APIType)
	if err != nil {
		return err
	}

	wsHandler, err := factory.CreateWSHandler(nr.configs.Flags.APIType, hub)
	if err != nil {
		return err
	}

	argsEventsHandler := factory.ArgsEventsHandlerFactory{
		APIConfig:    nr.configs.ConnectorApi,
		Locker:       lockService,
		MqPublisher:  publisher,
		HubPublisher: hub,
		APIType:      nr.configs.Flags.APIType,
	}
	eventsHandler, err := factory.CreateEventsHandler(argsEventsHandler)
	if err != nil {
		return err
	}

	eventsInterceptor, err := factory.CreateEventsInterceptor()
	if err != nil {
		return err
	}

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler:     eventsHandler,
		APIConfig:         nr.configs.ConnectorApi,
		WSHandler:         wsHandler,
		EventsInterceptor: eventsInterceptor,
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
