package factory

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
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
