package disabled

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// Hub defines a disabled hub component
type Hub struct {
}

// RegisterListener does nothing
func (h *Hub) RegisterListener() error {
	return nil
}

// Publish -
func (h *Hub) Publish(events data.BlockEvents) {
}

// PublishRevert -
func (h *Hub) PublishRevert(revertBlock data.RevertBlock) {
}

// PublishFinalized -
func (h *Hub) PublishFinalized(finalizedBlock data.FinalizedBlock) {
}

// PublishTxs -
func (h *Hub) PublishTxs(blockTxs data.BlockTxs) {
}

// PublishScrs -
func (h *Hub) PublishScrs(blockScrs data.BlockScrs) {
}

// PublishBlockEventsWithOrder -
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
