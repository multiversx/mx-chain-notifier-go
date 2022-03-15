package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/ElrondNetwork/notifier-go/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestNewRedlockWrapper(t *testing.T) {
	t.Parallel()

	t.Run("nil redlock client, should fail", func(t *testing.T) {
		t.Parallel()

		redlock, err := redis.NewRedlockWrapper(ctx, nil)
		assert.True(t, check.IfNil(redlock))
		assert.Equal(t, redis.ErrNilRedlockClient, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		redlock, err := redis.NewRedlockWrapper(ctx, &mocks.RedisClientStub{})
		assert.Nil(t, err)
		assert.False(t, check.IfNil(redlock))
	})
}

func TestRedlockWrapper_IsBlockProcessed(t *testing.T) {
	t.Parallel()

	t.Run("set should work", func(t *testing.T) {
		client := &mocks.RedisClientStub{
			SetEntryCalled: func(key string, value bool, ttl time.Duration) (bool, error) {
				return true, nil
			},
		}

		redlock, err := redis.NewRedlockWrapper(ctx, client)
		require.Nil(t, err)

		ok, err := redlock.IsBlockProcessed("randStr")
		require.Nil(t, err)
		require.True(t, ok)
	})

	t.Run("should work for new key, but fail for existing key", func(t *testing.T) {
		t.Parallel()

		existingKey := "exists"

		client := &mocks.RedisClientStub{
			SetEntryCalled: func(key string, value bool, ttl time.Duration) (bool, error) {
				if key == existingKey {
					return false, nil
				}

				return true, nil
			},
		}

		redlock, err := redis.NewRedlockWrapper(ctx, client)
		require.Nil(t, err)

		ok, err := redlock.IsBlockProcessed("randStr")
		require.Nil(t, err)
		require.True(t, ok)

		ok, err = redlock.IsBlockProcessed(existingKey)
		require.Nil(t, err)
		require.False(t, ok)
	})
}

func TestRedlockWrapper_HasConnection(t *testing.T) {
	t.Parallel()

	client := &mocks.RedisClientStub{
		IsConnectedCalled: func() bool {
			return true
		},
	}

	redlock, err := redis.NewRedlockWrapper(ctx, client)
	require.Nil(t, err)

	ok := redlock.HasConnection()
	require.True(t, ok)
}
