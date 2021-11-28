package pubsub

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisClientWrapper struct {
	redis *redis.Client
}

func NewRedisClientWrapper(redis *redis.Client) *redisClientWrapper {
	return &redisClientWrapper{
		redis: redis,
	}
}

func (rc *redisClientWrapper) SetNX(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error) {
	return rc.redis.SetNX(ctx, key, value, ttl).Result()
}

func (rc *redisClientWrapper) Publish(ctx context.Context, channel string, message interface{}) error {
	return rc.redis.Publish(ctx, channel, message).Err()
}

func (rc *redisClientWrapper) SubscribeGetChan(ctx context.Context, channel string) <-chan *redis.Message {
	subChan := rc.redis.Subscribe(ctx, channel)
	return subChan.Channel()
}

func (rc *redisClientWrapper) Ping(ctx context.Context) (string, error) {
	return rc.redis.Ping(ctx).Result()
}
