package facade

import (
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

var log = logger.GetOrCreate("facade")

// ArgsNotifierFacade defines the arguments necessary for notifierFacade creation
type ArgsNotifierFacade struct {
	APIConfig         config.ConnectorApiConfig
	EventsHandler     EventsHandler
	WSHandler         dispatcher.WSHandler
	EventsInterceptor EventsInterceptor
}

type notifierFacade struct {
	config            config.ConnectorApiConfig
	eventsHandler     EventsHandler
	wsHandler         dispatcher.WSHandler
	eventsInterceptor EventsInterceptor
}

// NewNotifierFacade creates a new notifier facade instance
func NewNotifierFacade(args ArgsNotifierFacade) (*notifierFacade, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &notifierFacade{
		eventsHandler:     args.EventsHandler,
		config:            args.APIConfig,
		wsHandler:         args.WSHandler,
		eventsInterceptor: args.EventsInterceptor,
	}, nil
}

func checkArgs(args ArgsNotifierFacade) error {
	if check.IfNil(args.EventsHandler) {
		return ErrNilEventsHandler
	}
	if check.IfNil(args.WSHandler) {
		return ErrNilWSHandler
	}
	if check.IfNil(args.EventsInterceptor) {
		return ErrNilEventsInterceptor
	}

	return nil
}

// HandlePushEventsV2 will handle push events received from observer
// It splits block data and handles log, txs and srcs events separately
func (nf *notifierFacade) HandlePushEventsV2(allEvents data.ArgsSaveBlockData) error {
	eventsData, err := nf.eventsInterceptor.ProcessBlockEvents(&allEvents)
	if err != nil {
		return err
	}

	pushEvents := data.BlockEvents{
		Hash:   eventsData.Hash,
		Events: eventsData.LogEvents,
	}
	err = nf.eventsHandler.HandlePushEvents(pushEvents)
	if err != nil {
		return err
	}

	txs := data.BlockTxs{
		Hash: eventsData.Hash,
		Txs:  eventsData.Txs,
	}
	nf.eventsHandler.HandleBlockTxs(txs)

	scrs := data.BlockScrs{
		Hash: eventsData.Hash,
		Scrs: eventsData.Scrs,
	}
	nf.eventsHandler.HandleBlockScrs(scrs)

	return nil
}

// HandlePushEventsV1 will handle push events received from observer
// It splits block data and handles log, txs and srcs events separately
// TODO: remove this implementation
func (nf *notifierFacade) HandlePushEventsV1(eventsData data.SaveBlockData) error {
	pushEvents := data.BlockEvents{
		Hash:   eventsData.Hash,
		Events: eventsData.LogEvents,
	}
	err := nf.eventsHandler.HandlePushEvents(pushEvents)
	if err != nil {
		return err
	}

	txs := data.BlockTxs{
		Hash: eventsData.Hash,
		Txs:  eventsData.Txs,
	}
	nf.eventsHandler.HandleBlockTxs(txs)

	scrs := data.BlockScrs{
		Hash: eventsData.Hash,
		Scrs: eventsData.Scrs,
	}
	nf.eventsHandler.HandleBlockScrs(scrs)

	return nil
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

// IsInterfaceNil returns true if there is no value under the interface
func (nf *notifierFacade) IsInterfaceNil() bool {
	return nf == nil
}
