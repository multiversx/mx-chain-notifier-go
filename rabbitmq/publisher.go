package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/streadway/amqp"
)

const (
	reconnectRetryMs = 500
)

type rabbitMqPublisher struct {
	dispatcher.Hub

	broadcast chan []data.Event
	connErrCh chan *amqp.Error

	conn   *amqp.Connection
	ch     *amqp.Channel
	amqurl string

	ctx context.Context
}

func NewRabbitMqPublisher(ctx context.Context) *rabbitMqPublisher {
	rp := &rabbitMqPublisher{
		broadcast: make(chan []data.Event),
		amqurl:    "amqp://guest:guest@localhost:5672",
		ctx:       ctx,
	}

	err := rp.connect()
	if err != nil {
		fmt.Println("connect err", err.Error())
	}

	return rp
}

// Run is launched as a goroutine and listens for events on the exposed channels
func (rp *rabbitMqPublisher) Run() {
	for {
		select {
		case events := <-rp.broadcast:
			rp.publishToExchanges(events)
		case err := <-rp.connErrCh:
			if err != nil {
				fmt.Println("conn err", err.Error())
				rp.reconnect()
			}
		case <-rp.ctx.Done():
			fmt.Println("ctx done. closing")
			err := rp.conn.Close()
			if err != nil {
				fmt.Println("failed to close conn", err.Error())
			}
		}
	}
}

// BroadcastChan returns a receive-only channel on which events are pushed by producers
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastChan() chan<- []data.Event {
	return rp.broadcast
}

func (rp *rabbitMqPublisher) publishToExchanges(events []data.Event) {

}

func (rp *rabbitMqPublisher) connect() error {
	conn, err := amqp.Dial(rp.amqurl)
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
			fmt.Println("err while trying to reconnect", err.Error())
			continue
		}

		fmt.Println("connection established after reconnect")
		break
	}
}
