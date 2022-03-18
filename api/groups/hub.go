package groups

import (
	"net/http"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
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

type hubHandler struct {
	*baseGroup
	notifierHub           dispatcher.Hub
	additionalMiddlewares []gin.HandlerFunc
}

// NewHubHandler registers handlers for the /hub group
// It only registers the specified hub implementation and its corresponding dispatchers
func NewHubHandler(
	config *config.GeneralConfig,
	notifierHub dispatcher.Hub,
) (*hubHandler, error) {
	h := &hubHandler{
		baseGroup:             &baseGroup{},
		notifierHub:           notifierHub,
		additionalMiddlewares: make([]gin.HandlerFunc, 0),
	}

	endpoints := h.getDispatchHandlers(config.ConnectorApi.DispatchType)

	h.endpoints = endpoints

	return h, nil
}

// GetAdditionalMiddlewares return additional middlewares for this group
func (h *hubHandler) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return h.additionalMiddlewares
}

func (h *hubHandler) getDispatchHandlers(dispatchType string) []*shared.EndpointHandlerData {
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
	h.endpoints = append(h.endpoints, &shared.EndpointHandlerData{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *hubHandler) IsInterfaceNil() bool {
	return h == nil
}
