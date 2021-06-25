package factory

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ElrondNetwork/notifier-go"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
)

const (
	separator   = ":"
	hubCommon   = "common"
	dispatchAll = "dispatch:*"
	dispatchWs  = "websocket"
	dispatchGql = "graphql"
	defaultPort = "5000"
)

var availableHubDelegates = map[string]func() dispatcher.Hub{
	"default-custom": func() dispatcher.Hub {
		return nil
	},
}

type EventNotifierFactoryArgs struct {
	TcpPort      string
	Username     string
	Password     string
	HubType      string
	DispatchType string
	Marshalizer  marshal.Marshalizer
}

func CreateEventNotifier(args *EventNotifierFactoryArgs) (process.Indexer, error) {
	if err := checkInputArgs(args); err != nil {
		return nil, err
	}

	notifierHub, err := launchHubWithDispatchers(args)
	if err != nil {
		return nil, err
	}

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

func launchHubWithDispatchers(args *EventNotifierFactoryArgs) (dispatcher.Hub, error) {
	notifierHub, err := makeHub(args.HubType)
	if err != nil {
		return nil, err
	}

	go notifierHub.Run()

	var registerGQLHandler = func() {
		gqlHandler := gql.NewGraphQLServer(notifierHub)
		http.Handle("/query", gqlHandler)
		http.Handle("/subscription", gqlHandler)
	}

	var registerWsHandler = func() {
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			ws.Serve(notifierHub, w, r)
		})
		fmt.Println("a")
	}

	switch args.DispatchType {
	case dispatchAll:
		registerGQLHandler()
		registerWsHandler()
		break
	case dispatchWs:
		registerWsHandler()
		break
	case dispatchGql:
		registerGQLHandler()
		break
	default:
		return nil, ErrInvalidDispatchType
	}

	if args.TcpPort == "" {
		args.TcpPort = fmt.Sprintf(":%s", defaultPort)
	}
	if !strings.Contains(args.TcpPort, separator) {
		args.TcpPort = fmt.Sprintf(":%s", args.TcpPort)
	}

	go func(_port string) {
		if httpErr := http.ListenAndServe(_port, nil); httpErr != nil {
			panic(httpErr)
		}
	}(args.TcpPort)

	return notifierHub, nil
}

func makeHub(hubType string) (dispatcher.Hub, error) {
	if hubType == hubCommon {
		commonHub := hub.NewWsHub()
		return commonHub, nil
	}

	hubConfig := strings.Split(hubType, separator)
	return tryMakeCustomHubForID(hubConfig[1])
}

func tryMakeCustomHubForID(id string) (dispatcher.Hub, error) {
	if makeHubFunc, ok := availableHubDelegates[id]; ok {
		customHub := makeHubFunc()
		return customHub, nil
	}

	return nil, ErrInvalidCustomHubInput
}

func checkInputArgs(args *EventNotifierFactoryArgs) error {
	if check.IfNil(args.Marshalizer) {
		return core.ErrNilMarshalizer
	}

	return nil
}
