package hub

import "errors"

// ErrNilEventFilter signals that a nil event filter has been provided
var ErrNilEventFilter = errors.New("nil events filter")

// ErrNilSubscriptionMapper signals that a nil subscription mapper has been provided
var ErrNilSubscriptionMapper = errors.New("nil subscription mapper")
