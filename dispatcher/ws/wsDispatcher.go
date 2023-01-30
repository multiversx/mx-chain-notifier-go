package ws

import (
	"bytes"
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

var log = logger.GetOrCreate("websocket")

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 1024 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// argsWebSocketDispatcher defines the arguments needed for ws dispatcher
type argsWebSocketDispatcher struct {
	Hub  dispatcher.Hub
	Conn dispatcher.WSConnection
}

type websocketDispatcher struct {
	id   uuid.UUID
	wg   sync.WaitGroup
	send chan []byte
	conn dispatcher.WSConnection
	hub  dispatcher.Hub
}

// newWebSocketDispatcher createa a new ws dispatcher instance
func newWebSocketDispatcher(args argsWebSocketDispatcher) (*websocketDispatcher, error) {
	if check.IfNil(args.Hub) {
		return nil, ErrNilHubHandler
	}
	if args.Conn == nil {
		return nil, ErrNilWSConn
	}

	return &websocketDispatcher{
		id:   uuid.New(),
		send: make(chan []byte, 256),
		conn: args.Conn,
		hub:  args.Hub,
	}, nil
}

// GetID returns the id corresponding to this dispatcher instance
func (wd *websocketDispatcher) GetID() uuid.UUID {
	return wd.id
}

// PushEvents receives an events slice and processes it before pushing to socket
func (wd *websocketDispatcher) PushEvents(events []data.Event) {
	eventBytes, err := json.Marshal(events)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wsEvent := &data.WSEvent{
		Type: common.PushLogsAndEvents,
		Data: eventBytes,
	}
	wsEventBytes, err := json.Marshal(wsEvent)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wd.send <- wsEventBytes
}

// RevertEvent receives a reverted block event and process it before pushing to socket
func (wd *websocketDispatcher) RevertEvent(event data.RevertBlock) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}
	wsEvent := &data.WSEvent{
		Type: common.RevertBlockEvents,
		Data: eventBytes,
	}
	wsEventBytes, err := json.Marshal(wsEvent)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wd.send <- wsEventBytes
}

// FinalizedEvent receives a finalized block event and process it before pushing to socket
func (wd *websocketDispatcher) FinalizedEvent(event data.FinalizedBlock) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}
	wsEvent := &data.WSEvent{
		Type: common.FinalizedBlockEvents,
		Data: eventBytes,
	}
	wsEventBytes, err := json.Marshal(wsEvent)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wd.send <- wsEventBytes
}

// TxsEvent receives a block txs event and process it before pushing to socket
func (wd *websocketDispatcher) TxsEvent(event data.BlockTxs) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}
	wsEvent := &data.WSEvent{
		Type: common.BlockTxs,
		Data: eventBytes,
	}
	wsEventBytes, err := json.Marshal(wsEvent)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wd.send <- wsEventBytes
}

// BlockEvents receives block events with data and process it before pushing to socket
func (wd *websocketDispatcher) BlockEvents(event data.BlockEventsWithOrder) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}
	wsEvent := &data.WSEvent{
		Type: common.BlockEvents,
		Data: eventBytes,
	}
	wsEventBytes, err := json.Marshal(wsEvent)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wd.send <- wsEventBytes
}

// ScrsEvent receives a block scrs event and process it before pushing to socket
func (wd *websocketDispatcher) ScrsEvent(event data.BlockScrs) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}
	wsEvent := &data.WSEvent{
		Type: common.BlockScrs,
		Data: eventBytes,
	}
	wsEventBytes, err := json.Marshal(wsEvent)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	wd.send <- wsEventBytes
}

// writePump listens on the send-channel and pushes data on the socket stream
func (wd *websocketDispatcher) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := wd.conn.Close(); err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Debug("close attempt on closed connection", "err", err.Error())
				return
			}

			log.Error("failed to close socket", "err", err.Error())
		}
	}()

	nextWriterWrap := func(msgType int, data []byte) error {
		writer, err := wd.conn.NextWriter(msgType)
		if err != nil {
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
		return writer.Close()
	}

	for {
		select {
		case message, ok := <-wd.send:
			if err := wd.setSocketWriteLimits(); err != nil {
				log.Error("channel: failed to set socket write limits", "err", err.Error())
				return
			}

			if !ok {
				if err := wd.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Debug("failed to write close message", "err", err.Error())
					return
				}
			}

			if err := nextWriterWrap(websocket.TextMessage, message); err != nil {
				log.Error("failed to write text message", "err", err.Error())
				return
			}
		case <-ticker.C:
			if err := wd.setSocketWriteLimits(); err != nil {
				log.Error("ticker: failed to set socket write limits", "err", err.Error())
			}
			if err := wd.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Error("ticker: failed to write ping message", "err", err.Error())
				return
			}
		}
	}
}

// readPump listens for incoming events and reads the content from the socket stream
func (wd *websocketDispatcher) readPump() {
	defer func() {
		wd.hub.UnregisterEvent(wd)
		if err := wd.conn.Close(); err != nil {
			log.Error("failed to close socket on defer", "err", err.Error())
		}
		close(wd.send)
	}()

	if err := wd.setSocketReadLimits(); err != nil {
		log.Error("failed to set socket read limits", "err", err.Error())
	}

	for {
		_, msg, innerErr := wd.conn.ReadMessage()
		if innerErr != nil {
			log.Debug("failed reading socket", "err", innerErr.Error())

			if websocket.IsUnexpectedCloseError(
				innerErr,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Debug("received unexpected socket close", "err", innerErr.Error())
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		wd.trySendSubscribeEvent(msg)
	}
}

func (wd *websocketDispatcher) trySendSubscribeEvent(eventBytes []byte) {
	var subscribeEvent data.SubscribeEvent
	err := json.Unmarshal(eventBytes, &subscribeEvent)
	if err != nil {
		log.Error("failure unmarshalling subscribe event", "err", err.Error())
		return
	}
	subscribeEvent.DispatcherID = wd.id
	wd.hub.Subscribe(subscribeEvent)
}

func (wd *websocketDispatcher) setSocketWriteLimits() error {
	if err := wd.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return nil
}

func (wd *websocketDispatcher) setSocketReadLimits() error {
	wd.conn.SetReadLimit(maxMsgSize)
	if err := wd.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}
	wd.conn.SetPongHandler(func(string) error {
		return wd.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	return nil
}
