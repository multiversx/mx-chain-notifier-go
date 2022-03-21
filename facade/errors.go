package facade

import "errors"

// ErrNilEventsHandler signals that a nil events handler was provided
var ErrNilEventsHandler = errors.New("nil events handler")

// ErrNilHubHandler signals that a nil hub handler was provided
var ErrNilHubHandler = errors.New("nil hub handler")
