package dispatcher

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionMap_Subscriptions(t *testing.T) {
	t.Parallel()

	subMap := NewSubscriptionMapper()

	subEvents := generateSubscribeEvents(10)

	for _, subEvent := range subEvents {
		subMap.MatchSubscribeEvent(subEvent)
	}

	subs := subMap.Subscriptions()

	require.NotEmpty(t, subs[common.PushLogsAndEvents])
}

func TestSubscriptionsMap_ShouldMatchAllForEmptySubscriptionEntry(t *testing.T) {
	t.Parallel()

	entry := data.SubscriptionEntry{}

	subMap := NewSubscriptionMapper()

	require.True(t, subMap.matchLevelFromInput(entry) == MatchAll)
}

func TestSubscriptionMapper_MatchSubscribeEventResultsInCorrectSet(t *testing.T) {
	t.Parallel()

	subEvents := generateSubscribeEvents(10)

	subMap := NewSubscriptionMapper()

	for _, subEvent := range subEvents {
		subMap.MatchSubscribeEvent(subEvent)
	}

	var subsFromEntries []data.Subscription
	matchAllAdded := make(map[uuid.UUID]bool)
	for _, subEvent := range subEvents {
		if len(subEvent.SubscriptionEntries) == 0 && !matchAllAdded[subEvent.DispatcherID] {
			// identical subscriptions of the same dispatcher are stored only once
			matchAllAdded[subEvent.DispatcherID] = true
			subsFromEntries = append(subsFromEntries, data.Subscription{
				DispatcherID: subEvent.DispatcherID,
				MatchLevel:   MatchAll,
				EventType:    common.PushLogsAndEvents,
			})
		}
		for _, entry := range subEvent.SubscriptionEntries {
			subsFromEntries = append(subsFromEntries, data.Subscription{
				Address:      entry.Address,
				Identifier:   entry.Identifier,
				Topics:       entry.Topics,
				DispatcherID: subEvent.DispatcherID,
				MatchLevel:   subMap.matchLevelFromInput(entry),
				EventType:    entry.EventType,
			})
		}
	}

	subsFromMap := subMap.Subscriptions()

	require.True(t, len(subsFromMap[common.PushLogsAndEvents]) == len(subsFromEntries))

	for _, subs := range subsFromMap {
		for _, sub1 := range subs {
			found := false
			for _, sub2 := range subsFromEntries {
				if reflect.DeepEqual(sub1, sub2) {
					found = true
				}
			}
			require.True(t, found)
		}
	}
}

func TestSubscriptionMap_MatchSubscribeEventCorrectMatchLevel(t *testing.T) {
	t.Parallel()

	addr1 := "erd111"
	addr2 := "erd222"
	addr3 := "erd333"
	dispatcherId := uuid.New()
	subEvent := data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockEvents,
				Address:   addr1,
			},
			{
				EventType:  common.BlockEvents,
				Address:    addr2,
				Identifier: "ESDTTransfer",
			},
			{
				EventType:  common.BlockTxs,
				Identifier: "wrapEGLD",
			},
			{
				EventType:  common.BlockTxs,
				Address:    addr3,
				Identifier: "withdraw",
				Topics:     []string{"1", "2"},
			},
		},
		DispatcherID: dispatcherId,
	}

	subMap := NewSubscriptionMapper()
	subMap.MatchSubscribeEvent(subEvent)

	subs := subMap.Subscriptions()

	require.NotEmpty(t, subs[common.BlockEvents][0])
	require.True(t, subs[common.BlockEvents][0].MatchLevel == MatchAddress)

	require.NotEmpty(t, subs[common.BlockEvents][1])
	require.True(t, subs[common.BlockEvents][1].MatchLevel == MatchAddressIdentifier)

	require.NotEmpty(t, subs[common.BlockTxs][0])
	require.True(t, subs[common.BlockTxs][0].MatchLevel == MatchIdentifier)

	require.NotEmpty(t, subs[common.BlockTxs][1])
	require.True(t, subs[common.BlockTxs][1].MatchLevel == MatchTopics)
}

func TestSubscriptionMapper_RemoveSubscriptions(t *testing.T) {
	t.Parallel()

	subEvents := generateSubscribeEvents(10)

	subMap := NewSubscriptionMapper()

	for _, subEvent := range subEvents {
		subMap.MatchSubscribeEvent(subEvent)
	}

	rmDispatcherID := subEvents[2].DispatcherID

	subMap.RemoveSubscriptions(rmDispatcherID)

	subsFromMap := subMap.Subscriptions()

	for _, subs := range subsFromMap {
		for _, sub := range subs {
			if sub.DispatcherID == rmDispatcherID {
				t.Error("should clear all subscriptions")
			}
		}
	}
}

func TestSubscriptionMapper_IdenticalSubscriptionsAreStoredOnce(t *testing.T) {
	t.Parallel()

	dispatcherID := uuid.New()
	entry := data.SubscriptionEntry{Address: "erd1", Identifier: "id", Topics: []string{"t"}}
	subMap := NewSubscriptionMapper()

	subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID})
	subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID})
	subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID, SubscriptionEntries: []data.SubscriptionEntry{entry, entry, {}}})
	subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID, SubscriptionEntries: []data.SubscriptionEntry{entry}})

	// the empty entry is the same match-all subscription as the empty subscribe event
	require.Len(t, subMap.Subscriptions()[common.PushLogsAndEvents], 2)
}

func TestSubscriptionMapper_SubscriptionsPerDispatcherAreBounded(t *testing.T) {
	t.Parallel()

	newEntries := func(start, num int) []data.SubscriptionEntry {
		entries := make([]data.SubscriptionEntry, 0, num)
		for i := start; i < start+num; i++ {
			entries = append(entries, data.SubscriptionEntry{Identifier: fmt.Sprintf("id%d", i)})
		}
		return entries
	}

	t.Run("oversized subscribe event is rejected", func(t *testing.T) {
		t.Parallel()

		subMap := NewSubscriptionMapper()
		subMap.MatchSubscribeEvent(data.SubscribeEvent{
			DispatcherID:        uuid.New(),
			SubscriptionEntries: newEntries(0, MaxSubscriptionsPerDispatcher+1),
		})
		require.Empty(t, subMap.Subscriptions())
	})

	t.Run("limit is enforced across subscribe events", func(t *testing.T) {
		t.Parallel()

		dispatcherID := uuid.New()
		otherDispatcherID := uuid.New()
		subMap := NewSubscriptionMapper()

		subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID, SubscriptionEntries: newEntries(0, MaxSubscriptionsPerDispatcher-1)})
		require.Len(t, subMap.Subscriptions()[common.PushLogsAndEvents], MaxSubscriptionsPerDispatcher-1)

		// would exceed the limit: nothing from this event is stored
		subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID, SubscriptionEntries: newEntries(1000, 2)})
		require.Len(t, subMap.Subscriptions()[common.PushLogsAndEvents], MaxSubscriptionsPerDispatcher-1)

		for i := 0; i < 10; i++ {
			subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: dispatcherID, SubscriptionEntries: newEntries(2000+i, 1)})
		}
		require.Len(t, subMap.Subscriptions()[common.PushLogsAndEvents], MaxSubscriptionsPerDispatcher)

		// the limit is per dispatcher
		subMap.MatchSubscribeEvent(data.SubscribeEvent{DispatcherID: otherDispatcherID})
		require.Len(t, subMap.Subscriptions()[common.PushLogsAndEvents], MaxSubscriptionsPerDispatcher+1)
	})
}

func generateSubscribeEvents(num int) []data.SubscribeEvent {
	var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

	idsLen := num / 4
	var dispatcherIDs []uuid.UUID
	for i := 0; i < idsLen; i++ {
		dispatcherIDs = append(dispatcherIDs, uuid.New())
	}

	var events []data.SubscribeEvent
	for i := 0; i < num; i++ {
		dispatchID := dispatcherIDs[randSeed.Intn(idsLen)]
		numEntries := randSeed.Intn(3)

		var subEntries []data.SubscriptionEntry
		if numEntries > 0 {
			for entryIdx := 0; entryIdx < numEntries; entryIdx++ {
				var topics []string
				if randSeed.Intn(2) == 1 {
					topics = []string{randStr(12), randStr(60)}
				}

				entry := data.SubscriptionEntry{
					Address:    fmt.Sprintf("erd%s", randStr(30)),
					Identifier: randStr(12),
					Topics:     topics,
					EventType:  common.PushLogsAndEvents,
				}
				subEntries = append(subEntries, entry)
			}
		}
		event := data.SubscribeEvent{
			DispatcherID:        dispatchID,
			SubscriptionEntries: subEntries,
		}
		events = append(events, event)
	}
	return events
}

func randStr(length int) string {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randSeed.Intn(len(charset))]
	}
	return string(b)
}
