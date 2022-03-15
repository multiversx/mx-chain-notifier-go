package mocks

import (
	"context"
	"time"
)

// RedisClientStub -
type RedisClientStub struct {
	SetEntryCalled    func(key string, value bool, ttl time.Duration) (bool, error)
	PingCalled        func() (string, error)
	IsConnectedCalled func() bool
}

// SetEntry -
func (rc *RedisClientStub) SetEntry(_ context.Context, key string, value bool, ttl time.Duration) (bool, error) {
	if rc.SetEntryCalled != nil {
		return rc.SetEntryCalled(key, value, ttl)
	}

	return false, nil
}

// Ping -
func (rc *RedisClientStub) Ping(_ context.Context) (string, error) {
	if rc.PingCalled != nil {
		return rc.PingCalled()
	}

	return "", nil
}

// IsConnected -
func (rc *RedisClientStub) IsConnected(_ context.Context) bool {
	if rc.IsConnectedCalled != nil {
		return rc.IsConnectedCalled()
	}

	return false
}

// IsInterfaceNil -
func (rc *RedisClientStub) IsInterfaceNil() bool {
	return false
}
