package rabbitmq

import "errors"

// ErrNilRabbitMqClient signals that a nil rabbitmq client has been provided
var ErrNilRabbitMqClient = errors.New("nil rabbitmq client")

// ErrInvalidRabbitMqExchangeName signals that an empty rabbitmq exchange name has been provided
var ErrInvalidRabbitMqExchangeName = errors.New("invalid rabbitmq exchange name")
