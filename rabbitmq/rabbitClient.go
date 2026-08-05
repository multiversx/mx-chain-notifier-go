package rabbitmq

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	reconnectRetryMs = 500

	// notifyChanBufferSize is the buffer size used for the amqp notification channels.
	// It has to be at least 1: on connection/channel shutdown the amqp library sends the
	// close notification while holding its internal mutex, so an unbuffered channel with
	// no reader (no publish in progress) would block the whole shutdown.
	notifyChanBufferSize = 1

	// ExchangeDeclare constants
	isDurable  = true
	autoDelete = false
	isInternal = false
	noWait     = false
)

// exchangeDeclaration holds the data needed to (re)declare an exchange
type exchangeDeclaration struct {
	name string
	kind string
}

// connState holds a consistent snapshot of the currently used connection resources
type connState struct {
	ch        *amqp.Channel
	connErrCh chan *amqp.Error
	chanErr   chan *amqp.Error
	ackCh     chan uint64
	nackCh    chan uint64
}

type rabbitMqClient struct {
	url    string
	pubMut sync.Mutex

	// dialFunc and retryInterval are set on creation, they are overridden only in tests
	dialFunc      func(url string) (*amqp.Connection, error)
	retryInterval time.Duration

	connMut sync.RWMutex
	conn    *amqp.Connection
	ch      *amqp.Channel

	connErrCh chan *amqp.Error
	chanErr   chan *amqp.Error
	ackCh     chan uint64
	nackCh    chan uint64

	exchanges []exchangeDeclaration

	closeChan chan struct{}
	closeOnce sync.Once
}

// NewRabbitMQClient creates a new rabbitMQ client instance
func NewRabbitMQClient(url string) (*rabbitMqClient, error) {
	rc := &rabbitMqClient{
		url:           url,
		pubMut:        sync.Mutex{},
		dialFunc:      amqp.Dial,
		retryInterval: time.Millisecond * reconnectRetryMs,
		closeChan:     make(chan struct{}),
	}

	err := rc.tryConnect()
	if err != nil {
		return nil, err
	}

	return rc, nil
}

// ExchangeDeclare will declare an exchange. The declaration is saved, so that it can be
// applied again on every newly opened channel. This is needed because a reconnect may
// end up on a different rabbitmq instance (for example if the server IP behind the
// configured host has changed), which does not have the exchanges declared yet.
func (rc *rabbitMqClient) ExchangeDeclare(name, kind string) error {
	rc.connMut.RLock()
	ch := rc.ch
	rc.connMut.RUnlock()

	if ch == nil {
		return ErrClientClosed
	}

	err := declareExchange(ch, name, kind)
	if err != nil {
		return err
	}

	rc.saveExchangeDeclaration(name, kind)

	return nil
}

func declareExchange(ch *amqp.Channel, name, kind string) error {
	return ch.ExchangeDeclare(
		name,
		kind,
		isDurable,
		autoDelete,
		isInternal,
		noWait,
		nil,
	)
}

func (rc *rabbitMqClient) saveExchangeDeclaration(name, kind string) {
	rc.connMut.Lock()
	defer rc.connMut.Unlock()

	for idx, exchange := range rc.exchanges {
		if exchange.name == name {
			rc.exchanges[idx].kind = kind
			return
		}
	}

	rc.exchanges = append(rc.exchanges, exchangeDeclaration{name: name, kind: kind})
}

func (rc *rabbitMqClient) currentState() connState {
	rc.connMut.RLock()
	defer rc.connMut.RUnlock()

	return connState{
		ch:        rc.ch,
		connErrCh: rc.connErrCh,
		chanErr:   rc.chanErr,
		ackCh:     rc.ackCh,
		nackCh:    rc.nackCh,
	}
}

// Publish will publich an item on the rabbitMq channel
func (rc *rabbitMqClient) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	rc.pubMut.Lock()
	defer rc.pubMut.Unlock()

	// In order to avoid losing any event, check rabbitmq ack event for the
	// published message. If not-acknowledged, check if there is a connection or
	// channel issue, and after that is solved try again.  This was done to
	// make sure no event is lost, for example, if rabbitmq connection is not
	// closing gracefully (port disabled from firewall), it may happen that the
	// main loop will not catch the conn err event, and it will still try to
	// publish the message.
	for {
		if rc.isClosed() {
			return ErrClientClosed
		}

		state := rc.currentState()
		if state.ch == nil {
			return ErrClientClosed
		}

		err := state.ch.Publish(
			exchange,
			key,
			mandatory,
			immediate,
			msg,
		)
		if err != nil {
			// no confirmation will ever be received for a failed publish, so the
			// connection/channel has to be restored before trying again, otherwise
			// this would spin without ever making progress
			log.Error("failed to publish message", "exchange", exchange, "error", err.Error())
			rc.recoverConnection()
			continue
		}

		select {
		case deliveryTag, ok := <-state.ackCh:
			if !ok {
				// the confirmations channel is closed together with the amqp channel,
				// so this is a failed publish, not an acknowledgement
				log.Debug("Publish: confirmations channel closed, will retry to publish message", "exchange", exchange)
				rc.recoverConnection()
				continue
			}

			log.Debug("Publish: published message ack", "deliveryTag", deliveryTag)
			return nil
		case deliveryTag, ok := <-state.nackCh:
			if !ok {
				log.Debug("Publish: confirmations channel closed, will retry to publish message", "exchange", exchange)
				rc.recoverConnection()
				continue
			}

			log.Debug("Publish: published message nack, will retry to publish message", "deliveryTag", deliveryTag)
		case amqpErr, ok := <-state.connErrCh:
			logAmqpFailure("rabbitMQ connection failure", amqpErr, ok)
			rc.Reconnect()
		case amqpErr, ok := <-state.chanErr:
			logAmqpFailure("rabbitMQ channel failure", amqpErr, ok)

			// a connection failure is broadcast on the channel notification as well, so
			// the recovery has to check what was actually lost
			rc.recoverConnection()
		}
	}
}

func logAmqpFailure(message string, amqpErr *amqp.Error, okChanRead bool) {
	if !okChanRead || amqpErr == nil {
		log.Error(message, "err", "notification channel closed")
		return
	}

	log.Error(message, "err", amqpErr.Error())
}

// ConnErrChan will return connection error channel
func (rc *rabbitMqClient) ConnErrChan() chan *amqp.Error {
	rc.connMut.RLock()
	defer rc.connMut.RUnlock()

	return rc.connErrCh
}

// CloseErrChan will return closing error channel
func (rc *rabbitMqClient) CloseErrChan() chan *amqp.Error {
	rc.connMut.RLock()
	defer rc.connMut.RUnlock()

	return rc.chanErr
}

// connect will create a new connection and a new channel, closing the previous ones if
// still around. It has to be called while holding the connMut write lock.
func (rc *rabbitMqClient) connect() error {
	rc.closeConnectionAndChannel()

	conn, err := rc.dialFunc(rc.url)
	if err != nil {
		return err
	}

	if rc.isClosed() {
		// the client was closed while dialing, do not leak the new connection
		closeConnection(conn)
		return ErrClientClosed
	}

	rc.conn = conn

	rc.connErrCh = make(chan *amqp.Error, notifyChanBufferSize)
	rc.conn.NotifyClose(rc.connErrCh)

	return rc.openChannel()
}

// openChannel will open a new channel on the existing connection, closing the previous
// one if still around. It has to be called while holding the connMut write lock.
func (rc *rabbitMqClient) openChannel() error {
	rc.closeChannel()

	if rc.conn == nil {
		return ErrClientClosed
	}

	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}
	rc.ch = ch

	rc.chanErr = make(chan *amqp.Error, notifyChanBufferSize)
	rc.ch.NotifyClose(rc.chanErr)
	rc.ackCh, rc.nackCh = rc.ch.NotifyConfirm(
		make(chan uint64, notifyChanBufferSize),
		make(chan uint64, notifyChanBufferSize),
	)

	err = rc.ch.Confirm(false)
	if err != nil {
		return err
	}

	return rc.redeclareExchanges()
}

// redeclareExchanges has to be called while holding the connMut write lock
func (rc *rabbitMqClient) redeclareExchanges() error {
	for _, exchange := range rc.exchanges {
		err := declareExchange(rc.ch, exchange.name, exchange.kind)
		if err != nil {
			return err
		}

		log.Debug("re-declared rabbitMQ exchange", "name", exchange.name, "type", exchange.kind)
	}

	return nil
}

func (rc *rabbitMqClient) tryConnect() error {
	rc.connMut.Lock()
	defer rc.connMut.Unlock()

	return rc.connect()
}

func (rc *rabbitMqClient) tryOpenChannel() error {
	rc.connMut.Lock()
	defer rc.connMut.Unlock()

	return rc.openChannel()
}

// recover will restore whatever was lost: if the connection is down, re-opening only the
// channel would never succeed, since a channel cannot be opened on a closed connection
func (rc *rabbitMqClient) recoverConnection() {
	if rc.connectionIsDown() {
		rc.Reconnect()
		return
	}

	rc.ReopenChannel()
}

func (rc *rabbitMqClient) connectionIsDown() bool {
	rc.connMut.RLock()
	defer rc.connMut.RUnlock()

	return rc.conn == nil || rc.conn.IsClosed()
}

// Reconnect will try to reconnect to rabbitmq
func (rc *rabbitMqClient) Reconnect() {
	for {
		if rc.waitBeforeRetry() {
			log.Debug("client closed, stopping reconnect attempts")
			return
		}

		err := rc.tryConnect()
		if err != nil {
			log.Debug("could not reconnect", "err", err.Error())
			continue
		}

		log.Info("connection established after reconnect attempts")
		return
	}
}

// ReopenChannel will try to reopen communication channel
func (rc *rabbitMqClient) ReopenChannel() {
	for {
		if rc.waitBeforeRetry() {
			log.Debug("client closed, stopping channel reopen attempts")
			return
		}

		if rc.connectionIsDown() {
			log.Debug("connection is down, will reconnect before re-opening the channel")
			rc.Reconnect()
			return
		}

		err := rc.tryOpenChannel()
		if err != nil {
			log.Debug("could not re-open channel", "err", err.Error())
			continue
		}

		log.Info("channel opened after reconnect attempts")
		return
	}
}

// waitBeforeRetry will wait for the retry interval to pass, returning true if the client
// has been closed in the meantime, case in which the retry loop has to be stopped
func (rc *rabbitMqClient) waitBeforeRetry() bool {
	timer := time.NewTimer(rc.retryInterval)
	defer timer.Stop()

	select {
	case <-rc.closeChan:
		return true
	case <-timer.C:
		return false
	}
}

func (rc *rabbitMqClient) isClosed() bool {
	select {
	case <-rc.closeChan:
		return true
	default:
		return false
	}
}

// Close will close rabbitMq client connection
func (rc *rabbitMqClient) Close() {
	rc.closeOnce.Do(func() {
		close(rc.closeChan)
	})

	rc.connMut.Lock()
	defer rc.connMut.Unlock()

	rc.closeConnectionAndChannel()
}

// closeConnectionAndChannel has to be called while holding the connMut write lock
func (rc *rabbitMqClient) closeConnectionAndChannel() {
	rc.closeChannel()

	if rc.conn == nil {
		return
	}

	closeConnection(rc.conn)
	rc.conn = nil
}

// closeChannel has to be called while holding the connMut write lock
func (rc *rabbitMqClient) closeChannel() {
	if rc.ch == nil {
		return
	}

	err := rc.ch.Close()
	if err != nil && err != amqp.ErrClosed {
		log.Debug("failed to close rabbitMQ channel", "err", err.Error())
	}
	rc.ch = nil
}

func closeConnection(conn *amqp.Connection) {
	err := conn.Close()
	if err != nil && err != amqp.ErrClosed {
		log.Debug("failed to close rabbitMQ connection", "err", err.Error())
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (rc *rabbitMqClient) IsInterfaceNil() bool {
	return rc == nil
}
