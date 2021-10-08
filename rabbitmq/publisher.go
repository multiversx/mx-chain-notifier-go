package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/streadway/amqp"
)

const (
	reconnectRetryMs = 500

	emptyStr = ""
)

var log = logger.GetOrCreate("rabbitMQPublisher")

type rabbitMqPublisher struct {
	dispatcher.Hub

	broadcast       chan data.BlockEvents
	broadcastRevert chan data.RevertBlock
	connErrCh       chan *amqp.Error

	conn *amqp.Connection
	ch   *amqp.Channel
	cfg  config.RabbitMQConfig

	ctx context.Context
}

func NewRabbitMqPublisher(
	ctx context.Context,
	cfg config.RabbitMQConfig,
) (*rabbitMqPublisher, error) {
	rp := &rabbitMqPublisher{
		broadcast:       make(chan data.BlockEvents),
		broadcastRevert: make(chan data.RevertBlock),
		cfg:             cfg,
		ctx:             ctx,
	}

	err := rp.connect()
	if err != nil {
		return nil, err
	}

	return rp, nil
}

// Run is launched as a goroutine and listens for events on the exposed channels
func (rp *rabbitMqPublisher) Run() {
	for {
		select {
		case events := <-rp.broadcast:
			rp.publishToExchanges(events)
		case revertBlock := <-rp.broadcastRevert:
			rp.publishRevertToExchange(revertBlock)
		case err := <-rp.connErrCh:
			if err != nil {
				log.Error("rabbitMQ connection failure", "err", err.Error())
				rp.reconnect()
			}
		case <-rp.ctx.Done():
			err := rp.ch.Close()
			if err != nil {
				log.Error("failed to close rabbitMQ channel", "err", err.Error())
			}
			err = rp.conn.Close()
			if err != nil {
				log.Error("failed to close rabbitMQ channel", "err", err.Error())
			}
		}
	}
}

// BroadcastChan returns a receive-only channel on which events are pushed by producers
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastChan() chan<- data.BlockEvents {
	return rp.broadcast
}

// BroadcastRevertChan returns a receive-only channel on which revert events are pushed by producers
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastRevertChan() chan<- data.RevertBlock {
	return rp.broadcastRevert
}

func (rp *rabbitMqPublisher) publishToExchanges(events data.BlockEvents) {
	if rp.cfg.EventsExchange != "" {
		eventsBytes, err := json.Marshal(events)
		if err != nil {
			log.Error("could not marshal events", "err", err.Error())
			return
		}

		err = rp.publishFanout(rp.cfg.EventsExchange, eventsBytes)
		if err != nil {
			log.Error("failed to publish to rabbitMQ", "err", err.Error())
		}
	}
}

func (rp *rabbitMqPublisher) publishRevertToExchange(revertBlock data.RevertBlock) {
}

func (rp *rabbitMqPublisher) publishFanout(exchangeName string, payload []byte) error {
	err := rp.ch.Publish(
		exchangeName,
		emptyStr,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			Body: payload,
		},
	)

	return err
}

func (rp *rabbitMqPublisher) connect() error {
	conn, err := amqp.Dial(rp.cfg.Url)
	if err != nil {
		return err
	}
	rp.conn = conn

	rp.connErrCh = make(chan *amqp.Error)
	rp.conn.NotifyClose(rp.connErrCh)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	rp.ch = ch

	return nil
}

func (rp *rabbitMqPublisher) reconnect() {
	for {
		time.Sleep(time.Millisecond * reconnectRetryMs)

		err := rp.connect()
		if err != nil {
			log.Debug("could not reconnect", "err", err.Error())
		} else {
			log.Debug("connection established after reconnect attempts")
			break
		}
	}
}
