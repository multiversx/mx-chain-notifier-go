package ws

import (
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// ArgsWebSocketProcessor defines the argument needed to create a websocketHandler
type ArgsWebSocketProcessor struct {
	Dispatcher dispatcher.Dispatcher
	Upgrader   dispatcher.WSUpgrader
	Marshaller marshal.Marshalizer
}

type websocketProcessor struct {
	dispatcher dispatcher.Dispatcher
	upgrader   dispatcher.WSUpgrader
	marshaller marshal.Marshalizer
}

// NewWebSocketProcessor creates a new websocketProcessor component
func NewWebSocketProcessor(args ArgsWebSocketProcessor) (*websocketProcessor, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &websocketProcessor{
		dispatcher: args.Dispatcher,
		upgrader:   args.Upgrader,
		marshaller: args.Marshaller,
	}, nil
}

func checkArgs(args ArgsWebSocketProcessor) error {
	if check.IfNil(args.Dispatcher) {
		return ErrNilDispatcher
	}
	if args.Upgrader == nil {
		return ErrNilWSUpgrader
	}
	if check.IfNil(args.Marshaller) {
		return common.ErrNilMarshaller
	}

	return nil
}

// ServeHTTP is the entry point used by a http server to serve the websocket upgrader
func (wh *websocketProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed upgrading connection", "err", err.Error())
		return
	}

	args := argsWebSocketDispatcher{
		Dispatcher: wh.dispatcher,
		Conn:       conn,
		Marshaller: wh.marshaller,
	}
	wsDispatcher, err := newWebSocketDispatcher(args)
	if err != nil {
		log.Error("failed creating a new websocket dispatcher", "err", err.Error())
		return
	}
	wsDispatcher.dispatcher.RegisterEvent(wsDispatcher)

	go wsDispatcher.writePump()
	go wsDispatcher.readPump()
}

// IsInterfaceNil returns true if there is no value under the interface
func (wh *websocketProcessor) IsInterfaceNil() bool {
	return wh == nil
}
