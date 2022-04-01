package hub

import (
	"testing"

	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/ElrondNetwork/notifier-go/mocks"
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
		require.NotNil(t, hub)
	})
}

func TestCommonHub_RegisterDispatcher(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	dispatcher1 := mocks.NewDispatcherMock(nil, hub)
	dispatcher2 := mocks.NewDispatcherMock(nil, hub)

	hub.registerDispatcher(dispatcher1)
	require.True(t, hub.dispatchers[dispatcher1.GetID()] == dispatcher1)

	hub.registerDispatcher(dispatcher2)
	require.True(t, hub.dispatchers[dispatcher2.GetID()] == dispatcher2)
}

func TestCommonHub_UnregisterDispatcher(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	dispatcher1 := mocks.NewDispatcherMock(nil, hub)
	dispatcher2 := mocks.NewDispatcherMock(nil, hub)

	hub.registerDispatcher(dispatcher1)
	hub.registerDispatcher(dispatcher2)

	hub.unregisterDispatcher(dispatcher1)

	require.True(t, hub.dispatchers[dispatcher1.GetID()] == nil)
	require.True(t, hub.dispatchers[dispatcher2.GetID()] == dispatcher2)
}

func TestCommonHub_HandleBroadcastDispatcherReceivesEvents(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	consumer := mocks.NewConsumerMock()
	dispatcher1 := mocks.NewDispatcherMock(consumer, hub)

	hub.registerDispatcher(dispatcher1)
	hub.Subscribe(data.SubscribeEvent{
		DispatcherID:        dispatcher1.GetID(),
		SubscriptionEntries: []data.SubscriptionEntry{},
	})

	blockEvents := getEvents()

	hub.handleBroadcast(blockEvents)

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

	hub.registerDispatcher(dispatcher1)
	hub.registerDispatcher(dispatcher2)

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

	blockEvents := getEvents()

	hub.handleBroadcast(blockEvents)

	require.True(t, consumer1.HasEvent(blockEvents.Events[0]))
	require.True(t, consumer2.HasEvent(blockEvents.Events[1]))
}

func TestCommonHubRun(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	hub.Run()
	err = hub.Close()
	require.Nil(t, err)
}

func TestCommonHub_HandleRevertBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	wasCalled := false
	hub.registerDispatcher(&mocks.DispatcherStub{
		RevertEventCalled: func(event data.RevertBlock) {
			wasCalled = true
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

	hub.handleRevertBroadcast(blockEvents)

	assert.True(t, wasCalled)
}

func TestCommonHub_HandleFinalizedBroadcast(t *testing.T) {
	t.Parallel()

	args := createMockCommonHubArgs()
	hub, err := NewCommonHub(args)
	require.Nil(t, err)

	wasCalled := false
	hub.registerDispatcher(&mocks.DispatcherStub{
		FinalizedEventCalled: func(event data.FinalizedBlock) {
			wasCalled = true
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

	hub.handleFinalizedBroadcast(blockEvents)

	assert.True(t, wasCalled)
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
