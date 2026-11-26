package process

import (
	"context"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

type publisher struct {
	handler PublisherHandler

	// requests is served by a single goroutine, so that the handler publishes one item at a time
	requests chan publishRequest

	cancelFunc func()
	closeChan  chan struct{}
	mutState   sync.RWMutex
}

// publishRequest holds a publish operation and the channel on which its result is sent back
type publishRequest struct {
	publish func() error
	result  chan error
}

// NewPublisher will create a new publisher component
func NewPublisher(handler PublisherHandler) (*publisher, error) {
	if check.IfNil(handler) {
		return nil, ErrNilPublisherHandler
	}

	p := &publisher{
		handler:   handler,
		requests:  make(chan publishRequest),
		closeChan: make(chan struct{}),
	}

	return p, nil
}

// Run creates a goroutine and listens for events on the exposed channels
func (p *publisher) Run() error {
	p.mutState.Lock()
	defer p.mutState.Unlock()

	if p.cancelFunc != nil {
		return common.ErrLoopAlreadyStarted
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
		case req := <-p.requests:
			req.result <- req.publish()
		}
	}
}

// broadcast hands the publish operation to the publishing goroutine and waits for its result,
// so that the caller knows whether the item was published or not
func (p *publisher) broadcast(publish func() error) error {
	req := publishRequest{
		publish: publish,
		result:  make(chan error, 1),
	}

	select {
	case <-p.closeChan:
		return ErrPublisherClosed
	default:
	}

	select {
	case p.requests <- req:
	case <-p.closeChan:
		return ErrPublisherClosed
	}

	select {
	case err := <-req.result:
		return err
	case <-p.closeChan:
		return ErrPublisherClosed
	}
}

// Broadcast will handle the block events pushed by producers
func (p *publisher) Broadcast(events data.BlockEvents) error {
	return p.broadcast(func() error {
		return p.handler.Publish(events)
	})
}

// BroadcastRevert will handle the revert event pushed by producers
func (p *publisher) BroadcastRevert(events data.RevertBlock) error {
	return p.broadcast(func() error {
		return p.handler.PublishRevert(events)
	})
}

// BroadcastFinalized will handle the finalized event pushed by producers
func (p *publisher) BroadcastFinalized(events data.FinalizedBlock) error {
	return p.broadcast(func() error {
		return p.handler.PublishFinalized(events)
	})
}

// BroadcastTxs will handle the txs event pushed by producers
func (p *publisher) BroadcastTxs(events data.BlockTxs) error {
	return p.broadcast(func() error {
		return p.handler.PublishTxs(events)
	})
}

// BroadcastScrs will handle the scrs event pushed by producers
func (p *publisher) BroadcastScrs(events data.BlockScrs) error {
	return p.broadcast(func() error {
		return p.handler.PublishScrs(events)
	})
}

// BroadcastBlockEventsWithOrder will handle the full block events pushed by producers
func (p *publisher) BroadcastBlockEventsWithOrder(events data.BlockEventsWithOrder) error {
	return p.broadcast(func() error {
		return p.handler.PublishBlockEventsWithOrder(events)
	})
}

// BroadcastStateAccesses will handle state accesses pushed by producers
func (p *publisher) BroadcastStateAccesses(events data.BlockStateAccesses) error {
	return p.broadcast(func() error {
		return p.handler.PublishStateAccesses(events)
	})
}

// CheckSaveBlockMessages checks that all the messages of a saved block can be published, so that
// a block is not published partially. It does not go through the publishing goroutine, since it
// does not publish anything.
func (p *publisher) CheckSaveBlockMessages(messages data.SaveBlockMessages) error {
	return p.handler.CheckSaveBlockMessages(messages)
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
