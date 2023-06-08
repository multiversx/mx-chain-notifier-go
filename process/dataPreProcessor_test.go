package process_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/mock"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/stretchr/testify/require"
)

func createMockEventsDataPreProcessorArgs() process.ArgsEventsDataPreProcessor {
	return process.ArgsEventsDataPreProcessor{
		Marshaller:         &mock.MarshalizerMock{},
		InternalMarshaller: &mock.MarshalizerMock{},
		Facade:             &mocks.FacadeStub{},
	}
}

func TestNewEventsDataPreProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()
		args.Marshaller = nil

		dp, err := process.NewEventsDataPreProcessor(args)
		require.Nil(t, dp)
		require.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("nil internal marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()
		args.InternalMarshaller = nil

		dp, err := process.NewEventsDataPreProcessor(args)
		require.Nil(t, dp)
		require.True(t, errors.Is(err, common.ErrNilMarshaller))
	})

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()
		args.Facade = nil

		dp, err := process.NewEventsDataPreProcessor(args)
		require.Nil(t, dp)
		require.Equal(t, common.ErrNilFacadeHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)
		require.NotNil(t, dp)
		require.False(t, dp.IsInterfaceNil())
	})
}

func TestSaveBlock(t *testing.T) {
	t.Parallel()

	t.Run("nil block data", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.BlockData = nil

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.SaveBlock(outportBlock)
		require.Equal(t, process.ErrNilBlockData, err)
	})

	t.Run("nil transaction pool", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.TransactionPool = nil

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.SaveBlock(outportBlock)
		require.Equal(t, process.ErrNilTransactionPool, err)
	})

	t.Run("nil header gas consumption", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.HeaderGasConsumption = nil

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.SaveBlock(outportBlock)
		require.Equal(t, process.ErrNilHeaderGasConsumption, err)
	})

	t.Run("failed to get header from bytes", func(t *testing.T) {
		t.Parallel()

		outportBlock := createDefaultOutportBlock()
		outportBlock.BlockData.HeaderType = "invalid"

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.SaveBlock(outportBlock)
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

		outportBlock := createDefaultOutportBlock()

		dp, err := process.NewEventsDataPreProcessor(args)
		require.Nil(t, err)

		err = dp.SaveBlock(outportBlock)
		require.Equal(t, expectedErr, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		outportBlock := createDefaultOutportBlock()

		err = dp.SaveBlock(outportBlock)
		require.Nil(t, err)
	})
}

func TestRevertIndexerBlock(t *testing.T) {
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

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.RevertIndexedBlock(blockData)
		require.Equal(t, coreData.ErrInvalidHeaderType, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		b := &block.Header{
			Nonce: 1,
		}
		blockBytes, _ := json.Marshal(b)

		blockData := &outport.BlockData{
			HeaderBytes: blockBytes,
			HeaderType:  "Header",
		}

		dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
		require.Nil(t, err)

		err = dp.RevertIndexedBlock(blockData)
		require.Nil(t, err)
	})
}

func TestFinalizedBlock(t *testing.T) {
	t.Parallel()

	finalizedBlock := &outport.FinalizedBlock{
		HeaderHash: []byte("headerHash1"),
	}

	dp, err := process.NewEventsDataPreProcessor(createMockEventsDataPreProcessorArgs())
	require.Nil(t, err)

	err = dp.FinalizedBlock(finalizedBlock)
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
