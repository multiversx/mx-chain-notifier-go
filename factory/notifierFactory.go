package factory

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/notifier-go"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
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
		BaseUrl:          args.ProxyUrl,
		Marshalizer:      args.Marshalizer,
	})

	notifierArgs := notifier.EventNotifierArgs{
		HttpClient:  httpClient,
		Marshalizer: args.Marshalizer,
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
