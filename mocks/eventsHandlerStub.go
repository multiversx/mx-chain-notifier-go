package mocks

import "github.com/multiversx/mx-chain-notifier-go/data"

// EventsHandlerStub implements EventsHandler interface
type EventsHandlerStub struct {
	HandleSaveBlockEventsCalled func(allEvents data.ArgsSaveBlockData) error
	HandleRevertEventsCalled    func(revertBlock data.RevertBlock)
	HandleFinalizedEventsCalled func(finalizedBlock data.FinalizedBlock)
}

// HandleSaveBlockEvents -
func (e *EventsHandlerStub) HandleSaveBlockEvents(events data.ArgsSaveBlockData) error {
	if e.HandleSaveBlockEventsCalled != nil {
		return e.HandleSaveBlockEventsCalled(events)
	}

	return nil
}

// HandleRevertEvents -
func (e *EventsHandlerStub) HandleRevertEvents(revertBlock data.RevertBlock) error {
	if e.HandleRevertEventsCalled != nil {
		e.HandleRevertEventsCalled(revertBlock)
	}

	return nil
}

// HandleFinalizedEvents -
func (e *EventsHandlerStub) HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) error {
	if e.HandleFinalizedEventsCalled != nil {
		e.HandleFinalizedEventsCalled(finalizedBlock)
	}

	return nil
}

// IsInterfaceNil -
func (e *EventsHandlerStub) IsInterfaceNil() bool {
	return e == nil
}
