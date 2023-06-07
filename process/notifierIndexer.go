package process

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
)

var errNilDataIndexer = errors.New("nil data indexer")

type indexer struct {
	marshaller marshal.Marshalizer
	di         DataIndexer
	actions    map[string]func(marshalledData []byte) error
}

// NewNotifierIndexer will create a new instance of notifier indexer
func NewNotifierIndexer(marshaller marshal.Marshalizer, dataIndexer DataIndexer) (*indexer, error) {
	if check.IfNil(marshaller) {
		return nil, common.ErrNilMarshaller
	}
	if check.IfNil(dataIndexer) {
		return nil, errNilDataIndexer
	}

	payloadIndexer := &indexer{
		marshaller: marshaller,
		di:         dataIndexer,
	}
	payloadIndexer.initActionsMap()

	return payloadIndexer, nil
}

// GetOperationsMap returns the map with all the operations that will index data
func (i *indexer) initActionsMap() {
	i.actions = map[string]func(d []byte) error{
		outport.TopicSaveBlock:             i.saveBlock,
		outport.TopicRevertIndexedBlock:    i.revertIndexedBlock,
		outport.TopicSaveRoundsInfo:        i.saveRounds,
		outport.TopicSaveValidatorsRating:  i.saveValidatorsRating,
		outport.TopicSaveValidatorsPubKeys: i.saveValidatorsPubKeys,
		outport.TopicSaveAccounts:          i.saveAccounts,
		outport.TopicFinalizedBlock:        i.finalizedBlock,
	}
}

// ProcessPayload will proces the provided payload based on the topic
func (i *indexer) ProcessPayload(payload []byte, topic string) error {
	payloadTypeAction, ok := i.actions[topic]
	if !ok {
		log.Warn("invalid payload type", "topic", topic)
		return nil
	}

	return payloadTypeAction(payload)
}

func (i *indexer) saveBlock(marshalledData []byte) error {
	outportBlock := &outport.OutportBlock{}
	err := i.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return err
	}

	return i.di.SaveBlock(outportBlock)
}

func (i *indexer) revertIndexedBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := i.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	return i.di.RevertIndexedBlock(blockData)
}

func (i *indexer) finalizedBlock(marshalledData []byte) error {
	finalizedBlock := &outport.FinalizedBlock{}
	err := i.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	return i.di.FinalizedBlock(finalizedBlock)
}

func (i *indexer) saveRounds(marshalledData []byte) error {
	return nil
}

func (i *indexer) saveValidatorsRating(marshalledData []byte) error {
	return nil
}

func (i *indexer) saveValidatorsPubKeys(marshalledData []byte) error {
	return nil
}

func (i *indexer) saveAccounts(marshalledData []byte) error {
	return nil
}

// Close will close the indexer
func (i *indexer) Close() error {
	return i.di.Close()
}

// IsInterfaceNil returns true if underlying object is nil
func (i *indexer) IsInterfaceNil() bool {
	return i == nil
}
