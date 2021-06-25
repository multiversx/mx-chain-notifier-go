package factory

import "errors"

var ErrInvalidCustomHubInput = errors.New("failed to make custom hub")

var ErrInvalidDispatchType = errors.New("invalid dispatch type. failed to register dispatchers")
