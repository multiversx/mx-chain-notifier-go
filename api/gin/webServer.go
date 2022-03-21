package gin

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	apiErrors "github.com/ElrondNetwork/notifier-go/api/errors"
	"github.com/ElrondNetwork/notifier-go/api/groups"
	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var log = logger.GetOrCreate("api/gin")

// ArgsWebServerHandler holds the arguments needed to create a web server handler
type ArgsWebServerHandler struct {
	Facade shared.FacadeHandler
	Config config.ConnectorApiConfig
	Type   string
}

// webServer is a wrapper for gin.Engine, holding additional components
type webServer struct {
	sync.RWMutex
	facade     shared.FacadeHandler
	httpServer shared.HTTPServerCloser
	config     config.ConnectorApiConfig
	groups     map[string]shared.GroupHandler
	apiType    string
	cancelFunc func()
}

// NewWebServerHandler creates and configures an instance of webServer
func NewWebServerHandler(args ArgsWebServerHandler) (*webServer, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &webServer{
		facade:  args.Facade,
		config:  args.Config,
		apiType: args.Type,
		groups:  make(map[string]shared.GroupHandler),
	}, nil
}

func checkArgs(args ArgsWebServerHandler) error {
	if check.IfNil(args.Facade) {
		return apiErrors.ErrNilFacadeHandler
	}
	if args.Type == "" {
		return common.ErrInvalidAPIType
	}

	return nil
}

// Run starts the server and the Hub as goroutines
// It returns an instance of http.Server
func (w *webServer) Run() error {
	w.Lock()
	defer w.Unlock()

	var err error

	port := w.config.Port
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
	groupsMap := make(map[string]shared.GroupHandler)

	eventsGroup, err := groups.NewEventsGroup(w.facade)
	if err != nil {
		return err
	}
	groupsMap["events"] = eventsGroup

	if w.apiType == common.WSAPIType {
		hubHandler, err := groups.NewHubGroup(w.facade)
		if err != nil {
			return err
		}
		groupsMap["hub"] = hubHandler
	}

	w.groups = groupsMap

	return nil
}

func (w *webServer) registerRoutes(ginEngine *gin.Engine) {
	for groupName, groupHandler := range w.groups {
		log.Info("registering API group", "group name", groupName)
		ginGroup := ginEngine.Group(fmt.Sprintf("/%s", groupName)).Use(groupHandler.GetAdditionalMiddlewares()...)
		groupHandler.RegisterRoutes(ginGroup)
	}
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
