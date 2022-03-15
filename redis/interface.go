package redis

import (
	"context"
	"time"
)

// RedLockClient defines the behaviour of a cache handler component
type RedLockClient interface {
	SetEntry(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error)
	Ping(ctx context.Context) (string, error)
	IsConnected(ctx context.Context) bool
	IsInterfaceNil() bool
}
