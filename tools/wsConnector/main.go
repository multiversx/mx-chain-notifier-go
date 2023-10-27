package main

import (
	"fmt"

	wsData "github.com/multiversx/mx-chain-communication-go/websocket/data"
	wsFactory "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/tools"
)

func main() {
	marshaller := &marshal.GogoProtoMarshalizer{}
	wsClient, err := newWSObsClient(marshaller)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	blockData, err := tools.NewBlockData(marshaller)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = wsClient.PushEventsRequest(blockData.OutportBlockV1())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = wsClient.RevertEventsRequest(blockData.RevertBlockV1())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = wsClient.FinalizedEventsRequest(blockData.FinalizedBlockV1())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// senderHost defines the actions that a host sender should do
type senderHost interface {
	Send(payload []byte, topic string) error
	Close() error
	IsInterfaceNil() bool
}

type wsObsClient struct {
	marshaller marshal.Marshalizer
	senderHost senderHost
}

// newWSObsClient will create a new instance of observer websocket client
func newWSObsClient(marshaller marshal.Marshalizer) (*wsObsClient, error) {
	var log = logger.GetOrCreate("hostdriver")

	port := 22111
	wsHost, err := wsFactory.CreateWebSocketHost(wsFactory.ArgsWebSocketHost{
		WebSocketConfig: wsData.WebSocketConfig{
			URL:                     "localhost:" + fmt.Sprintf("%d", port),
			WithAcknowledge:         true,
			Mode:                    "client",
			RetryDurationInSec:      5,
			BlockingAckOnError:      false,
			AcknowledgeTimeoutInSec: 60,
			Version:                 1,
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
