package rabbitmq

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/streadway/amqp"
)

// RabbitMqClient defines the behaviour of a rabbitMq client
type RabbitMqClient interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	IsInterfaceNil() bool
}

// PublisherService defines the behaviour of a publisher component which should be
// able to publish received events and broadcast them to channels
type PublisherService interface {
	Run()
	BroadcastChan() chan<- data.BlockEvents
	BroadcastRevertChan() chan<- data.RevertBlock
	BroadcastFinalizedChan() chan<- data.FinalizedBlock
	IsInterfaceNil() bool
}
