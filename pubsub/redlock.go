package pubsub

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var expiry = time.Minute * 30

type RedlockWrapper struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedlockWrapper(ctx context.Context, client *redis.Client) *RedlockWrapper {
	return &RedlockWrapper{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedlockWrapper) IsBlockProcessed(blockHash string) (bool, error) {
	ok, err := r.client.SetNX(r.ctx, blockHash, true, expiry).Result()
	return ok, err
}

func (r *RedlockWrapper) HasConnection() bool {
	pong, err := r.client.Ping(r.ctx).Result()

	return err == nil && pong == pongValue
}
