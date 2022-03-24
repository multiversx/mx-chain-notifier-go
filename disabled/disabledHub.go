package disabled

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// Hub defines a disabled hub component
type Hub struct {
	RunCalled                func()
	BroadcastCalled          func(events data.BlockEvents)
	BroadcastRevertCalled    func(event data.RevertBlock)
	BroadcastFinalizedCalled func(event data.FinalizedBlock)
	RegisterEventCalled      func(event dispatcher.EventDispatcher)
	UnregisterEventCalled    func(event dispatcher.EventDispatcher)
	SubscribeCalled          func(event dispatcher.SubscribeEvent)
	CloseCalled              func() error
}

// Run -
func (h *Hub) Run() {
	if h.RunCalled != nil {
		h.RunCalled()
	}
}

// Broadcast -
func (h *Hub) Broadcast(events data.BlockEvents) {
	if h.BroadcastCalled != nil {
		h.BroadcastCalled(events)
	}
}

// BroadcastRevert -
func (h *Hub) BroadcastRevert(event data.RevertBlock) {
	if h.BroadcastRevertCalled != nil {
		h.BroadcastRevertCalled(event)
	}
}

// BroadcastFinalized -
func (h *Hub) BroadcastFinalized(event data.FinalizedBlock) {
	if h.BroadcastFinalizedCalled != nil {
		h.BroadcastFinalizedCalled(event)
	}
}

// RegisterEvent -
func (h *Hub) RegisterEvent(event dispatcher.EventDispatcher) {
	if h.RegisterEventCalled != nil {
		h.RegisterEventCalled(event)
	}
}

// UnregisterEvent -
func (h *Hub) UnregisterEvent(event dispatcher.EventDispatcher) {
	if h.UnregisterEventCalled != nil {
		h.UnregisterEventCalled(event)
	}
}

// Subscribe -
func (h *Hub) Subscribe(event dispatcher.SubscribeEvent) {
	if h.SubscribeCalled != nil {
		h.SubscribeCalled(event)
	}
}

// Close -
func (h *Hub) Close() error {
	return nil
}

// IsInterfaceNil -
func (h *Hub) IsInterfaceNil() bool {
	return h == nil
}
