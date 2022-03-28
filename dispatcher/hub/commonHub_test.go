package hub

import (
	"testing"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/stretchr/testify/require"
)

func createMockCommonHubArgs() ArgsCommonHub {
	return ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
}

func makeHub() *commonHub {
	args := createMockCommonHubArgs()
	hub, _ := NewCommonHub(args)
	return hub
}

func TestCommonHub_RegisterDispatcher(t *testing.T) {
	t.Parallel()

	hub := makeHub()

	dispatcher1 := mocks.NewDispatcherMock(nil, hub)
	dispatcher2 := mocks.NewDispatcherMock(nil, hub)

	hub.registerDispatcher(dispatcher1)
	require.True(t, hub.dispatchers[dispatcher1.GetID()] == dispatcher1)

	hub.registerDispatcher(dispatcher2)
	require.True(t, hub.dispatchers[dispatcher2.GetID()] == dispatcher2)
}

func TestCommonHub_UnregisterDispatcher(t *testing.T) {
	t.Parallel()

	hub := makeHub()

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

	hub := makeHub()

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

	hub := makeHub()

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
