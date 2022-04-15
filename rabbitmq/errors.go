package rabbitmq

import "errors"

// ErrNilRabbitMqClient signals that a nil rabbitmq client has been provided
var ErrNilRabbitMqClient = errors.New("nil rabbitmq client")
