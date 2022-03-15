package rabbitmq

import "github.com/streadway/amqp"

// Client defines the behaviour of a rabbitMq client
type Client interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Dial(url string) (*amqp.Connection, error)
	IsInterfaceNil() bool
}
