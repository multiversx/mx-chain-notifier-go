package preprocess

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/core"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// TODO: dismiss this implementation after http integration is fully deprecated
type eventsPreProcessorV0 struct {
	*baseEventsPreProcessor
}

// NewEventsPreProcessorV0 will create a new events data preprocessor instance
func NewEventsPreProcessorV0(args ArgsEventsPreProcessor) (*eventsPreProcessorV0, error) {
	baseEventsPreProcessor, err := newBaseEventsPreProcessor(args)
	if err != nil {
		return nil, err
	}

	return &eventsPreProcessorV0{
		baseEventsPreProcessor: baseEventsPreProcessor,
	}, nil
}

// SaveBlock will handle the block info data
func (d *eventsPreProcessorV0) SaveBlock(marshalledData []byte) error {
	argsBlockS := struct {
		HeaderHash             []byte
		Body                   *block.Body
		TransactionsPool       *outport.TransactionPool
		SignersIndexes         []uint64
		NotarizedHeadersHashes []string
		HeaderGasConsumption   outport.HeaderGasConsumption
		AlteredAccounts        map[string]*alteredAccount.AlteredAccount
		NumberOfShards         uint32
		IsImportDB             bool
	}{}
	err := json.Unmarshal(marshalledData, &argsBlockS)
	if err != nil {
		return err
	}

	header, err := d.getHeader(marshalledData)
	if err != nil {
		return err
	}

	saveBlockData := &data.ArgsSaveBlockData{
		HeaderHash:             argsBlockS.HeaderHash,
		Body:                   argsBlockS.Body,
		SignersIndexes:         argsBlockS.SignersIndexes,
		NotarizedHeadersHashes: argsBlockS.NotarizedHeadersHashes,
		HeaderGasConsumption:   &argsBlockS.HeaderGasConsumption,
		AlteredAccounts:        argsBlockS.AlteredAccounts,
		NumberOfShards:         argsBlockS.NumberOfShards,
		TransactionsPool:       argsBlockS.TransactionsPool,
		Header:                 header,
	}

	err = d.facade.HandlePushEventsV2(*saveBlockData)
	if err != nil {
		return err
	}

	return nil
}

func (d *eventsPreProcessorV0) getHeader(marshaledData []byte) (nodeData.HeaderHandler, error) {
	headerTypeStruct := struct {
		HeaderType core.HeaderType
	}{}

	err := json.Unmarshal(marshaledData, &headerTypeStruct)
	if err != nil {
		return nil, err
	}

	header, err := d.getHeaderFromBytes(headerTypeStruct.HeaderType, marshaledData)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (d *eventsPreProcessorV0) getHeaderFromBytes(headerType core.HeaderType, headerBytes []byte) (header coreData.HeaderHandler, err error) {
	creator, err := d.emptyBlockCreator.Get(headerType)
	if err != nil {
		return nil, err
	}

	return block.GetHeaderFromBytes(d.marshaller, creator, headerBytes)
}

// RevertIndexedBlock will handle the revert block event
func (d *eventsPreProcessorV0) RevertIndexedBlock(marshalledData []byte) error {
	revertBlock := &data.RevertBlock{}
	err := d.marshaller.Unmarshal(revertBlock, marshalledData)
	if err != nil {
		return err
	}

	d.facade.HandleRevertEvents(*revertBlock)

	return nil
}

// FinalizedBlock will handler the finalized block event
func (d *eventsPreProcessorV0) FinalizedBlock(marshalledData []byte) error {
	finalizedBlock := &data.FinalizedBlock{}
	err := d.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	d.facade.HandleFinalizedEvents(*finalizedBlock)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *eventsPreProcessorV0) IsInterfaceNil() bool {
	return d == nil
}
