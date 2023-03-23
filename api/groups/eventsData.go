package groups

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/unmarshal"
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

// UnmarshallBlockData will try to unmarshal block data v2
func UnmarshallBlockData(marshalledData []byte) (*data.ArgsSaveBlockData, error) {
	var argsBlockS *outport.OutportBlock
	err := json.Unmarshal(marshalledData, &argsBlockS)
	if err != nil {
		return nil, err
	}

	header, err := unmarshal.GetHeaderFromBytes(&marshal.JsonMarshalizer{}, core.HeaderType(argsBlockS.BlockData.HeaderType), argsBlockS.BlockData.HeaderBytes)
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
