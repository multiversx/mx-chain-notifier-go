package mocks

import "github.com/ElrondNetwork/notifier-go/data"

// EventsHandlerStub implements EventsHandler interface
type EventsHandlerStub struct {
	HandlePushEventsCalled      func(events data.BlockEvents)
	HandleRevertEventsCalled    func(revertBlock data.RevertBlock)
	HandleFinalizedEventsCalled func(finalizedBlock data.FinalizedBlock)
}

// HandlePushEvents -
func (e *EventsHandlerStub) HandlePushEvents(events data.BlockEvents) {
	if e.HandlePushEventsCalled != nil {
		e.HandlePushEventsCalled(events)
	}
}

// HandleRevertEvents -
func (e *EventsHandlerStub) HandleRevertEvents(revertBlock data.RevertBlock) {
	if e.HandleRevertEventsCalled != nil {
		e.HandleRevertEventsCalled(revertBlock)
	}
}

// HandleFinalizedEvents -
func (e *EventsHandlerStub) HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) {
	if e.HandleFinalizedEventsCalled != nil {
		e.HandleFinalizedEventsCalled(finalizedBlock)
	}
}

// IsInterfaceNil -
func (e *EventsHandlerStub) IsInterfaceNil() bool {
	return e == nil
}
