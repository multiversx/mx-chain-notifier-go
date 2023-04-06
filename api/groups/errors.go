package groups

import "errors"

var errNilBlockData = errors.New("nil block data")
var errNilTransactionPool = errors.New("nil transaction pool")
var errNilHeaderGasConsumption = errors.New("nil header gas consumption")

// ErrNilEventsDataHandler signals that a nil events data handler was provided
var ErrNilEventsDataHandler = errors.New("nil events data handler")
