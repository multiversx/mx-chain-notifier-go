package mocks

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// ConsumerMock -
type ConsumerMock struct {
	id              uuid.UUID
	mut             sync.Mutex
	collectedEvents []data.Event
}

// NewConsumerMock -
func NewConsumerMock() *ConsumerMock {
	return &ConsumerMock{
		id:              uuid.New(),
		mut:             sync.Mutex{},
		collectedEvents: []data.Event{},
	}
}

// Receive -
func (c *ConsumerMock) Receive(events []data.Event) {
	c.mut.Lock()
	defer c.mut.Unlock()
	for _, event := range events {
		c.collectedEvents = append(c.collectedEvents, event)
	}
}

// CollectedEvents -
func (c *ConsumerMock) CollectedEvents() []data.Event {
	c.mut.Lock()
	defer c.mut.Unlock()

	return c.collectedEvents
}

// HasEvents -
func (c *ConsumerMock) HasEvents(events []data.Event) bool {
	c.mut.Lock()
	defer c.mut.Unlock()

	for _, event := range events {
		exists := false
		for _, rcvEvent := range c.collectedEvents {
			if reflect.DeepEqual(event, rcvEvent) {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}

	return true
}

// HasEvent -
func (c *ConsumerMock) HasEvent(event data.Event) bool {
	c.mut.Lock()
	defer c.mut.Unlock()

	for _, rcvEvent := range c.collectedEvents {
		if reflect.DeepEqual(event, rcvEvent) {
			return true
		}
	}

	return false
}
