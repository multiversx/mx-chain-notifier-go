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

func (r *RedlockWrapper) IsBlockProcessed(blockHash []byte) (bool, error) {
	ok, err := r.client.SetNX(r.ctx, string(blockHash), true, expiry).Result()
	return ok, err
}
