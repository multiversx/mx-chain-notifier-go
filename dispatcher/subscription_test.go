package dispatcher

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/ElrondNetwork/notifier-go/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

func TestSubscriptionMap_Subscriptions(t *testing.T) {
	t.Parallel()

	subMap := NewSubscriptionMapper()

	subEvents := generateSubscribeEvents(1000)

	for _, subEvent := range subEvents {
		subMap.MatchSubscribeEvent(subEvent)
	}

	subs := subMap.Subscriptions()

	require.True(t, len(subs) >= len(subEvents))
}

func TestSubscriptionsMap_ShouldMatchAllForEmptySubscriptionEntry(t *testing.T) {
	t.Parallel()

	entry := SubscriptionEntry{}

	subMap := NewSubscriptionMapper()

	require.True(t, subMap.matchLevelFromInput(entry) == MatchAll)
}

func TestSubscriptionMapper_MatchSubscribeEventResultsInCorrectSet(t *testing.T) {
	t.Parallel()

	subEvents := generateSubscribeEvents(1000)

	subMap := NewSubscriptionMapper()

	for _, subEvent := range subEvents {
		subMap.MatchSubscribeEvent(subEvent)
	}

	var subsFromEntries []Subscription
	for _, subEvent := range subEvents {
		if len(subEvent.SubscriptionEntries) == 0 {
			subsFromEntries = append(subsFromEntries, Subscription{
				DispatcherID: subEvent.DispatcherID,
				MatchLevel:   MatchAll,
			})
		}
		for _, entry := range subEvent.SubscriptionEntries {
			subsFromEntries = append(subsFromEntries, Subscription{
				Address:      entry.Address,
				Identifier:   entry.Identifier,
				Topics:       entry.Topics,
				DispatcherID: subEvent.DispatcherID,
				MatchLevel:   subMap.matchLevelFromInput(entry),
			})
		}
	}

	subsFromMap := subMap.Subscriptions()

	require.True(t, len(subsFromMap) == len(subsFromEntries))

	for _, sub1 := range subsFromMap {
		found := false
		for _, sub2 := range subsFromEntries {
			if reflect.DeepEqual(sub1, sub2) {
				found = true
			}
		}
		require.True(t, found)
	}
}

func TestSubscriptionMap_MatchSubscribeEventCorrectMatchLevel(t *testing.T) {
	t.Parallel()

	addr1 := "erd111"
	addr2 := "erd222"
	addr3 := "erd333"
	dispatcherId := uuid.New()
	subEvent := SubscribeEvent{
		SubscriptionEntries: []SubscriptionEntry{
			{
				Address: addr1,
			},
			{
				Address:    addr2,
				Identifier: "ESDTTransfer",
			},
			{
				Identifier: "wrapEGLD",
			},
			{
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

	require.NotEmpty(t, subs[0])
	require.True(t, subs[0].MatchLevel == MatchAddress)

	require.NotEmpty(t, subs[1])
	require.True(t, subs[1].MatchLevel == MatchAddressIdentifier)

	require.NotEmpty(t, subs[2])
	require.True(t, subs[2].MatchLevel == MatchIdentifier)

	require.NotEmpty(t, subs[3])
	require.True(t, subs[3].MatchLevel == MatchTopics)
}

func TestSubscriptionMapper_RemoveSubscriptions(t *testing.T) {
	t.Parallel()

	subEvents := generateSubscribeEvents(10_000)

	subMap := NewSubscriptionMapper()

	for _, subEvent := range subEvents {
		subMap.MatchSubscribeEvent(subEvent)
	}

	rmDispatcherID := subEvents[1200].DispatcherID

	subMap.RemoveSubscriptions(rmDispatcherID)

	subsFromMap := subMap.Subscriptions()

	for _, sub := range subsFromMap {
		if sub.DispatcherID == rmDispatcherID {
			t.Error("should clear all subscriptions")
		}
	}
}

func generateSubscribeEvents(num int) []SubscribeEvent {
	idsLen := num / 4
	var dispatcherIDs []uuid.UUID
	for i := 0; i < idsLen; i++ {
		dispatcherIDs = append(dispatcherIDs, uuid.New())
	}

	var events []SubscribeEvent
	for i := 0; i < num; i++ {
		dispatchID := dispatcherIDs[randSeed.Intn(idsLen)]
		numEntries := randSeed.Intn(3)

		var subEntries []SubscriptionEntry
		if numEntries > 0 {
			for entryIdx := 0; entryIdx < numEntries; entryIdx++ {
				var topics []string
				if randSeed.Intn(2) == 1 {
					topics = []string{test.RandStr(12), test.RandStr(60)}
				}

				entry := SubscriptionEntry{
					Address:    fmt.Sprintf("erd%s", test.RandStr(30)),
					Identifier: test.RandStr(12),
					Topics:     topics,
				}
				subEntries = append(subEntries, entry)
			}
		}
		event := SubscribeEvent{
			DispatcherID:        dispatchID,
			SubscriptionEntries: subEntries,
		}
		events = append(events, event)
	}
	return events
}
