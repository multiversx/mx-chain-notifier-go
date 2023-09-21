package process_test

import (
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-chain-communication-go/testscommon"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
	"github.com/stretchr/testify/require"
)

func createDefaultDataProcessors() map[uint32]process.DataProcessor {
	eventsProcessors := make(map[uint32]process.DataProcessor)

	dataPreProcessorArgs := preprocess.ArgsEventsPreProcessor{
		Marshaller: &mock.MarshalizerMock{},
		Facade:     &mocks.FacadeStub{},
	}

	eventsProcessorV0, _ := preprocess.NewEventsPreProcessorV0(dataPreProcessorArgs)
	eventsProcessors[common.PayloadV0] = eventsProcessorV0

	eventsProcessorV1, _ := preprocess.NewEventsPreProcessorV1(dataPreProcessorArgs)
	eventsProcessors[common.PayloadV1] = eventsProcessorV1

	return eventsProcessors
}

func TestNewPayloadHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		ei, err := process.NewPayloadHandler(nil, createDefaultDataProcessors())
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

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, createDefaultDataProcessors())
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

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, createDefaultDataProcessors())
		err = ei.ProcessPayload([]byte("payload"), "invalid topic", 1)
		require.Nil(t, err)
	})

	t.Run("save block", func(t *testing.T) {
		t.Parallel()

		wasCalled := false

		eventsProcessors := make(map[uint32]process.DataProcessor)
		dp := &mocks.EventsDataProcessorStub{
			SaveBlockCalled: func(marshalledData []byte) error {
				wasCalled = true
				return nil
			},
		}
		eventsProcessors[common.PayloadV1] = dp

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, eventsProcessors)
		require.Nil(t, err)

		err = ei.ProcessPayload([]byte("payload"), outport.TopicSaveBlock, common.PayloadV0)
		require.NotNil(t, err)

		outportBlock := &outport.OutportBlock{}
		outportBlock.BlockData = &outport.BlockData{HeaderHash: []byte("hash")}
		outportBlockBytes, err := marshaller.Marshal(outportBlock)
		require.Nil(t, err)

		err = ei.ProcessPayload(outportBlockBytes, outport.TopicSaveBlock, common.PayloadV1)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})

	t.Run("revert indexed block", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		eventsProcessors := make(map[uint32]process.DataProcessor)
		dp := &mocks.EventsDataProcessorStub{
			RevertIndexedBlockCalled: func(marshalledData []byte) error {
				wasCalled = true
				return nil
			},
		}
		eventsProcessors[common.PayloadV1] = dp

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, eventsProcessors)
		require.Nil(t, err)

		err = ei.ProcessPayload([]byte("payload"), outport.TopicRevertIndexedBlock, common.PayloadV0)
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

		err = ei.ProcessPayload(blockDataBytes, outport.TopicRevertIndexedBlock, common.PayloadV1)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})

	t.Run("finalized block", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		eventsProcessors := make(map[uint32]process.DataProcessor)
		dp := &mocks.EventsDataProcessorStub{
			FinalizedBlockCalled: func(marshalledData []byte) error {
				wasCalled = true
				return nil
			},
		}
		eventsProcessors[common.PayloadV1] = dp

		ei, err := process.NewPayloadHandler(&mock.MarshalizerMock{}, eventsProcessors)
		require.Nil(t, err)

		err = ei.ProcessPayload([]byte("payload"), outport.TopicFinalizedBlock, common.PayloadV0)
		require.NotNil(t, err)

		finalizedBlock := &outport.FinalizedBlock{
			HeaderHash: []byte("headerHash1"),
		}
		finalizedBlockBytes, err := marshaller.Marshal(finalizedBlock)
		require.Nil(t, err)

		err = ei.ProcessPayload(finalizedBlockBytes, outport.TopicFinalizedBlock, common.PayloadV1)
		require.Nil(t, err)
		require.True(t, wasCalled)
	})
}
