package tools

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	notifierData "github.com/multiversx/mx-chain-notifier-go/data"
)

type blockData struct {
	marshaller marshal.Marshalizer
}

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
	return &notifierData.ArgsSaveBlock{
		HeaderType: "Header",
		ArgsSaveBlockData: notifierData.ArgsSaveBlockData{
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
			TransactionsPool: &outport.TransactionPool{
				Transactions: map[string]*outport.TxInfo{
					"txHash1": {
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
					"scrHash1": {
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
				Logs: []*outport.LogData{
					{
						Log: &transaction.Log{
							Address: []byte("logaddr1"),
							Events:  []*transaction.Event{},
						},
						TxHash: "logHash1",
					},
				},
			},
			NumberOfShards: 2,
		},
	}
}

// OutportBlockV1 -
func (bd *blockData) OutportBlockV1() *outport.OutportBlock {
	header := &block.Header{
		ShardID:   1,
		TimeStamp: 1234,
	}
	headerBytes, _ := bd.marshaller.Marshal(header)

	return &outport.OutportBlock{
		BlockData: &outport.BlockData{
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
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
		TransactionPool: &outport.TransactionPool{
			Transactions: map[string]*outport.TxInfo{
				"txHash1": {
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
				"scrHash1": {
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
			Logs: []*outport.LogData{
				{
					Log: &transaction.Log{
						Address: []byte("logaddr1"),
						Events:  []*transaction.Event{},
					},
					TxHash: "logHash1",
				},
			},
		},
		NumberOfShards: 2,
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
