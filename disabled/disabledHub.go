package disabled

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// Hub defines a disabled hub component
type Hub struct {
}

// Run does nothing
func (h *Hub) Run() {
}

// Broadcast does nothing
func (h *Hub) Broadcast(_ data.BlockEvents) {
}

// BroadcastRevert does nothing
func (h *Hub) BroadcastRevert(_ data.RevertBlock) {
}

// BroadcastFinalized does nothing
func (h *Hub) BroadcastFinalized(_ data.FinalizedBlock) {
}

// BroadcastTxs does nothing
func (h *Hub) BroadcastTxs(_ data.BlockTxs) {
}

// BroadcastScrs does nothing
func (h *Hub) BroadcastScrs(_ data.BlockScrs) {
}

// RegisterEvent does nothing
func (h *Hub) RegisterEvent(_ dispatcher.EventDispatcher) {
}

// UnregisterEvent does nothing
func (h *Hub) UnregisterEvent(_ dispatcher.EventDispatcher) {
}

// Subscribe does nothing
func (h *Hub) Subscribe(_ data.SubscribeEvent) {
}

// Close returns nil
func (h *Hub) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *Hub) IsInterfaceNil() bool {
	return h == nil
}
