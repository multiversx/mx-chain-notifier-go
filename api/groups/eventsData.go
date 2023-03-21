package groups

import (
	"encoding/json"
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/unmarshal"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
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

// UnmarshallBlockDataV2 will try to unmarshal block data v2
func UnmarshallBlockDataV2(marshalledData []byte) (*data.ArgsSaveBlockData, error) {
	argsBlockS := struct {
		BlockData              *outport.BlockData
		TransactionsPool       *data.TransactionsPool
		SignersIndexes         []uint64
		NotarizedHeadersHashes []string
		HeaderGasConsumption   *outport.HeaderGasConsumption
		AlteredAccounts        map[string]*alteredAccount.AlteredAccount
		NumberOfShards         uint32
		IsImportDB             bool
	}{}
	err := json.Unmarshal(marshalledData, &argsBlockS)
	if err != nil {
		return nil, err
	}

	header, err := unmarshal.GetHeaderFromBytes(&marshal.GogoProtoMarshalizer{}, core.HeaderType(argsBlockS.BlockData.HeaderType), argsBlockS.BlockData.HeaderBytes)
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
		TransactionsPool:       argsBlockS.TransactionsPool,
		Header:                 header,
	}, nil
}

func getHeader(marshaledData []byte) (nodeData.HeaderHandler, error) {
	headerTypeStruct := struct {
		HeaderType core.HeaderType
	}{}

	err := json.Unmarshal(marshaledData, &headerTypeStruct)
	if err != nil {
		return nil, err
	}

	switch headerTypeStruct.HeaderType {
	case core.MetaHeader:
		hStruct := struct {
			H1 *block.MetaBlock `json:"Header"`
		}{}
		err = json.Unmarshal(marshaledData, &hStruct)
		return hStruct.H1, err
	case core.ShardHeaderV1:
		hStruct := struct {
			H1 *block.Header `json:"Header"`
		}{}
		err = json.Unmarshal(marshaledData, &hStruct)
		return hStruct.H1, err
	case core.ShardHeaderV2:
		hStruct := struct {
			H1 *block.HeaderV2 `json:"Header"`
		}{}
		err = json.Unmarshal(marshaledData, &hStruct)
		return hStruct.H1, err
	default:
		return nil, errors.New("invalid header type")
	}
}
