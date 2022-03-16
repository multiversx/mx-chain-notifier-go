package rabbitmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

const (
	reconnectRetryMs = 500
)

type rabbitMqClient struct {
	url string

	ctx context.Context

	conn *amqp.Connection
	ch   *amqp.Channel

	connErrCh chan *amqp.Error
	chanErr   chan *amqp.Error
}

// NewRabbitMQClient creates a new rabbitMQ client instance
func NewRabbitMQClient(ctx context.Context, url string) (*rabbitMqClient, error) {
	rc := &rabbitMqClient{
		ctx: ctx,
		url: url,
	}

	err := rc.connect()
	if err != nil {
		return nil, err
	}

	go rc.connListener()

	return rc, nil
}

// Publish will publich an item on the rabbitMq channel
func (rc *rabbitMqClient) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return rc.ch.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

// dial will return a rabbitMq connection
func (rc *rabbitMqClient) dial(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func (rc *rabbitMqClient) connListener() {
	for {
		select {
		case err := <-rc.connErrCh:
			if err != nil {
				log.Error("rabbitMQ connection failure", "err", err.Error())
				rc.reconnect()
			}
		case err := <-rc.chanErr:
			if err != nil {
				log.Error("rabbitMQ channel failure", "err", err.Error())
				rc.reopenChannel()
			}
		case <-rc.ctx.Done():
			err := rc.ch.Close()
			if err != nil {
				log.Error("failed to close rabbitMQ channel", "err", err.Error())
			}
			err = rc.conn.Close()
			if err != nil {
				log.Error("failed to close rabbitMQ channel", "err", err.Error())
			}
		}
	}
}

func (rc *rabbitMqClient) connect() error {
	conn, err := rc.dial(rc.url)
	if err != nil {
		return err
	}
	rc.conn = conn

	rc.connErrCh = make(chan *amqp.Error)
	rc.conn.NotifyClose(rc.connErrCh)

	err = rc.openChannel()
	if err != nil {
		return err
	}

	return nil
}

func (rc *rabbitMqClient) openChannel() error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}
	rc.ch = ch

	rc.chanErr = make(chan *amqp.Error)
	rc.ch.NotifyClose(rc.chanErr)

	return nil
}

func (rc *rabbitMqClient) reconnect() {
	for {
		time.Sleep(time.Millisecond * reconnectRetryMs)

		err := rc.connect()
		if err != nil {
			log.Debug("could not reconnect", "err", err.Error())
		} else {
			log.Debug("connection established after reconnect attempts")
			break
		}
	}
}

func (rc *rabbitMqClient) reopenChannel() {
	for {
		time.Sleep(time.Millisecond * reconnectRetryMs)

		err := rc.openChannel()
		if err != nil {
			log.Debug("could not re-open channel", "err", err.Error())
		} else {
			log.Debug("channel opened after reconnect attempts")
			break
		}
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (rc *rabbitMqClient) IsInterfaceNil() bool {
	return rc == nil
}
