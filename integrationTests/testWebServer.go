package integrationTests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/stretchr/testify/assert"
)

// TestWebServer defines a test web server instance
type TestWebServer struct {
	facade     shared.FacadeHandler
	apiType    string
	marshaller marshal.Marshalizer
	ws         *gin.Engine
	mutWs      sync.Mutex
}

// NewTestWebServer creates a new test web server
func NewTestWebServer(facade shared.FacadeHandler, apiType string) *TestWebServer {
	webServer := &TestWebServer{
		facade:     facade,
		apiType:    apiType,
		marshaller: &marshal.JsonMarshalizer{},
	}

	ws := gin.New()
	ws.Use(cors.Default())

	groupsMap := webServer.createGroups()
	for groupName, groupHandler := range groupsMap {
		ginGroup := ws.Group(groupName)
		groupHandler.RegisterRoutes(ginGroup)
	}

	webServer.ws = ws

	return webServer
}

// DoRequest preforms a test request on the web server, returning the response ready to be parsed
func (w *TestWebServer) DoRequest(request *http.Request) *httptest.ResponseRecorder {
	w.mutWs.Lock()
	defer w.mutWs.Unlock()

	resp := httptest.NewRecorder()
	w.ws.ServeHTTP(resp, request)

	return resp
}

func (w *TestWebServer) createGroups() map[string]shared.GroupHandler {
	groupsMap := make(map[string]shared.GroupHandler)

	eventsGroup, err := groups.NewEventsGroup(w.facade, w.marshaller)
	if err == nil {
		groupsMap["events"] = eventsGroup
	}

	if w.apiType == common.WSAPIType {
		hubHandler, err := groups.NewHubGroup(w.facade)
		if err == nil {
			groupsMap["hub"] = hubHandler
		}
	}

	return groupsMap
}

// PushEventsRequest will send a http request for push events
func (w *TestWebServer) PushEventsRequest(events *outport.OutportBlock) *httptest.ResponseRecorder {
	jsonBytes, _ := json.Marshal(events)

	req, _ := http.NewRequest("POST", "/events/push", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := w.DoRequest(req)
	return resp
}

// RevertEventsRequest will send a http request for revert event
func (w *TestWebServer) RevertEventsRequest(events *data.RevertBlock) *httptest.ResponseRecorder {
	jsonBytes, _ := json.Marshal(events)

	req, _ := http.NewRequest("POST", "/events/revert", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := w.DoRequest(req)
	return resp
}

// FinalizedEventsRequest will send a http request for finalized event
func (w *TestWebServer) FinalizedEventsRequest(events *data.FinalizedBlock) *httptest.ResponseRecorder {
	jsonBytes, _ := json.Marshal(events)

	req, _ := http.NewRequest("POST", "/events/finalized", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := w.DoRequest(req)
	return resp
}

func loadResponse(t *testing.T, rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)

	assert.Nil(t, err)
}
