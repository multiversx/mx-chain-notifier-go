package rabbitmq

import (
	"context"
	"encoding/json"
	"sync"
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

	broadcast          chan data.BlockEvents
	broadcastRevert    chan data.RevertBlock
	broadcastFinalized chan data.FinalizedBlock
	connErrCh          chan *amqp.Error
	chanErr            chan *amqp.Error

	conn *amqp.Connection
	ch   *amqp.Channel
	cfg  config.RabbitMQConfig

	pubMut sync.Mutex

	ctx context.Context
}

func NewRabbitMqPublisher(
	ctx context.Context,
	cfg config.RabbitMQConfig,
) (*rabbitMqPublisher, error) {
	rp := &rabbitMqPublisher{
		broadcast:          make(chan data.BlockEvents),
		broadcastRevert:    make(chan data.RevertBlock),
		broadcastFinalized: make(chan data.FinalizedBlock),
		cfg:                cfg,
		ctx:                ctx,
		pubMut:             sync.Mutex{},
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
		case err := <-rp.chanErr:
			if err != nil {
				log.Error("rabbitMQ channel failure", "err", err.Error())
				rp.reopenChannel()
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

// BroadcastFinalizedChan returns a receive-only channel on which finalized events are pushed
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastFinalizedChan() chan<- data.FinalizedBlock {
	return rp.broadcastFinalized
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
			log.Error("failed to publish events to rabbitMQ", "err", err.Error())
		}
	}
}

func (rp *rabbitMqPublisher) publishRevertToExchange(revertBlock data.RevertBlock) {
	if rp.cfg.RevertEventsExchange != "" {
		revertBlockBytes, err := json.Marshal(revertBlock)
		if err != nil {
			log.Error("could not marshal revert event", "err", err.Error())
			return
		}

		err = rp.publishFanout(rp.cfg.RevertEventsExchange, revertBlockBytes)
		if err != nil {
			log.Error("failed to publish revert event to rabbitMQ", "err", err.Error())
		}
	}
}

func (rp *rabbitMqPublisher) publishFanout(exchangeName string, payload []byte) error {
	rp.pubMut.Lock()
	defer rp.pubMut.Unlock()

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

	err = rp.openChannel()
	if err != nil {
		return err
	}

	return nil
}

func (rp *rabbitMqPublisher) openChannel() error {
	ch, err := rp.conn.Channel()
	if err != nil {
		return err
	}
	rp.ch = ch

	rp.chanErr = make(chan *amqp.Error)
	rp.ch.NotifyClose(rp.chanErr)

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

func (rp *rabbitMqPublisher) reopenChannel() {
	for {
		time.Sleep(time.Millisecond * reconnectRetryMs)

		err := rp.openChannel()
		if err != nil {
			log.Debug("could not re-open channel", "err", err.Error())
		} else {
			log.Debug("channel opened after reconnect attempts")
			break
		}
	}
}
