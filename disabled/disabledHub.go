package disabled

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// Hub defines a disabled hub component
type Hub struct{}

// Run does nothing
func (dh *Hub) Run() {
}

// BroadcastChan returns a nil channel
func (dh *Hub) BroadcastChan() chan<- data.BlockEvents {
	return nil
}

// BroadcastRevertChan returns a nil channel
func (dh *Hub) BroadcastRevertChan() chan<- data.RevertBlock {
	return nil
}

// BroadcastFinalizedChan returns a nil channel
func (dh *Hub) BroadcastFinalizedChan() chan<- data.FinalizedBlock {
	return nil
}

// RegisterChan returns a nil channel
func (dh *Hub) RegisterChan() chan<- dispatcher.EventDispatcher {
	return nil
}

// UnregisterChan returns a nil channel
func (dh *Hub) UnregisterChan() chan<- dispatcher.EventDispatcher {
	return nil
}

// Subscribe returns nothing
func (dh *Hub) Subscribe(event dispatcher.SubscribeEvent) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (dh *Hub) IsInterfaceNil() bool {
	return false
}
