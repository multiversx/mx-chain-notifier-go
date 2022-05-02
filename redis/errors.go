package redis

import "errors"

// ErrNilRedlockClient signals that a nil redlock client has been provided
var ErrNilRedlockClient = errors.New("nil redlock client")

// ErrRedisConnectionFailed signals that connection to redis failed
var ErrRedisConnectionFailed = errors.New("error connecting to redis")

// ErrZeroValueReceived signals that a zero value has been received
var ErrZeroValueReceived = errors.New("zero value received")
