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
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = logger.GetOrCreate("integrationTests/rabbitmq")

func TestNotifierWithRabbitMQ(t *testing.T) {
	t.Run("with http observer connnector + payload version 0", func(t *testing.T) {
		testNotifierWithRabbitMQ(t, common.HTTPConnectorType, common.PayloadV0)
	})

	t.Run("with http observer connnector", func(t *testing.T) {
		testNotifierWithRabbitMQ(t, common.HTTPConnectorType, common.PayloadV1)
	})

	t.Run("with ws observer connnector", func(t *testing.T) {
		testNotifierWithRabbitMQ(t, common.WSObsConnectorType, common.PayloadV1)
	})
}

func testNotifierWithRabbitMQ(t *testing.T, observerType string, payloadVersion uint32) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.MainConfig.General.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithRabbitMq(cfg.MainConfig)
	require.Nil(t, err)

	client, err := integrationTests.CreateObserverConnector(notifier.Facade, observerType, common.MessageQueuePublisherType, payloadVersion)
	require.Nil(t, err)

	_ = notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go pushEventsRequest(wg, client)
	go pushRevertRequest(wg, client)
	go pushFinalizedRequest(wg, client)

	// send requests again
	go pushEventsRequest(wg, client)
	go pushRevertRequest(wg, client)

	integrationTests.WaitTimeout(t, wg, time.Second*2)

	assert.Equal(t, 3, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, 6, len(notifier.RabbitMQClient.GetEntries()))
}

func pushEventsRequest(wg *sync.WaitGroup, webServer integrationTests.ObserverConnector) {
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
				MiniBlocks: []*block.MiniBlock{
					&block.MiniBlock{},
				},
			},
		},
		TransactionPool:      txPool,
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
	}

	err := webServer.PushEventsRequest(saveBlockData)
	log.LogIfError(err)

	if err == nil {
		wg.Done()
	}
}

func pushRevertRequest(wg *sync.WaitGroup, webServer integrationTests.ObserverConnector) {
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
	}
	err := webServer.RevertEventsRequest(blockData)
	log.LogIfError(err)

	if err == nil {
		wg.Done()
	}
}

func pushFinalizedRequest(wg *sync.WaitGroup, webServer integrationTests.ObserverConnector) {
	blockEvents := &outport.FinalizedBlock{
		HeaderHash: []byte("headerHash3"),
	}
	err := webServer.FinalizedEventsRequest(blockEvents)
	log.LogIfError(err)

	if err == nil {
		wg.Done()
	}
}
