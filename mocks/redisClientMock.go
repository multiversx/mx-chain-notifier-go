package mocks

import (
	"context"
	"sync"
	"time"
)

// RedisClientMock -
type RedisClientMock struct {
	mut     sync.Mutex
	entries map[string]bool
}

// NewRedisClientMock -
func NewRedisClientMock() *RedisClientMock {
	return &RedisClientMock{
		entries: make(map[string]bool),
	}
}

// SetEntry -
func (rc *RedisClientMock) SetEntry(_ context.Context, key string, value bool, ttl time.Duration) (bool, error) {
	rc.mut.Lock()
	defer rc.mut.Unlock()

	willSet := true
	for k, val := range rc.entries {
		if k == key && val == value {
			willSet = false
			break
		}
	}

	if willSet {
		rc.entries[key] = value
		return true, nil
	}

	return false, nil
}

// GetEntries -
func (rc *RedisClientMock) GetEntries() map[string]bool {
	rc.mut.Lock()
	defer rc.mut.Unlock()

	return rc.entries
}

// Ping -
func (rc *RedisClientMock) Ping(_ context.Context) (string, error) {
	return "PONG", nil
}

// IsConnected -
func (rc *RedisClientMock) IsConnected(_ context.Context) bool {
	return true
}

// IsInterfaceNil -
func (rc *RedisClientMock) IsInterfaceNil() bool {
	return rc == nil
}
