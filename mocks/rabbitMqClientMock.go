package mocks

import (
	"sync"

	"github.com/streadway/amqp"
)

// RabbitClientMock -
type RabbitClientMock struct {
	mut    sync.Mutex
	events map[string]amqp.Publishing
}

// NewRabbitClientMock -
func NewRabbitClientMock() *RabbitClientMock {
	return &RabbitClientMock{
		events: make(map[string]amqp.Publishing),
	}
}

// Publish -
func (rc *RabbitClientMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	rc.mut.Lock()
	defer rc.mut.Unlock()

	rc.events[exchange] = msg

	return nil
}

// GetEntries -
func (rc *RabbitClientMock) GetEntries() map[string]amqp.Publishing {
	rc.mut.Lock()
	defer rc.mut.Unlock()

	return rc.events
}

// IsInterfaceNil -
func (rc *RabbitClientMock) IsInterfaceNil() bool {
	return rc == nil
}
