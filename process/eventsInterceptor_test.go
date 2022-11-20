package process_test

import (
	"encoding/hex"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/ElrondNetwork/notifier-go/process"
	"github.com/stretchr/testify/require"
)

func createMockEventsInterceptorArgs() process.ArgsEventsInterceptor {
	return process.ArgsEventsInterceptor{
		PubKeyConverter: &mocks.PubkeyConverterMock{},
	}
}

func TestNewEventsInterceptor(t *testing.T) {
	t.Parallel()

	t.Run("nil pub key converter", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsInterceptorArgs()
		args.PubKeyConverter = nil

		eventsInterceptor, err := process.NewEventsInterceptor(args)
		require.Nil(t, eventsInterceptor)
		require.Equal(t, process.ErrNilPubKeyConverter, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		eventsInterceptor, err := process.NewEventsInterceptor(createMockEventsInterceptorArgs())
		require.Nil(t, err)
		require.False(t, check.IfNil(eventsInterceptor))
	})
}

func TestProcessBlockEvents(t *testing.T) {
	t.Parallel()

	eventsInterceptor, _ := process.NewEventsInterceptor(createMockEventsInterceptorArgs())
}

func TestGetLogEventsFromTransactionsPool(t *testing.T) {
	t.Parallel()

	txHash1 := "txHash1"
	txHash2 := "txHash2"

	events := []*transaction.Event{
		{
			Address:    []byte("addr1"),
			Identifier: []byte("identifier1"),
		},
		{
			Address:    []byte("addr2"),
			Identifier: []byte("identifier2"),
		},
		{
			Address:    []byte("addr3"),
			Identifier: []byte("identifier3"),
		},
	}

	logs := []*data.LogData{
		{
			LogHandler: &transaction.Log{
				Events: []*transaction.Event{
					events[0],
					events[1],
				},
			},
			TxHash: txHash1,
		},
		{
			LogHandler: &transaction.Log{
				Events: []*transaction.Event{
					events[2],
				},
			},
			TxHash: txHash2,
		},
	}

	args := createMockEventsInterceptorArgs()
	en, _ := process.NewEventsInterceptor(args)

	receivedEvents := en.GetLogEventsFromTransactionsPool(logs)

	for i, event := range receivedEvents {
		require.Equal(t, hex.EncodeToString(events[i].Address), event.Address)
		require.Equal(t, string(events[i].Identifier), event.Identifier)
	}

	require.Equal(t, len(events), len(receivedEvents))
	require.Equal(t, hex.EncodeToString([]byte(txHash1)), receivedEvents[0].TxHash)
	require.Equal(t, hex.EncodeToString([]byte(txHash1)), receivedEvents[1].TxHash)
	require.Equal(t, hex.EncodeToString([]byte(txHash2)), receivedEvents[2].TxHash)
}
