package handlers

import (
	"net/http"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"

	"github.com/gin-gonic/gin"
)

const (
	baseEventsEndpoint = "/events"
	pushEventsEndpoint = "/push"
)

type eventsHandler struct {
	notifierHub dispatcher.Hub
	config      config.ConnectorApiConfig
}

// NewEventsHandler registers handlers for the /events group
func NewEventsHandler(
	notifierHub dispatcher.Hub,
	groupHandler *GroupHandler,
	config config.ConnectorApiConfig,
) error {
	h := &eventsHandler{
		notifierHub: notifierHub,
		config:      config,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: pushEventsEndpoint, HandlerFunc: h.pushEvents},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseEventsEndpoint,
		Middlewares:      h.createMiddlewares(),
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)

	return nil
}

func (h *eventsHandler) pushEvents(c *gin.Context) {
	var blockEvents data.BlockEvents

	err := c.Bind(&blockEvents)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if blockEvents.Events != nil {
		h.notifierHub.BroadcastChan() <- blockEvents.Events
	}

	JsonResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) createMiddlewares() []gin.HandlerFunc {
	var middleware []gin.HandlerFunc

	if h.config.Username != "" && h.config.Password != "" {
		basicAuth := gin.BasicAuth(gin.Accounts{
			h.config.Username: h.config.Password,
		})
		middleware = append(middleware, basicAuth)
	}

	return middleware
}
