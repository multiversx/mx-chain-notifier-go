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

var errNilBlockData = errors.New("nil block data")
var errNilTransactionPool = errors.New("nil transaction pool")
var errNilHeaderGasConsumption = errors.New("nil header gas consumption")

type ArgsDataIndexer struct {
	Marshaller         marshal.Marshalizer
	InternalMarshaller marshal.Marshalizer
	Facade             EventsFacadeHandler
}

type dataIndexer struct {
	marshaller         marshal.Marshalizer
	internalMarshaller marshal.Marshalizer
	emptyBlockCreator  EmptyBlockCreatorContainer
	facade             EventsFacadeHandler
}

func NewDataIndexerHandler(args ArgsDataIndexer) (*dataIndexer, error) {
	err := checkDataIndexerArgs(args)
	if err != nil {
		return nil, err
	}

	di := &dataIndexer{
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

func checkDataIndexerArgs(args ArgsDataIndexer) error {
	if check.IfNil(args.Marshaller) {
		return common.ErrNilMarshaller
	}
	if check.IfNil(args.InternalMarshaller) {
		return fmt.Errorf("%w (internal)", common.ErrNilMarshaller)
	}

	return nil
}

func (d *dataIndexer) SaveBlock(outportBlock *outport.OutportBlock) error {
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

func (d *dataIndexer) getHeaderFromBytes(headerType core.HeaderType, headerBytes []byte) (header coreData.HeaderHandler, err error) {
	creator, err := d.emptyBlockCreator.Get(headerType)
	if err != nil {
		return nil, err
	}

	return block.GetHeaderFromBytes(d.internalMarshaller, creator, headerBytes)
}

func (d *dataIndexer) getEmptyHeaderCreator(headerType string) (block.EmptyBlockCreator, error) {
	return d.emptyBlockCreator.Get(core.HeaderType(headerType))
}

func checkBlockDataValid(block *outport.OutportBlock) error {
	if block.BlockData == nil {
		return errNilBlockData
	}
	if block.TransactionPool == nil {
		return errNilTransactionPool
	}
	if block.HeaderGasConsumption == nil {
		return errNilHeaderGasConsumption
	}

	return nil
}

func (d *dataIndexer) RevertIndexedBlock(blockData *outport.BlockData) error {
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

func (d *dataIndexer) FinalizedBlock(finalizedBlock *outport.FinalizedBlock) error {
	finalizedData := data.FinalizedBlock{
		Hash: string(finalizedBlock.GetHeaderHash()),
	}

	d.facade.HandleFinalizedEvents(finalizedData)

	return nil
}

func (d *dataIndexer) Close() error {
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

func (d *dataIndexer) IsInterfaceNil() bool {
	return d == nil
}
