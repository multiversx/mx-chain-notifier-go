package handlers

import "github.com/gin-gonic/gin"

type EndpointHandler struct {
	Path        string
	Method      string
	HandlerFunc gin.HandlerFunc
}

type groupHandler struct {
	endpointHandlersMap map[string][]EndpointHandler
}

func NewGroupHandler() *groupHandler {
	return &groupHandler{
		endpointHandlersMap: make(map[string][]EndpointHandler),
	}
}

func (g *groupHandler) RegisterEndpoints(r *gin.Engine) {
	for route, handlersGroup := range g.endpointHandlersMap {
		routerGroup := r.Group(route)
		{
			for _, h := range handlersGroup {
				routerGroup.Handle(h.Method, h.Path, h.HandlerFunc)
			}
		}
	}
}

func (g *groupHandler) AddEndpointHandlers(route string, handlers []EndpointHandler) {
	g.endpointHandlersMap[route] = handlers
}

type apiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func JsonResponse(c *gin.Context, status int, data interface{}, error string) {
	c.JSON(status, apiResponse{
		Data:  data,
		Error: error,
	})
}
