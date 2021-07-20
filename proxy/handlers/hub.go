package handlers

import (
	"net/http"
	"strings"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/gin-gonic/gin"
)

const (
	separator   = ":"
	hubCommon   = "common"
	dispatchAll = "dispatch:*"
	dispatchWs  = "websocket"
	dispatchGql = "graphql"

	baseHubEndpoint         = "/hub"
	websocketEndpoint       = "/ws"
	gqlQueryEndpoint        = "/query"
	gqlSubscriptionEndpoint = "/subscription"
)

var availableHubDelegates = map[string]func() dispatcher.Hub{
	"default-custom": func() dispatcher.Hub {
		return nil
	},
}

type hubHandler struct {
	notifierHub dispatcher.Hub
	endpoints   []EndpointHandler
}

// NewHubHandler registers handlers for the /hub group
// It only registers the specified hub implementation and its corresponding dispatchers
func NewHubHandler(
	config *config.GeneralConfig,
	groupHandler *groupHandler,
) (*hubHandler, error) {
	notifierHub, err := makeHub(config.ConnectorApi.HubType)
	if err != nil {
		return nil, err
	}

	h := &hubHandler{
		notifierHub: notifierHub,
		endpoints:   []EndpointHandler{},
	}

	handlers := h.getDispatchHandlers(config.ConnectorApi.DispatchType)
	groupHandler.AddEndpointHandlers(baseHubEndpoint, handlers)
	return h, nil
}

func (h *hubHandler) GetHub() dispatcher.Hub {
	return h.notifierHub
}

func (h *hubHandler) getDispatchHandlers(dispatchType string) []EndpointHandler {
	gqlServer := gql.NewGraphQLServer(h.notifierHub)

	var registerGql = func() {
		h.appendEndpointHandler(http.MethodPost, gqlQueryEndpoint, h.gqlHandler(gqlServer))
		h.appendEndpointHandler(http.MethodPost, gqlSubscriptionEndpoint, h.gqlHandler(gqlServer))
	}

	var registerWs = func() {
		h.appendEndpointHandler(http.MethodGet, websocketEndpoint, h.wsHandler)
	}

	switch dispatchType {
	case dispatchAll:
		registerWs()
		registerGql()
	case dispatchGql:
		registerGql()
	case dispatchWs:
		registerWs()
	}

	return h.endpoints
}

func (h *hubHandler) wsHandler(c *gin.Context) {
	ws.Serve(h.notifierHub, c.Writer, c.Request)
}

func (h *hubHandler) gqlHandler(gqlServer *gqlHandler.Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		gqlServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *hubHandler) appendEndpointHandler(method, path string, handler gin.HandlerFunc) {
	h.endpoints = append(h.endpoints, EndpointHandler{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
	})
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
