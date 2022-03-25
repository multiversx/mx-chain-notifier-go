package disabled

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// Hub defines a disabled hub component
type Hub struct {
}

// Run does nothing
func (h *Hub) Run() {
}

// Broadcast does nothing
func (h *Hub) Broadcast(events data.BlockEvents) {
}

// BroadcastRevert does nothing
func (h *Hub) BroadcastRevert(event data.RevertBlock) {
}

// BroadcastFinalized does nothing
func (h *Hub) BroadcastFinalized(event data.FinalizedBlock) {
}

// RegisterEvent does nothing
func (h *Hub) RegisterEvent(event dispatcher.EventDispatcher) {
}

// UnregisterEvent does nothing
func (h *Hub) UnregisterEvent(event dispatcher.EventDispatcher) {
}

// Subscribe does nothing
func (h *Hub) Subscribe(event dispatcher.SubscribeEvent) {
}

// Close returns nil
func (h *Hub) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *Hub) IsInterfaceNil() bool {
	return h == nil
}
