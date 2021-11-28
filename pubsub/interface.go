package pubsub

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type LockService interface {
	IsBlockProcessed(blockHash string) (bool, error)
	HasConnection() bool
}

type RedisClient interface {
	SetNX(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error)
	Publish(ctx context.Context, channel string, message interface{}) error
	SubscribeGetChan(ctx context.Context, channel string) <-chan *redis.Message
	Ping(ctx context.Context) (string, error)
}
