package mocks

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClientStub struct {
	SetNXCalled            func(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error)
	PublishCalled          func(ctx context.Context, channel string, message interface{}) error
	SubscribeGetChanCalled func(ctx context.Context, channel string) <-chan *redis.Message
	PingCalled             func(ctx context.Context) (string, error)
}

func (rc *RedisClientStub) SetNX(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error) {
	if rc.SetNXCalled != nil {
		return rc.SetNXCalled(ctx, key, value, ttl)
	}

	return false, nil
}

func (rc *RedisClientStub) Publish(ctx context.Context, channel string, message interface{}) error {
	if rc.PublishCalled != nil {
		return rc.PublishCalled(ctx, channel, message)
	}

	return nil
}

func (rc *RedisClientStub) SubscribeGetChan(ctx context.Context, channel string) <-chan *redis.Message {
	if rc.SubscribeGetChanCalled != nil {
		return rc.SubscribeGetChanCalled(ctx, channel)
	}

	return nil
}

func (rc *RedisClientStub) Ping(ctx context.Context) (string, error) {
	if rc.PingCalled != nil {
		return rc.PingCalled(ctx)
	}

	return "", nil
}
