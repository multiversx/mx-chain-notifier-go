package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/ElrondNetwork/notifier-go/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockRedlockWrapperArgs() redis.ArgsRedlockWrapper {
	return redis.ArgsRedlockWrapper{
		Client:       &mocks.RedisClientMock{},
		TTLInMinutes: 30,
	}
}

func TestNewRedlockWrapper(t *testing.T) {
	t.Parallel()

	t.Run("nil redlock client, should fail", func(t *testing.T) {
		t.Parallel()

		args := createMockRedlockWrapperArgs()
		args.Client = nil

		redlock, err := redis.NewRedlockWrapper(args)
		assert.True(t, check.IfNil(redlock))
		assert.Equal(t, redis.ErrNilRedlockClient, err)
	})

	t.Run("invalid ttl value", func(t *testing.T) {
		t.Parallel()

		args := createMockRedlockWrapperArgs()
		args.TTLInMinutes = 0

		redlock, err := redis.NewRedlockWrapper(args)
		assert.True(t, check.IfNil(redlock))
		assert.True(t, errors.Is(err, common.ErrInvalidValue))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		redlock, err := redis.NewRedlockWrapper(createMockRedlockWrapperArgs())
		assert.Nil(t, err)
		assert.False(t, check.IfNil(redlock))
	})
}

func TestRedlockWrapper_IsBlockProcessed(t *testing.T) {
	t.Parallel()

	t.Run("set should work", func(t *testing.T) {
		t.Parallel()

		args := createMockRedlockWrapperArgs()
		args.Client = &mocks.RedisClientStub{
			SetEntryCalled: func(key string, value bool, ttl time.Duration) (bool, error) {
				return true, nil
			},
		}

		redlock, err := redis.NewRedlockWrapper(args)
		require.Nil(t, err)

		ok, err := redlock.IsEventProcessed(context.Background(), "randStr")
		require.Nil(t, err)
		require.True(t, ok)
	})

	t.Run("should work for new key, but fail for existing key", func(t *testing.T) {
		t.Parallel()

		existingKey := "exists"

		args := createMockRedlockWrapperArgs()
		args.Client = &mocks.RedisClientStub{
			SetEntryCalled: func(key string, value bool, ttl time.Duration) (bool, error) {
				if key == existingKey {
					return false, nil
				}

				return true, nil
			},
		}

		redlock, err := redis.NewRedlockWrapper(args)
		require.Nil(t, err)

		ok, err := redlock.IsEventProcessed(context.Background(), "randStr")
		require.Nil(t, err)
		require.True(t, ok)

		ok, err = redlock.IsEventProcessed(context.Background(), existingKey)
		require.Nil(t, err)
		require.False(t, ok)
	})
}

func TestRedlockWrapper_HasConnection(t *testing.T) {
	t.Parallel()

	args := createMockRedlockWrapperArgs()
	args.Client = &mocks.RedisClientStub{
		IsConnectedCalled: func() bool {
			return true
		},
	}

	redlock, err := redis.NewRedlockWrapper(args)
	require.Nil(t, err)

	ok := redlock.HasConnection(context.Background())
	require.True(t, ok)
}
