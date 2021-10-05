package rabbitmq

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()
var cfg = config.RabbitMQConfig{
	Url:            "amqp://guest:guest@localhost:5672",
	EventsExchange: "events",
}

func wait() {
	c := make(chan int)
	<-c
}

func TestNewRabbitMqPublisher(t *testing.T) {
	t.Parallel()

	r, err := NewRabbitMqPublisher(ctx, cfg)
	require.Nil(t, err)

	go r.Run()

	wait()
}

func TestRabbitMqPublisher_BroadcastChan(t *testing.T) {
	t.Parallel()

	r, err := NewRabbitMqPublisher(ctx, cfg)
	require.Nil(t, err)

	events := []data.Event{
		{
			Address:    "erd1",
			Identifier: "id1",
			Topics:     [][]byte{[]byte("topic1"), []byte("topic2")},
		},
	}

	r.publishToExchanges(events)
}
