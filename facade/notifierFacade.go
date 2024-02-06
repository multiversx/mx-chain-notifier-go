package facade

import (
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

var log = logger.GetOrCreate("facade")

// ArgsNotifierFacade defines the arguments necessary for notifierFacade creation
type ArgsNotifierFacade struct {
	APIConfig            config.ConnectorApiConfig
	EventsHandler        EventsHandler
	WSHandler            dispatcher.WSHandler
	StatusMetricsHandler common.StatusMetricsHandler
}

type notifierFacade struct {
	config        config.ConnectorApiConfig
	eventsHandler EventsHandler
	wsHandler     dispatcher.WSHandler
	statusMetrics common.StatusMetricsHandler
}

// NewNotifierFacade creates a new notifier facade instance
func NewNotifierFacade(args ArgsNotifierFacade) (*notifierFacade, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &notifierFacade{
		eventsHandler: args.EventsHandler,
		config:        args.APIConfig,
		wsHandler:     args.WSHandler,
		statusMetrics: args.StatusMetricsHandler,
	}, nil
}

func checkArgs(args ArgsNotifierFacade) error {
	if check.IfNil(args.EventsHandler) {
		return ErrNilEventsHandler
	}
	if check.IfNil(args.WSHandler) {
		return ErrNilWSHandler
	}
	if check.IfNil(args.StatusMetricsHandler) {
		return common.ErrNilStatusMetricsHandler
	}

	return nil
}

// HandlePushEvents will handle push events received from observer
// It splits block data and handles log, txs and srcs events separately
func (nf *notifierFacade) HandlePushEvents(allEvents data.ArgsSaveBlockData) error {
	return nf.eventsHandler.HandleSaveBlockEvents(allEvents)
}

// HandleRevertEvents will handle revents events received from observer
func (nf *notifierFacade) HandleRevertEvents(events data.RevertBlock) {
	nf.eventsHandler.HandleRevertEvents(events)
}

// HandleFinalizedEvents will handle finalized events received from observer
func (nf *notifierFacade) HandleFinalizedEvents(events data.FinalizedBlock) {
	nf.eventsHandler.HandleFinalizedEvents(events)
}

// ServeHTTP will handle a websocket request
func (nf *notifierFacade) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nf.wsHandler.ServeHTTP(w, r)
}

// GetConnectorUserAndPass will return username and password (for basic authentication)
// from config
func (nf *notifierFacade) GetConnectorUserAndPass() (string, string) {
	return nf.config.Username, nf.config.Password
}

// GetMetrics will return metrics in json format
func (nf *notifierFacade) GetMetrics() map[string]*data.EndpointMetricsResponse {
	return nf.statusMetrics.GetAll()
}

// GetMetricsForPrometheus will return metrics in prometheus format
func (nf *notifierFacade) GetMetricsForPrometheus() string {
	return nf.statusMetrics.GetMetricsForPrometheus()
}

// IsInterfaceNil returns true if there is no value under the interface
func (nf *notifierFacade) IsInterfaceNil() bool {
	return nf == nil
}
