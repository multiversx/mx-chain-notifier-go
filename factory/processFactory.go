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
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
)

var log = logger.GetOrCreate("factory")

const bech32PubkeyConverterType = "bech32"

// ArgsEventsHandlerFactory defines the arguments needed for events handler creation
type ArgsEventsHandlerFactory struct {
	APIConfig            config.ConnectorApiConfig
	Locker               process.LockService
	MqPublisher          process.Publisher
	HubPublisher         process.Publisher
	APIType              string
	StatusMetricsHandler common.StatusMetricsHandler
}

// CreateEventsHandler will create an events handler processor
func CreateEventsHandler(args ArgsEventsHandlerFactory) (process.EventsHandler, error) {
	publisher, err := getPublisher(args.APIType, args.MqPublisher, args.HubPublisher)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Config:               args.APIConfig,
		Locker:               args.Locker,
		Publisher:            publisher,
		StatusMetricsHandler: args.StatusMetricsHandler,
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
func CreatePayloadHandler(cfg config.Configs, facade process.EventsFacadeHandler) (websocket.PayloadHandler, error) {
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

func createConnectorMarshaller(cfg config.Configs) (marshal.Marshalizer, error) {
	switch cfg.Flags.ConnectorType {
	case common.WSObsConnectorType:
		return marshalFactory.NewMarshalizer(cfg.MainConfig.WebSocketConnector.DataMarshallerType)
	case common.HTTPConnectorType:
		return &marshal.JsonMarshalizer{}, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createHeaderMarshaller(cfg config.Configs) (marshal.Marshalizer, error) {
	switch cfg.Flags.ConnectorType {
	case common.WSObsConnectorType:
		return marshalFactory.NewMarshalizer(cfg.MainConfig.WebSocketConnector.DataMarshallerType)
	case common.HTTPConnectorType:
		return marshalFactory.NewMarshalizer(cfg.MainConfig.General.InternalMarshaller.Type)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createPayloadHandler(marshaller, headerMarshaller marshal.Marshalizer, facade process.EventsFacadeHandler) (websocket.PayloadHandler, error) {
	dataPreProcessorArgs := preprocess.ArgsEventsPreProcessor{
		Marshaller: headerMarshaller,
		Facade:     facade,
	}
	dataPreProcessors, err := createEventsDataPreProcessors(dataPreProcessorArgs)
	if err != nil {
		return nil, err
	}

	payloadHandler, err := process.NewPayloadHandler(marshaller, dataPreProcessors)
	if err != nil {
		return nil, err
	}

	return payloadHandler, nil
}

func createEventsDataPreProcessors(dataPreProcessorArgs preprocess.ArgsEventsPreProcessor) (map[common.PayloadVersion]process.DataProcessor, error) {
	eventsProcessors := make(map[common.PayloadVersion]process.DataProcessor)

	eventsProcessorV0, err := preprocess.NewEventsPreProcessorV0(dataPreProcessorArgs)
	if err != nil {
		return nil, err
	}
	eventsProcessors[common.PayloadV0] = eventsProcessorV0

	eventsProcessorV1, err := preprocess.NewEventsPreProcessorV1(dataPreProcessorArgs)
	if err != nil {
		return nil, err
	}
	eventsProcessors[common.PayloadV1] = eventsProcessorV1

	return eventsProcessors, nil
}
