package mocks

import "github.com/multiversx/mx-chain-notifier-go/data"

// PublisherHandlerStub -
type PublisherHandlerStub struct {
	PublishCalled                     func(events data.BlockEvents)
	PublishRevertCalled               func(revertBlock data.RevertBlock)
	PublishFinalizedCalled            func(finalizedBlock data.FinalizedBlock)
	PublishTxsCalled                  func(blockTxs data.BlockTxs)
	PublishScrsCalled                 func(blockScrs data.BlockScrs)
	PublishBlockEventsWithOrderCalled func(blockTxs data.BlockEventsWithOrder)
	PublishBlockStateAccessesCalled   func(stateAccesses data.BlockStateAccesses)
	CheckSaveBlockMessagesCalled      func(messages data.SaveBlockMessages) error
	CloseCalled                       func() error
}

// Publish -
func (p *PublisherHandlerStub) Publish(events data.BlockEvents) error {
	if p.PublishCalled != nil {
		p.PublishCalled(events)
	}

	return nil
}

// PublishRevert -
func (p *PublisherHandlerStub) PublishRevert(revertBlock data.RevertBlock) error {
	if p.PublishRevertCalled != nil {
		p.PublishRevertCalled(revertBlock)
	}

	return nil
}

// PublishFinalized -
func (p *PublisherHandlerStub) PublishFinalized(finalizedBlock data.FinalizedBlock) error {
	if p.PublishFinalizedCalled != nil {
		p.PublishFinalizedCalled(finalizedBlock)
	}

	return nil
}

// PublishTxs -
func (p *PublisherHandlerStub) PublishTxs(blockTxs data.BlockTxs) error {
	if p.PublishTxsCalled != nil {
		p.PublishTxsCalled(blockTxs)
	}

	return nil
}

// PublishScrs -
func (p *PublisherHandlerStub) PublishScrs(blockScrs data.BlockScrs) error {
	if p.PublishScrsCalled != nil {
		p.PublishScrsCalled(blockScrs)
	}

	return nil
}

// PublishBlockEventsWithOrder -
func (p *PublisherHandlerStub) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) error {
	if p.PublishBlockEventsWithOrderCalled != nil {
		p.PublishBlockEventsWithOrderCalled(blockTxs)
	}

	return nil
}

// PublishStateAccesses -
func (p *PublisherHandlerStub) PublishStateAccesses(stateAccesses data.BlockStateAccesses) error {
	if p.PublishBlockStateAccessesCalled != nil {
		p.PublishBlockStateAccessesCalled(stateAccesses)
	}

	return nil
}

// CheckSaveBlockMessages -
func (p *PublisherHandlerStub) CheckSaveBlockMessages(messages data.SaveBlockMessages) error {
	if p.CheckSaveBlockMessagesCalled != nil {
		return p.CheckSaveBlockMessagesCalled(messages)
	}

	return nil
}

// Close -
func (p *PublisherHandlerStub) Close() error {
	if p.CloseCalled != nil {
		return p.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (p *PublisherHandlerStub) IsInterfaceNil() bool {
	return p == nil
}
