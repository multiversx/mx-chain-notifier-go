package factory

import (
	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/marshal"
	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

var log = logger.GetOrCreate("factory")

const bech32PubkeyConverterType = "bech32"

// ArgsEventsHandlerFactory defines the arguments needed for events handler creation
type ArgsEventsHandlerFactory struct {
	APIConfig    config.ConnectorApiConfig
	Locker       process.LockService
	MqPublisher  process.Publisher
	HubPublisher process.Publisher
	APIType      string
}

// CreateEventsHandler will create an events handler processor
func CreateEventsHandler(args ArgsEventsHandlerFactory) (process.EventsHandler, error) {
	publisher, err := getPublisher(args.APIType, args.MqPublisher, args.HubPublisher)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Config:    args.APIConfig,
		Locker:    args.Locker,
		Publisher: publisher,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	return eventsHandler, nil
}

func getPublisher(
	apiType string,
	mqPublisher process.Publisher,
	hubPublisher process.Publisher,
) (process.Publisher, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return mqPublisher, nil
	case common.WSAPIType:
		return hubPublisher, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}

// CreateEventsInterceptor will create the events interceptor
func CreateEventsInterceptor(cfg config.GeneralConfig) (process.EventsInterceptor, error) {
	pubKeyConverter, err := getPubKeyConverter(cfg)
	if err != nil {
		return nil, err
	}

	argsEventsInterceptor := process.ArgsEventsInterceptor{
		PubKeyConverter: pubKeyConverter,
	}

	return process.NewEventsInterceptor(argsEventsInterceptor)
}

func getPubKeyConverter(cfg config.GeneralConfig) (core.PubkeyConverter, error) {
	switch cfg.AddressConverter.Type {
	case bech32PubkeyConverterType:
		return pubkeyConverter.NewBech32PubkeyConverter(cfg.AddressConverter.Length, cfg.AddressConverter.Prefix)
	default:
		return nil, common.ErrInvalidPubKeyConverterType
	}
}

// CreatePayloadHandler will create a new instance of payload handler
func CreatePayloadHandler(cfg config.Config, facade process.EventsFacadeHandler) (websocket.PayloadHandler, error) {
	headerMarshaller, err := createHeaderMarshaller(cfg)
	if err != nil {
		return nil, err
	}

	connectorMarshaller, err := createConnectorMarshaller(cfg)
	if err != nil {
		return nil, err
	}

	return createPayloadHandler(connectorMarshaller, headerMarshaller, facade)
}

func createConnectorMarshaller(cfg config.Config) (marshal.Marshalizer, error) {
	switch cfg.Flags.ConnectorType {
	case common.WSObsConnectorType:
		return marshalFactory.NewMarshalizer(cfg.WebSocketConnector.DataMarshallerType)
	case common.HTTPConnectorType:
		return &marshal.JsonMarshalizer{}, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createHeaderMarshaller(cfg config.Config) (marshal.Marshalizer, error) {
	switch cfg.Flags.ConnectorType {
	case common.WSObsConnectorType:
		return marshalFactory.NewMarshalizer(cfg.WebSocketConnector.DataMarshallerType)
	case common.HTTPConnectorType:
		return marshalFactory.NewMarshalizer(cfg.General.InternalMarshaller.Type)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createPayloadHandler(marshaller, headerMarshaller marshal.Marshalizer, facade process.EventsFacadeHandler) (websocket.PayloadHandler, error) {
	dataIndexerArgs := process.ArgsEventsDataPreProcessor{
		Marshaller: headerMarshaller,
		Facade:     facade,
	}
	dataPreProcessor, err := process.NewEventsDataPreProcessor(dataIndexerArgs)
	if err != nil {
		return nil, err
	}

	payloadHandler, err := process.NewPayloadHandler(marshaller, dataPreProcessor)
	if err != nil {
		return nil, err
	}

	return payloadHandler, nil
}
