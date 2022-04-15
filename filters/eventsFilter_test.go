package filters

import (
	"testing"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/stretchr/testify/require"
)

var events = []data.Event{
	{
		Address:    "erd1",
		Identifier: "swap",
		Topics:     [][]byte{},
	},
	{
		Address:    "erd2",
		Identifier: "addLiquidity",
		Topics:     [][]byte{},
	},
	{
		Address:    "erd3",
		Identifier: "setValue",
		Topics:     [][]byte{},
	},
	{
		Address:    "erd4",
		Identifier: "getValue",
		Topics:     [][]byte{},
	},
}

var filter = NewDefaultFilter()

func TestDefaultFilter_MatchEventMatchAll(t *testing.T) {
	t.Parallel()

	s := data.Subscription{
		MatchLevel: dispatcher.MatchAll,
	}

	for _, e := range events {
		require.True(t, filter.MatchEvent(s, e))
	}
}

func TestDefaultFilter_MatchEventMatchAddress(t *testing.T) {
	t.Parallel()

	s := data.Subscription{
		Address:    "erd2",
		MatchLevel: dispatcher.MatchAddress,
	}

	require.True(t, filter.MatchEvent(s, events[1]))
}

func TestDefaultFilter_MatchEventMatchAddressIdentifier(t *testing.T) {
	t.Parallel()

	s := data.Subscription{
		Address:    "erd1",
		Identifier: "swap",
		MatchLevel: dispatcher.MatchAddressIdentifier,
	}

	require.True(t, filter.MatchEvent(s, events[0]))
}

func TestDefaultFilter_MatchEventMatchIdentifier(t *testing.T) {
	t.Parallel()

	s := data.Subscription{
		Identifier: "setValue",
		MatchLevel: dispatcher.MatchIdentifier,
	}

	require.True(t, filter.MatchEvent(s, events[2]))
}
