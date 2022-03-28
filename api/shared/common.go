package shared

import "github.com/gin-gonic/gin"

// EndpointHandlerData holds the items needed for creating a new gin HTTP endpoint
type EndpointHandlerData struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

type apiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// JSONResponse is a wrapper for gin.Context JSON payload
func JSONResponse(c *gin.Context, status int, data interface{}, error string) {
	c.JSON(status, apiResponse{
		Data:  data,
		Error: error,
	})
}
