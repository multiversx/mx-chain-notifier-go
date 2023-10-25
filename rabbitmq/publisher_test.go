package rabbitmq_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func createMockArgsRabbitMqPublisher() rabbitmq.ArgsRabbitMqPublisher {
	return rabbitmq.ArgsRabbitMqPublisher{
		Client: &mocks.RabbitClientStub{},
		Config: config.RabbitMQConfig{
			EventsExchange: config.RabbitMQExchangeConfig{
				Name: "allevents",
				Type: "fanout",
			},
			RevertEventsExchange: config.RabbitMQExchangeConfig{
				Name: "revert",
				Type: "fanout",
			},
			FinalizedEventsExchange: config.RabbitMQExchangeConfig{
				Name: "finalized",
				Type: "fanout",
			},
			BlockTxsExchange: config.RabbitMQExchangeConfig{
				Name: "blocktxs",
				Type: "fanout",
			},
			BlockScrsExchange: config.RabbitMQExchangeConfig{
				Name: "blockscrs",
				Type: "fanout",
			},
			BlockEventsExchange: config.RabbitMQExchangeConfig{
				Name: "blockeventswithorder",
				Type: "fanout",
			},
		},
		Marshaller: &mock.MarshalizerMock{},
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

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Marshaller = nil

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("invalid events exchange name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Config.EventsExchange.Name = ""

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeName))
	})

	t.Run("invalid revert exchange name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Config.RevertEventsExchange.Name = ""

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeName))
	})

	t.Run("invalid finalized exchange name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Config.FinalizedEventsExchange.Name = ""

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeName))
	})

	t.Run("invalid txs exchange name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Config.BlockTxsExchange.Name = ""

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeName))
	})

	t.Run("invalid scrs exchange name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Config.BlockScrsExchange.Name = ""

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeName))
	})

	t.Run("invalid exchange type", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()
		args.Config.EventsExchange.Type = ""

		client, err := rabbitmq.NewRabbitMqPublisher(args)

		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeType))

		args.Config.RevertEventsExchange.Type = ""
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeType))

		args.Config.FinalizedEventsExchange.Type = ""
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeType))

		args.Config.BlockTxsExchange.Type = ""
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeType))

		args.Config.BlockScrsExchange.Type = ""
		require.True(t, check.IfNil(client))
		require.True(t, errors.Is(err, rabbitmq.ErrInvalidRabbitMqExchangeType))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsRabbitMqPublisher()

		wasCalled := false
		args.Client = &mocks.RabbitClientStub{
			ExchangeDeclareCalled: func(name, kind string) error {
				wasCalled = true
				return nil
			},
		}

		client, err := rabbitmq.NewRabbitMqPublisher(args)
		require.Nil(t, err)
		require.NotNil(t, client)
		require.True(t, wasCalled)
	})
}

func TestPublish(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			wasCalled = true
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.Publish(data.BlockEvents{})

	require.True(t, wasCalled)
}

func TestPublishRevert(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			wasCalled = true
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.PublishRevert(data.RevertBlock{})

	require.True(t, wasCalled)
}

func TestBroadcastFinalized(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			wasCalled = true
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.PublishFinalized(data.FinalizedBlock{})

	require.True(t, wasCalled)
}

func TestBroadcastTxs(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			wasCalled = true
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.PublishTxs(data.BlockTxs{})

	require.True(t, wasCalled)
}

func TestBroadcastScrs(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			wasCalled = true
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.PublishScrs(data.BlockScrs{})

	require.True(t, wasCalled)
}

func TestBroadcastBlockEventsWithOrder(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			wasCalled = true
			return nil
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.PublishBlockEventsWithOrder(data.BlockEventsWithOrder{})

	require.True(t, wasCalled)
}

func TestClose(t *testing.T) {
	t.Parallel()

	wasCalled := false
	client := &mocks.RabbitClientStub{
		CloseCalled: func() {
			wasCalled = true
		},
	}

	args := createMockArgsRabbitMqPublisher()
	args.Client = client

	rabbitmq, err := rabbitmq.NewRabbitMqPublisher(args)
	require.Nil(t, err)

	rabbitmq.Close()
	require.True(t, wasCalled)
}
