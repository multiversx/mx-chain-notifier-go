package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-notifier-go/data"
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
	HandlePushEventsV2(events data.ArgsSaveBlockData) error
	HandlePushEventsV1(events data.SaveBlockData) error
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetMetrics() map[string]*data.EndpointMetricsResponse
	GetMetricsForPrometheus() string
	IsInterfaceNil() bool
}

// WebServerHandler defines the behavior of a web server
type WebServerHandler interface {
	Run() error
	Close() error
	IsInterfaceNil() bool
}
