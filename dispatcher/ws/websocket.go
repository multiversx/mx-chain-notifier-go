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

func NewWebsocketDispatcher(
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

// Serve is the entry point used by a http server to serve the ws upgrader
func Serve(hub dispatcher.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsDispatcher := NewWebsocketDispatcher(conn, hub)
	wsDispatcher.hub.RegisterChan() <- wsDispatcher

	go wsDispatcher.writePump()
	go wsDispatcher.readPump()
}

func (wd *websocketDispatcher) GetID() uuid.UUID {
	return wd.id
}

func (wd *websocketDispatcher) PushEvents(events []data.Event) {
	eventBytes, err := json.Marshal(events)
	if err != nil {
		log.Println(err)
		return
	}
	wd.send <- eventBytes
}

// writePump listens on the send channel and pushes data on the ws stream
func (wd *websocketDispatcher) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := wd.conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case msg, ok := <-wd.send:
			_ = wd.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = wd.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := wd.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}

			_, _ = w.Write(msg)

			if err = w.Close(); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			_ = wd.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wd.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// readPump listens for incoming events and reads the content
func (wd *websocketDispatcher) readPump() {
	defer func() {
		wd.hub.UnregisterChan() <- wd
		if err := wd.conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	wd.conn.SetReadLimit(maxMsgSize)
	_ = wd.conn.SetReadDeadline(time.Now().Add(pongWait))
	wd.conn.SetPongHandler(func(string) error {
		_ = wd.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := wd.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println(err)
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
