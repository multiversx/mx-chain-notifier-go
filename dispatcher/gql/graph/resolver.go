package graph

import (
	"sync"

	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/model"
	"github.com/google/uuid"
)

type Resolver struct {
	mut          sync.Mutex
	hub          dispatcher.Hub
	chanMappings map[uuid.UUID]chan<- []*model.Event
}

func NewResolver(hub dispatcher.Hub) *Resolver {
	return &Resolver{
		hub:          hub,
		mut:          sync.Mutex{},
		chanMappings: make(map[uuid.UUID]chan<- []*model.Event),
	}
}

func (r *Resolver) MapDispatchChan(dispatchID uuid.UUID, dispatchChan chan<- []*model.Event) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.chanMappings[dispatchID] = dispatchChan
}

func (r *Resolver) RemoveDispatchChan(dispatchID uuid.UUID) {
	if _, ok := r.chanMappings[dispatchID]; ok {
		return
	}
	r.mut.Lock()
	defer r.mut.Unlock()
	delete(r.chanMappings, dispatchID)
}
