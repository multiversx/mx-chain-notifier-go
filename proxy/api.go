package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/proxy/handlers"
	"github.com/ElrondNetwork/notifier-go/pubsub"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

type WebServer struct {
	router        *gin.Engine
	notifierHub   dispatcher.Hub
	groupHandler  *handlers.GroupHandler
	generalConfig *config.GeneralConfig
}

// NewWebServer creates and configures an instance of WebServer
func newWebServer(generalConfig *config.GeneralConfig) *WebServer {
	router := gin.Default()
	router.Use(cors.Default())

	return &WebServer{
		router:        router,
		generalConfig: generalConfig,
		groupHandler:  handlers.NewGroupHandler(),
	}
}

// NewNotifierApi -
func NewNotifierApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	hubHandler, err := handlers.NewHubHandler(config, server.groupHandler)
	if err != nil {
		return nil, err
	}

	notifierHub := hubHandler.GetHub()
	server.notifierHub = notifierHub

	err = handlers.NewEventsHandler(notifierHub, server.groupHandler, config.ConnectorApi)
	if err != nil {
		return nil, err
	}

	server.groupHandler.RegisterEndpoints(server.router)

	return server, nil
}

// NewObserverApi -
func NewObserverApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	pubsubHub := pubsub.NewHubPublisher(ctx, config.PubSub)
	server.notifierHub = pubsubHub

	err := handlers.NewEventsHandler(pubsubHub, server.groupHandler, config.ConnectorApi)
	if err != nil {
		return nil, err
	}

	server.groupHandler.RegisterEndpoints(server.router)

	return server, nil
}

// NewClientApi -
func NewClientApi(config *config.GeneralConfig) (*WebServer, error) {
	server := newWebServer(config)

	hubHandler, err := handlers.NewHubHandler(config, server.groupHandler)
	if err != nil {
		return nil, err
	}

	notifierHub := hubHandler.GetHub()
	server.notifierHub = notifierHub

	pubsubListener := pubsub.NewHubSubscriber(ctx, config.PubSub, notifierHub)
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
