package preprocess_test

import (
	"encoding/json"
	"errors"
	"testing"

	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/data"
	notifierData "github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
	"github.com/multiversx/mx-chain-notifier-go/testdata"
	"github.com/stretchr/testify/require"
)

func TestPreProcessorV0_SaveBlock(t *testing.T) {
	t.Parallel()

	marshaller := &marshal.JsonMarshalizer{}
	blockData, err := testdata.NewBlockData(marshaller)
	require.Nil(t, err)

	t.Run("nil block data", func(t *testing.T) {
		t.Parallel()

		outportBlock := blockData.OutportBlockV0()
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
			HandlePushEventsCalled: func(events data.ArgsSaveBlockData) error {
				return expectedErr
			},
		}

		outportBlock := blockData.OutportBlockV0()

		dp, err := preprocess.NewEventsPreProcessorV0(args)
		require.Nil(t, err)

		marshalledBlock, _ := json.Marshal(outportBlock)
		err = dp.SaveBlock(marshalledBlock)
		require.Equal(t, expectedErr, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()

		wasCalled := false
		args.Facade = &mocks.FacadeStub{
			HandlePushEventsCalled: func(events data.ArgsSaveBlockData) error {
				wasCalled = true
				return nil
			},
		}

		dp, err := preprocess.NewEventsPreProcessorV0(args)
		require.Nil(t, err)

		outportBlock := blockData.OutportBlockV0()

		marshalledBlock, err := json.Marshal(outportBlock)
		require.Nil(t, err)
		err = dp.SaveBlock(marshalledBlock)
		require.Nil(t, err)

		require.True(t, wasCalled)
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
