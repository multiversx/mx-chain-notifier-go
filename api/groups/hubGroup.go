package groups

import (
	"fmt"
	"net/http"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/api/errors"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
	"github.com/gin-gonic/gin"
)

const (
	dispatchAll = "dispatch:*"
	dispatchWs  = "websocket"
	dispatchGql = "graphql"

	websocketEndpoint       = "/ws"
	gqlQueryEndpoint        = "/query"
	gqlSubscriptionEndpoint = "/subscription"
)

type hubGroup struct {
	*baseGroup
	facade                HubFacadeHandler
	additionalMiddlewares []gin.HandlerFunc
}

// NewHubGroup registers handlers for the /hub group
// It only registers the specified hub implementation and its corresponding dispatchers
func NewHubGroup(facade HubFacadeHandler) (*hubGroup, error) {
	if check.IfNil(facade) {
		return nil, fmt.Errorf("%w for hub group", errors.ErrNilFacadeHandler)
	}

	h := &hubGroup{
		baseGroup:             &baseGroup{},
		facade:                facade,
		additionalMiddlewares: make([]gin.HandlerFunc, 0),
	}

	endpoints, err := h.getDispatchHandlers()
	if err != nil {
		return nil, err
	}

	h.endpoints = endpoints

	return h, nil
}

// GetAdditionalMiddlewares return additional middlewares for this group
func (h *hubGroup) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return h.additionalMiddlewares
}

func (h *hubGroup) getDispatchHandlers() ([]*shared.EndpointHandlerData, error) {
	// TODO: handle graphql server in factory
	gqlServer := gql.NewGraphQLServer(h.facade.GetHub())

	var registerGql = func() {
		h.appendEndpointHandler(http.MethodPost, gqlQueryEndpoint, h.gqlHandler(gqlServer))
		h.appendEndpointHandler(http.MethodPost, gqlSubscriptionEndpoint, h.gqlHandler(gqlServer))
	}

	var registerWs = func() {
		h.appendEndpointHandler(http.MethodGet, websocketEndpoint, h.wsHandler)
	}

	switch h.facade.GetDispatchType() {
	case dispatchAll:
		registerWs()
		registerGql()
	case dispatchGql:
		registerGql()
	case dispatchWs:
		registerWs()
	default:
		return nil, common.ErrInvalidDispatchType
	}

	return h.endpoints, nil
}

func (h *hubGroup) wsHandler(c *gin.Context) {
	ws.Serve(h.facade.GetHub(), c.Writer, c.Request)
}

func (h *hubGroup) gqlHandler(gqlServer *gqlHandler.Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		gqlServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *hubGroup) appendEndpointHandler(method, path string, handler gin.HandlerFunc) {
	h.endpoints = append(h.endpoints, &shared.EndpointHandlerData{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *hubGroup) IsInterfaceNil() bool {
	return h == nil
}
