package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const eventsPath = "/events"

func createMockEventsGroupArgs() groups.ArgsEventsGroup {
	return groups.ArgsEventsGroup{
		Facade:            &mocks.FacadeStub{},
		EventsDataHandler: &mocks.EventsDataHandlerStub{},
	}
}

func TestNewEventsGroup(t *testing.T) {
	t.Parallel()

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsGroupArgs()
		args.Facade = nil

		eg, err := groups.NewEventsGroup(args)
		require.True(t, errors.Is(err, apiErrors.ErrNilFacadeHandler))
		require.True(t, check.IfNil(eg))
	})

	t.Run("nil events data handler", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsGroupArgs()
		args.EventsDataHandler = nil

		eg, err := groups.NewEventsGroup(args)
		require.True(t, errors.Is(err, groups.ErrNilEventsDataHandler))
		require.True(t, check.IfNil(eg))
	})

	t.Run("without basic auth middleware, should work", func(t *testing.T) {
		t.Parallel()

		eg, err := groups.NewEventsGroup(createMockEventsGroupArgs())
		require.NoError(t, err)
		require.NotNil(t, eg)

		require.Equal(t, 0, len(eg.GetAdditionalMiddlewares()))
	})

	t.Run("with basic auth middleware, should work", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsGroupArgs()
		args.Facade = &mocks.FacadeStub{
			GetConnectorUserAndPassCalled: func() (string, string) {
				return "user", "pass"
			},
		}
		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		require.Equal(t, 1, len(eg.GetAdditionalMiddlewares()))

	})
}

func TestEventsGroup_PushEvents(t *testing.T) {
	t.Parallel()

	t.Run("facade error, will try push events v2, should fail", func(t *testing.T) {
		t.Parallel()

		argsSaveBlockData := outport.OutportBlock{
			BlockData: &outport.BlockData{
				HeaderType: string(core.ShardHeaderV2),
				HeaderHash: []byte("headerHash"),
			},
			TransactionPool: &outport.TransactionPool{
				Transactions: map[string]*outport.TxInfo{
					"hash2": {
						Transaction: &transaction.Transaction{
							Nonce: 2,
						},
						ExecutionOrder: 1,
					},
				},
			},
		}

		jsonBytes, _ := json.Marshal(argsSaveBlockData)

		wasCalled := false
		wasCalledV2 := false

		args := createMockEventsGroupArgs()
		args.Facade = &mocks.FacadeStub{
			HandlePushEventsV1Called: func(events data.SaveBlockData) error {
				wasCalled = true
				return errors.New("facade error")
			},
			HandlePushEventsV2Called: func(events data.ArgsSaveBlockData) error {
				wasCalledV2 = true
				return errors.New("facade error")
			},
		}

		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath)

		req, _ := http.NewRequest("POST", "/events/push", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.True(t, wasCalled)
		assert.True(t, wasCalledV2)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		header := &block.HeaderV2{
			Header: &block.Header{
				Nonce: 1,
			},
		}
		headerBytes, _ := json.Marshal(header)

		argsSaveBlockData := outport.OutportBlock{
			BlockData: &outport.BlockData{
				HeaderBytes: headerBytes,
				HeaderType:  string(core.ShardHeaderV2),
				HeaderHash:  []byte("headerHash"),
				Body: &block.Body{
					MiniBlocks: []*block.MiniBlock{
						{SenderShardID: 1},
					},
				},
			},
			TransactionPool: &outport.TransactionPool{
				Transactions: map[string]*outport.TxInfo{
					"hash2": {
						Transaction: &transaction.Transaction{
							Nonce: 2,
						},
						ExecutionOrder: 1,
					},
				},
				SmartContractResults: map[string]*outport.SCRInfo{
					"hash3": {
						SmartContractResult: &smartContractResult.SmartContractResult{
							Nonce: 3,
						},
						ExecutionOrder: 1,
					},
				},
				Logs: []*outport.LogData{
					{
						Log: &transaction.Log{
							Address: []byte("logaddr1"),
							Events:  []*transaction.Event{},
						},
						TxHash: "logHash1",
					},
				},
			},
			HeaderGasConsumption: &outport.HeaderGasConsumption{},
		}

		jsonBytes, _ := json.Marshal(argsSaveBlockData)

		wasCalled := false
		args := createMockEventsGroupArgs()
		args.Facade = &mocks.FacadeStub{
			HandlePushEventsV1Called: func(eventsData data.SaveBlockData) error {
				return common.ErrReceivedEmptyEvents
			},
			HandlePushEventsV2Called: func(events data.ArgsSaveBlockData) error {
				wasCalled = true
				return nil
			},
		}

		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath)

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

		args := createMockEventsGroupArgs()
		args.EventsDataHandler = &mocks.EventsDataHandlerStub{
			UnmarshallRevertDataCalled: func(marshalledData []byte) (*data.RevertBlock, error) {
				return nil, errors.New("expected err")
			},
		}

		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath)

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
		args := createMockEventsGroupArgs()
		args.EventsDataHandler = &mocks.EventsDataHandlerStub{
			UnmarshallRevertDataCalled: func(marshalledData []byte) (*data.RevertBlock, error) {
				return &revertBlockEvents, nil
			},
		}
		args.Facade = &mocks.FacadeStub{
			HandleRevertEventsCalled: func(events data.RevertBlock) {
				wasCalled = true
				assert.Equal(t, revertBlockEvents, events)
			},
		}

		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath)

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

		args := createMockEventsGroupArgs()
		args.EventsDataHandler = &mocks.EventsDataHandlerStub{
			UnmarshallFinalizedDataCalled: func(marshalledData []byte) (*data.FinalizedBlock, error) {
				return nil, errors.New("expected err")
			},
		}

		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath)

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
		args := createMockEventsGroupArgs()
		args.EventsDataHandler = &mocks.EventsDataHandlerStub{
			UnmarshallFinalizedDataCalled: func(marshalledData []byte) (*data.FinalizedBlock, error) {
				return &finalizedBlockEvents, nil
			},
		}
		args.Facade = &mocks.FacadeStub{
			HandleFinalizedEventsCalled: func(events data.FinalizedBlock) {
				wasCalled = true
				assert.Equal(t, finalizedBlockEvents, events)
			},
		}

		eg, err := groups.NewEventsGroup(args)
		require.Nil(t, err)

		ws := startWebServer(eg, eventsPath)

		req, _ := http.NewRequest("POST", "/events/finalized", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		assert.True(t, wasCalled)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
