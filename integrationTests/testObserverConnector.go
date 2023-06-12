package integrationTests

import (
	"errors"
	"fmt"
	"time"

	wsData "github.com/multiversx/mx-chain-communication-go/websocket/data"
	wsFactory "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/factory"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

func CreateObserverConnector(facade shared.FacadeHandler, connType string, apiType string) (ObserverConnector, process.WSClient, error) {
	marshaller := &marshal.JsonMarshalizer{}
	switch connType {
	case "http":
		return NewTestWebServer(facade, apiType), nil, nil
	case "ws":
		return NewTestWSServer(facade, marshaller)
	default:
		return nil, nil, errors.New("invalid observer connector type")
	}
}

func NewTestWSServer(facade process.EventsFacadeHandler, marshaller marshal.Marshalizer) (ObserverConnector, process.WSClient, error) {
	conf := config.WebSocketConfig{
		URL:                "localhost:22111",
		WithAcknowledge:    true,
		Mode:               "server",
		RetryDurationInSec: 5,
		BlockingAckOnError: false,
	}

	wsServer, err := factory.CreateWSObserverConnector(conf, marshaller, facade)
	if err != nil {
		return nil, nil, err
	}

	// wait for ws server to start
	time.Sleep(2 * time.Second)

	wsClient, err := NewWSObsClient(marshaller)
	if err != nil {
		return nil, nil, err
	}

	return wsClient, wsServer, nil
}

type wsObsConnector struct {
	wsClient *wsObsClient
}

// SenderHost defines the actions that a host sender should do
type SenderHost interface {
	Send(payload []byte, topic string) error
	Close() error
	IsInterfaceNil() bool
}

type wsObsClient struct {
	marshaller marshal.Marshalizer
	senderHost SenderHost
}

// NewWSObsClient will create a new instance of observer websocket client
func NewWSObsClient(marshaller marshal.Marshalizer) (*wsObsClient, error) {
	var log = logger.GetOrCreate("hostdriver")

	wsHost, err := wsFactory.CreateWebSocketHost(wsFactory.ArgsWebSocketHost{
		WebSocketConfig: wsData.WebSocketConfig{
			URL:                "localhost:22111",
			WithAcknowledge:    true,
			Mode:               "client",
			RetryDurationInSec: 5,
			BlockingAckOnError: false,
		},
		Marshaller: marshaller,
		Log:        log,
	})
	if err != nil {
		return nil, err
	}

	return &wsObsClient{
		marshaller: marshaller,
		senderHost: wsHost,
	}, nil
}

// SaveBlock will handle the saving of block
func (o *wsObsClient) PushEventsRequest(outportBlock *outport.OutportBlock) error {
	return o.handleAction(outportBlock, outport.TopicSaveBlock)
}

// RevertIndexedBlock will handle the action of reverting the indexed block
func (o *wsObsClient) RevertEventsRequest(blockData *outport.BlockData) error {
	return o.handleAction(blockData, outport.TopicRevertIndexedBlock)
}

// FinalizedBlock will handle the finalized block
func (o *wsObsClient) FinalizedEventsRequest(finalizedBlock *outport.FinalizedBlock) error {
	return o.handleAction(finalizedBlock, outport.TopicFinalizedBlock)
}

func (o *wsObsClient) handleAction(args interface{}, topic string) error {
	marshalledPayload, err := o.marshaller.Marshal(args)
	if err != nil {
		return fmt.Errorf("%w while marshaling block for topic %s", err, topic)
	}

	err = o.senderHost.Send(marshalledPayload, topic)
	if err != nil {
		return fmt.Errorf("%w while sending data on route for topic %s", err, topic)
	}

	return nil
}

// Close will close sender host connection
func (o *wsObsClient) Close() error {
	return o.senderHost.Close()
}
