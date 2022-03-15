package redis

import (
	"context"
	"time"
)

// LockService defines the behaviour of a lock service component.
// It makes sure that a duplicated entry is not processed multiple times,
// it lockes an item once it has been processed.
type LockService interface {
	IsBlockProcessed(blockHash string) (bool, error)
	HasConnection() bool
	IsInterfaceNil() bool
}

// CacheHandler defines the behaviour of a chace handler component
type CacheHandler interface {
	SetEntry(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error)
	Ping(ctx context.Context) (string, error)
	IsConnected(ctx context.Context) bool
	IsInterfaceNil() bool
}
