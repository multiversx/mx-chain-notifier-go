package pubsub

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/stretchr/testify/require"
)

func TestCreatePubsubClient_PingShouldConnectToDefault(t *testing.T) {
	t.Parallel()

	r := CreatePubsubClient(config.PubSubConfig{
		Url: "redis://localhost:6379/0",
	})

	ctx := context.Background()

	pong, err := r.Ping(ctx).Result()

	require.Nil(t, err)
	require.True(t, pong == "PONG")
}
