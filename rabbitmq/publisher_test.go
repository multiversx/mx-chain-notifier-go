package rabbitmq

import (
	"encoding/json"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/test/mocks"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	cfg = config.RabbitMQConfig{
		Url:                     "amqp://guest:guest@localhost:5672",
		EventsExchange:          "events",
		RevertEventsExchange:    "revert",
		FinalizedEventsExchange: "finalized",
	}

	event = data.BlockEvents{
		Hash: "abcdef",
		Events: []data.Event{
			{
				Address:    "erd1",
				Identifier: "test",
			},
		},
	}

	revertEvent = data.RevertBlock{
		Hash:  "abcdef",
		Nonce: 1,
		Round: 100,
		Epoch: 2,
	}

	finalizedEvent = data.FinalizedBlock{
		Hash: "abcdef",
	}
)

func TestRabbitMqPublisher_PublishEventsShouldWork(t *testing.T) {
	t.Parallel()

	var received data.BlockEvents
	rc := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			_ = json.Unmarshal(msg.Body, &received)
			require.Equal(t, event, received)

			return nil
		},
	}

	rp := NewRabbitMqPublisher(rc, cfg)

	rp.publishToExchanges(event)
}

func TestRabbitMqPublisher_PublishEventsNoEventsExchangeSpecified(t *testing.T) {
	t.Parallel()

	calledCnt := 0

	rc := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			calledCnt++
			return nil
		},
	}

	cfg.EventsExchange = ""
	rp := NewRabbitMqPublisher(rc, cfg)

	rp.publishToExchanges(event)
	require.True(t, calledCnt == 0)
}

func TestRabbitMqPublisher_PublishRevertEventsShouldWork(t *testing.T) {
	t.Parallel()

	var received data.RevertBlock
	rc := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			_ = json.Unmarshal(msg.Body, &received)
			require.Equal(t, revertEvent, received)

			return nil
		},
	}

	rp := NewRabbitMqPublisher(rc, cfg)

	rp.publishRevertToExchange(revertEvent)
}

func TestRabbitMqPublisher_PublishRevertEventsNoRevertEventsExchangeSpecified(t *testing.T) {
	t.Parallel()

	calledCnt := 0

	rc := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			calledCnt++
			return nil
		},
	}

	cfg.RevertEventsExchange = ""
	rp := NewRabbitMqPublisher(rc, cfg)

	rp.publishRevertToExchange(revertEvent)
	require.True(t, calledCnt == 0)
}

func TestRabbitMqPublisher_PublishFinalizedEventsShouldWork(t *testing.T) {
	t.Parallel()

	var received data.FinalizedBlock
	rc := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			_ = json.Unmarshal(msg.Body, &received)
			require.Equal(t, finalizedEvent, received)

			return nil
		},
	}

	rp := NewRabbitMqPublisher(rc, cfg)

	rp.publishFinalizedToExchange(finalizedEvent)
}

func TestRabbitMqPublisher_PublishFinalizedEventsNoEventsExchangeSpecified(t *testing.T) {
	t.Parallel()

	calledCnt := 0

	rc := &mocks.RabbitClientStub{
		PublishCalled: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			calledCnt++
			return nil
		},
	}

	cfg.FinalizedEventsExchange = ""
	rp := NewRabbitMqPublisher(rc, cfg)

	rp.publishFinalizedToExchange(finalizedEvent)
	require.True(t, calledCnt == 0)
}
