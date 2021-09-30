package pubsub

import (
	"context"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

const expiry = time.Minute * 5

type RedlockWrapper struct {
	client *goredislib.Client
	rsync  *redsync.Redsync
	ctx    context.Context
}

func NewRedlockWrapper(client *goredislib.Client) *RedlockWrapper {
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	return &RedlockWrapper{
		client: client,
		rsync:  rs,
	}
}

func (r *RedlockWrapper) IsBlockProcessed(blockHash []byte) bool {
	ok, err := r.client.SetNX(r.ctx, string(blockHash), true, expiry).Result()
	return err == nil && ok
}
