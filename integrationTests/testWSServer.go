package integrationTests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/gorilla/websocket"
)

type wsServer struct {
	wsConn     dispatcher.WSConnection
	httpServer *httptest.Server
}

// NewWSServer creates a new http server with websocket handler
func NewWSServer(h http.Handler) (*wsServer, error) {
	s := httptest.NewServer(h)
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	return &wsServer{
		wsConn:     ws,
		httpServer: s,
	}, nil
}

func (ws *wsServer) SendSubscribeMessage(subscribeEvent *data.SubscribeEvent) error {
	m, err := json.Marshal(subscribeEvent)
	if err != nil {
		return err
	}

	if err := ws.wsConn.WriteMessage(websocket.BinaryMessage, m); err != nil {
		return err
	}

	return nil
}

func (ws *wsServer) ReadMessage() ([]byte, error) {
	_, m, err := ws.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (ws *wsServer) ReceiveEvent() (interface{}, error) {
	m, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	var reply data.WSEvent
	err = json.Unmarshal(m, &reply)
	if err != nil {
		return nil, err
	}

	for {
		switch reply.Type {
		case common.PushBlockEvents:
			var pushEvent []data.Event
			_ = json.Unmarshal(reply.Data, &pushEvent)
			return pushEvent, nil
		case common.RevertBlockEvents:
			var pushEvent *data.RevertBlock
			_ = json.Unmarshal(reply.Data, &pushEvent)
			return pushEvent, nil
		case common.FinalizedBlockEvents:
			var pushEvent *data.FinalizedBlock
			_ = json.Unmarshal(reply.Data, &pushEvent)
			return pushEvent, nil
		default:
			return nil, errors.New("invalid reply type")
		}
	}
}

func (ws *wsServer) ReceiveEvents() ([]data.Event, error) {
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

func (ws *wsServer) ReceiveRevertBlock() (*data.RevertBlock, error) {
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

func (ws *wsServer) ReceiveFinalized() (*data.FinalizedBlock, error) {
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

func (ws *wsServer) Close() {
	ws.httpServer.Close()
	ws.wsConn.Close()
}
