package groups

import (
	"encoding/json"
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

func GetBlockDataV1(marshalledData []byte) (*data.SaveBlockData, error) {
	var saveBlockData data.SaveBlockData

	err := json.Unmarshal(marshalledData, &saveBlockData)
	if err != nil {
		return nil, err
	}

	return &saveBlockData, nil
}

func GetBlockDataV2(marshalledData []byte) (*data.ArgsSaveBlockData, error) {
	argsBlockS := struct {
		HeaderHash             []byte
		SignersIndexes         []uint64
		NotarizedHeadersHashes []string
		HeaderGasConsumption   outport.HeaderGasConsumption
		AlteredAccounts        map[string]*outport.AlteredAccount
		NumberOfShards         uint32
		IsImportDB             bool
	}{}
	err := json.Unmarshal(marshalledData, &argsBlockS)
	if err != nil {
		return nil, err
	}

	body, err := getBody(marshalledData)
	if err != nil {
		return nil, err
	}

	header, err := getHeader(marshalledData)
	if err != nil {
		return nil, err
	}

	txsPool, err := getTransactionsPool(marshalledData)
	if err != nil {
		return nil, err
	}

	return &data.ArgsSaveBlockData{
		HeaderHash:             argsBlockS.HeaderHash,
		Body:                   body,
		SignersIndexes:         argsBlockS.SignersIndexes,
		NotarizedHeadersHashes: argsBlockS.NotarizedHeadersHashes,
		HeaderGasConsumption:   argsBlockS.HeaderGasConsumption,
		AlteredAccounts:        argsBlockS.AlteredAccounts,
		NumberOfShards:         argsBlockS.NumberOfShards,
		IsImportDB:             argsBlockS.IsImportDB,
		Header:                 header,
		TransactionsPool:       txsPool,
	}, nil
}

func getTransactionsPool(marshaledData []byte) (*data.TransactionsPool, error) {
	txPoolStruct := struct {
		TransactionsPool *data.TransactionsPool
	}{}

	err := json.Unmarshal(marshaledData, &txPoolStruct)
	return txPoolStruct.TransactionsPool, err
}

func getBody(marshaledData []byte) (nodeData.BodyHandler, error) {
	bodyStruct := struct {
		Body *block.Body
	}{}

	err := json.Unmarshal(marshaledData, &bodyStruct)
	return bodyStruct.Body, err
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
