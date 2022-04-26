package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
)

// DispatcherMock -
type DispatcherMock struct {
	id       uuid.UUID
	consumer *ConsumerMock
	hub      dispatcher.Hub
}

// NewDispatcherMock -
func NewDispatcherMock(consumer *ConsumerMock, hub dispatcher.Hub) *DispatcherMock {
	return &DispatcherMock{
		id:       uuid.New(),
		consumer: consumer,
		hub:      hub,
	}
}

// GetID -
func (d *DispatcherMock) GetID() uuid.UUID {
	return d.id
}

// PushEvents -
func (d *DispatcherMock) PushEvents(events []data.Event) {
	d.consumer.Receive(events)
}

// RevertEvent -
func (d *DispatcherMock) RevertEvent(event data.RevertBlock) {
}

// FinalizedEvent -
func (d *DispatcherMock) FinalizedEvent(event data.FinalizedBlock) {
}

// TxsEvent -
func (d *DispatcherMock) TxsEvent(event data.BlockTxs) {
}

// ScrsEvent -
func (d *DispatcherMock) ScrsEvent(event data.BlockScrs) {
}

// Subscribe -
func (d *DispatcherMock) Subscribe(event data.SubscribeEvent) {
	d.hub.Subscribe(event)
}

// Register -
func (d *DispatcherMock) Register() {
	d.hub.RegisterEvent(d)
}

// Unregister -
func (d *DispatcherMock) Unregister() {
	d.hub.UnregisterEvent(d)
}
