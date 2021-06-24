package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type websocketDispatcher struct {
	id   uuid.UUID
	wg   sync.WaitGroup
	send chan []byte
	conn *websocket.Conn
	hub  dispatcher.Hub
}

func newWebsocketDispatcher(
	conn *websocket.Conn,
	hub dispatcher.Hub,
) *websocketDispatcher {
	return &websocketDispatcher{
		id:   uuid.New(),
		send: make(chan []byte, 256),
		conn: conn,
		hub:  hub,
	}
}

// Serve is the entry point used by a http server to serve the websocket upgrader
func Serve(hub dispatcher.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsDispatcher := newWebsocketDispatcher(conn, hub)
	wsDispatcher.hub.RegisterChan() <- wsDispatcher

	go wsDispatcher.writePump()
	go wsDispatcher.readPump()
}

// GetID returns the id corresponding to this dispatcher instance
func (wd *websocketDispatcher) GetID() uuid.UUID {
	return wd.id
}

// PushEvents receives an events slice and processes it before pushing to socket
func (wd *websocketDispatcher) PushEvents(events []data.Event) {
	eventBytes, err := json.Marshal(events)
	if err != nil {
		log.Println(err)
		return
	}
	wd.send <- eventBytes
}

// writePump listens on the send channel and pushes data on the socket stream
func (wd *websocketDispatcher) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := wd.conn.Close(); err != nil {
			log.Println(err)
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
				log.Println("channel: failed to set socket write limits", err.Error())
				return
			}

			if !ok {
				if err := wd.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println("failed to write close message", err.Error())
					return
				}
			}

			if err := nextWriterWrap(websocket.TextMessage, message); err != nil {
				log.Println("failed to write text message", err.Error())
				return
			}
		case <-ticker.C:
			if err := wd.setSocketWriteLimits(); err != nil {
				log.Println("ticker: failed to set socket write limits")
			}
			if err := wd.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ticker: failed to write ping message", err.Error())
				return
			}
		}
	}
}

// readPump listens for incoming events and reads the content from the socket stream
func (wd *websocketDispatcher) readPump() {
	defer func() {
		wd.hub.UnregisterChan() <- wd
		if err := wd.conn.Close(); err != nil {
			log.Println(err)
		}
		close(wd.send)
	}()

	if err := wd.setSocketReadLimits(); err != nil {
		log.Println("failed to set socket read limits", err.Error())

	}

	for {
		_, msg, innerErr := wd.conn.ReadMessage()
		if innerErr != nil {
			log.Println("failure reading socket", innerErr.Error())
			if websocket.IsUnexpectedCloseError(
				innerErr,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("unexpected socket close error", innerErr.Error())
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		wd.trySendSubscribeEvent(msg)
	}
}

func (wd *websocketDispatcher) trySendSubscribeEvent(data []byte) {
	var subscribeEvent dispatcher.SubscribeEvent
	err := json.Unmarshal(data, &subscribeEvent)
	if err != nil {
		log.Println(err)
		return
	}
	subscribeEvent.Caller = wd
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
