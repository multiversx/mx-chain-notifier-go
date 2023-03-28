package groups

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// UnmarshallBlockDataV1 will try to unmarshal block data with old format
func UnmarshallBlockDataV1(marshalledData []byte) (*data.SaveBlockData, error) {
	var saveBlockData data.SaveBlockData

	err := json.Unmarshal(marshalledData, &saveBlockData)
	if err != nil {
		return nil, err
	}

	return &saveBlockData, nil
}

// UnmarshallBlockData will try to unmarshal block data
func UnmarshallBlockData(marshaller marshal.Marshalizer, marshalledData []byte) (*data.ArgsSaveBlockData, error) {
	var argsBlockS *outport.OutportBlock
	err := json.Unmarshal(marshalledData, &argsBlockS)
	if err != nil {
		return nil, err
	}

	headerCreator, err := getHeaderCreator(argsBlockS.BlockData.HeaderType)
	if err != nil {
		return nil, err
	}

	header, err := block.GetHeaderFromBytes(marshaller, headerCreator, argsBlockS.BlockData.HeaderBytes)
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

func getHeaderCreator(headerType string) (block.EmptyBlockCreator, error) {
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

	return container.Get(core.HeaderType(headerType))
}
