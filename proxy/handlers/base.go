package handlers

import "github.com/gin-gonic/gin"

// EndpointGroupHandler defines the components of a gin.Engine Group
// It allows for multiple middleware configuration
type EndpointGroupHandler struct {
	Root             string
	Middlewares      []gin.HandlerFunc
	EndpointHandlers []EndpointHandler
}

// EndpointHandler defines the needed components for registering an endpoint
type EndpointHandler struct {
	Path        string
	Method      string
	HandlerFunc gin.HandlerFunc
}

type groupHandler struct {
	endpointHandlersMap map[string]EndpointGroupHandler
}

// NewGroupHandler creates an instance of a groupHandler
func NewGroupHandler() *groupHandler {
	return &groupHandler{
		endpointHandlersMap: make(map[string]EndpointGroupHandler),
	}
}

// RegisterEndpoints registers the endpoints groups and the corresponding handlers
// It should be called after all the handlers have been defined
func (g *groupHandler) RegisterEndpoints(r *gin.Engine) {
	for groupRoot, handlersGroup := range g.endpointHandlersMap {
		routerGroup := r.Group(groupRoot).Use(handlersGroup.Middlewares...)
		{
			for _, h := range handlersGroup.EndpointHandlers {
				routerGroup.Handle(h.Method, h.Path, h.HandlerFunc)
			}
		}
	}
}

// AddEndpointGroupHandler inserts an EndpointGroupHandler instance to the map
// The key of the endpointHandlersMap is the base path of the group
// The method is not thread-safe and does not validate inputs
func (g *groupHandler) AddEndpointGroupHandler(endpointHandler EndpointGroupHandler) {
	g.endpointHandlersMap[endpointHandler.Root] = endpointHandler
}

type apiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// JsonResponse is a wrapper for gin.Context JSON payload
func JsonResponse(c *gin.Context, status int, data interface{}, error string) {
	c.JSON(status, apiResponse{
		Data:  data,
		Error: error,
	})
}
