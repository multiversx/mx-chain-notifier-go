package rabbitmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

type RabbitClient interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type rabbitClientWrapper struct {
	url string

	ctx context.Context

	conn *amqp.Connection
	ch   *amqp.Channel

	connErrCh chan *amqp.Error
	chanErr   chan *amqp.Error
}

func NewRabbitClientWrapper(ctx context.Context, url string) (*rabbitClientWrapper, error) {
	rc := &rabbitClientWrapper{
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

func (rc *rabbitClientWrapper) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return rc.ch.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

func (rc *rabbitClientWrapper) connListener() {
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

func (rc *rabbitClientWrapper) connect() error {
	conn, err := amqp.Dial(rc.url)
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

func (rc *rabbitClientWrapper) openChannel() error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}
	rc.ch = ch

	rc.chanErr = make(chan *amqp.Error)
	rc.ch.NotifyClose(rc.chanErr)

	return nil
}

func (rc *rabbitClientWrapper) reconnect() {
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

func (rc *rabbitClientWrapper) reopenChannel() {
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
