package pubsub

import (
	"context"
	"errors"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/go-redis/redis/v8"
)

const (
	pongValue = "PONG"
)

var (
	log = logger.GetOrCreate("notifier/pubsub")

	ErrRedisConnectionFailed = errors.New("error connecting to redis")
)

func CreatePubsubClient(cfg config.PubSubConfig) (RedisClient, error) {
	opt, err := redis.ParseURL(cfg.Url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return NewRedisClientWrapper(client), nil
}

func CreateFailoverClient(cfg config.PubSubConfig) (RedisClient, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: []string{cfg.SentinelUrl},
	})

	rc := NewRedisClientWrapper(client)

	ok := isConnected(rc)
	if !ok {
		return nil, ErrRedisConnectionFailed
	}

	return rc, nil
}

func isConnected(client RedisClient) bool {
	pong, err := client.Ping(context.Background())
	return err == nil && pong == pongValue
}
