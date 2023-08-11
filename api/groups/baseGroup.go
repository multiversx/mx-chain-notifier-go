package groups

import (
	"github.com/gin-gonic/gin"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
)

var log = logger.GetOrCreate("api/groups")

type baseGroup struct {
	endpoints []*shared.EndpointHandlerData
}

// RegisterRoutes will register all the endpoints to the given web server
func (bg *baseGroup) RegisterRoutes(
	ws gin.IRoutes,
) {
	for _, handlerData := range bg.endpoints {
		ws.Handle(handlerData.Method, handlerData.Path, handlerData.Handler)
	}
}
