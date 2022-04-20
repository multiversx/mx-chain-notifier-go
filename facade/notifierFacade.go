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
func (nf *notifierFacade) HandlePushEvents(allEvents data.SaveBlockData) {
	pushEvents := data.BlockEvents{
		Hash:   allEvents.Hash,
		Events: allEvents.LogEvents,
	}
	nf.eventsHandler.HandlePushEvents(pushEvents)

	txs := data.BlockTxs{
		Hash: allEvents.Hash,
		Txs:  allEvents.Txs,
	}
	nf.eventsHandler.HandleTxsEvents(txs)

	scrs := data.BlockScrs{
		Hash: allEvents.Hash,
		Scrs: allEvents.Scrs,
	}
	nf.eventsHandler.HandleScrsEvents(scrs)
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
