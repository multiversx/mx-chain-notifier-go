package rabbitmq

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/streadway/amqp"
)

// RabbitMqClient defines the behaviour of a rabbitMq client
type RabbitMqClient interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	ExchangeDeclare(name, kind string) error
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
	Run() error
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	BroadcastTxs(event data.BlockTxs)
	BroadcastScrs(event data.BlockScrs)
	BroadcastBlockEventsWithOrder(event data.BlockEventsWithOrder)
	Close() error
	IsInterfaceNil() bool
}
