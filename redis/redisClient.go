package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	pongValue = "PONG"
)

type redisClientWrapper struct {
	redis *redis.Client
}

// NewRedisClientWrapper creates a wrapper over redis client instance. This is an abstraction
// over different redis connections: single instance, sentinel, cluster.
func NewRedisClientWrapper(redis *redis.Client) *redisClientWrapper {
	return &redisClientWrapper{
		redis: redis,
	}
}

// SetEntry will try to update a key value entry in redis database
func (rc *redisClientWrapper) SetEntry(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error) {
	return rc.redis.SetNX(ctx, key, value, ttl).Result()
}

// Ping will check if Redis instance is reponding
func (rc *redisClientWrapper) Ping(ctx context.Context) (string, error) {
	return rc.redis.Ping(ctx).Result()
}

// IsConnected will check if Redis is connected
func (rc *redisClientWrapper) IsConnected(ctx context.Context) bool {
	pong, err := rc.Ping(context.Background())
	return err == nil && pong == pongValue
}

// IsInterfaceNil returns true if there is no value under the interface
func (rc *redisClientWrapper) IsInterfaceNil() bool {
	return rc == nil
}
