package ws_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/ws"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgsWSHandler() ws.ArgsWebSocketHandler {
	return ws.ArgsWebSocketHandler{
		Hub:      &disabled.Hub{},
		Upgrader: &mocks.WSUpgraderStub{},
	}
}

func TestNewWebSocketHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil hub handler", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		args.Hub = nil

		wh, err := ws.NewWebSocketHandler(args)
		require.True(t, check.IfNil(wh))
		assert.Equal(t, ws.ErrNilHubHandler, err)
	})

	t.Run("nil ws upgrader", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		args.Upgrader = nil

		wh, err := ws.NewWebSocketHandler(args)
		require.True(t, check.IfNil(wh))
		assert.Equal(t, ws.ErrNilWSUpgrader, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWSHandler()
		wh, err := ws.NewWebSocketHandler(args)
		require.False(t, check.IfNil(wh))
		require.Nil(t, err)
	})
}

func TestServeHTTP(t *testing.T) {
	t.Parallel()

	args := createMockArgsWSHandler()

	wasCalled := false
	args.Upgrader = &mocks.WSUpgraderStub{
		UpgradeCalled: func(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (dispatcher.WSConnection, error) {
			wasCalled = true
			return &mocks.WSConnStub{}, nil
		},
	}

	wh, err := ws.NewWebSocketHandler(args)
	require.Nil(t, err)

	wh.ServeHTTP(httptest.NewRecorder(), &http.Request{})
	assert.True(t, wasCalled)
}
