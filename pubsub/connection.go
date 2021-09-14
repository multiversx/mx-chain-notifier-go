package pubsub

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/go-redis/redis/v8"
)

var log = logger.GetOrCreate("notifier/pubsub")

func CreatePubsubClient(cfg config.PubSubConfig) *redis.Client {
	opt, err := redis.ParseURL(cfg.Url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}
