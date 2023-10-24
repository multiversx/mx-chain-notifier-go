package process

import (
	"context"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

type publisher struct {
	handler PublisherHandler

	broadcast                     chan data.BlockEvents
	broadcastRevert               chan data.RevertBlock
	broadcastFinalized            chan data.FinalizedBlock
	broadcastTxs                  chan data.BlockTxs
	broadcastBlockEventsWithOrder chan data.BlockEventsWithOrder
	broadcastScrs                 chan data.BlockScrs

	cancelFunc func()
	closeChan  chan struct{}
	mutState   sync.RWMutex
}

// NewPublisher will create a new publisher component
func NewPublisher(handler PublisherHandler) (*publisher, error) {
	if check.IfNil(handler) {
		return nil, ErrNilPublisherHandler
	}

	p := &publisher{
		handler:                       handler,
		broadcast:                     make(chan data.BlockEvents),
		broadcastRevert:               make(chan data.RevertBlock),
		broadcastFinalized:            make(chan data.FinalizedBlock),
		broadcastTxs:                  make(chan data.BlockTxs),
		broadcastScrs:                 make(chan data.BlockScrs),
		broadcastBlockEventsWithOrder: make(chan data.BlockEventsWithOrder),
		closeChan:                     make(chan struct{}),
	}

	return p, nil
}

// Run creates a goroutine and listens for events on the exposed channels
func (p *publisher) Run() error {
	p.mutState.Lock()
	defer p.mutState.Unlock()

	if p.cancelFunc != nil {
		return ErrLoopAlreadyStarted
	}

	var ctx context.Context
	ctx, p.cancelFunc = context.WithCancel(context.Background())

	go p.run(ctx)

	return nil
}

func (p *publisher) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.handler.Close()
			return
		case events := <-p.broadcast:
			p.handler.Publish(events)
		case revertBlock := <-p.broadcastRevert:
			p.handler.PublishRevert(revertBlock)
		case finalizedBlock := <-p.broadcastFinalized:
			p.handler.PublishFinalized(finalizedBlock)
		case blockTxs := <-p.broadcastTxs:
			p.handler.PublishTxs(blockTxs)
		case blockScrs := <-p.broadcastScrs:
			p.handler.PublishScrs(blockScrs)
		case blockEvents := <-p.broadcastBlockEventsWithOrder:
			p.handler.PublishBlockEventsWithOrder(blockEvents)
		}
	}
}

// Broadcast will handle the block events pushed by producers
func (p *publisher) Broadcast(events data.BlockEvents) {
	select {
	case p.broadcast <- events:
	case <-p.closeChan:
	}
}

// BroadcastRevert will handle the revert event pushed by producers
func (p *publisher) BroadcastRevert(events data.RevertBlock) {
	select {
	case p.broadcastRevert <- events:
	case <-p.closeChan:
	}
}

// BroadcastFinalized will handle the finalized event pushed by producers
func (p *publisher) BroadcastFinalized(events data.FinalizedBlock) {
	select {
	case p.broadcastFinalized <- events:
	case <-p.closeChan:
	}
}

// BroadcastTxs will handle the txs event pushed by producers
func (p *publisher) BroadcastTxs(events data.BlockTxs) {
	select {
	case p.broadcastTxs <- events:
	case <-p.closeChan:
	}
}

// BroadcastScrs will handle the scrs event pushed by producers
func (p *publisher) BroadcastScrs(events data.BlockScrs) {
	select {
	case p.broadcastScrs <- events:
	case <-p.closeChan:
	}
}

// BroadcastBlockEventsWithOrder will handle the full block events pushed by producers
func (p *publisher) BroadcastBlockEventsWithOrder(events data.BlockEventsWithOrder) {
	select {
	case p.broadcastBlockEventsWithOrder <- events:
	case <-p.closeChan:
	}
}

// Close will close the channels
func (p *publisher) Close() error {
	p.mutState.RLock()
	defer p.mutState.RUnlock()

	if p.cancelFunc != nil {
		p.cancelFunc()
	}

	close(p.closeChan)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (p *publisher) IsInterfaceNil() bool {
	return p == nil
}
