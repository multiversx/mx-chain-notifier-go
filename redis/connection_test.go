package redis

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/stretchr/testify/require"
)

var client = CreatePubsubClient(config.PubSubConfig{
	Url: "redis://localhost:6379/0",
})

func TestCreatePubsubClient_PingShouldConnectToDefault(t *testing.T) {
	t.Parallel()

	// TODO remove this skip
	t.Skip("fix this test by adding an abstraction layer between the Redis client and what we use")

	ctx := context.Background()

	pong, err := pingRedis(ctx, client)

	require.Nil(t, err)
	require.True(t, pong == "PONG")
}

func TestCreatePubsubClient(t *testing.T) {
	t.Parallel()

	// TODO remove this skip
	t.Skip("fix this test by adding an abstraction layer between the Redis client and what we use")

	ctx := context.Background()

	rl := NewRedlockWrapper(ctx, client)

	require.True(t, rl.HasConnection())
}

func TestCreateFailoverClient_PingShouldConnectToDefault(t *testing.T) {
	t.Parallel()

	// TODO remove this skip
	t.Skip("fix this test by adding an abstraction layer between the Redis client and what we use")

	_, err := CreateFailoverClient(config.PubSubConfig{
		MasterName:  "my_master",
		SentinelUrl: ":26379",
	})

	require.Nil(t, err)
}
