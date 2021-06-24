package hub

import (
	"sync"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
)

type wsHub struct {
	rwMut           sync.RWMutex
	subscriptionMap *dispatcher.SubscriptionMap
	dispatchers     map[uuid.UUID]dispatcher.EventDispatcher
	register        chan dispatcher.EventDispatcher
	unregister      chan dispatcher.EventDispatcher
	broadcast       chan []data.Event
}

func NewWsHub() *wsHub {
	return &wsHub{
		rwMut:           sync.RWMutex{},
		subscriptionMap: dispatcher.NewSubscriptionMap(),
		dispatchers:     make(map[uuid.UUID]dispatcher.EventDispatcher),
		register:        make(chan dispatcher.EventDispatcher),
		unregister:      make(chan dispatcher.EventDispatcher),
		broadcast:       make(chan []data.Event),
	}
}

func (wh *wsHub) Run() {
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

func (wh *wsHub) Subscribe(event dispatcher.SubscribeEvent) {
	wh.subscriptionMap.MatchSubscribeEvent(event)
}

func (wh *wsHub) BroadcastChan() chan<- []data.Event {
	return wh.broadcast
}

func (wh *wsHub) RegisterChan() chan<- dispatcher.EventDispatcher {
	return wh.register
}

func (wh *wsHub) UnregisterChan() chan<- dispatcher.EventDispatcher {
	return wh.unregister
}

func (wh *wsHub) handleBroadcast(events []data.Event) {
	subscriptions := wh.subscriptionMap.Subscriptions()

	for _, subscription := range subscriptions[dispatcher.MatchAll] {
		subscription.Subscriber.PushEvents(events)
	}

	var filterableEvents []data.Event
	for _, event := range events {
		if _, ok := subscriptions[event.Address]; ok {
			filterableEvents = append(filterableEvents, event)
		}
	}

	dispatchersMap := make(map[dispatcher.EventDispatcher][]data.Event)
	mapEventToDispatcher := func(d dispatcher.EventDispatcher, e data.Event) {
		dispatchersMap[d] = append(dispatchersMap[d], e)
	}

	for _, event := range filterableEvents {
		subscriptionEntries := subscriptions[event.Address]
		for _, subEntry := range subscriptionEntries {
			switch subEntry.MatchLevel {
			case dispatcher.MatchAddress:
				mapEventToDispatcher(subEntry.Subscriber, event)
				break
			case dispatcher.MatchIdentifier:
				if event.Identifier == subEntry.Identifier {
					mapEventToDispatcher(subEntry.Subscriber, event)
				}
				break
			case dispatcher.MatchTopics:
				break
			}
		}
	}

	for dispatch, eventValues := range dispatchersMap {
		dispatch.PushEvents(eventValues)
	}
}

func (wh *wsHub) registerDispatcher(d dispatcher.EventDispatcher) {
	if _, ok := wh.dispatchers[d.GetID()]; ok {
		return
	}
	wh.rwMut.Lock()
	defer wh.rwMut.Unlock()
	wh.dispatchers[d.GetID()] = d
}

func (wh *wsHub) unregisterDispatcher(d dispatcher.EventDispatcher) {
	wh.rwMut.Lock()
	defer wh.rwMut.Unlock()
	if _, ok := wh.dispatchers[d.GetID()]; ok {
		delete(wh.dispatchers, d.GetID())
	}
}
