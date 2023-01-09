package hub

import (
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/google/uuid"
)

func (ch *commonHub) CheckDispatcherByID(uuid uuid.UUID, dispatcher dispatcher.EventDispatcher) bool {
	ch.rwMut.Lock()
	defer ch.rwMut.Unlock()

	return ch.dispatchers[uuid] == dispatcher
}
