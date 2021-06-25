package dispatcher

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

type EventDispatcher interface {
	GetID() uuid.UUID
	PushEvents(events []data.Event)
}

type Hub interface {
	BroadcastChan() chan []data.Event
	RegisterChan() chan EventDispatcher
	UnregisterChan() chan EventDispatcher
	Subscribe(event SubscribeEvent)
}
