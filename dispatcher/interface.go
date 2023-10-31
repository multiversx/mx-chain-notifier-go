package dispatcher

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

// EventDispatcher defines the behaviour of a event dispatcher component
type EventDispatcher interface {
	GetID() uuid.UUID
	PushEvents(events []data.Event)
	RevertEvent(event data.RevertBlock)
	FinalizedEvent(event data.FinalizedBlock)
	TxsEvent(event data.BlockTxs)
	BlockEvents(event data.BlockEventsWithOrder)
	ScrsEvent(event data.BlockScrs)
}

// Hub defines the behaviour of a component which should be able to receive events
// and publish them to subscribers
type Hub interface {
	process.PublisherHandler
	Dispatcher
}

// Dispatcher defines the behaviour of a dispatcher component which should be able to register
// and unregister dispatching events
type Dispatcher interface {
	RegisterEvent(event EventDispatcher)
	UnregisterEvent(event EventDispatcher)
	Subscribe(event data.SubscribeEvent)
	IsInterfaceNil() bool
}

// WSHandler defines the behaviour of a websocket handler. It will serve http requests
type WSHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	IsInterfaceNil() bool
}

// WSConnection defines the behaviour of a websocket connection
type WSConnection interface {
	NextWriter(messageType int) (io.WriteCloser, error)
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	SetWriteDeadline(t time.Time) error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	Close() error
}

// WSUpgrader defines the behaviour of a websocket upgrader
type WSUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (WSConnection, error)
}

// SubscriptionMapperHandler defines the behaviour of a subscription mapper
type SubscriptionMapperHandler interface {
	MatchSubscribeEvent(event data.SubscribeEvent)
	RemoveSubscriptions(dispatcherID uuid.UUID)
	Subscriptions() []data.Subscription
	IsInterfaceNil() bool
}
