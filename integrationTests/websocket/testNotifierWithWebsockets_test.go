package websocket

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

	header := &block.HeaderV2{
		Header: &block.Header{
			ShardID:   1,
			TimeStamp: 1234,
		},
	}
	headerBytes, _ := json.Marshal(header)
	saveBlockData := &outport.OutportBlock{
		TransactionPool: &outport.TransactionPool{
			Logs: []*outport.LogData{
				{
					Log: &transaction.Log{
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
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  []byte("headerHash"),
			Body: &block.Body{
				MiniBlocks: make([]*block.MiniBlock, 1),
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
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

	err = webServer.PushEventsRequest(saveBlockData)
	require.Nil(t, err)

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

	header := &block.HeaderV2{
		Header: &block.Header{
			ShardID:   1,
			TimeStamp: 1234,
		},
	}
	headerBytes, _ := json.Marshal(header)
	saveBlockData := &outport.OutportBlock{
		TransactionPool: &outport.TransactionPool{
			Logs: []*outport.LogData{
				{
					Log: &transaction.Log{
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
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  headerHash,
			Body: &block.Body{
				MiniBlocks: make([]*block.MiniBlock, 1),
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
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

	err = webServer.PushEventsRequest(saveBlockData)
	require.Nil(t, err)

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

	header := &block.HeaderV2{
		Header: &block.Header{
			Nonce: 1,
		},
	}
	headerBytes, _ := json.Marshal(header)
	blockEvents := &outport.BlockData{
		HeaderBytes: headerBytes,
		HeaderType:  string(core.ShardHeaderV2),
		HeaderHash:  []byte("hash1"),
	}

	expReply := &data.RevertBlock{
		Hash:  hex.EncodeToString([]byte("hash1")),
		Nonce: 1,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveRevertBlock()
		require.Nil(t, err)

		require.Equal(t, expReply, reply)
		wg.Done()
	}()

	time.Sleep(time.Second)

	err = webServer.RevertEventsRequest(blockEvents)
	require.Nil(t, err)

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

	expReply := &data.FinalizedBlock{
		Hash: "hash1",
	}

	blockEvents := &outport.FinalizedBlock{
		HeaderHash: []byte("hash1"),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveFinalized()
		require.Nil(t, err)

		require.Equal(t, expReply, reply)
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
	txs := map[string]*outport.TxInfo{
		"hash1": {
			Transaction: &transaction.Transaction{
				Nonce: 1,
			},
		},
	}

	header := &block.HeaderV2{
		Header: &block.Header{
			ShardID:   1,
			TimeStamp: 1234,
		},
	}
	headerBytes, _ := json.Marshal(header)
	saveBlockData := &outport.OutportBlock{
		TransactionPool: &outport.TransactionPool{
			Transactions: txs,
		},
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  blockHash,
			Body: &block.Body{
				MiniBlocks: make([]*block.MiniBlock, 1),
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
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

	webServer.PushEventsRequest(saveBlockData)

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
	scrs := map[string]*outport.SCRInfo{
		"hash2": {
			SmartContractResult: &smartContractResult.SmartContractResult{
				Nonce: 2,
			},
		},
	}
	header := &block.HeaderV2{
		Header: &block.Header{
			ShardID:   1,
			TimeStamp: 1234,
		},
	}
	headerBytes, _ := json.Marshal(header)
	blockEvents := &outport.OutportBlock{
		TransactionPool: &outport.TransactionPool{
			SmartContractResults: scrs,
		},
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  blockHash,
			Body: &block.Body{
				MiniBlocks: make([]*block.MiniBlock, 1),
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
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

	header := &block.HeaderV2{
		Header: &block.Header{
			Nonce: 1,
		},
	}
	headerBytes, _ := json.Marshal(header)
	revertBlock := &outport.BlockData{
		HeaderBytes: headerBytes,
		HeaderType:  string(core.ShardHeaderV2),
		HeaderHash:  []byte("hash1"),
	}
	expRevertBlock := &data.RevertBlock{
		Hash:  hex.EncodeToString([]byte("hash1")),
		Nonce: 1,
	}

	finalizedBlock := &outport.FinalizedBlock{
		HeaderHash: []byte("hash1"),
	}
	expFinalizedBlock := &data.FinalizedBlock{
		Hash: "hash1",
	}

	addr := []byte("addr1")
	events := []data.Event{
		{
			Address: hex.EncodeToString(addr),
		},
	}

	txs := map[string]*outport.TxInfo{
		"hash1": {
			Transaction: &transaction.Transaction{
				Nonce: 1,
			},
		},
	}
	scrs := map[string]*outport.SCRInfo{
		"hash2": {
			SmartContractResult: &smartContractResult.SmartContractResult{
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

	expTxsWithOrder := map[string]*outport.TxInfo{
		"hash1": {
			Transaction: &transaction.Transaction{
				Nonce: 1,
			},
		},
	}
	expScrsWithOrder := map[string]*outport.SCRInfo{
		"hash2": {
			SmartContractResult: &smartContractResult.SmartContractResult{
				Nonce: 2,
			},
		},
	}
	expBlockEvents := data.BlockEventsWithOrder{
		Hash:      hex.EncodeToString(blockHash),
		ShardID:   1,
		TimeStamp: 1234,
		Events:    events,
		Txs:       expTxsWithOrder,
		Scrs:      expScrsWithOrder,
	}

	header = &block.HeaderV2{
		Header: &block.Header{
			ShardID:   1,
			TimeStamp: 1234,
		},
	}
	headerBytes, _ = json.Marshal(header)
	blockEvents := &outport.OutportBlock{
		TransactionPool: &outport.TransactionPool{
			Transactions:         txs,
			SmartContractResults: scrs,
			Logs: []*outport.LogData{
				{
					Log: &transaction.Log{
						Events: []*transaction.Event{
							{
								Address: addr,
							},
						},
					},
				},
			},
		},
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV2),
			HeaderHash:  blockHash,
			Body: &block.Body{
				MiniBlocks: make([]*block.MiniBlock, 1),
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
	}

	numEvents := 6
	wg := &sync.WaitGroup{}
	wg.Add(numEvents)

	go func(wg *sync.WaitGroup) {
		for i := 0; i < numEvents; i++ {
			m, err := ws.ReadMessage()
			require.Nil(t, err)

			var reply data.WebSocketEvent
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
				assert.Equal(t, expRevertBlock, event)
				wg.Done()
			case common.BlockEvents:
				var event data.BlockEventsWithOrder
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, expBlockEvents, event)
				wg.Done()
			case common.FinalizedBlockEvents:
				var event *data.FinalizedBlock
				_ = json.Unmarshal(reply.Data, &event)
				assert.Equal(t, expFinalizedBlock, event)
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
