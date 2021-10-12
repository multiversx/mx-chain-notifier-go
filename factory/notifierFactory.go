package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
)

var log = logger.GetOrCreate("outport/eventNotifierFactory")

const (
	pubkeyLen = 32
)

type EventNotifierFactoryArgs struct {
	Enabled          bool
	UseAuthorization bool
	ProxyUrl         string
	Username         string
	Password         string
	Marshalizer      marshal.Marshalizer
	Hasher           hashing.Hasher
}

func CreateEventNotifier(args *EventNotifierFactoryArgs) (notifier.Indexer, error) {
	if err := checkInputArgs(args); err != nil {
		return nil, err
	}

	httpClient := client.NewHttpClient(client.HttpClientArgs{
		UseAuthorization: args.UseAuthorization,
		Username:         args.Username,
		Password:         args.Password,
		BaseUrl:          args.ProxyUrl,
	})

	pubkeyConv, err := pubkeyConverter.NewBech32PubkeyConverter(pubkeyLen, log)
	if err != nil {
		return nil, err
	}

	notifierArgs := notifier.EventNotifierArgs{
		HttpClient:      httpClient,
		Marshalizer:     args.Marshalizer,
		Hasher:          args.Hasher,
		PubKeyConverter: pubkeyConv,
	}

	eventNotifier, err := notifier.NewEventNotifier(notifierArgs)
	if err != nil {
		return nil, err
	}

	return eventNotifier, nil
}

func checkInputArgs(args *EventNotifierFactoryArgs) error {
	if check.IfNil(args.Marshalizer) {
		return core.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return core.ErrNilHasher
	}

	return nil
}
