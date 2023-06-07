package notifier

import (
	"os"
	"os/signal"

	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/gin"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/facade"
	"github.com/multiversx/mx-chain-notifier-go/factory"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
)

var log = logger.GetOrCreate("notifierRunner")

type notifierRunner struct {
	configs *config.Config
}

// NewNotifierRunner create a new notifierRunner instance
func NewNotifierRunner(cfgs *config.Config) (*notifierRunner, error) {
	if cfgs == nil {
		return nil, ErrNilConfigs
	}

	return &notifierRunner{
		configs: cfgs,
	}, nil
}

// Start will trigger the notifier service
func (nr *notifierRunner) Start() error {
	externalMarshaller, err := marshalFactory.NewMarshalizer(nr.configs.General.ExternalMarshaller.Type)
	if err != nil {
		return err
	}
	internalMarshaller, err := marshalFactory.NewMarshalizer(nr.configs.General.InternalMarshaller.Type)
	if err != nil {
		return err
	}

	lockService, err := factory.CreateLockService(nr.configs.ConnectorApi.CheckDuplicates, nr.configs.Redis)
	if err != nil {
		return err
	}

	publisher, err := factory.CreatePublisher(nr.configs.Flags.APIType, nr.configs, externalMarshaller)
	if err != nil {
		return err
	}

	hub, err := factory.CreateHub(nr.configs.Flags.APIType)
	if err != nil {
		return err
	}

	wsHandler, err := factory.CreateWSHandler(nr.configs.Flags.APIType, hub, externalMarshaller)
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

	eventsInterceptor, err := factory.CreateEventsInterceptor(nr.configs.General)
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
		Facade:             facade,
		Config:             nr.configs.ConnectorApi,
		Type:               nr.configs.Flags.APIType,
		Marshaller:         externalMarshaller,
		InternalMarshaller: internalMarshaller,
	}
	webServer, err := gin.NewWebServerHandler(webServerArgs)
	if err != nil {
		return err
	}

	wsConnector, err := factory.CreateWSIndexer(nr.configs.WebSocketConnector, externalMarshaller, facade)
	if err != nil {
		return err
	}

	startHandlers(hub, publisher)

	err = webServer.Run()
	if err != nil {
		return err
	}

	err = waitForGracefulShutdown(webServer, publisher, hub, wsConnector)
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
	wsConnector process.WSClient,
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

	err = wsConnector.Close()
	if err != nil {
		return err
	}

	return nil
}
