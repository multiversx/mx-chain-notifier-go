package process

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
)

// ErrNilDataProcessor signals that a nil data processor has been provided
var ErrNilDataProcessor = errors.New("nil data processor")

// ErrInvalidPayloadType signals that an invalid payload type has been provided
var ErrInvalidPayloadType = errors.New("invalid payload type")

type payloadHandler struct {
	marshaller marshal.Marshalizer
	dp         DataProcessor
	actions    map[string]func(marshalledData []byte) error
}

// NewPayloadHandler will create a new instance of events indexer
func NewPayloadHandler(marshaller marshal.Marshalizer, dataProcessor DataProcessor) (*payloadHandler, error) {
	if check.IfNil(marshaller) {
		return nil, common.ErrNilMarshaller
	}
	if check.IfNil(dataProcessor) {
		return nil, ErrNilDataProcessor
	}

	payloadIndexer := &payloadHandler{
		marshaller: marshaller,
		dp:         dataProcessor,
	}
	payloadIndexer.initActionsMap()

	return payloadIndexer, nil
}

// GetOperationsMap returns the map with all the operations that will index data
func (ph *payloadHandler) initActionsMap() {
	ph.actions = map[string]func(d []byte) error{
		outport.TopicSaveBlock:             ph.saveBlock,
		outport.TopicRevertIndexedBlock:    ph.revertIndexedBlock,
		outport.TopicSaveRoundsInfo:        ph.saveRounds,
		outport.TopicSaveValidatorsRating:  ph.saveValidatorsRating,
		outport.TopicSaveValidatorsPubKeys: ph.saveValidatorsPubKeys,
		outport.TopicSaveAccounts:          ph.saveAccounts,
		outport.TopicFinalizedBlock:        ph.finalizedBlock,
	}
}

// ProcessPayload will proces the provided payload based on the topic
func (ph *payloadHandler) ProcessPayload(payload []byte, topic string, version string) error {
	payloadTypeAction, ok := ph.actions[topic]
	if !ok {
		log.Warn("invalid payload type", "topic", topic)
		return nil
	}

	return payloadTypeAction(payload)
}

func (ph *payloadHandler) saveBlock(marshalledData []byte) error {
	outportBlock := &outport.OutportBlock{}
	err := ph.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return err
	}

	return ph.dp.SaveBlock(outportBlock)
}

func (ph *payloadHandler) revertIndexedBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := ph.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	return ph.dp.RevertIndexedBlock(blockData)
}

func (ph *payloadHandler) finalizedBlock(marshalledData []byte) error {
	finalizedBlock := &outport.FinalizedBlock{}
	err := ph.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	return ph.dp.FinalizedBlock(finalizedBlock)
}

func (ph *payloadHandler) saveRounds(marshalledData []byte) error {
	return nil
}

func (ph *payloadHandler) saveValidatorsRating(marshalledData []byte) error {
	return nil
}

func (ph *payloadHandler) saveValidatorsPubKeys(marshalledData []byte) error {
	return nil
}

func (ph *payloadHandler) saveAccounts(marshalledData []byte) error {
	return nil
}

// Close will close the indexer
func (ph *payloadHandler) Close() error {
	return nil
}

// IsInterfaceNil returns true if underlying object is nil
func (ph *payloadHandler) IsInterfaceNil() bool {
	return ph == nil
}
