package process

import (
	"errors"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// ErrNilBlockData signals that a nil block data has been provided
var ErrNilBlockData = errors.New("nil block data")

// ErrNilTransactionPool signals that a nil transaction pool has been provided
var ErrNilTransactionPool = errors.New("nil transaction pool")

// ErrNilHeaderGasConsumption signals that a nil header gas consumption has been provided
var ErrNilHeaderGasConsumption = errors.New("nil header gas consumption")

// ArgsEventsDataPreProcessor defines the arguments needed to create a new events data preprocessor
type ArgsEventsDataPreProcessor struct {
	Marshaller         marshal.Marshalizer
	InternalMarshaller marshal.Marshalizer
	Facade             EventsFacadeHandler
}

type eventsDataPreProcessor struct {
	marshaller         marshal.Marshalizer
	internalMarshaller marshal.Marshalizer
	emptyBlockCreator  EmptyBlockCreatorContainer
	facade             EventsFacadeHandler
}

// NewEventsDataPreProcessor will create a new events data preprocessor instance
func NewEventsDataPreProcessor(args ArgsEventsDataPreProcessor) (*eventsDataPreProcessor, error) {
	err := checkDataIndexerArgs(args)
	if err != nil {
		return nil, err
	}

	di := &eventsDataPreProcessor{
		marshaller:         args.Marshaller,
		internalMarshaller: args.InternalMarshaller,
		facade:             args.Facade,
	}

	emptyBlockContainer, err := createEmptyBlockCreatorContainer()
	if err != nil {
		return nil, err
	}

	di.emptyBlockCreator = emptyBlockContainer

	return di, nil
}

func checkDataIndexerArgs(args ArgsEventsDataPreProcessor) error {
	if check.IfNil(args.Marshaller) {
		return common.ErrNilMarshaller
	}
	if check.IfNil(args.InternalMarshaller) {
		return fmt.Errorf("%w (internal)", common.ErrNilMarshaller)
	}
	if check.IfNil(args.Facade) {
		return common.ErrNilFacadeHandler
	}

	return nil
}

// SaveBlock will handle the block info data
func (d *eventsDataPreProcessor) SaveBlock(outportBlock *outport.OutportBlock) error {
	err := checkBlockDataValid(outportBlock)
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

	err = d.facade.HandlePushEventsV2(*saveBlockData)
	if err != nil {
		return err
	}

	return nil
}

func (d *eventsDataPreProcessor) getHeaderFromBytes(headerType core.HeaderType, headerBytes []byte) (header coreData.HeaderHandler, err error) {
	creator, err := d.emptyBlockCreator.Get(headerType)
	if err != nil {
		return nil, err
	}

	return block.GetHeaderFromBytes(d.internalMarshaller, creator, headerBytes)
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
func (d *eventsDataPreProcessor) RevertIndexedBlock(blockData *outport.BlockData) error {
	header, err := d.getHeaderFromBytes(core.HeaderType(blockData.HeaderType), blockData.HeaderBytes)
	if err != nil {
		return err
	}

	revertData := &data.RevertBlock{
		Hash:  string(header.GetRootHash()),
		Nonce: header.GetNonce(),
		Round: header.GetRound(),
		Epoch: header.GetEpoch(),
	}

	d.facade.HandleRevertEvents(*revertData)

	return nil
}

// FinalizedBlock will handler the finalized block event
func (d *eventsDataPreProcessor) FinalizedBlock(finalizedBlock *outport.FinalizedBlock) error {
	finalizedData := data.FinalizedBlock{
		Hash: string(finalizedBlock.GetHeaderHash()),
	}

	d.facade.HandleFinalizedEvents(finalizedData)

	return nil
}

func createEmptyBlockCreatorContainer() (EmptyBlockCreatorContainer, error) {
	container := block.NewEmptyBlockCreatorsContainer()

	err := container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())
	if err != nil {
		return nil, err
	}

	err = container.Add(core.ShardHeaderV2, block.NewEmptyHeaderV2Creator())
	if err != nil {
		return nil, err
	}

	err = container.Add(core.MetaHeader, block.NewEmptyMetaBlockCreator())
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (d *eventsDataPreProcessor) IsInterfaceNil() bool {
	return d == nil
}
