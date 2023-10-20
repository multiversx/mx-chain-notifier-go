package factory

import (
	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/marshal"
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
	CheckDuplicates      bool
	Locker               process.LockService
	Publisher            process.Publisher
	StatusMetricsHandler common.StatusMetricsHandler
}

// CreateEventsHandler will create an events handler processor
func CreateEventsHandler(args ArgsEventsHandlerFactory) (process.EventsHandler, error) {
	argsEventsHandler := process.ArgsEventsHandler{
		Locker:               args.Locker,
		Publisher:            args.Publisher,
		StatusMetricsHandler: args.StatusMetricsHandler,
		CheckDuplicates:      args.CheckDuplicates,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	return eventsHandler, nil
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

func createPayloadHandler(marshaller marshal.Marshalizer, facade process.EventsFacadeHandler) (websocket.PayloadHandler, error) {
	dataPreProcessorArgs := preprocess.ArgsEventsPreProcessor{
		Marshaller: marshaller,
		Facade:     facade,
	}
	dataPreProcessors, err := createEventsDataPreProcessors(dataPreProcessorArgs)
	if err != nil {
		return nil, err
	}

	payloadHandler, err := process.NewPayloadHandler(dataPreProcessors)
	if err != nil {
		return nil, err
	}

	return payloadHandler, nil
}

func createEventsDataPreProcessors(dataPreProcessorArgs preprocess.ArgsEventsPreProcessor) (map[uint32]process.DataProcessor, error) {
	eventsProcessors := make(map[uint32]process.DataProcessor)

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
