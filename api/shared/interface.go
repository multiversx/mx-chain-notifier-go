package shared

import (
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
	RegisterRoutes(ws *gin.RouterGroup)
	IsInterfaceNil() bool
}

// FacadeHandler
type FacadeHandler interface {
	IsInterfaceNil() bool
}

// HTTPServerHandler defines the behavior of a web server
type HTTPServerHandler interface {
	Run() error
	Close() error
	AddGroup(groupName string, group GroupHandler)
	IsInterfaceNil() bool
}
