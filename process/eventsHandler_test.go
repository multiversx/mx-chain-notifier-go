package process_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

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

func TestHandleSaveBlockEvents_ShouldFail(t *testing.T) {
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

		err = eventsHandler.HandleSaveBlockEvents(data.ArgsSaveBlockData{
			Header: &block.HeaderV2{},
		})
		require.Nil(t, err)
	})

	t.Run("nil events header, should fail", func(t *testing.T) {
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

		blockData := data.ArgsSaveBlockData{
			Header: nil,
		}

		err = eventsHandler.HandleSaveBlockEvents(blockData)
		require.Equal(t, process.ErrNilBlockHeader, err)
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

		blockData := data.ArgsSaveBlockData{
			Header: &block.HeaderV2{},
		}

		err = eventsHandler.HandleSaveBlockEvents(blockData)
		require.Equal(t, expectedErr, err)
	})
}

func TestHandleSaveBlockEvents_ShouldWork(t *testing.T) {
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
	logData := []*transaction.LogData{
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

	t.Run("should work before header v3", func(t *testing.T) {
		t.Parallel()

		header := &block.HeaderV2{
			Header: &block.Header{
				ShardID: 2,
			},
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

		blockData := data.ArgsSaveBlockData{
			HeaderHash: []byte(blockHash),
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				Logs:                 logData,
			},
			Header: header,
		}

		err = eventsHandler.HandleSaveBlockEvents(blockData)
		require.Nil(t, err)

		assert.True(t, pushWasCalled)
		assert.True(t, txsWasCalled)
		assert.True(t, scrsWasCalled)
		assert.True(t, blockEventsWithOrderWasCalled)
	})

	t.Run("should work with header v3", func(t *testing.T) {
		t.Parallel()

		header := &block.HeaderV3{
			ShardID: 2,
		}

		pushWasCalled := false
		txsWasCalled := false
		scrsWasCalled := false
		blockEventsWithOrderWasCalled := false

		args := createMockEventsHandlerArgs()

		args.EventsInterceptor = &mocks.EventsInterceptorStub{
			ProcessBlockEventsCalled: func(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error) {
				assert.Fail(t, "should have not been called")
				return &data.InterceptorBlockData{}, nil
			},
			ProcessBlockEventsV3Called: func(eventsData *data.ArgsSaveBlockData) ([]*data.InterceptorBlockData, error) {
				return []*data.InterceptorBlockData{
					{
						Hash:          blockHash,
						Header:        header,
						Txs:           expTxs,
						Scrs:          expScrs,
						LogEvents:     logEvents,
						TxsWithOrder:  expTxsWithOrder,
						ScrsWithOrder: expScrsWithOrder,
					},
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

		blockData := data.ArgsSaveBlockData{
			HeaderHash: []byte(blockHash),
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				Logs:                 logData,
			},
			Header: header,
		}

		err = eventsHandler.HandleSaveBlockEvents(blockData)
		require.Nil(t, err)

		assert.True(t, pushWasCalled)
		assert.True(t, txsWasCalled)
		assert.True(t, scrsWasCalled)
		assert.True(t, blockEventsWithOrderWasCalled)
	})
}

func TestHandleSaveBlockEventsV3_PartialFailure(t *testing.T) {
	t.Parallel()

	header := &block.HeaderV3{ShardID: 2}

	// this entry has a nil Header, which makes handleSaveBlockEvents fail
	// with ErrNilBlockHeader - simulating one execution block in the batch
	// erroring out
	failingBlock := &data.InterceptorBlockData{
		Hash:  "execHash1",
		Nonce: 1,
	}

	okLogEvents := []data.Event{
		{Address: "addr1"},
	}
	okBlock := &data.InterceptorBlockData{
		Hash:      "execHash2",
		Header:    header,
		LogEvents: okLogEvents,
		Nonce:     2,
	}

	expOkPushEvents := data.BlockEvents{
		Hash:    "execHash2",
		ShardID: 2,
		Events:  okLogEvents,
	}

	args := createMockEventsHandlerArgs()
	args.CheckDuplicates = true

	claimed := make(map[string]bool)
	args.Locker = &mocks.LockerStub{
		IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
			if claimed[blockHash] {
				return false, nil
			}
			claimed[blockHash] = true
			return true, nil
		},
		HasConnectionCalled: func(ctx context.Context) bool {
			return true
		},
	}

	args.EventsInterceptor = &mocks.EventsInterceptorStub{
		ProcessBlockEventsV3Called: func(eventsData *data.ArgsSaveBlockData) ([]*data.InterceptorBlockData, error) {
			// failingBlock is listed first, so a bug that aborts the loop on
			// the first error would never even attempt okBlock
			return []*data.InterceptorBlockData{failingBlock, okBlock}, nil
		},
	}

	pushCalls := 0
	args.Publisher = &mocks.PublisherStub{
		BroadcastCalled: func(events data.BlockEvents) {
			pushCalls++
			require.Equal(t, expOkPushEvents, events)
		},
	}

	eventsHandler, err := process.NewEventsHandler(args)
	require.Nil(t, err)

	blockData := data.ArgsSaveBlockData{
		HeaderHash: []byte("proposedHeaderHash"),
		Header:     header,
	}

	err = eventsHandler.HandleSaveBlockEvents(blockData)
	require.Equal(t, process.ErrNilBlockHeader, err)
	require.Equal(t, 0, pushCalls)

	err = eventsHandler.HandleSaveBlockEvents(blockData)
	require.Nil(t, err)
	require.Equal(t, 1, pushCalls)
}

func TestHandleSaveBlockEventsV3_ConcurrentDuplicateDeliveries_NoInterleaving(t *testing.T) {
	t.Parallel()

	header := &block.HeaderV3{ShardID: 1}

	const numExecResults = 5
	execResults := make([]*data.InterceptorBlockData, 0, numExecResults)
	for i := 1; i <= numExecResults; i++ {
		execResults = append(execResults, &data.InterceptorBlockData{
			Hash:      fmt.Sprintf("execHash%d", i),
			Header:    header,
			LogEvents: []data.Event{{Address: fmt.Sprintf("addr%d", i)}},
			Nonce:     uint64(i),
		})
	}

	args := createMockEventsHandlerArgs()
	args.CheckDuplicates = true

	var lockMu sync.Mutex
	locked := false

	var claimMu sync.Mutex
	claimed := make(map[string]bool)

	args.Locker = &mocks.LockerStub{
		TryLockCalled: func(ctx context.Context, key string) (bool, error) {
			lockMu.Lock()
			defer lockMu.Unlock()

			if locked {
				return false, nil
			}
			locked = true
			return true, nil
		},
		UnlockCalled: func(ctx context.Context, key string) error {
			lockMu.Lock()
			defer lockMu.Unlock()

			locked = false
			return nil
		},
		IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
			claimMu.Lock()
			defer claimMu.Unlock()

			if claimed[blockHash] {
				return false, nil
			}
			claimed[blockHash] = true
			return true, nil
		},
		HasConnectionCalled: func(ctx context.Context) bool {
			return true
		},
	}

	args.EventsInterceptor = &mocks.EventsInterceptorStub{
		ProcessBlockEventsV3Called: func(eventsData *data.ArgsSaveBlockData) ([]*data.InterceptorBlockData, error) {
			return execResults, nil
		},
	}

	var publishMu sync.Mutex
	var publishedNonces []uint64
	args.Publisher = &mocks.PublisherStub{
		BroadcastCalled: func(events data.BlockEvents) {
			// give a racing goroutine that (incorrectly) skipped the lock a
			// chance to interleave its own publishes here
			time.Sleep(time.Millisecond)

			publishMu.Lock()
			defer publishMu.Unlock()

			for _, execResult := range execResults {
				if execResult.Hash == events.Hash {
					publishedNonces = append(publishedNonces, execResult.Nonce)
				}
			}
		},
	}

	eventsHandler, err := process.NewEventsHandler(args)
	require.Nil(t, err)

	blockData := data.ArgsSaveBlockData{
		HeaderHash: []byte("proposedHeaderHash"),
		Header:     header,
	}

	const numConcurrentDeliveries = 5
	wg := &sync.WaitGroup{}
	wg.Add(numConcurrentDeliveries)
	for i := 0; i < numConcurrentDeliveries; i++ {
		go func() {
			defer wg.Done()
			_ = eventsHandler.HandleSaveBlockEvents(blockData)
		}()
	}
	wg.Wait()

	require.Equal(t, numExecResults, len(publishedNonces))
	for i, nonce := range publishedNonces {
		require.Equal(t, uint64(i+1), nonce)
	}
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
