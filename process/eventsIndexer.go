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

type eventsIndexer struct {
	marshaller marshal.Marshalizer
	dp         DataProcessor
	actions    map[string]func(marshalledData []byte) error
}

// NewEventsIndexer will create a new instance of events indexer
func NewEventsIndexer(marshaller marshal.Marshalizer, dataProcessor DataProcessor) (*eventsIndexer, error) {
	if check.IfNil(marshaller) {
		return nil, common.ErrNilMarshaller
	}
	if check.IfNil(dataProcessor) {
		return nil, ErrNilDataProcessor
	}

	payloadIndexer := &eventsIndexer{
		marshaller: marshaller,
		dp:         dataProcessor,
	}
	payloadIndexer.initActionsMap()

	return payloadIndexer, nil
}

// GetOperationsMap returns the map with all the operations that will index data
func (ei *eventsIndexer) initActionsMap() {
	ei.actions = map[string]func(d []byte) error{
		outport.TopicSaveBlock:             ei.saveBlock,
		outport.TopicRevertIndexedBlock:    ei.revertIndexedBlock,
		outport.TopicSaveRoundsInfo:        ei.saveRounds,
		outport.TopicSaveValidatorsRating:  ei.saveValidatorsRating,
		outport.TopicSaveValidatorsPubKeys: ei.saveValidatorsPubKeys,
		outport.TopicSaveAccounts:          ei.saveAccounts,
		outport.TopicFinalizedBlock:        ei.finalizedBlock,
	}
}

// ProcessPayload will proces the provided payload based on the topic
func (ei *eventsIndexer) ProcessPayload(payload []byte, topic string) error {
	payloadTypeAction, ok := ei.actions[topic]
	if !ok {
		log.Warn("invalid payload type", "topic", topic)
		return ErrInvalidPayloadType
	}

	return payloadTypeAction(payload)
}

func (ei *eventsIndexer) saveBlock(marshalledData []byte) error {
	outportBlock := &outport.OutportBlock{}
	err := ei.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return err
	}

	return ei.dp.SaveBlock(outportBlock)
}

func (ei *eventsIndexer) revertIndexedBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := ei.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	return ei.dp.RevertIndexedBlock(blockData)
}

func (ei *eventsIndexer) finalizedBlock(marshalledData []byte) error {
	finalizedBlock := &outport.FinalizedBlock{}
	err := ei.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	return ei.dp.FinalizedBlock(finalizedBlock)
}

func (ei *eventsIndexer) saveRounds(marshalledData []byte) error {
	return nil
}

func (ei *eventsIndexer) saveValidatorsRating(marshalledData []byte) error {
	return nil
}

func (ei *eventsIndexer) saveValidatorsPubKeys(marshalledData []byte) error {
	return nil
}

func (ei *eventsIndexer) saveAccounts(marshalledData []byte) error {
	return nil
}

// Close will close the indexer
func (ei *eventsIndexer) Close() error {
	return nil
}

// IsInterfaceNil returns true if underlying object is nil
func (ei *eventsIndexer) IsInterfaceNil() bool {
	return ei == nil
}
