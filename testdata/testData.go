package testdata

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/stateChange"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	notifierData "github.com/multiversx/mx-chain-notifier-go/data"
)

type blockData struct {
	marshaller marshal.Marshalizer
}

// NewBlockData will create block data component for testing
func NewBlockData(marshaller marshal.Marshalizer) (*blockData, error) {
	if check.IfNil(marshaller) {
		return nil, common.ErrNilMarshaller
	}

	return &blockData{marshaller: marshaller}, nil
}

// OldSaveBlockData defines block events data before initial refactoring
func (bd *blockData) OldSaveBlockData() *notifierData.SaveBlockData {
	return &notifierData.SaveBlockData{
		Hash: "blockHash",
		Txs: map[string]*transaction.Transaction{
			"hash1": {
				Nonce: 1,
			},
		},
		Scrs: map[string]*smartContractResult.SmartContractResult{
			"hash2": {
				Nonce: 2,
			},
		},
		LogEvents: []notifierData.Event{
			{
				Address: "logaddr1",
			},
		},
	}
}

// OutportBlockV0 -
func (bd *blockData) OutportBlockV0() *notifierData.ArgsSaveBlock {
	stateAccesses := make(map[string]*stateChange.StateAccesses)
	stateAccesses["txHash1"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:           stateChange.Write,
				MainTrieKey:    []byte("mainTrieKey1"),
				MainTrieVal:    []byte("mainTrieVal1"),
				TxHash:         []byte("txHash1"),
				AccountChanges: 8,
			},
			&stateChange.StateAccess{
				Type:           stateChange.Write,
				MainTrieKey:    []byte("mainTrieKey2"),
				MainTrieVal:    []byte("mainTrieVal2"),
				TxHash:         []byte("txHash1"),
				AccountChanges: 4,
			},
		},
	}
	stateAccesses["txHash2"] = &stateChange.StateAccesses{}

	saveBlockData := notifierData.OutportBlockDataOld{
		HeaderHash: []byte("headerHash3"),
		Body: &block.Body{
			MiniBlocks: []*block.MiniBlock{
				{
					TxHashes:        [][]byte{},
					ReceiverShardID: 1,
					SenderShardID:   1,
				},
			},
		},
		TransactionsPool: &notifierData.TransactionsPool{
			Txs: map[string]*notifierData.NodeTransaction{
				"txHash1": {
					TransactionHandler: &transaction.Transaction{
						Nonce:    1,
						GasPrice: 1,
						GasLimit: 1,
					},
					FeeInfo: outport.FeeInfo{
						GasUsed: 1,
					},
					ExecutionOrder: 2,
				},
			},
			Scrs: map[string]*notifierData.NodeSmartContractResult{
				"scrHash1": {
					TransactionHandler: &smartContractResult.SmartContractResult{
						Nonce:    2,
						GasLimit: 2,
						GasPrice: 2,
						CallType: 2,
					},
					FeeInfo: outport.FeeInfo{
						GasUsed: 2,
					},
					ExecutionOrder: 0,
				},
			},
			Logs: []*notifierData.LogData{
				{
					LogHandler: &transaction.Log{
						Address: []byte("logaddr1"),
						Events:  []*transaction.Event{},
					},
					TxHash: "logHash1",
				},
			},
		},
		NumberOfShards: 2,
		StateAccesses:  stateAccesses,
	}

	return &notifierData.ArgsSaveBlock{
		HeaderType:          "Header",
		OutportBlockDataOld: saveBlockData,
	}
}

// OutportBlockV1 -
func (bd *blockData) OutportBlockV1() *outport.OutportBlock {
	header := &block.Header{
		ShardID:   1,
		TimeStamp: 1234,
	}
	headerBytes, _ := bd.marshaller.Marshal(header)
	headerHash := []byte("headerHash1")
	stateAccessesForBlock := getStateAccessesForBlock(headerHash)

	return &outport.OutportBlock{
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  "Header",
			HeaderHash:  headerHash,
			Body: &block.Body{
				MiniBlocks: []*block.MiniBlock{
					{
						TxHashes:        [][]byte{},
						ReceiverShardID: 1,
						SenderShardID:   1,
					},
				},
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
		TransactionPool: &outport.TransactionPool{
			Transactions: map[string]*outport.TxInfo{
				hex.EncodeToString([]byte("txHash1")): {
					Transaction: &transaction.Transaction{
						Nonce:    1,
						GasPrice: 1,
						GasLimit: 1,
					},
					FeeInfo: &outport.FeeInfo{
						GasUsed: 1,
					},
					ExecutionOrder: 2,
				},
			},
			SmartContractResults: map[string]*outport.SCRInfo{
				hex.EncodeToString([]byte("scrHash1")): {
					SmartContractResult: &smartContractResult.SmartContractResult{
						Nonce:    2,
						GasLimit: 2,
						GasPrice: 2,
						CallType: 2,
					},
					FeeInfo: &outport.FeeInfo{
						GasUsed: 2,
					},
					ExecutionOrder: 0,
				},
			},
			Logs: []*transaction.LogData{
				{
					Log: &transaction.Log{
						Address: []byte("logaddr1"),
						Events:  []*transaction.Event{},
					},
					TxHash: "logHash1",
				},
			},
		},
		StateAccessesForBlock: stateAccessesForBlock,
		NumberOfShards:        2,
	}
}

func getStateAccessesForBlock(headerHash []byte) map[string]*outport.StateAccessesForBlock {
	stateAccesses := make(map[string]*stateChange.StateAccesses)
	stateAccesses["txHash1"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			&stateChange.StateAccess{
				Type:           stateChange.Write,
				MainTrieKey:    []byte("mainTrieKey1"),
				MainTrieVal:    []byte("mainTrieVal1"),
				TxHash:         []byte("txHash1"),
				AccountChanges: 8,
			},
			&stateChange.StateAccess{
				Type:           stateChange.Write,
				MainTrieKey:    []byte("mainTrieKey2"),
				MainTrieVal:    []byte("mainTrieVal2"),
				TxHash:         []byte("txHash1"),
				AccountChanges: 4,
			},
		},
	}
	stateAccesses["txHash2"] = &stateChange.StateAccesses{}
	stateAccessesForBlock := map[string]*outport.StateAccessesForBlock{}
	stateAccessesForBlock[hex.EncodeToString(headerHash)] = &outport.StateAccessesForBlock{StateAccesses: stateAccesses}
	return stateAccessesForBlock
}

// OutportBlockV2 -
func (bd *blockData) OutportBlockV2() *outport.OutportBlock {
	header := &block.HeaderV3{
		ShardID:     1,
		TimestampMs: 1234,
	}
	headerBytes, _ := bd.marshaller.Marshal(header)

	execBlockHash := []byte("execBlockHash1")
	stateAccessesForBlock := getStateAccessesForBlock(execBlockHash)

	blockBody := &block.Body{
		MiniBlocks: []*block.MiniBlock{
			{
				TxHashes:        [][]byte{},
				ReceiverShardID: 1,
				SenderShardID:   1,
			},
		},
	}

	execResTxPool := &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			hex.EncodeToString([]byte("txHash1")): {
				Transaction: &transaction.Transaction{
					Nonce:    1,
					GasPrice: 1,
					GasLimit: 1,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 1,
				},
				ExecutionOrder: 2,
			},
		},
		SmartContractResults: map[string]*outport.SCRInfo{
			hex.EncodeToString([]byte("scrHash1")): {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce:    2,
					GasLimit: 2,
					GasPrice: 2,
					CallType: 2,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 2,
				},
				ExecutionOrder: 0,
			},
		},
		Logs: []*transaction.LogData{
			{
				Log: &transaction.Log{
					Address: []byte("logaddr1"),
					Events: []*transaction.Event{
						{
							Address: []byte("logaddr1"),
						},
					},
				},
				TxHash: "txHash1",
			},
		},
	}

	execBlockHash2 := []byte("execBlockHash2")

	execResults := map[string]*outport.ExecutionResultData{
		hex.EncodeToString(execBlockHash): {
			Body:                 blockBody,
			TransactionPool:      execResTxPool,
			HeaderGasConsumption: &outport.HeaderGasConsumption{},
		},
		hex.EncodeToString(execBlockHash2): {
			Body:                 blockBody,
			TransactionPool:      execResTxPool,
			HeaderGasConsumption: &outport.HeaderGasConsumption{},
		},
	}

	return &outport.OutportBlock{
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  "HeaderV3",
			HeaderHash:  []byte("headerHash1"),
			Body: &block.Body{
				MiniBlocks: []*block.MiniBlock{
					{
						TxHashes:        [][]byte{},
						ReceiverShardID: 1,
						SenderShardID:   1,
					},
				},
			},
			Results: execResults,
		},
		HeaderGasConsumption:  &outport.HeaderGasConsumption{},
		NumberOfShards:        2,
		StateAccessesForBlock: stateAccessesForBlock,
	}
}

// RevertBlockV0 -
func (bd *blockData) RevertBlockV0() *notifierData.RevertBlock {
	return &notifierData.RevertBlock{
		Hash:  "headerHash1",
		Nonce: 1,
		Round: 1,
		Epoch: 1,
	}
}

// RevertBlockV1 -
func (bd *blockData) RevertBlockV1() *outport.BlockData {
	header := &block.Header{
		ShardID:   1,
		TimeStamp: 1234,
	}
	headerBytes, _ := bd.marshaller.Marshal(header)

	return &outport.BlockData{
		ShardID:     1,
		HeaderBytes: headerBytes,
		HeaderType:  "Header",
		HeaderHash:  []byte("headerHash1"),
		Body: &block.Body{
			MiniBlocks: []*block.MiniBlock{
				{
					TxHashes:        [][]byte{},
					ReceiverShardID: 1,
					SenderShardID:   1,
				},
			},
		},
	}
}

// FinalizedBlockV0 -
func (bd *blockData) FinalizedBlockV0() *notifierData.FinalizedBlock {
	return &notifierData.FinalizedBlock{
		Hash: "headerHash1",
	}
}

// FinalizedBlockV1 -
func (bd *blockData) FinalizedBlockV1() *outport.FinalizedBlock {
	return &outport.FinalizedBlock{
		ShardID:    1,
		HeaderHash: []byte("headerHash1"),
	}
}
