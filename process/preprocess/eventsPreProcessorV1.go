package preprocess

import (
	"encoding/hex"
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

var (
	// ErrNilBlockData signals that a nil block data has been provided
	ErrNilBlockData = errors.New("nil block data")

	// ErrNilTransactionPool signals that a nil transaction pool has been provided
	ErrNilTransactionPool = errors.New("nil transaction pool")

	// ErrNilHeaderGasConsumption signals that a nil header gas consumption has been provided
	ErrNilHeaderGasConsumption = errors.New("nil header gas consumption")
)

type eventsPreProcessorV1 struct {
	*baseEventsPreProcessor
}

// NewEventsPreProcessorV1 will create a new events data preprocessor instance
func NewEventsPreProcessorV1(args ArgsEventsPreProcessor) (*eventsPreProcessorV1, error) {
	baseEventsPreProcessor, err := newBaseEventsPreProcessor(args)
	if err != nil {
		return nil, err
	}

	return &eventsPreProcessorV1{
		baseEventsPreProcessor: baseEventsPreProcessor,
	}, nil
}

// SaveBlock will handle the block info data
func (d *eventsPreProcessorV1) SaveBlock(marshalledData []byte) error {
	outportBlock := &outport.OutportBlock{}
	err := d.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return err
	}

	err = checkBlockDataValid(outportBlock)
	if err != nil {
		return err
	}

	header, err := d.getHeaderFromBytes(core.HeaderType(outportBlock.BlockData.HeaderType), outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return err
	}

	saveBlockData := &data.ArgsSaveBlockData{
		HeaderHash:             outportBlock.BlockData.HeaderHash,
		Body:                   outportBlock.BlockData.Body,
		SignersIndexes:         outportBlock.SignersIndexes,
		NotarizedHeadersHashes: outportBlock.NotarizedHeadersHashes,
		HeaderGasConsumption:   outportBlock.HeaderGasConsumption,
		AlteredAccounts:        outportBlock.AlteredAccounts,
		NumberOfShards:         outportBlock.NumberOfShards,
		TransactionsPool:       outportBlock.TransactionPool,
		Header:                 header,
	}

	err = d.facade.HandlePushEvents(*saveBlockData)
	if err != nil {
		return err
	}

	return nil
}

func checkBlockDataValid(block *outport.OutportBlock) error {
	if block.BlockData == nil {
		return ErrNilBlockData
	}
	if block.TransactionPool == nil {
		return ErrNilTransactionPool
	}
	if block.HeaderGasConsumption == nil {
		return ErrNilHeaderGasConsumption
	}

	return nil
}

// RevertIndexedBlock will handle the revert block event
func (d *eventsPreProcessorV1) RevertIndexedBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := d.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	header, err := d.getHeaderFromBytes(core.HeaderType(blockData.HeaderType), blockData.HeaderBytes)
	if err != nil {
		return err
	}

	revertData := &data.RevertBlock{
		Hash:      hex.EncodeToString(blockData.GetHeaderHash()),
		Nonce:     header.GetNonce(),
		Round:     header.GetRound(),
		Epoch:     header.GetEpoch(),
		ShardID:   blockData.GetShardID(),
		TimeStamp: header.GetTimeStamp(),
	}

	d.facade.HandleRevertEvents(*revertData)

	return nil
}

// FinalizedBlock will handle the finalized block event
func (d *eventsPreProcessorV1) FinalizedBlock(marshalledData []byte) error {
	finalizedBlock := &outport.FinalizedBlock{}
	err := d.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	finalizedData := data.FinalizedBlock{
		Hash: hex.EncodeToString(finalizedBlock.GetHeaderHash()),
	}

	d.facade.HandleFinalizedEvents(finalizedData)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *eventsPreProcessorV1) IsInterfaceNil() bool {
	return d == nil
}
