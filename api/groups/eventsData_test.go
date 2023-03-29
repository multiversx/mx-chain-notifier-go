package groups_test

import (
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/stretchr/testify/require"
)

func TestNewEventsDataHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		edh, err := groups.NewEventsDataHandler(nil)
		require.Nil(t, edh)
		require.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		edh, err := groups.NewEventsDataHandler(&mock.MarshalizerMock{})
		require.Nil(t, err)
		require.False(t, edh.IsInterfaceNil())
	})
}

func TestEventsDataHandler_UnmarshallBlockDataV1(t *testing.T) {
	t.Parallel()

	edh, err := groups.NewEventsDataHandler(&mock.MarshalizerMock{})
	require.Nil(t, err)

	blockEvents := &data.SaveBlockData{
		Hash: "hash1",
		Txs: map[string]*transaction.Transaction{
			"hash2": {
				Nonce: 2,
			},
		},
		Scrs: map[string]*smartContractResult.SmartContractResult{
			"hash3": {
				Nonce: 3,
			},
		},
		LogEvents: []data.Event{},
	}

	jsonBytes, _ := json.Marshal(blockEvents)

	retBlockEvents, err := edh.UnmarshallBlockDataV1(jsonBytes)
	require.Nil(t, err)
	require.Equal(t, blockEvents, retBlockEvents)
}

func TestEventsDataHandler_UnmarshallBlockData(t *testing.T) {
	t.Parallel()

	edh, err := groups.NewEventsDataHandler(&mock.MarshalizerMock{})
	require.Nil(t, err)

	header := &block.HeaderV2{
		Header: &block.Header{
			Nonce: 1,
		},
	}
	headerBytes, _ := json.Marshal(header)

	txPool := &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			"hash1": {
				Transaction: &transaction.Transaction{
					Nonce: 1,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 1,
				},
				ExecutionOrder: 1,
			},
		},
		SmartContractResults: map[string]*outport.SCRInfo{
			"hash2": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 2,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 2,
				},
				ExecutionOrder: 3,
			},
		},
		Logs: []*outport.LogData{
			{
				Log: &transaction.Log{
					Address: []byte("logaddr1"),
					Events:  []*transaction.Event{},
				},
				TxHash: "logHash1",
			},
		},
	}

	argsSaveBlockData := outport.OutportBlock{
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  []byte("headerHash"),
			Body: &block.Body{
				MiniBlocks: []*block.MiniBlock{
					{SenderShardID: 1},
				},
			},
		},
		TransactionPool:      txPool,
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
	}
	jsonBytes, _ := json.Marshal(argsSaveBlockData)

	expArgsSaveBlockData := &data.ArgsSaveBlockData{
		HeaderHash:           []byte("headerHash"),
		Header:               header,
		TransactionsPool:     txPool,
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
		Body: &block.Body{
			MiniBlocks: []*block.MiniBlock{
				{SenderShardID: 1},
			},
		},
	}

	retBlockEvents, err := edh.UnmarshallBlockData(jsonBytes)
	require.Nil(t, err)
	require.Equal(t, expArgsSaveBlockData, retBlockEvents)
}
