package redis

import (
	"context"
	"time"
)

const expiry = time.Minute * 30

type redlockWrapper struct {
	client CacheHandler
	ctx    context.Context
}

func NewRedlockWrapper(ctx context.Context, client CacheHandler) *redlockWrapper {
	return &redlockWrapper{
		client: client,
		ctx:    ctx,
	}
}

func (r *redlockWrapper) IsBlockProcessed(blockHash string) (bool, error) {
	return r.client.SetEntry(r.ctx, blockHash, true, expiry)
}

func (r *redlockWrapper) HasConnection() bool {
	return r.client.IsConnected(r.ctx)
}

func (r *redlockWrapper) IsInterfaceNil() bool {
	return r == nil
}
