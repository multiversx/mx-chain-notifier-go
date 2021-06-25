package hub

import (
	"sync"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
)

type commonHub struct {
	rwMut           sync.RWMutex
	subscriptionMap *dispatcher.SubscriptionMap
	dispatchers     map[uuid.UUID]dispatcher.EventDispatcher
	register        chan dispatcher.EventDispatcher
	unregister      chan dispatcher.EventDispatcher
	broadcast       chan []data.Event
}

func NewCommonHub() *commonHub {
	return &commonHub{
		rwMut:           sync.RWMutex{},
		subscriptionMap: dispatcher.NewSubscriptionMap(),
		dispatchers:     make(map[uuid.UUID]dispatcher.EventDispatcher),
		register:        make(chan dispatcher.EventDispatcher),
		unregister:      make(chan dispatcher.EventDispatcher),
		broadcast:       make(chan []data.Event),
	}
}

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

func (wh *commonHub) Subscribe(event dispatcher.SubscribeEvent) {
	wh.subscriptionMap.MatchSubscribeEvent(event)
}

func (wh *commonHub) BroadcastChan() chan []data.Event {
	return wh.broadcast
}

func (wh *commonHub) RegisterChan() chan dispatcher.EventDispatcher {
	return wh.register
}

func (wh *commonHub) UnregisterChan() chan dispatcher.EventDispatcher {
	return wh.unregister
}

func (wh *commonHub) handleBroadcast(events []data.Event) {
	subscriptions := wh.subscriptionMap.Subscriptions()

	for _, subscription := range subscriptions[dispatcher.MatchAll] {
		wh.dispatchers[subscription.DispatcherID].PushEvents(events)
	}

	var filterableEvents []data.Event
	for _, event := range events {
		if _, ok := subscriptions[event.Address]; ok {
			filterableEvents = append(filterableEvents, event)
		}
	}

	dispatchersMap := make(map[uuid.UUID][]data.Event)
	mapEventToDispatcher := func(id uuid.UUID, e data.Event) {
		dispatchersMap[id] = append(dispatchersMap[id], e)
	}

	for _, event := range filterableEvents {
		subscriptionEntries := subscriptions[event.Address]
		for _, subEntry := range subscriptionEntries {
			switch subEntry.MatchLevel {
			case dispatcher.MatchAddress:
				mapEventToDispatcher(subEntry.DispatcherID, event)
				break
			case dispatcher.MatchIdentifier:
				if event.Identifier == subEntry.Identifier {
					mapEventToDispatcher(subEntry.DispatcherID, event)
				}
				break
			case dispatcher.MatchTopics:
				break
			}
		}
	}

	for id, eventValues := range dispatchersMap {
		wh.dispatchers[id].PushEvents(eventValues)
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
