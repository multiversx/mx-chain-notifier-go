package notifier

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/api/gin"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/facade"
	"github.com/ElrondNetwork/notifier-go/factory"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
)

var log = logger.GetOrCreate("notifierRunner")

var (
	backgroundContextTimeout = 5 * time.Second
)

type notifierRunner struct {
	configs *config.GeneralConfig
	// TODO: refactor cmd/main to include apiType (cli args) in a separate config struct
	apiType common.APIType
}

// NewNotifierRunner create a new notifierRunner instance
func NewNotifierRunner(typeValue common.APIType, cfgs *config.GeneralConfig) (*notifierRunner, error) {
	if cfgs == nil {
		return nil, fmt.Errorf("nil configs provided")
	}

	return &notifierRunner{
		configs: cfgs,
		apiType: typeValue,
	}, nil
}

// Start will trigger the notifier service
func (nr *notifierRunner) Start() error {
	lockService, err := factory.CreateLockService(nr.apiType, nr.configs)
	if err != nil {
		return err
	}

	publisher, err := factory.CreatePublisher(nr.apiType, nr.configs)
	if err != nil {
		return err
	}

	hubType := common.HubType(nr.configs.ConnectorApi.HubType)
	hub, err := factory.CreateHub(nr.apiType, hubType)
	if err != nil {
		return err
	}

	argsEventsHandler := ArgsEventsHandler{
		Config:    nr.configs.ConnectorApi,
		Locker:    lockService,
		Publisher: publisher,
	}
	eventsHandler, err := NewEventsHandler(argsEventsHandler)
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
		Type:   nr.apiType,
	}
	webServer, err := gin.NewWebServerHandler(webServerArgs)
	if err != nil {
		return err
	}

	startHubHandler(hub, publisher)

	err = webServer.Run()
	if err != nil {
		return err
	}

	err = waitForGracefulShutdown(webServer)
	if err != nil {
		return err
	}
	log.Debug("closing eventNotifier proxy...")

	return nil
}

func startHubHandler(hub dispatcher.Hub, publisher rabbitmq.PublisherService) {
	go hub.Run()
	go publisher.Run()
}

func waitForGracefulShutdown(server shared.HTTPServerHandler) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	err := server.Close()
	if err != nil {
		return err
	}

	return nil
}
