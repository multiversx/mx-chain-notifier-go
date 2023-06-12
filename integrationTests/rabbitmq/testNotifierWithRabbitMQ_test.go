package rabbitmq

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifierWithRabbitMQ(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.ConnectorApi.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithRabbitMq(cfg)
	require.Nil(t, err)

	client, server, err := integrationTests.CreateObserverConnector(notifier.Facade, "http", common.MessageQueueAPIType)
	require.Nil(t, err)
	defer func() {
		if server != nil {
			server.Close()
		}
		if client != nil {
			client.Close()
		}
	}()

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	mutResponses := &sync.Mutex{}
	responses := make(map[int]int)

	go pushEventsRequest(client, mutResponses, responses)
	go pushRevertRequest(client, mutResponses, responses)
	go pushFinalizedRequest(client, mutResponses, responses)

	// send requests again
	go pushEventsRequest(client, mutResponses, responses)
	go pushRevertRequest(client, mutResponses, responses)

	time.Sleep(time.Second)

	mutResponses.Lock()
	assert.Equal(t, 5, responses[200])
	mutResponses.Unlock()

	assert.Equal(t, 6, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, 6, len(notifier.RabbitMQClient.GetEntries()))
}

func pushEventsRequest(webServer integrationTests.ObserverConnector, mutResponses *sync.Mutex, responses map[int]int) {
	header := &block.HeaderV2{
		Header: &block.Header{
			Nonce: 1,
		},
	}
	headerBytes, _ := json.Marshal(header)

	txPool := &outport.TransactionPool{
		Transactions: map[string]*outport.TxInfo{
			"hash1": {
				Transaction: &transaction.Transaction{
					Nonce: 1,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 1,
				},
				ExecutionOrder: 1,
			},
		},
		SmartContractResults: map[string]*outport.SCRInfo{
			"hash2": {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 2,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 2,
				},
				ExecutionOrder: 3,
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
	}

	saveBlockData := &outport.OutportBlock{
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  []byte("headerHash1"),
			Body: &block.Body{
				MiniBlocks: make([]*block.MiniBlock, 1),
			},
		},
		TransactionPool:      txPool,
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
	}

	err := webServer.PushEventsRequest(saveBlockData)

	mutResponses.Lock()
	if err == nil {
		responses[200]++
	}
	mutResponses.Unlock()
}

func pushRevertRequest(webServer integrationTests.ObserverConnector, mutResponses *sync.Mutex, responses map[int]int) {
	header := &block.HeaderV2{
		Header: &block.Header{
			Nonce: 1,
		},
	}
	headerBytes, _ := json.Marshal(header)
	blockData := &outport.BlockData{
		HeaderBytes: headerBytes,
		HeaderType:  string(core.ShardHeaderV2),
		HeaderHash:  []byte("headerHash2"),
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
	}
	err := webServer.RevertEventsRequest(blockData)

	mutResponses.Lock()
	if err == nil {
		responses[200]++
	}
	mutResponses.Unlock()
}

func pushFinalizedRequest(webServer integrationTests.ObserverConnector, mutResponses *sync.Mutex, responses map[int]int) {
	blockEvents := &outport.FinalizedBlock{
		HeaderHash: []byte("headerHash3"),
	}
	err := webServer.FinalizedEventsRequest(blockEvents)

	mutResponses.Lock()
	if err == nil {
		responses[200]++
	}
	mutResponses.Unlock()
}
