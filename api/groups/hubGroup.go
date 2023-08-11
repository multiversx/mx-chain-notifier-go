package groups

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
)

const (
	websocketEndpoint = "/ws"
)

type hubGroup struct {
	*baseGroup
	facade HubFacadeHandler
}

// NewHubGroup registers handlers for the /hub group
// It only registers the specified hub implementation and its corresponding dispatchers
func NewHubGroup(facade HubFacadeHandler) (*hubGroup, error) {
	if check.IfNil(facade) {
		return nil, fmt.Errorf("%w for hub group", errors.ErrNilFacadeHandler)
	}

	h := &hubGroup{
		facade: facade,
		baseGroup: &baseGroup{
			additionalMiddlewares: make([]gin.HandlerFunc, 0),
		},
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

func (h *hubGroup) wsHandler(c *gin.Context) {
	h.facade.ServeHTTP(c.Writer, c.Request)
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *hubGroup) IsInterfaceNil() bool {
	return h == nil
}
