package main

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data"
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

func getShardV1Data() *outport.ArgsSaveBlockData {
	return &outport.ArgsSaveBlockData{
		HeaderHash: []byte("hash3"),
		Body: &block.Body{
			MiniBlocks: []*block.MiniBlock{
				{
					TxHashes:        [][]byte{},
					ReceiverShardID: 1,
					SenderShardID:   1,
					Type:            0,
				},
			},
		},
		Header: &block.Header{
			ShardID:   1,
			TimeStamp: 1234,
		},
		SignersIndexes:         []uint64{},
		NotarizedHeadersHashes: []string{},
		HeaderGasConsumption:   outport.HeaderGasConsumption{},
		TransactionsPool: &outport.Pool{
			Txs: map[string]data.TransactionHandlerWithGasUsedAndFee{
				"txhash1": &outport.TransactionHandlerWithGasAndFee{
					TransactionHandler: &transaction.Transaction{
						Nonce: 1,
					},
				},
			},
			Scrs: map[string]data.TransactionHandlerWithGasUsedAndFee{
				"scrHash1": &outport.TransactionHandlerWithGasAndFee{
					TransactionHandler: &transaction.Transaction{
						Nonce: 1,
					},
				},
			},
			Logs: []*data.LogData{
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
	}
}

func getShardV2Data() *notifierData.ArgsSaveBlock {
	return &notifierData.ArgsSaveBlock{
		HeaderType: "Header",
		ArgsSaveBlockData: notifierData.ArgsSaveBlockData{
			HeaderHash: []byte("headerHash93"),
			Body: &block.Body{
				MiniBlocks: []*block.MiniBlock{
					{
						TxHashes:        [][]byte{},
						ReceiverShardID: 1,
						SenderShardID:   1,
					},
				},
			},
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
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
		},
	}
}
