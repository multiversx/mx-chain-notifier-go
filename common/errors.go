package common

import "errors"

// ErrInvalidValue signals that an invalid value has been provided
var ErrInvalidValue = errors.New("invalid value")

// ErrInvalidAPIType signals that an invalid api type has been provided
var ErrInvalidAPIType = errors.New("invalid api type provided")

// ErrInvalidDispatchType signals that an invalid dispatch type has been provided
var ErrInvalidDispatchType = errors.New("invalid dispatch type")

// ErrInvalidRedisConnType signals that an invalid redis connection type has been provided
var ErrInvalidRedisConnType = errors.New("invalid redis connection type")
