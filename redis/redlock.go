package redis

import (
	"context"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
)

const expiry = time.Minute * 30

type redlockWrapper struct {
	client RedLockClient
}

// NewRedlockWrapper create a new redLock based on a cache instance
func NewRedlockWrapper(client RedLockClient) (*redlockWrapper, error) {
	if check.IfNil(client) {
		return nil, ErrNilRedlockClient
	}

	return &redlockWrapper{
		client: client,
	}, nil
}

// IsEventProcessed returns wether the item is already locked
func (r *redlockWrapper) IsEventProcessed(ctx context.Context, blockHash string) (bool, error) {
	return r.client.SetEntry(ctx, blockHash, true, expiry)
}

// HasConnection return true if the redis client is connected
func (r *redlockWrapper) HasConnection(ctx context.Context) bool {
	return r.client.IsConnected(ctx)
}

// IsInterfaceNil returns true if there is no value under the interface
func (r *redlockWrapper) IsInterfaceNil() bool {
	return r == nil
}
