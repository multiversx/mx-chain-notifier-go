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
)

var log = logger.GetOrCreate("notifierRunner")

var (
	backgroundContextTimeout = 5 * time.Second
)

type notifierRunner struct {
	configs *config.GeneralConfig
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

	pubsubHandler, err := factory.CreatePubSubHandler(nr.apiType, nr.configs)
	if err != nil {
		return err
	}

	argsEventsHandler := ArgsEventsHandler{
		Config:    nr.configs.ConnectorApi,
		Locker:    lockService,
		Publisher: pubsubHandler,
	}
	eventsHandler, err := NewEventsHandler(argsEventsHandler)

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
		APIConfig:     nr.configs.ConnectorApi,
		Hub:           pubsubHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)

	webServerArgs := gin.ArgsWebServerHandler{
		Facade: facade,
		Config: nr.configs.ConnectorApi,
		Type:   nr.apiType,
	}
	webServer, err := gin.NewWebServerHandler(webServerArgs)
	if err != nil {
		return err
	}

	startHubHandler(pubsubHandler)

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

func startHubHandler(notifierHub dispatcher.Hub) {
	go notifierHub.Run()
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
