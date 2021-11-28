package pubsub

import (
	"context"
	"time"
)

var expiry = time.Minute * 30

type RedlockWrapper struct {
	client RedisClient
	ctx    context.Context
}

func NewRedlockWrapper(ctx context.Context, client RedisClient) *RedlockWrapper {
	return &RedlockWrapper{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedlockWrapper) IsBlockProcessed(blockHash string) (bool, error) {
	ok, err := r.client.SetNX(r.ctx, blockHash, true, expiry)
	return ok, err
}

func (r *RedlockWrapper) HasConnection() bool {
	pong, err := r.client.Ping(r.ctx)

	return err == nil && pong == pongValue
}
