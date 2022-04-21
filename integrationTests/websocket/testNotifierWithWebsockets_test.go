package websocket

import (
	"encoding/json"
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
				EventType: "all_events",
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	events := []data.Event{
		{
			Address: "addr1",
		},
	}
	blockEvents := &data.SaveBlockData{
		Hash:      "hash1",
		LogEvents: events,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveEvents()
		require.Nil(t, err)

		assert.Equal(t, events, reply)
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
				EventType: "revert_events",
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

		assert.Equal(t, blockEvents, reply)
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
				EventType: "finalized_events",
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

		assert.Equal(t, blockEvents, reply)
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
				EventType: "txs_events",
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	blockHash := "hash1"
	txs := map[string]transaction.Transaction{
		"txhash1": {
			Nonce: 1,
		},
	}
	blockEvents := &data.SaveBlockData{
		Hash: blockHash,
		Txs:  txs,
	}
	expBlockTxs := &data.BlockTxs{
		Hash: blockHash,
		Txs:  txs,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveTxs()
		require.Nil(t, err)

		assert.Equal(t, expBlockTxs, reply)
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
				EventType: "scrs_events",
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	blockHash := "hash1"
	scrs := map[string]smartContractResult.SmartContractResult{
		"hash2": {
			Nonce: 2,
		},
	}
	blockEvents := &data.SaveBlockData{
		Hash: blockHash,
		Scrs: scrs,
	}
	expBlockScrs := &data.BlockScrs{
		Hash: blockHash,
		Scrs: scrs,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveScrs()
		require.Nil(t, err)

		assert.Equal(t, expBlockScrs, reply)
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
				EventType: "all_events",
			},
			{
				EventType: "revert_events",
			},
			{
				EventType: "finalized_events",
			},
			{
				EventType: "txs_events",
			},
			{
				EventType: "scrs_events",
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

	events := []data.Event{
		{
			Address: "addr1",
		},
	}

	txs := map[string]transaction.Transaction{
		"txhash1": {
			Nonce: 1,
		},
	}
	scrs := map[string]smartContractResult.SmartContractResult{
		"txhash2": {
			Nonce: 2,
		},
	}
	blockHash := "hash1"
	blockTxs := &data.BlockTxs{
		Hash: blockHash,
		Txs:  txs,
	}
	blockScrs := &data.BlockScrs{
		Hash: blockHash,
		Scrs: scrs,
	}
	blockEvents := &data.SaveBlockData{
		Hash:      blockHash,
		LogEvents: events,
		Txs:       txs,
		Scrs:      scrs,
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
			case common.BlockTxsEvents:
				var event *data.BlockTxs
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, blockTxs, event)
				wg.Done()
			case common.BlockScrsEvents:
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

	assert.Equal(t, 5, len(notifier.RedisClient.GetEntries()))
}
