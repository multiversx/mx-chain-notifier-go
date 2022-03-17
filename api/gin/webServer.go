package gin

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	apiErrors "github.com/ElrondNetwork/notifier-go/api/errors"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var log = logger.GetOrCreate("api/gin")

// ArgsWebServerHandler holds the arguments needed to create a web server handler
type ArgsWebServerHandler struct {
	Facade shared.FacadeHandler
	Config *config.GeneralConfig
}

// webServer is a wrapper for gin.Engine, holding additional components
type webServer struct {
	sync.RWMutex
	facade     shared.FacadeHandler
	httpServer shared.HTTPServerCloser
	config     *config.GeneralConfig
	groups     map[string]shared.GroupHandler
	cancelFunc func()
}

// NewWebServerHandler creates and configures an instance of webServer
func NewWebServerHandler(args ArgsWebServerHandler) (*webServer, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &webServer{
		facade: args.Facade,
		config: args.Config,
		groups: make(map[string]shared.GroupHandler),
	}, nil
}

func checkArgs(args ArgsWebServerHandler) error {
	if check.IfNil(args.Facade) {
		return apiErrors.ErrNilFacadeHandler
	}
	// TODO: check configs

	return nil
}

// Run starts the server and the Hub as goroutines
// It returns an instance of http.Server
func (w *webServer) Run() error {
	w.Lock()
	defer w.Unlock()

	var err error

	port := w.config.ConnectorApi.Port
	if !strings.Contains(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	engine := gin.Default()
	engine.Use(cors.Default())

	err = w.createGroups()
	if err != nil {
		return err
	}

	w.registerRoutes(engine)

	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	w.httpServer, err = NewHTTPServer(server)
	if err != nil {
		return err
	}

	go w.httpServer.Start()

	return nil
}

func (w *webServer) createGroups() error {
	// groupsMap := make(map[string]shared.GroupHandler)
	// vmValuesGroup, err := groups.NewEventsHandler()
	// if err != nil {
	// 	return err
	// }
	// groupsMap["vm-values"] = vmValuesGroup

	// w.groups = groupsMap

	return nil
}

func (w *webServer) registerRoutes(ginEngine *gin.Engine) {
	for groupName, groupHandler := range w.groups {
		log.Debug("registering API group", "group name", groupName)
		ginGroup := ginEngine.Group(fmt.Sprintf("/%s", groupName))
		groupHandler.RegisterRoutes(ginGroup)
	}
}

// TODO: remove this after further refactoring
func (w *webServer) AddGroup(groupName string, group shared.GroupHandler) {
	w.Lock()
	defer w.Unlock()

	w.groups[groupName] = group
}

// Close will handle the closing of inner components
func (w *webServer) Close() error {
	if w.cancelFunc != nil {
		w.cancelFunc()
	}

	w.Lock()
	err := w.httpServer.Close()
	w.Unlock()

	if err != nil {
		err = fmt.Errorf("%w while closing the http server in gin/webServer", err)
	}

	return err
}

// IsInterfaceNil returns true if there is no value under the interface
func (w *webServer) IsInterfaceNil() bool {
	return w == nil
}
