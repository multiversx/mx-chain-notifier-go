package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const eventsPath = "/events"

func TestNewEventsGroup(t *testing.T) {
	t.Parallel()

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		eg, err := groups.NewEventsGroup(nil)
		require.True(t, errors.Is(err, apiErrors.ErrNilFacadeHandler))
		require.True(t, check.IfNil(eg))
	})

	t.Run("without basic auth middleware, should work", func(t *testing.T) {
		t.Parallel()

		eg, err := groups.NewEventsGroup(&mocks.FacadeStub{})
		require.NoError(t, err)
		require.NotNil(t, eg)

		require.Equal(t, 0, len(eg.GetAdditionalMiddlewares()))
	})

	t.Run("with basic auth middleware, should work", func(t *testing.T) {
		t.Parallel()

		facade := &mocks.FacadeStub{
			GetConnectorUserAndPassCalled: func() (string, string) {
				return "user", "pass"
			},
		}
		eg, err := groups.NewEventsGroup(facade)
		require.Nil(t, err)

		require.NotNil(t, eg.GetAuthMiddleware())
	})
}

func TestEventsGroup_PushEvents(t *testing.T) {
	t.Parallel()

	t.Run("invalid data, bad request", func(t *testing.T) {
		t.Parallel()

		eg, err := groups.NewEventsGroup(&mocks.FacadeStub{})
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/push", bytes.NewBuffer([]byte("invalid data")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("facade error, will try push events v2, should fail", func(t *testing.T) {
		t.Parallel()

		blockEvents := &data.SaveBlockData{
			Hash: "hash1",
			Txs: map[string]*transaction.Transaction{
				"hash2": {
					Nonce: 2,
				},
			},
			Scrs: map[string]*smartContractResult.SmartContractResult{
				"hash3": {
					Nonce: 3,
				},
			},
			LogEvents: []data.Event{},
		}

		jsonBytes, _ := json.Marshal(blockEvents)

		wasCalled := false
		facade := &mocks.FacadeStub{
			HandlePushEventsV1Called: func(events data.SaveBlockData) error {
				wasCalled = true
				return errors.New("facade error")
			},
		}

		eg, err := groups.NewEventsGroup(facade)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/push", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.True(t, wasCalled)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		argsSaveBlockData := data.ArgsSaveBlockData{
			HeaderHash: []byte{},
			Body:       &block.Body{},
			Header:     &block.HeaderV2{},
			TransactionsPool: &data.TransactionsPool{
				Txs: map[string]*data.NodeTransaction{
					"hash2": {
						TransactionHandler: &transaction.Transaction{
							Nonce: 2,
						},
						ExecutionOrder: 1,
					},
				},
				Scrs: map[string]*data.NodeSmartContractResult{
					"hash3": {
						TransactionHandler: &smartContractResult.SmartContractResult{
							Nonce: 3,
						},
						ExecutionOrder: 1,
					},
				},
				Logs: []*data.LogData{
					{
						LogHandler: &transaction.Log{
							Address: []byte("logaddr1"),
							Events:  []*transaction.Event{},
						},
						TxHash: "logHash1",
					},
				},
			},
		}
		blockEvents := data.ArgsSaveBlock{
			HeaderType:        "HeaderV2",
			ArgsSaveBlockData: argsSaveBlockData,
		}

		jsonBytes, _ := json.Marshal(blockEvents)

		wasCalled := false
		facade := &mocks.FacadeStub{
			HandlePushEventsV1Called: func(eventsData data.SaveBlockData) error {
				return common.ErrReceivedEmptyEvents
			},
			HandlePushEventsV2Called: func(events data.ArgsSaveBlockData) error {
				wasCalled = true
				assert.Equal(t, argsSaveBlockData, events)
				return nil
			},
		}

		eg, err := groups.NewEventsGroup(facade)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/push", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.True(t, wasCalled)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestEventsGroup_RevertEvents(t *testing.T) {
	t.Parallel()

	t.Run("invalid data, bad request", func(t *testing.T) {
		t.Parallel()

		eg, err := groups.NewEventsGroup(&mocks.FacadeStub{})
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/revert", bytes.NewBuffer([]byte("invalid data")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		revertBlockEvents := data.RevertBlock{
			Hash:  "hash1",
			Nonce: 1,
		}
		jsonBytes, _ := json.Marshal(revertBlockEvents)

		wasCalled := false
		facade := &mocks.FacadeStub{
			HandleRevertEventsCalled: func(events data.RevertBlock) {
				wasCalled = true
				assert.Equal(t, revertBlockEvents, events)
			},
		}

		eg, err := groups.NewEventsGroup(facade)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/revert", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.True(t, wasCalled)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestEventsGroup_FinalizedEvents(t *testing.T) {
	t.Parallel()

	t.Run("invalid data, bad request", func(t *testing.T) {
		t.Parallel()

		eg, err := groups.NewEventsGroup(&mocks.FacadeStub{})
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/finalized", bytes.NewBuffer([]byte("invalid data")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		finalizedBlockEvents := data.FinalizedBlock{
			Hash: "hash1",
		}
		jsonBytes, _ := json.Marshal(finalizedBlockEvents)

		wasCalled := false
		facade := &mocks.FacadeStub{
			HandleFinalizedEventsCalled: func(events data.FinalizedBlock) {
				wasCalled = true
				assert.Equal(t, finalizedBlockEvents, events)
			},
		}

		eg, err := groups.NewEventsGroup(facade)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath, getEventsRoutesConfig())

		req, _ := http.NewRequest("POST", "/events/finalized", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.True(t, wasCalled)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func getEventsRoutesConfig() config.APIRoutesConfig {
	return config.APIRoutesConfig{
		APIPackages: map[string]config.APIPackageConfig{
			"events": {
				Routes: []config.RouteConfig{
					{Name: "/push", Open: true},
					{Name: "/revert", Open: true},
					{Name: "/finalized", Open: true},
				},
			},
		},
	}
}
