package factory

import (
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/disabled"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher/ws"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// CreateWSHandler creates websocket handler component based on api type
func CreateWSHandler(apiType string, hub dispatcher.Hub, marshaller marshal.Marshalizer) (dispatcher.WSHandler, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return &disabled.WSHandler{}, nil
	case common.WSAPIType:
		return createWSHandler(hub, marshaller)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createWSHandler(hub dispatcher.Hub, marshaller marshal.Marshalizer) (dispatcher.WSHandler, error) {
	upgrader, err := ws.NewWSUpgraderWrapper(readBufferSize, writeBufferSize)
	if err != nil {
		return nil, err
	}

	args := ws.ArgsWebSocketProcessor{
		Hub:        hub,
		Upgrader:   upgrader,
		Marshaller: marshaller,
	}
	return ws.NewWebSocketProcessor(args)
}
