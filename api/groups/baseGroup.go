package groups

import (
	"fmt"
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
	authMiddleware        gin.HandlerFunc
	hasAuthMiddleware     bool
}

func newBaseGroup() *baseGroup {
	return &baseGroup{
		additionalMiddlewares: make([]gin.HandlerFunc, 0),
		authMiddleware:        func(ctx *gin.Context) {},
		hasAuthMiddleware:     false,
	}
}

// setAuthMiddleware installs a real auth middleware, marking the group as
// having one so RegisterRoutes can detect (and refuse) an Auth = true route
// that would otherwise fall back to the no-op default middleware.
func (bg *baseGroup) setAuthMiddleware(middleware gin.HandlerFunc) {
	bg.authMiddleware = middleware
	bg.hasAuthMiddleware = true
}

// RegisterRoutes will register all the endpoints to the given web server.
// It fails instead of registering a route whose config requires
// authentication if no real auth middleware was ever configured for the
// group - otherwise the route would silently be served with the no-op
// default middleware, i.e. with no authentication at all.
func (bg *baseGroup) RegisterRoutes(
	ws *gin.RouterGroup,
	apiConfig config.APIRoutesConfig,
) error {
	for _, handlerData := range bg.endpoints {
		isOpen, isAuthEnabled := getEndpointStatus(ws, handlerData.Path, apiConfig)
		if !isOpen {
			log.Debug("endpoint is closed", "path", handlerData.Path)
			continue
		}

		handlers := make([]gin.HandlerFunc, 0)

		if isAuthEnabled {
			if !bg.hasAuthMiddleware {
				return fmt.Errorf("%w: path %s", ErrAuthEnabledWithoutMiddleware, handlerData.Path)
			}
			handlers = append(handlers, bg.GetAuthMiddleware())
		}

		handlers = append(handlers, bg.GetAdditionalMiddlewares()...)
		handlers = append(handlers, handlerData.Handler)

		ws.Handle(handlerData.Method, handlerData.Path, handlers...)
	}

	return nil
}

// GetAdditionalMiddlewares returns additional middlewares
func (bg *baseGroup) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return bg.additionalMiddlewares
}

// GetAuthMiddleware returns auth middleware
func (bg *baseGroup) GetAuthMiddleware() gin.HandlerFunc {
	return bg.authMiddleware
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
