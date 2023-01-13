package mocks

import "github.com/ElrondNetwork/notifier-go/data"

// PublisherStub implements PublisherService interface
type PublisherStub struct {
	RunCalled                   func()
	BroadcastCalled             func(events data.BlockEvents)
	BroadcastRevertCalled       func(event data.RevertBlock)
	BroadcastFinalizedCalled    func(event data.FinalizedBlock)
	BroadcastTxsCalled          func(event data.BlockTxs)
	BroadcastScrsCalled         func(event data.BlockScrs)
	BroadcastTxsWithOrderCalled func(event data.BlockTxsWithOrder)
}

// Run -
func (ps *PublisherStub) Run() {
	if ps.RunCalled != nil {
		ps.RunCalled()
	}
}

// Broadcast -
func (ps *PublisherStub) Broadcast(events data.BlockEvents) {
	if ps.BroadcastCalled != nil {
		ps.BroadcastCalled(events)
	}
}

// BroadcastRevert -
func (ps *PublisherStub) BroadcastRevert(event data.RevertBlock) {
	if ps.BroadcastRevertCalled != nil {
		ps.BroadcastRevertCalled(event)
	}
}

// BroadcastFinalized -
func (ps *PublisherStub) BroadcastFinalized(event data.FinalizedBlock) {
	if ps.BroadcastFinalizedCalled != nil {
		ps.BroadcastFinalizedCalled(event)
	}
}

// BroadcastTxs -
func (ps *PublisherStub) BroadcastTxs(event data.BlockTxs) {
	if ps.BroadcastTxsCalled != nil {
		ps.BroadcastTxsCalled(event)
	}
}

// BroadcastScrs -
func (ps *PublisherStub) BroadcastScrs(event data.BlockScrs) {
	if ps.BroadcastScrsCalled != nil {
		ps.BroadcastScrsCalled(event)
	}
}

// BroadcastTxsWithOrder -
func (ps *PublisherStub) BroadcastTxsWithOrder(event data.BlockTxsWithOrder) {
	if ps.BroadcastTxsWithOrderCalled != nil {
		ps.BroadcastTxsWithOrderCalled(event)
	}
}

// IsInterfaceNil -
func (ps *PublisherStub) IsInterfaceNil() bool {
	return ps == nil
}
