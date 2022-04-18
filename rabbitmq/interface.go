package rabbitmq

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/streadway/amqp"
)

// RabbitMqClient defines the behaviour of a rabbitMq client
type RabbitMqClient interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	ConnErrChan() chan *amqp.Error
	CloseErrChan() chan *amqp.Error
	Reconnect()
	ReopenChannel()
	Close()
	IsInterfaceNil() bool
}

// PublisherService defines the behaviour of a publisher component which should be
// able to publish received events and broadcast them to channels
type PublisherService interface {
	Run()
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	BroadcastTxs(event data.BlockTxs)
	Close() error
	IsInterfaceNil() bool
}
