package hub

import (
	"github.com/google/uuid"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

func (ch *commonHub) CheckDispatcherByID(uuid uuid.UUID, dispatcher dispatcher.EventDispatcher) bool {
	ch.mutDispatchers.Lock()
	defer ch.mutDispatchers.Unlock()

	return ch.dispatchers[uuid] == dispatcher
}
