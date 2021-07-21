package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/notifier-go"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
)

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
}

func CreateEventNotifier(args *EventNotifierFactoryArgs) (process.Indexer, error) {
	if err := checkInputArgs(args); err != nil {
		return nil, err
	}

	httpClient := client.NewHttpClient(client.HttpClientArgs{
		UseAuthorization: args.UseAuthorization,
		Username:         args.Username,
		Password:         args.Password,
		BaseUrl:          args.ProxyUrl,
	})

	pubkeyConv, err := pubkeyConverter.NewBech32PubkeyConverter(pubkeyLen)
	if err != nil {
		return nil, err
	}

	notifierArgs := notifier.EventNotifierArgs{
		HttpClient:      httpClient,
		Marshalizer:     args.Marshalizer,
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

	return nil
}
