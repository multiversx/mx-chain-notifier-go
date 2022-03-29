package factory

import (
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/redis"
)

// CreateLockService creates lock service component based on config
func CreateLockService(checkDuplicates bool, config config.RedisConfig) (redis.LockService, error) {
	if !checkDuplicates {
		return disabled.NewDisabledRedlockWrapper(), nil
	}

	redisClient, err := createRedisClient(config)
	if err != nil {
		return nil, err
	}

	lockService, err := redis.NewRedlockWrapper(redisClient)
	if err != nil {
		return nil, err
	}

	return lockService, nil
}

func createRedisClient(cfg config.RedisConfig) (redis.RedLockClient, error) {
	switch cfg.ConnectionType {
	case common.RedisInstanceConnType:
		return redis.CreateSimpleClient(cfg)
	case common.RedisSentinelConnType:
		return redis.CreateFailoverClient(cfg)
	default:
		return nil, common.ErrInvalidRedisConnType
	}
}
