package factory

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/notifier-go"
)

type EventNotifierFactoryArgs struct {
	TcpPort     string
	Marshalizer marshal.Marshalizer
}

func NewEventNotifier(args *EventNotifierFactoryArgs) (process.Indexer, error) {
	if err := checkInputArgs(args); err != nil {
		return nil, err
	}

	notifierHub := launchHubWithDispatchers(args.TcpPort)

	notifierArgs := notifier.EventNotifierArgs{
		Hub:         notifierHub,
		Marshalizer: args.Marshalizer,
	}

	eventNotifier, err := notifier.NewEventNotifier(notifierArgs)
	if err != nil {
		return nil, err
	}

	return eventNotifier, nil
}

func launchHubWithDispatchers(port string) dispatcher.Hub {
	if port == "" {
		port = ":5000"
	}
	if !strings.Contains(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	commonHub := hub.NewWsHub()
	go commonHub.Run()

	gqlSrv := gql.NewGraphQLServer(commonHub)

	http.Handle("/query", gqlSrv)
	http.Handle("/subscription", gqlSrv)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Serve(commonHub, w, r)
	})

	go func(_port string) {
		if err := http.ListenAndServe(_port, nil); err != nil {
			panic(err)
		}
	}(port)

	return commonHub
}

func checkInputArgs(args *EventNotifierFactoryArgs) error {
	if check.IfNil(args.Marshalizer) {
		return core.ErrNilMarshalizer
	}

	return nil
}
