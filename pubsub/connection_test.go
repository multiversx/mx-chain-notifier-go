package pubsub

import (
	"testing"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/stretchr/testify/require"
)

func TestCreatePubsubClient(t *testing.T) {
	t.Skip("this test requires running a redis server")

	_, err := CreatePubsubClient(config.PubSubConfig{
		Url: "redis://localhost:6379/0",
	})
	require.Nil(t, err)
}

func TestCreateFailoverClient(t *testing.T) {
	t.Skip("this test requires running a redis server")

	_, err := CreateFailoverClient(config.PubSubConfig{
		MasterName:  "my_master",
		SentinelUrl: ":26379",
	})
	require.Nil(t, err)
}
