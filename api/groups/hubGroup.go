package groups

import (
	"fmt"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/gin-gonic/gin"
)

const (
	websocketEndpoint = "/ws"
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

	endpoints := []*shared.EndpointHandlerData{
		{
			Method:  http.MethodGet,
			Path:    websocketEndpoint,
			Handler: h.wsHandler,
		},
	}

	h.endpoints = endpoints

	return h, nil
}

// GetAdditionalMiddlewares return additional middlewares for this group
func (h *hubGroup) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return h.additionalMiddlewares
}

func (h *hubGroup) wsHandler(c *gin.Context) {
	h.facade.ServeHTTP(c.Writer, c.Request)
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *hubGroup) IsInterfaceNil() bool {
	return h == nil
}
