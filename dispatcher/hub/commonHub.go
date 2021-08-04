package hub

import (
	"sync"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/google/uuid"
)

type commonHub struct {
	rwMut           sync.RWMutex
	filter          filters.EventFilter
	subscriptionMap *dispatcher.SubscriptionMap
	dispatchers     map[uuid.UUID]dispatcher.EventDispatcher
	register        chan dispatcher.EventDispatcher
	unregister      chan dispatcher.EventDispatcher
	broadcast       chan []data.Event
}

func NewCommonHub(eventFilter filters.EventFilter) *commonHub {
	return &commonHub{
		rwMut:           sync.RWMutex{},
		filter:          eventFilter,
		subscriptionMap: dispatcher.NewSubscriptionMap(),
		dispatchers:     make(map[uuid.UUID]dispatcher.EventDispatcher),
		register:        make(chan dispatcher.EventDispatcher),
		unregister:      make(chan dispatcher.EventDispatcher),
		broadcast:       make(chan []data.Event),
	}
}

// Run is launched as a goroutine and listens for events on the exposed channels
func (wh *commonHub) Run() {
	for {
		select {
		case events := <-wh.broadcast:
			wh.handleBroadcast(events)

		case dispatcherClient := <-wh.register:
			wh.registerDispatcher(dispatcherClient)

		case dispatcherClient := <-wh.unregister:
			wh.unregisterDispatcher(dispatcherClient)
		}
	}
}

// Subscribe is used by a dispatcher to send a dispatcher.SubscribeEvent
func (wh *commonHub) Subscribe(event dispatcher.SubscribeEvent) {
	wh.subscriptionMap.MatchSubscribeEvent(event)
}

// BroadcastChan returns a receive only channel on which events are pushed by producers
// Upon reading the channel, the hub notifies the registered dispatchers, if any
func (wh *commonHub) BroadcastChan() chan<- []data.Event {
	return wh.broadcast
}

// RegisterChan returns a receive only channel used to register dispatchers
func (wh *commonHub) RegisterChan() chan<- dispatcher.EventDispatcher {
	return wh.register
}

// UnregisterChan return a receive only channel used by a dispatcher to signal it has disconnected
func (wh *commonHub) UnregisterChan() chan<- dispatcher.EventDispatcher {
	return wh.unregister
}

func (wh *commonHub) handleBroadcast(events []data.Event) {
	subscriptions := wh.subscriptionMap.Subscriptions()

	for _, subscription := range subscriptions[dispatcher.MatchAll] {
		if _, ok := wh.dispatchers[subscription.DispatcherID]; ok {
			wh.dispatchers[subscription.DispatcherID].PushEvents(events)
		}
	}

	dispatchersMap := make(map[uuid.UUID][]data.Event)
	mapEventToDispatcher := func(id uuid.UUID, e data.Event) {
		dispatchersMap[id] = append(dispatchersMap[id], e)
	}

	for _, event := range events {
		subscriptionEntries := subscriptions[event.Address]
		for _, subEntry := range subscriptionEntries {
			if wh.filter.MatchEvent(subEntry, event) {
				mapEventToDispatcher(subEntry.DispatcherID, event)
			}
		}
	}

	wh.rwMut.RLock()
	defer wh.rwMut.RUnlock()
	for id, eventValues := range dispatchersMap {
		if d, ok := wh.dispatchers[id]; ok {
			d.PushEvents(eventValues)
		}
	}
}

func (wh *commonHub) registerDispatcher(d dispatcher.EventDispatcher) {
	if _, ok := wh.dispatchers[d.GetID()]; ok {
		return
	}
	wh.rwMut.Lock()
	defer wh.rwMut.Unlock()
	wh.dispatchers[d.GetID()] = d
}

func (wh *commonHub) unregisterDispatcher(d dispatcher.EventDispatcher) {
	wh.rwMut.Lock()
	defer wh.rwMut.Unlock()
	if _, ok := wh.dispatchers[d.GetID()]; ok {
		delete(wh.dispatchers, d.GetID())
	}
}
