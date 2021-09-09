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
	Run()
	BroadcastChan() chan<- []data.Event
	RegisterChan() chan<- EventDispatcher
	UnregisterChan() chan<- EventDispatcher
	Subscribe(event SubscribeEvent)
}

type H struct {
	Hub

}

func (h *H) BroadcastChan() chan<-[]data.Event {
	h := &H{}

	return nil
}

