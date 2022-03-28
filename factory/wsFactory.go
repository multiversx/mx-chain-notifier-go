package factory

import (
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// CreateWSHandler creates websocket handler component based on api type
func CreateWSHandler(apiType string, hub dispatcher.Hub) (dispatcher.WSHandler, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return &disabled.WSHandler{}, nil
	case common.WSAPIType:
		return createWSHandler(hub)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createWSHandler(hub dispatcher.Hub) (dispatcher.WSHandler, error) {
	upgrader, err := ws.NewWSUpgraderWrapper(readBufferSize, writeBufferSize)
	if err != nil {
		return nil, err
	}

	args := ws.ArgsWebSocketHandler{
		Hub:      hub,
		Upgrader: upgrader,
	}
	wsHandler, err := ws.NewWebSocketHandler(args)
	if err != nil {
		return nil, err
	}

	return wsHandler, nil
}
