package notifier

import (
	"os"
	"os/signal"

	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/facade"
	"github.com/multiversx/mx-chain-notifier-go/factory"
	"github.com/multiversx/mx-chain-notifier-go/metrics"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
)

var log = logger.GetOrCreate("notifierRunner")

type notifierRunner struct {
	configs config.Configs
}

// NewNotifierRunner create a new notifierRunner instance
func NewNotifierRunner(cfgs *config.Configs) (*notifierRunner, error) {
	if cfgs == nil {
		return nil, ErrNilConfigs
	}

	return &notifierRunner{
		configs: *cfgs,
	}, nil
}

// Start will trigger the notifier service
func (nr *notifierRunner) Start() error {
	externalMarshaller, err := marshalFactory.NewMarshalizer(nr.configs.MainConfig.General.ExternalMarshaller.Type)
	if err != nil {
		return err
	}

	lockService, err := factory.CreateLockService(nr.configs.MainConfig.General.CheckDuplicates, nr.configs.MainConfig.Redis)
	if err != nil {
		return err
	}

	publisher, err := factory.CreatePublisher(nr.configs.Flags.PublisherType, nr.configs.MainConfig, externalMarshaller)
	if err != nil {
		return err
	}

	wsHandler, err := factory.CreateWSHandler(nr.configs.Flags.PublisherType, publisher, externalMarshaller)
	if err != nil {
		return err
	}

	statusMetricsHandler := metrics.NewStatusMetrics()

	argsEventsHandler := factory.ArgsEventsHandlerFactory{
		CheckDuplicates:      nr.configs.MainConfig.General.CheckDuplicates,
		Locker:               lockService,
		Publisher:            publisher,
		StatusMetricsHandler: statusMetricsHandler,
	}
	eventsHandler, err := factory.CreateEventsHandler(argsEventsHandler)
	if err != nil {
		return err
	}

	eventsInterceptor, err := factory.CreateEventsInterceptor(nr.configs.MainConfig.General)
	if err != nil {
		return err
	}

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler:        eventsHandler,
		APIConfig:            nr.configs.MainConfig.ConnectorApi,
		WSHandler:            wsHandler,
		EventsInterceptor:    eventsInterceptor,
		StatusMetricsHandler: statusMetricsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return err
	}

	webServer, err := factory.CreateWebServerHandler(facade, nr.configs)
	if err != nil {
		return err
	}

	wsConnector, err := factory.CreateWSObserverConnector(nr.configs.MainConfig.WebSocketConnector, facade)
	if err != nil {
		return err
	}

	publisher.Run()

	err = webServer.Run()
	if err != nil {
		return err
	}

	err = waitForGracefulShutdown(webServer, publisher, wsConnector)
	if err != nil {
		return err
	}
	log.Debug("closing eventNotifier proxy...")

	return nil
}

func waitForGracefulShutdown(
	server shared.WebServerHandler,
	publisher rabbitmq.PublisherService,
	wsConnector process.WSClient,
) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	err := server.Close()
	if err != nil {
		return err
	}

	err = wsConnector.Close()
	if err != nil {
		return err
	}

	err = publisher.Close()
	if err != nil {
		return err
	}

	return nil
}
