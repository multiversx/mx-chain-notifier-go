package dispatcher

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

// EventDispatcher defines the behaviour of a event dispatcher component
type EventDispatcher interface {
	GetID() uuid.UUID
	PushEvents(events []data.Event)
}

// Hub defines the behaviour of a hub component which should be able to register
// and unregister dispatching events
type Hub interface {
	Run()
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	RegisterEvent(event EventDispatcher)
	UnregisterEvent(event EventDispatcher)
	Subscribe(event SubscribeEvent)
	Close() error
	IsInterfaceNil() bool
}
