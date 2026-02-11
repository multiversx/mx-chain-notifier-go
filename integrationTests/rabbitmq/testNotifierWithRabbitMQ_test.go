package rabbitmq

import (
	"encoding/hex"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/stateChange"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/integrationTests"
	"github.com/multiversx/mx-chain-notifier-go/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// number of expected redis events
	// one event for each outport driver method: Save, Revert, Finalized
	numExpRedisEvents = 3

	// number of exected rabbitmq events
	// 5 events (logs & events, txs, scrs, full blocks events, state accesses) for Save method
	// + one event for each other outport driver method: Revert, Finalized
	numExpRabbitMQEvents = 7
)

var log = logger.GetOrCreate("integrationTests/rabbitmq")

func TestNotifierWithRabbitMQ(t *testing.T) {
	t.Run("with http observer connnector", func(t *testing.T) {
		testNotifierWithRabbitMQ(t, common.HTTPConnectorType, common.PayloadV1)
	})

	t.Run("with ws observer connnector", func(t *testing.T) {
		testNotifierWithRabbitMQ(t, common.WSObsConnectorType, common.PayloadV1)
	})
}

func TestNotifierWithRabbitMQV3(t *testing.T) {
	t.Run("with http observer connnector", func(t *testing.T) {
		testNotifierWithRabbitMQV3(t, common.HTTPConnectorType, common.PayloadV1)
	})

	t.Run("with ws observer connnector", func(t *testing.T) {
		testNotifierWithRabbitMQV3(t, common.WSObsConnectorType, common.PayloadV1)
	})
}

func testNotifierWithRabbitMQ(t *testing.T, observerType string, payloadVersion uint32) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.MainConfig.General.CheckDuplicates = true
	cfg.MainConfig.General.WithReadStateChanges = true
	notifier, err := integrationTests.NewTestNotifierWithRabbitMq(cfg.MainConfig)
	require.Nil(t, err)

	client, err := integrationTests.CreateObserverConnector(notifier.Facade, observerType, common.MessageQueuePublisherType, payloadVersion)
	require.Nil(t, err)

	// wait for components to start
	time.Sleep(time.Second * 5)

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

	integrationTests.WaitTimeout(t, wg, time.Second*5)

	assert.Equal(t, numExpRedisEvents, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, numExpRabbitMQEvents, len(notifier.RabbitMQClient.GetEntries()))
}

func testNotifierWithRabbitMQV3(t *testing.T, observerType string, payloadVersion uint32) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.MainConfig.General.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithRabbitMq(cfg.MainConfig)
	require.Nil(t, err)

	client, err := integrationTests.CreateObserverConnector(notifier.Facade, observerType, common.MessageQueuePublisherType, payloadVersion)
	require.Nil(t, err)

	// wait for components to start
	time.Sleep(time.Second * 5)

	_ = notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go pushEventsRequestV3(wg, client)
	go pushRevertRequestV3(wg, client)
	go pushFinalizedRequest(wg, client)

	// send requests again
	go pushEventsRequestV3(wg, client)
	go pushRevertRequestV3(wg, client)

	integrationTests.WaitTimeout(t, wg, time.Second*5)

	assert.Equal(t, numExpRedisEvents, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, numExpRabbitMQEvents, len(notifier.RabbitMQClient.GetEntries()))
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
			hex.EncodeToString([]byte("hash1")): {
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
			hex.EncodeToString([]byte("hash2")): {
				SmartContractResult: &smartContractResult.SmartContractResult{
					Nonce: 2,
				},
				FeeInfo: &outport.FeeInfo{
					GasUsed: 2,
				},
				ExecutionOrder: 3,
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
	}

	stateAccesses := make(map[string]*stateChange.StateAccesses)
	stateAccesses["txHash1"] = &stateChange.StateAccesses{
		StateAccess: []*stateChange.StateAccess{
			{
				MainTrieKey: []byte("mainTrieKey1"),
				MainTrieVal: []byte("mainTrieVal1"),
			},
			{
				MainTrieKey: []byte("mainTrieKey2"),
				MainTrieVal: []byte("mainTrieVal2"),
			},
		},
	}
	stateAccesses["txHash2"] = &stateChange.StateAccesses{}

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
		StateAccessesForBlock: map[string]*outport.StateAccessesForBlock{
			hex.EncodeToString([]byte("headerHash1")): {
				StateAccesses: stateAccesses,
			},
		},
	}

	err := webServer.PushEventsRequest(saveBlockData)
	log.LogIfError(err)

	if err == nil {
		wg.Done()
	}
}

func pushEventsRequestV3(
	wg *sync.WaitGroup,
	webServer integrationTests.ObserverConnector,
) {
	marshaller := &marshal.JsonMarshalizer{}
	blockData, err := testdata.NewBlockData(marshaller)
	log.LogIfError(err)

	saveBlockData := blockData.OutportBlockV2()

	err = webServer.PushEventsRequest(saveBlockData)
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

func pushRevertRequestV3(wg *sync.WaitGroup, webServer integrationTests.ObserverConnector) {
	header := &block.HeaderV3{
		Nonce: 1,
	}
	headerBytes, _ := json.Marshal(header)
	blockData := &outport.BlockData{
		HeaderBytes: headerBytes,
		HeaderType:  string(core.ShardHeaderV3),
		HeaderHash:  []byte("headerHash3"),
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
