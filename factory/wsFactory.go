package factory

import (
	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-communication-go/websocket/data"
	factoryHost "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/disabled"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher/ws"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// CreateWSHandler creates websocket handler component based on api type
func CreateWSHandler(apiType string, hub dispatcher.Hub, marshaller marshal.Marshalizer) (dispatcher.WSHandler, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return &disabled.WSHandler{}, nil
	case common.WSPublisherType:
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

// CreateWSObserverConnector will create the web socket connector for observer node communication
func CreateWSObserverConnector(
	connectorType string,
	config config.WebSocketConfig,
	marshaller marshal.Marshalizer,
	payloadHandler websocket.PayloadHandler,
) (process.WSClient, error) {
	switch connectorType {
	case common.WSObsConnectorType:
		return createWsObsConnector(config, marshaller, payloadHandler)
	case common.HTTPConnectorType:
		return &disabled.WSHandler{}, nil
	default:
		return nil, common.ErrInvalidConnectorType
	}
}

func createWsObsConnector(
	config config.WebSocketConfig,
	marshaller marshal.Marshalizer,
	payloadHandler websocket.PayloadHandler,
) (process.WSClient, error) {
	host, err := createWsHost(config, marshaller)
	if err != nil {
		return nil, err
	}

	err = host.SetPayloadHandler(payloadHandler)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func createWsHost(wsConfig config.WebSocketConfig, wsMarshaller marshal.Marshalizer) (factoryHost.FullDuplexHost, error) {
	return factoryHost.CreateWebSocketHost(factoryHost.ArgsWebSocketHost{
		WebSocketConfig: data.WebSocketConfig{
			URL:                wsConfig.URL,
			WithAcknowledge:    wsConfig.WithAcknowledge,
			Mode:               wsConfig.Mode,
			RetryDurationInSec: int(wsConfig.RetryDurationInSec),
			BlockingAckOnError: wsConfig.BlockingAckOnError,
		},
		Marshaller: wsMarshaller,
		Log:        log,
	})
}
