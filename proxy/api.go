package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/proxy/handlers"

	"github.com/gin-gonic/gin"
)

type webServer struct {
	router        *gin.Engine
	notifierHub   dispatcher.Hub
	generalConfig *config.GeneralConfig
}

func NewWebServer(generalConfig *config.GeneralConfig) (*webServer, error) {
	router := gin.Default()
	groupHandler := handlers.NewGroupHandler()

	hubHandler, err := handlers.NewHubHandler(generalConfig, groupHandler)
	if err != nil {
		return nil, err
	}

	notifierHub := hubHandler.GetHub()

	err = handlers.NewEventsHandler(notifierHub, groupHandler)
	if err != nil {
		return nil, err
	}

	groupHandler.RegisterEndpoints(router)

	return &webServer{
		router:        router,
		notifierHub:   notifierHub,
		generalConfig: generalConfig,
	}, nil
}

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
