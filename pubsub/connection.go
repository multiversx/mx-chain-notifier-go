package pubsub

import (
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/go-redis/redis/v8"
)

func CreatePubsubClient(cfg config.PubSubConfig) *redis.Client {
	opt, err := redis.ParseURL(cfg.Url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}
