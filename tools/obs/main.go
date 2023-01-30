package main

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
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

	data := getShardV1Data()

	err = httpClient.Post("/events/push", data)
	if err != nil {
		fmt.Println(fmt.Errorf("%w in eventNotifier.SaveBlock while posting block data", err))
		return
	}

	fmt.Println(data)
}

func getShardV1Data() *outport.ArgsSaveBlockData {
	return &outport.ArgsSaveBlockData{
		HeaderHash: []byte("hash2"),
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
