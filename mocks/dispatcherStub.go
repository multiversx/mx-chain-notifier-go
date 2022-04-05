package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

// DispatcherStub implements dispatcher EventDispatcher interface
type DispatcherStub struct {
	GetIDCalled          func() uuid.UUID
	PushEventsCalled     func(events []data.Event)
	RevertEventCalled    func(event data.RevertBlock)
	FinalizedEventCalled func(event data.FinalizedBlock)
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
