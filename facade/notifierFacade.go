package facade

import (
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/data"
)

// TODO: comments update

type ArgsNotifierFacade struct {
	EventsHandler EventsHandler
}

type notifierFacade struct {
	eventsHandler EventsHandler
}

func NewNotifierFacade(args ArgsNotifierFacade) (*notifierFacade, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &notifierFacade{
		eventsHandler: args.EventsHandler,
	}, nil
}

func checkArgs(args ArgsNotifierFacade) error {
	if check.IfNil(args.EventsHandler) {
		return ErrNilEventsHandler
	}

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
