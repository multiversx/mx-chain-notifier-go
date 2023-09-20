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

// ErrInvalidPayloadVersion signals that an invalid payload version has been provided
var ErrInvalidPayloadVersion = errors.New("invalid payload version")

type payloadHandler struct {
	marshaller marshal.Marshalizer
	dp         map[uint32]DataProcessor
	actions    map[string]func(marshalledData []byte, version uint32) error
}

// NewPayloadHandler will create a new instance of events indexer
func NewPayloadHandler(marshaller marshal.Marshalizer, dataProcessors map[uint32]DataProcessor) (*payloadHandler, error) {
	if check.IfNil(marshaller) {
		return nil, common.ErrNilMarshaller
	}
	if len(dataProcessors) == 0 {
		return nil, ErrNilDataProcessor
	}

	payloadIndexer := &payloadHandler{
		marshaller: marshaller,
		dp:         dataProcessors,
	}
	payloadIndexer.initActionsMap()

	return payloadIndexer, nil
}

// GetOperationsMap returns the map with all the operations that will index data
func (ph *payloadHandler) initActionsMap() {
	ph.actions = map[string]func(d []byte, v uint32) error{
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
func (ph *payloadHandler) ProcessPayload(payload []byte, topic string, version uint32) error {
	payloadTypeAction, ok := ph.actions[topic]
	if !ok {
		log.Warn("invalid payload type", "topic", topic)
		return nil
	}

	return payloadTypeAction(payload, version)
}

func (ph *payloadHandler) saveBlock(marshalledData []byte, version uint32) error {
	dataProcessor, ok := ph.dp[version]
	if !ok {
		log.Warn("invalid provided version", "version", version)
		return ErrInvalidPayloadType
	}
	log.Debug("PayloadHandler", "version", version)

	return dataProcessor.SaveBlock(marshalledData)
}

func (ph *payloadHandler) revertIndexedBlock(marshalledData []byte, version uint32) error {
	dataProcessor, ok := ph.dp[version]
	if !ok {
		log.Warn("invalid provided version", "version", version)
		return ErrInvalidPayloadType
	}
	log.Debug("PayloadHandler", "version", version)

	return dataProcessor.RevertIndexedBlock(marshalledData)
}

func (ph *payloadHandler) finalizedBlock(marshalledData []byte, version uint32) error {
	dataProcessor, ok := ph.dp[version]
	if !ok {
		log.Warn("invalid provided version", "version", version)
		return ErrInvalidPayloadType
	}

	return dataProcessor.FinalizedBlock(marshalledData)
}

func (ph *payloadHandler) saveRounds(marshalledData []byte, version uint32) error {
	return nil
}

func (ph *payloadHandler) saveValidatorsRating(marshalledData []byte, version uint32) error {
	return nil
}

func (ph *payloadHandler) saveValidatorsPubKeys(marshalledData []byte, version uint32) error {
	return nil
}

func (ph *payloadHandler) saveAccounts(marshalledData []byte, version uint32) error {
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
