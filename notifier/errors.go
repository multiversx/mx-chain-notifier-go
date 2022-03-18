package notifier

import "errors"

// ErrInvalidAPIType signals that an invalida api type has been provided
var ErrInvalidAPIType = errors.New("invalid api type provided")

var ErrNilLockService = errors.New("nil lock service")
var ErrNilPublisherService = errors.New("nil publisher service")
var ErrNilConnectorAPIConfig = errors.New("nil connector api config")
