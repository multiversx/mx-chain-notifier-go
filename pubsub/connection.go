package pubsub

import (
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

func CreateFailoverClient(cfg config.PubSubConfig) *redis.Client {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: []string{cfg.SentinelUrl},
	})
}
