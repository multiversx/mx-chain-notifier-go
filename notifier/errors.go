package notifier

import "errors"

// ErrInvalidAPIType signals that an invalida api type has been provided
var ErrInvalidAPIType = errors.New("invalid api type provided")

// ErrNilLockService signals that a nil lock service has been provided
var ErrNilLockService = errors.New("nil lock service")

// ErrNilPublisherService signals that a nil publisher service has been provided
var ErrNilPublisherService = errors.New("nil publisher service")

// ErrInvalidValue signals that an invalid value has been provided
var ErrInvalidValue = errors.New("invalid value")
