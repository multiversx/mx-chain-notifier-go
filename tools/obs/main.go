package main

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go/outport/notifier"
)

func main() {
	args := notifier.HTTPClientWrapperArgs{
		UseAuthorization:  false,
		BaseUrl:           "http://localhost:5000",
		RequestTimeoutSec: 10,
	}
	httpClient, err := notifier.NewHTTPWrapperClient(args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := &outport.ArgsSaveBlockData{
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
		Header: &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
			},
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

	err = httpClient.Post("/events/push", data)
	if err != nil {
		fmt.Println(fmt.Errorf("%w in eventNotifier.SaveBlock while posting block data", err))
		return
	}

	fmt.Println(data)
}
