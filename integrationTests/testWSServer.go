package integrationTests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/gorilla/websocket"
)

type wsClient struct {
	wsConn     dispatcher.WSConnection
	httpServer *httptest.Server
}

// NewWSClient creates a new http server with websocket handler
func NewWSClient(h http.Handler) (*wsClient, error) {
	s := httptest.NewServer(h)
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	return &wsClient{
		wsConn:     ws,
		httpServer: s,
	}, nil
}

// SendSubscribeMessage will send subscribe message
func (ws *wsClient) SendSubscribeMessage(subscribeEvent *data.SubscribeEvent) error {
	m, err := json.Marshal(subscribeEvent)
	if err != nil {
		return err
	}

	if err := ws.wsConn.WriteMessage(websocket.BinaryMessage, m); err != nil {
		return err
	}

	return nil
}

func (ws *wsClient) ReadMessage() ([]byte, error) {
	_, m, err := ws.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ReceiveEvents will try to receive events
func (ws *wsClient) ReceiveEvents() ([]data.Event, error) {
	m, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	var reply data.WSEvent
	err = json.Unmarshal(m, &reply)
	if err != nil {
		return nil, err
	}

	var event []data.Event
	err = json.Unmarshal(reply.Data, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// ReceiveRevertBlock will try to receive revert block event
func (ws *wsClient) ReceiveRevertBlock() (*data.RevertBlock, error) {
	m, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	var reply data.WSEvent
	err = json.Unmarshal(m, &reply)
	if err != nil {
		return nil, err
	}

	var event data.RevertBlock
	err = json.Unmarshal(reply.Data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// ReceiveFinalized will try to receive finalized block event
func (ws *wsClient) ReceiveFinalized() (*data.FinalizedBlock, error) {
	m, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	var reply data.WSEvent
	err = json.Unmarshal(m, &reply)
	if err != nil {
		return nil, err
	}

	var event data.FinalizedBlock
	err = json.Unmarshal(reply.Data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// Close will close connection
func (ws *wsClient) Close() {
	ws.httpServer.Close()
	ws.wsConn.Close()
}
