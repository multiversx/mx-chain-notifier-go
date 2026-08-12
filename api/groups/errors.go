package groups

import "errors"

var errNilBlockData = errors.New("nil block data")
var errNilTransactionPool = errors.New("nil transaction pool")
var errNilHeaderGasConsumption = errors.New("nil header gas consumption")

// ErrNilEventsDataHandler signals that a nil events data handler was provided
var ErrNilEventsDataHandler = errors.New("nil events data handler")

// ErrAuthEnabledWithoutMiddleware signals that a route was configured with Auth = true
// but no real auth middleware was ever set up for its group, which would otherwise
// silently register the route without any authentication
var ErrAuthEnabledWithoutMiddleware = errors.New("route has auth enabled but no auth middleware is configured")
