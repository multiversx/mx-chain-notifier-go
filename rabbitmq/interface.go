package rabbitmq

import "github.com/streadway/amqp"

type RabbitClient interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Dial(url string) (*amqp.Connection, error)
	IsInterfaceNil() bool
}
