package hub

import (
	"testing"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/ElrondNetwork/notifier-go/test/mocks"
	"github.com/stretchr/testify/require"
)

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
	hub.Subscribe(dispatcher.SubscribeEvent{
		DispatcherID:        dispatcher1.GetID(),
		SubscriptionEntries: []dispatcher.SubscriptionEntry{},
	})

	events := getEvents()

	hub.handleBroadcast(events)

	require.True(t, consumer.HasEvents(events))
	require.True(t, len(consumer.CollectedEvents()) == len(events))
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

	hub.Subscribe(dispatcher.SubscribeEvent{
		DispatcherID: dispatcher1.GetID(),
		SubscriptionEntries: []dispatcher.SubscriptionEntry{
			{
				Address:    "erd1",
				Identifier: "swap",
			},
		},
	})
	hub.Subscribe(dispatcher.SubscribeEvent{
		DispatcherID: dispatcher2.GetID(),
		SubscriptionEntries: []dispatcher.SubscriptionEntry{
			{
				Address:    "erd2",
				Identifier: "lock",
			},
		},
	})

	events := getEvents()

	hub.handleBroadcast(events)

	require.True(t, consumer1.HasEvent(events[0]))
	require.True(t, consumer2.HasEvent(events[1]))
}

func getEvents() []data.Event {
	return []data.Event{
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
	}
}

func makeHub() *commonHub {
	return NewCommonHub(filters.NewDefaultFilter())
}
