package rabbitmq

import "errors"

// ErrNilRabbitMqClient signals that a nil rabbitmq client has been provided
var ErrNilRabbitMqClient = errors.New("nil rabbitmq client")

// ErrInvalidRabbitMqExchangeName signals that an empty rabbitmq exchange name has been provided
var ErrInvalidRabbitMqExchangeName = errors.New("invalid rabbitmq exchange name")

// ErrInvalidRabbitMqExchangeType signals that an empty rabbitmq exchange type has been provided
var ErrInvalidRabbitMqExchangeType = errors.New("invalid rabbitmq exchange type")
