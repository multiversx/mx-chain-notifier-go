package integrationTests

import (
	"errors"
	"fmt"
	"net"
	"time"

	wsData "github.com/multiversx/mx-chain-communication-go/websocket/data"
	wsFactory "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/factory"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
)

// CreateObserverConnector will create observer connector component
func CreateObserverConnector(facade shared.FacadeHandler, connType string, apiType string, payloadVersion uint32) (ObserverConnector, error) {
	marshaller := &marshal.JsonMarshalizer{}
	preProcessorArgs := preprocess.ArgsEventsPreProcessor{
		Marshaller: marshaller,
		Facade:     facade,
	}

	eventsProcessors := make(map[uint32]process.DataProcessor)
	dataPreProcessor, err := preprocess.NewEventsPreProcessorV1(preProcessorArgs)
	if err != nil {
		return nil, err
	}

	eventsProcessors[payloadVersion] = dataPreProcessor
	payloadHandler, err := process.NewPayloadHandler(eventsProcessors)
	if err != nil {
		return nil, err
	}

	switch connType {
	case common.HTTPConnectorType:
		return NewTestWebServer(facade, apiType, payloadHandler, payloadVersion), nil
	case common.WSObsConnectorType:
		return newTestWSServer(facade, marshaller)
	default:
		return nil, errors.New("invalid observer connector type")
	}
}

// newTestWSServer will create a new test ws server
func newTestWSServer(facade shared.FacadeHandler, marshaller marshal.Marshalizer) (ObserverConnector, error) {
	port := getRandomPort()
	conf := config.WebSocketConfig{
		Enabled:                 true,
		URL:                     "localhost:" + fmt.Sprintf("%d", port),
		WithAcknowledge:         true,
		Mode:                    "server",
		RetryDurationInSec:      5,
		BlockingAckOnError:      false,
		AcknowledgeTimeoutInSec: 60,
		DataMarshallerType:      "json",
	}

	_, err := factory.CreateWSObserverConnector(conf, facade)
	if err != nil {
		return nil, err
	}

	// wait for ws server to start
	time.Sleep(4 * time.Second)

	wsClient, err := newWSObsClient(marshaller, conf.URL)
	if err != nil {
		return nil, err
	}

	// wait for ws client to start
	time.Sleep(4 * time.Second)

	return wsClient, nil
}

func getRandomPort() int {
	// Listen on port 0 to get a free port
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = l.Close()
	}()

	// Get the port number that was assigned
	addr := l.Addr().(*net.TCPAddr)
	return addr.Port
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
func newWSObsClient(marshaller marshal.Marshalizer, url string) (*wsObsClient, error) {
	var log = logger.GetOrCreate("hostdriver")

	wsHost, err := wsFactory.CreateWebSocketHost(wsFactory.ArgsWebSocketHost{
		WebSocketConfig: wsData.WebSocketConfig{
			URL:                     url,
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
