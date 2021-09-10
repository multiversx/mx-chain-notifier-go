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

const (
	observerApiType = "observer-api"
	clientApiType   = "client-api"
	notifierType    = "notifier"
)

var ctx = context.Background()

type webServer struct {
	router        *gin.Engine
	notifierHub   dispatcher.Hub
	generalConfig *config.GeneralConfig
}

// NewWebServer creates and configures an instance of webServer
func NewWebServer(generalConfig *config.GeneralConfig) (*webServer, error) {
	router := gin.Default()
	router.Use(cors.Default())

	server := &webServer{
		router:        router,
		generalConfig: generalConfig,
	}

	err := server.registerHandlers()
	if err != nil {
		return nil, err
	}

	return server, nil
}

// Run starts the server and the Hub as goroutines
// It returns an instance of http.Server
func (w *webServer) Run() *http.Server {
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

func (w *webServer) registerHandlers() error {
	apiType := w.generalConfig.ConnectorApi.ApiType

	groupHandler := handlers.NewGroupHandler()

	if apiType == observerApiType {
		pubsubHub := pubsub.NewHubPublisher(ctx, w.generalConfig.PubSub)
		w.notifierHub = pubsubHub

		err := handlers.NewEventsHandler(pubsubHub, groupHandler, w.generalConfig.ConnectorApi)
		if err != nil {
			return err
		}
	}

	if apiType == clientApiType {
		hubHandler, err := handlers.NewHubHandler(w.generalConfig, groupHandler)
		if err != nil {
			return err
		}

		notifierHub := hubHandler.GetHub()
		w.notifierHub = notifierHub

		pubsubListener := pubsub.NewHubSubscriber(ctx, w.generalConfig.PubSub, notifierHub)
		pubsubListener.Subscribe()
	}

	if apiType == notifierType {
		hubHandler, err := handlers.NewHubHandler(w.generalConfig, groupHandler)
		if err != nil {
			return err
		}

		notifierHub := hubHandler.GetHub()
		w.notifierHub = notifierHub

		err = handlers.NewEventsHandler(notifierHub, groupHandler, w.generalConfig.ConnectorApi)
		if err != nil {
			return err
		}
	}

	groupHandler.RegisterEndpoints(w.router)

	return nil
}
