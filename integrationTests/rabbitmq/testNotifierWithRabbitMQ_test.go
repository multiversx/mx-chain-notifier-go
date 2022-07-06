package rabbitmq

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifierWithRabbitMQ(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.ConnectorApi.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithRabbitMq(cfg)
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

	assert.Equal(t, 5, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, 5, len(notifier.RabbitMQClient.GetEntries()))
}

func pushEventsRequest(webServer *integrationTests.TestWebServer, mutResponses *sync.Mutex, responses map[int]int) {
	blockEvents := &data.SaveBlockData{
		Hash: "hash1",
		Txs: map[string]transaction.Transaction{
			"txHash1": {
				Nonce: 1,
			},
		},
		Scrs: map[string]smartContractResult.SmartContractResult{
			"scrHash1": {
				Nonce: 2,
			},
		},
		LogEvents: []data.Event{
			{
				Address: "addr1",
				TxHash:  "txHash1",
			},
		},
	}
	resp := webServer.PushEventsRequest(blockEvents)

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
