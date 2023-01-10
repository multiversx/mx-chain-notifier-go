package ws

import (
	"fmt"
	"net/http"

	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/gorilla/websocket"
)

type wsUpgraderWrapper struct {
	upgrader *websocket.Upgrader
}

// NewWSUpgraderWrapper creates a websocket upgrader wrapper
func NewWSUpgraderWrapper(readBuffSize int, writeBuffSize int) (dispatcher.WSUpgrader, error) {
	if readBuffSize <= 0 {
		return nil, fmt.Errorf("invalid buffer size provided: %d", readBuffSize)
	}
	if writeBuffSize <= 0 {
		return nil, fmt.Errorf("invalid buffer size provided: %d", writeBuffSize)
	}

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  readBuffSize,
		WriteBufferSize: writeBuffSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &wsUpgraderWrapper{
		upgrader: upgrader,
	}, nil
}

// Upgrade upgrades the HTTP server connection to the websocket protocol
func (wuw *wsUpgraderWrapper) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (dispatcher.WSConnection, error) {
	return wuw.upgrader.Upgrade(w, r, responseHeader)
}
