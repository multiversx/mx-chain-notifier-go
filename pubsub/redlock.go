package pubsub

import (
	"context"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var expiry = time.Minute * 5

type RedlockWrapper struct {
	client *goredislib.Client
	rsync  *redsync.Redsync
	ctx    context.Context
}

func NewRedlockWrapper(ctx context.Context, client *goredislib.Client) *RedlockWrapper {
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	return &RedlockWrapper{
		client: client,
		rsync:  rs,
		ctx:    ctx,
	}
}

func (r *RedlockWrapper) IsBlockProcessed(blockHash []byte) (bool, error) {
	ok, err := r.client.SetNX(r.ctx, string(blockHash), true, expiry).Result()
	return ok, err
}
