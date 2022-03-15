package redis

import (
	"context"
	"time"
)

// CacheHandler defines the behaviour of a chace handler component
type CacheHandler interface {
	SetEntry(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error)
	Ping(ctx context.Context) (string, error)
	IsConnected(ctx context.Context) bool
	IsInterfaceNil() bool
}
