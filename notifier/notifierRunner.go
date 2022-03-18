package notifier

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/api/gin"
	"github.com/ElrondNetwork/notifier-go/api/groups"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/facade"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
	"github.com/ElrondNetwork/notifier-go/redis"
)

var log = logger.GetOrCreate("notifierRunner")

const (
	separator     = ":"
	hubCommon     = "common"
	rabbitAPIType = "rabbit-api"
	notifierType  = "notifier"
)

var availableHubDelegates = map[string]func() dispatcher.Hub{
	"default-custom": func() dispatcher.Hub {
		return nil
	},
}

var (
	backgroundContextTimeout = 5 * time.Second
)

// ErrMakeCustomHub -
var ErrMakeCustomHub = errors.New("failed to make custom hub")

type notifierRunner struct {
	configs *config.GeneralConfig
	apiType string
}

// NewNotifierRunner create a new notifierRunner instance
func NewNotifierRunner(typeValue string, cfgs *config.GeneralConfig) (*notifierRunner, error) {
	if cfgs == nil {
		return nil, fmt.Errorf("nil configs provided")
	}

	return &notifierRunner{
		configs: cfgs,
		apiType: typeValue,
	}, nil
}

func (nr *notifierRunner) Start() error {
	webServer, notifierHub, err := nr.initWebserver(nr.configs)
	if err != nil {
		return err
	}

	notifierHubRun(notifierHub)

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

func notifierHubRun(notifierHub dispatcher.Hub) {
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

func (nr *notifierRunner) initWebserver(cfg *config.GeneralConfig) (
	shared.HTTPServerHandler,
	dispatcher.Hub,
	error,
) {
	switch nr.apiType {
	case rabbitAPIType:
		return NewObserverToRabbitAPI(cfg)
	case notifierType:
		return NewNotifierAPI(cfg)
	default:
		return nil, nil, ErrInvalidAPIType
	}
}

type disabledFacade struct{}

func (df *disabledFacade) IsInterfaceNil() bool { return df == nil }

// NewNotifierAPI launches a notifier api - exposing a clients hub
func NewNotifierAPI(config *config.GeneralConfig) (shared.HTTPServerHandler, dispatcher.Hub, error) {
	notifierHub, err := makeHub(config.ConnectorApi.HubType)
	if err != nil {
		return nil, nil, err
	}

	hubHandler, err := groups.NewHubHandler(config, notifierHub)
	if err != nil {
		return nil, nil, err
	}

	disabledLockService := disabled.NewDisabledRedlockWrapper()

	argsEventsHandler := ArgsEventsHandler{
		Config:    config.ConnectorApi,
		Locker:    disabledLockService,
		Publisher: notifierHub,
	}
	eventsHandler, err := NewEventsHandler(argsEventsHandler)

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)

	eventsGroup, err := groups.NewEventsHandler(
		facade,
		config.ConnectorApi,
	)
	if err != nil {
		return nil, nil, err
	}

	webServerArgs := gin.ArgsWebServerHandler{
		Facade: &disabledFacade{},
		Config: config,
	}
	server, err := gin.NewWebServerHandler(webServerArgs)
	if err != nil {
		return nil, nil, err
	}

	server.AddGroup("events", eventsGroup)
	server.AddGroup("hub", hubHandler)

	return server, notifierHub, nil
}

// NewObserverToRabbitAPI launches an observer api - pushing data to rabbitMQ exchanges
func NewObserverToRabbitAPI(config *config.GeneralConfig) (shared.HTTPServerHandler, dispatcher.Hub, error) {
	pubsubClient, err := redis.CreateFailoverClient(config.PubSub)
	if err != nil {
		return nil, nil, err
	}

	rabbitClient, err := rabbitmq.NewRabbitMQClient(context.Background(), config.RabbitMQ.Url)
	if err != nil {
		return nil, nil, err
	}
	rabbitPublisher := rabbitmq.NewRabbitMqPublisher(rabbitClient, config.RabbitMQ)

	var lockService groups.LockService
	if config.ConnectorApi.CheckDuplicates {
		lockService, err = redis.NewRedlockWrapper(pubsubClient)
		if err != nil {
			return nil, nil, err
		}
	} else {
		lockService = disabled.NewDisabledRedlockWrapper()
	}

	argsEventsHandler := ArgsEventsHandler{
		Config:    config.ConnectorApi,
		Locker:    lockService,
		Publisher: rabbitPublisher,
	}
	eventsHandler, err := NewEventsHandler(argsEventsHandler)

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)

	eventsGroup, err := groups.NewEventsHandler(
		facade,
		config.ConnectorApi,
	)
	if err != nil {
		return nil, nil, err
	}

	webServerArgs := gin.ArgsWebServerHandler{
		Facade: &disabledFacade{},
		Config: config,
	}
	server, err := gin.NewWebServerHandler(webServerArgs)
	if err != nil {
		return nil, nil, err
	}

	server.AddGroup("events", eventsGroup)

	return server, rabbitPublisher, nil
}

func makeHub(hubType string) (dispatcher.Hub, error) {
	if hubType == hubCommon {
		eventFilter := filters.NewDefaultFilter()
		commonHub := hub.NewCommonHub(eventFilter)
		return commonHub, nil
	}

	hubConfig := strings.Split(hubType, separator)
	return tryMakeCustomHubForID(hubConfig[1])
}

func tryMakeCustomHubForID(id string) (dispatcher.Hub, error) {
	if makeHubFunc, ok := availableHubDelegates[id]; ok {
		customHub := makeHubFunc()
		return customHub, nil
	}

	return nil, ErrMakeCustomHub
}
