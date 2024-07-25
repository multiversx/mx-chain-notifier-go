package preprocess_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
	"github.com/stretchr/testify/require"
)

func TestPreProcessorV1_SaveBlock(t *testing.T) {
	t.Parallel()

	t.Run("nil block data", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.BlockData = nil
		marshalledBlock, _ := json.Marshal(outportBlock)

		dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, preprocess.ErrNilBlockData, err)
	})

	t.Run("nil transaction pool", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.TransactionPool = nil

		dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, preprocess.ErrNilTransactionPool, err)
	})

	t.Run("nil header gas consumption", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.HeaderGasConsumption = nil

		dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, preprocess.ErrNilHeaderGasConsumption, err)
	})

	t.Run("failed to get header from bytes", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.BlockData.HeaderType = "invalid"

		dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, coreData.ErrInvalidHeaderType, err)
	})

	t.Run("failed to handle push events", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()

		expectedErr := errors.New("exp error")
		args.Facade = &mocks.FacadeStub{
			HandlePushEventsCalled: func(events data.ArgsSaveBlockData) error {
				return expectedErr
			},
		}

		outportBlock := createDefaultOutportBlock()

		dp, err := preprocess.NewEventsPreProcessorV1(args)
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, expectedErr, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		outportBlock := createDefaultOutportBlock()

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Nil(t, err)
	})
}

func TestPreProcessorV1_RevertIndexerBlock(t *testing.T) {
	t.Parallel()

	t.Run("failed to get header from bytes, invalid header type", func(t *testing.T) {
		t.Parallel()

		b := &block.Header{
			Nonce: 1,
		}
		blockBytes, _ := json.Marshal(b)

		blockData := &outport.BlockData{
			HeaderBytes: blockBytes,
			HeaderType:  "invalid",
		}

		dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(blockData)
		err = dp.RevertIndexedBlock(marshalledBlock)
		require.Equal(t, coreData.ErrInvalidHeaderType, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		shardID := uint32(2)
		nonce := uint64(11)
		round := uint64(12)
		epoch := uint32(3)
		timestamp := uint64(1234)
		headerHash := []byte("headerHash1")

		b := &block.Header{
			Nonce:     nonce,
			Round:     round,
			Epoch:     epoch,
			TimeStamp: timestamp,
		}
		blockBytes, _ := json.Marshal(b)

		blockData := &outport.BlockData{
			HeaderBytes: blockBytes,
			HeaderHash:  headerHash,
			HeaderType:  "Header",
			ShardID:     shardID,
		}

		expRevertBlock := data.RevertBlock{
			Hash:      hex.EncodeToString(headerHash),
			Nonce:     nonce,
			Round:     round,
			Epoch:     epoch,
			ShardID:   shardID,
			TimeStamp: timestamp,
		}

		args := createMockEventsDataPreProcessorArgs()

		revertCalled := false
		args.Facade = &mocks.FacadeStub{
			HandleRevertEventsCalled: func(events data.RevertBlock) {
				revertCalled = true
				require.Equal(t, expRevertBlock, events)
			},
		}

		dp, err := preprocess.NewEventsPreProcessorV1(args)
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(blockData)
		err = dp.RevertIndexedBlock(marshalledBlock)
		require.Nil(t, err)

		require.True(t, revertCalled)
	})
}

func TestPreProcessorV1_FinalizedBlock(t *testing.T) {
	t.Parallel()

	finalizedBlock := &outport.FinalizedBlock{
		HeaderHash: []byte("headerHash1"),
	}

	dp, err := preprocess.NewEventsPreProcessorV1(createMockEventsDataPreProcessorArgs())
	require.Nil(t, err)

	marshalledBlock, _ := json.Marshal(finalizedBlock)
	err = dp.FinalizedBlock(marshalledBlock)
	require.Nil(t, err)
}

func createDefaultOutportBlock() *outport.OutportBlock {
	b := &block.Header{
		Nonce: 1,
	}
	blockBytes, _ := json.Marshal(b)

	outportBlock := &outport.OutportBlock{
		BlockData: &outport.BlockData{
			HeaderBytes: blockBytes,
			HeaderType:  "Header",
		},
		TransactionPool: &outport.TransactionPool{
			Transactions: map[string]*outport.TxInfo{
				"txHash1": {
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
				"scrHash1": {
					SmartContractResult: &smartContractResult.SmartContractResult{
						Nonce: 2,
					},
					FeeInfo: &outport.FeeInfo{
						GasUsed: 2,
					},
					ExecutionOrder: 2,
				},
			},
			Logs: []*outport.LogData{},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{
			GasProvided:    3,
			GasRefunded:    3,
			GasPenalized:   3,
			MaxGasPerBlock: 3,
		},
	}

	return outportBlock
}
