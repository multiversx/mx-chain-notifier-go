package process_test

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/stateChange"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/stretchr/testify/require"
)

func createMockEventsInterceptorArgs() process.ArgsEventsInterceptor {
	return process.ArgsEventsInterceptor{
		PubKeyConverter:      &mocks.PubkeyConverterMock{},
		WithReadStateChanges: false,
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

func TestEventsInterceptor_baseNilChecks(t *testing.T) {
	t.Parallel()

	t.Run("nil block events data", func(t *testing.T) {
		t.Parallel()

		err := process.BaseNilEventsDataCheks(nil)
		require.Equal(t, process.ErrNilBlockEvents, err)
	})

	t.Run("nil block body", func(t *testing.T) {
		t.Parallel()

		eventsData := &data.ArgsSaveBlockData{
			HeaderHash:       []byte("headerHash"),
			TransactionsPool: &outport.TransactionPool{},
			Body:             nil,
		}
		err := process.BaseNilEventsDataCheks(eventsData)
		require.Equal(t, process.ErrNilBlockBody, err)
	})

	t.Run("nil block header", func(t *testing.T) {
		t.Parallel()

		eventsData := &data.ArgsSaveBlockData{
			HeaderHash:       []byte("headerHash"),
			TransactionsPool: &outport.TransactionPool{},
			Body:             &block.Body{},
			Header:           nil,
		}
		err := process.BaseNilEventsDataCheks(eventsData)
		require.Equal(t, process.ErrNilBlockHeader, err)
	})
}

func TestProcessBlockEvents_WithoutExecutionResults(t *testing.T) {
	t.Parallel()

	t.Run("nil state accesses, should return empty map", func(t *testing.T) {
		t.Parallel()

		eventsInterceptor, _ := process.NewEventsInterceptor(createMockEventsInterceptorArgs())

		eventsData := &data.ArgsSaveBlockData{
			HeaderHash:       []byte("headerHash"),
			TransactionsPool: &outport.TransactionPool{},
			Body:             &block.Body{},
			Header:           &block.HeaderV2{},
			StateAccesses:    nil,
		}
		events, err := eventsInterceptor.ProcessBlockEvents(eventsData)
		require.Nil(t, err)

		expInterceptorData := &data.InterceptorBlockData{
			Hash:                     hex.EncodeToString([]byte("headerHash")),
			Body:                     &block.Body{},
			Header:                   &block.HeaderV2{},
			Txs:                      map[string]*transaction.Transaction{},
			TxsWithOrder:             map[string]*outport.TxInfo(nil),
			Scrs:                     map[string]*smartContractResult.SmartContractResult{},
			ScrsWithOrder:            map[string]*outport.SCRInfo(nil),
			LogEvents:                []data.Event{},
			StateAccessesPerAccounts: map[string]*stateChange.StateAccesses{},
		}

		require.Equal(t, expInterceptorData, events)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		eventsInterceptor, _ := process.NewEventsInterceptor(createMockEventsInterceptorArgs())

		txs := map[string]*outport.TxInfo{
			"hash2": {
				Transaction: &transaction.Transaction{
					Nonce: 2,
				},
				ExecutionOrder: 1,
			},
		}
		scrs := map[string]*outport.SCRInfo{
			"hash3": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 3,
				},
				ExecutionOrder: 1,
			},
		}
		addr := []byte("addr1")

		blockBody := &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		}
		blockHeader := &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
			},
		}

		logs := []*transaction.LogData{
			{
				Log: &transaction.Log{
					Address: addr,
					Events: []*transaction.Event{
						{
							Address: addr,
						},
					},
				},
			},
		}

		blockHash := []byte("blockHash")
		blockEvents := data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			Body:       blockBody,
			Header:     blockHeader,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				Logs:                 logs,
			},
			StateAccesses:         make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: make(map[string]*outport.StateAccessesForBlock),
		}

		expTxs := map[string]*transaction.Transaction{
			"hash2": {
				Nonce: 2,
			},
		}
		expTxsWithOrder := map[string]*outport.TxInfo{
			"hash2": {
				Transaction: &transaction.Transaction{
					Nonce: 2,
				},
				ExecutionOrder: 1,
			},
		}
		expScrs := map[string]*smartContractResult.SmartContractResult{
			"hash3": {
				Nonce: 3,
			},
		}
		expScrsWithOrder := map[string]*outport.SCRInfo{
			"hash3": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 3,
				},
				ExecutionOrder: 1,
			},
		}

		expEvents := &data.InterceptorBlockData{
			Hash:          hex.EncodeToString(blockHash),
			Body:          blockBody,
			Header:        blockHeader,
			Txs:           expTxs,
			TxsWithOrder:  expTxsWithOrder,
			Scrs:          expScrs,
			ScrsWithOrder: expScrsWithOrder,
			LogEvents: []data.Event{
				{
					Address:    hex.EncodeToString(addr),
					Identifier: "",
					Data:       make([]byte, 0),
					Topics:     make([][]byte, 0),
				},
			},
			StateAccessesPerAccounts: make(map[string]*stateChange.StateAccesses),
		}

		events, err := eventsInterceptor.ProcessBlockEvents(&blockEvents)
		require.Nil(t, err)
		require.Equal(t, expEvents, events)
	})

	t.Run("nil event fields should be returned as empty", func(t *testing.T) {
		t.Parallel()

		eventsInterceptor, _ := process.NewEventsInterceptor(createMockEventsInterceptorArgs())

		addr := []byte("addr1")

		blockBody := &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		}
		blockHeader := &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
			},
		}

		logs := []*transaction.LogData{
			{
				Log: &transaction.Log{
					Address: addr,
					Events: []*transaction.Event{
						{
							Address:    addr,
							Topics:     nil,
							Data:       nil,
							Identifier: nil,
						},
					},
				},
			},
		}

		blockHash := []byte("blockHash")
		blockEvents := data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			Body:       blockBody,
			Header:     blockHeader,
			TransactionsPool: &outport.TransactionPool{
				Logs: logs,
			},
			StateAccesses:         make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: make(map[string]*outport.StateAccessesForBlock),
		}

		expEvents := &data.InterceptorBlockData{
			Hash:   hex.EncodeToString(blockHash),
			Body:   blockBody,
			Header: blockHeader,
			Txs:    make(map[string]*transaction.Transaction),
			Scrs:   make(map[string]*smartContractResult.SmartContractResult),
			LogEvents: []data.Event{
				{
					Address:    hex.EncodeToString(addr),
					Identifier: "",
					Data:       make([]byte, 0),
					Topics:     make([][]byte, 0),
				},
			},
			StateAccessesPerAccounts: make(map[string]*stateChange.StateAccesses),
		}

		events, err := eventsInterceptor.ProcessBlockEvents(&blockEvents)
		require.Nil(t, err)
		require.Equal(t, expEvents, events)
	})
}

func TestProcessBlockEvents_WithExecutionResults(t *testing.T) {
	t.Parallel()

	// TODO: add test for state accesses when implemented on the node

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		eventsInterceptor, _ := process.NewEventsInterceptor(createMockEventsInterceptorArgs())

		txs := map[string]*outport.TxInfo{
			"hash2": {
				Transaction: &transaction.Transaction{
					Nonce: 2,
				},
				ExecutionOrder: 1,
			},
		}
		scrs := map[string]*outport.SCRInfo{
			"hash3": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 3,
				},
				ExecutionOrder: 1,
			},
		}
		addr := []byte("addr1")

		blockBody := &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		}
		blockHeader := &block.HeaderV3{
			ShardID:     1,
			TimestampMs: 1234,
		}

		logs := []*transaction.LogData{
			{
				Log: &transaction.Log{
					Address: addr,
					Events: []*transaction.Event{
						{
							Address: addr,
						},
					},
				},
			},
		}

		proposedTxPool := &outport.TransactionPool{
			Transactions:         nil,
			SmartContractResults: nil,
			Logs:                 nil,
		}

		execTxPool := &outport.TransactionPool{
			Transactions:         txs,
			SmartContractResults: scrs,
			Logs:                 logs,
		}

		blockHash := []byte("blockHash")

		execResults := map[string]*outport.ExecutionResultData{
			hex.EncodeToString(blockHash): {
				Body:            blockBody,
				TransactionPool: execTxPool,
			},
		}

		blockEvents := data.ArgsSaveBlockData{
			HeaderHash:            blockHash,
			Body:                  blockBody,
			Header:                blockHeader,
			TransactionsPool:      proposedTxPool,
			StateAccesses:         make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: make(map[string]*outport.StateAccessesForBlock),
			Results:               execResults,
		}

		expTxs := map[string]*transaction.Transaction{
			"hash2": {
				Nonce: 2,
			},
		}
		expTxsWithOrder := map[string]*outport.TxInfo{
			"hash2": {
				Transaction: &transaction.Transaction{
					Nonce: 2,
				},
				ExecutionOrder: 1,
			},
		}
		expScrs := map[string]*smartContractResult.SmartContractResult{
			"hash3": {
				Nonce: 3,
			},
		}
		expScrsWithOrder := map[string]*outport.SCRInfo{
			"hash3": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 3,
				},
				ExecutionOrder: 1,
			},
		}

		expEvents := []*data.InterceptorBlockData{
			{
				Hash:          hex.EncodeToString(blockHash),
				Body:          blockBody,
				Header:        blockHeader,
				Txs:           expTxs,
				TxsWithOrder:  expTxsWithOrder,
				Scrs:          expScrs,
				ScrsWithOrder: expScrsWithOrder,
				LogEvents: []data.Event{
					{
						Address:    hex.EncodeToString(addr),
						Identifier: "",
						Data:       make([]byte, 0),
						Topics:     make([][]byte, 0),
					},
				},
				StateAccessesPerAccounts: make(map[string]*stateChange.StateAccesses),
			},
		}

		events, err := eventsInterceptor.ProcessBlockEventsV3(&blockEvents)
		require.Nil(t, err)
		require.Equal(t, expEvents, events)
	})

	t.Run("nil event fields should be returned as empty", func(t *testing.T) {
		t.Parallel()

		eventsInterceptor, _ := process.NewEventsInterceptor(createMockEventsInterceptorArgs())

		addr := []byte("addr1")

		blockBody := &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		}
		blockHeader := &block.HeaderV3{
			ShardID:     1,
			TimestampMs: 1234,
		}

		logs := []*transaction.LogData{
			{
				Log: &transaction.Log{
					Address: addr,
					Events: []*transaction.Event{
						{
							Address:    addr,
							Topics:     nil,
							Data:       nil,
							Identifier: nil,
						},
					},
				},
			},
		}

		proposedTxPool := &outport.TransactionPool{
			Logs: nil,
		}

		execTxPool := &outport.TransactionPool{
			Logs: logs,
		}

		blockHash := []byte("blockHash")

		execResults := map[string]*outport.ExecutionResultData{
			hex.EncodeToString(blockHash): {
				Body:            blockBody,
				TransactionPool: execTxPool,
			},
		}

		blockEvents := data.ArgsSaveBlockData{
			HeaderHash:            blockHash,
			Body:                  blockBody,
			Header:                blockHeader,
			TransactionsPool:      proposedTxPool,
			StateAccesses:         make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: make(map[string]*outport.StateAccessesForBlock),
			Results:               execResults,
		}

		expEvents := []*data.InterceptorBlockData{
			{
				Hash:   hex.EncodeToString(blockHash),
				Body:   blockBody,
				Header: blockHeader,
				Txs:    make(map[string]*transaction.Transaction),
				Scrs:   make(map[string]*smartContractResult.SmartContractResult),
				LogEvents: []data.Event{
					{
						Address:    hex.EncodeToString(addr),
						Identifier: "",
						Data:       make([]byte, 0),
						Topics:     make([][]byte, 0),
					},
				},
				StateAccessesPerAccounts: make(map[string]*stateChange.StateAccesses),
			},
		}

		events, err := eventsInterceptor.ProcessBlockEventsV3(&blockEvents)
		require.Nil(t, err)
		require.Equal(t, expEvents, events)
	})
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

	logs := []*transaction.LogData{
		{
			Log: &transaction.Log{
				Events: []*transaction.Event{
					events[0],
					events[1],
				},
			},
			TxHash: txHash1,
		},
		{
			Log: &transaction.Log{
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
	require.Equal(t, txHash1, receivedEvents[0].TxHash)
	require.Equal(t, txHash1, receivedEvents[1].TxHash)
	require.Equal(t, txHash2, receivedEvents[2].TxHash)
}

func TestEventsInterceptor_GetStateAccessesPerAccounts(t *testing.T) {
	t.Parallel()

	txs := map[string]*outport.TxInfo{
		hex.EncodeToString([]byte("txHash1")): {
			Transaction: &transaction.Transaction{
				Nonce: 2,
			},
			ExecutionOrder: 1,
		},
	}
	scrs := map[string]*outport.SCRInfo{
		hex.EncodeToString([]byte("txHash2")): {
			SmartContractResult: &smartContractResult.SmartContractResult{
				Nonce: 3,
			},
			ExecutionOrder: 2,
		},
	}
	invalidTxs := map[string]*outport.TxInfo{
		hex.EncodeToString([]byte("txHash0")): {
			Transaction: &transaction.Transaction{
				Nonce: 1,
			},
			ExecutionOrder: 0,
		},
	}

	stateAccessesRead := make(map[string]*stateChange.StateAccesses)
	stateAccessesRead["txHash1"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:        stateChange.Read,
				MainTrieKey: []byte("mainTrieKey1"),
				MainTrieVal: []byte("mainTrieVal1"),
			},
			&stateChange.StateAccess{
				Type:        stateChange.Read,
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal2"),
			},
		},
	}
	stateAccessesRead["txHash2"] = &stateChange.StateAccesses{}
	stateAccessesRead["txHash0"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:        stateChange.Read,
				MainTrieKey: []byte("mainTrieKey3"),
				MainTrieVal: []byte("mainTrieVal3"),
			},
			&stateChange.StateAccess{
				Type:        stateChange.Read,
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal4"),
			},
		},
	}

	stateAccessesWrite := make(map[string]*stateChange.StateAccesses)
	stateAccessesWrite["txHash1"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:        stateChange.Write,
				MainTrieKey: []byte("mainTrieKey1"),
				MainTrieVal: []byte("mainTrieVal1"),
			},
			&stateChange.StateAccess{
				Type:        stateChange.Write,
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal2"),
			},
		},
	}
	stateAccessesWrite["txHash2"] = &stateChange.StateAccesses{}
	stateAccessesWrite["txHash0"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:        stateChange.Write,
				MainTrieKey: []byte("mainTrieKey3"),
				MainTrieVal: []byte("mainTrieVal3"),
			},
			&stateChange.StateAccess{
				Type:        stateChange.Write,
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal4"),
			},
		},
	}

	stateAccessesReadWrite := make(map[string]*stateChange.StateAccesses)
	stateAccessesReadWrite["txHash1"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:        stateChange.Read,
				MainTrieKey: []byte("mainTrieKey1"),
				MainTrieVal: []byte("mainTrieVal1"),
			},
			&stateChange.StateAccess{
				Type:        stateChange.Write,
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal2"),
			},
		},
	}
	stateAccessesReadWrite["txHash2"] = &stateChange.StateAccesses{}
	stateAccessesReadWrite["txHash0"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:        stateChange.Write,
				MainTrieKey: []byte("mainTrieKey3"),
				MainTrieVal: []byte("mainTrieVal3"),
			},
			&stateChange.StateAccess{
				Type:        stateChange.Read,
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal4"),
			},
		},
	}

	blockHash := []byte("blockHash")

	// before header V3

	t.Run("with write operations", func(t *testing.T) {
		t.Parallel()

		stateAccesses := make(map[string]*stateChange.StateAccesses)
		stateAccesses["txHash1"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		stateAccesses["txHash2"] = &stateChange.StateAccesses{}
		stateAccesses["txHash0"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
			},
		}

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: stateAccesses,
		}

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey1"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey2"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey3"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
			},
		}

		args := createMockEventsInterceptorArgs()
		en, _ := process.NewEventsInterceptor(args)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccounts(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	t.Run("with read operations, but not enabled from config", func(t *testing.T) {
		t.Parallel()

		stateAccesses := make(map[string]*stateChange.StateAccesses)
		stateAccesses["txHash1"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		stateAccesses["txHash2"] = &stateChange.StateAccesses{}
		stateAccesses["txHash0"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
			},
		}

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: stateAccesses,
		}

		args := createMockEventsInterceptorArgs()
		en, _ := process.NewEventsInterceptor(args)

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccounts(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	t.Run("with read (not enabled from config) and write operations", func(t *testing.T) {
		t.Parallel()

		stateAccesses := make(map[string]*stateChange.StateAccesses)
		stateAccesses["txHash1"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		stateAccesses["txHash2"] = &stateChange.StateAccesses{}
		stateAccesses["txHash0"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
			},
		}

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			Header:     &block.HeaderV3{},
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: stateAccesses,
		}

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey2"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey3"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
			},
		}

		args := createMockEventsInterceptorArgs()
		en, _ := process.NewEventsInterceptor(args)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccounts(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	t.Run("with read and write operations", func(t *testing.T) {
		t.Parallel()

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: stateAccessesReadWrite,
		}

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey1"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				&stateChange.StateAccess{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey2"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				&stateChange.StateAccess{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
				&stateChange.StateAccess{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey3"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				&stateChange.StateAccess{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
			},
		}

		args := createMockEventsInterceptorArgs()
		args.WithReadStateChanges = true
		en, _ := process.NewEventsInterceptor(args)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccounts(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	// with header V3

	t.Run("with write operations", func(t *testing.T) {
		t.Parallel()

		stateAccesses := make(map[string]*stateChange.StateAccesses)
		stateAccesses["txHash1"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		stateAccesses["txHash2"] = &stateChange.StateAccesses{}
		stateAccesses["txHash0"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
			},
		}

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: map[string]*outport.StateAccessesForBlock{
				hex.EncodeToString(blockHash): {stateAccesses},
			},
		}

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey1"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey2"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey3"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
			},
		}

		args := createMockEventsInterceptorArgs()
		en, _ := process.NewEventsInterceptor(args)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccountsV3(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	t.Run("with read operations, but not enabled from config", func(t *testing.T) {
		t.Parallel()

		stateAccesses := make(map[string]*stateChange.StateAccesses)
		stateAccesses["txHash1"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		stateAccesses["txHash2"] = &stateChange.StateAccesses{}
		stateAccesses["txHash0"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
			},
		}

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: map[string]*outport.StateAccessesForBlock{
				hex.EncodeToString(blockHash): {stateAccesses},
			},
		}

		args := createMockEventsInterceptorArgs()
		en, _ := process.NewEventsInterceptor(args)

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccountsV3(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	t.Run("with read (not enabled from config) and write operations", func(t *testing.T) {
		t.Parallel()

		stateAccesses := make(map[string]*stateChange.StateAccesses)
		stateAccesses["txHash1"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		stateAccesses["txHash2"] = &stateChange.StateAccesses{}
		stateAccesses["txHash0"] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
				{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
			},
		}

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			Header:     &block.HeaderV3{},
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: map[string]*outport.StateAccessesForBlock{
				hex.EncodeToString(blockHash): {stateAccesses},
			},
		}

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey2"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey3"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
			},
		}

		args := createMockEventsInterceptorArgs()
		en, _ := process.NewEventsInterceptor(args)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccountsV3(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})

	t.Run("with read and write operations", func(t *testing.T) {
		t.Parallel()

		blockEvents := &data.ArgsSaveBlockData{
			HeaderHash: blockHash,
			TransactionsPool: &outport.TransactionPool{
				Transactions:         txs,
				SmartContractResults: scrs,
				InvalidTxs:           invalidTxs,
			},
			StateAccesses: make(map[string]*stateChange.StateAccesses),
			StateAccessesForBlock: map[string]*outport.StateAccessesForBlock{
				hex.EncodeToString(blockHash): {stateAccessesReadWrite},
			},
		}

		expStateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey1"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				&stateChange.StateAccess{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey1"),
					MainTrieVal: []byte("mainTrieVal1"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey2"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				&stateChange.StateAccess{
					Type:        stateChange.Read,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal4"),
				},
				&stateChange.StateAccess{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey2"),
					MainTrieVal: []byte("mainTrieVal2"),
				},
			},
		}
		expStateAccessesPerAccounts[hex.EncodeToString([]byte("mainTrieKey3"))] = &stateChange.StateAccesses{
			StateAccess: []*stateChange.StateAccess{
				&stateChange.StateAccess{
					Type:        stateChange.Write,
					MainTrieKey: []byte("mainTrieKey3"),
					MainTrieVal: []byte("mainTrieVal3"),
				},
			},
		}

		args := createMockEventsInterceptorArgs()
		args.WithReadStateChanges = true
		en, _ := process.NewEventsInterceptor(args)

		stateAccessesPerAccounts := en.GetStateAccessesPerAccountsV3(blockEvents)

		require.Equal(t, expStateAccessesPerAccounts, stateAccessesPerAccounts)
	})
}

func TestEventsInterceptor_GetTxsWithOrder(t *testing.T) {
	t.Parallel()

	transactionPool := &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			"hash1": {
				ExecutionOrder: 0,
			},
			"hash2": {
				ExecutionOrder: 1,
			},
		},
		SmartContractResults: map[string]*outport.SCRInfo{
			"hash3": {
				ExecutionOrder: 2,
			},
		},
		Rewards: map[string]*outport.RewardInfo{
			"hash4": {
				ExecutionOrder: 3,
			},
		},
		InvalidTxs: map[string]*outport.TxInfo{
			"hash1": {
				ExecutionOrder: 0,
			},
		},
	}

	txsWithOrder := process.GetTxsWithOrder(transactionPool)

	// we expect one entry for each execution order 0..3,
	// with the duplicate "hash1" treated as invalid
	require.Len(t, txsWithOrder, 4)
	var hashes []string
	invalidCount := 0
	for _, tx := range txsWithOrder {
		hashes = append(hashes, tx.Hash)
		if tx.TxType == 3 { //invalid tx
			invalidCount++
		}
	}
	// execution order should be preserved: 0,1,2,3 -> hash1,hash2,hash3,hash4
	require.Equal(t, []string{"hash1", "hash2", "hash3", "hash4"}, hashes)
	// "hash1" should appear exactly once and be marked invalid
	hash1Count := 0
	for _, h := range hashes {
		if h == "hash1" {
			hash1Count++
		}
	}
	require.Equal(t, 1, hash1Count)
	require.Equal(t, invalidCount, 1)
}
