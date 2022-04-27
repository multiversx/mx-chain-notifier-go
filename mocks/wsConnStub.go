package mocks

import (
	"io"
	"time"
)

// WSConnStub implements dispatcher.WSConnection interface
type WSConnStub struct {
	NextWriterCalled       func(messageType int) (io.WriteCloser, error)
	WriteMessageCalled     func(messageType int, data []byte) error
	ReadMessageCalled      func() (messageType int, p []byte, err error)
	SetWriteDeadlineCalled func(t time.Time) error
	SetReadLimitCalled     func(limit int64)
	SetReadDeadlineCalled  func(t time.Time) error
	SetPongHandlerCalled   func(h func(appData string) error)
	CloseCalled            func() error
}

// NextWriter -
func (w *WSConnStub) NextWriter(messageType int) (io.WriteCloser, error) {
	if w.NextWriterCalled != nil {
		return w.NextWriterCalled(messageType)
	}

	return nil, nil
}

// WriteMessage -
func (w *WSConnStub) WriteMessage(messageType int, data []byte) error {
	if w.WriteMessageCalled != nil {
		return w.WriteMessageCalled(messageType, data)
	}

	return nil
}

// ReadMessage -
func (w *WSConnStub) ReadMessage() (messageType int, p []byte, err error) {
	if w.ReadMessageCalled != nil {
		return w.ReadMessageCalled()
	}

	return 0, nil, nil
}

// SetWriteDeadline -
func (w *WSConnStub) SetWriteDeadline(t time.Time) error {
	if w.SetWriteDeadlineCalled != nil {
		return w.SetWriteDeadlineCalled(t)
	}

	return nil
}

// SetReadLimit -
func (w *WSConnStub) SetReadLimit(limit int64) {
	if w.SetReadLimitCalled != nil {
		w.SetReadLimitCalled(limit)
	}
}

// SetReadDeadline -
func (w *WSConnStub) SetReadDeadline(t time.Time) error {
	if w.SetReadDeadlineCalled != nil {
		return w.SetReadDeadlineCalled(t)
	}

	return nil
}

// SetPongHandler -
func (w *WSConnStub) SetPongHandler(h func(appData string) error) {
	if w.SetPongHandlerCalled != nil {
		w.SetPongHandlerCalled(h)
	}
}

// Close -
func (w *WSConnStub) Close() error {
	if w.CloseCalled != nil {
		return w.CloseCalled()
	}

	return nil
}
