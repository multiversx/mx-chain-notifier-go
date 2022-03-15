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

// Subscribe -
func (d *DispatcherMock) Subscribe(event dispatcher.SubscribeEvent) {
	d.hub.Subscribe(event)
}

// Register -
func (d *DispatcherMock) Register() {
	d.hub.RegisterChan() <- d
}

// Unregister -
func (d *DispatcherMock) Unregister() {
	d.hub.UnregisterChan() <- d
}
