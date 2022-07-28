package mocks

import "github.com/streadway/amqp"

// RabbitClientStub -
type RabbitClientStub struct {
	PublishCalled         func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	ExchangeDeclareCalled func(name, kind string) error
	ConnErrChanCalled     func() chan *amqp.Error
	CloseErrChanCalled    func() chan *amqp.Error
	ReconnectCalled       func()
	ReopenChannelCalled   func()
	CloseCalled           func()
}

// Publish -
func (rc *RabbitClientStub) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if rc.PublishCalled != nil {
		return rc.PublishCalled(exchange, key, mandatory, immediate, msg)
	}
	return nil
}

// ExchangeDeclare -
func (rc *RabbitClientStub) ExchangeDeclare(name, kind string) error {
	if rc.ExchangeDeclareCalled != nil {
		return rc.ExchangeDeclareCalled(name, kind)
	}
	return nil
}

// ConnErrChan -
func (rc *RabbitClientStub) ConnErrChan() chan *amqp.Error {
	if rc.ConnErrChanCalled != nil {
		return rc.ConnErrChanCalled()
	}
	return nil
}

// CloseErrChan -
func (rc *RabbitClientStub) CloseErrChan() chan *amqp.Error {
	if rc.CloseErrChanCalled != nil {
		return rc.CloseErrChanCalled()
	}
	return nil
}

// Reconnect -
func (rc *RabbitClientStub) Reconnect() {
	if rc.ReconnectCalled != nil {
		rc.ReconnectCalled()
	}
}

// ReopenChannel -
func (rc *RabbitClientStub) ReopenChannel() {
	if rc.ReopenChannelCalled != nil {
		rc.ReopenChannelCalled()
	}
}

// Close -
func (rc *RabbitClientStub) Close() {
	if rc.CloseCalled != nil {
		rc.CloseCalled()
	}
}

// IsInterfaceNil -
func (rc *RabbitClientStub) IsInterfaceNil() bool {
	return rc == nil
}
