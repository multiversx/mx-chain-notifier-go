package websocket

import (
	"encoding/hex"
	"encoding/json"
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
				EventType: common.PushLogsAndEvents,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	addr := []byte("addr1")
	events := []data.Event{
		{
			Address: hex.EncodeToString(addr),
			TxHash:  "txHash1",
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
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
		Header: &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
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

	if waitTimeout(t, wg, time.Second*2) {
		assert.Fail(t, "timeout when handling websocket events")
	}
}

func TestNotifierWithWebsockets_BlockEvents(t *testing.T) {
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
				EventType: common.BlockEvents,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	headerHash := []byte("hash1")
	addr := []byte("addr1")
	events := []data.Event{
		{
			Address: hex.EncodeToString(addr),
			TxHash:  "txHash1",
		},
	}
	expBlockEvents := &data.BlockEventsWithOrder{
		Hash:      hex.EncodeToString(headerHash),
		ShardID:   1,
		TimeStamp: 1234,
		Events:    events,
	}

	saveBlockData := data.ArgsSaveBlockData{
		HeaderHash: headerHash,
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
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
		Header: &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
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
		reply, err := ws.ReceiveBlockEventsData()
		require.Nil(t, err)

		require.Equal(t, expBlockEvents, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	resp := webServer.PushEventsRequest(blockEvents)
	require.NotNil(t, resp)

	if waitTimeout(t, wg, time.Second*2) {
		assert.Fail(t, "timeout when handling websocket events")
	}
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

	if waitTimeout(t, wg, time.Second*2) {
		assert.Fail(t, "timeout when handling websocket events")
	}
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

	if waitTimeout(t, wg, time.Second*2) {
		assert.Fail(t, "timeout when handling websocket events")
	}
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
	txs := map[string]data.TransactionWrapped{
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
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
		Header: &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
			},
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

	if waitTimeout(t, wg, time.Second*2) {
		assert.Fail(t, "timeout when handling websocket events")
	}
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
	scrs := map[string]data.SmartContractResultWrapped{
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
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
		Header: &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
			},
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

	if waitTimeout(t, wg, time.Second*2) {
		assert.Fail(t, "timeout when handling websocket events")
	}
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
				EventType: common.PushLogsAndEvents,
			},
			{
				EventType: common.BlockEvents,
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

	txs := map[string]data.TransactionWrapped{
		"hash1": {
			TransactionHandler: &transaction.Transaction{
				Nonce: 1,
			},
		},
	}
	scrs := map[string]data.SmartContractResultWrapped{
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

	expBlockEvents := data.BlockEventsWithOrder{
		Hash:      hex.EncodeToString(blockHash),
		ShardID:   1,
		TimeStamp: 1234,
		Events:    events,
		Txs:       txs,
		Scrs:      scrs,
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
		Body: &block.Body{
			MiniBlocks: make([]*block.MiniBlock, 1),
		},
		Header: &block.HeaderV2{
			Header: &block.Header{
				ShardID:   1,
				TimeStamp: 1234,
			},
		},
	}
	blockEvents := &data.ArgsSaveBlock{
		HeaderType:        "HeaderV2",
		ArgsSaveBlockData: saveBlockData,
	}

	numEvents := 6
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
			case common.PushLogsAndEvents:
				var event []data.Event
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, events, event)
				wg.Done()
			case common.RevertBlockEvents:
				var event *data.RevertBlock
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, revertBlock, event)
				wg.Done()
			case common.BlockEvents:
				var event data.BlockEventsWithOrder
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, expBlockEvents, event)
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

	if waitTimeout(t, wg, time.Second*4) {
		assert.Fail(t, "timeout when handling websocket events")
	}

	assert.Equal(t, numEvents, len(notifier.RedisClient.GetEntries()))
	assert.Equal(t, numEvents, len(notifier.RedisClient.GetEntries()))
}

// waitTimeout returns true if work group waiting timed out
func waitTimeout(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ch:
		return false
	case <-timer.C:
		return true
	}
}
