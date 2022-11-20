package process

import "errors"

// ErrNilLockService signals that a nil lock service has been provided
var ErrNilLockService = errors.New("nil lock service")

// ErrNilPublisherService signals that a nil publisher service has been provided
var ErrNilPublisherService = errors.New("nil publisher service")

// ErrInvalidValue signals that an invalid value has been provided
var ErrInvalidValue = errors.New("invalid value")

// ErrNilPubKeyConverter signals that a nil pubkey converter has been provided
var ErrNilPubKeyConverter = errors.New("nil pubkey converter")
