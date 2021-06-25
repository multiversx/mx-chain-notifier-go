package dispatcher

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubscriptionMap_Subscriptions(t *testing.T) {
	t.Parallel()

	subMap := NewSubscriptionMap()
	subMap.MatchSubscribeEvent(SubscribeEvent{SubscriptionEntries: []SubscriptionEntry{}})

	subs := subMap.Subscriptions()

	require.True(t, len(subs[MatchAll]) == 1)
}

func TestSubscriptionMap_MatchSubscribeEvent(t *testing.T) {
	t.Parallel()
	addr1 := "erd111"
	addr2 := "erd222"
	addr3 := "erd333"
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
				Address:    addr3,
				Identifier: "withdraw",
				Topics:     []string{"1", "2"},
			},
		},
	}

	subMap := NewSubscriptionMap()
	subMap.MatchSubscribeEvent(subEvent)

	subs := subMap.Subscriptions()

	sub1 := subs[addr1][0]
	require.NotEmpty(t, sub1)
	require.True(t, sub1.MatchLevel == MatchAddress)

	sub2 := subs[addr2][0]
	require.NotEmpty(t, sub2)
	require.True(t, sub2.MatchLevel == MatchIdentifier)

	sub3 := subs[addr3][0]
	require.NotEmpty(t, sub3)
	require.True(t, sub3.MatchLevel == MatchTopics)
}
