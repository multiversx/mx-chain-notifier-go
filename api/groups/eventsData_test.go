package groups_test

import (
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	coreData "github.com/multiversx/mx-chain-core-go/websocketOutportDriver/data"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/stretchr/testify/require"
)

func TestUnmarshallBlockDataV1(t *testing.T) {
	t.Parallel()

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

	retBlockEvents, err := groups.UnmarshallBlockDataV1(jsonBytes)
	require.Nil(t, err)
	require.Equal(t, blockEvents, retBlockEvents)
}

func TestUnmarshallBlockDataV2(t *testing.T) {
	t.Parallel()

	argsSaveBlockData := data.ArgsSaveBlockData{
		HeaderHash: []byte{},
		Body:       &block.Body{},
		Header:     &block.HeaderV2{},
		TransactionsPool: &data.TransactionsPool{
			Txs: map[string]*data.NodeTransaction{
				"hash2": {
					TransactionHandler: &transaction.Transaction{
						Nonce: 2,
					},
					ExecutionOrder: 1,
				},
			},
			Scrs: map[string]*data.NodeSmartContractResult{
				"hash3": {
					TransactionHandler: &smartContractResult.SmartContractResult{
						Nonce: 3,
					},
					ExecutionOrder: 1,
				},
			},
			Logs: []*data.LogData{
				{
					LogHandler: &transaction.Log{
						Address: []byte("logaddr1"),
						Events:  []*transaction.Event{},
					},
					TxHash: "logHash1",
				},
			},
		},
	}
	blockEvents := data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: argsSaveBlockData,
	}
	jsonBytes, _ := json.Marshal(blockEvents)

	retBlockEvents, err := groups.UnmarshallBlockDataV2(jsonBytes)
	require.Nil(t, err)
	require.Equal(t, &argsSaveBlockData, retBlockEvents)
}

func TestGetHeader(t *testing.T) {
	t.Parallel()

	t.Run("header v2", func(t *testing.T) {
		t.Parallel()

		header := &block.HeaderV2{
			Header: &block.Header{
				Nonce:     1,
				ShardID:   1,
				TimeStamp: 2,
			},
			ScheduledGasProvided:  1,
			ScheduledGasPenalized: 1,
			ScheduledGasRefunded:  1,
		}
		saveBlockData := &coreData.ArgsSaveBlock{
			HeaderType: "HeaderV2",
			ArgsSaveBlockData: outport.ArgsSaveBlockData{
				Header: header,
			},
		}
		headerBytes, _ := json.Marshal(saveBlockData)

		h, err := groups.GetHeader(headerBytes)
		require.Nil(t, err)
		require.Equal(t, header, h)
	})

	t.Run("header v1", func(t *testing.T) {
		t.Parallel()

		header := &block.Header{
			Nonce:     1,
			ShardID:   1,
			TimeStamp: 1,
		}
		saveBlockData := &coreData.ArgsSaveBlock{
			HeaderType: "Header",
			ArgsSaveBlockData: outport.ArgsSaveBlockData{
				Header: header,
			},
		}
		headerBytes, _ := json.Marshal(saveBlockData)

		h, err := groups.GetHeader(headerBytes)
		require.Nil(t, err)
		require.Equal(t, header, h)
	})

	t.Run("metablock", func(t *testing.T) {
		t.Parallel()

		header := &block.MetaBlock{
			Nonce:     2,
			TimeStamp: 2,
		}
		saveBlockData := &coreData.ArgsSaveBlock{
			HeaderType: "MetaBlock",
			ArgsSaveBlockData: outport.ArgsSaveBlockData{
				Header: header,
			},
		}
		headerBytes, _ := json.Marshal(saveBlockData)

		h, err := groups.GetHeader(headerBytes)
		require.Nil(t, err)
		require.Equal(t, header, h)
	})
}
