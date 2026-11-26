package mocks

import "github.com/multiversx/mx-chain-notifier-go/data"

// PublisherStub implements PublisherService interface
type PublisherStub struct {
	RunCalled                           func() error
	BroadcastCalled                     func(events data.BlockEvents)
	BroadcastRevertCalled               func(event data.RevertBlock)
	BroadcastFinalizedCalled            func(event data.FinalizedBlock)
	BroadcastTxsCalled                  func(event data.BlockTxs)
	BroadcastScrsCalled                 func(event data.BlockScrs)
	BroadcastBlockEventsWithOrderCalled func(event data.BlockEventsWithOrder)
	BroadcastStateAccessesCalled        func(event data.BlockStateAccesses)
	CheckSaveBlockMessagesCalled        func(messages data.SaveBlockMessages) error
	CloseCalled                         func() error
}

// Run -
func (ps *PublisherStub) Run() error {
	if ps.RunCalled != nil {
		return ps.RunCalled()
	}

	return nil
}

// Broadcast -
func (ps *PublisherStub) Broadcast(events data.BlockEvents) error {
	if ps.BroadcastCalled != nil {
		ps.BroadcastCalled(events)
	}

	return nil
}

// BroadcastRevert -
func (ps *PublisherStub) BroadcastRevert(event data.RevertBlock) error {
	if ps.BroadcastRevertCalled != nil {
		ps.BroadcastRevertCalled(event)
	}

	return nil
}

// BroadcastFinalized -
func (ps *PublisherStub) BroadcastFinalized(event data.FinalizedBlock) error {
	if ps.BroadcastFinalizedCalled != nil {
		ps.BroadcastFinalizedCalled(event)
	}

	return nil
}

// BroadcastTxs -
func (ps *PublisherStub) BroadcastTxs(event data.BlockTxs) error {
	if ps.BroadcastTxsCalled != nil {
		ps.BroadcastTxsCalled(event)
	}

	return nil
}

// BroadcastScrs -
func (ps *PublisherStub) BroadcastScrs(event data.BlockScrs) error {
	if ps.BroadcastScrsCalled != nil {
		ps.BroadcastScrsCalled(event)
	}

	return nil
}

// BroadcastBlockEventsWithOrder -
func (ps *PublisherStub) BroadcastBlockEventsWithOrder(event data.BlockEventsWithOrder) error {
	if ps.BroadcastBlockEventsWithOrderCalled != nil {
		ps.BroadcastBlockEventsWithOrderCalled(event)
	}

	return nil
}

// BroadcastStateAccesses -
func (ps *PublisherStub) BroadcastStateAccesses(event data.BlockStateAccesses) error {
	if ps.BroadcastStateAccessesCalled != nil {
		ps.BroadcastStateAccessesCalled(event)
	}

	return nil
}

// CheckSaveBlockMessages -
func (ps *PublisherStub) CheckSaveBlockMessages(messages data.SaveBlockMessages) error {
	if ps.CheckSaveBlockMessagesCalled != nil {
		return ps.CheckSaveBlockMessagesCalled(messages)
	}

	return nil
}

// Close -
func (ps *PublisherStub) Close() error {
	if ps.CloseCalled != nil {
		return ps.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (ps *PublisherStub) IsInterfaceNil() bool {
	return ps == nil
}
