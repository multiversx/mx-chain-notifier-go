package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// TODO: change this after dispatcher refactoring

// HubStub implements Hub interface
type HubStub struct {
}

// Run -
func (h *HubStub) Run() {
	panic("not implemented") // TODO: Implement
}

// BroadcastChan -
func (h *HubStub) BroadcastChan() chan<- data.BlockEvents {
	panic("not implemented") // TODO: Implement
}

// BroadcastRevertChan -
func (h *HubStub) BroadcastRevertChan() chan<- data.RevertBlock {
	panic("not implemented") // TODO: Implement
}

// BroadcastFinalizedChan -
func (h *HubStub) BroadcastFinalizedChan() chan<- data.FinalizedBlock {
	panic("not implemented") // TODO: Implement
}

// RegisterChan -
func (h *HubStub) RegisterChan() chan<- dispatcher.EventDispatcher {
	panic("not implemented") // TODO: Implement
}

// UnregisterChan -
func (h *HubStub) UnregisterChan() chan<- dispatcher.EventDispatcher {
	panic("not implemented") // TODO: Implement
}

// Subscribe -
func (h *HubStub) Subscribe(event dispatcher.SubscribeEvent) {
	panic("not implemented") // TODO: Implement
}

// IsInterfaceNil -
func (h *HubStub) IsInterfaceNil() bool {
	return h == nil
}
