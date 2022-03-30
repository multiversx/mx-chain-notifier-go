package rabbitmq_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgsRabbitMqPublisher() rabbitmq.ArgsRabbitMqPublisher {
	return rabbitmq.ArgsRabbitMqPublisher{
		Client: &mocks.RabbitClientStub{},
		Config: config.RabbitMQConfig{
			EventsExchange:          "allevents",
			RevertEventsExchange:    "revert",
			FinalizedEventsExchange: "finalized",
		},
	}
}

func TestRabbitMqPublisher(t *testing.T) {
	t.Parallel()

	t.Run("nil rabbitmq client", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Client = nil

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.Equal(t, rabbitmq.ErrNilRabbitMqClient, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.Nil(t, err)
		require.NotNil(t, client)
	})
}

func TestBroadcast(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.Run()
	defer rabbitmq.Close()
	wg.Add(1)

	rabbitmq.Broadcast(data.BlockEvents{})

	wg.Wait()

	assert.Equal(t, uint32(1), numCalls)
}

func TestBroadcastRevert(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.Run()
	defer rabbitmq.Close()
	wg.Add(1)

	rabbitmq.BroadcastRevert(data.RevertBlock{})

	wg.Wait()

	assert.Equal(t, uint32(1), numCalls)
}

func TestBroadcastFinalized(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.Run()
	defer rabbitmq.Close()
	wg.Add(1)

	rabbitmq.BroadcastFinalized(data.FinalizedBlock{})

	wg.Wait()

	assert.Equal(t, uint32(1), numCalls)
}

func TestClose(t *testing.T) {
	t.Parallel()

	args := createMockArgsRabbitMqPublisher()

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.Run()

	err = rabbitmq.Close()
	require.Nil(t, err)
}
