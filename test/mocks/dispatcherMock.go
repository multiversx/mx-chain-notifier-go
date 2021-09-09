package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
)

type DispatcherMock struct {
	id       uuid.UUID
	consumer *ConsumerMock
	hub      dispatcher.Hub
}

func NewDispatcherMock(consumer *ConsumerMock, hub dispatcher.Hub) *DispatcherMock {
	return &DispatcherMock{
		id:       uuid.New(),
		consumer: consumer,
		hub:      hub,
	}
}

func (d *DispatcherMock) GetID() uuid.UUID {
	return d.id
}

func (d *DispatcherMock) PushEvents(events []data.Event) {
	d.consumer.Receive(events)
}

func (d *DispatcherMock) Subscribe(event dispatcher.SubscribeEvent) {
	d.hub.Subscribe(event)
}

func (d *DispatcherMock) Register() {
	d.hub.RegisterChan() <- d
}

func (d *DispatcherMock) Unregister() {
	d.hub.UnregisterChan() <- d
}
