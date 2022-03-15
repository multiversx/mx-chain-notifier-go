package redis

import "errors"

// ErrRedisConnectionFailed signals that connection to redis failed
var ErrRedisConnectionFailed = errors.New("error connecting to redis")
