package rabbitmq

import (
	"context"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

var (
	ctx = context.Background()
	url = "amqp://guest:guest@localhost:5672"
)

func TestNewRabbitClientWrapper(t *testing.T) {
	if testing.Short() {
		t.Skip("this test requires running a rabbitMQ server")
	}

	rc, err := NewRabbitClientWrapper(ctx, url)
	require.Nil(t, err)
	require.True(t, !rc.conn.IsClosed())
}

func TestRabbitClientWrapper_Dial(t *testing.T) {
	if testing.Short() {
		t.Skip("this test requires running a rabbitMQ server")
	}

	rc := &rabbitClientWrapper{}

	c, err := rc.Dial(url)
	require.Nil(t, err)
	require.True(t, !c.IsClosed())
}

func TestRabbitClientWrapper_Publish(t *testing.T) {
	if testing.Short() {
		t.Skip("this test requires running a rabbitMQ server")
	}

	rc, err := NewRabbitClientWrapper(ctx, url)
	require.Nil(t, err)

	payload := amqp.Publishing{Body: []byte("test")}
	err = rc.Publish("events", "", true, false, payload)
	require.Nil(t, err)
}
