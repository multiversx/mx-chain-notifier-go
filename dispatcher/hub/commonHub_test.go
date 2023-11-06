package hub

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/filters"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockCommonHubArgs() ArgsCommonHub {
	return ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
}

func TestNewCommonHub(t *testing.T) {
	t.Parallel()

	t.Run("nil event filter", func(t *testing.T) {
		t.Parallel()

		args := createMockCommonHubArgs()
		args.Filter = nil

		hub, err := NewCommonHub(args)
		require.Nil(t, hub)
		assert.Equal(t, ErrNilEventFilter, err)
	})

	t.Run("nil subscription mapper", func(t *testing.T) {
		t.Parallel()

		args := createMockCommonHubArgs()
		args.SubscriptionMapper = nil

		hub, err := NewCommonHub(args)
		require.Nil(t, hub)
		assert.Equal(t, ErrNilSubscriptionMapper, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockCommonHubArgs()
		hub, err := NewCommonHub(args)
		require.Nil(t, err)
		require.False(t, hub.IsInterfaceNil())
	})
}

func TestCommonHub_RegisterDispatcher(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	dispatcher1 := mocks.NewDispatcherMock(nil, hub)
	dispatcher2 := mocks.NewDispatcherMock(nil, hub)

	hub.RegisterEvent(dispatcher1)
	hub.RegisterEvent(dispatcher2)

	time.Sleep(time.Millisecond * 100)

	require.True(t, hub.CheckDispatcherByID(dispatcher1.GetID(), dispatcher1))
	require.True(t, hub.CheckDispatcherByID(dispatcher2.GetID(), dispatcher2))
}

func TestCommonHub_UnregisterDispatcher(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	dispatcher1 := mocks.NewDispatcherMock(nil, hub)
	dispatcher2 := mocks.NewDispatcherMock(nil, hub)

	hub.RegisterEvent(dispatcher1)
	hub.RegisterEvent(dispatcher2)

	hub.UnregisterEvent(dispatcher1)

	time.Sleep(time.Millisecond * 100)

	require.True(t, hub.CheckDispatcherByID(dispatcher1.GetID(), nil))
	require.True(t, hub.CheckDispatcherByID(dispatcher2.GetID(), dispatcher2))
}

func TestCommonHub_HandleBroadcastDispatcherReceivesEvents(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	consumer := mocks.NewConsumerMock()
	dispatcher1 := mocks.NewDispatcherMock(consumer, hub)

	hub.RegisterEvent(dispatcher1)
	hub.Subscribe(data.SubscribeEvent{
		DispatcherID:        dispatcher1.GetID(),
		SubscriptionEntries: []data.SubscriptionEntry{},
	})

	blockEvents := getEvents()

	hub.Publish(blockEvents)

	time.Sleep(time.Millisecond * 100)

	require.True(t, consumer.HasEvents(blockEvents.Events))
	require.True(t, len(consumer.CollectedEvents()) == len(blockEvents.Events))
}

func TestCommonHub_HandleBroadcastMultipleDispatchers(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	consumer1 := mocks.NewConsumerMock()
	dispatcher1 := mocks.NewDispatcherMock(consumer1, hub)
	consumer2 := mocks.NewConsumerMock()
	dispatcher2 := mocks.NewDispatcherMock(consumer2, hub)
	consumer3 := mocks.NewConsumerMock()
	dispatcher3 := mocks.NewDispatcherMock(consumer3, hub)

	hub.RegisterEvent(dispatcher1)
	hub.RegisterEvent(dispatcher2)
	hub.RegisterEvent(dispatcher3)

	time.Sleep(time.Millisecond * 100)

	hub.Subscribe(data.SubscribeEvent{
		DispatcherID: dispatcher1.GetID(),
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				Address:    "erd1",
				Identifier: "swap",
			},
		},
	})
	hub.Subscribe(data.SubscribeEvent{
		DispatcherID: dispatcher2.GetID(),
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				Address:    "erd2",
				Identifier: "lock",
			},
		},
	})
	hub.Subscribe(data.SubscribeEvent{
		DispatcherID: dispatcher3.GetID(),
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				Address:    "erd3",
				Identifier: "random",
			},
		},
	})

	blockEvents := getEvents()

	hub.Publish(blockEvents)

	time.Sleep(time.Millisecond * 100)

	require.True(t, consumer2.HasEvent(blockEvents.Events[1]))
	require.True(t, consumer1.HasEvent(blockEvents.Events[0]))
	require.True(t, consumer3.HasEvent(blockEvents.Events[2]))
}

func TestCommonHub_HandleRevertBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	numCalls := uint32(0)
	hub.registerDispatcher(&mocks.DispatcherStub{
		RevertEventCalled: func(event data.RevertBlock) {
			atomic.AddUint32(&numCalls, 1)
		},
	})

	hub.Subscribe(data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.RevertBlockEvents,
			},
		},
	})

	blockEvents := data.RevertBlock{
		Hash:  "hash1",
		Nonce: 1,
	}

	hub.PublishRevert(blockEvents)

	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestCommonHub_HandleFinalizedBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	numCalls := uint32(0)
	hub.registerDispatcher(&mocks.DispatcherStub{
		FinalizedEventCalled: func(event data.FinalizedBlock) {
			atomic.AddUint32(&numCalls, 1)
		},
	})

	hub.Subscribe(data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.FinalizedBlockEvents,
			},
		},
	})

	blockEvents := data.FinalizedBlock{
		Hash: "hash1",
	}

	hub.PublishFinalized(blockEvents)

	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestCommonHub_HandleTxsBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.NoError(t, err)

	numCalls := uint32(0)
	hub.registerDispatcher(&mocks.DispatcherStub{
		TxsEventCalled: func(event data.BlockTxs) {
			atomic.AddUint32(&numCalls, 1)
		},
	})

	hub.Subscribe(data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockTxs,
			},
		},
	})

	blockEvents := data.BlockTxs{
		Hash: "hash1",
	}

	hub.PublishTxs(blockEvents)

	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestCommonHub_HandleBlockEventsBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.NoError(t, err)

	numCalls := uint32(0)
	hub.registerDispatcher(&mocks.DispatcherStub{
		BlockEventsCalled: func(event data.BlockEventsWithOrder) {
			atomic.AddUint32(&numCalls, 1)
		},
	})

	hub.Subscribe(data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockEvents,
			},
		},
	})

	blockEvents := data.BlockEventsWithOrder{
		Hash: "hash1",
	}

	hub.PublishBlockEventsWithOrder(blockEvents)

	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func TestCommonHub_HandleScrsBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	numCalls := uint32(0)
	hub.registerDispatcher(&mocks.DispatcherStub{
		ScrsEventCalled: func(event data.BlockScrs) {
			atomic.AddUint32(&numCalls, 1)
		},
	})

	hub.Subscribe(data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockScrs,
			},
		},
	})

	blockEvents := data.BlockScrs{
		Hash: "hash1",
	}

	hub.PublishScrs(blockEvents)

	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
}

func getEvents() data.BlockEvents {
	return data.BlockEvents{
		Hash: "374d75573060d840257045add9cd104b70180065f2406808ebabe02a1a3cb5f8",
		Events: []data.Event{
			{
				Address:    "erd1",
				Identifier: "swap",
				Topics:     [][]byte{},
				Data:       []byte{},
			},
			{
				Address:    "erd2",
				Identifier: "lock",
				Topics:     [][]byte{},
				Data:       []byte{},
			},
			{
				Address:    "erd3",
				Identifier: "random",
				Topics:     [][]byte{},
				Data:       []byte{},
			},
		},
	}
}
