package ws

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ArgsWebSocketHandler defines the argument needed to create a websocketHandler
type ArgsWebSocketHandler struct {
	Hub      dispatcher.Hub
	Upgrader dispatcher.WSUpgrader
}

type websocketHandler struct {
	hub      dispatcher.Hub
	upgrader dispatcher.WSUpgrader
}

// NewWebSocketHandler creates a new websocketHandler component
func NewWebSocketHandler(args ArgsWebSocketHandler) (*websocketHandler, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &websocketHandler{
		hub:      args.Hub,
		upgrader: args.Upgrader,
	}, nil
}

func checkArgs(args ArgsWebSocketHandler) error {
	if check.IfNil(args.Hub) {
		return ErrNilHubHandler
	}
	if check.IfNil(args.Upgrader) {
		return ErrNilWSUpgrader
	}

	return nil
}

// ServeHTTP is the entry point used by a http server to serve the websocket upgrader
func (wh *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed upgrading connection", "err", err.Error())
		return
	}
	wsDispatcher, err := newWebsocketDispatcher(conn, wh.hub)
	if err != nil {
		log.Error("failed creating a new websocket dispatcher", "err", err.Error())
		return
	}
	wsDispatcher.hub.RegisterEvent(wsDispatcher)

	go wsDispatcher.writePump()
	go wsDispatcher.readPump()
}

// // Serve is the entry point used by a http server to serve the websocket upgrader
// func Serve(hub dispatcher.Hub, w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Error("failed upgrading connection", "err", err.Error())
// 		return
// 	}
// 	wsDispatcher, err := NewWebsocketDispatcher(conn, hub)
// 	if err != nil {
// 		log.Error("failed creating a new websocket dispatcher", "err", err.Error())
// 		return
// 	}
// 	wsDispatcher.hub.RegisterEvent(wsDispatcher)

// 	go wsDispatcher.writePump()
// 	go wsDispatcher.readPump()
// }
