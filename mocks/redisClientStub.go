package mocks

import (
	"context"
	"time"
)

// RedisClientStub -
type RedisClientStub struct {
	SetEntryCalled    func(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error)
	PingCalled        func(ctx context.Context) (string, error)
	IsConnectedCalled func(ctx context.Context) bool
}

// SetEntry -
func (rc *RedisClientStub) SetEntry(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error) {
	if rc.SetEntryCalled != nil {
		return rc.SetEntryCalled(ctx, key, value, ttl)
	}

	return false, nil
}

// Ping -
func (rc *RedisClientStub) Ping(ctx context.Context) (string, error) {
	if rc.PingCalled != nil {
		return rc.PingCalled(ctx)
	}

	return "", nil
}

// IsInterfaceNil -
func (rc *RedisClientStub) IsInterfaceNil() bool {
	return false
}
