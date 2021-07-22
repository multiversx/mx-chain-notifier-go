package handlers

import (
	"net/http"

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
	endpoints   []EndpointHandler
}

// NewEventsHandler registers handlers for the /events group
func NewEventsHandler(
	notifierHub dispatcher.Hub,
	groupHandler *groupHandler,
) error {
	h := &eventsHandler{notifierHub: notifierHub}

	h.endpoints = []EndpointHandler{
		{Method: http.MethodPost, Path: pushEventsEndpoint, HandlerFunc: h.pushEvents},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseEventsEndpoint,
		Middleware:       []gin.HandlerFunc{},
		EndpointHandlers: h.endpoints,
	}

	groupHandler.AddGroupHandler(endpointGroupHandler)

	return nil
}

func (h *eventsHandler) pushEvents(c *gin.Context) {
	var events []data.Event

	err := c.Bind(&events)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if events != nil {
		h.notifierHub.BroadcastChan() <- events
	}

	JsonResponse(c, http.StatusOK, nil, "")
}
