package notifier

import "errors"

// ErrNilConfigs signals that nil config has been provided
var ErrNilConfigs = errors.New("nil configs provided")
