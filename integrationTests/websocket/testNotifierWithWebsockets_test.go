package websocket

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

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

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.MessageQueueAPIType)

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
	blockEvents := &data.BlockEvents{
		Hash:   "hash1",
		Events: events,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reply, err := ws.ReceiveEvents()
		require.Nil(t, err)

		assert.Equal(t, events, reply)
		wg.Done()
	}()

	resp := webServer.PushEventsRequest(blockEvents)
	require.NotNil(t, resp)

	wg.Wait()
}

func TestNotifierWithWebsockets_RevertEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.MessageQueueAPIType)

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

	resp := webServer.RevertEventsRequest(blockEvents)
	require.NotNil(t, resp)

	wg.Wait()
}

func TestNotifierWithWebsockets_FinalizedEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.MessageQueueAPIType)

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

	webServer.FinalizedEventsRequest(blockEvents)

	wg.Wait()
}

func TestNotifierWithWebsockets_AllEvents(t *testing.T) {
	cfg := integrationTests.GetDefaultConfigs()
	cfg.ConnectorApi.CheckDuplicates = true
	notifier, err := integrationTests.NewTestNotifierWithWS(cfg)
	require.Nil(t, err)

	webServer := integrationTests.NewTestWebServer(notifier.Facade, common.MessageQueueAPIType)

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
	blockEvents := &data.BlockEvents{
		Hash:   "hash1",
		Events: events,
	}

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func(wg *sync.WaitGroup) {
		for i := 0; i < 3; i++ {
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

	assert.Equal(t, 3, len(notifier.RedisClient.GetEntries()))
}
