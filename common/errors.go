package common

import "errors"

// ErrInvalidAPIType signals that an invalid api type has been provided
var ErrInvalidAPIType = errors.New("invalid api type provided")

// ErrInvalidDispatchType signals that an invalid dispatch type has been provided
var ErrInvalidDispatchType = errors.New("invalid dispatch type")

// ErrInvalidRedisConnType signals that an invalid redis connection type has been provided
var ErrInvalidRedisConnType = errors.New("invalid redis connection type")

// ErrReceivedEmptyEvents signals that empty events have been received
var ErrReceivedEmptyEvents = errors.New("received empty events")
