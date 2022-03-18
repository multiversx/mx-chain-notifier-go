package facade

import (
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// TODO: comments update

type ArgsNotifierFacade struct {
	EventsHandler EventsHandler
	APIConfig     config.ConnectorApiConfig
	Hub           dispatcher.Hub
}

type notifierFacade struct {
	eventsHandler EventsHandler
	config        config.ConnectorApiConfig
	hub           dispatcher.Hub
}

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
	// TODO: more checks

	return nil
}

func (nf *notifierFacade) HandlePushEvents(events data.BlockEvents) {
	nf.eventsHandler.HandlePushEvents(events)
}

func (nf *notifierFacade) HandleRevertEvents(events data.RevertBlock) {
	nf.eventsHandler.HandleRevertEvents(events)
}

func (nf *notifierFacade) HandleFinalizedEvents(events data.FinalizedBlock) {
	nf.eventsHandler.HandleFinalizedEvents(events)
}

func (nf *notifierFacade) GetDispatchType() string {
	return nf.config.DispatchType
}

func (nf *notifierFacade) GetHub() dispatcher.Hub {
	return nf.hub
}

func (nf *notifierFacade) GetConnectorUserAndPass() (string, string) {
	//return nf.config.Username, nf.config.Password
	return "", ""
}

// IsInterfaceNil returns true if there is no value under the interface
func (nf *notifierFacade) IsInterfaceNil() bool {
	return nf == nil
}
