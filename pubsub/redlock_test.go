package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/ElrondNetwork/notifier-go/test"
	"github.com/ElrondNetwork/notifier-go/test/mocks"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestRedlockWrapper_IsBlockProcessedShouldSET(t *testing.T) {
	client := &mocks.RedisClientStub{
		SetNXCalled: func(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error) {
			return true, nil
		},
	}

	redlock := NewRedlockWrapper(ctx, client)
	ok, err := redlock.IsBlockProcessed(test.RandStr(2))
	require.Nil(t, err)
	require.True(t, ok)
}

func TestRedlockWrapper_IsBlockProcessedShouldSETForNewKeyFailForExisting(t *testing.T) {
	existingKey := "exists"

	client := &mocks.RedisClientStub{
		SetNXCalled: func(ctx context.Context, key string, value bool, ttl time.Duration) (bool, error) {
			if key == existingKey {
				return false, nil
			}

			return true, nil
		},
	}

	redlock := NewRedlockWrapper(ctx, client)
	ok, err := redlock.IsBlockProcessed(test.RandStr(2))
	require.Nil(t, err)
	require.True(t, ok)

	ok, err = redlock.IsBlockProcessed(existingKey)
	require.Nil(t, err)
	require.False(t, ok)
}

func TestRedlockWrapper_HasConnection(t *testing.T) {
	client := &mocks.RedisClientStub{
		PingCalled: func(ctx context.Context) (string, error) {
			return "PONG", nil
		},
	}

	redlock := NewRedlockWrapper(ctx, client)
	ok := redlock.HasConnection()
	require.True(t, ok)
}
