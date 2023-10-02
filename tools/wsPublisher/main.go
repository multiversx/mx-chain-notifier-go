package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

const (
	hostURL = "127.0.0.1:5000"
	wsPath  = "/hub/ws"
)

func main() {
	ws, err := NewWSClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockEvents,
			},
			{
				EventType: common.BlockScrs,
			},
			{
				EventType: common.BlockTxs,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	for {
		m, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		var reply data.WebSocketEvent
		err = json.Unmarshal(m, &reply)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		switch reply.Type {
		case common.BlockEvents:
			var event data.BlockEventsWithOrder
			_ = json.Unmarshal(reply.Data, &event)
			fmt.Println(event)
		case common.RevertBlockEvents:
			var event *data.RevertBlock
			_ = json.Unmarshal(reply.Data, &event)
			fmt.Println(event)
		case common.FinalizedBlockEvents:
			var event *data.FinalizedBlock
			_ = json.Unmarshal(reply.Data, &event)
			fmt.Println(event)
		case common.BlockTxs:
			var event *data.BlockTxs
			_ = json.Unmarshal(reply.Data, &event)
			fmt.Println(event)
		case common.BlockScrs:
			var event data.BlockScrs
			_ = json.Unmarshal(reply.Data, &event)
			fmt.Println(event.Hash)
		default:
			fmt.Println("invalid message type")
		}
	}
}

type wsClient struct {
	wsConn    dispatcher.WSConnection
	mutWsConn sync.RWMutex
}

// NewWSClient creates a new websocket client
func NewWSClient() (*wsClient, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   hostURL,
		Path:   wsPath,
	}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &wsClient{
		wsConn: ws,
	}, nil
}

// SendSubscribeMessage will send subscribe message
func (ws *wsClient) SendSubscribeMessage(subscribeEvent *data.SubscribeEvent) error {
	ws.mutWsConn.Lock()
	defer ws.mutWsConn.Unlock()

	m, err := json.Marshal(subscribeEvent)
	if err != nil {
		return err
	}

	return ws.wsConn.WriteMessage(websocket.BinaryMessage, m)
}

// ReadMessage will read the received message
func (ws *wsClient) ReadMessage() ([]byte, error) {
	ws.mutWsConn.Lock()
	defer ws.mutWsConn.Unlock()

	_, m, err := ws.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Close will close connection
func (ws *wsClient) Close() {
	ws.mutWsConn.Lock()
	defer ws.mutWsConn.Unlock()

	ws.wsConn.Close()
}
