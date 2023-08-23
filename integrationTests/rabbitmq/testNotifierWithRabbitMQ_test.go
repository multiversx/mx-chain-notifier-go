package rabbitmq

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifierWithRabbitMQ(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.GeneralConfig.ConnectorApi.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithRabbitMq(cfg.GeneralConfig)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.MessageQueueAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	mutResponses := &sync.Mutex{}
	responses := make(map[int]int)

	go pushEventsRequest(webServer, mutResponses, responses)
	go pushRevertRequest(webServer, mutResponses, responses)
	go pushFinalizedRequest(webServer, mutResponses, responses)

	// send requests again
	go pushEventsRequest(webServer, mutResponses, responses)
	go pushRevertRequest(webServer, mutResponses, responses)

	time.Sleep(time.Second)

	mutResponses.Lock()
	assert.Equal(t, 5, responses[http.StatusOK])
	mutResponses.Unlock()

	assert.Equal(t, 6, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, 6, len(notifier.RabbitMQClient.GetEntries()))
}

func pushEventsRequest(webServer *integrationTests.TestWebServer, mutResponses *sync.Mutex, responses map[int]int) {
	blockEvents := data.ArgsSaveBlockData{
		HeaderHash: []byte("hash1"),
		Header: &block.HeaderV2{
			Header: &block.Header{
				Nonce:     1,
				ShardID:   2,
				TimeStamp: 1234,
			},
		},
		TransactionsPool: &data.TransactionsPool{
			Txs: map[string]*data.NodeTransaction{
				"txHash1": {
					TransactionHandler: &transaction.Transaction{
						Nonce: 1,
					},
					ExecutionOrder: 1,
				},
			},
			Scrs: map[string]*data.NodeSmartContractResult{
				"scrHash1": {
					TransactionHandler: &smartContractResult.SmartContractResult{
						Nonce: 2,
					},
				},
			},
			Logs: []*data.LogData{
				{
					LogHandler: &transaction.Log{
						Address: []byte("addr1"),
						Events:  []*transaction.Event{},
					},
					TxHash: "txHash1",
				},
			},
		},
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
	}
	saveBlockData := &data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: blockEvents,
	}

	resp := webServer.PushEventsRequest(saveBlockData)

	mutResponses.Lock()
	responses[resp.Code]++
	mutResponses.Unlock()
}

func pushRevertRequest(webServer *integrationTests.TestWebServer, mutResponses *sync.Mutex, responses map[int]int) {
	blockEvents := &data.RevertBlock{
		Hash:  "hash2",
		Nonce: 1,
	}
	resp := webServer.RevertEventsRequest(blockEvents)

	mutResponses.Lock()
	responses[resp.Code]++
	mutResponses.Unlock()
}

func pushFinalizedRequest(webServer *integrationTests.TestWebServer, mutResponses *sync.Mutex, responses map[int]int) {
	blockEvents := &data.FinalizedBlock{
		Hash: "hash3",
	}
	resp := webServer.FinalizedEventsRequest(blockEvents)

	mutResponses.Lock()
	responses[resp.Code]++
	mutResponses.Unlock()
}
