package hub

import (
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/google/uuid"
)

var log = logger.GetOrCreate("hub")

type commonHub struct {
	rwMut              sync.RWMutex
	filter             filters.EventFilter
	subscriptionMapper *dispatcher.SubscriptionMapper
	dispatchers        map[uuid.UUID]dispatcher.EventDispatcher
	register           chan dispatcher.EventDispatcher
	unregister         chan dispatcher.EventDispatcher
	broadcast          chan data.BlockEvents
	broadcastRevert    chan data.RevertBlock
}

// NewCommonHub creates a new commonHub instance
func NewCommonHub(eventFilter filters.EventFilter) *commonHub {
	return &commonHub{
		rwMut:              sync.RWMutex{},
		filter:             eventFilter,
		subscriptionMapper: dispatcher.NewSubscriptionMapper(),
		dispatchers:        make(map[uuid.UUID]dispatcher.EventDispatcher),
		register:           make(chan dispatcher.EventDispatcher),
		unregister:         make(chan dispatcher.EventDispatcher),
		broadcast:          make(chan data.BlockEvents),
		broadcastRevert:    make(chan data.RevertBlock),
	}
}

// Run is launched as a goroutine and listens for events on the exposed channels
func (wh *commonHub) Run() {
	for {
		select {
		case events := <-wh.broadcast:
			wh.handleBroadcast(events)

		case revertEvent := <-wh.broadcastRevert:
			wh.handleRevertBroadcast(revertEvent)

		case dispatcherClient := <-wh.register:
			wh.registerDispatcher(dispatcherClient)

		case dispatcherClient := <-wh.unregister:
			wh.unregisterDispatcher(dispatcherClient)
		}
	}
}

// Subscribe is used by a dispatcher to send a dispatcher.SubscribeEvent
func (wh *commonHub) Subscribe(event dispatcher.SubscribeEvent) {
	wh.subscriptionMapper.MatchSubscribeEvent(event)
}

// BroadcastChan returns a receive-only channel on which events are pushed by producers
// Upon reading the channel, the hub notifies the registered dispatchers, if any
func (wh *commonHub) BroadcastChan() chan<- data.BlockEvents {
	return wh.broadcast
}

// BroadcastRevertChan returns a receive-only channel on which revert events are pushed
// Upon reading the channel, the hub notifies the registered dispatchers, if any
func (wh *commonHub) BroadcastRevertChan() chan<- data.RevertBlock {
	return wh.broadcastRevert
}

// RegisterChan returns a receive-only channel used to register dispatchers
func (wh *commonHub) RegisterChan() chan<- dispatcher.EventDispatcher {
	return wh.register
}

// UnregisterChan return a receive-only channel used by a dispatcher to signal it has disconnected
func (wh *commonHub) UnregisterChan() chan<- dispatcher.EventDispatcher {
	return wh.unregister
}

func (wh *commonHub) handleBroadcast(blockEvents data.BlockEvents) {
	subscriptions := wh.subscriptionMapper.Subscriptions()

	dispatchersMap := make(map[uuid.UUID][]data.Event)
	mapEventToDispatcher := func(id uuid.UUID, e data.Event) {
		dispatchersMap[id] = append(dispatchersMap[id], e)
	}

	for _, event := range blockEvents.Events {
		for _, subscription := range subscriptions {
			if wh.filter.MatchEvent(subscription, event) {
				mapEventToDispatcher(subscription.DispatcherID, event)
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

func (wh *commonHub) handleRevertBroadcast(revertBlock data.RevertBlock) {
}

func (wh *commonHub) registerDispatcher(d dispatcher.EventDispatcher) {
	wh.rwMut.Lock()
	defer wh.rwMut.Unlock()

	if _, ok := wh.dispatchers[d.GetID()]; ok {
		return
	}

	wh.dispatchers[d.GetID()] = d

	log.Info("registered new dispatcher", "dispatcherID", d.GetID())
}

func (wh *commonHub) unregisterDispatcher(d dispatcher.EventDispatcher) {
	wh.rwMut.Lock()
	defer wh.rwMut.Unlock()

	if _, ok := wh.dispatchers[d.GetID()]; ok {
		delete(wh.dispatchers, d.GetID())
	}

	log.Info("unregistered dispatcher", "dispatcherID", d.GetID(), "unsubscribing", true)

	wh.subscriptionMapper.RemoveSubscriptions(d.GetID())
}
