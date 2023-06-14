package gin_test

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/gin"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/require"
)

func createMockArgsWebServerHandler() gin.ArgsWebServerHandler {
	return gin.ArgsWebServerHandler{
		Facade: &mocks.FacadeStub{},
		Config: config.ConnectorApiConfig{
			Port: "8080",
		},
		Type:               "notifier",
		ConnectorType:      "http",
		Marshaller:         &mock.MarshalizerMock{},
		InternalMarshaller: &mock.MarshalizerMock{},
	}
}

func TestNewWebServerHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Facade = nil

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, apiErrors.ErrNilFacadeHandler, err)
	})

	t.Run("nil external marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Marshaller = nil

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("nil internal marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.InternalMarshaller = nil

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, common.ErrNilInternalMarshaller, err)
	})

	t.Run("invalid api type", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Type = ""

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, common.ErrInvalidAPIType, err)
	})

	t.Run("invalid obs connector type", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.ConnectorType = ""

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, common.ErrInvalidConnectorType, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()

		ws, err := gin.NewWebServerHandler(args)
		require.Nil(t, err)
		require.NotNil(t, ws)

		err = ws.Run()
		require.Nil(t, err)

		err = ws.Close()
		require.Nil(t, err)
	})
}
