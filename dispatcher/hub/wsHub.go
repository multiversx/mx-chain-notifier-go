package hub

import (
	"sync"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
)

type Subscription struct {
	Address     string
	Identifier  string
	Subscribers []dispatcher.EventDispatcher
}

type wsHub struct {
	rwMut       sync.RWMutex
	dispatchers map[uuid.UUID]dispatcher.EventDispatcher
	register    chan dispatcher.EventDispatcher
	unregister  chan dispatcher.EventDispatcher
	broadcast   chan []data.Event
}

func NewWsHub() *wsHub {
	return &wsHub{
		rwMut:       sync.RWMutex{},
		dispatchers: make(map[uuid.UUID]dispatcher.EventDispatcher),
		register:    make(chan dispatcher.EventDispatcher),
		unregister:  make(chan dispatcher.EventDispatcher),
		broadcast:   make(chan []data.Event),
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

func (wh *wsHub) BroadcastChan() chan []data.Event {
	return wh.broadcast
}

func (wh *wsHub) RegisterChan() chan dispatcher.EventDispatcher {
	return wh.register
}

func (wh *wsHub) UnregisterChan() chan dispatcher.EventDispatcher {
	return wh.unregister
}

// TODO: events will be filtered before broadcasting
func (wh *wsHub) handleBroadcast(events []data.Event) {
	wh.rwMut.RLock()
	defer wh.rwMut.RUnlock()

	for _, d := range wh.dispatchers {
		d.PushEvents(events)
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
	if _, ok := wh.dispatchers[d.GetID()]; ok {
		delete(wh.dispatchers, d.GetID())
	}
}
