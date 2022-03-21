package common

import "errors"

// ErrInvalidAPIType signals that an invalid api type has been provided
var ErrInvalidAPIType = errors.New("invalid api type provided")
