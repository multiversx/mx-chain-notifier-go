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
	CloseCalled                       func() error
}

// Publish -
func (p *PublisherHandlerStub) Publish(events data.BlockEvents) {
	if p.PublishCalled != nil {
		p.PublishCalled(events)
	}
}

// PublishRevert -
func (p *PublisherHandlerStub) PublishRevert(revertBlock data.RevertBlock) {
	if p.PublishRevertCalled != nil {
		p.PublishRevertCalled(revertBlock)
	}
}

// PublishFinalized -
func (p *PublisherHandlerStub) PublishFinalized(finalizedBlock data.FinalizedBlock) {
	if p.PublishFinalizedCalled != nil {
		p.PublishFinalizedCalled(finalizedBlock)
	}
}

// PublishTxs -
func (p *PublisherHandlerStub) PublishTxs(blockTxs data.BlockTxs) {
	if p.PublishTxsCalled != nil {
		p.PublishTxsCalled(blockTxs)
	}
}

// PublishScrs -
func (p *PublisherHandlerStub) PublishScrs(blockScrs data.BlockScrs) {
	if p.PublishScrsCalled != nil {
		p.PublishScrsCalled(blockScrs)
	}
}

// PublishBlockEventsWithOrder -
func (p *PublisherHandlerStub) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	if p.PublishBlockEventsWithOrderCalled != nil {
		p.PublishBlockEventsWithOrderCalled(blockTxs)
	}
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
