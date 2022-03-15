package redis

import (
	"context"
	"time"
)

// TODO: update this after LockService interface change (if needed)

const expiry = time.Minute * 30

type redlockWrapper struct {
	client Client
	ctx    context.Context
}

// NewRedlockWrapper create a new redLock based on a chance instance
func NewRedlockWrapper(ctx context.Context, client Client) *redlockWrapper {
	return &redlockWrapper{
		client: client,
		ctx:    ctx,
	}
}

// IsBlockProcessed returns wether the item is already locked
func (r *redlockWrapper) IsBlockProcessed(blockHash string) (bool, error) {
	return r.client.SetEntry(r.ctx, blockHash, true, expiry)
}

func (r *redlockWrapper) HasConnection() bool {
	return r.client.IsConnected(r.ctx)
}

func (r *redlockWrapper) IsInterfaceNil() bool {
	return r == nil
}
