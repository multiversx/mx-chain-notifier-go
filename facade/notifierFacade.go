package facade

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// ArgsNotifierFacade defines the arguments necessary for notifierFacade creation
type ArgsNotifierFacade struct {
	APIConfig     config.ConnectorApiConfig
	EventsHandler EventsHandler
	WSHandler     dispatcher.WSHandler
}

type notifierFacade struct {
	config        config.ConnectorApiConfig
	eventsHandler EventsHandler
	wsHandler     dispatcher.WSHandler
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
	}, nil
}

func checkArgs(args ArgsNotifierFacade) error {
	if check.IfNil(args.EventsHandler) {
		return ErrNilEventsHandler
	}
	if check.IfNil(args.WSHandler) {
		return ErrNilWSHandler
	}

	return nil
}

// HandlePushEvents will handle push events received from observer
// It splits block data and handles log, txs and srcs events separately
func (nf *notifierFacade) HandlePushEvents(allEvents data.ArgsSaveBlockData) {
	hash := string(allEvents.HeaderHash)
	pushEvents := data.BlockEvents{
		Hash:   hash,
		Events: nil,
	}
	nf.eventsHandler.HandlePushEvents(pushEvents)

	txs := data.BlockTxs{
		Hash: hash,
		Txs:  nil,
	}
	nf.eventsHandler.HandleBlockTxs(txs)

	scrs := data.BlockScrs{
		Hash: hash,
		Scrs: nil,
	}
	nf.eventsHandler.HandleBlockScrs(scrs)
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
