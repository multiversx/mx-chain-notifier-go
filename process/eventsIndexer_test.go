package process_test

import (
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-chain-communication-go/testscommon"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/stretchr/testify/require"
)

func TestNewEventsIndexer(t *testing.T) {
	t.Parallel()

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		ei, err := process.NewPayloadHandler(nil, &mocks.EventsDataProcessorStub{})
		require.Nil(t, ei)
		require.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("nil data processor", func(t *testing.T) {
		t.Parallel()

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, nil)
		require.Nil(t, ei)
		require.Equal(t, process.ErrNilDataProcessor, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, &mocks.EventsDataProcessorStub{})
		require.Nil(t, err)
		require.NotNil(t, ei)
		require.False(t, ei.IsInterfaceNil())
	})
}

func TestProcessPayload(t *testing.T) {
	t.Parallel()

	marshaller := &testscommon.MarshallerMock{}

	t.Run("invalid topic, should return nil", func(t *testing.T) {
		t.Parallel()

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, &mocks.EventsDataProcessorStub{})
		err = ei.ProcessPayload([]byte("payload"), "invalid topic")
		require.Nil(t, err)
	})

	t.Run("save block", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		dp := &mocks.EventsDataProcessorStub{
			SaveBlockCalled: func(outportBlock *outport.OutportBlock) error {
				wasCalled = true
				return nil
			},
		}

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, dp)
		require.Nil(t, err)

		err = ei.ProcessPayload([]byte("payload"), outport.TopicSaveBlock)
		require.NotNil(t, err)

		outportBlock := &outport.OutportBlock{}
		outportBlock.BlockData = &outport.BlockData{HeaderHash: []byte("hash")}
		outportBlockBytes, err := marshaller.Marshal(outportBlock)
		require.Nil(t, err)

		err = ei.ProcessPayload(outportBlockBytes, outport.TopicSaveBlock)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})

	t.Run("revert indexed block", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		dp := &mocks.EventsDataProcessorStub{
			RevertIndexedBlockCalled: func(blockData *data.RevertBlock) error {
				wasCalled = true
				return nil
			},
		}

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, dp)
		require.Nil(t, err)

		err = ei.ProcessPayload([]byte("payload"), outport.TopicRevertIndexedBlock)
		require.NotNil(t, err)

		b := &block.Header{
			Nonce: 1,
		}
		blockBytes, _ := json.Marshal(b)

		blockData := &outport.BlockData{
			HeaderBytes: blockBytes,
			HeaderType:  "Header",
		}
		blockDataBytes, err := marshaller.Marshal(blockData)
		require.Nil(t, err)

		err = ei.ProcessPayload(blockDataBytes, outport.TopicRevertIndexedBlock)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})

	t.Run("finalized block", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		dp := &mocks.EventsDataProcessorStub{
			FinalizedBlockCalled: func(finalizedBlock *outport.FinalizedBlock) error {
				wasCalled = true
				return nil
			},
		}

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, dp)
		require.Nil(t, err)

		err = ei.ProcessPayload([]byte("payload"), outport.TopicFinalizedBlock)
		require.NotNil(t, err)

		finalizedBlock := &outport.FinalizedBlock{
			HeaderHash: []byte("headerHash1"),
		}
		finalizedBlockBytes, err := marshaller.Marshal(finalizedBlock)
		require.Nil(t, err)

		err = ei.ProcessPayload(finalizedBlockBytes, outport.TopicFinalizedBlock)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})
}
