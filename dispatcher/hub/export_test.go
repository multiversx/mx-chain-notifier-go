package hub

import (
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/google/uuid"
)

func (ch *commonHub) CheckDispatcherByID(uuid uuid.UUID, dispatcher dispatcher.EventDispatcher) bool {
	ch.mutDispatchers.Lock()
	defer ch.mutDispatchers.Unlock()

	return ch.dispatchers[uuid] == dispatcher
}
