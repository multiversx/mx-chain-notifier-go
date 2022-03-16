package rabbitmq

import "github.com/streadway/amqp"

// RabbitMqClient defines the behaviour of a rabbitMq client
type RabbitMqClient interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	IsInterfaceNil() bool
}
