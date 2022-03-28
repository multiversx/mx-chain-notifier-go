package facade

import (
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// ArgsNotifierFacade defines the arguments necessary for notifierFacade creation
type ArgsNotifierFacade struct {
	APIConfig     config.ConnectorApiConfig
	EventsHandler EventsHandler
	Hub           HubHandler
}

type notifierFacade struct {
	config        config.ConnectorApiConfig
	eventsHandler EventsHandler
	hub           HubHandler
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
		hub:           args.Hub,
	}, nil
}

func checkArgs(args ArgsNotifierFacade) error {
	if check.IfNil(args.EventsHandler) {
		return ErrNilEventsHandler
	}
	if check.IfNil(args.Hub) {
		return ErrNilHubHandler
	}

	return nil
}

// HandlePushEvents will handle push events received from observer
func (nf *notifierFacade) HandlePushEvents(events data.BlockEvents) {
	nf.eventsHandler.HandlePushEvents(events)
}

// HandleRevertEvents will handle revents events received from observer
func (nf *notifierFacade) HandleRevertEvents(events data.RevertBlock) {
	nf.eventsHandler.HandleRevertEvents(events)
}

// HandleFinalizedEvents will handle finalized events received from observer
func (nf *notifierFacade) HandleFinalizedEvents(events data.FinalizedBlock) {
	nf.eventsHandler.HandleFinalizedEvents(events)
}

// GetDispatchType will provide dispatch type based on config
func (nf *notifierFacade) GetDispatchType() string {
	return nf.config.DispatchType
}

// GetHub will return the hub component
func (nf *notifierFacade) GetHub() dispatcher.Hub {
	return nf.hub
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
