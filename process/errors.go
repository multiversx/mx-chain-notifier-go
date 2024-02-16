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

// ErrNilBlockEvents signals that a nil block events struct has been provided
var ErrNilBlockEvents = errors.New("nil block events provided")

// ErrNilTransactionsPool signals that a nil transactions pool has been provided
var ErrNilTransactionsPool = errors.New("nil transactions pool provided")

// ErrNilBlockBody signals that a nil block body has been provided
var ErrNilBlockBody = errors.New("nil block body provided")

// ErrNilBlockHeader signals that a nil block header has been provided
var ErrNilBlockHeader = errors.New("nil block header provided")

// ErrNilPublisherHandler signals that a nil publisher handler has been provided
var ErrNilPublisherHandler = errors.New("nil publisher handler provided")

// ErrNilEventsInterceptor signals that a nil events interceptor was provided
var ErrNilEventsInterceptor = errors.New("nil events interceptor")
