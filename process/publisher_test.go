package process_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/stretchr/testify/require"
)

func TestNewPublisher(t *testing.T) {
	t.Parallel()

	t.Run("nil handler", func(t *testing.T) {
		t.Parallel()

		p, err := process.NewPublisher(nil)
		require.Nil(t, p)
		require.Equal(t, process.ErrNilPublisherHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		p, err := process.NewPublisher(&mocks.PublisherHandlerStub{})
		require.Nil(t, err)
		require.False(t, p.IsInterfaceNil())
	})
}

func TestBroadcast(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	ph := &mocks.PublisherHandlerStub{
		PublishCalled: func(events data.BlockEvents) {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
		},
	}

	p, err := process.NewPublisher(ph)
	require.Nil(t, err)

	p.Run()
	defer p.Close()
	wg.Add(1)

	p.Broadcast(data.BlockEvents{})

	wg.Wait()

	require.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestBroadcastRevert(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	ph := &mocks.PublisherHandlerStub{
		PublishRevertCalled: func(revertBlock data.RevertBlock) {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
		},
	}

	p, err := process.NewPublisher(ph)
	require.Nil(t, err)

	p.Run()
	defer p.Close()
	wg.Add(1)

	p.BroadcastRevert(data.RevertBlock{})

	wg.Wait()

	require.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestBroadcastFinalized(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	ph := &mocks.PublisherHandlerStub{
		PublishFinalizedCalled: func(finalizedBlock data.FinalizedBlock) {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
		},
	}

	p, err := process.NewPublisher(ph)
	require.Nil(t, err)

	p.Run()
	defer p.Close()
	wg.Add(1)

	p.BroadcastFinalized(data.FinalizedBlock{})

	wg.Wait()

	require.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestBroadcastTxs(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	ph := &mocks.PublisherHandlerStub{
		PublishTxsCalled: func(blockTxs data.BlockTxs) {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
		},
	}

	p, err := process.NewPublisher(ph)
	require.Nil(t, err)

	p.Run()
	defer p.Close()
	wg.Add(1)

	p.BroadcastTxs(data.BlockTxs{})

	wg.Wait()

	require.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestBroadcastScrs(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	ph := &mocks.PublisherHandlerStub{
		PublishScrsCalled: func(blockScrs data.BlockScrs) {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
		},
	}

	p, err := process.NewPublisher(ph)
	require.Nil(t, err)

	p.Run()
	defer p.Close()
	wg.Add(1)

	p.BroadcastScrs(data.BlockScrs{})

	wg.Wait()

	require.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestBroadcastBlockEventsWithOrder(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	numCalls := uint32(0)

	ph := &mocks.PublisherHandlerStub{
		PublishBlockEventsWithOrderCalled: func(blockTxs data.BlockEventsWithOrder) {
			atomic.AddUint32(&numCalls, 1)
			wg.Done()
		},
	}

	p, err := process.NewPublisher(ph)
	require.Nil(t, err)

	p.Run()
	defer p.Close()
	wg.Add(1)

	p.BroadcastBlockEventsWithOrder(data.BlockEventsWithOrder{})

	wg.Wait()

	require.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestClose(t *testing.T) {
	t.Parallel()

	p, err := process.NewPublisher(&mocks.PublisherHandlerStub{})
	require.Nil(t, err)

	p.Run()

	err = p.Close()
	require.Nil(t, err)
}
