package groups

import (
	"strings"

	"github.com/gin-gonic/gin"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/config"
)

var log = logger.GetOrCreate("api/groups")

type baseGroup struct {
	endpoints             []*shared.EndpointHandlerData
	additionalMiddlewares []gin.HandlerFunc
}

// RegisterRoutes will register all the endpoints to the given web server
func (bg *baseGroup) RegisterRoutes(
	ws *gin.RouterGroup,
	apiConfig config.APIRoutesConfig,
) {
	for _, handlerData := range bg.endpoints {
		isOpen, isAuthEnabled := getEndpointStatus(ws, handlerData.Path, apiConfig)
		if !isOpen {
			log.Debug("endpoint is closed", "path", handlerData.Path)
			continue
		}

		if isAuthEnabled {
			ws.Use(bg.GetAdditionalMiddlewares()...)
		}

		ws.Handle(handlerData.Method, handlerData.Path, handlerData.Handler)
	}
}

// GetAdditionalMiddlewares return additional middlewares
func (bg *baseGroup) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return bg.additionalMiddlewares
}

func getEndpointStatus(
	ws *gin.RouterGroup,
	path string,
	apiConfig config.APIRoutesConfig,
) (bool, bool) {
	basePath := ws.BasePath()

	// ws.BasePath will return paths like /group
	// so we need the last token after splitting by /
	splitPath := strings.Split(basePath, "/")
	basePath = splitPath[len(splitPath)-1]

	group, ok := apiConfig.APIPackages[basePath]
	if !ok {
		return false, false
	}

	for _, route := range group.Routes {
		if route.Name == path {
			return route.Open, route.Auth
		}
	}

	return false, false
}
