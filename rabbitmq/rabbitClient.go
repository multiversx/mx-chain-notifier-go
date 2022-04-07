package rabbitmq

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	reconnectRetryMs = 500
)

type rabbitMqClient struct {
	url    string
	pubMut sync.Mutex

	conn *amqp.Connection
	ch   *amqp.Channel

	connErrCh chan *amqp.Error
	chanErr   chan *amqp.Error
}

// NewRabbitMQClient creates a new rabbitMQ client instance
func NewRabbitMQClient(url string) (*rabbitMqClient, error) {
	rc := &rabbitMqClient{
		url:    url,
		pubMut: sync.Mutex{},
	}

	err := rc.connect()
	if err != nil {
		return nil, err
	}

	return rc, nil
}

// Publish will publich an item on the rabbitMq channel
func (rc *rabbitMqClient) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	rc.pubMut.Lock()
	defer rc.pubMut.Unlock()

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

// ConnErrChan will return connection error channel
func (rc *rabbitMqClient) ConnErrChan() chan *amqp.Error {
	return rc.connErrCh
}

// CloseErrChan will return closing error channel
func (rc *rabbitMqClient) CloseErrChan() chan *amqp.Error {
	return rc.chanErr
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

// Reconnect will try to reconnect to rabbitmq
func (rc *rabbitMqClient) Reconnect() {
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

// ReopenChannel will try to reopen communication channel
func (rc *rabbitMqClient) ReopenChannel() {
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

// Close will close rabbitMq client connection
func (rc *rabbitMqClient) Close() {
	err := rc.ch.Close()
	if err != nil {
		log.Error("failed to close rabbitMQ channel", "err", err.Error())
	}
	err = rc.conn.Close()
	if err != nil {
		log.Error("failed to close rabbitMQ channel", "err", err.Error())
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (rc *rabbitMqClient) IsInterfaceNil() bool {
	return rc == nil
}
