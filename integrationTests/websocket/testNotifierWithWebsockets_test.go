package websocket

import (
	"encoding/hex"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/integrationTests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifierWithWebsockets_PushEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.WSAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	ws, err := integrationTests.NewWSClient(notifier.WSHandler)
	require.Nil(t, err)
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.PushBlockEvents,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	addr := []byte("addr1")
	events := []data.Event{
		{
			Address: hex.EncodeToString(addr),
			TxHash:  hex.EncodeToString([]byte("txHash1")),
		},
	}
	saveBlockData := data.ArgsSaveBlockData{
		HeaderHash: []byte("hash1"),
		TransactionsPool: &data.TransactionsPool{
			Logs: []*data.LogData{
				{
					LogHandler: &transaction.Log{
						Events: []*transaction.Event{
							{
								Address: addr,
							},
						},
					},
					TxHash: "txHash1",
				},
			},
		},
	}
	blockEvents := &data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: saveBlockData,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveEvents()
		require.Nil(t, err)

		require.Equal(t, events, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	resp := webServer.PushEventsRequest(blockEvents)
	require.NotNil(t, resp)

	wg.Wait()
}

func TestNotifierWithWebsockets_RevertEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.WSAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	ws, err := integrationTests.NewWSClient(notifier.WSHandler)
	require.Nil(t, err)
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.RevertBlockEvents,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	blockEvents := &data.RevertBlock{
		Hash:  "hash1",
		Nonce: 1,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveRevertBlock()
		require.Nil(t, err)

		require.Equal(t, blockEvents, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	resp := webServer.RevertEventsRequest(blockEvents)
	require.NotNil(t, resp)

	wg.Wait()
}

func TestNotifierWithWebsockets_FinalizedEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.WSAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	ws, err := integrationTests.NewWSClient(notifier.WSHandler)
	require.Nil(t, err)
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.FinalizedBlockEvents,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	blockEvents := &data.FinalizedBlock{
		Hash: "hash1",
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveFinalized()
		require.Nil(t, err)

		require.Equal(t, blockEvents, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	webServer.FinalizedEventsRequest(blockEvents)

	wg.Wait()
}

func TestNotifierWithWebsockets_TxsEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.WSAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	ws, err := integrationTests.NewWSClient(notifier.WSHandler)
	require.Nil(t, err)
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockTxs,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	blockHash := []byte("hash1")
	txs := map[string]data.TransactionWithOrder{
		"hash1": {
			TransactionHandler: &transaction.Transaction{
				Nonce: 1,
			},
		},
	}
	saveBlockData := data.ArgsSaveBlockData{
		HeaderHash: blockHash,
		TransactionsPool: &data.TransactionsPool{
			Txs: txs,
		},
	}
	blockEvents := &data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: saveBlockData,
	}

	expTxs := map[string]*transaction.Transaction{
		"hash1": {
			Nonce: 1,
		},
	}
	expBlockTxs := &data.BlockTxs{
		Hash: hex.EncodeToString(blockHash),
		Txs:  expTxs,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveTxs()
		require.Nil(t, err)

		require.Equal(t, expBlockTxs, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	webServer.PushEventsRequest(blockEvents)

	wg.Wait()
}

func TestNotifierWithWebsockets_ScrsEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.WSAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	ws, err := integrationTests.NewWSClient(notifier.WSHandler)
	require.Nil(t, err)
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockScrs,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	blockHash := []byte("hash1")
	scrs := map[string]data.SmartContractResultWithOrder{
		"hash2": {
			TransactionHandler: &smartContractResult.SmartContractResult{
				Nonce: 2,
			},
		},
	}
	saveBlockData := data.ArgsSaveBlockData{
		HeaderHash: blockHash,
		TransactionsPool: &data.TransactionsPool{
			Scrs: scrs,
		},
	}
	blockEvents := &data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: saveBlockData,
	}

	expScrs := map[string]*smartContractResult.SmartContractResult{
		"hash2": {
			Nonce: 2,
		},
	}
	expBlockScrs := &data.BlockScrs{
		Hash: hex.EncodeToString(blockHash),
		Scrs: expScrs,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveScrs()
		require.Nil(t, err)

		require.Equal(t, expBlockScrs, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	webServer.PushEventsRequest(blockEvents)

	wg.Wait()
}

func TestNotifierWithWebsockets_AllEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.ConnectorApi.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.WSAPIType)

	notifier.Publisher.Run()
	defer notifier.Publisher.Close()

	ws, err := integrationTests.NewWSClient(notifier.WSHandler)
	require.Nil(t, err)
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.PushBlockEvents,
			},
			{
				EventType: common.RevertBlockEvents,
			},
			{
				EventType: common.FinalizedBlockEvents,
			},
			{
				EventType: common.BlockTxs,
			},
			{
				EventType: common.BlockScrs,
			},
		},
	}

	err = ws.SendSubscribeMessage(subscribeEvent)
	require.Nil(t, err)

	revertBlock := &data.RevertBlock{
		Hash:  "hash1",
		Nonce: 1,
	}
	finalizedBlock := &data.FinalizedBlock{
		Hash: "hash1",
	}

	addr := []byte("addr1")
	events := []data.Event{
		{
			Address: hex.EncodeToString(addr),
		},
	}

	txs := map[string]data.TransactionWithOrder{
		"hash1": {
			TransactionHandler: &transaction.Transaction{
				Nonce: 1,
			},
		},
	}
	scrs := map[string]data.SmartContractResultWithOrder{
		"hash2": {
			TransactionHandler: &smartContractResult.SmartContractResult{
				Nonce: 2,
			},
		},
	}
	blockHash := []byte("hash1")

	expTxs := map[string]*transaction.Transaction{
		"hash1": {
			Nonce: 1,
		},
	}
	blockTxs := &data.BlockTxs{
		Hash: hex.EncodeToString(blockHash),
		Txs:  expTxs,
	}

	expScrs := map[string]*smartContractResult.SmartContractResult{
		"hash2": {
			Nonce: 2,
		},
	}
	blockScrs := &data.BlockScrs{
		Hash: hex.EncodeToString(blockHash),
		Scrs: expScrs,
	}

	saveBlockData := data.ArgsSaveBlockData{
		HeaderHash: []byte(blockHash),
		TransactionsPool: &data.TransactionsPool{
			Txs:  txs,
			Scrs: scrs,
			Logs: []*data.LogData{
				{
					LogHandler: &transaction.Log{
						Events: []*transaction.Event{
							{
								Address: addr,
							},
						},
					},
				},
			},
		},
	}
	blockEvents := &data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: saveBlockData,
	}

	numEvents := 5
	wg := &sync.WaitGroup{}
	wg.Add(numEvents)

	go func(wg *sync.WaitGroup) {
		for i := 0; i < numEvents; i++ {
			m, err := ws.ReadMessage()
			require.Nil(t, err)

			var reply data.WSEvent
			err = json.Unmarshal(m, &reply)
			require.Nil(t, err)

			switch reply.Type {
			case common.PushBlockEvents:
				var event []data.Event
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, events, event)
				wg.Done()
			case common.RevertBlockEvents:
				var event *data.RevertBlock
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, revertBlock, event)
				wg.Done()
			case common.FinalizedBlockEvents:
				var event *data.FinalizedBlock
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, finalizedBlock, event)
				wg.Done()
			case common.BlockTxs:
				var event *data.BlockTxs
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, blockTxs, event)
				wg.Done()
			case common.BlockScrs:
				var event *data.BlockScrs
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, blockScrs, event)
				wg.Done()
			default:
				t.Errorf("invalid message type")
			}
		}
	}(wg)

	time.Sleep(time.Second)

	go webServer.PushEventsRequest(blockEvents)
	go webServer.FinalizedEventsRequest(finalizedBlock)
	go webServer.RevertEventsRequest(revertBlock)

	wg.Wait()

	assert.Equal(t, numEvents, len(notifier.RedisClient.GetEntries()))
}
