package dispatcher

import (
	"io"
	"net/http"
	"time"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

// EventDispatcher defines the behaviour of a event dispatcher component
type EventDispatcher interface {
	GetID() uuid.UUID
	PushEvents(events []data.Event)
	RevertEvent(event data.RevertBlock)
	FinalizedEvent(event data.FinalizedBlock)
	TxsEvent(event data.BlockTxs)
	ScrsEvent(event data.BlockScrs)
}

// Hub defines the behaviour of a hub component which should be able to register
// and unregister dispatching events
type Hub interface {
	Run()
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	BroadcastTxs(event data.BlockTxs)
	BroadcastScrs(event data.BlockScrs)
	RegisterEvent(event EventDispatcher)
	UnregisterEvent(event EventDispatcher)
	Subscribe(event data.SubscribeEvent)
	Close() error
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
