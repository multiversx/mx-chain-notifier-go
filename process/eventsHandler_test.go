package process_test

import (
	"context"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockEventsHandlerArgs() process.ArgsEventsHandler {
	return process.ArgsEventsHandler{
		Locker: &mocks.LockerStub{
			HasConnectionCalled: func(ctx context.Context) bool {
				return true
			},
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return true, nil
			},
		},
		Publisher:            &mocks.PublisherStub{},
		StatusMetricsHandler: &mocks.StatusMetricsStub{},
		EventsInterceptor:    &mocks.EventsInterceptorStub{},
	}
}

func TestNewEventsHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil locker service", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.Locker = nil

		eventsHandler, err := process.NewEventsHandler(args)
		require.Equal(t, process.ErrNilLockService, err)
		require.True(t, check.IfNil(eventsHandler))
	})

	t.Run("nil publisher", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.Publisher = nil

		eventsHandler, err := process.NewEventsHandler(args)
		require.Equal(t, process.ErrNilPublisherService, err)
		require.Nil(t, eventsHandler)
	})

	t.Run("nil status metrics handler", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.StatusMetricsHandler = nil

		eventsHandler, err := process.NewEventsHandler(args)
		require.Equal(t, common.ErrNilStatusMetricsHandler, err)
		require.Nil(t, eventsHandler)
	})

	t.Run("nil events interceptor", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.EventsInterceptor = nil

		eventsHandler, err := process.NewEventsHandler(args)
		require.Equal(t, process.ErrNilEventsInterceptor, err)
		require.Nil(t, eventsHandler)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)
		require.NotNil(t, eventsHandler)
	})
}

func TestHandleSaveBlockEvents(t *testing.T) {
	t.Parallel()

	t.Run("duplicated events, should return early", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true

		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		args.EventsInterceptor = &mocks.EventsInterceptorStub{
			ProcessBlockEventsCalled: func(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error) {
				require.Fail(t, "should have not been called")
				return nil, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		err = eventsHandler.HandleSaveBlockEvents(data.ArgsSaveBlockData{})
		require.Nil(t, err)
	})

	t.Run("failed to pre-process events, should fail", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true

		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return true, nil
			},
		}

		expectedErr := errors.New("expected err")
		args.EventsInterceptor = &mocks.EventsInterceptorStub{
			ProcessBlockEventsCalled: func(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error) {
				return nil, expectedErr
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		err = eventsHandler.HandleSaveBlockEvents(data.ArgsSaveBlockData{})
		require.Equal(t, expectedErr, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		blockHash := "blockHash1"
		txs := map[string]*outport.TxInfo{
			"hash1": {
				Transaction: &transaction.Transaction{
					Nonce: 1,
				},
				ExecutionOrder: 1,
			},
		}
		scrs := map[string]*outport.SCRInfo{
			"hash2": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 2,
				},
			},
		}
		logData := []*outport.LogData{
			{
				Log: &transaction.Log{
					Address: []byte("logaddr1"),
					Events:  []*transaction.Event{},
				},
				TxHash: "logHash1",
			},
		}

		logEvents := []data.Event{
			{
				Address: "addr1",
			},
		}

		header := &block.HeaderV2{
			Header: &block.Header{
				ShardID: 2,
			},
		}
		blockData := data.ArgsSaveBlockData{
			HeaderHash: []byte(blockHash),
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				Logs:                 logData,
			},
			Header: &block.HeaderV2{},
		}

		expTxs := map[string]*transaction.Transaction{
			"hash1": {
				Nonce: 1,
			},
		}
		expScrs := map[string]*smartContractResult.SmartContractResult{
			"hash2": {
				Nonce: 2,
			},
		}

		expTxsData := data.BlockTxs{
			Hash: blockHash,
			Txs:  expTxs,
		}
		expScrsData := data.BlockScrs{
			Hash: blockHash,
			Scrs: expScrs,
		}
		expLogEvents := data.BlockEvents{
			Hash:    blockHash,
			Events:  logEvents,
			ShardID: 2,
		}

		expTxsWithOrder := map[string]*outport.TxInfo{
			"hash1": {
				Transaction: &transaction.Transaction{
					Nonce: 1,
				},
				ExecutionOrder: 1,
			},
		}
		expScrsWithOrder := map[string]*outport.SCRInfo{
			"hash2": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 2,
				},
			},
		}
		expTxsWithOrderData := data.BlockEventsWithOrder{
			Hash:    blockHash,
			ShardID: 2,
			Txs:     expTxsWithOrder,
			Scrs:    expScrsWithOrder,
			Events:  logEvents,
		}

		pushWasCalled := false
		txsWasCalled := false
		scrsWasCalled := false
		blockEventsWithOrderWasCalled := false

		args := createMockEventsHandlerArgs()

		args.EventsInterceptor = &mocks.EventsInterceptorStub{
			ProcessBlockEventsCalled: func(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error) {
				return &data.InterceptorBlockData{
					Hash:          blockHash,
					Header:        header,
					Txs:           expTxs,
					Scrs:          expScrs,
					LogEvents:     logEvents,
					TxsWithOrder:  expTxsWithOrder,
					ScrsWithOrder: expScrsWithOrder,
				}, nil
			},
		}

		args.Publisher = &mocks.PublisherStub{
			BroadcastCalled: func(events data.BlockEvents) {
				pushWasCalled = true
				assert.Equal(t, expLogEvents, events)
			},
			BroadcastTxsCalled: func(event data.BlockTxs) {
				txsWasCalled = true
				assert.Equal(t, expTxsData, event)
			},
			BroadcastScrsCalled: func(event data.BlockScrs) {
				scrsWasCalled = true
				assert.Equal(t, expScrsData, event)
			},
			BroadcastBlockEventsWithOrderCalled: func(event data.BlockEventsWithOrder) {
				blockEventsWithOrderWasCalled = true
				assert.Equal(t, expTxsWithOrderData, event)
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		err = eventsHandler.HandleSaveBlockEvents(blockData)
		require.Nil(t, err)

		assert.True(t, pushWasCalled)
		assert.True(t, txsWasCalled)
		assert.True(t, scrsWasCalled)
		assert.True(t, blockEventsWithOrderWasCalled)
	})
}

func TestShouldProcessSaveBlockEvents(t *testing.T) {
	t.Parallel()

	t.Run("should process", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true

		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return true, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		shouldProcess := eventsHandler.ShouldProcessSaveBlockEvents("blockHash1")
		require.True(t, shouldProcess)
	})

	t.Run("duplicated events, should not process", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true

		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		shouldProcess := eventsHandler.ShouldProcessSaveBlockEvents("blockHash1")
		require.False(t, shouldProcess)
	})
}

func TestHandlePushEvents(t *testing.T) {
	t.Parallel()

	t.Run("empty hash should return error", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockEvents{
			Hash:      "",
			ShardID:   1,
			TimeStamp: 1234,
			Events:    []data.Event{},
		}

		err = eventsHandler.HandlePushEvents(events)
		require.Equal(t, common.ErrReceivedEmptyEvents, err)
	})

	t.Run("broadcast event was called", func(t *testing.T) {
		t.Parallel()

		events := data.BlockEvents{
			Hash:      "hash1",
			ShardID:   1,
			TimeStamp: 1234,
			Events:    []data.Event{},
		}

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastCalled: func(evs data.BlockEvents) {
				require.Equal(t, events, evs)
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		err = eventsHandler.HandlePushEvents(events)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})
}

func TestHandleRevertEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastRevertCalled: func(events data.RevertBlock) {
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.RevertBlock{
			Hash:  "hash1",
			Nonce: 1,
		}

		eventsHandler.HandleRevertEvents(events)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		revertEvents := data.RevertBlock{
			Hash:  "hash1",
			Nonce: 1,
		}

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastRevertCalled: func(events data.RevertBlock) {
				require.Equal(t, revertEvents, events)
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		eventsHandler.HandleRevertEvents(revertEvents)
		require.False(t, wasCalled)
	})
}

func TestHandleFinalizedEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast finalized event was called", func(t *testing.T) {
		t.Parallel()

		finalizedEvents := data.FinalizedBlock{
			Hash: "hash1",
		}

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastFinalizedCalled: func(events data.FinalizedBlock) {
				require.Equal(t, finalizedEvents, events)
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		eventsHandler.HandleFinalizedEvents(finalizedEvents)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastFinalizedCalled: func(events data.FinalizedBlock) {
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.FinalizedBlock{
			Hash: "hash1",
		}

		eventsHandler.HandleFinalizedEvents(events)
		require.False(t, wasCalled)
	})
}

func TestHandleTxsEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast txs event was called", func(t *testing.T) {
		t.Parallel()

		blockTxs := data.BlockTxs{
			Hash: "hash1",
			Txs: map[string]*transaction.Transaction{
				"hash1": {
					Nonce: 1,
				},
			},
		}

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastTxsCalled: func(event data.BlockTxs) {
				require.Equal(t, blockTxs, event)
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		eventsHandler.HandleBlockTxs(blockTxs)
		require.True(t, wasCalled)
	})
}

func TestHandleScrsEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast scrs event was called", func(t *testing.T) {
		t.Parallel()

		blockScrs := data.BlockScrs{
			Hash: "hash1",
			Scrs: map[string]*smartContractResult.SmartContractResult{
				"hash2": {
					Nonce: 2,
				},
			},
		}

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastScrsCalled: func(event data.BlockScrs) {
				require.Equal(t, blockScrs, event)
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		eventsHandler.HandleBlockScrs(blockScrs)
		require.True(t, wasCalled)
	})
}

func TestHandleBlockEventsWithOrderEvents(t *testing.T) {
	t.Parallel()

	events := data.BlockEventsWithOrder{
		Hash: "hash1",
		Txs: map[string]*outport.TxInfo{
			"hash1": {
				Transaction: &transaction.Transaction{
					Nonce: 1,
				},
				ExecutionOrder: 2,
			},
		},
	}

	t.Run("broadcast block events with order event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastBlockEventsWithOrderCalled: func(event data.BlockEventsWithOrder) {
				require.Equal(t, events, events)
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		eventsHandler.HandleBlockEventsWithOrder(events)
		require.True(t, wasCalled)
	})
}

func TestTryCheckProcessedWithRetry(t *testing.T) {
	t.Parallel()

	hash := "hash1"
	prefix := "prefix_"

	t.Run("event is NOT already processed", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		ok := eventsHandler.TryCheckProcessedWithRetry(prefix, hash)
		require.False(t, ok)
	})

	t.Run("event is already processed", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.CheckDuplicates = true
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return true, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		ok := eventsHandler.TryCheckProcessedWithRetry(prefix, hash)
		require.True(t, ok)
	})

	t.Run("locker service is failing on first try, has no connection, works on second try", func(t *testing.T) {
		t.Parallel()

		numCallsHasConnection := 0
		numCallsIsProcessed := 0

		args := createMockEventsHandlerArgs()
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				if numCallsIsProcessed > 0 {
					return true, nil
				}
				numCallsIsProcessed++
				return false, errors.New("fail to process")
			},
			HasConnectionCalled: func(ctx context.Context) bool {
				if numCallsHasConnection > 0 {
					return true
				}
				numCallsHasConnection++
				return false
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		ok := eventsHandler.TryCheckProcessedWithRetry(prefix, hash)
		require.True(t, ok)
	})
}
