package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// HubStub implements Hub interface
type HubStub struct {
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
func (h *HubStub) Run() {
	if h.RunCalled != nil {
		h.RunCalled()
	}
}

// Broadcast -
func (h *HubStub) Broadcast(events data.BlockEvents) {
	if h.BroadcastCalled != nil {
		h.BroadcastCalled(events)
	}
}

// BroadcastRevert -
func (h *HubStub) BroadcastRevert(event data.RevertBlock) {
	if h.BroadcastRevertCalled != nil {
		h.BroadcastRevertCalled(event)
	}
}

// BroadcastFinalized -
func (h *HubStub) BroadcastFinalized(event data.FinalizedBlock) {
	if h.BroadcastFinalizedCalled != nil {
		h.BroadcastFinalizedCalled(event)
	}
}

// RegisterEvent -
func (h *HubStub) RegisterEvent(event dispatcher.EventDispatcher) {
	if h.RegisterEventCalled != nil {
		h.RegisterEventCalled(event)
	}
}

// UnregisterEvent -
func (h *HubStub) UnregisterEvent(event dispatcher.EventDispatcher) {
	if h.UnregisterEventCalled != nil {
		h.UnregisterEventCalled(event)
	}
}

// Subscribe -
func (h *HubStub) Subscribe(event dispatcher.SubscribeEvent) {
	if h.SubscribeCalled != nil {
		h.SubscribeCalled(event)
	}
}

// Close -
func (h *HubStub) Close() error {
	return nil
}

// IsInterfaceNil -
func (h *HubStub) IsInterfaceNil() bool {
	return h == nil
}
