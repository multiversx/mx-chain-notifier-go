package groups

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/gin-gonic/gin"
)

var log = logger.GetOrCreate("api/groups")

type baseGroup struct {
	endpoints []*shared.EndpointHandlerData
}

// GetEndpoints returns all the endpoints specific to the group
func (bg *baseGroup) GetEndpoints() []*shared.EndpointHandlerData {
	return bg.endpoints
}

// RegisterRoutes will register all the endpoints to the given web server
func (bg *baseGroup) RegisterRoutes(
	ws gin.IRoutes,
) {
	for _, handlerData := range bg.endpoints {
		ws.Handle(handlerData.Method, handlerData.Path, handlerData.Handler)
	}
}

// AddEndpointGroupHandler inserts an EndpointGroupHandler instance to the map
// The key of the endpointHandlersMap is the base path of the group
// The method is not thread-safe and does not validate inputs
func (bg *baseGroup) AddEndpointGroupHandler(endpointHandler *shared.EndpointHandlerData) {
	bg.endpoints = append(bg.endpoints, endpointHandler)
}
