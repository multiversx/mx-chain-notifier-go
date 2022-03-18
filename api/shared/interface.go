package shared

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/gin-gonic/gin"
)

// HTTPServerCloser defines the basic actions of starting and closing that a web server should be able to do
type HTTPServerCloser interface {
	Start()
	Close() error
	IsInterfaceNil() bool
}

// GroupHandler defines the actions needed to be performed by an gin API group
type GroupHandler interface {
	RegisterRoutes(ws gin.IRoutes)
	GetAdditionalMiddlewares() []gin.HandlerFunc
	IsInterfaceNil() bool
}

// FacadeHandler defines the behavior of a notifier base facade handler
type FacadeHandler interface {
	HandlePushEvents(events data.BlockEvents)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	GetHub() dispatcher.Hub
	GetDispatchType() string
	IsInterfaceNil() bool
}

// HTTPServerHandler defines the behavior of a web server
type HTTPServerHandler interface {
	Run() error
	Close() error
	AddGroup(groupName string, group GroupHandler)
	IsInterfaceNil() bool
}
