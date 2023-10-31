package disabled

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// Hub defines a disabled hub component
type Hub struct {
}

// Publish does nothing
func (h *Hub) Publish(events data.BlockEvents) {
}

// PublishRevert does nothing
func (h *Hub) PublishRevert(revertBlock data.RevertBlock) {
}

// PublishFinalized does nothing
func (h *Hub) PublishFinalized(finalizedBlock data.FinalizedBlock) {
}

// PublishTxs does nothing
func (h *Hub) PublishTxs(blockTxs data.BlockTxs) {
}

// PublishScrs does nothing
func (h *Hub) PublishScrs(blockScrs data.BlockScrs) {
}

// PublishBlockEventsWithOrder does nothing
func (h *Hub) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
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
