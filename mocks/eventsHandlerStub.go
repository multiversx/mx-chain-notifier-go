package mocks

import "github.com/ElrondNetwork/notifier-go/data"

// EventsHandlerStub implements EventsHandler interface
type EventsHandlerStub struct {
	HandlePushEventsCalled        func(events data.BlockEvents) error
	HandleRevertEventsCalled      func(revertBlock data.RevertBlock)
	HandleFinalizedEventsCalled   func(finalizedBlock data.FinalizedBlock)
	HandleBlockTxsCalled          func(blockTxs data.BlockTxs)
	HandleBlockScrsCalled         func(blockScrs data.BlockScrs)
	HandleBlockTxsWithOrderCalled func(blockTxs data.BlockTxsWithOrder)
}

// HandlePushEvents -
func (e *EventsHandlerStub) HandlePushEvents(events data.BlockEvents) error {
	if e.HandlePushEventsCalled != nil {
		return e.HandlePushEventsCalled(events)
	}

	return nil
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

// HandleBlockTxs -
func (e *EventsHandlerStub) HandleBlockTxs(blockTxs data.BlockTxs) {
	if e.HandleBlockTxsCalled != nil {
		e.HandleBlockTxsCalled(blockTxs)
	}
}

// HandleBlockScrs -
func (e *EventsHandlerStub) HandleBlockScrs(blockScrs data.BlockScrs) {
	if e.HandleBlockScrsCalled != nil {
		e.HandleBlockScrsCalled(blockScrs)
	}
}

// HandleBlockTxsWithOrder -
func (e *EventsHandlerStub) HandleBlockTxsWithOrder(blockTxs data.BlockTxsWithOrder) {
	if e.HandleBlockTxsWithOrderCalled != nil {
		e.HandleBlockTxsWithOrderCalled(blockTxs)
	}
}

// IsInterfaceNil -
func (e *EventsHandlerStub) IsInterfaceNil() bool {
	return e == nil
}
