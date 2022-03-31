package gin_test

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	apiErrors "github.com/ElrondNetwork/notifier-go/api/errors"
	"github.com/ElrondNetwork/notifier-go/api/gin"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPServer(t *testing.T) {
	t.Parallel()

	t.Run("nil http server", func(t *testing.T) {
		t.Parallel()

		server, err := gin.NewHTTPServerWrapper(nil)
		require.True(t, check.IfNil(server))
		require.Equal(t, apiErrors.ErrNilHTTPServer, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		listenAndServerWasCalled := false
		shutdownWasCalled := false
		httpServer := &mocks.HTTPServerStub{
			ListenAndServeCalled: func() error {
				listenAndServerWasCalled = true
				return nil
			},
			ShutdownCalled: func(ctx context.Context) error {
				shutdownWasCalled = true
				return nil
			},
		}

		server, err := gin.NewHTTPServerWrapper(httpServer)
		require.Nil(t, err)

		server.Start()
		require.Nil(t, err)
		assert.True(t, listenAndServerWasCalled)

		err = server.Close()
		require.Nil(t, err)
		assert.True(t, shutdownWasCalled)
	})
}
