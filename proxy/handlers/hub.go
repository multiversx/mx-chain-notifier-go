package handlers

import (
	"net/http"

	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/gin-gonic/gin"
)

const (
	baseHubEndpoint   = "/hub"
	websocketEndpoint = "/ws"
)

type hubHandler struct {
	notifierHub dispatcher.Hub
}

// NewHubHandler registers handlers for the /hub group
// It only registers the specified hub implementation and its corresponding dispatchers
func NewHubHandler(groupHandler *GroupHandler) *hubHandler {
	notifierHub := makeHub()

	h := &hubHandler{
		notifierHub: notifierHub,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodGet, Path: websocketEndpoint, HandlerFunc: h.wsHandler},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseHubEndpoint,
		Middlewares:      []gin.HandlerFunc{},
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)

	return h
}

func (h *hubHandler) GetHub() dispatcher.Hub {
	return h.notifierHub
}

func (h *hubHandler) wsHandler(c *gin.Context) {
	ws.Serve(h.notifierHub, c.Writer, c.Request)
}

func makeHub() dispatcher.Hub {
	eventFilter := filters.NewDefaultFilter()

	return hub.NewCommonHub(eventFilter)
}
