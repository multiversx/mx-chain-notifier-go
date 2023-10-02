package preprocess_test

import (
	"encoding/json"
	"errors"
	"testing"

	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/data"
	notifierData "github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
	"github.com/stretchr/testify/require"
)

func TestPreProcessorV0_SaveBlock(t *testing.T) {
	t.Parallel()

	t.Run("nil block data", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlockV0()
		outportBlock.HeaderType = "invalid"
		marshalledBlock, _ := json.Marshal(outportBlock)

		dp, err := preprocess.NewEventsPreProcessorV0(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, coreData.ErrInvalidHeaderType, err)
	})

	t.Run("failed to handle push events", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()

		expectedErr := errors.New("exp error")
		args.Facade = &mocks.FacadeStub{
			HandlePushEventsV2Called: func(events data.ArgsSaveBlockData) error {
				return expectedErr
			},
		}

		outportBlock := createDefaultOutportBlockV0()

		dp, err := preprocess.NewEventsPreProcessorV0(args)
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, expectedErr, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := preprocess.NewEventsPreProcessorV0(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		outportBlock := createDefaultOutportBlockV0()

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Nil(t, err)
	})
}

func TestPreProcessorV0_RevertIndexerBlock(t *testing.T) {
	t.Parallel()

	blockData := &notifierData.RevertBlock{
		Hash:  "hash1",
		Nonce: 1,
		Round: 1,
		Epoch: 1,
	}

	dp, err := preprocess.NewEventsPreProcessorV0(createMockEventsDataPreProcessorArgs())
	require.Nil(t, err)

	marshalledBlock, _ := json.Marshal(blockData)
	err = dp.RevertIndexedBlock(marshalledBlock)
	require.Nil(t, err)
}

func TestPreProcessorV0_FinalizedBlock(t *testing.T) {
	t.Parallel()

	finalizedBlock := &data.FinalizedBlock{
		Hash: "headerHash1",
	}

	dp, err := preprocess.NewEventsPreProcessorV0(createMockEventsDataPreProcessorArgs())
	require.Nil(t, err)

	marshalledBlock, _ := json.Marshal(finalizedBlock)
	err = dp.FinalizedBlock(marshalledBlock)
	require.Nil(t, err)
}

func createDefaultOutportBlockV0() *notifierData.ArgsSaveBlock {
	outportBlock := &notifierData.ArgsSaveBlock{
		HeaderType: "Header",
		ArgsSaveBlockData: notifierData.ArgsSaveBlockData{
			HeaderHash: []byte("headerHash3"),
			Body: &block.Body{
				MiniBlocks: []*block.MiniBlock{
					{
						TxHashes:        [][]byte{},
						ReceiverShardID: 1,
						SenderShardID:   1,
					},
				},
			},
			TransactionsPool: &outport.TransactionPool{
				Transactions: map[string]*outport.TxInfo{
					"txHash1": {
						Transaction: &transaction.Transaction{
							Nonce:    1,
							GasPrice: 1,
							GasLimit: 1,
						},
						FeeInfo: &outport.FeeInfo{
							GasUsed: 1,
						},
						ExecutionOrder: 2,
					},
				},
				SmartContractResults: map[string]*outport.SCRInfo{
					"scrHash1": {
						SmartContractResult: &smartContractResult.SmartContractResult{
							Nonce:    2,
							GasLimit: 2,
							GasPrice: 2,
							CallType: 2,
						},
						FeeInfo: &outport.FeeInfo{
							GasUsed: 2,
						},
						ExecutionOrder: 0,
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
			},
			NumberOfShards: 2,
		},
	}

	return outportBlock
}
