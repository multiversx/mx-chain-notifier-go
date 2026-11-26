package preprocess

import (
	"encoding/hex"
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	coreData "github.com/multiversx/mx-chain-core-go/data"
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

	if outportBlock.BlockData == nil {
		return ErrNilBlockData
	}

	headerType := core.HeaderType(outportBlock.BlockData.HeaderType)

	header, err := d.getHeaderFromBytes(headerType, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return err
	}

	err = checkHeaderGasConsumption(header, outportBlock)
	if err != nil {
		return err
	}

	var executionResults map[string]*outport.ExecutionResultData
	if header.IsHeaderV3() {
		executionResults = outportBlock.BlockData.Results
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
		HeaderTimeStampMs:      outportBlock.BlockData.GetTimestampMs(),
		StateAccesses:          outportBlock.GetStateAccesses(),
		StateAccessesForBlock:  outportBlock.GetStateAccessesForBlock(),
		Results:                executionResults,
	}

	err = d.facade.HandlePushEvents(*saveBlockData)
	if err != nil {
		return err
	}

	return nil
}

func checkHeaderGasConsumption(header coreData.HeaderHandler, block *outport.OutportBlock) error {
	if !header.IsHeaderV3() {
		if block.HeaderGasConsumption == nil {
			return ErrNilHeaderGasConsumption
		}

		return nil
	}

	for _, execRes := range block.BlockData.Results {
		if execRes == nil {
			continue
		}
		if execRes.HeaderGasConsumption == nil {
			return ErrNilHeaderGasConsumption
		}
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

	headerTimeStamp := header.GetTimeStamp()
	headerTimeStampMs := blockData.GetTimestampMs()

	revertData := &data.RevertBlock{
		Hash:        hex.EncodeToString(blockData.GetHeaderHash()),
		Nonce:       header.GetNonce(),
		Round:       header.GetRound(),
		Epoch:       header.GetEpoch(),
		ShardID:     blockData.GetShardID(),
		TimeStamp:   headerTimeStamp,
		TimeStampMs: headerTimeStampMs,
	}

	return d.facade.HandleRevertEvents(*revertData)
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

	return d.facade.HandleFinalizedEvents(finalizedData)
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *eventsPreProcessorV1) IsInterfaceNil() bool {
	return d == nil
}
