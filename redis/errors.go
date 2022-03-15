package redis

import "errors"

// ErrNilRedlockClient signals that a nil redlock client has been provided
var ErrNilRedlockClient = errors.New("nil redlock client")

// ErrRedisConnectionFailed signals that connection to redis failed
var ErrRedisConnectionFailed = errors.New("error connecting to redis")
