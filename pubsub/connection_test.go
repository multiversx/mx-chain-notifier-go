package pubsub

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

	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()

	require.Nil(t, err)
	require.True(t, pong == "PONG")
}

func TestCreatePubsubClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	rl := NewRedlockWrapper(ctx, client)

	require.True(t, rl.HasConnection())
}
