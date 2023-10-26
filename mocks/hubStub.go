package mocks

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// HubStub implements Hub interface
type HubStub struct {
	RegisterListenerCalled            func() error
	PublishCalled                     func(events data.BlockEvents)
	PublishRevertCalled               func(revertBlock data.RevertBlock)
	PublishFinalizedCalled            func(finalizedBlock data.FinalizedBlock)
	PublishTxsCalled                  func(blockTxs data.BlockTxs)
	PublishScrsCalled                 func(blockScrs data.BlockScrs)
	PublishBlockEventsWithOrderCalled func(blockTxs data.BlockEventsWithOrder)
	RegisterEventCalled               func(event dispatcher.EventDispatcher)
	UnregisterEventCalled             func(event dispatcher.EventDispatcher)
	SubscribeCalled                   func(event data.SubscribeEvent)
	CloseCalled                       func() error
}

// RegisterListener -
func (h *HubStub) RegisterListener() error {
	if h.RegisterListenerCalled != nil {
		return h.RegisterListenerCalled()
	}

	return nil
}

// Publish -
func (h *HubStub) Publish(events data.BlockEvents) {
	if h.PublishCalled != nil {
		h.PublishCalled(events)
	}
}

// PublishRevert -
func (h *HubStub) PublishRevert(revertBlock data.RevertBlock) {
	if h.PublishRevertCalled != nil {
		h.PublishRevertCalled(revertBlock)
	}
}

// PublishFinalized -
func (h *HubStub) PublishFinalized(finalizedBlock data.FinalizedBlock) {
	if h.PublishFinalizedCalled != nil {
		h.PublishFinalizedCalled(finalizedBlock)
	}
}

// PublishTxs -
func (h *HubStub) PublishTxs(blockTxs data.BlockTxs) {
	if h.PublishTxsCalled != nil {
		h.PublishTxsCalled(blockTxs)
	}
}

// PublishScrs -
func (h *HubStub) PublishScrs(blockScrs data.BlockScrs) {
	if h.PublishScrsCalled != nil {
		h.PublishScrsCalled(blockScrs)
	}
}

// PublishBlockEventsWithOrder -
func (h *HubStub) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	if h.PublishBlockEventsWithOrderCalled != nil {
		h.PublishBlockEventsWithOrderCalled(blockTxs)
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
func (h *HubStub) Subscribe(event data.SubscribeEvent) {
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
