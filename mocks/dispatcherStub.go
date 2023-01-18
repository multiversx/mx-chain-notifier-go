package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

// DispatcherStub implements dispatcher EventDispatcher interface
type DispatcherStub struct {
	GetIDCalled             func() uuid.UUID
	PushEventsCalled        func(events []data.Event)
	BlockEventsCalled       func(event data.BlockEvents)
	RevertEventCalled       func(event data.RevertBlock)
	FinalizedEventCalled    func(event data.FinalizedBlock)
	TxsEventCalled          func(event data.BlockTxs)
	TxsWithOrderEventCalled func(event data.BlockTxsWithOrder)
	ScrsEventCalled         func(event data.BlockScrs)
}

// GetID -
func (d *DispatcherStub) GetID() uuid.UUID {
	if d.GetIDCalled != nil {
		return d.GetIDCalled()
	}

	return uuid.UUID{}
}

// PushEvents -
func (d *DispatcherStub) PushEvents(events []data.Event) {
	if d.PushEventsCalled != nil {
		d.PushEventsCalled(events)
	}
}

// BlockEvents -
func (d *DispatcherStub) BlockEvents(events data.BlockEvents) {
	if d.BlockEventsCalled != nil {
		d.BlockEventsCalled(events)
	}
}

// RevertEvent -
func (d *DispatcherStub) RevertEvent(event data.RevertBlock) {
	if d.RevertEventCalled != nil {
		d.RevertEventCalled(event)
	}
}

// FinalizedEvent -
func (d *DispatcherStub) FinalizedEvent(event data.FinalizedBlock) {
	if d.FinalizedEventCalled != nil {
		d.FinalizedEventCalled(event)
	}
}

// TxsEvent -
func (d *DispatcherStub) TxsEvent(event data.BlockTxs) {
	if d.TxsEventCalled != nil {
		d.TxsEventCalled(event)
	}
}

// TxsWithOrderEvent -
func (d *DispatcherStub) TxsWithOrderEvent(event data.BlockTxsWithOrder) {
	if d.TxsWithOrderEventCalled != nil {
		d.TxsWithOrderEventCalled(event)
	}
}

// ScrsEvent -
func (d *DispatcherStub) ScrsEvent(event data.BlockScrs) {
	if d.ScrsEventCalled != nil {
		d.ScrsEventCalled(event)
	}
}
