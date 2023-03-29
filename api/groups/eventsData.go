package groups

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

type eventsDataHandler struct {
	marshaller        marshal.Marshalizer
	emptyBlockCreator EmptyBlockCreatorContainer
}

// NewEventsDataHandler will create a new data handler component
func NewEventsDataHandler(marshaller marshal.Marshalizer) (*eventsDataHandler, error) {
	if check.IfNil(marshaller) {
		return nil, common.ErrNilMarshaller
	}

	edh := &eventsDataHandler{
		marshaller: marshaller,
	}

	emptyBlockContainer, err := createEmptyBlockCreatorContainer()
	if err != nil {
		return nil, err
	}

	edh.emptyBlockCreator = emptyBlockContainer

	return edh, nil
}

// UnmarshallBlockDataV1 will try to unmarshal block data with old format
func (edh *eventsDataHandler) UnmarshallBlockDataV1(marshalledData []byte) (*data.SaveBlockData, error) {
	var saveBlockData data.SaveBlockData

	err := json.Unmarshal(marshalledData, &saveBlockData)
	if err != nil {
		return nil, err
	}

	return &saveBlockData, nil
}

// UnmarshallBlockData will try to unmarshal block data
func (edh *eventsDataHandler) UnmarshallBlockData(marshalledData []byte) (*data.ArgsSaveBlockData, error) {
	var argsBlockS *outport.OutportBlock
	err := json.Unmarshal(marshalledData, &argsBlockS)
	if err != nil {
		return nil, err
	}

	err = checkBlockDataValid(argsBlockS)
	if err != nil {
		return nil, err
	}

	headerCreator, err := edh.getEmptyHeader(argsBlockS.BlockData.HeaderType)
	if err != nil {
		return nil, err
	}

	header, err := block.GetHeaderFromBytes(edh.marshaller, headerCreator, argsBlockS.BlockData.HeaderBytes)
	if err != nil {
		return nil, err
	}

	return &data.ArgsSaveBlockData{
		HeaderHash:             argsBlockS.BlockData.HeaderHash,
		Body:                   argsBlockS.BlockData.Body,
		SignersIndexes:         argsBlockS.SignersIndexes,
		NotarizedHeadersHashes: argsBlockS.NotarizedHeadersHashes,
		HeaderGasConsumption:   argsBlockS.HeaderGasConsumption,
		AlteredAccounts:        argsBlockS.AlteredAccounts,
		NumberOfShards:         argsBlockS.NumberOfShards,
		IsImportDB:             argsBlockS.IsImportDB,
		TransactionsPool:       argsBlockS.TransactionPool,
		Header:                 header,
	}, nil
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

func (edh *eventsDataHandler) getEmptyHeader(headerType string) (block.EmptyBlockCreator, error) {
	return edh.emptyBlockCreator.Get(core.HeaderType(headerType))
}

// IsInterfaceNil returns true if there is no value under the interface
func (edh *eventsDataHandler) IsInterfaceNil() bool {
	return edh == nil
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
