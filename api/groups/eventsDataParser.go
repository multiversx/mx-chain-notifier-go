package groups

import (
	"encoding/json"

	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

func GetSaveBlockData(marshalledData []byte) (*data.SaveBlockData, error) {
	var saveBlockData *data.SaveBlockData

	err := json.Unmarshal(marshalledData, saveBlockData)
	if err != nil {
		return nil, err
	}

	return saveBlockData, nil
}

func GetArgsSaveBlockData(marshalledData []byte) (*data.ArgsSaveBlockData, error) {
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

	body, err := GetBody(marshalledData)
	if err != nil {
		return nil, err
	}

	header, err := GetHeader(marshalledData)
	if err != nil {
		return nil, err
	}

	txsPool, err := GetTransactionsPool(marshalledData)
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

func GetTransactionsPool(marshaledData []byte) (*data.TransactionsPool, error) {
	txPoolStruct := struct {
		TransactionsPool *data.TransactionsPool
	}{}

	err := json.Unmarshal(marshaledData, &txPoolStruct)
	return txPoolStruct.TransactionsPool, err
}

func GetBody(marshaledData []byte) (nodeData.BodyHandler, error) {
	bodyStruct := struct {
		Body *block.Body
	}{}

	err := json.Unmarshal(marshaledData, &bodyStruct)
	return bodyStruct.Body, err
}

func GetHeader(marshaledData []byte) (nodeData.HeaderHandler, error) {
	h2Struct := struct {
		H *block.HeaderV2 `json:"Header"`
	}{}
	err := json.Unmarshal(marshaledData, &h2Struct)
	if err == nil {
		if isHeaderV2(h2Struct.H) {
			return h2Struct.H, nil
		}
	}

	h1Struct := struct {
		H *block.Header `json:"Header"`
	}{}
	err = json.Unmarshal(marshaledData, &h1Struct)
	if err == nil {
		if !isMetaBlock(h1Struct.H) {
			return h1Struct.H, err
		}
	}

	mStruct := struct {
		H *block.MetaBlock `json:"Header"`
	}{}
	err = json.Unmarshal(marshaledData, &mStruct)
	if err == nil {
		return mStruct.H, err
	}

	return nil, err
}

func isHeaderV2(header *block.HeaderV2) bool {
	if header == nil {
		return false
	}
	if header.Header == nil {
		return false
	}

	return true
}

func isMetaBlock(header *block.Header) bool {
	if header.EpochStartMetaHash == nil {
		return true
	}

	return false
}
