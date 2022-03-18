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
