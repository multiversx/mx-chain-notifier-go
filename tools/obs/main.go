package main

import (
	"encoding/json"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	notifierData "github.com/multiversx/mx-chain-notifier-go/data"
)

func main() {
	args := HTTPClientWrapperArgs{
		UseAuthorization:  false,
		BaseUrl:           "http://localhost:5000",
		RequestTimeoutSec: 10,
	}
	httpClient, err := NewHTTPWrapperClient(args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := getShardV2Data()

	err = httpClient.Post("/events/push", data)
	if err != nil {
		fmt.Println(fmt.Errorf("%w in eventNotifier.SaveBlock while posting block data", err))
		return
	}

	fmt.Println(data)
}

func getShardV1Data() *notifierData.SaveBlockData {
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

func getShardV2Data() *outport.OutportBlock {
	header := &block.Header{
		ShardID:   1,
		TimeStamp: 1234,
	}
	headerBytes, _ := json.Marshal(header)

	return &outport.OutportBlock{
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  "Header",
			HeaderHash:  []byte("headerHash"),
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
