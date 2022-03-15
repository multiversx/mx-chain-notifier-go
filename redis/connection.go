package redis

import (
	"context"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/go-redis/redis/v8"
)

var log = logger.GetOrCreate("redis")

// TODO: - handle these using factory pattern
// 	     - create LockManager instance in main factory based on Redis

// CreateSimpleClient will create a redis client for a redis setup with one instance
func CreateSimpleClient(cfg config.PubSubConfig) (CacheHandler, error) {
	opt, err := redis.ParseURL(cfg.Url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	rc := NewRedisClientWrapper(client)
	ok := rc.IsConnected(context.Background())
	if !ok {
		return nil, ErrRedisConnectionFailed
	}

	return rc, nil
}

// CreateFailoverClient will create a redis client for a redis setup with sentinel
func CreateFailoverClient(cfg config.PubSubConfig) (CacheHandler, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: []string{cfg.SentinelUrl},
	})
	rc := NewRedisClientWrapper(client)

	ok := rc.IsConnected(context.Background())
	if !ok {
		return nil, ErrRedisConnectionFailed
	}

	return rc, nil
}
