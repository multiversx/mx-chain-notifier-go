package pubsub

import (
	"context"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/go-redis/redis/v8"
)

const (
	pongValue = "PONG"
)

var log = logger.GetOrCreate("notifier/pubsub")

func CreatePubsubClient(cfg config.PubSubConfig) *redis.Client {
	opt, err := redis.ParseURL(cfg.Url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}

func CreateFailoverClient(cfg config.PubSubConfig) (*redis.Client, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: []string{cfg.SentinelUrl},
	})

	pong, err := pingRedis(context.Background(), client)
	if err != nil || pong != pongValue {
		return nil, err
	}

	return client, nil
}

func pingRedis(ctx context.Context, redis *redis.Client) (string, error) {
	return redis.Ping(ctx).Result()
}
