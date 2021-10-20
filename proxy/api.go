package proxy

import (
	"context"
	"fmt"
	"github.com/ElrondNetwork/notifier-go/pubsub/disabled"
	"net/http"
	"strings"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/proxy/handlers"
	"github.com/ElrondNetwork/notifier-go/pubsub"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// WebServer is a wrapper for gin.Engine, holding additional components
type WebServer struct {
	router        *gin.Engine
	notifierHub   dispatcher.Hub
	groupHandler  *handlers.GroupHandler
	generalConfig *config.GeneralConfig
}

// newWebServer creates and configures an instance of WebServer
func newWebServer(generalConfig *config.GeneralConfig) *WebServer {
	router := gin.Default()
	router.Use(cors.Default())

	return &WebServer{
		router:        router,
		generalConfig: generalConfig,
		groupHandler:  handlers.NewGroupHandler(),
	}
}

// NewNotifierApi launches a notifier api - exposing a clients hub
func NewNotifierApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	hubHandler, err := handlers.NewHubHandler(config, server.groupHandler)
	if err != nil {
		return nil, err
	}

	notifierHub := hubHandler.GetHub()
	server.notifierHub = notifierHub

	disabledLockService := disabled.NewDisabledRedlockWrapper()
	err = handlers.NewEventsHandler(
		notifierHub,
		server.groupHandler,
		config.ConnectorApi,
		disabledLockService,
	)
	if err != nil {
		return nil, err
	}

	server.groupHandler.RegisterEndpoints(server.router)

	return server, nil
}

// NewObserverToRabbitApi launches an observer api - pushing data to rabbitMQ exchanges
func NewObserverToRabbitApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	pubsubClient, err := pubsub.CreateFailoverClient(config.PubSub)
	if err != nil {
		return nil, err
	}

	rabbitPublisher, err := rabbitmq.NewRabbitMqPublisher(ctx, config.RabbitMQ)
	if err != nil {
		return nil, err
	}

	server.notifierHub = rabbitPublisher

	var lockService pubsub.LockService
	if config.ConnectorApi.CheckDuplicates {
		lockService = pubsub.NewRedlockWrapper(ctx, pubsubClient)
	} else {
		lockService = disabled.NewDisabledRedlockWrapper()
	}

	err = handlers.NewEventsHandler(
		rabbitPublisher,
		server.groupHandler,
		config.ConnectorApi,
		lockService,
	)
	if err != nil {
		return nil, err
	}

	server.groupHandler.RegisterEndpoints(server.router)

	return server, nil
}

// NewObserverApi launches an observer api - exposing a pubsub publisher hub
func NewObserverApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	pubsubClient := pubsub.CreatePubsubClient(config.PubSub)

	pubsubHub := pubsub.NewHubPublisher(ctx, config.PubSub, pubsubClient)
	server.notifierHub = pubsubHub

	redlock := pubsub.NewRedlockWrapper(ctx, pubsubClient)

	err := handlers.NewEventsHandler(pubsubHub, server.groupHandler, config.ConnectorApi, redlock)
	if err != nil {
		return nil, err
	}

	server.groupHandler.RegisterEndpoints(server.router)

	return server, nil
}

// NewClientApi launches a client api - exposing a pubsub subscriber hub
func NewClientApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	hubHandler, err := handlers.NewHubHandler(config, server.groupHandler)
	if err != nil {
		return nil, err
	}

	notifierHub := hubHandler.GetHub()
	server.notifierHub = notifierHub

	pubsubClient := pubsub.CreatePubsubClient(config.PubSub)

	pubsubListener := pubsub.NewHubSubscriber(ctx, config.PubSub, notifierHub, pubsubClient)
	pubsubListener.Subscribe()

	server.groupHandler.RegisterEndpoints(server.router)

	return server, nil
}

// Run starts the server and the Hub as goroutines
// It returns an instance of http.Server
func (w *WebServer) Run() *http.Server {
	go w.notifierHub.Run()

	port := w.generalConfig.ConnectorApi.Port
	if !strings.Contains(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	server := &http.Server{
		Addr:    port,
		Handler: w.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server
}
