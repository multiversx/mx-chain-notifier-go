package ws

import "github.com/gorilla/websocket"

type WebSocketsWrapper struct {
	conn *websocket.Conn
}
